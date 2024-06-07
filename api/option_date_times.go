package api

import (
	"fmt"
	"log"
	"net/http"

	db "github.com/makuo12/ghost_server/db/sqlc"
	"github.com/makuo12/ghost_server/tools"

	"github.com/gin-gonic/gin"
)

func (server *Server) ListOptionDateItems(ctx *gin.Context) {
	var req OptionDateParams

	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at ListOptionDateItems OptionDateParams in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Println(req)
	requestID, err := tools.StringToUuid(req.OptionID)
	if err != nil {
		log.Printf("Error at ListOptionDateItems at tools.StringToUuid: %v, OptionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, _, option, _, _, err := HandleGetCompleteOptionCalender(requestID, ctx, server, true)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var resAvailableData []OptionDateItem
	var resBookData []DateOptionBookedItem
	// We want to verify access to this resource
	optionData, err := server.store.GetShortletDateTimeByOption(ctx, db.GetShortletDateTimeByOptionParams{
		OptionID:   option.ID,
		ID:         user.ID,
		IsComplete: true,
	})
	if err != nil {
		log.Printf("Error at  ListOptionDateItems in GetShortletDateTimeByOption err: %v, user: %v\n", err, user.ID)
		err = fmt.Errorf("you don't have access to this resource")
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}
	optionDateTimes, err := server.store.ListAllOptionDateTime(ctx, requestID)
	if err != nil || len(optionDateTimes) == 0 {
		if err != nil {
			log.Printf("Error at ListOptionDateItems at tools.StringToUuid: %v, OptionID: %v \n", err.Error(), req.OptionID)
		}
		resAvailableData = []OptionDateItem{{"none", "none", "none", true, "none", true}}
	} else {
		for i := 0; i < len(optionDateTimes); i++ {
			var exist bool
			match, err := tools.AreDatesInSameMonthAndYear(tools.ConvertDateOnlyToString(optionDateTimes[i].Date), req.Date)
			if err != nil {
				log.Printf("Error at ListOptionDateItems at tools.AreDatesInSameMonthAndYear: %v, OptionID: %v \n", err.Error(), req.OptionID)
			} else {
				if match {
					exist = true
				}
			}
			if exist {
				data := OptionDateItem{
					ID:        tools.UuidToString(optionDateTimes[i].ID),
					OptionID:  tools.UuidToString(optionDateTimes[i].OptionID),
					Date:      tools.ConvertDateOnlyToString(optionDateTimes[i].Date),
					Available: optionDateTimes[i].Available,
					IsEmpty:   false,
					Price:     tools.IntToMoneyString(optionDateTimes[i].Price),
				}
				resAvailableData = append(resAvailableData, data)
			}
		}
	}
	dates := GetOptionDatesBooked(ctx, server, option.OptionUserID, "ListOptionDateBooked", tools.UuidToString(user.ID))

	if len(dates) > 0 {
		for _, d := range dates {
			var exist bool
			myDates := tools.GenerateDateListStringFromTime(d.StartDate, d.EndDate)
			for _, m := range myDates {
				match, err := tools.AreDatesInSameMonthAndYear(m, req.Date)
				if err != nil {
					log.Printf("Error at ListOptionDateItems at tools.AreDatesInSameMonthAndYear: %v, OptionID: %v \n", err.Error(), req.OptionID)
				} else {
					if match {
						exist = true
					}
				}
			}
			if exist {
				data := DateOptionBookedItem{
					ID:           tools.UuidToString(d.ReferenceID),
					StartDate:    tools.ConvertDateOnlyToString(d.StartDate),
					EndDate:      tools.ConvertDateOnlyToString(d.EndDate),
					ProfilePhoto: d.Photo,
					Booked:       true,
					FirstName:    d.FirstName,
					UserID:       tools.UuidToString(d.UserID),
					IsEmpty:      false,
				}
				resBookData = append(resBookData, data)
			}

		}
	} else {
		resBookData = []DateOptionBookedItem{{"none", "none", "none", "none", false, "none", "none", true}}
	}
	log.Println("resAvailableData", resAvailableData)
	log.Println("resBookData", resBookData)
	res := ListOptionDateItem{
		List:         resAvailableData,
		ListBooked:   resBookData,
		BasePrice:    tools.IntToMoneyString(optionData.Price),
		WeekendPrice: tools.IntToMoneyString(optionData.WeekendPrice),
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) CreateUpdateOptionDateTime(ctx *gin.Context) {
	var req CreateUpdateOptionDateTimeParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at CreateUpdateOptionDateTimeParams in ShouldBindJSON: %v, EventDateTimeID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("please make sure you select at least a start date")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Println(req)
	requestID, err := tools.StringToUuid(req.OptionID)
	if err != nil {
		log.Printf("Error at tools.StringToUuid: %v, OptionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, _, option, isCoHost, userCoHost, err := HandleGetCompleteOptionCalender(requestID, ctx, server, true)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	optionData, err := server.store.GetShortletDateTimeByOption(ctx, db.GetShortletDateTimeByOptionParams{
		OptionID:   option.ID,
		ID:         user.ID,
		IsComplete: true,
	})
	if err != nil {
		log.Printf("Error at  CreateUpdateOptionDateTime in GetShortletDateTimeByOption err: %v, user: %v\n", err, user.ID)
		err = fmt.Errorf("you don't have access to this resource")
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}
	if len(req.Price) == 0 {
		req.Price = tools.IntToMoneyString(optionData.Price)
	}
	count, err := server.store.GetOptionDateTimeCount(ctx, requestID)
	if err != nil {
		log.Printf("Error at  CreateUpdateOptionDateTime in GetOptionDateTimeByOption err: %v, user: %v\n", err, user.ID)
		count = 0
	}
	err = HandleUpdateCreateOptionDate(ctx, server, requestID, req.Price, req, int(count), user.ID, optionData.Currency)
	if err != nil {
		log.Printf("Error at  CreateUpdateOptionDateTime in HandleUpdateCreateOptionDate err: %v, user: %v\n", err, user.ID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var resAvailableData []OptionDateItem
	var resBookData []DateOptionBookedItem
	optionDateTimes, err := server.store.ListAllOptionDateTime(ctx, requestID)
	if err != nil || len(optionDateTimes) == 0 {
		if err != nil {
			log.Printf("Error at ListOptionDateItems at tools.StringToUuid: %v, OptionID: %v \n", err.Error(), req.OptionID)
		}
		resAvailableData = []OptionDateItem{{"none", "none", "none", true, "none", true}}
	} else {
		for i := 0; i < len(optionDateTimes); i++ {
			data := OptionDateItem{
				ID:        tools.UuidToString(optionDateTimes[i].ID),
				OptionID:  tools.UuidToString(optionDateTimes[i].OptionID),
				Date:      tools.ConvertDateOnlyToString(optionDateTimes[i].Date),
				Available: optionDateTimes[i].Available,
				IsEmpty:   false,
				Price:     tools.IntToMoneyString(optionDateTimes[i].Price),
			}
			resAvailableData = append(resAvailableData, data)
		}
	}
	dates := GetOptionDatesBooked(ctx, server, option.OptionUserID, "ListOptionDateBooked", tools.UuidToString(user.ID))

	if len(dates) > 0 {
		for _, d := range dates {
			data := DateOptionBookedItem{
				ID:           tools.UuidToString(d.ReferenceID),
				StartDate:    tools.ConvertDateOnlyToString(d.StartDate),
				EndDate:      tools.ConvertDateOnlyToString(d.EndDate),
				ProfilePhoto: d.Photo,
				Booked:       true,
				FirstName:    d.FirstName,
				UserID:       tools.UuidToString(d.UserID),
				IsEmpty:      false,
			}
			resBookData = append(resBookData, data)
		}
	} else {
		resBookData = []DateOptionBookedItem{{"none", "none", "none", "none", false, "none", "none", true}}
	}
	if isCoHost {
		HandleCoHostUpdateMsg(ctx, server, userCoHost, user, option, "ListEventDateBooked", "event date times", "create update event date times")
	}
	res := ListOptionDateItem{
		List:         resAvailableData,
		ListBooked:   resBookData,
		BasePrice:    tools.IntToMoneyString(optionData.Price),
		WeekendPrice: tools.IntToMoneyString(optionData.WeekendPrice),
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) GetOptionDateNote(ctx *gin.Context) {
	var req OptionSelectedDateParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at  OptionDateParams in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionDateTimeID)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Println(req)
	requestID, err := tools.StringToUuid(req.OptionID)
	if err != nil {
		log.Printf("Error at tools.StringToUuid: %v, OptionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	id, err := tools.StringToUuid(req.OptionDateTimeID)
	if err != nil {
		log.Printf("Error at tools.StringToUuid: %v, OptionDateTimeID: %v \n", err.Error(), req.OptionDateTimeID)
		err = fmt.Errorf("error occurred while processing your id")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, _, option, _, _, err := HandleGetCompleteOptionCalender(requestID, ctx, server, true)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	// We want to verify access to this resource
	data, err := server.store.GetOptionDateTimeNoteByOption(ctx, db.GetOptionDateTimeNoteByOptionParams{
		ID:         id,
		ID_2:       user.ID,
		ID_3:       option.ID,
		IsComplete: true,
	})

	if err != nil {
		log.Printf("Error at  GetOptionDateNote in GetOptionDateTimeNoteByOption err: %v, user: %v\n", err, user.ID)
		err = fmt.Errorf("could not get the note")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var note string
	if tools.ServerStringEmpty(data.Note) {
		note = "none"
	} else {
		note = data.Note
	}
	res := GetOptionDateNoteRes{
		OptionDateTimeID: tools.UuidToString(data.ID),
		Note:             note,
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) RemoveOptionDateTime(ctx *gin.Context) {
	var req OptionSelectedDateParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at  OptionDateParams in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionDateTimeID)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Println(req)
	requestID, err := tools.StringToUuid(req.OptionID)
	if err != nil {
		log.Printf("Error at tools.StringToUuid: %v, OptionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	id, err := tools.StringToUuid(req.OptionDateTimeID)
	if err != nil {
		log.Printf("Error at tools.StringToUuid: %v, OptionDateTimeID: %v \n", err.Error(), req.OptionDateTimeID)
		err = fmt.Errorf("error occurred while processing your id")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, _, option, _, _, err := HandleGetCompleteOptionCalender(requestID, ctx, server, true)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	// We want to verify access to this resource
	err = server.store.RemoveOptionDateTime(ctx, db.RemoveOptionDateTimeParams{
		ID:       id,
		OptionID: option.ID,
	})
	if err != nil {
		log.Printf("Error at RemoveOptionDateTime in RemoveOptionDateTime err: %v, user: %v\n", err, user.ID)
		err = fmt.Errorf("you don't have access to this resource")
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}
	res := UserResponseMsg{
		Success: true,
	}
	ctx.JSON(http.StatusOK, res)
}
