package api

import (
	"context"
	"log"

	"github.com/makuo12/ghost_server/constants"
	db "github.com/makuo12/ghost_server/db/sqlc"
	"github.com/makuo12/ghost_server/tools"
)

func HandleExEventFilterPrice(ctx context.Context, server *Server, funcName string, req ExControlEventRequest) (startPrice string, endPrice string) {
	if req.Filter.OnAddPrice {
		endPrice = req.Filter.AddMaxPrice
		startPrice = req.Filter.AddMinPrice
	} else {
		endPrice = req.Filter.MaxPrice
		startPrice = req.Filter.MinPrice
	}
	return
}

func HandleExEvent(ctx context.Context, server *Server, funcName string, req ExControlEventRequest, events []db.ListEventRow, startPrice string, endPrice string) (resData []ExperienceEventData) {
	dollarToCAD := server.config.DollarToCAD
	dollarToNaira := server.config.DollarToNaira
	for _, e := range events {
		var ticketAvailable bool
		var locationList []ExperienceEventLocation
		var price float64
		var startDateData string
		var endDateData string
		var hasFreeTicket bool
		// Before we get to event dates we need to handle filters related to event only
		eventCategories := []string{e.Category, e.CategoryTwo, e.CategoryThree}
		if !req.IsFilterEmpty {
			// First we do category types
			if len(req.Filter.CategoryType) != 0 {
				var exist = false
				for _, ca := range req.Filter.CategoryType {
					if tools.IsInList(eventCategories, ca) {
						exist = true
					}
				}
				if !exist {
					continue
				}
			}
			// We do sub category
			if len(req.Filter.SubCategory) != 0 {
				if !tools.IsInList(req.Filter.SubCategory, e.SubCategoryType) {
					continue
				}
			}
			// For price range we do it with tickets
		}
		if req.Filter.TicketAvailable && !req.IsFilterEmpty {
			dateTimes, err := server.store.ListEventDateTimeEx(ctx, db.ListEventDateTimeExParams{
				EventInfoID: e.ID,
			})
			if err != nil || len(dateTimes) == 0 {
				// We want to continue because there a no events dates matching the location put for this event
				continue
			} else {
				// If there are event dates
				// We want to check if filter was set on
				dateTimes = HandleExSearchEventDates(dateTimes, tools.UuidToString(e.ID))
				startTime := dateTimes[0].StartDate
				endTime := dateTimes[0].StartDate
				for _, d := range dateTimes {
					if d.StartDate.Before(startTime) {
						startTime = d.StartDate
					}
					if d.StartDate.After(endTime) {
						endTime = d.StartDate
					}
					// We only want to check ticket if ticketAvailable is still false
					tickets, err := server.store.ListEventDateTicketUser(ctx, d.EventDateTimeID)
					if err != nil {
						log.Printf("Error at  HandleListEventExperience in ListEventDateTicket err: %v, user: %v\n", err, tools.UuidToString(e.ID))
					} else {
						if len(tickets) > 0 {
							if HandleTicketInRange(startPrice, endPrice, tickets, e.Currency, req.Filter.TicketType, req.Search.Currency, dollarToNaira, dollarToCAD, e.ID, funcName) {
								continue
							}
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
		} else {
			dateTimes, err := server.store.ListEventDateTimeEx(ctx, db.ListEventDateTimeExParams{
				EventInfoID: e.ID,
			})
			if err != nil || len(dateTimes) == 0 {
				startDateData = "none"
				endDateData = "none"
				ticketAvailable = false
				price = 0.00
				hasFreeTicket = false
				locationList = []ExperienceEventLocation{{"none", "none", true}}
			} else {
				// If there are event dates
				// We want to check if filter was set on
				dateTimes = HandleExSearchEventDates(dateTimes, tools.UuidToString(e.ID))
				startTime := dateTimes[0].StartDate
				endTime := dateTimes[0].StartDate
				for _, d := range dateTimes {
					if d.StartDate.Before(startTime) {
						startTime = d.StartDate
					}
					if d.StartDate.After(endTime) {
						endTime = d.StartDate
					}

					// We only want to check ticket if ticketAvailable is still false
					tickets, err := server.store.ListEventDateTicketUser(ctx, d.EventDateTimeID)
					if err != nil {
						log.Printf("Error at  HandleListEventExperience in ListEventDateTicket err: %v, user: %v\n", err, tools.UuidToString(e.ID))
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

		priceFloat, err := tools.ConvertPrice(tools.ConvertFloatToString(price), e.Currency, req.Search.Currency, dollarToNaira, dollarToCAD, e.ID)
		if err != nil {
			priceFloat = 0.0
			log.Printf("Error at  HandleExEvent in ConvertPrice err: %v, user: %v\n", err, e.ID)
		}
		_, mainUrl := tools.GetImageItem(e.MainImage)
		_, userUrl := tools.GetImageItem(e.HostImage)
		_, urls := tools.GetImageListItem(e.Images)
		newData := ExperienceEventData{
			UserOptionID:      tools.UuidToString(e.OptionUserID),
			Name:              e.HostNameOption,
			IsVerified:        e.IsVerified,
			TicketAvailable:   ticketAvailable,
			SubEventType:      e.SubCategoryType,
			TicketLowestPrice: tools.ConvertFloatToString(priceFloat),
			EventStartDate:    startDateData,
			EventEndDate:      endDateData,
			Location:          locationList,
			HostAsIndividual:  e.HostAsIndividual,
			HostName:          e.FirstName,
			HostJoined:        tools.ConvertDateOnlyToString(e.CreatedAt),
			HostVerified:      e.IsVerified_2,
			Category:          e.Category,
			HasFreeTicket:     hasFreeTicket,
			MainUrl:           mainUrl,
			HostUrl:           userUrl,
			Urls:              urls,
		}
		resData = append(resData, newData)

	}
	return
}

func HandleExEventLocation(ctx context.Context, server *Server, funcName string, req ExControlEventRequest, events []db.ListEventRow, startPrice string, endPrice string) (resData []ExperienceEventData) {
	dollarToCAD := server.config.DollarToCAD
	dollarToNaira := server.config.DollarToNaira
	lat := tools.ConvertLocationStringToFloat(req.Search.Lat, 9)
	lng := tools.ConvertLocationStringToFloat(req.Search.Lng, 9)

	for _, e := range events {
		var ticketAvailable bool
		var locationList []ExperienceEventLocation
		var price float64
		var startDateData string
		var endDateData string
		var hasFreeTicket bool
		// Before we get to event dates we need to handle filters related to event only
		eventCategories := []string{e.Category, e.CategoryTwo, e.CategoryThree}
		if !req.IsFilterEmpty {
			// First we do category types
			if len(req.Filter.CategoryType) != 0 {
				var exist = false
				for _, ca := range req.Filter.CategoryType {
					if tools.IsInList(eventCategories, ca) {
						exist = true
					}
				}
				if !exist {
					continue
				}
			}
			// We do sub category
			if len(req.Filter.SubCategory) != 0 {
				if !tools.IsInList(req.Filter.SubCategory, e.SubCategoryType) {
					continue
				}
			}
			// For price range we do it with tickets
		}
		if req.Filter.TicketAvailable && !req.IsFilterEmpty {
			dateTimes, err := server.store.ListEventDateTimeExLocation(ctx, db.ListEventDateTimeExLocationParams{
				EventInfoID: e.ID,
				State:       req.Search.State,
				Country:     req.Search.Country,
				City:        req.Search.City,
				Street:      req.Search.Street,
				LlToEarth:   lat,
				LlToEarth_2: lng,
				Column8:     10000.0,
			})
			if err != nil || len(dateTimes) == 0 {
				// We want to continue because there a no events dates matching the location put for this event
				continue
			} else {
				// If there are event dates
				// We want to check if filter was set on
				dateTimes = HandleExSearchEventLocationDates(dateTimes, tools.UuidToString(e.ID))
				startTime := dateTimes[0].StartDate
				endTime := dateTimes[0].StartDate
				for _, d := range dateTimes {
					if d.StartDate.Before(startTime) {
						startTime = d.StartDate
					}
					if d.StartDate.After(endTime) {
						endTime = d.StartDate
					}
					// We only want to check ticket if ticketAvailable is still false
					tickets, err := server.store.ListEventDateTicketUser(ctx, d.EventDateTimeID)
					if err != nil {
						log.Printf("Error at  HandleListEventExperience in ListEventDateTicket err: %v, user: %v\n", err, tools.UuidToString(e.ID))
					} else {
						if len(tickets) > 0 {
							if HandleTicketInRange(startPrice, endPrice, tickets, e.Currency, req.Filter.TicketType, req.Search.Currency, dollarToNaira, dollarToCAD, e.ID, funcName) {
								continue
							}
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
		} else {
			dateTimes, err := server.store.ListEventDateTimeEx(ctx, db.ListEventDateTimeExParams{
				EventInfoID: e.ID,
			})
			if err != nil || len(dateTimes) == 0 {
				startDateData = "none"
				endDateData = "none"
				ticketAvailable = false
				price = 0.00
				hasFreeTicket = false
				locationList = []ExperienceEventLocation{{"none", "none", true}}
			} else {
				// If there are event dates
				// We want to check if filter was set on
				dateTimes = HandleExSearchEventDates(dateTimes, tools.UuidToString(e.ID))
				startTime := dateTimes[0].StartDate
				endTime := dateTimes[0].StartDate
				for _, d := range dateTimes {
					if d.StartDate.Before(startTime) {
						startTime = d.StartDate
					}
					if d.StartDate.After(endTime) {
						endTime = d.StartDate
					}

					// We only want to check ticket if ticketAvailable is still false
					tickets, err := server.store.ListEventDateTicketUser(ctx, d.EventDateTimeID)
					if err != nil {
						log.Printf("Error at  HandleListEventExperience in ListEventDateTicket err: %v, user: %v\n", err, tools.UuidToString(e.ID))
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

		priceFloat, err := tools.ConvertPrice(tools.ConvertFloatToString(price), e.Currency, req.Search.Currency, dollarToNaira, dollarToCAD, e.ID)
		if err != nil {
			priceFloat = 0.0
			log.Printf("Error at  HandleExEventLocation in ConvertPrice err: %v, user: %v\n", err, e.ID)
		}
		_, mainUrl := tools.GetImageItem(e.MainImage)
		_, userUrl := tools.GetImageItem(e.HostImage)
		_, urls := tools.GetImageListItem(e.Images)
		newData := ExperienceEventData{
			UserOptionID:       tools.UuidToString(e.OptionUserID),
			Name:               e.HostNameOption,
			IsVerified:         e.IsVerified,
			TicketAvailable:    ticketAvailable,
			SubEventType:       e.SubCategoryType,
			TicketLowestPrice:  tools.ConvertFloatToString(priceFloat),
			EventStartDate:     startDateData,
			EventEndDate:       endDateData,
			Location:           locationList,
			HostAsIndividual:   e.HostAsIndividual,
			HostName:           e.FirstName,
			HostJoined:         tools.ConvertDateOnlyToString(e.CreatedAt),
			HostVerified:       e.IsVerified_2,
			Category:           e.Category,
			HasFreeTicket:      hasFreeTicket,
			MainUrl:           mainUrl,
			HostUrl:           userUrl,
			Urls:              urls,
		}
		resData = append(resData, newData)

	}
	return
}

func HandleExSearchEventLocationDates(eventDates []db.ListEventDateTimeExLocationRow, id string) []db.ListEventDateTimeExLocationRow {
	none := constants.NONE
	data := []db.ListEventDateTimeExLocationRow{}
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
			}
		} else {
			data = append(data, eventDate)
		}
	}
	return data
}

func HandleExSearchEventDates(eventDates []db.ListEventDateTimeExRow, id string) []db.ListEventDateTimeExRow {
	none := constants.NONE
	data := []db.ListEventDateTimeExRow{}
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
			}
		} else {
			data = append(data, eventDate)
		}
	}
	return data
}
