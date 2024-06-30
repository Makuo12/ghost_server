package api

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	db "github.com/makuo12/ghost_server/db/sqlc"
	"github.com/makuo12/ghost_server/payment"
	"github.com/makuo12/ghost_server/tools"

	"github.com/gin-gonic/gin"
)

func (server *Server) CreateEventReserveDetail(ctx *gin.Context) {

	var req ReserveEventParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("error at CreateEventReserveDetail in ShouldBindJSON: %v \n", err)
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
		err = errors.New("you cannot book your own event for now")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	reference, err := HandleReserveEvent(user, server, ctx, req.Tickets, req.OptionUserID, req.UserCurrency)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	userID := tools.UuidToString(user.ID)
	// Card Details
	defaultCardID, cardDetail, hasCard := HandleReserveCard(ctx, server, user, "CreateEventReserveDetail")

	reserveData, err := HandleEventReserveRedisData(userID, reference)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	res := CreateEventReserveDetailRes{
		ReserveData:    reserveData,
		DefaultCardID:  defaultCardID,
		HasCard:        hasCard,
		CardDetail:     cardDetail,
		EventReference: reference,
	}
	log.Printf("CreateEventReserveDetail successfully (%v) id: %v\n", user.Email, user.ID)
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) FinalEventReserveDetail(ctx *gin.Context) {
	var req FinalOptionReserveDetailParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at FinalEventReserveDetail in ShouldBindJSON: %v, currency: %v \n", err.Error(), req.Reference)
		err = fmt.Errorf("select one of the provided options")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	successUrl := server.config.PaymentSuccessUrl
	failureUrl := server.config.PaymentFailUrl
	cardID, err := tools.StringToUuid(req.ID)
	if err != nil {
		log.Printf("Error at FinalEventReserveDetail in StringToUuid: %v, reqID: %v \n", err.Error(), req.ID)
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
		log.Printf("Error at FinalEventReserveDetail in GetCard: %v, cardID: %v, userID: %v \n", err.Error(), cardID, user.ID)
		err = fmt.Errorf("this payment option does not exist")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	uID := tools.UuidToString(user.ID)
	reserveData, err := HandleEventReserveRedisData(uID, req.Reference)
	if err != nil {
		log.Printf("Error at FinalEventReserveDetail in HandleOptionReserveRedisData: %v, cardID: %v, userID: %v \n", err.Error(), cardID, user.ID)
		err = fmt.Errorf("reference has expired, please try again")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if reserveData.Currency != card.Currency {
		err = fmt.Errorf("selected currency doesn't match card currency")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	resData, resChallengeData, resChallenged, err := payment.HandlePaystackChargeAuthorization(ctx, successUrl, failureUrl, server.config.PaystackSecretLiveKey, card, reserveData.TotalFee)
	if err != nil {
		log.Printf("Error at FinalEventReserveDetail in HandlePaystackChargeAuthorization: %v, cardID: %v, userID: %v \n", err.Error(), cardID, user.ID)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var res FinalOptionReserveDetailRes
	// If the transaction was not challenge we expect status to be success
	if !resChallenged {
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
				Reference:         req.Reference,
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
			err = HandleEventReserveComplete(server, ctx, reserveData, resData.Data.Reference, user, req.Message, req.Reference)
			if err != nil {
				log.Printf("Error at FinalEventReserveDetail in HandleEventReserveComplete: %v, cardID: %v, userID: %v \n", err.Error(), cardID, user.ID)
				err = fmt.Errorf("your payment was successful, but we were unable to create a receipt pls try contacting us because you would need this to verify your self and post vids")
				ctx.JSON(http.StatusBadRequest, errorResponse(err))
				return
			}
			// We want to send no error to the user because the payment was successful
			err = nil
		}
	} else {
		// This means payment was unsuccessful because the payment was challenged
		// We want to make sure sure the challenged data is not empty
		if resChallengeData.Status {
			// This means that is not empty
			res = FinalOptionReserveDetailRes{
				Reference:         req.Reference,
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
	log.Printf("FinalEventReserveDetail sent successfully (%v) id: %v\n", user.Email, user.ID)
	ctx.JSON(http.StatusOK, res)
}

// This is after the two factor verification
func (server *Server) FinalEventReserveVerificationDetail(ctx *gin.Context) {
	var req FinalOptionReserveDetailVerificationParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at FinalEventReserveVerificationDetail in ShouldBindJSON: %v, currency: %v \n", err.Error(), req.Reference)
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
	userID := tools.UuidToString(user.ID)
	reserveData, err := HandleEventReserveRedisData(userID, req.Reference)
	if err != nil {
		log.Printf("Error at FinalEventReserveVerificationDetail in HandleOptionReserveRedisData: %v, req.PaymentReference: %v, userID: %v \n", err.Error(), req.PaymentReference, user.ID)
		err = fmt.Errorf("reference has expired, please try again")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	// Creating snapshot
	err = HandleEventReserveComplete(server, ctx, reserveData, req.Message, user, req.Message, req.Reference)
	if err != nil {
		log.Printf("Error at FinalEventReserveDetail in HandleEventReserveComplete: %v, paymentReference: %v, userID: %v \n", err.Error(), req.PaymentReference, user.ID)
		err = fmt.Errorf("your payment was successful, but we were unable to create a receipt pls try contacting us because you would need this to verify your self and post vids")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	res := FinalOptionReserveDetailRes{
		Reference:         req.Reference,
		AuthorizationUrl:  "none",
		AccessCode:        "none",
		PaymentReference:  req.PaymentReference,
		Paused:            false,
		PaymentSuccess:    true,
		PaymentChallenged: false,
	}
	log.Printf("FinalEventReserveVerificationDetail sent successfully (%v) id: %v\n", user.Email, user.ID)
	ctx.JSON(http.StatusOK, res)
}
