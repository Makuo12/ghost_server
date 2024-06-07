package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/makuo12/ghost_server/constants"
	db "github.com/makuo12/ghost_server/db/sqlc"
	"github.com/makuo12/ghost_server/tools"

	"github.com/gin-gonic/gin"
)

func (server *Server) InitPayment(ctx *gin.Context) {
	var req InitPaymentParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at InitAddCard in ShouldBindJSON: %v, Reference: %v \n", err.Error(), req.Reference)
		err = fmt.Errorf("select one of the provided options")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := HandleGetUser(ctx, server)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var res InitPaymentRes
	switch req.MainOptionType {
	case "options":
		detailRes, hasResData, totalFee, resRef, _, _, _, err := HandleFinalOptionReserveDetail(server, ctx, req.Reference, user, req.Reference, req.Message)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		if hasResData {
			// This means we are meant to send a reservation request response
			ctx.JSON(http.StatusAccepted, detailRes)
			return
		}
		res = InitPaymentRes{
			Reference: resRef,
			Charge:    tools.ConvertToPaystackCharge(totalFee),
			Reason:    constants.USER_OPTION_PAYMENT,
			Email:     user.Email,
		}
	case "events":
		res, err = HandleInitPaymentEvent(user, req.Reference)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
	}
	log.Printf("InitAddCard sent successfully (%v) id: %v\n", user.Email, user.ID)
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) VerifyPaymentReference(ctx *gin.Context) {
	var req VerifyPaymentReferenceParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at VerifyPaymentReference in ShouldBindJSON: %v, currency: %v \n", err.Error(), req.Reference)
		err = fmt.Errorf("reference invalid")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := HandleGetUser(ctx, server)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	// HandlePaystackReference we try to verify payment here
	resData, err := HandlePaystackReference(server, ctx, req.Reference, user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	switch req.MainOptionType {
	case "options":
		var fromCharge bool
		var reserveOptionData ExperienceReserveOModel
		var referenceCharge string
		charge, err := server.store.GetChargeOptionReferenceByRef(ctx, db.GetChargeOptionReferenceByRefParams{
			Reference: req.Reference,
			UserID:    user.UserID,
		})
		if err != nil {
			log.Printf("Error at FinalOptionReserveDetail in GetChargeOptionReferenceByRef: %v, reference: %v, userID: %v \n", err.Error(), req.Reference, user.ID)
		} else {
			fromCharge = true
			referenceCharge = charge.Reference
			if resData.Data.Amount != int(charge.TotalFee) {
				// If amount is not equal we send back a refund and an error
				HandleInitRefund(ctx, server, user, reserveOptionData.Reference, constants.USER_OPTION_INVALID_PAYMENT_AMOUNT, resData.Data.Amount, resData.Data.Currency, "VerifyPaymentReference")
				err = fmt.Errorf("amount payed is not the total amount for the listing")
				ctx.JSON(http.StatusBadRequest, errorResponse(err))
				return
			}
		}
		log.Println("fromCharge, ", fromCharge)
		if !fromCharge {
			reserveOptionData, err = HandleOptionReserveRedisData(tools.UuidToString(user.ID), req.Reference)
			if err != nil {
				log.Printf("Error at FinalOptionReserveDetail in HandleOptionReserveRedisData: %v, reference: %v, userID: %v \n", err.Error(), req.Reference, user.ID)
				err = fmt.Errorf("req.Reference has expired, please try again")
				return
			}
			if resData.Data.Amount != tools.ConvertToPaystackCharge(reserveOptionData.TotalFee) {
				// If amount is not equal we send back a refund and an error
				HandleInitRefund(ctx, server, user, reserveOptionData.Reference, constants.USER_OPTION_INVALID_PAYMENT_AMOUNT, resData.Data.Amount, resData.Data.Currency, "VerifyPaymentReference")
				err = fmt.Errorf("amount payed is not the total amount for the listing")
				ctx.JSON(http.StatusBadRequest, errorResponse(err))
				return
			}
		}

		// We Store the receipt and snap shot
		err = HandleOptionReserveComplete(server, ctx, reserveOptionData, referenceCharge, resData.Data.Reference, user, req.Message, fromCharge)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}

	case "events":
		reserveEventData, err := HandleEventReserveRedisData(tools.UuidToString(user.ID), req.Reference)
		if err != nil {
			log.Printf("Error at HandleInitPaymentEvent in HandleOptionReserveRedisData: %v, req.Reference: %v, userID: %v \n", err.Error(), req.Reference, user.ID)
			err = fmt.Errorf("req.Reference has expired, please try again")
			return
		}
		if resData.Data.Amount != tools.ConvertToPaystackCharge(reserveEventData.TotalFee) {
			// If amount is not equal we send back a refund and an error
			HandleInitRefund(ctx, server, user, req.Reference, constants.USER_EVENT_INVALID_PAYMENT_AMOUNT, resData.Data.Amount, resData.Data.Currency, "VerifyPaymentReference")
			err = fmt.Errorf("amount payed is not the total amount for the event")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		// Creating snapshot
		err = HandleEventReserveComplete(server, ctx, reserveEventData, resData.Data.Reference, user, req.Message, req.Reference)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
	}
	res := VerifyPaymentReferenceRes{
		Reference:      req.Reference,
		MainOptionType: req.MainOptionType,
		Success:        true,
	}

	ctx.JSON(http.StatusOK, res)
}
