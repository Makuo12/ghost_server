package api

import (
	"log"

	db "github.com/makuo12/ghost_server/db/sqlc"
	"github.com/makuo12/ghost_server/tools"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func UpdateEventCurrency(currencyToConvert string, ctx *gin.Context, server *Server, option db.OptionsInfo, oldCurrency string, eventDateTimeID uuid.UUID, userID uuid.UUID) {
	dollarToNaira := server.config.DollarToNaira
	dollarToCAD := server.config.DollarToCAD
	tickets, err := server.store.ListEventDateTicket(ctx, eventDateTimeID)
	log.Println("ticketsTR", tickets)
	if err != nil {
		log.Printf("There an error at UpdateEventCurrency at ListEventDateTicket: %v, optionID: %v, eventDateTimeID: %v, userID: %v \n", err.Error(), option.ID, eventDateTimeID, userID)
		return
	}
	for _, t := range tickets {

		price, err := tools.ConvertPrice(tools.IntToMoneyString(t.Price), oldCurrency, currencyToConvert, dollarToNaira, dollarToCAD, userID)
		if err != nil {
			log.Printf("There an error at UpdateEventCurrency at tools.ConvertPrice: %v, optionID: %v, eventDateTimeID: %v, userID: %v \n", err.Error(), option.ID, eventDateTimeID, userID)
			continue
		}
		_, err = server.store.UpdateEventDateTicket(ctx, db.UpdateEventDateTicketParams{
			Price: pgtype.Int8{
				Int64: tools.MoneyFloatToInt(price),
				Valid: true,
			},
			ID:              t.ID,
			EventDateTimeID: t.EventDateTimeID,
		})
		if err != nil {
			log.Printf("There an error at UpdateEventCurrency at UpdateEventDateTicket: %v, optionID: %v, eventDateTimeID: %v, userID: %v \n", err.Error(), option.ID, eventDateTimeID, userID)
			continue
		}
	}

}

func UpdateOptionDataCurrency(currencyToConvert string, ctx *gin.Context, server *Server, option db.OptionsInfo, oldCurrency string, userID uuid.UUID) {
	dollarToNaira := server.config.DollarToNaira
	dollarToCAD := server.config.DollarToCAD
	p, err := server.store.GetOptionPrice(ctx, option.ID)
	if err != nil {
		log.Printf("There an error at UpdateOptionDataCurrency at GetOptionPrice: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, userID)
		return
	}
	price, err := tools.ConvertPrice(tools.IntToMoneyString(p.Price), oldCurrency, currencyToConvert, dollarToNaira, dollarToCAD, userID)
	if err != nil {
		log.Printf("There an error at UpdateOptionDataCurrency at tools.ConvertPrice: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, userID)
		return

	}
	weekendPrice, err := tools.ConvertPrice(tools.IntToMoneyString(p.WeekendPrice), oldCurrency, currencyToConvert, dollarToNaira, dollarToCAD, userID)
	if err != nil {
		log.Printf("There an error at UpdateOptionDataCurrency at tools.ConvertPrice: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, userID)
		return

	}
	// We also try to do weekend price
	_, err = server.store.UpdateOptionPrice(ctx, db.UpdateOptionPriceParams{
		Price: pgtype.Int8{
			Int64: tools.MoneyFloatToInt(price),
			Valid: true,
		},
		WeekendPrice: pgtype.Int8{
			Int64: tools.MoneyFloatToInt(weekendPrice),
			Valid: true,
		},
		OptionID: option.ID,
	})
	if err != nil {
		log.Printf("There an error at UpdateOptionDataCurrency at UpdateOptionPrice: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, userID)
		return
	}

}

func UpdateOptionAddChargeCurrency(currencyToConvert string, ctx *gin.Context, server *Server, option db.OptionsInfo, oldCurrency string, userID uuid.UUID) {
	dollarToNaira := server.config.DollarToNaira
	dollarToCAD := server.config.DollarToCAD
	charges, err := server.store.ListOptionAddCharge(ctx, option.ID)
	if err != nil {
		log.Printf("There an error at UpdateOptionAddChargeCurrency at ListOptionAddCharge: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, userID)
		return
	}
	for _, char := range charges {
		mainFee, err := tools.ConvertPrice(tools.IntToMoneyString(char.MainFee), oldCurrency, currencyToConvert, dollarToNaira, dollarToCAD, userID)
		if err != nil {
			log.Printf("There an error at UpdateOptionAddChargeCurrency at tools.NairaToDollar: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, userID)
			continue
		}
		extraFee, err := tools.ConvertPrice(tools.IntToMoneyString(char.ExtraFee), oldCurrency, currencyToConvert, dollarToNaira, dollarToCAD, userID)
		if err != nil {
			log.Printf("There an error at UpdateOptionAddChargeCurrency extraFee at tools.NairaToDollar: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, userID)
			continue
		}
		_, err = server.store.UpdateOptionAddCharge(ctx, db.UpdateOptionAddChargeParams{
			MainFee: pgtype.Int8{
				Int64: tools.MoneyFloatToInt(mainFee),
				Valid: true,
			},

			ExtraFee: pgtype.Int8{
				Int64: tools.MoneyFloatToInt(extraFee),
				Valid: true,
			},
			OptionID: option.ID,
			ID:       char.ID,
		})
		if err != nil {
			log.Printf("There an error at UpdateOptionAddChargeCurrency at UpdateOptionAddCharge: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, userID)
			continue
		}
	}
}
