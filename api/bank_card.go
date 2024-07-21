package api

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/makuo12/ghost_server/constants"
	db "github.com/makuo12/ghost_server/db/sqlc"
	"github.com/makuo12/ghost_server/payment"
	"github.com/makuo12/ghost_server/tools"
)

func (server *Server) SetDefaultCard(ctx *gin.Context) {
	var req payment.SetDefaultCardParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at SetDefaultCard in ShouldBindJSON: %v, currency: %v \n", err.Error(), req.ID)
		err = fmt.Errorf("the params do not meet the requirement")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := HandleGetUser(ctx, server)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	requestID, err := tools.StringToUuid(req.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	card, err := server.store.GetCard(ctx, db.GetCardParams{
		ID:     requestID,
		UserID: user.ID,
	})
	if err != nil {
		log.Printf("Error at SetDefaultCard in GetCard for Reference: %v and user: %v \n", requestID, user.ID)
		err = fmt.Errorf("this card doesn't exist")
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	switch req.Type {
	case "payout":
		_, err = server.store.UpdateUser(ctx, db.UpdateUserParams{
			DefaultPayoutCard: pgtype.Text{
				String: tools.UuidToString(card.ID),
				Valid:  true,
			},
			ID: user.ID,
		})
		if err != nil {
			log.Printf("Error at SetDefaultCard Payout in UpdateUse for Reference: %v and user: %v \n", requestID, user.ID)
			err = fmt.Errorf("could not set this card as your default card")
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
	case "payment":
		_, err = server.store.UpdateUser(ctx, db.UpdateUserParams{
			DefaultCard: pgtype.Text{
				String: tools.UuidToString(card.ID),
				Valid:  true,
			},
			ID: user.ID,
		})
		if err != nil {
			log.Printf("Error at SetDefaultCard payment in UpdateUse for Reference: %v and user: %v \n", requestID, user.ID)
			err = fmt.Errorf("could not set this card as your default card")
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
	}
	res := payment.SetDefaultCardRes{
		Success: true,
		ID:      req.ID,
		Type:    req.Type,
	}
	log.Printf("InitCardCharge sent successfully (%v) id: %v\n", user.Email, user.ID)
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) InitCardCharge(ctx *gin.Context) {
	var charge string
	var hasObjectReference bool = true
	var reason = constants.PAY_CARD_REASON
	var objectReference uuid.UUID
	var req payment.InitCardChargeParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at InitCardCharge in ShouldBindJSON: %v, currency: %v \n", err.Error(), req.Currency)
		err = fmt.Errorf("select one of the provided options")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := HandleGetUser(ctx, server)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	_, err = server.store.GetCardByLast4(ctx, db.GetCardByLast4Params{
		Last4:    req.CardLast4,
		Currency: req.Currency,
		UserID:   user.ID,
	})
	if err == nil {
		// If error is nil we expect that the card already exist
		err = fmt.Errorf("this already exist")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	switch req.Currency {
	case "USD":
		charge = server.config.AddCardChargeDollar
	case "NGN":
		charge = server.config.AddCardChargeNaira
	default:
		err = fmt.Errorf("currency is invalid")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	objectReference, err = tools.StringToUuid(req.ObjectReference)
	if err != nil {
		hasObjectReference = false
		objectReference = uuid.New()
		reason = constants.ADD_CARD_REASON
		err = nil
	}
	reference := tools.UuidToString(uuid.New())
	paymentReference := tools.UuidToString(uuid.New())
	_, err = CreateChargeReference(ctx, server, user.UserID, reference, paymentReference, objectReference, hasObjectReference, reason, req.Currency, "none", charge, "InitCardCharge")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	paystackCharge := tools.ConvertToPaystackCharge(charge)
	res := payment.InitCardChargeRes{
		Reference: paymentReference,
		Reason:    constants.ADD_CARD_REASON,
		Charge:    paystackCharge,
		Currency:  req.Currency,
		Email:     user.Email,
	}
	log.Printf("InitCardCharge sent successfully (%v) id: %v\n", user.Email, user.ID)
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) InitRemoveCard(ctx *gin.Context) {
	var req payment.InitRemoveCardParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at InitRemoveCard in ShouldBindJSON: %v, ID: %v \n", err.Error(), req.ID)
		err = fmt.Errorf("select one of the provided options")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	cardID, err := tools.StringToUuid(req.ID)
	if err != nil {
		log.Printf("Error at InitRemoveCard in StringToUuid: %v, req.ID: %v \n", err.Error(), req.ID)
		err = fmt.Errorf("card not found")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := HandleGetUser(ctx, server)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err = server.store.RemoveCard(ctx, db.RemoveCardParams{
		UserID: user.ID,
		ID:     cardID,
	})
	if err != nil {
		log.Printf("Error at InitRemoveCard in RemoveCard: %v, userID: %v \n", err.Error(), user.ID)
		err = fmt.Errorf("sorry an error ocurred, please try again later")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	res := payment.InitRemoveCardRes{
		Success: true,
		ID:      req.ID,
	}
	log.Printf("InitRemoveCard sent successfully (%v) id: %v\n", user.Email, user.ID)
	ctx.JSON(http.StatusOK, res)
}

func CreateCard(ctx context.Context, server *Server, resData payment.PaystackVerifyResponse, user db.User, willRefund bool, cardStoreType string, chargeID uuid.UUID, currency string) (res payment.CardAddResponse, err error) {
	accountName := fmt.Sprintf("%v", resData.Data.Authorization.AccountName)
	argsCard := db.CreateCardParams{
		UserID:            user.ID,
		Email:             resData.Data.Customer.Email,
		AuthorizationCode: resData.Data.Authorization.AuthorizationCode,
		CardType:          resData.Data.Authorization.CardType,
		Last4:             resData.Data.Authorization.Last4,
		ExpMonth:          resData.Data.Authorization.ExpMonth,
		ExpYear:           resData.Data.Authorization.ExpYear,
		Bank:              resData.Data.Authorization.Bank,
		CountryCode:       resData.Data.Authorization.CountryCode,
		Reusable:          resData.Data.Authorization.Reusable,
		Channel:           resData.Data.Authorization.Channel,
		AccountName:       accountName,
		Bin:               resData.Data.Authorization.Bin,
		CardSignature:     resData.Data.Authorization.Signature,
		Currency:          resData.Data.Currency,
	}
	card, err := server.store.CreateCard(ctx, argsCard)
	if err != nil {
		log.Printf("Error at VerifyPaymentReference in CreateCard. Error is %v for AuthorizationCode: %v and user: %v \n", err, resData.Data.Authorization.AuthorizationCode, user.ID)
		err = fmt.Errorf("there was an error while verifying your card. Please don't try just again we are working on it")
		return
	}
	// update user default card
	switch cardStoreType {
	case "payout":
		_, err = server.store.UpdateUser(ctx, db.UpdateUserParams{
			DefaultPayoutCard: pgtype.Text{
				String: tools.UuidToString(card.ID),
				Valid:  true,
			},
			ID: user.ID,
		})
		if err != nil {
			log.Printf("Error at VerifyPaymentReference payout in UpdateUser. Error is %v for DefaultID: %v and user: %v \n", err, card.ID, user.ID)
		}

	case "payment":
		_, err = server.store.UpdateUser(ctx, db.UpdateUserParams{
			DefaultCard: pgtype.Text{
				String: tools.UuidToString(card.ID),
				Valid:  true,
			},
			ID: user.ID,
		})
		if err != nil {
			log.Printf("Error at VerifyPaymentReference payment in UpdateUser. Error is %v for DefaultID: %v and user: %v \n", err, card.ID, user.ID)
		}

	}

	// We create a refund to send all the money back
	_, err = server.store.CreateMainRefund(ctx, db.CreateMainRefundParams{
		ChargeID:    chargeID,
		UserPercent: 100,
		HostPercent: 0,
		ChargeType:  constants.CHARGE_REFERENCE,
		Type:        constants.ADD_CARD,
	})

	if err != nil {
		log.Printf("Error at VerifyPaymentReference in CreateMainRefund. Error is %v for AuthorizationCode: %v and user: %v \n", err, resData.Data.Authorization.AuthorizationCode, user.ID)
	}
	// If card is not reusable then we send an error
	if !resData.Data.Authorization.Reusable {
		err = fmt.Errorf("there was an error while verifying your card. Please note that this card is not reusable so it cannot be used as a payment method. You would be sent your charge refund shortly")
		return
	}

	resCard := payment.CardDetailResponse{
		CardID:    tools.UuidToString(card.ID),
		CardLast4: card.Last4,
		CardType:  card.CardType,
		ExpMonth:  card.ExpMonth,
		ExpYear:   card.ExpYear,
		Currency:  card.Currency,
	}

	res = payment.CardAddResponse{
		CardDetail:      resCard,
		Account:         "0.00",
		AccountCurrency: currency,
		DefaultID:       user.DefaultCard,
	}
	return
}
