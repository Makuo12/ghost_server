package api

import (
	db "flex_server/db/sqlc"
	"flex_server/tools"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func HandleWishlistOptionExperience(ctx *gin.Context, server *Server, user db.User, req WishlistOffsetParams) (res ListExperienceWishlistOptionRes, err error, hasData bool) {
	dollarToNaira := server.config.DollarToNaira
	dollarToCAD := server.config.DollarToCAD
	var onLastIndex bool
	hasData = true
	wishlistID, err := tools.StringToUuid(req.WishlistID)
	if err != nil {
		log.Printf("Error at  HandleWishlistOptionExperience in StringToUuid err: %v, user: %v\n", err, user.ID)
		hasData = false
		err = fmt.Errorf("could not perform your request")
		return
	}
	count, err := server.store.GetWishlistItemCount(ctx, db.GetWishlistItemCountParams{
		WishlistID:     wishlistID,
		UserID:         user.ID,
		MainOptionType: "options",
	})
	if err != nil {
		log.Printf("Error at  HandleWishlistOptionExperience in GetOptionExperienceCount err: %v, user: %v\n", err, user.ID)
		hasData = false
		err = fmt.Errorf("could not perform your request")
		return
	}
	if count <= int64(req.OptionOffset) {
		err = nil
		hasData = false
		return
	}
	optionUserIDs, err := server.store.ListWishlistItem(ctx, db.ListWishlistItemParams{
		UserID:         user.ID,
		WishlistID:     wishlistID,
		MainOptionType: "options",
		Limit:          10,
		Offset:         int32(req.OptionOffset),
	})
	if err != nil {
		if err == db.ErrorRecordNotFound {
			err = nil
			hasData = false
			return
		} else {
			log.Printf("Error at  HandleWishlistOptionExperience in ListOptionInfo err: %v, user: %v\n", err, user.ID)
			hasData = false
			err = fmt.Errorf("an error occurred while getting your data")
			return
		}
	}

	var resData []ExperienceOptionData
	for _, id := range optionUserIDs {
		data, err := server.store.GetOptionExperienceByOptionUserID(ctx, db.GetOptionExperienceByOptionUserIDParams{
			OptionUserID:    id,
			IsActive:        true,
			IsComplete:      true,
			IsActive_2:      true,
			OptionStatusOne: "list",
			OptionStatusTwo: "staged",
		})
		if err != nil {
			log.Printf("Error at  HandleWishlistOptionExperience in GetOptionExperienceCount err: %v, user: %v\n", err, user.ID)
			continue
		}
		basePrice, err := tools.ConvertPrice(tools.IntToMoneyString(data.Price), data.Currency, req.Currency, dollarToNaira, dollarToCAD, id)
		if err != nil {
			log.Printf("Error at  basePrice HandleWishlistOptionExperience in ConvertPrice err: %v, user: %v\n", err, user.ID)
			basePrice = 0.0

		}
		weekendPrice, err := tools.ConvertPrice(tools.IntToMoneyString(data.WeekendPrice), data.Currency, req.Currency, dollarToNaira, dollarToCAD, id)
		if err != nil {
			log.Printf("Error at weekendPrice HandleWishlistOptionExperience in ConvertPrice err: %v, user: %v\n", err, user.ID)
			weekendPrice = 0.0
		}
		newData := ExperienceOptionData{
			UserOptionID:     tools.UuidToString(data.OptionUserID),
			Name:             data.HostNameOption,
			IsVerified:       data.IsVerified,
			CoverImage:       data.CoverImage,
			HostAsIndividual: data.HostAsIndividual,
			BasePrice:        tools.ConvertFloatToString(basePrice),

			WeekendPrice:   tools.ConvertFloatToString(weekendPrice),
			Photos:         data.Photo,
			TypeOfShortlet: data.TypeOfShortlet,
			State:          data.State,
			Country:        data.Country,
			ProfilePhoto:   data.Photo_2,
			HostName:       data.FirstName,
			HostJoined:     tools.ConvertDateOnlyToString(data.CreatedAt),
			HostVerified:   data.IsVerified_2,
			Category:       data.Category,
		}
		resData = append(resData, newData)

	}
	err = nil
	if err == nil && hasData {
		if count <= int64(req.OptionOffset+len(optionUserIDs)) {
			onLastIndex = true
		}
		res = ListExperienceWishlistOptionRes{
			List:         resData,
			OptionOffset: req.OptionOffset + len(optionUserIDs),
			OnLastIndex:  onLastIndex,
			WishlistID:   req.WishlistID,
		}
	}
	return
}

func HandleWishlistEventExperience(ctx *gin.Context, server *Server, user db.User, req WishlistOffsetParams) (res ListExperienceWishlistEventRes, err error, hasData bool) {
	dollarToNaira := server.config.DollarToNaira
	dollarToCAD := server.config.DollarToCAD
	var onLastIndex bool
	hasData = true
	wishlistID, err := tools.StringToUuid(req.WishlistID)
	if err != nil {
		log.Printf("Error at  HandleWishlistEventExperience in StringToUuid err: %v, user: %v\n", err, user.ID)
		hasData = false
		err = fmt.Errorf("could not perform your request")
		return
	}
	log.Println("wishlistID Check", wishlistID)
	if err != nil {
		log.Printf("Error at  HandleListEventExperience in GetOptionExperienceCount err: %v, user: %v\n", err, user.ID)
		hasData = false
		err = fmt.Errorf("could not perform your request")
		return
	}
	count, err := server.store.GetWishlistItemCount(ctx, db.GetWishlistItemCountParams{
		WishlistID:     wishlistID,
		UserID:         user.ID,
		MainOptionType: "events",
	})
	if err != nil {
		log.Printf("Error at  HandleWishlistEventExperience in GetOptionExperienceCount err: %v, user: %v\n", err, user.ID)
		hasData = false
		err = fmt.Errorf("could not perform your request")
		return
	}
	if count <= int64(req.OptionOffset) {
		err = nil
		hasData = false
		return
	}

	optionUserIDs, err := server.store.ListWishlistItem(ctx, db.ListWishlistItemParams{
		UserID:         user.ID,
		WishlistID:     wishlistID,
		MainOptionType: "events",
		Limit:          10,
		Offset:         int32(req.OptionOffset),
	})
	if err != nil {
		if err == db.ErrorRecordNotFound {
			err = nil
			hasData = false
			return
		} else {
			log.Printf("Error at  HandleWishlistEventExperience in .ListWishlistItem err: %v, user: %v\n", err, user.ID)
			hasData = false
			err = fmt.Errorf("an error occurred while getting your data")
			return
		}
	}

	var resData []ExperienceEventData
	for _, id := range optionUserIDs {
		log.Println("wishlistID optionUserID", id)
		data, err := server.store.GetEventExperienceByOptionUserID(ctx, db.GetEventExperienceByOptionUserIDParams{
			OptionUserID:    id,
			IsComplete:      true,
			IsActive:        true,
			IsActive_2:      true,
			OptionStatusOne: "list",
			OptionStatusTwo: "staged",
		})
		if err != nil {
			log.Printf("Error at  HandleWishlistOptionExperience in GetEventExperienceByOptionUserID err: %v, user: %v\n", err, user.ID)
			continue
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
			startDateData = "none"
			endDateData = "none"
			ticketAvailable = false
			price = 0.00
			hasFreeTicket = false
			locationList = []ExperienceEventLocation{{"none", "none", true}}
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
						log.Printf("Error at  HandleListEventExperience in ListEventDateTicket err: %v, user: %v\n", err, ctx.ClientIP())
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
		priceFloat, err := tools.ConvertPrice(tools.ConvertFloatToString(price), data.Currency, req.Currency, dollarToNaira, dollarToCAD, data.ID)
		if err != nil {
			priceFloat = 0.0
			log.Printf("Error at  HandleListEventExperience in ConvertPrice err: %v, user: %v\n", err, ctx.ClientIP())
		}

		newData := ExperienceEventData{
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
		resData = append(resData, newData)
	}
	if err == nil && hasData {
		if count <= int64(req.OptionOffset+len(optionUserIDs)) {
			onLastIndex = true
		}
		log.Println("wishlist resData ", resData)
		res = ListExperienceWishlistEventRes{
			List:         resData,
			OptionOffset: req.OptionOffset + len(optionUserIDs),
			OnLastIndex:  onLastIndex,
			WishlistID:   req.WishlistID,
		}
	}
	return
}
