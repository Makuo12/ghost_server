package api

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/makuo12/ghost_server/constants"
	db "github.com/makuo12/ghost_server/db/sqlc"
	"github.com/makuo12/ghost_server/tools"
	"github.com/makuo12/ghost_server/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func (server *Server) GetWallet(ctx *gin.Context) {
	var usdAccount string
	var ngnAccount string
	var hasCard bool
	var cardDetail []CardDetailResponse
	user, err := HandleGetUser(ctx, server)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	accounts, err := server.store.ListAccount(ctx, user.ID)
	if err != nil {
		log.Printf("There was an error at GetWallet in GetAccount %v \n", err.Error())
		err = fmt.Errorf("could not get your wallet")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	for _, a := range accounts {
		if a.Currency == utils.USD {
			usdAccount = tools.IntToMoneyString(a.Balance)
		} else if a.Currency == utils.NGN {
			ngnAccount = tools.IntToMoneyString(a.Balance)
		}
	}
	cards, err := server.store.ListCard(ctx, user.ID)
	if err != nil || len(cards) == 0 {
		if err != nil {
			log.Printf("There was an error at GetWallet in ListCard %v \n", err.Error())
		}
		hasCard = false
		cardDetail = []CardDetailResponse{{"none", "none", "none", "none", "none", "none"}}
	} else {
		hasCard = true
		for _, c := range cards {
			data := CardDetailResponse{
				CardID:    tools.UuidToString(c.ID),
				CardType:  c.CardType,
				ExpMonth:  c.ExpMonth,
				ExpYear:   c.ExpYear,
				Currency:  c.Currency,
				CardLast4: c.Last4,
			}
			cardDetail = append(cardDetail, data)
		}
	}
	res := GetWalletResponse{
		Cards:           cardDetail,
		USDAccount:      usdAccount,
		NGNAccount:      ngnAccount,
		DefaultID:       user.DefaultCard,
		DefaultPayoutID: user.DefaultPayoutCard,
		HasCard:         hasCard,
	}
	log.Println("GetWallet: ", res.Cards)
	log.Println("user wallet as been sent successfully to ", user.ID)
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) SetDefaultCard(ctx *gin.Context) {
	var req SetDefaultCardParams
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
	res := SetDefaultCardRes{
		Success: true,
		ID:      req.ID,
		Type:    req.Type,
	}
	log.Printf("InitAddCard sent successfully (%v) id: %v\n", user.Email, user.ID)
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) InitAddCard(ctx *gin.Context) {
	var req InitAddCardParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at InitAddCard in ShouldBindJSON: %v, currency: %v \n", err.Error(), req.Currency)
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

	var charge string
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
	// We set it in redis instead
	paymentReference := tools.UuidToString(uuid.New())
	chargeData := fmt.Sprintf("%v&%v&%v&%v", user.UserID, req.Currency, charge, constants.ADD_CARD_REASON)
	err = RedisClient.Set(RedisContext, paymentReference, chargeData, time.Hour*30).Err()
	if err != nil {
		log.Printf("Error at InitAddCard in RedisClient.Set: %v, userID: %v \n", err.Error(), user.ID)
		err = fmt.Errorf("could not initialize a reference, card was not added")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	paystackCharge := tools.ConvertToPaystackCharge(charge)
	res := InitAddCardRes{
		Reference: paymentReference,
		Reason:    constants.ADD_CARD_REASON,
		Charge:    paystackCharge,
		Currency:  req.Currency,
		Email:     user.Email,
	}
	log.Printf("InitAddCard sent successfully (%v) id: %v\n", user.Email, user.ID)
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) InitRemoveCard(ctx *gin.Context) {
	var req InitRemoveCardParams
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
	res := InitRemoveCardRes{
		Success: true,
		ID:      req.ID,
	}
	log.Printf("InitRemoveCard sent successfully (%v) id: %v\n", user.Email, user.ID)
	ctx.JSON(http.StatusOK, res)
}
