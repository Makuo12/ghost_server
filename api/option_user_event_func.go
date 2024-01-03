package api

import (
	"flex_server/constants"
	db "flex_server/db/sqlc"
	"flex_server/tools"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

//func HandleListEventExperience(ctx *gin.Context, server *Server, req ExperienceOffsetParams) (res ListExperienceEventRes, err error, hasData bool) {
//	dollarToNaira := server.config.DollarToNaira
//	dollarToCAD := server.config.DollarToCAD
//	var onLastIndex bool
//	hasData = true
//	count, err := server.store.GetEventExperienceCount(ctx, db.GetEventExperienceCountParams{
//		IsComplete:      true,
//		IsActive:        true,
//		IsActive_2:      true,
//		MainOptionType:  "events",
//		Category:        req.Type,
//		OptionStatusOne: "list",
//		OptionStatusTwo: "staged",
//	})
//	if err != nil {
//		log.Printf("Error at  HandleListEventExperience in GetOptionExperienceCount err: %v, user: %v\n", err, ctx.ClientIP())
//		hasData = false
//		err = fmt.Errorf("could not perform your request")
//		return
//	}
//	if count <= int64(req.OptionOffset) {
//		err = nil
//		hasData = false
//		return
//	}
//	var optionInfos []db.ListEventExperienceByLocationRow
//	optionInfos, err = server.store.ListEventExperienceByLocation(ctx, db.ListEventExperienceByLocationParams{
//		IsComplete:      true,
//		IsActive:        true,
//		IsActive_2:      true,
//		MainOptionType:  "events",
//		Category:        req.Type,
//		OptionStatusOne: "list",
//		OptionStatusTwo: "staged",
//		Country:         req.Country,
//		State:           req.State,
//		Limit:           10,
//		Offset:          int32(req.OptionOffset),
//	})
//	if err != nil {
//		if err == db.ErrorRecordNotFound {
//			err = nil
//			hasData = false
//			return
//		} else {
//			log.Printf("Error at  HandleListEventExperience in ListEventExperienceByLocation err: %v, user: %v\n", err, ctx.ClientIP())
//			hasData = false
//			err = fmt.Errorf("an error occurred while getting your data")
//			return
//		}
//	}

//	var resData []ExperienceEventData
//	for _, data := range optionInfos {
//		_, locationList, price, ticketAvailable, startDateData, endDateData, hasFreeTicket := SetupExperienceEventData(ctx, server, db.ListEventExperienceRow{}, data, false, "HandleEventExperienceToRedis")
//		priceFloat, err := tools.ConvertPrice(tools.ConvertFloatToString(price), optionInfos[0].Currency, req.Currency, dollarToNaira, dollarToCAD, optionInfos[0].ID)
//		if err != nil {
//			priceFloat = 0.0
//		}
//		newData := ExperienceEventData{
//			UserOptionID:      tools.UuidToString(data.OptionUserID),
//			Name:              data.HostNameOption,
//			IsVerified:        data.IsVerified,
//			CoverImage:        data.CoverImage,
//			Photos:            data.Photo,
//			TicketAvailable:   ticketAvailable,
//			SubEventType:      data.SubCategoryType,
//			TicketLowestPrice: tools.ConvertFloatToString(priceFloat),
//			EventStartDate:    startDateData,
//			EventEndDate:      endDateData,
//			Location:          locationList,
//			ProfilePhoto:      data.Photo_2,
//			HostAsIndividual:  data.HostAsIndividual,
//			HostName:          data.FirstName,
//			HostJoined:        tools.ConvertDateOnlyToString(data.CreatedAt),
//			HostVerified:      data.IsVerified_2,
//			Category:          data.Category,
//			HasFreeTicket:     hasFreeTicket,
//		}
//		resData = append(resData, newData)
//	}
//	if err == nil && hasData {
//		if count <= int64(req.OptionOffset+len(optionInfos)) {
//			onLastIndex = true
//		}
//		res = ListExperienceEventRes{
//			List:         resData,
//			OptionOffset: req.OptionOffset + len(optionInfos),
//			OnLastIndex:  onLastIndex,
//			Category:     req.Type,
//		}
//	}
//	return
//}

func HandleGetTicketStartPrice(tickets []db.EventDateTicket) (price float64) {
	// We want to loop through the price to get the lowest
	price = tools.ConvertStringToFloat(tools.IntToMoneyString(tickets[0].Price))
	for _, t := range tickets {
		p := tools.ConvertStringToFloat(tools.IntToMoneyString(t.Price))
		if p < price {
			price = p
		}
	}
	return
}

func HandleGetFreeTicket(tickets []db.EventDateTicket) (hasFreeTicket bool) {

	for _, t := range tickets {
		if t.Type == "free" {
			hasFreeTicket = true
		}
	}
	return
}

func HandleExEventDates(eventDates []db.ListEventDateTimeOnSaleRow, id string) []db.ListEventDateTimeOnSaleRow {
	none := constants.NONE
	data := []db.ListEventDateTimeOnSaleRow{}
	for _, eventDate := range eventDates {
		if eventDate.Type != "single" {
			// We want to create a new data base on the list of dates
			for _, date := range eventDate.EventDates {
				startDate, err := tools.ConvertDateOnlyStringToDate(date)
				if err != nil {
					log.Printf("Error at  HandleReserveEventHostDates in ConvertDateOnlyStringToDate err: %v, user: %v, startTimeType: %v\n", err, id, "leave_before")
					continue
				}
				eventDate.StartDate = startDate
				eventDate.EndDate = startDate
				eventDate.EventDates = []string{none}
				data = append(data, eventDate)
				//if eventDate.StartDate.Add(time.Hour * 10).After(time.Now()) {
				//	eventDate.StartDate = startDate
				//	eventDate.EndDate = startDate
				//	eventDate.EventDates = []string{none}
				//	data = append(data, eventDate)
				//}
			}
		} else {
			data = append(data, eventDate)
			//if eventDate.StartDate.Add(time.Hour * 10).After(time.Now()) {
			//	data = append(data, eventDate)
			//}
		}
	}
	return data
}

func HandleDetailEventExperience(ctx *gin.Context, server *Server, req ExperienceDetailParams) (res ExperienceEventDetailRes, hasData bool, err error) {

	optionUserID, err := tools.StringToUuid(req.OptionUserID)
	if err != nil {
		log.Printf("Error at  HandleDetailOptionExperience in StringToUuid err: %v, user: %v, optionID: %v\n", err, ctx.ClientIP(), optionUserID)
		hasData = false
		return
	}
	option, err := server.store.GetOptionInfoByOptionUserID(ctx, optionUserID)
	if err != nil {
		log.Printf("Error at  HandleDetailOptionExperience in GetOptionInfoByOptionUserID err: %v, user: %v, optionID: %v\n", err, ctx.ClientIP(), option.ID)
		hasData = false
		return
	}

	hostLanguages, hostBio, userID, cohost, book, caption := HandleExCommon(option, server, ctx)

	question, des, cancelPolicy, _, totalReviewCount := HandleExCommonTwo(option, server, ctx, userID)

	dates := HandleExEventDateTime(option, server, ctx)

	res = ExperienceEventDetailRes{
		HostLanguages:    hostLanguages,
		CoHost:           cohost,
		Des:              des,
		CancelPolicy:     cancelPolicy,
		BookMethod:       book,
		EventDateTimes:   dates,
		Question:         question,
		TotalReviewCount: totalReviewCount,
		HostBio:          hostBio,
		Captions:         caption,
	}
	log.Println("resEvent: ", res)
	hasData = true
	err = nil
	return
}

func HandleExEventDateTime(option db.OptionsInfo, server *Server, ctx *gin.Context) (dates []ExEventDateTimes) {
	none := constants.NONE
	locationEmpty := ExEventDateTimesLocation{none, none, none, none, none, none, true}
	dateEmpty := ExEventDateTimes{none, none, none, none, none, none, none, none, locationEmpty, none, none, true}
	dateTimes, err := server.store.ListEventDateTimeOnSale(ctx, db.ListEventDateTimeOnSaleParams{
		EventInfoID: option.ID,
		Status:      "on_sale",
	})
	log.Println("see dates ", dateTimes)
	if err != nil || len(dateTimes) == 0 {
		if err != nil {
			log.Printf("Error at  HandleDetailOptionExperience in GetOptionInfoByOptionUserID err: %v, user: %v, optionID: %v\n", err, ctx.ClientIP(), option.ID)
		}
		log.Println("dates empty")
		dates = []ExEventDateTimes{dateEmpty}
		return
	}
	currentDate := time.Now().Add(time.Hour)
	dateTimes = HandleExEventDates(dateTimes, ctx.ClientIP())
	log.Println("see dates 2", dateTimes)
	for _, d := range dateTimes {
		// We need to check if the date is available
		// ! this would convert the true to false
		log.Println("startDate ", d.StartDate)
		log.Println("currentDate ", currentDate)
		log.Println("endDate ", d.EndDate)
		if tools.ConvertDateOnlyToString(d.EndDate) == tools.ConvertDateOnlyToString(currentDate) || currentDate.Before(d.EndDate) {
			var location ExEventDateTimesLocation
			if tools.ServerStringEmpty(d.State) || tools.ServerStringEmpty(d.Country) {
				location = locationEmpty
			} else {
				location = ExEventDateTimesLocation{
					State:   d.State,
					Country: d.Country,
					Street:  d.Street,
					City:    d.City,
					Lat:     tools.ConvertFloatToLocationString(d.Geolocation.P.Y, 9),
					Lng:     tools.ConvertFloatToLocationString(d.Geolocation.P.X, 9),
					IsEmpty: false,
				}
			}
			date := ExEventDateTimes{
				Name:      d.Name,
				StartDate: tools.ConvertDateOnlyToString(d.StartDate),
				EndDate:   tools.ConvertDateOnlyToString(d.EndDate),
				Type:      d.Type,
				StartTime: d.StartTime,
				EndTime:   d.EndTime,
				MainID:    tools.UuidToString(d.ID),
				Timezone:  d.TimeZone,
				Location:  location,
				RandomID:  tools.UuidToString(uuid.New()),
				Status:    d.Status,
				IsEmpty:   false,
			}
			dates = append(dates, date)
		}

	}
	if len(dates) < 1 {
		dates = []ExEventDateTimes{dateEmpty}
	}
	return
}
