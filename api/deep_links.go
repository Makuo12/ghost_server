package api

import (
	"fmt"
	"log"
	"net/http"

	db "github.com/makuo12/ghost_server/db/sqlc"
	"github.com/makuo12/ghost_server/tools"

	"github.com/gin-gonic/gin"
)

func (server *Server) GetOptionDeepLink(ctx *gin.Context) {
	var req GetOptionParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at GetUHMOptionData in GetOptionParams in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("a title is required")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Println(req)
	requestID, err := tools.StringToUuid(req.OptionID)
	if err != nil {
		log.Printf("Error at tools.StringToUuid: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	_, _, option, _, _, err := HandleGetCompleteOptionEditOptionInfo(requestID, ctx, server, true)
	if err != nil {
		err = fmt.Errorf("you cannot access this resource")
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}
	deepLink := fmt.Sprintf("https://flizzup.com/flex?type=%v&id=%v", option.MainOptionType, option.DeepLinkID)
	res := GetDeepLinkRes{
		DeepLink: deepLink,
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) GetEventDateDeepLink(ctx *gin.Context) {
	var req GetEventDateDeepLink
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at GetEventDateDeepLink in ShouldBindJSON: %v, EventInfoID: %v \n", err.Error(), req.EventDateTimeID)
		err = fmt.Errorf("error occurred while processing your request please try again")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Println(req)
	eventID, err := tools.StringToUuid(req.EventInfoID)
	if err != nil {
		log.Printf("Error at tools.StringToUuid: %v, EventInfoID: %v \n", err.Error(), req.EventInfoID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	eventDateTimeID, err := tools.StringToUuid(req.EventDateTimeID)
	if err != nil {
		log.Printf("Error at tools.StringToUuid: %v, EventDateTimeID: %v \n", err.Error(), req.EventDateTimeID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, _, option, _, _, err := HandleGetCompleteOptionEditEventDateTimes(eventID, ctx, server, true)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	eventDate, err := server.store.GetEventDateTime(ctx, eventDateTimeID)
	if err != nil {
		log.Printf("There an error at GetEventDateDeepLink at GetEventDateTimeByOption: %v, EventInfoID: %v, userID: %v \n", err.Error(), eventDateTimeID, user.ID)
		err = fmt.Errorf("could not preform this update please try again")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	deepLink := fmt.Sprintf("https://flizzup.com/flex?type=%v&oid=%v&id=%v", "event_date", eventDate.DeepLinkID, option.DeepLinkID)
	res := GetDeepLinkRes{
		DeepLink: deepLink,
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) GetOptionDeepLinkExperience(ctx *gin.Context) {
	var req GetDeepLinkExperienceParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at GetDeepLinkExperienceParams in GetOptionParams in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.DeepLinkID)
		err = fmt.Errorf("a title is required")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Println(req)
	linkID, err := tools.StringToUuid(req.DeepLinkID)
	if err != nil {
		log.Printf("Error at tools.StringToUuid: %v, DeepLinkID: %v \n", err.Error(), req.DeepLinkID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	data, err := server.store.GetOptionExperienceByDeepLinkID(ctx, db.GetOptionExperienceByDeepLinkIDParams{
		DeepLinkID: linkID,
		IsActive:   true,
		IsComplete: true,
		IsActive_2: true,
	})
	if err != nil {
		log.Printf("Error at GetDeepLinkExperience in .GetOptionExperienceByDeepLinkID err: %v, user: %v\n", err, ctx.ClientIP())
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	basePrice, err := tools.ConvertPrice(tools.IntToMoneyString(data.Price), data.Currency, req.Currency, server.config.DollarToNaira, server.config.DollarToCAD, linkID)
	if err != nil {
		log.Printf("Error at  basePrice GetDeepLinkExperience in ConvertPrice err: %v, user: %v\n", err, ctx.ClientIP())
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return

	}
	weekendPrice, err := tools.ConvertPrice(tools.IntToMoneyString(data.WeekendPrice), data.Currency, req.Currency, server.config.DollarToNaira, server.config.DollarToCAD, linkID)
	if err != nil {
		log.Printf("Error at weekendPrice GetDeepLinkExperience in ConvertPrice err: %v, user: %v\n", err, ctx.ClientIP())
		weekendPrice = 0.0
	}
	addDateFound, startDateBook, endDateBook, addPrice := HandleOptionRedisExAddPrice(ctx, server, data.ID, data.OptionUserID, data.PreparationTime, data.AvailabilityWindow, data.AdvanceNotice, data.Price, data.WeekendPrice)
	addedPrice, err := tools.ConvertPrice(tools.IntToMoneyString(addPrice), data.Currency, req.Currency, server.config.DollarToNaira, server.config.DollarToCAD, linkID)
	if err != nil {
		log.Printf("Error at addedPrice GetDeepLinkExperience in ConvertPrice err: %v, user: %v\n", err, ctx.ClientIP())
		addedPrice = 0.0
	}
	res := ExperienceOptionData{
		UserOptionID:     tools.UuidToString(data.OptionUserID),
		Name:             data.HostNameOption,
		IsVerified:       data.IsVerified,
		CoverImage:       data.CoverImage,
		HostAsIndividual: data.HostAsIndividual,
		BasePrice:        tools.ConvertFloatToString(basePrice),
		WeekendPrice:     tools.ConvertFloatToString(weekendPrice),
		Photos:           data.Photo,
		TypeOfShortlet:   data.TypeOfShortlet,
		State:            data.State,
		Country:          data.Country,
		ProfilePhoto:     data.Photo_2,
		HostName:         data.FirstName,
		HostJoined:       tools.ConvertDateOnlyToString(data.CreatedAt),
		HostVerified:     data.IsVerified_2,
		Category:         data.Category,
		AddedPrice:       tools.ConvertFloatToString(addedPrice),
		AddPriceFound:    addDateFound,
		StartDate:        tools.ConvertDateOnlyToString(startDateBook),
		EndDate:          tools.ConvertDateOnlyToString(endDateBook),
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) GetEventDeepLinkExperience(ctx *gin.Context) {
	var req GetDeepLinkExperienceParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at GetDeepLinkExperienceParams in GetOptionParams in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.DeepLinkID)
		err = fmt.Errorf("a title is required")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Println(req)
	linkID, err := tools.StringToUuid(req.DeepLinkID)
	if err != nil {
		log.Printf("Error at GetEventDeepLinkExperience tools.StringToUuid: %v, DeepLinkID: %v \n", err.Error(), req.DeepLinkID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	data, err := server.store.GetEventExperienceByDeepLinkID(ctx, db.GetEventExperienceByDeepLinkIDParams{
		DeepLinkID: linkID,
		IsComplete: true,
		IsActive:   true,
		IsActive_2: true,
	})
	if err != nil {
		log.Printf("Error at GetEventDeepLinkExperience in GetEventExperienceByDeepLinkID err: %v, user: %v\n", err, ctx.ClientIP())
		err = fmt.Errorf("error occurred while processing your request, event not found")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var ticketAvailable bool
	var locationList []ExperienceEventLocation
	var price float64
	var startDateData string
	var endDateData string
	var hasFreeTicket bool
	dateTimes, err := server.store.ListEventDateTimeOnSale(ctx, db.ListEventDateTimeOnSaleParams{
		EventInfoID: data.ID,
		Status:      "on_sale",
	})
	if err != nil || len(dateTimes) == 0 {
		if err != nil {
			log.Printf("Error at  GetEventDeepLinkExperience in ListEventDateTimeOnSale err: %v, user: %v\n", err, ctx.ClientIP())
			startDateData = "none"
			endDateData = "none"
			ticketAvailable = false
			price = 0.00
			hasFreeTicket = false
			locationList = []ExperienceEventLocation{{"none", "none", true}}
		}
	} else {
		dateTimes = HandleExEventDates(dateTimes, ctx.ClientIP())
		if len(dateTimes) == 0 {
			startDateData = "none"
			endDateData = "none"
			ticketAvailable = false
			price = 0.00
			hasFreeTicket = false
			locationList = []ExperienceEventLocation{{"none", "none", true}}
		} else {
			startTime := dateTimes[0].StartDate
			endTime := dateTimes[0].StartDate
			for _, d := range dateTimes {
				if d.StartDate.Before(startTime) {
					startTime = d.StartDate
				}
				if d.StartDate.After(endTime) {
					endTime = d.StartDate
				}
				// We only care about dates with location when setting up the location
				if !tools.ServerStringEmpty(d.Country) && !tools.ServerStringEmpty(d.State) {
					locationList = append(locationList, ExperienceEventLocation{d.State, d.Country, false})
				}
				// We only want to check ticket if ticketAvailable is still false
				tickets, err := server.store.ListEventDateTicketUser(ctx, d.EventDateTimeID)
				if err != nil {
					log.Printf("Error at GetEventDeepLinkExperience in ListEventDateTicket err: %v, user: %v\n", err, ctx.ClientIP())
				} else {
					if len(tickets) > 0 {
						priceData := HandleGetTicketStartPrice(tickets)
						if price == 0.00 {
							price = priceData
						} else {
							if priceData < price {
								price = priceData
							}
						}
						hasFreeTicket = HandleGetFreeTicket(tickets)
						ticketAvailable = true
					}
				}
			}
			startDateData = tools.ConvertDateOnlyToString(startTime)
			endDateData = tools.ConvertDateOnlyToString(endTime)
		}
	}
	priceFloat, err := tools.ConvertPrice(tools.ConvertFloatToString(price), data.Currency, req.Currency, server.config.DollarToNaira, server.config.DollarToCAD, data.ID)
	if err != nil {
		priceFloat = 0.0
		log.Printf("Error at GetEventDeepLinkExperience in ConvertPrice err: %v, user: %v\n", err, ctx.ClientIP())
	}
	res := ExperienceEventData{
		UserOptionID:      tools.UuidToString(data.OptionUserID),
		Name:              data.HostNameOption,
		IsVerified:        data.IsVerified,
		CoverImage:        data.CoverImage,
		Photos:            data.Photo,
		TicketAvailable:   ticketAvailable,
		SubEventType:      data.SubCategoryType,
		TicketLowestPrice: tools.ConvertFloatToString(priceFloat),
		EventStartDate:    startDateData,
		EventEndDate:      endDateData,
		Location:          locationList,
		ProfilePhoto:      data.Photo_2,
		HostAsIndividual:  data.HostAsIndividual,
		HostName:          data.FirstName,
		HostJoined:        tools.ConvertDateOnlyToString(data.CreatedAt),
		HostVerified:      data.IsVerified_2,
		Category:          data.Category,
		HasFreeTicket:     hasFreeTicket,
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) GetEventDateDeepLinkExperience(ctx *gin.Context) {
	var req GetEventDateDeepLinkExperienceParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at GetEventDateDeepLinkExperienceParams in GetOptionParams in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.DeepLinkID)
		err = fmt.Errorf("a title is required")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Println(req)
	linkID, err := tools.StringToUuid(req.DeepLinkID)
	if err != nil {
		log.Printf("Error at GetEventDateDeepLinkExperience tools.StringToUuid: %v, DeepLinkID: %v \n", err.Error(), req.DeepLinkID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	data, err := server.store.GetEventExperienceByDeepLinkID(ctx, db.GetEventExperienceByDeepLinkIDParams{
		DeepLinkID: linkID,
		IsComplete: true,
		IsActive:   true,
		IsActive_2: true,
	})
	if err != nil {
		log.Printf("Error at GetEventDateDeepLinkExperience in GetEventExperienceByDeepLinkID err: %v, user: %v\n", err, ctx.ClientIP())
		err = fmt.Errorf("error occurred while processing your request, event not found")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var ticketAvailable bool
	var locationList []ExperienceEventLocation
	var price float64
	var startDateData string
	var endDateData string
	var hasFreeTicket bool
	var eventDateTimeID string
	dateTimes, err := server.store.ListEventDateTimeOnSale(ctx, db.ListEventDateTimeOnSaleParams{
		EventInfoID: data.ID,
		Status:      "on_sale",
	})
	if err != nil || len(dateTimes) == 0 {
		if err != nil {
			log.Printf("Error at  GetEventDateDeepLinkExperience in ListEventDateTimeOnSale err: %v, user: %v\n", err, ctx.ClientIP())
			startDateData = "none"
			endDateData = "none"
			ticketAvailable = false
			price = 0.00
			hasFreeTicket = false
			locationList = []ExperienceEventLocation{{"none", "none", true}}
		}
	} else {
		dateTimes = HandleExEventDates(dateTimes, ctx.ClientIP())
		if len(dateTimes) == 0 {
			startDateData = "none"
			endDateData = "none"
			ticketAvailable = false
			price = 0.00
			hasFreeTicket = false
			locationList = []ExperienceEventLocation{{"none", "none", true}}
		} else {
			startTime := dateTimes[0].StartDate
			endTime := dateTimes[0].StartDate
			for _, d := range dateTimes {
				if req.EventLinkID == tools.UuidToString(d.DeepLinkID) {
					eventDateTimeID = tools.UuidToString(d.EventDateTimeID)
				}
				if d.StartDate.Before(startTime) {
					startTime = d.StartDate
				}
				if d.StartDate.After(endTime) {
					endTime = d.StartDate
				}
				// We only care about dates with location when setting up the location
				if !tools.ServerStringEmpty(d.Country) && !tools.ServerStringEmpty(d.State) {
					locationList = append(locationList, ExperienceEventLocation{d.State, d.Country, false})
				}
				// We only want to check ticket if ticketAvailable is still false
				tickets, err := server.store.ListEventDateTicketUser(ctx, d.EventDateTimeID)
				if err != nil {
					log.Printf("Error at GetEventDateDeepLinkExperience in ListEventDateTicket err: %v, user: %v\n", err, ctx.ClientIP())
				} else {
					if len(tickets) > 0 {
						priceData := HandleGetTicketStartPrice(tickets)
						if price == 0.00 {
							price = priceData
						} else {
							if priceData < price {
								price = priceData
							}
						}
						hasFreeTicket = HandleGetFreeTicket(tickets)
						ticketAvailable = true
					}
				}
			}
			startDateData = tools.ConvertDateOnlyToString(startTime)
			endDateData = tools.ConvertDateOnlyToString(endTime)
		}
	}
	priceFloat, err := tools.ConvertPrice(tools.ConvertFloatToString(price), data.Currency, req.Currency, server.config.DollarToNaira, server.config.DollarToCAD, data.ID)
	if err != nil {
		priceFloat = 0.0
		log.Printf("Error at GetEventDateDeepLinkExperience in ConvertPrice err: %v, user: %v\n", err, ctx.ClientIP())
	}
	resData := ExperienceEventData{
		UserOptionID:      tools.UuidToString(data.OptionUserID),
		Name:              data.HostNameOption,
		IsVerified:        data.IsVerified,
		CoverImage:        data.CoverImage,
		Photos:            data.Photo,
		TicketAvailable:   ticketAvailable,
		SubEventType:      data.SubCategoryType,
		TicketLowestPrice: tools.ConvertFloatToString(priceFloat),
		EventStartDate:    startDateData,
		EventEndDate:      endDateData,
		Location:          locationList,
		ProfilePhoto:      data.Photo_2,
		HostAsIndividual:  data.HostAsIndividual,
		HostName:          data.FirstName,
		HostJoined:        tools.ConvertDateOnlyToString(data.CreatedAt),
		HostVerified:      data.IsVerified_2,
		Category:          data.Category,
		HasFreeTicket:     hasFreeTicket,
	}
	res := GetEventDateDeepLinkExperienceRes{
		EventDateTimeID: eventDateTimeID,
		Data:            resData,
	}
	ctx.JSON(http.StatusOK, res)
}
