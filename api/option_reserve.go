package api

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/makuo12/ghost_server/constants"
	db "github.com/makuo12/ghost_server/db/sqlc"
	"github.com/makuo12/ghost_server/payment"
	"github.com/makuo12/ghost_server/tools"
)

func (server *Server) CreateOptionReserveDetail(ctx *gin.Context) {
	var req ReserveOptionParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("error at CreateOptionReserveDetail in ShouldBindJSON: %v \n", err)
		err = fmt.Errorf("there was an error while processing your inputs please try again later")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := HandleGetUser(ctx, server)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	optionUserID, err := tools.StringToUuid(req.OptionUserID)
	if err != nil {
		log.Printf("error at CreateOptionReserveDetail at optionUserID, err := tools.StringToUuid: %v, userID: %v\n", err.Error(), user.ID)
		err = errors.New("this list does not exist")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	option, err := server.store.GetOptionInfoByOptionUserID(ctx, optionUserID)
	if err != nil {
		log.Printf("error at CreateOptionReserveDetail at store.GetOptionInfoByOptionUserID: %v, userID: %v\n", err.Error(), user.ID)
		err = errors.New("this list does not exist")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if option.HostID == user.ID {
		err = errors.New("you cannot book your own listing")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	reference, err := HandleReserveOption(user, server, ctx, req.StartDate, req.EndDate, req.Guests, req.OptionUserID, req.UserCurrency)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	userID := tools.UuidToString(user.ID)
	// Card Details
	defaultCardID, cardDetail, hasCard := HandleReserveCard(ctx, server, user, "CreateOptionReserveDetail")
	reserveData, err := HandleOptionReserveRedisData(userID, reference)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	res := CreateOptionReserveDetailRes{
		ReserveData:   reserveData,
		DefaultCardID: defaultCardID,
		HasCard:       hasCard,
		CardDetail:    cardDetail,
	}
	log.Printf("CreateOptionReserveDetail successfully (%v) id: %v\n", user.Email, user.ID)
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) FinalOptionReserveDetail(ctx *gin.Context) {
	var req FinalOptionReserveDetailParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at FinalOptionReserveDetail in ShouldBindJSON: %v, reference: %v \n", err.Error(), req.Reference)
		err = fmt.Errorf("select one of the provided options")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	successUrl := server.config.PaymentSuccessUrl
	failureUrl := server.config.PaymentFailUrl
	cardID, err := tools.StringToUuid(req.ID)
	if err != nil {
		log.Printf("Error at FinalOptionReserveDetail in StringToUuid: %v, reqID: %v \n", err.Error(), req.ID)
		err = fmt.Errorf("this payment option does not exist")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := HandleGetUser(ctx, server)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	card, err := server.store.GetCard(ctx, db.GetCardParams{
		ID:     cardID,
		UserID: user.ID,
	})
	if err != nil {
		log.Printf("Error at FinalOptionReserveDetail in GetCard: %v, cardID: %v, userID: %v \n", err.Error(), cardID, user.ID)
		err = fmt.Errorf("this payment option does not exist")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	detailRes, hasResData, totalFee, refRes, reserveData, fromCharge, chargeData, err := HandleFinalOptionReserveDetail(server, ctx, req.Reference, user, tools.UuidToString(card.ID), req.Message)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if hasResData {
		// This means we are meant to send a reservation request response
		ctx.JSON(http.StatusOK, detailRes)
		return
	}
	if reserveData.Currency != card.Currency && !fromCharge {
		err = fmt.Errorf("payment did not go through, please selected currency doesn't go with the card")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if chargeData.Currency != card.Currency && fromCharge {
		err = fmt.Errorf("payment did not go through, please selected currency doesn't go with the card")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if fromCharge {
		// We want to check if the dates are still available
		option, err := server.store.GetOptionInfoCustomer(ctx, db.GetOptionInfoCustomerParams{
			OptionUserID:    chargeData.OptionUserID,
			IsComplete:      true,
			IsActive:        true, // Option is active
			IsActive_2:      true, // Host is active
			OptionStatusOne: "list",
			OptionStatusTwo: "staged",
		})
		if err != nil {
			log.Printf("Error at FinalOptionReserveDetail in store.GetOptionInfoCustomer: %v, cardID: %v, userID: %v \n", err.Error(), cardID, user.ID)
			err = fmt.Errorf("your selected dates are no more available, please try again")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		available, err := HandleChargeDatesAvailable(ctx, server, option, tools.ConvertDateOnlyToString(chargeData.StartDate), tools.ConvertDateOnlyToString(chargeData.EndDate), "FinalOptionReserveDetail")
		if err != nil || !available {
			if err != nil {
				log.Printf("Error at FinalOptionReserveDetail in HandleChargeDatesAvailable: %v, cardID: %v, userID: %v \n", err.Error(), cardID, user.ID)
			}
			err = fmt.Errorf("your selected dates are no more available, please try again")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
	}
	resData, resChallengeData, resChallenged, err := payment.HandlePaystackChargeAuthorization(ctx, successUrl, failureUrl, server.config.PaystackSecretLiveKey, card, totalFee)
	if err != nil {
		log.Printf("Error at FinalOptionReserveDetail in HandlePaystackChargeAuthorization: %v, cardID: %v, userID: %v \n", err.Error(), cardID, user.ID)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var res FinalOptionReserveDetailRes
	// If the transaction was not challenge we expect status to be success
	if !resChallenged {
		log.Println("pay_data ", resData.Data)
		if resData.Data.Status != "success" {
			log.Printf("Error at HandlePaystackChargeAuthorization payment did not go through")
			err = fmt.Errorf(resData.Data.GatewayResponse)
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		} else {
			// Payment was successful
			// We want to save a recept in the database
			// We also want to store a snap shot of what the shortlet looks like
			res = FinalOptionReserveDetailRes{
				Reference:         refRes,
				AuthorizationUrl:  "none",
				AccessCode:        "none",
				PaymentReference:  resData.Data.Reference,
				Paused:            false,
				PaymentSuccess:    true,
				PaymentChallenged: false,
				SuccessUrl:        successUrl,
				FailureUrl:        failureUrl,
			}
			// We want to send information back to the user
			host, err := server.store.GetOptionInfoUserIDByUserID(ctx, chargeData.OptionUserID)
			if err != nil {
				log.Printf("Error at HandleOptionReserveRequest in GetOptionInfoByUserID: %v optionID: %v referenceID: %v, pay_method_reference: %v\n", err.Error(), chargeData.OptionUserID, chargeData.Reference, chargeData.PaymentReference)
				err = nil
			} else {
				header := "Payment and Reservation confirmed"
				msg := "Thank you for using Flizzup"
				checkIn := tools.HandleReadableDate(chargeData.StartDate, tools.DateDMMYyyy)
				checkout := tools.HandleReadableDate(chargeData.EndDate, tools.DateDMMYyyy)
				BrevoOptionPaymentSuccess(ctx, server, header, msg, "FinalOptionReserveDetail", chargeData.ID, host.Email, host.FirstName, host.LastName, tools.UuidToString(chargeData.ID), tools.UuidToString(host.UserID), user.Email, user.FirstName, user.LastName, tools.UuidToString(user.UserID), host.HostNameOption, checkIn, checkout)
				// Notification for Guest
				msg = "Payment received, and reservation confirmed! Check your email for further information about your booking, including scheduling an inspection and our 100% refund policy if the property does not match the app's description."
				header = fmt.Sprintf("Hey %v, booking confirmed", user.FirstName)
				CreateTypeNotification(ctx, server, chargeData.ID, user.UserID, constants.OPTION_PAYMENT_SUCCESSFUL, msg, false, header)
				// Notification for Host
				msg = fmt.Sprintf("You have a new booking at %v! Check your email for further details.", host.HostNameOption)
				header = fmt.Sprintf("Hey %v", host.FirstName)
				CreateTypeNotification(ctx, server, chargeData.ID, host.UserID, constants.HOST_OPTION_PAYMENT_SUCCESSFUL, msg, false, header)
			}
			// Creating snapshot
			err = HandleOptionReserveComplete(server, ctx, reserveData, refRes, resData.Data.Reference, user, req.Message, fromCharge)
			if err != nil {
				log.Printf("Error at FinalOptionReserveDetail in HandleOptionReserveStore: %v, cardID: %v, userID: %v \n", err.Error(), cardID, user.ID)
				err = fmt.Errorf("your payment was successful, but we were unable to create a receipt pls try contacting us because you would need this to verify your self and post vids")
				ctx.JSON(http.StatusBadRequest, errorResponse(err))
				return
			}
			// We want to send no error to the user because the payment was successful

		}
	} else {
		// This means payment was unsuccessful because the payment was challenged
		// We want to make sure sure the challenged data is not empty
		if resChallengeData.Status {
			// This means that is not empty
			res = FinalOptionReserveDetailRes{
				Reference:         reserveData.Reference,
				AuthorizationUrl:  resChallengeData.Data.AuthorizationUrl,
				AccessCode:        resChallengeData.Data.AccessCode,
				PaymentReference:  resChallengeData.Data.Reference,
				Paused:            resChallengeData.Data.Paused,
				PaymentSuccess:    false,
				PaymentChallenged: true,
				SuccessUrl:        successUrl,
				FailureUrl:        failureUrl,
			}
		} else {
			// Because no challenged data we want to send an error
			log.Printf("Error at HandlePaystackChargeAuthorization payment did not go through")
			err = fmt.Errorf(resData.Data.GatewayResponse)
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
	}
	log.Printf("FinalOptionReserveDetail sent successfully (%v) id: %v\n", user.Email, user.ID)
	ctx.JSON(http.StatusOK, res)
}

// This is after the two factor verification
func (server *Server) FinalOptionReserveVerificationDetail(ctx *gin.Context) {
	var req FinalOptionReserveDetailVerificationParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at FinalOptionReserveVerificationDetail in ShouldBindJSON: %v, reference: %v \n", err.Error(), req.Reference)
		err = fmt.Errorf("select one of the provided options")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := HandleGetUser(ctx, server)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if !req.Successful {
		//// TimerRemoveOptionReserveUser we call this function to remove the user
		//TimerRemoveOptionReserveUser(tools.UuidToString(user.ID), req.Reference)()
		err = fmt.Errorf("payment was unsuccessful, please contact us if your having any issues with paying")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	detailRes, hasResData, _, refRes, reserveData, fromCharge, chargeData, err := HandleFinalOptionReserveDetail(server, ctx, req.Reference, user, "no card at 2 factor verification", req.Message)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if hasResData {
		// This means we are meant to send a reservation request response
		ctx.JSON(http.StatusOK, detailRes)
		return
	}
	// Creating snapshot
	err = HandleOptionReserveComplete(server, ctx, reserveData, refRes, req.PaymentReference, user, req.Message, fromCharge)
	if err != nil {
		log.Printf("Error at FinalOptionReserveDetail in HandleOptionReserveStore: %v, userID: %v \n", err.Error(), user.ID)
		err = fmt.Errorf("your payment was successful, but we were unable to create a receipt pls try contacting us because you would need this to verify your self and post vids")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}
	res := FinalOptionReserveDetailRes{
		Reference:         reserveData.Reference,
		AuthorizationUrl:  "none",
		AccessCode:        "none",
		PaymentReference:  req.PaymentReference,
		Paused:            false,
		PaymentSuccess:    true,
		PaymentChallenged: false,
	}
	// We want to send information back to the user
	if fromCharge {
		host, err := server.store.GetOptionInfoUserIDByUserID(ctx, chargeData.OptionUserID)
		if err != nil {
			log.Printf("Error at HandleOptionReserveRequest in GetOptionInfoByUserID: %v optionID: %v referenceID: %v, pay_method_reference: %v\n", err.Error(), chargeData.OptionUserID, chargeData.Reference, chargeData.PaymentReference)
			err = nil
		} else {
			header := "Payment and Reservation confirmed"
			msg := "Thank you for using Flizzup"
			checkIn := tools.HandleReadableDate(chargeData.StartDate, tools.DateDMMYyyy)
			checkout := tools.HandleReadableDate(chargeData.EndDate, tools.DateDMMYyyy)
			BrevoOptionPaymentSuccess(ctx, server, header, msg, "FinalOptionReserveDetail", chargeData.ID, host.Email, host.FirstName, host.LastName, tools.UuidToString(chargeData.ID), tools.UuidToString(host.UserID), user.Email, user.FirstName, user.LastName, tools.UuidToString(user.UserID), host.HostNameOption, checkIn, checkout)
			// Notification for Guest
			msg = "Payment received, and reservation confirmed! Check your email for further information about your booking, including scheduling an inspection and our 100% refund policy if the property does not match the app's description."
			header = fmt.Sprintf("Hey %v, booking confirmed", user.FirstName)
			CreateTypeNotification(ctx, server, chargeData.ID, user.UserID, constants.OPTION_PAYMENT_SUCCESSFUL, msg, false, header)
			// Notification for Host
			msg = fmt.Sprintf("You have a new booking at %v! Check your email for further details.", host.HostNameOption)
			header = fmt.Sprintf("Hey %v", host.FirstName)
			CreateTypeNotification(ctx, server, chargeData.ID, host.UserID, constants.HOST_OPTION_PAYMENT_SUCCESSFUL, msg, false, header)
		}
	}
	log.Printf("FinalOptionReserveVerificationDetail sent successfully (%v) id: %v\n", user.Email, user.ID)
	ctx.JSON(http.StatusOK, res)
}

func ReservePaymentMethod(ctx context.Context, server *Server, arg InitMethodPaymentParams, user db.User) (paystackBankCharge payment.PaystackBankAccountMainRes, paystackPWT payment.PaystackPWTMainRes, paystackUSSD payment.PaystackUSSDRes, paystackCard payment.InitCardChargeRes, detailRes FinalOptionReserveRequestDetailRes, hasReqData bool, reason string, err error) {
	var objectReference = uuid.New()
	var hasObjectReference = false
	var currency string
	var fee string
	var reference string
	switch arg.MainOptionType {
	case "options":
		detailDataRes, hasResData, totalFee, resRef, resData, fromCharge, chargeData, errData := HandleFinalOptionReserveDetail(server, ctx, arg.Reference, user, arg.Reference, arg.Message)
		if errData != nil {
			err = errData
			return
		}
		if hasResData {
			detailRes = detailDataRes
			// This means we are meant to send a reservation request response
			hasReqData = hasResData
			return
		}
		if fromCharge {
			objectReference = chargeData.ID
			hasObjectReference = true
		}
		reason = constants.USER_OPTION_PAYMENT
		currency = resData.Currency
		fee = totalFee
		reference = resRef
	case "events":
		resData, errData := HandleEventReserveRedisData(tools.UuidToString(user.ID), arg.Reference)
		if errData != nil {
			err = errData
			return
		}
		reason = constants.USER_EVENT_PAYMENT
		currency = resData.Currency
		fee = resData.TotalFee
		reference = arg.Reference
	}
	if err != nil {
		return
	}
	_, err = CreateChargeReference(ctx, server, user.UserID, reference, objectReference, hasObjectReference, reason, currency, arg.MainOptionType, fee, "ReservePaymentMethod")
	if err != nil {
		return
	}
	paystackBankCharge, paystackPWT, paystackUSSD, paystackCard, err = ReservePaymentChannel(ctx, server, arg, user, fee, reference, reason)
	return
}

func ReservePaymentChannel(ctx context.Context, server *Server, arg InitMethodPaymentParams, user db.User, charge string, reference string, reason string) (paystackBankCharge payment.PaystackBankAccountMainRes, paystackPWT payment.PaystackPWTMainRes, paystackUSSD payment.PaystackUSSDRes, paystackCard payment.InitCardChargeRes, err error) {
	switch arg.PaymentChannel {
	case constants.PAYSTACK_BANK_ACCOUNT:
		res, errData := payment.HandlePaystackBankAccount(ctx, server.config.PaystackSecretLiveKey, charge, reference, arg.PaystackBankAccount, user.Email)
		if errData != nil {
			err = errData
		} else {
			paystackBankCharge = payment.PaystackBankAccountMainRes{
				Reference:   res.Data.Reference,
				DisplayText: res.Data.DisplayText,
			}
		}
	case constants.PAYSTACK_PWT:
		res, errData := payment.HandlePaystackPWT(ctx, server.config.PaystackSecretLiveKey, charge, reference, user.Email)
		if errData != nil {
			err = errData
		} else {
			paystackPWT = payment.PaystackPWTMainRes{
				Reference:     res.Data.Reference,
				Slug:          res.Data.Bank.Slug,
				AccountName:   res.Data.AccountName,
				AccountNumber: res.Data.AccountNumber,
				ExpiresAt:     res.Data.AccountExpiresAt,
			}
		}
	case constants.PAYSTACK_CARD:
		res, errData := payment.HandlePaystackCard(ctx, server.config.PaystackSecretLiveKey, charge, reference, arg.Currency, user.Email, reason)
		if errData != nil {
			err = errData
		} else {
			paystackCard = res
		}
	case constants.PAYSTACK_USSD:
		res, errData := payment.HandlePaystackUSSD(ctx, server.config.PaystackSecretLiveKey, charge, reference, arg.PaystackUSSD, user.Email, user.FirstName)
		if errData != nil {
			err = errData
		} else {
			paystackUSSD = payment.PaystackUSSDRes{
				Reference:   res.Data.Reference,
				DisplayText: res.Data.DisplayText,
				USSDCode:    res.Data.UssdCode,
			}
		}
	}
	return
}
