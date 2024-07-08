package api

import (
	"log"
	"time"

	db "github.com/makuo12/ghost_server/db/sqlc"
	"github.com/makuo12/ghost_server/tools"

	"github.com/gin-gonic/gin"

	"github.com/google/uuid"
)

func HandleReserveEventHost(selection string, server *Server, ctx *gin.Context, user db.User) (eventDateData []DateHostItem, hasData bool, err error) {
	// Went want to get all event dates for this tickets
	eventDates, err := server.store.ListEventDateTimeHost(ctx, db.ListEventDateTimeHostParams{
		CoUserID: tools.UuidToString(user.UserID),
		HostID:   user.ID,
	})
	if err != nil {
		log.Printf("Error at  HandleReserveEventHost in server.store.ListEventDateTimeHost err: %v, user: %v, startTimeType: %v\n", err, user.ID, "event_date")
		return
	}
	for _, d := range eventDates {
		var startDate time.Time
		var endDate time.Time
		currentTime := time.Now()
		currentDateOnlyString := tools.ConvertTimeToStringDateOnly(currentTime)
		switch d.Type {
		case "recurring":
			startDate, err = tools.ConvertDateOnlyStringToDate(d.EventDate)
			if err != nil {
				log.Printf("Error at  HandleReserveEventHost in ConvertDateOnlyStringToDate err: %v, user: %v, startTimeType: %v\n", err, user.ID, "startDate")
				continue
			}
			endDate = startDate
		case "single":
			startDate = d.StartDate
			endDate = d.EndDate

		}
		hasItem, dateItem, err := HandleEventTime(selection, server, ctx, user, d, currentTime, currentDateOnlyString, startDate, endDate, d.ScanCode, d.Reservations, tools.UuidToString(d.OptionID))
		if err != nil {
			log.Printf("Error at  HandleReserveEventHost in HandleEventTime err: %v, user: %v, startTimeType: %v\n", err, user.ID, "event_date")
			continue
		}
		if hasItem {
			eventDateData = append(eventDateData, dateItem)
		}
	}
	if len(eventDateData) > 0 {
		hasData = true
	}
	return
}

func HandleEventTime(selection string, server *Server, ctx *gin.Context, user db.User, eventDate db.ListEventDateTimeHostRow, currentTime time.Time, currentDateOnlyString string, startDate time.Time, endDate time.Time, canScanCode bool, canReserve bool, dOptionID string) (hasData bool, dateItem DateHostItem, err error) {
	startTimeString := "none"
	endTimeString := "none"
	startDateString := tools.ConvertDateOnlyToString(startDate)
	endDateString := tools.ConvertDateOnlyToString(endDate)
	if !tools.ServerStringEmpty(eventDate.StartTime) {
		startTimeString = eventDate.StartTime
	}
	if !tools.ServerStringEmpty(eventDate.EndTime) {
		endTimeString = eventDate.EndTime
	}
	startTime, err := tools.ConvertDateTimeStringToTime(startDateString, startTimeString)
	if err != nil {
		log.Printf("Error at  HandleEventTime in ConvertDateTimeStringToTime err: %v, user: %v, startTimeType: %v\n", err, user.ID, "start_time")
		startTime = startDate
	}
	if !tools.ServerStringEmpty(eventDate.TimeZone) {
		startTimezone, err := tools.ConvertToTimeZone(startTime, eventDate.TimeZone)
		if err != nil {
			log.Printf("Error at  HandleEventTime in ConvertToTimeZone err: %v, user: %v, startTimeType: %v\n", err, user.ID, "start_time_zone")
		} else {
			startTime = startTimezone
		}
	}
	endTime, err := tools.ConvertDateTimeStringToTime(endDateString, endTimeString)
	if err != nil {
		log.Printf("Error at  HandleEventTime in ConvertDateTimeStringToTime err: %v, user: %v, endTimeType: %v\n", err, user.ID, "end_time")
		endTime = endDate
	}
	if !tools.ServerStringEmpty(eventDate.TimeZone) {
		endTimezone, err := tools.ConvertToTimeZone(endTime, eventDate.TimeZone)
		if err != nil {
			log.Printf("Error at  HandleEventTime in ConvertToTimeZone err: %v, user: %v, endTimeType: %v\n", err, user.ID, "end_time_zone")
		} else {
			endTime = endTimezone
		}
	}
	ticketData := []TicketHostItem{}
	// We add four hours to time because events never end on time
	endTime = tools.AddHoursToTime(endTime, 6)
	ticketEmpty := TicketHostItem{"none", 0, 0, "none", "none", true}
	tickets, err := server.store.ListEventDateTicket(ctx, eventDate.EventDateTimeID)
	if err != nil {
		log.Printf("Error at  HandleEventTime in ListEventDateTicket err: %v, user: %v, eventDateID: %v\n", err, user.ID, eventDate.EventDateTimeID)
		ticketData = append(ticketData, ticketEmpty)
	} else {
		for _, t := range tickets {
			bookCount := HandleTicketCapacity(server, ctx, startDate, t.ID, user, eventDate.EventDateTimeID)
			data := TicketHostItem{
				Grade:        t.Level,
				Capacity:     int(t.Capacity),
				CapacityBook: bookCount,
				TicketType:   t.TicketType,
				Type:         t.Type,
				IsEmpty:      false,
			}
			ticketData = append(ticketData, data)
		}
	}
	switch selection {
	case "happening_today":
		if currentTime.After(startTime) && currentTime.Before(endTime) {
			dateItem = DateHostItem{
				StartDate:         tools.ConvertDateOnlyToString(startDate),
				EndDate:           tools.ConvertDateOnlyToString(endDate),
				StartTime:         eventDate.StartTime,
				EndTime:           eventDate.EndTime,
				HostNameOption:    eventDate.HostNameOption,
				EventID:           dOptionID,
				EventDateTimeID:   tools.UuidToString(eventDate.EventDateTimeID),
				EventDateTimeType: eventDate.Type,
				Tickets:           ticketData,
				Status:            eventDate.EventStatus,
				Timezone:          eventDate.TimeZone,
				HostMethod:        eventDate.HostType,
				CanReserve:        canReserve,
				CanScanCode:       canScanCode,
				OptionStatus:      eventDate.OptionStatus,
				MainImage:        eventDate.MainImage,
			}
			hasData = true
		}
	case "upcoming":
		if currentTime.Before(startTime) {
			dateItem = DateHostItem{
				StartDate:         tools.ConvertDateOnlyToString(startDate),
				EndDate:           tools.ConvertDateOnlyToString(endDate),
				StartTime:         eventDate.StartTime,
				EndTime:           eventDate.EndTime,
				HostNameOption:    eventDate.HostNameOption,
				EventID:           dOptionID,
				EventDateTimeID:   tools.UuidToString(eventDate.EventDateTimeID),
				EventDateTimeType: eventDate.Type,
				Tickets:           ticketData,
				Status:            eventDate.EventStatus,
				Timezone:          eventDate.TimeZone,
				HostMethod:        eventDate.HostType,
				CanReserve:        canReserve,
				CanScanCode:       canScanCode,
				OptionStatus:      eventDate.OptionStatus,
				MainImage:        eventDate.MainImage,
			}
			hasData = true
		}
	case "occurred":
		if currentTime.After(endTime) {
			dateItem = DateHostItem{
				StartDate:         tools.ConvertDateOnlyToString(startDate),
				EndDate:           tools.ConvertDateOnlyToString(endDate),
				StartTime:         eventDate.StartTime,
				EndTime:           eventDate.EndTime,
				HostNameOption:    eventDate.HostNameOption,
				EventID:           dOptionID,
				EventDateTimeID:   tools.UuidToString(eventDate.EventDateTimeID),
				EventDateTimeType: eventDate.Type,
				Tickets:           ticketData,
				Status:            eventDate.EventStatus,
				Timezone:          eventDate.TimeZone,
				HostMethod:        eventDate.HostType,
				CanReserve:        canReserve,
				CanScanCode:       canScanCode,
				OptionStatus:      eventDate.OptionStatus,
				MainImage:        eventDate.MainImage,
			}
			hasData = true
		}
	}

	return
}

func HandleTicketCapacity(server *Server, ctx *gin.Context, startTime time.Time, ticketID uuid.UUID, user db.User, eventDateID uuid.UUID) int {
	ticketCount, err := server.store.CountChargeTicketReferenceByStartDate(ctx, db.CountChargeTicketReferenceByStartDateParams{
		Date:        startTime,
		TicketID:    ticketID,
		Cancelled:   false,
		IsComplete:  true,
		EventDateID: eventDateID,
	})
	if err != nil {
		log.Printf("Error at  HandleTicketCapacity in CountChargeTicketReferenceByStartDate err: %v, user: %v, ticketID: %v\n", err, user.ID, ticketID)
		return 0
	}
	return int(ticketCount)
}

func HandleReserveEventHostOptionDates(eventDates []db.ListEventDateTimeByOptionRow, user db.User) []db.ListEventDateTimeByUserRow {
	data := []db.ListEventDateTimeByUserRow{}
	for _, eventDate := range eventDates {
		if eventDate.Type != "single" {
			// We want to create a new data base on the list of dates
			for _, date := range eventDate.EventDates {
				startDate, err := tools.ConvertDateOnlyStringToDate(date)
				if err != nil {
					log.Printf("Error at  HandleReserveEventHostDates in ConvertDateOnlyStringToDate err: %v, user: %v, startTimeType: %v\n", err, user.ID, "leave_before")
					continue
				}
				eventDate.StartDate = startDate
				eventDate.EndDate = startDate
				eventDate.EventDates = []string{"none"}
				data = append(data, db.ListEventDateTimeByUserRow(eventDate))
			}
		} else {
			data = append(data, db.ListEventDateTimeByUserRow(eventDate))
		}
	}
	return data
}

func ReserveEventHostItemOffset(data []DateHostItem, offset int, limit int) []DateHostItem {
	// If offset is greater than or equal to the length of data, return an empty slice.
	if offset >= len(data) {
		return []DateHostItem{}
	}

	// Calculate the end index based on offset and limit.
	end := offset + limit

	// If the end index is greater than the length of data, set it to the length of data.
	if end > len(data) {
		end = len(data)
	}

	// Return a subset of data starting from the offset and up to the end index.
	return data[offset:end]
}
