package api

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/makuo12/ghost_server/constants"
	db "github.com/makuo12/ghost_server/db/sqlc"
	"github.com/makuo12/ghost_server/tools"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func (server *Server) GetEventDateIsBooked(ctx *gin.Context) {
	var req GetEventDateDetailParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at GetEventDateIsBookedParams in ShouldBindJSON: %v, EventDateTimeID: %v \n", err.Error(), req.EventDateTimeID)
		err = fmt.Errorf("please make sure yuo select at least a start date")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Println(req)
	eventInfoID, err := tools.StringToUuid(req.EventInfoID)
	if err != nil {
		log.Printf("Error GetEventDateIsBooked at tools.StringToUuid: %v, EventInfoID: %v \n", err.Error(), req.EventInfoID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	requestID, err := tools.StringToUuid(req.EventDateTimeID)
	if err != nil {
		log.Printf("Error at tools.StringToUuid: %v, EventDateTimeID: %v \n", err.Error(), req.EventDateTimeID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, _, option, _, _, err := HandleGetCompleteOptionEditEventDateTimes(eventInfoID, ctx, server, true)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	_, err = server.store.GetEventDateTimeByOption(ctx, db.GetEventDateTimeByOptionParams{
		ID:         requestID,
		ID_2:       user.ID,
		ID_3:       option.ID,
		IsComplete: true,
	})
	if err != nil {
		log.Printf("There an error at GetEventDateIsBooked at GetEventDateTimeByOption: %v, EventDateTimeID: %v, userID: %v \n", err.Error(), requestID, user.ID)
		err = fmt.Errorf("error occurred")
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}
	isBooked := EventDateIsBooked(ctx, server, requestID, "GetEventDateIsBooked", tools.UuidToString(user.ID))

	res := GetEventDateIsBookedRes{
		IsBooked: isBooked,
	}
	log.Println(res)
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) ListEventDateBooked(ctx *gin.Context) {
	var req GetEventDateDetailParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at ListEventDateBooked in ShouldBindJSON: %v, EventDateTimeID: %v \n", err.Error(), req.EventDateTimeID)
		err = fmt.Errorf("please make sure yuo select at least a start date")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Println(req)
	eventInfoID, err := tools.StringToUuid(req.EventInfoID)
	if err != nil {
		log.Printf("Error ListEventDateBooked at tools.StringToUuid: %v, EventInfoID: %v \n", err.Error(), req.EventInfoID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	requestID, err := tools.StringToUuid(req.EventDateTimeID)
	if err != nil {
		log.Printf("Error ListEventDateBooked at tools.StringToUuid: %v, EventDateTimeID: %v \n", err.Error(), req.EventDateTimeID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, _, option, isCoHost, userCoHost, err := HandleGetCompleteOptionEditEventDateTimes(eventInfoID, ctx, server, true)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	_, err = server.store.GetEventDateTimeByOption(ctx, db.GetEventDateTimeByOptionParams{
		ID:         requestID,
		ID_2:       user.ID,
		ID_3:       option.ID,
		IsComplete: true,
	})
	if err != nil {
		log.Printf("There an error at ListEventDateBooked at GetEventDateTimeByOption: %v, EventDateTimeID: %v, userID: %v \n", err.Error(), requestID, user.ID)
		err = fmt.Errorf("error occurred")
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}
	dates := GetEventDatesBooked(ctx, server, requestID, user, "ListEventDateBooked")
	if len(dates) == 0 {
		if isCoHost {
			HandleCoHostUpdateMsg(ctx, server, userCoHost, user, option, "ListEventDateBooked", "event date status", "update event date status")
		}
		ctx.JSON(http.StatusNoContent, "none")
		return
	}
	var resData []DateEventBookedItem
	for _, d := range dates {
		data := DateEventBookedItem{
			FakeID:    tools.UuidToString(uuid.New()),
			StartDate: tools.ConvertDateOnlyToString(d.StartDate),
			EndDate:   tools.ConvertDateOnlyToString(d.EndDate),
			Booked:    true,
			Count:     int(d.ItemCount),
		}
		resData = append(resData, data)
	}
	res := ListDateEventBookedRes{
		List: resData,
	}
	log.Println(res)
	if isCoHost {
		HandleCoHostUpdateMsg(ctx, server, userCoHost, user, option, "ListEventDateBooked", "event date status", "update event date status")
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) OptionDateIsBooked(ctx *gin.Context) {
	var req OptionSelectedDateParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at OptionDateIsBooked at  OptionDateParams in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionDateTimeID)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Println(req)
	requestID, err := tools.StringToUuid(req.OptionID)
	if err != nil {
		log.Printf("Error at OptionDateIsBooked at tools.StringToUuid: %v, OptionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, _, option, _, _, err := HandleGetCompleteOptionCalender(requestID, ctx, server, true)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	isBooked := OptionDateIsBooked(ctx, server, option.OptionUserID, "OptionDateIsBooked", tools.UuidToString(user.ID))

	res := GetOptionDateIsBookedRes{
		IsBooked: isBooked,
	}
	log.Println(res)
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) UpdateEventDatesBooked(ctx *gin.Context) {
	var req UpdateHostEventDateBookedParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at UpdateEventDatesBooked  OptionDateParams in ShouldBindJSON: %v, eventDateTimeID: %v \n", err.Error(), req.EventDateID)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Println("req change date booked ", req)
	eventDateID, err := tools.StringToUuid(req.EventDateID)
	if err != nil {
		log.Printf("FuncName: %v. There an error at HandleEventHostCancel at tools.StringToUuid: %v, startDate: %v, eventDateID: %v \n", "UpdateEventDatesBooked", err.Error(), req.StartDate, req.EventDateID)
		err = fmt.Errorf("this event date cannot be found")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	eventID, err := tools.StringToUuid(req.EventID)
	if err != nil {
		log.Printf("FuncName: %v. There an error at HandleEventHostCancel at tools.StringToUuid(req.EventID): %v, startDate: %v, eventDateID: %v \n", "UpdateEventDatesBooked", err.Error(), req.StartDate, req.EventID)
		err = fmt.Errorf("this event date cannot be found")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, _, option, isCoHost, userCoHost, err := HandleGetCompleteOptionReservation(eventID, ctx, server, true)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// If new dates and current dates match we don't want to update because there is not point
	if tools.DatesMatchString(req.StartDate, req.EndDate, req.NewStartDate, req.NewEndDate) {
		err = fmt.Errorf("error occurred, your current dates and new dates match")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	startDate, err := tools.ConvertDateOnlyStringToDate(req.StartDate)
	if err != nil {
		err = fmt.Errorf("this date for the event is not in the right format")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	endDate, err := tools.ConvertDateOnlyStringToDate(req.EndDate)
	if err != nil {
		err = fmt.Errorf("this date for the event is not in the right format")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	newStartDate, err := tools.ConvertDateOnlyStringToDate(req.NewStartDate)
	if err != nil {
		err = fmt.Errorf("this date for the event is not in the right format")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	newEndDate, err := tools.ConvertDateOnlyStringToDate(req.NewEndDate)
	if err != nil {
		err = fmt.Errorf("this date for the event is not in the right format")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	// If an event is about to end you cannot also cancel it
	if time.Now().After(endDate.Add(time.Hour * -5)) {
		err = fmt.Errorf("this date cannot be changed")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if newStartDate.Before(time.Now().Add(time.Hour)) || tools.ConvertTimeToString(newStartDate) == tools.ConvertTimeToString(time.Now().Add(time.Hour)) {
		err = fmt.Errorf("event date must be change to a future date")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	eventDate, err := server.store.GetEventDateTimeByUID(ctx, db.GetEventDateTimeByUIDParams{
		EventDateTimeID: eventDateID,
		UID:             user.ID,
		EventInfoID:     eventID,
	})
	if err != nil {
		log.Printf("FuncName: %v. There an error at UpdateEventDatesBooked at .CreateRefund: %v, startDate: %v, userID: %v \n", "UpdateEventDatesBooked", err.Error(), req.StartDate, user.ID)
		err = fmt.Errorf("this event date was not found")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	// Lets check if the particular event has data in redis that needs to be cancelled or updated
	keyCancel := fmt.Sprintf("%v&%v&%v&%v", eventDate.EventDateTimeID, eventDate.Type, req.StartDate, "cancel")
	log.Println("keyCancel ", keyCancel)
	result, err := RedisClient.HExists(RedisContext, keyCancel, constants.REFERENCE).Result()
	if err != nil {
		log.Printf("FuncName: %v. Cancel There an error at RedisClient.SMembers: %v, startDate: %v, userID: %v \n", "UpdateEventDatesBooked", err.Error(), req.StartDate, user.ID)
		err = fmt.Errorf("an error occurred while checking details")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if result {
		err = fmt.Errorf("there is already a cancellation process in progress for this event, please try again is the cancellation fails")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	keyUpdate := fmt.Sprintf("%v&%v&%v&%v", eventDate.EventDateTimeID, eventDate.Type, req.StartDate, "update")
	log.Println("keyUpdate ", keyUpdate)
	result, err = RedisClient.HExists(RedisContext, keyUpdate, constants.REFERENCE).Result()
	if err != nil {
		log.Printf("FuncName: %v. Update There an error at RedisClient.SMembers RedisClient.SMembers(RedisContext, keyUpdate: %v, startDate: %v, userID: %v \n", "UpdateEventDatesBooked", err.Error(), req.StartDate, user.ID)
		err = fmt.Errorf("an error occurred while checking details")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if result {
		err = fmt.Errorf("there is already an update process in progress for this event, please try again if the process fails")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	// Lets get all the current charge ticket ids for this event
	chargeIDs, err := server.store.ListChargeTicketReferenceIDByStartDate(ctx, db.ListChargeTicketReferenceIDByStartDateParams{
		Date:        startDate,
		Cancelled:   false,
		IsComplete:  true,
		EventDateID: eventDate.EventDateTimeID,
	})
	if err != nil {
		log.Printf("FuncName: %v. There an error at UpdateEventDatesBooked at .ListChargeTicketReferenceIDByStartDate: %v, startDate: %v, userID: %v \n", "UpdateEventDatesBooked", err.Error(), req.StartDate, user.ID)
		err = fmt.Errorf("unable to find booked tickets for this event date")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	switch eventDate.Type {
	case "single":
		// We need to make sure the startDate and endDate match before update
		if !tools.DatesMatch(req.StartDate, req.EndDate, eventDate.StartDate, eventDate.EndDate) {
			err = fmt.Errorf("current dates not found, please try again later")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		_, err = server.store.UpdateEventDateTime(ctx, db.UpdateEventDateTimeParams{
			StartDate: pgtype.Date{
				Time:  newStartDate,
				Valid: true,
			},
			EndDate: pgtype.Date{
				Time:  newEndDate,
				Valid: true,
			},
			ID:   eventDate.EventDateTimeID,
			Type: "single",
		})
		if err != nil {
			log.Printf("FuncName: %v. There an error at UpdateEventDatesBooked at UpdateEventDateTime: %v, startDate: %v, userID: %v \n", "UpdateEventDatesBooked", err.Error(), req.StartDate, user.ID)
			err = fmt.Errorf("unable to cancel this single event")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
	case "recurring":
		newDates, errRecur := tools.UpdateRecurDate(req.NewStartDate, req.StartDate, eventDate.EventDates)
		if errRecur != nil {
			err = errRecur
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		_, err = server.store.UpdateEventDateTimeDates(ctx, db.UpdateEventDateTimeDatesParams{
			EventDates: newDates,
			ID:         eventDate.EventDateTimeID,
		})
		if err != nil {
			log.Printf("FuncName: %v. There an error at UpdateEventDatesBooked at UpdateEventDateTimeDates: %v, startDate: %v, userID: %v \n", "UpdateEventDatesBooked", err.Error(), req.StartDate, user.ID)
			err = fmt.Errorf("unable to cancel this recurring event")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
	}
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if len(chargeIDs) == 0 {
		log.Println("chargeID empty ", chargeIDs)
		res := UpdateHostEventDateBookedRes{
			StartDate:   req.NewStartDate,
			EndDate:     req.NewEndDate,
			EventDateID: req.EventDateID,
		}
		ctx.JSON(http.StatusOK, res)
	}
	uniqueID := fmt.Sprintf("%v&%v&%v&%v", tools.UuidToString(eventDate.EventDateTimeID), eventDate.Type, req.NewStartDate, "update")
	chargeReference := tools.UuidToString(uuid.New())
	data := []string{
		constants.START_DATE,
		req.NewStartDate,
		constants.END_DATE,
		req.NewEndDate,
		constants.MESSAGE,
		req.Message,
		constants.REASON_ONE,
		req.ReasonOne,
		constants.REFERENCE,
		chargeReference,
	}

	err = RedisClient.SAdd(RedisContext, constants.CHARGE_TICKET_ID_UPDATE, uniqueID).Err()
	if err != nil {
		log.Printf("FuncName: %v. There an error at UpdateEventDatesBooked at RedisClient.SAdd: %v, startDate: %v, userID: %v, chargeIDs: %v \n", "UpdateEventDatesBooked", err.Error(), req.StartDate, user.ID, chargeIDs)
		err = nil
	}
	err = RedisClient.HSet(RedisContext, uniqueID, data).Err()
	if err != nil {
		log.Printf("FuncName: %v. There an error at UpdateEventDatesBooked at RedisClient.HSet: %v, startDate: %v, userID: %v, chargeIDs: %v \n", "UpdateEventDatesBooked", err.Error(), req.StartDate, user.ID, chargeIDs)
		err = nil
	}
	// We want to store the chargeIDs in redis so that we can later give full refund and set it to cancel
	err = RedisClient.SAdd(RedisContext, chargeReference, tools.ListUuidToString(chargeIDs)).Err()
	if err != nil {
		log.Printf("FuncName: %v. There an error at UpdateEventDatesBooked at RedisClient.SAdd: %v, startDate: %v, userID: %v, chargeIDs: %v \n", "UpdateEventDatesBooked", err.Error(), req.StartDate, user.ID, chargeIDs)
		err = nil
	}
	if isCoHost {
		HandleCoHostUpdateMsg(ctx, server, userCoHost, user, option, "ListEventDateBooked", "event date status", "update event date status")
	}
	res := UpdateHostEventDateBookedRes{
		StartDate:   req.NewStartDate,
		EndDate:     req.NewEndDate,
		EventDateID: req.EventDateID,
	}
	ctx.JSON(http.StatusOK, res)
}
