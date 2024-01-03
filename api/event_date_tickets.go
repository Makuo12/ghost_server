package api

import (
	"errors"
	db "flex_server/db/sqlc"
	"flex_server/tools"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

func (server *Server) CreateEventDateTicket(ctx *gin.Context) {
	var req CreateEventDateTicketParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at CreateEventDateTicketParams in ShouldBindJSON: %v, eventID: %v \n", err.Error(), req.EventDateTimeID)
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
		log.Printf("There an error at CreateUpdateEventDateTicket at GetEventDateTimeByOption: %v, EventDateTimeID: %v, userID: %v \n", err.Error(), eventDateTimeID, user.ID)
		err = fmt.Errorf("could not create ticket for this event date, please try again")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
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
	startTime, err := tools.ConvertStringToTimeOnly(req.StartTime)
	log.Println("Show start time", startTime)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	endTime, err := tools.ConvertStringToTimeOnly(req.EndTime)
	log.Println("Show end time", endTime)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if req.Type == "paid" && tools.MoneyStringToInt(req.Price) < 1 {
		err = fmt.Errorf("ticket price is too low for a paid event")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}
	// Before we create the ticket lets know if there exist a ticket of that grade
	_, err = server.store.GetEventDateTicketByGrade(ctx, db.GetEventDateTicketByGradeParams{
		EventDateTimeID: eventDateTimeID,
		Level:           req.Level,
	})
	if err == nil {
		// If error is equal to nil we want to send a error saying that only one grade of this ticket can be created
		err = errors.New("the grade level of this ticket has already been created please try another grade type that you haven't created before for this event")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	} else {
		log.Printf("There an error at CreateUpdateEventDateTicket at GetEventDateTimeByOption: %v, EventDateTimeID: %v, userID: %v \n", err.Error(), eventDateTimeID, user.ID)
	}
	eventDateTicket, err := server.store.CreateEventDateTicket(ctx, db.CreateEventDateTicketParams{
		EventDateTimeID: eventDateTimeID,
		StartDate:       startDate,
		EndDate:         endDate,
		StartTime:       startTime,
		EndTime:         endTime,
		Name:            req.Name,
		Capacity:        int32(req.Capacity),
		Type:            req.Type,
		Level:           req.Level,
		Price:           tools.MoneyStringToInt(req.Price),
		AbsorbFees:      req.AbsorbFees,
		Description:     req.Description,
		TicketType:      req.TicketType,
		NumOfSeats:      int32(req.NumOfSeats),
		FreeRefreshment: req.FreeRefreshment,
	})
	if err != nil {
		log.Printf("There an error at CreateEventDateTicket at CreateEventDateTicket: %v, eventID: %v, userID: %v \n", err.Error(), eventDateTimeID, user.ID)

		err = fmt.Errorf("could not create ticket for this event date, please try again")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	ticketCount, err := server.store.GetEventDateTicketCount(ctx, eventDateTicket.EventDateTimeID)
	if err != nil {
		log.Printf("There an error at CreateEventDateTicket at GetEventDateTicketCount: %v, eventID: %v, userID: %v \n", err.Error(), eventDateTimeID, user.ID)
		ticketCount = 0
	}
	res := CreateEventDateTicketRes{
		EventDateTimeID: tools.UuidToString(eventDateTicket.EventDateTimeID),
		ID:              tools.UuidToString(eventDateTicket.ID),
		StartDate:       tools.ConvertDateOnlyToString(eventDateTicket.StartDate),
		EndDate:         tools.ConvertDateOnlyToString(eventDateTicket.EndDate),
		StartTime:       tools.ConvertTimeOnlyToString(eventDateTicket.StartTime),
		EndTime:         tools.ConvertTimeOnlyToString(eventDateTicket.EndTime),
		Name:            eventDateTicket.Name,
		Price:           req.Price,
		AbsorbFees:      eventDateTicket.AbsorbFees,
		Description:     eventDateTicket.Description,
		Capacity:        int(eventDateTicket.Capacity),
		Type:            eventDateTicket.Type,
		Level:           eventDateTicket.Level,
		TicketType:      eventDateTicket.TicketType,
		NumOfSeats:      int(eventDateTicket.NumOfSeats),
		FreeRefreshment: eventDateTicket.FreeRefreshment,
		TicketCount:     int(ticketCount),
	}
	if isCoHost {
		HandleCoHostUpdateMsg(ctx, server, userCoHost, user, option, "CreateEventDateTicket", "event date ticket", "create event date ticket")
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) ListEventDateTicket(ctx *gin.Context) {
	var req ListEventDateTicketParams
	var onLastIndex = false
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at CreateEventDateTicketParams in ShouldBindJSON: %v, eventID: %v \n", err.Error(), req.EventDateTimeID)
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
	user, _, option, _, _, err := HandleGetCompleteOptionEditEventDateTimes(eventID, ctx, server, true)
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
		log.Printf("There an error at ListEventDateTicket at GetEventDateTimeByOption: %v, eventID: %v, userID: %v \n", err.Error(), eventDateTimeID, user.ID)
		err = fmt.Errorf("could not create ticket for this event date, please try again")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	count, err := server.store.GetEventDateTicketCount(ctx, eventDateTimeID)
	if err != nil {
		log.Printf("Error at  ListEventDateTicket in GetEventDateTimeCount err: %v, user: %v\n", err, user.ID)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if count <= int64(req.OptionOffset) {
		res := "none"
		ctx.JSON(http.StatusNoContent, res)
		return
	}
	eventDateTickets, err := server.store.ListEventDateTicketOffset(ctx, db.ListEventDateTicketOffsetParams{
		EventDateTimeID: eventDateTimeID,
		Limit:           5,
		Offset:          int32(req.OptionOffset),
	})
	if err != nil {
		log.Printf("Error at  ListEventDateTicket in ListEventDateTicket err: %v, user: %v\n", err, user.ID)
		res := "none"
		ctx.JSON(http.StatusNoContent, res)
		return
	}
	ticketCount, err := server.store.GetEventDateTicketCount(ctx, eventDateTimeID)
	if err != nil {
		log.Printf("Error at  ListEventDateTicket in GetEventDateTicketCount err: %v, user: %v\n", err, user.ID)
		ticketCount = 0
	}
	var grades []string
	var res ListEventDateTicketRes
	var resData []CreateEventDateTicketRes
	for i := 0; i < len(eventDateTickets); i++ {
		data := CreateEventDateTicketRes{
			ID:              tools.UuidToString(eventDateTickets[i].ID),
			EventDateTimeID: tools.UuidToString(eventDateTickets[i].EventDateTimeID),
			StartDate:       tools.ConvertDateOnlyToString(eventDateTickets[i].StartDate),
			EndDate:         tools.ConvertDateOnlyToString(eventDateTickets[i].EndDate),
			StartTime:       tools.ConvertTimeOnlyToString(eventDateTickets[i].StartTime),
			EndTime:         tools.ConvertTimeOnlyToString(eventDateTickets[i].EndTime),
			Name:            eventDateTickets[i].Name,
			Capacity:        int(eventDateTickets[i].Capacity),
			Price:           tools.IntToMoneyString(eventDateTickets[i].Price),
			Type:            eventDateTickets[i].Type,
			Level:           eventDateTickets[i].Level,
			TicketType:      eventDateTickets[i].TicketType,
			NumOfSeats:      int(eventDateTickets[i].NumOfSeats),
			FreeRefreshment: eventDateTickets[i].FreeRefreshment,
			TicketCount:     int(ticketCount),
			Description:     eventDateTickets[i].Description,
		}
		grades = append(grades, eventDateTickets[i].Level)
		resData = append(resData, data)
	}
	if count <= int64(req.OptionOffset+len(eventDateTickets)) {
		onLastIndex = true
	}
	res = ListEventDateTicketRes{
		List:         resData,
		OptionOffset: req.OptionOffset + len(eventDateTickets),
		OnLastIndex:  onLastIndex,
		Grades:       grades,
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) UpdateEventDateTicket(ctx *gin.Context) {
	var req UpdateEventDateTicketParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at UpdateEventDateTicketParams in ShouldBindJSON: %v, eventID: %v \n", err.Error(), req.EventDateTimeID)
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
	ticketID, err := tools.StringToUuid(req.TicketID)
	if err != nil {
		log.Printf("Error at tools.StringToUuid: %v, TicketID: %v \n", err.Error(), req.TicketID)
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
		log.Printf("There an error at CreateUpdateEventDateTicket at GetEventDateTimeByOption: %v, eventID: %v, userID: %v \n", err.Error(), eventDateTimeID, user.ID)
		err = fmt.Errorf("could not create ticket for this event date, please try again")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

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
	startTime, err := tools.ConvertStringToTimeOnly(req.StartTime)
	log.Println("Show start time", startTime)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	endTime, err := tools.ConvertStringToTimeOnly(req.EndTime)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	capacityOk, err := TicketQuantity(server, ctx, req.Capacity, ticketID, eventDateTimeID)
	if err != nil || !capacityOk {
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))

		} else {
			err = fmt.Errorf("the ticket capacity does not meet the requirements")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
		}
		return
	}
	if req.Type == "paid" && tools.MoneyStringToInt(req.Price) < 1 {
		err = fmt.Errorf("ticket price is too low for a paid event")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}
	eventDateTicket, err := server.store.UpdateEventDateTicketTwo(ctx, db.UpdateEventDateTicketTwoParams{
		EventDateTimeID: eventDateTimeID,
		StartDate:       startDate,
		EndDate:         endDate,
		StartTime:       startTime,
		EndTime:         endTime,
		Name:            req.Name,
		Capacity:        int32(req.Capacity),
		Type:            req.Type,
		Level:           req.Level,
		Price:           tools.MoneyStringToInt(req.Price),
		AbsorbFees:      req.AbsorbFees,
		Description:     req.Description,
		TicketType:      req.TicketType,
		NumOfSeats:      int32(req.NumOfSeats),
		FreeRefreshment: req.FreeRefreshment,
		ID:              ticketID,
	})
	if err != nil {
		log.Printf("There an error at UpdateEventDateTicket at UpdateEventDateTicket: %v, eventID: %v, userID: %v \n", err.Error(), eventDateTimeID, user.ID)

		err = fmt.Errorf("could not create ticket for this event date, please try again")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	ticketCount, err := server.store.GetEventDateTicketCount(ctx, eventDateTicket.EventDateTimeID)
	if err != nil {
		log.Printf("There an error at UpdateEventDateTicket at GetEventDateTicketCount: %v, eventID: %v, userID: %v \n", err.Error(), eventDateTimeID, user.ID)
		ticketCount = 0
	}
	res := CreateEventDateTicketRes{
		EventDateTimeID: tools.UuidToString(eventDateTicket.EventDateTimeID),
		ID:              tools.UuidToString(eventDateTicket.ID),
		StartDate:       tools.ConvertDateOnlyToString(eventDateTicket.StartDate),
		EndDate:         tools.ConvertDateOnlyToString(eventDateTicket.EndDate),
		StartTime:       tools.ConvertTimeOnlyToString(eventDateTicket.StartTime),
		EndTime:         tools.ConvertTimeOnlyToString(eventDateTicket.EndTime),
		Name:            eventDateTicket.Name,
		Price:           req.Price,
		AbsorbFees:      eventDateTicket.AbsorbFees,
		Description:     eventDateTicket.Description,
		Capacity:        int(eventDateTicket.Capacity),
		Type:            eventDateTicket.Type,
		Level:           eventDateTicket.Level,
		TicketType:      eventDateTicket.TicketType,
		NumOfSeats:      int(eventDateTicket.NumOfSeats),
		FreeRefreshment: eventDateTicket.FreeRefreshment,
		TicketCount:     int(ticketCount),
	}
	if isCoHost {
		HandleCoHostUpdateMsg(ctx, server, userCoHost, user, option, "UpdateEventDateTicket", "event date ticket", "create event date ticket")
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) RemoveEventDateTicket(ctx *gin.Context) {
	var req RemoveEventDateTicketParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at RemoveEventDateTicketParams in ShouldBindJSON: %v, eventID: %v \n", err.Error(), req.EventDateTimeID)
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
	ticketID, err := tools.StringToUuid(req.TicketID)
	if err != nil {
		log.Printf("Error at tools.StringToUuid: %v, TicketID: %v \n", err.Error(), req.TicketID)
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
		log.Printf("There an error at RemoveEventDateTicket at GetEventDateTimeByOption: %v, eventID: %v, userID: %v \n", err.Error(), eventDateTimeID, user.ID)
		err = fmt.Errorf("could not create ticket for this event date, please try again")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// We check if the eventDateTicket has an active booking
	isBooked := EventDateTicketIsBooked(ctx, server, ticketID, "RemoveEventDateTicket", tools.UuidToString(user.ID))

	if isBooked {
		err = fmt.Errorf("this ticket currently has an active booking so it can't be deleted")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	// We want to know if there were any ticket booked even if it was not complete
	isBookedAny := EventDateTicketIsBookedAny(ctx, server, ticketID, "RemoveEventDateTicket", tools.UuidToString(user.ID))
	if isBookedAny {
		_, err = server.store.UpdateEventDateTicket(ctx, db.UpdateEventDateTicketParams{
			IsActive: pgtype.Bool{
				Bool:  false,
				Valid: true,
			},
			ID:              ticketID,
			EventDateTimeID: eventDateTimeID,
		})
		if err != nil {
			log.Printf("There an error at server.store.UpdateEventDateTicket at GetEventDateTimeByOption: %v, eventID: %v, userID: %v \n", err.Error(), eventDateTimeID, user.ID)
			err = fmt.Errorf("could not remove ticket for this event date, please try again")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
	} else {
		err = server.store.RemoveEventDateTicket(ctx, db.RemoveEventDateTicketParams{
			ID:              ticketID,
			EventDateTimeID: eventDateTimeID,
		})
		if err != nil {
			log.Printf("There an error at RemoveEventDateTicket at GetEventDateTimeByOption: %v, eventID: %v, userID: %v \n", err.Error(), eventDateTimeID, user.ID)
			err = fmt.Errorf("could not remove ticket for this event date, please try again")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
	}

	var ticketCount int
	ticketCo, err := server.store.GetEventDateTicketCount(ctx, eventDateTimeID)
	if err != nil {
		log.Printf("There an error at RemoveEventDateTicket at GetEventDateTicketCount: %v, eventID: %v, userID: %v \n", err.Error(), eventDateTimeID, user.ID)
	} else {
		ticketCount = int(ticketCo)
	}

	res := RemoveEventDateTicketRes{
		TicketID:    req.TicketID,
		TicketCount: ticketCount,
	}
	if isCoHost {
		HandleCoHostUpdateMsg(ctx, server, userCoHost, user, option, "RemoveEventDateTicket", "event date ticket removed", "remove event date ticket")
	}
	ctx.JSON(http.StatusOK, res)
}
