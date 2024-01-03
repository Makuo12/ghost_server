package api

import (
	db "flex_server/db/sqlc"
	"flex_server/tools"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

func (server *Server) CreateEventDateTime(ctx *gin.Context) {
	var req CreateEventDateTimeParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at CreateEventDateTimeParams in ShouldBindJSON: %v, EventInfoID: %v \n", err.Error(), req.EventInfoID)
		err = fmt.Errorf("please select the day you would want this event to start")
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
	user, _, option, isCoHost, userCoHost, err := HandleGetCompleteOptionEditEventDateTimes(eventID, ctx, server, true)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	_, err = server.store.GetEventInfoByOption(ctx, db.GetEventInfoByOptionParams{
		OptionID:   option.ID,
		ID:         user.ID,
		IsComplete: true,
	})
	if err != nil {
		log.Printf("There an error at CreateEventDateTime at GetEventInfoByOption: %v, EventInfoID: %v, userID: %v \n", err.Error(), eventID, user.ID)
		err = fmt.Errorf("could not create your date, please try again using the format on the app")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	publishDate, err := tools.ConvertDateOnlyStringToDate(tools.FakeDate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	publishTime, err := tools.ConvertStringToTimeOnly("00:00")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	// First we want to know the type
	var eventDT db.EventDateTime
	switch req.Type {
	case "single":
		startDate, err := tools.ConvertDateOnlyStringToDate(req.StartDate)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		endDate, err := tools.ConvertDateOnlyStringToDate(req.EndDate)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		eventDT, err = server.store.CreateEventDateTime(ctx, db.CreateEventDateTimeParams{
			EventInfoID: eventID,
			StartDate:   startDate,
			EndDate:     endDate,
			Type:        req.Type,
			EventDates:  []string{"none"},
		})
		if err != nil {
			log.Printf("There an error at single CreateEventDateTime at CreateEventDateTime: %v, EventInfoID: %v, userID: %v \n", err.Error(), eventID, user.ID)

			err = fmt.Errorf("could not create your date, please try again using the format on the app")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
	case "recurring":
		dates, err := HandleEventDatesList(req.EventDates)
		if err != nil {
			log.Printf("There an error at CreateEventDateTime at HandleEventDatesList: %v, EventInfoID: %v, userID: %v \n", err.Error(), eventID, user.ID)
			err = fmt.Errorf("could not create your date, please try again using the format on the app")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		minDate, maxDate, err := FindMinMaxDates(req.EventDates)
		if err != nil {
			log.Printf("There an error at recurring CreateEventDateTime at FindMinMaxDates: %v, EventInfoID: %v, userID: %v \n", err.Error(), eventID, user.ID)
			err = fmt.Errorf("could not create your date, please try again using the format on the app")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		eventDT, err = server.store.CreateEventDateTime(ctx, db.CreateEventDateTimeParams{
			EventInfoID: eventID,
			StartDate:   minDate,
			EndDate:     maxDate,
			Type:        req.Type,
			EventDates:  dates,
		})

		if err != nil {
			log.Printf("There an error at recurring CreateEventDateTime at CreateEventDateTime: %v, EventInfoID: %v, userID: %v \n", err.Error(), eventID, user.ID)
			err = fmt.Errorf("could not create your date, please try again using the format on the app")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
	}
	// Create an event date time publish

	_, err = server.store.CreateEventDatePublish(ctx, db.CreateEventDatePublishParams{
		EventDateTimeID:      eventDT.ID,
		EventGoingPublicDate: publishDate,
		EventGoingPublicTime: publishTime,
	})
	if err != nil {
		log.Printf("There an error at CreateEventDateTime at CreateEventDatePublish: %v, EventInfoID: %v, userID: %v \n", err.Error(), eventID, user.ID)
		err = server.store.RemoveEventDateTime(ctx, eventDT.ID)
		log.Printf("There an error at CreateEventDateTime at RemoveEventDateTime: %v, EventInfoID: %v, userID: %v \n", err.Error(), eventID, user.ID)
		err = fmt.Errorf("could not create your date, please try again using the format on the app")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Lets create event date detail
	_, err = server.store.CreateEventDateDetail(ctx, db.CreateEventDateDetailParams{
		EventDateTimeID: eventDT.ID,
		StartTime:       "none",
		EndTime:         "none",
		TimeZone:        req.TimeZone,
	})
	if err != nil {
		log.Printf("There an error at CreateEventDateTime at CreateEventDateDetail: %v, EventInfoID: %v, userID: %v \n", err.Error(), eventID, user.ID)
	}
	lng := tools.ConvertLocationStringToFloat("0.0", 9)
	lat := tools.ConvertLocationStringToFloat("0.0", 9)
	geolocation := pgtype.Point{
		P:     pgtype.Vec2{X: lng, Y: lat},
		Valid: true,
	}
	// Lets create event date location
	_, err = server.store.CreateEventDateLocation(ctx, db.CreateEventDateLocationParams{
		EventDateTimeID: eventDT.ID,
		Street:          "none",
		City:            "none",
		State:           "none",
		Country:         "none",
		Postcode:        "none",
		Geolocation:     geolocation,
	})
	if err != nil {
		log.Printf("There an error at CreateEventDateTime at CreateEventDateLocation: %v, EventInfoID: %v, userID: %v \n", err.Error(), eventID, user.ID)
	}
	res := CreateEventDateTimeRes{
		EventDateTimeID:  tools.UuidToString(eventDT.ID),
		StartDate:        tools.ConvertDateOnlyToString(eventDT.StartDate),
		EndDate:          tools.ConvertDateOnlyToString(eventDT.EndDate),
		Status:           eventDT.Status,
		Type:             eventDT.Type,
		EventDates:       eventDT.EventDates,
		NeedBands:        eventDT.NeedBands,
		NeedTickets:      eventDT.NeedTickets,
		AbsorbBandCharge: eventDT.AbsorbBandCharge,
	}
	if isCoHost {
		HandleCoHostUpdateMsg(ctx, server, userCoHost, user, option, "CreateEventDateTime", "event date time", "create event date time")
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) UpdateEventDateTime(ctx *gin.Context) {
	var req UpdateEventDateTimeParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at UpdateEventDateTimeParams in ShouldBindJSON: %v, EventInfoID: %v \n", err.Error(), req.EventDateTimeID)
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
	user, _, option, isCoHost, userCoHost, err := HandleGetCompleteOptionEditEventDateTimes(eventID, ctx, server, true)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	_, err = server.store.GetEventDateTimeByOption(ctx, db.GetEventDateTimeByOptionParams{
		ID:         eventDateTimeID,
		ID_2:       user.ID,
		ID_3:       option.ID,
		IsComplete: true,
	})
	if err != nil {
		log.Printf("There an error at UpdateEventDateTime at GetEventDateTimeByOption: %v, EventInfoID: %v, userID: %v \n", err.Error(), eventDateTimeID, user.ID)
		err = fmt.Errorf("could not preform this update please try again")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var eventDT db.EventDateTime
	switch req.Type {
	case "single":
		startDate, err := tools.ConvertDateOnlyStringToDate(req.StartDate)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		endDate, err := tools.ConvertDateOnlyStringToDate(req.EndDate)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}

		// For single we just need to know if this eventDateTimeID
		isBooked := EventDateIsBooked(ctx, server, eventDateTimeID, "UpdateEventDateTime", tools.UuidToString(user.ID))
		if isBooked {
			err = fmt.Errorf("this dates cannot be changed because tickets have already been sold. To modify the date of this event, go to the app's hosting area")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}

		eventDT, err = server.store.UpdateEventDateTime(ctx, db.UpdateEventDateTimeParams{
			StartDate: pgtype.Date{
				Time:  startDate,
				Valid: true,
			},
			EndDate: pgtype.Date{
				Time:  endDate,
				Valid: true,
			},
			ID:   eventDateTimeID,
			Type: "single",
		})
		if err != nil {
			log.Printf("There an error single at UpdateEventDateTime at CreateEventDateTime: %v, EventInfoID: %v, userID: %v \n", err.Error(), eventDateTimeID, user.ID)
			err = fmt.Errorf("could not create your date, please try again using the format on the app")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}

	case "recurring":
		log.Println("At recurring")
		// We want to know if you selected a date that is booked
		valid := IsUpdatedEventDatesBooked(ctx, server, req.EventDates, user, eventDateTimeID, "UpdateEventDateTime")
		if !valid {
			err = fmt.Errorf("this dates cannot be changed because tickets have already been sold. To modify the date of this event, go to the app's hosting area")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		// If you selected a date that is already selected we remove the date
		dates, err := HandleEventDatesUpdateList(server, ctx, eventDateTimeID, req.EventDates)
		log.Println("results", dates)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		if len(dates) == 0 {
			err = fmt.Errorf("could not update your dates. You must have at least a single date for recurring events")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}

		minDate, maxDate, err := FindMinMaxDates(req.EventDates)
		if err != nil {
			log.Printf("There an error at recurring UpdateEventDateTime at FindMinMaxDates: %v, EventInfoID: %v, userID: %v \n", err.Error(), eventDateTimeID, user.ID)
			err = fmt.Errorf("could not create your date, please try again using the format on the app")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		eventDT, err = server.store.UpdateEventDateTime(ctx, db.UpdateEventDateTimeParams{
			EventDates: dates,
			ID:         eventDateTimeID,
			Type:       "recurring",
			StartDate: pgtype.Date{
				Time:  minDate,
				Valid: true,
			},
			EndDate: pgtype.Date{
				Time:  maxDate,
				Valid: true,
			},
		})
		log.Println("At recurring gm")
		if err != nil {
			log.Printf("There an error at UpdateEventDateTime at UpdateEventDateTime: %v, EventInfoID: %v, userID: %v \n", err.Error(), eventDateTimeID, user.ID)
			err = fmt.Errorf("could not update your dates. You must have at least a single date for recurring events")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
	}
	res := CreateEventDateTimeRes{
		EventDateTimeID:  tools.UuidToString(eventDT.ID),
		StartDate:        tools.ConvertDateOnlyToString(eventDT.StartDate),
		EndDate:          tools.ConvertDateOnlyToString(eventDT.EndDate),
		Status:           eventDT.Status,
		Type:             eventDT.Type,
		EventDates:       eventDT.EventDates,
		NeedBands:        eventDT.NeedBands,
		NeedTickets:      eventDT.NeedTickets,
		AbsorbBandCharge: eventDT.AbsorbBandCharge,
	}
	if isCoHost {
		HandleCoHostUpdateMsg(ctx, server, userCoHost, user, option, "UpdateEventDateTime", "event date time", "update event date time")
	}
	ctx.JSON(http.StatusOK, res)
}

// Controls refers to needTickets, needBands, absorbFee
func (server *Server) UpdateEventDateTimeControls(ctx *gin.Context) {
	var req UpdateEventDateTimeControlParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at UpdateEventDateTimeControlParams in ShouldBindJSON: %v, EventInfoID: %v \n", err.Error(), req.EventDateTimeID)
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
	user, _, option, isCoHost, userCoHost, err := HandleGetCompleteOptionEditEventDateTimes(eventID, ctx, server, true)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	_, err = server.store.GetEventDateTimeByOption(ctx, db.GetEventDateTimeByOptionParams{
		ID:         eventDateTimeID,
		ID_2:       user.ID,
		ID_3:       option.ID,
		IsComplete: true,
	})
	if err != nil {
		log.Printf("There an error at UpdateEventDateTime at GetEventDateTimeByOption: %v, EventInfoID: %v, userID: %v \n", err.Error(), eventID, user.ID)
		err = fmt.Errorf("could not preform this update please try again")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	controls, err := server.store.UpdateEventDateTime(ctx, db.UpdateEventDateTimeParams{
		ID:   eventID,
		Type: req.Type,
		NeedBands: pgtype.Bool{
			Bool:  req.NeedBands,
			Valid: true,
		},
		NeedTickets: pgtype.Bool{
			Bool:  req.NeedTickets,
			Valid: true,
		},
		AbsorbBandCharge: pgtype.Bool{
			Bool:  req.AbsorbBandCharge,
			Valid: true,
		},
	})
	if err != nil {
		log.Printf("There an error at UpdateEventDateTime at UpdateEventDateTime: %v, EventInfoID: %v, userID: %v \n", err.Error(), eventID, user.ID)
		err = fmt.Errorf("could not preform this update please try again")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if isCoHost {
		HandleCoHostUpdateMsg(ctx, server, userCoHost, user, option, "UpdateEventDateTimeControls", "event date time", "update event date time controls")
	}
	res := UpdateEventDateTimeControlRes{
		NeedTickets:      controls.NeedTickets,
		NeedBands:        controls.NeedBands,
		AbsorbBandCharge: controls.AbsorbBandCharge,
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) ListEventDateItems(ctx *gin.Context) {
	var req EventDateParams
	var onLastIndex = false
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at  EventDateParams in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.ItemOffset)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Println(req)
	eventID, err := tools.StringToUuid(req.EventInfoID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, _, option, _, _, err := HandleGetCompleteOptionEditEventDateTimes(eventID, ctx, server, true)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// We want to verify access to this resource
	_, err = server.store.GetEventInfoByOption(ctx, db.GetEventInfoByOptionParams{
		OptionID:   option.ID,
		ID:         user.ID,
		IsComplete: true,
	})
	if err != nil {
		log.Printf("Error at  ListEventDateItems in GetEventDateTimeByOption err: %v, user: %v\n", err, user.ID)
		err = fmt.Errorf("you don't have access to this resource")
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}
	empty, uhmData, err := HandleListEventDates(ctx, server, option.ID, user, "ListEventDateItem", req.ItemOffset)
	if err != nil {
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}
	count, err := server.store.GetEventDateTimeCount(ctx, eventID)
	if err != nil {
		log.Printf("Error at  ListEventDateItems in GetEventDateTimeCount err: %v, user: %v\n", err, user.ID)
		res := ListEventDateItem{
			List:        empty,
			IsEmpty:     true,
			OptionData:  uhmData,
			ItemOffset:  0,
			OnLastIndex: onLastIndex,
		}
		ctx.JSON(http.StatusNoContent, res)
		return
	}
	if count <= int64(req.ItemOffset) || count == 0 {
		res := ListEventDateItem{
			List:        empty,
			IsEmpty:     true,
			OptionData:  uhmData,
			ItemOffset:  0,
			OnLastIndex: onLastIndex,
		}
		ctx.JSON(http.StatusNoContent, res)
		return
	}
	eventDateTimes, err := server.store.ListEventDateTime(ctx, db.ListEventDateTimeParams{
		EventInfoID: eventID,
		Limit:       10,
		Offset:      int32(req.ItemOffset),
	})
	if err != nil {
		log.Printf("Error at  ListEventDateItems in store.ListEventDateTime err: %v, user: %v\n", err, user.ID)
		res := ListEventDateItem{
			List:        empty,
			IsEmpty:     true,
			OptionData:  uhmData,
			ItemOffset:  0,
			OnLastIndex: onLastIndex,
		}
		ctx.JSON(http.StatusNoContent, res)
	}
	var res ListEventDateItem
	var resData []EventDateItem
	for i := 0; i < len(eventDateTimes); i++ {
		var data EventDateItem
		detailIsEmpty := false
		eventDateDetail, err := server.store.GetEventDateDetail(ctx, eventDateTimes[i].ID)
		// We check to error
		if err != nil {
			log.Printf("Error at  ListEventDateItems in .GetEventDateDetail err: %v, user: %v\n", err, user.ID)
			detailIsEmpty = true
		} //: IF End
		ticketCount, err := server.store.GetEventDateTicketCount(ctx, eventDateTimes[i].ID)
		// We check to error
		if err != nil {
			ticketCount = 0
			log.Printf("Error at  ListEventDateItems in ListEventDateTime err: %v, user: %v\n", err, user.ID)
		}
		if detailIsEmpty {
			data = EventDateItem{
				ID:               tools.UuidToString(eventDateTimes[i].ID),
				Name:             eventDateTimes[i].Name,
				StartTime:        "",
				EndTime:          "",
				StartDate:        tools.ConvertDateOnlyToString(eventDateTimes[i].StartDate),
				Status:           eventDateTimes[i].Status,
				EndDate:          tools.ConvertDateOnlyToString(eventDateTimes[i].EndDate),
				Tickets:          int(ticketCount),
				Note:             "",
				TimeZone:         "",
				Type:             eventDateTimes[i].Type,
				EventDates:       eventDateTimes[i].EventDates,
				NeedBands:        eventDateTimes[i].NeedBands,
				NeedTickets:      eventDateTimes[i].NeedTickets,
				AbsorbBandCharge: eventDateTimes[i].AbsorbBandCharge,
			}

		} else {
			data = EventDateItem{
				ID:               tools.UuidToString(eventDateTimes[i].ID),
				Name:             eventDateTimes[i].Name,
				StartTime:        eventDateDetail.StartTime,
				EndTime:          eventDateDetail.EndTime,
				Status:           eventDateTimes[i].Status,
				StartDate:        tools.ConvertDateOnlyToString(eventDateTimes[i].StartDate),
				EndDate:          tools.ConvertDateOnlyToString(eventDateTimes[i].EndDate),
				Tickets:          int(ticketCount),
				Note:             eventDateTimes[i].Note,
				TimeZone:         eventDateDetail.TimeZone,
				Type:             eventDateTimes[i].Type,
				EventDates:       eventDateTimes[i].EventDates,
				NeedBands:        eventDateTimes[i].NeedBands,
				NeedTickets:      eventDateTimes[i].NeedTickets,
				AbsorbBandCharge: eventDateTimes[i].AbsorbBandCharge,
			}
		}
		resData = append(resData, data)
	}

	if count <= int64(req.ItemOffset+len(eventDateTimes)) {
		onLastIndex = true
	}
	res = ListEventDateItem{
		List:        resData,
		IsEmpty:     false,
		OptionData:  uhmData,
		ItemOffset:  req.ItemOffset + len(eventDateTimes),
		OnLastIndex: onLastIndex,
	}

	ctx.JSON(http.StatusOK, res)
}

func (server *Server) ListEventDateNormalItems(ctx *gin.Context) {
	var req EventDateParams
	var onLastIndex = false
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at  EventDateParams in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.ItemOffset)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Println(req)
	eventID, err := tools.StringToUuid(req.EventInfoID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, _, option, _, _, err := HandleGetCompleteOptionEditEventDateTimes(eventID, ctx, server, true)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	// We want to verify access to this resource
	_, err = server.store.GetEventInfoByOption(ctx, db.GetEventInfoByOptionParams{
		OptionID:   option.ID,
		ID:         user.ID,
		IsComplete: true,
	})
	if err != nil {
		log.Printf("Error at  ListEventDateNormalItems in GetEventDateTimeByOption err: %v, user: %v\n", err, user.ID)
		err = fmt.Errorf("you don't have access to this resource")
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}
	count, err := server.store.GetEventDateTimeCount(ctx, eventID)
	if err != nil {
		log.Printf("Error at  ListEventDateNormalItems in GetEventDateTimeCount err: %v, user: %v\n", err, user.ID)
		res := "none"
		ctx.JSON(http.StatusNoContent, res)
		return
	}
	if count <= int64(req.ItemOffset) || count == 0 {
		res := "none"
		ctx.JSON(http.StatusNoContent, res)
		return
	}
	eventDateTimes, err := server.store.ListEventDateTime(ctx, db.ListEventDateTimeParams{
		EventInfoID: eventID,
		Limit:       5,
		Offset:      int32(req.ItemOffset),
	})
	if err != nil {
		res := "none"
		ctx.JSON(http.StatusNoContent, res)
	}
	var resData []EventDateNormalItem
	for i := 0; i < len(eventDateTimes); i++ {
		data := EventDateNormalItem{
			ID:         tools.UuidToString(eventDateTimes[i].ID),
			Name:       eventDateTimes[i].Name,
			StartDate:  tools.ConvertDateOnlyToString(eventDateTimes[i].StartDate),
			EndDate:    tools.ConvertDateOnlyToString(eventDateTimes[i].EndDate),
			Type:       eventDateTimes[i].Type,
			EventDates: eventDateTimes[i].EventDates,
		}
		resData = append(resData, data)
	}
	if count <= int64(req.ItemOffset+len(eventDateTimes)) {
		onLastIndex = true
	}
	res := ListEventDateNormalItem{
		List:        resData,
		ItemOffset:  req.ItemOffset + len(eventDateTimes),
		OnLastIndex: onLastIndex,
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) UpdateEventDateTimeNote(ctx *gin.Context) {
	var req UpdateEventDateTimeNoteParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at UpdateEventDateTimeNoteParams in ShouldBindJSON: %v, EventInfoID: %v \n", err.Error(), req.EventDateTimeID)
		err = fmt.Errorf("the note cannot be empty")
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
	user, _, option, isCoHost, userCoHost, err := HandleGetCompleteOptionEditEventDateTimes(eventID, ctx, server, true)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	_, err = server.store.GetEventDateTimeByOption(ctx, db.GetEventDateTimeByOptionParams{
		ID:         eventDateTimeID,
		ID_2:       user.ID,
		ID_3:       option.ID,
		IsComplete: true,
	})
	if err != nil {
		log.Printf("There an error at UpdateEventDateTimeNote at GetEventInfoByOption: %v, EventInfoID: %v, userID: %v \n", err.Error(), eventID, user.ID)
		err = fmt.Errorf("could not create your date, please try again using the format on the app")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	eventDT, err := server.store.UpdateEventDateTime(ctx, db.UpdateEventDateTimeParams{
		Note: pgtype.Text{
			String: req.Note,
			Valid:  true,
		},
		ID: eventID,
	})
	if err != nil {
		log.Printf("There an error at UpdateEventDateTimeNote at UpdateEventDateTime: %v, EventInfoID: %v, userID: %v \n", err.Error(), eventID, user.ID)
		err = fmt.Errorf("could not update your notes for this event date")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	res := UpdateEventDateTimeNoteParams{
		EventDateTimeID: tools.UuidToString(eventDT.ID),
		Note:            eventDT.Note,
	}
	if isCoHost {
		HandleCoHostUpdateMsg(ctx, server, userCoHost, user, option, "UpdateEventDateTimeNote", "event date time note", "update event date time note")
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) RemoveEventDateTime(ctx *gin.Context) {
	var req RemoveEventDateTimeParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at RemoveEventDateTimeParams in ShouldBindJSON: %v, EventInfoID: %v \n", err.Error(), req.EventDateTimeID)
		err = fmt.Errorf("please enter all the details for this ticket")
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
	user, _, option, isCoHost, userCoHost, err := HandleGetCompleteOptionEditEventDateTimes(eventID, ctx, server, true)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	// We want make sure the EventDateTimeID belongs to this user
	_, err = server.store.GetEventDateTimeByOption(ctx, db.GetEventDateTimeByOptionParams{
		ID:         eventDateTimeID,
		ID_2:       user.ID,
		ID_3:       option.ID,
		IsComplete: true,
	})
	if err != nil {
		log.Printf("There an error at RemoveEventDateTime at GetEventDateTimeByOption: %v, EventInfoID: %v, userID: %v \n", err.Error(), eventID, user.ID)
		err = fmt.Errorf("could not create ticket for this event date, please try again")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// We check if the eventDateTicket has an active booking
	isBooked := EventDateIsBooked(ctx, server, eventDateTimeID, "RemoveEventDateTime", tools.UuidToString(user.ID))
	log.Println("isBooked any ", isBooked)
	if isBooked {
		err = fmt.Errorf("this event date currently has an active booking so it can't be deleted")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	// We want to know if there were any ticket booked even if it was not complete
	isBookedAny := EventDateIsBookedAny(ctx, server, eventDateTimeID, "RemoveEventDateTime", tools.UuidToString(user.ID))
	log.Println("isBookend any ", isBookedAny)
	if isBookedAny {
		log.Println("isBookend any here ", isBookedAny)
		_, err = server.store.UpdateEventDateTimeActive(ctx, db.UpdateEventDateTimeActiveParams{
			IsActive: pgtype.Bool{
				Bool:  false,
				Valid: true,
			},
			ID: eventDateTimeID,
		})
		if err != nil {
			log.Printf("There an error at server.store.UpdateEventDateTimeActive at GetEventDateTimeByOption: %v, EventInfoID: %v, userID: %v \n", err.Error(), eventDateTimeID, user.ID)
			err = fmt.Errorf("could not remove ticket for this event date, please try again")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
	} else {
		log.Println("at else anybooked ", isBookedAny)
		// We need to remove all photos in check in steps
		photos, err := server.store.ListEventCheckInStepPhotos(ctx, eventDateTimeID)
		log.Println("isBookend any photos ", photos)
		if err != nil || len(photos) == 0 {
			if err != nil {
				log.Printf("There an error at store.RemoveEventDateDetail at ListEventCheckInStepPhotos: %v, EventInfoID: %v, userID: %v \n", err.Error(), eventDateTimeID, user.ID)
			}
			err = server.store.RemoveAllEventCheckInStep(ctx, eventDateTimeID)
			if err != nil {
				log.Printf("There an error at store.RemoveEventDateDetail at RemoveAllEventCheckInStep: %v, EventInfoID: %v, userID: %v \n", err.Error(), eventDateTimeID, user.ID)
			}
		} else {
			log.Println("isBookend any photos 66 ", photos)
			err = server.store.RemoveAllEventCheckInStep(ctx, eventDateTimeID)
			if err != nil {
				log.Printf("There an error at store.RemoveEventDateDetail at RemoveAllEventCheckInStep: %v, EventInfoID: %v, userID: %v \n", err.Error(), eventDateTimeID, user.ID)
			}
			// Lets remove photos from firebase
			for _, p := range photos {
				if !tools.ServerStringEmpty(p) {
					err = RemoveFirebasePhoto(server, ctx, p)
					if err != nil {
						log.Printf("There an error at RemoveEventCheckInStep at RemoveFirebasePhoto: %v, eventDateTimeID: %v, userID: %v, photo: %v \n", err.Error(), eventDateTimeID, user.ID, p)
						continue
					}
				}
			}
		}
		// Lets remove eventDateDetail
		err = server.store.RemoveEventDateDetail(ctx, eventDateTimeID)
		log.Println("RemoveEventDateDetail here ")
		if err != nil {
			log.Printf("There an error at store.RemoveEventDateDetail at RemoveEventDateDetail: %v, EventInfoID: %v, userID: %v \n", err.Error(), eventDateTimeID, user.ID)
		}
		// Lets remove location
		err = server.store.RemoveEventDateLocation(ctx, eventDateTimeID)
		if err != nil {
			log.Printf("There an error at RemoveEventDateLocation at GetEventDateTimeByOption: %v, EventInfoID: %v, userID: %v \n", err.Error(), eventDateTimeID, user.ID)
		}
		log.Println("RemoveEventDateLocation here ")
		err = server.store.RemoveAllEventDateTicket(ctx, eventDateTimeID)
		if err != nil {
			log.Printf("There an error at RemoveAllEventDateTicket at GetEventDateTimeByOption: %v, EventInfoID: %v, userID: %v \n", err.Error(), eventDateTimeID, user.ID)
		}
		err = server.store.RemoveEventDatePublish(ctx, eventDateTimeID)
		if err != nil {
			log.Printf("There an error at RemoveAllEventDateTicket at .RemoveEventDatePublish: %v, EventInfoID: %v, userID: %v \n", err.Error(), eventDateTimeID, user.ID)
		}
		err = server.store.RemoveAllEventDatePrivateAudience(ctx, eventDateTimeID)
		if err != nil {
			log.Printf("There an error at RemoveAllEventDateTicket at RemoveAllEventDatePrivateAudience: %v, EventInfoID: %v, userID: %v \n", err.Error(), eventDateTimeID, user.ID)
		}
		err = server.store.RemoveEventDateTime(ctx, eventDateTimeID)
		if err != nil {
			log.Printf("There an error at RemoveEventDateTime at GetEventDateTimeByOption: %v, EventInfoID: %v, userID: %v \n", err.Error(), eventDateTimeID, user.ID)
			err = fmt.Errorf("could not remove ticket for this event date, please try again")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
	}
	res := RemoveEventDateTimeRes{
		EventDateTimeID: req.EventDateTimeID,
	}
	if isCoHost {
		HandleCoHostUpdateMsg(ctx, server, userCoHost, user, option, "RemoveEventDateTime", "event date removed", "remove event date")
	}
	ctx.JSON(http.StatusOK, res)
}
