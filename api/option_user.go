package api

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/makuo12/ghost_server/constants"
	db "github.com/makuo12/ghost_server/db/sqlc"
	"github.com/makuo12/ghost_server/tools"
	"github.com/makuo12/ghost_server/val"

	"github.com/gin-gonic/gin"
)

func (server *Server) ListExperience(ctx *gin.Context) {
	var req ExperienceOffsetParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at  ListExperience in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionOffset)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Println(req)
	var resOption ListExperienceOptionRes
	var resEvent ListExperienceEventRes
	var hasData bool
	var err error
	switch req.MainOptionType {
	case "options":
		//resOption, err, hasData = HandleListOptionExperience(ctx, server, req)
		resOption, err, hasData = HandleRedisOptionExperience(ctx, server, req)

	case "events":
		//resEvent, err, hasData = HandleListEventExperience(ctx, server, req)
		resEvent, err, hasData = HandleRedisEventExperience(ctx, server, req)
	}
	if hasData && err == nil {
		switch req.MainOptionType {
		case "options":
			ctx.JSON(http.StatusOK, resOption)
			return
		case "events":
			ctx.JSON(http.StatusOK, resEvent)
			return
		}

	} else if !hasData && err == nil {
		res := ExperienceCategoryRes{
			Category: req.Type,
		}
		ctx.JSON(http.StatusNoContent, res)
		return
	} else {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

}

func (server *Server) CreateReportOptionUser(ctx *gin.Context) {
	var req CreateReportOptionParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("error at CreateReportOptionUser in ShouldBindJSON: %v \n", err)
		err = fmt.Errorf("there was an error while processing your inputs please try again later")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	optionUserID, err := tools.StringToUuid(req.OptionUserID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := HandleGetUser(ctx, server)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if req.TypeOne == "scam" {
		if !val.ValidateReportScamOption(req.TypeTwo) {
			log.Println("Report scam option not correct")
			err = fmt.Errorf("this scam option is not recorded on our app")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
	} else if req.TypeOne == "offensive" {
		log.Println("Report offensive option not correct")
		if !val.ValidateReportOffensiveOption(req.TypeTwo) {
			err = fmt.Errorf("this offensive option is not recorded on our app")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
	}
	err = server.store.CreateReportOption(ctx, db.CreateReportOptionParams{
		OptionUserID: optionUserID,
		UserID:       user.UserID,
		TypeOne:      tools.HandleStringTwo(req.TypeOne),
		TypeTwo:      tools.HandleStringTwo(req.TypeTwo),
		TypeThree:    tools.HandleStringTwo(req.TypeThree),
	})
	if err != nil {
		log.Printf("Error at CreateReportOptionUser in CreateReportOption err: %v, user: %v photoID: %v\n", err, user.ID, req.OptionUserID)
		err = fmt.Errorf("could not account your report, please try again or try connecting us")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	res := UserResponseMsg{
		Success: true,
	}
	log.Printf("CreateReportOptionUser successfully (%v) id: %v\n", user.Email, user.ID)
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) GetExperienceDetail(ctx *gin.Context) {
	var req ExperienceDetailParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at  ListExperience in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionUserID)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	none := constants.NONE
	log.Println(req)
	var resOption ExperienceOptionDetailRes
	var resEvent ExperienceEventDetailRes
	var hasData bool
	var err error
	switch req.MainOptionType {
	case "options":
		resOption, hasData, err = HandleDetailOptionExperience(ctx, server, req)

	case "events":
		resEvent, hasData, err = HandleDetailEventExperience(ctx, server, req)
	}
	if hasData && err == nil {
		switch req.MainOptionType {
		case "options":
			ctx.JSON(http.StatusOK, resOption)
			return
		case "events":
			ctx.JSON(http.StatusOK, resEvent)
			return
		}

	} else if !hasData && err == nil {

		ctx.JSON(http.StatusNoContent, none)
		return
	} else {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

}

func (server *Server) ListExperienceEventTickets(ctx *gin.Context) {
	var req ListExperienceEventTicketsParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at  ListExperienceEventTickets in ShouldBindJSON: %v, eventDateTimeID: %v \n", err.Error(), req.EventDateTimeID)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	dollarToNaira := server.config.DollarToNaira
	dollarToCAD := server.config.DollarToCAD
	none := constants.NONE
	eventDateTimeID, err := tools.StringToUuid(req.EventDateTimeID)
	if err != nil {
		log.Printf("Error at ListExperienceEventTickets in StringToUuid err: %v, eventDateTimeID: %v\n", err, req.EventDateTimeID)
		err = fmt.Errorf("this event date is not found")
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}
	hostCurrency, err := server.store.GetEventDateOptionInfo(ctx, eventDateTimeID)
	if err != nil {
		log.Printf("Error at ListExperienceEventTickets in StringToUuid err: %v, eventDateTimeID: %v\n", err, req.EventDateTimeID)
		err = fmt.Errorf("this event date is not found")
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}
	startDate, err := tools.ConvertDateOnlyStringToDate(req.StartDate)
	if err != nil {
		log.Printf("Error at ListExperienceEventTickets in StringToUuid err: %v, eventDateTimeID: %v\n", err, req.EventDateTimeID)
		err = fmt.Errorf("this event date is not found")
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}
	var res ListExperienceEventTicketsRes
	tickets, err := server.store.ListEventDateTicketUser(ctx, eventDateTimeID)
	if err != nil || len(tickets) == 0 {
		if err != nil {
			log.Printf("Error at ListExperienceEventTickets in ListEventDateTicketUser err: %v, eventDateTimeID: %v\n", err, req.EventDateTimeID)
		}
		list := []ExEventTicketData{{none, none, none, false, none, none, none, none, 0, false, true}}
		res = ListExperienceEventTicketsRes{
			List:            list,
			EventDateTimeID: req.EventDateTimeID,
			IsEmpty:         true,
		}
	} else {
		var resData []ExEventTicketData
		for _, t := range tickets {
			// We check if the ticket is available
			count, err := server.store.CountChargeTicketReferenceByStartDate(ctx, db.CountChargeTicketReferenceByStartDateParams{
				Date:        startDate,
				EventDateID: eventDateTimeID,
				Cancelled:   false,
				IsComplete:  true,
				TicketID:    t.ID,
			})
			if err != nil {
				log.Printf("Error at FuncName %v SetupExperienceEventData in CountChargeTicketReferenceByStartDate err: %v id: %v ticketID: %v\n", "ListExperienceEventTickets", err, eventDateTimeID, t.ID)
				err = nil
			} else {
				if int32(count) >= t.Capacity {
					continue
				}
			}

			price, errPrice := tools.ConvertPrice(tools.IntToMoneyString(t.Price), hostCurrency, req.Currency, dollarToNaira, dollarToCAD, t.ID)
			if errPrice != nil {
				log.Printf("Error at ListExperienceEventTickets in tools.ConvertPrice err: %v, eventDateTimeID: %v\n", err, req.EventDateTimeID)
				price = 0.0
			}
			data := ExEventTicketData{
				ID:               tools.UuidToString(t.ID),
				Name:             t.Name,
				Price:            tools.ConvertFloatToString(price),
				AbsorbFees:       t.AbsorbFees,
				Description:      t.Description,
				Type:             t.Type,
				Level:            t.Level,
				TicketType:       t.TicketType,
				NumOfSeats:       int(t.NumOfSeats),
				FreeRefreshments: t.FreeRefreshment,
				IsEmpty:          false,
			}
			resData = append(resData, data)
		}
		if len(resData) == 0 {
			resData = []ExEventTicketData{{none, none, none, false, none, none, none, none, 0, false, true}}
		}
		res = ListExperienceEventTicketsRes{
			List:            resData,
			EventDateTimeID: req.EventDateTimeID,
			IsEmpty:         false,
		}
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) ListExOptionDateTime(ctx *gin.Context) {
	var req ExperienceDetailParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at  ListExOptionDateTime in ShouldBindJSON: %v, OptionUserID: %v \n", err.Error(), req.OptionUserID)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	optionUserID, err := tools.StringToUuid(req.OptionUserID)
	if err != nil {
		log.Printf("Error at  ListExOptionDateTime in tools.StringToUuid: %v, OptionUserID: %v \n", err.Error(), req.OptionUserID)
		err = errors.New("this stay does not exist")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	option, err := server.store.GetOptionInfoByOptionWithPriceUserID(ctx, optionUserID)
	if err != nil {
		log.Printf("Error at ListExOptionDateTime in GetOptionInfoByOptionUserID err: %v, user: %v, optionID: %v\n", err, ctx.ClientIP(), option.ID)
		err = errors.New("this stay does not exist")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	available, err := server.store.GetOptionAvailabilitySetting(ctx, option.ID)
	if err != nil {
		log.Printf("Error at ListExOptionDateTime in GetOptionAvailabilitySetting err: %v, user: %v, optionID: %v\n", err, ctx.ClientIP(), option.ID)
		err = errors.New("this stay available settings not included something is wrong, try again")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var dates []ExOptionDateTimeItem
	// Lets first handle prepare time
	dates = append(dates, HandleExAvailable(ctx, server, option.OptionUserID, available.PreparationTime, available.AutoBlockDates)...)
	// Next we handle dates set by the host
	dates = append(dates, HandleHostSpecialDates(ctx, server, option, req.Currency)...)

	busyDates, priceDates, busyIsEmpty, priceIsEmpty := HandleImportantDates(dates, tools.IntToMoneyString(option.Price))

	res := ListExOptionDateTimeRes{
		BusyDates:    busyDates,
		PriceDates:   priceDates,
		BusyIsEmpty:  busyIsEmpty,
		PriceIsEmpty: priceIsEmpty,
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) GetExperienceAmDetail(ctx *gin.Context) {
	var req ExperienceDetailAmParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at  GetExperienceAmDetail in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionUserID)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	optionUserID, err := tools.StringToUuid(req.OptionUserID)
	if err != nil {
		log.Printf("Error at  GetExperienceAmDetail in tools.StringToUuid: %v, optionID: %v \n", err.Error(), req.OptionUserID)
		err = errors.New("this listing does not exist")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	am, err := server.store.GetAmenityDetailByOptionUserID(ctx, db.GetAmenityDetailByOptionUserIDParams{
		OptionUserID: optionUserID,
		Tag:          req.Tag,
	})
	if err != nil {
		log.Printf("Error at  GetExperienceAmDetail in GetAmenityDetailByOptionUserID: %v, optionID: %v \n", err.Error(), req.OptionUserID)
		err = errors.New("amenity not found")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	res := ExOptionAmDetail{
		TimeSet:            am.TimeSet,
		LocationOption:     am.LocationOption,
		SizeOption:         int(am.SizeOption),
		PrivacyOption:      am.PrivacyOption,
		TimeOption:         am.TimeOption,
		StartTime:          tools.ConvertTimeOnlyToString(am.StartTime),
		EndTime:            tools.ConvertTimeOnlyToString(am.EndTime),
		AvailabilityOption: am.AvailabilityOption,
		StartMonth:         am.StartMonth,
		EndMonth:           am.EndMonth,
		TypeOption:         am.TypeOption,
		PriceOption:        am.PriceOption,
		BrandOption:        am.BrandOption,
		ListOptions:        am.ListOptions,
	}
	ctx.JSON(http.StatusOK, res)
}
