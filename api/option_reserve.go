package api

import (
	"errors"
	db "flex_server/db/sqlc"
	"flex_server/tools"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
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
	resData, resChallengeData, resChallenged, err := HandlePaystackChargeAuthorization(server, ctx, card, totalFee)
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
			// Creating snapshot
			err := HandleOptionReserveComplete(server, ctx, reserveData, refRes, resData.Data.Reference, user, req.Message, fromCharge)
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
	detailRes, hasResData, _, refRes, reserveData, fromCharge, _, err := HandleFinalOptionReserveDetail(server, ctx, req.Reference, user, "no card at 2 factor verification", req.Message)
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
	log.Printf("FinalOptionReserveVerificationDetail sent successfully (%v) id: %v\n", user.Email, user.ID)
	ctx.JSON(http.StatusOK, res)
}
