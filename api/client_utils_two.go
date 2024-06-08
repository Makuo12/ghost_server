package api

import (
	"bytes"
	"encoding/json"
	"log"
	"strings"

	db "github.com/makuo12/ghost_server/db/sqlc"
	"github.com/makuo12/ghost_server/tools"
)

func EventExSearchText(ctx *connection, payload []byte) (data []byte, hasData bool, err error) {
	dollarToNaira := ctx.server.config.DollarToNaira
	dollarToCAD := ctx.server.config.DollarToCAD
	var search *EventSearchText = &EventSearchText{}
	err = json.Unmarshal(payload, search)
	if err != nil {
		log.Printf("error decoding EventExSearchText response: %v, user: %v", err, ctx.username)
		if e, ok := err.(*json.SyntaxError); ok {
			log.Printf("syntax error at byte offset %d\n", e.Offset)
		}
		return
	}
	searchName := "%" + search.Text + "%"
	events, err := ctx.server.store.ListEventSearch(ctx.ctx, strings.ToLower(searchName))
	if err != nil {
		log.Printf("error at ListEventSearch at ListOIDSearchByName err:%v, user: %v \n", err, ctx.username)
		return
	}
	var resData []ExperienceEventData
	for _, e := range events {
		var ticketAvailable bool
		var locationList []ExperienceEventLocation
		var price float64
		var startDateData string
		var endDateData string
		var hasFreeTicket bool
		dateTimes, err := ctx.server.store.ListEventDateTimeOnSale(ctx.ctx, db.ListEventDateTimeOnSaleParams{
			EventInfoID: e.ID,
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
			dateTimes = HandleExEventDates(dateTimes, ctx.ctx.ClientIP())
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
					tickets, err := ctx.server.store.ListEventDateTicketUser(ctx.ctx, d.EventDateTimeID)
					if err != nil {
						log.Printf("Error at  HandleListEventExperience in ListEventDateTicket err: %v, user: %v\n", err, ctx.ctx.ClientIP())
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
		priceFloat, err := tools.ConvertPrice(tools.ConvertFloatToString(price), e.Currency, search.Currency, dollarToNaira, dollarToCAD, e.ID)
		if err != nil {
			priceFloat = 0.0
			log.Printf("Error at  HandleListEventExperience in ConvertPrice err: %v, user: %v\n", err, ctx.ctx.ClientIP())
		}
		newData := ExperienceEventData{
			UserOptionID:       tools.UuidToString(e.OptionUserID),
			Name:               e.HostNameOption,
			IsVerified:         e.IsVerified,
			CoverImage:         e.CoverImage,
			Photos:             e.Photo,
			TicketAvailable:    ticketAvailable,
			SubEventType:       e.SubCategoryType,
			TicketLowestPrice:  tools.ConvertFloatToString(priceFloat),
			EventStartDate:     startDateData,
			EventEndDate:       endDateData,
			Location:           locationList,
			ProfilePhoto:       e.Photo_2,
			HostAsIndividual:   e.HostAsIndividual,
			HostName:           e.FirstName,
			HostJoined:         tools.ConvertDateOnlyToString(e.CreatedAt),
			HostVerified:       e.IsVerified_2,
			Category:           e.Category,
			HasFreeTicket:      hasFreeTicket,
			PublicCoverImage:   e.PublicCoverImage,
			PublicPhotos:       e.OptionPublicPhoto,
			PublicProfilePhoto: e.HostPublicPhoto,
		}
		resData = append(resData, newData)
	}
	if len(resData) > 0 {
		res := EventSearchTextRes{
			List: resData,
		}
		resBytes := new(bytes.Buffer)

		err = json.NewEncoder(resBytes).Encode(res)
		if err != nil {
			log.Println("error showing", err)
		}
		data = resBytes.Bytes()
		hasData = true
		return
	}
	return
}

//func HandleHostReserve(ctx *connection, payload []byte) (data []byte, err error) {
//	var reserve *ListReservationDetailParams = &ListReservationDetailParams{}
//	err = json.Unmarshal(payload, reserve)
//	if err != nil {
//		log.Printf("error decoding HandleHostReserve response: %v, user: %v", err, ctx.username)
//		if e, ok := err.(*json.SyntaxError); ok {
//			log.Printf("syntax error at byte offset %d\n", e.Offset)
//		}
//		log.Printf("HandleHostReserve response: %q", payload)
//		return
//	}
//	switch reserve.MainOption {
//	case "options":
//		mainHostRes, hasHostData, err := HandleReserveHost(reserve.Selection, ctx.server, ctx.ctx, ctx.user)
//		if err != nil {
//			log.Printf("Error at  HandleHostReserve in HandleReserveHost err: %v, user: %v\n", err, ctx.user.ID)
//			hasHostData = false
//		}
//		coHostRes, hasCoData, err := HandleReserveCoHost(reserve.Selection, ctx.server, ctx.ctx, ctx.user)
//		if err != nil {
//			log.Printf("Error at  HandleHostReserve in HandleReserveCoHost err: %v, user: %v\n", err, ctx.user.ID)
//			hasCoData = false
//		}
//		var list []ReserveHostItem
//		if hasCoData && hasHostData {
//			// Has co host data and main host data
//			list = ConcatSlicesReserveItem(mainHostRes, coHostRes)

//		} else if hasHostData && !hasCoData {
//			// Has main host data but no co host data
//			list = mainHostRes
//		} else if hasCoData && !hasHostData {
//			// Has co host data but not main host data
//			list = coHostRes
//		}
//		reserveIDs := tools.HandleListReq(reserve.ReferenceIDs)
//		list = HandleHostReserveOptionSelected(reserveIDs, list)
//		if len(list) > 0 {
//			res := ListReservationDetailRes{
//				List:      list,
//				Selection: reserve.Selection,
//			}
//			resBytes := new(bytes.Buffer)
//			json.NewEncoder(resBytes).Encode(res)
//			data = resBytes.Bytes()
//		}
//	case "events":
//		mainHostRes, hasHostData, err := HandleReserveEventHost(reserve.Selection, ctx.server, ctx.ctx, ctx.user)
//		if err != nil {
//			log.Printf("Error at  ListReservationDetail in HandleReserveEventHost err: %v, user: %v\n", err, ctx.user.ID)
//			hasHostData = false
//		}
//		coHostRes, hasCoData, err := HandleReserveEventCoHost(reserve.Selection, ctx.server, ctx.ctx, ctx.user)
//		if err != nil {
//			log.Printf("Error at  ListReservationDetail in HandleReserveEventCoHost err: %v, user: %v\n", err, ctx.user.ID)
//			hasCoData = false
//		}
//		var list []DateHostItem
//		if hasCoData && hasHostData {
//			// Has co host data and main host data
//			list = ConcatSlicesDateItem(mainHostRes, coHostRes)

//		} else if hasHostData && !hasCoData {
//			// Has main host data but no co host data
//			list = mainHostRes
//		} else if hasCoData && !hasHostData {
//			// Has co host data but not main host data
//			list = coHostRes
//		}
//		if len(list) > 0 {
//			res := ReserveEventHostItem{
//				List:      list,
//				Selection: reserve.Selection,
//			}
//			resBytes := new(bytes.Buffer)
//			json.NewEncoder(resBytes).Encode(res)
//			data = resBytes.Bytes()
//		}
//	}
//	return
//}

func HandleHostReserveOptionSelected(referenceIDs []string, list []ReserveHostItem) (res []ReserveHostItem) {
	for _, item := range list {
		exist := false
		for _, id := range referenceIDs {
			if item.ReferenceID == id {
				exist = true
			}
		}
		if !exist {
			res = append(res, item)
		}

	}
	return
}

func HandleNotificationListen(ctx *connection, payload []byte) (data []byte, hasData bool, err error) {
	//var resData []MessageContactItem
	var currentTime *CurrentTime = &CurrentTime{}
	err = json.Unmarshal(payload, currentTime)
	if err != nil {
		log.Printf("error decoding HandleNotificationListen response: %v, user: %v", err, ctx.userID)
		if e, ok := err.(*json.SyntaxError); ok {
			log.Printf("syntax error at byte offset %d\n", e.Offset)
		}
		return
	}
	current, err := tools.ConvertStringToTime(currentTime.Time)
	if err != nil {
		log.Printf("error HandleNotificationListen decoding ConvertStringToTime response: %v, user: %v", err, ctx.userID)
		return
	}
	// We first check the redis storage
	notifications, err := ctx.server.store.ListNotificationByTime(ctx.ctx, db.ListNotificationByTimeParams{
		UserID:    ctx.userID,
		CreatedAt: current,
	})
	if err != nil {
		log.Printf("error at HandleNotificationListen at ListNotificationByTime err:%v, user: %v \n", err, ctx.userID)
		return
	}
	var resData []NotificationItem
	for _, n := range notifications {
		dataNotification := NotificationItem{
			ID:        tools.UuidToString(n.ID),
			Type:      n.Type,
			Header:    n.Header,
			Message:   n.Message,
			Handled:   n.Handled,
			CreatedAt: tools.ConvertTimeToString(n.CreatedAt),
		}
		resData = append(resData, dataNotification)
	}
	if len(resData) > 0 {
		res := ListNotificationListenRes{
			List:        resData,
			CurrentTime: tools.ConvertTimeToString(current),
		}
		//log.Println("res", res)
		resBytes := new(bytes.Buffer)

		err = json.NewEncoder(resBytes).Encode(res)
		if err != nil {
			log.Println("error showing", err)
		}
		hasData = true
		data = resBytes.Bytes()
		return
	}
	return
}
