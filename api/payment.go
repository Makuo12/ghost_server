package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/makuo12/ghost_server/constants"
	db "github.com/makuo12/ghost_server/db/sqlc"
	"github.com/makuo12/ghost_server/payment"
)

func (server *Server) VerifyPaymentReference(ctx *gin.Context) {
	var req payment.ReferencePayment

	// If chargeData.MainObjectType is either options or events then we want willRefund to be true because when it gets to chargeData.MainObjectType if we are making payment it turns to false
	var willRefund bool = true
	var successReservation bool = false
	var cardData payment.CardAddResponse = payment.GetFakeCardRes()
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
	chargeData, err := server.store.GetChargeReference(ctx, db.GetChargeReferenceParams{
		UserID:    user.UserID,
		Reference: req.Reference,
	})
	if err != nil {
		log.Printf("error at VerifyPaymentReference at GetChargeReference for userID: %v, err: %v\n", user.UserID, err.Error())
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	resData, err := payment.HandlePaystackVerifyPayment(ctx, server.config.PaystackSecretLiveKey, chargeData.PaymentReference, user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	
	if !chargeData.IsComplete {
		chargeData, err = server.store.UpdateChargeReferenceComplete(ctx, db.UpdateChargeReferenceCompleteParams{
			IsComplete: true,
			UserID:     user.UserID,
			Reference:  req.Reference,
		})
	}
	switch chargeData.MainObjectType {
	case "options":
		willRefund = false
		success, errData := ObjectOptionPaymentReference(ctx, server, user, chargeData.Reference, chargeData.ObjectReference, int(chargeData.Charge), chargeData.Currency, req.Message)
		err = errData
		successReservation = success
	case "events":
		willRefund = false
		success, errData := ObjectEventPaymentReference(ctx, server, user, chargeData.Reference, chargeData.ObjectReference, int(chargeData.Charge), chargeData.Currency, req.Message)
		err = errData
		successReservation = success
	}
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if req.AddCard {
		cardData, err = CreateCard(ctx, server, resData, user, willRefund, req.Type, chargeData.ID, chargeData.Currency)
		// If willRefund is true that means we are not making any payment
		if err != nil && willRefund {
			log.Printf("error at VerifyPaymentReference at CreateCard for userID: %v, err: %v\n", user.UserID, err.Error())
			err = fmt.Errorf("card was not successfully created")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
	}
	res := payment.ReferencePaymentResponse{
		Verified: successReservation,
		Card:     cardData,
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) InitMethodPayment(ctx *gin.Context) {
	var req InitMethodPaymentParams
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
	if req.PaymentType == constants.PAYMENT_METHOD_FOR_RESERVATION {
		paystackBankCharge, paystackPWT, paystackUSSD, paystackCard, detailRes, hasReqData, _, paymentReference, err := ReservePaymentMethod(ctx, server, req, user)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		if hasReqData {
			ctx.JSON(http.StatusOK, detailRes)
			return
		}
		res := InitMethodPaymentRes{
			Reference:          req.Reference,
			PaymentType:        req.PaymentType,
			MainOptionType:     req.MainOptionType,
			PaymentReference:   paymentReference,
			PaymentMethod:      req.PaymentMethod,
			PaymentChannel:     req.PaymentChannel,
			PaystackBankCharge: paystackBankCharge,
			PaystackPWT:        paystackPWT,
			PaystackCard:       paystackCard,
			PaystackUSSD:       paystackUSSD,
		}
		fmt.Println("paystack: ", res)
		fmt.Println("paystackBankCharge: ", paystackBankCharge)
		fmt.Println("PaystackPWT: ", paystackPWT)
		fmt.Println("PaystackCard: ", paystackCard)
		fmt.Println("PaystackUSSD: ", paystackUSSD)
		ctx.JSON(http.StatusOK, res)
		return
	}
	ctx.JSON(http.StatusBadRequest, errorResponse(err))
}
