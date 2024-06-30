package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/makuo12/ghost_server/payment"
	"github.com/makuo12/ghost_server/tools"
	"github.com/makuo12/ghost_server/utils"
)

func (server *Server) GetWallet(ctx *gin.Context) {
	var usdAccount string
	var ngnAccount string
	var hasCard bool
	var cardDetail []payment.CardDetailResponse
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
		cardDetail = []payment.CardDetailResponse{{"none", "none", "none", "none", "none", "none"}}
	} else {
		hasCard = true
		for _, c := range cards {
			data := payment.CardDetailResponse{
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
	res := payment.GetWalletResponse{
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
