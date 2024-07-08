package api

import (
	"context"
	"fmt"
	"log"
	"net/http"

	db "github.com/makuo12/ghost_server/db/sqlc"
	"github.com/makuo12/ghost_server/tools"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func (server *Server) UpdateEventDateStatus(ctx *gin.Context) {
	var req UpdateEventDateStatusReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at UpdateEventDateStatus in ShouldBindJSON: %v, EventDateTimeID: %v \n", err.Error(), req.EventDateTimeID)
		err = fmt.Errorf("please make sure yuo select at least a start date")
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
		log.Printf("There an error at UpdateEventDateStatus at GetEventDateTimeByOption: %v, EventDateTimeID: %v, userID: %v \n", err.Error(), eventDateTimeID, user.ID)
		err = fmt.Errorf("could not set event date status, please try again using the format on the app")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	isBooked := EventDateIsBooked(ctx, server, eventDateTimeID, "UpdateEventDateStatus", tools.UuidToString(user.ID))
	if isBooked && (req.Status == "cancelled" || req.Status == "postponed") {
		err = fmt.Errorf("this status cannot be changed to postponed or cancelled because tickets have already been sold. To modify the date of this event, go to the app's hosting area")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	eventDate, err := server.store.UpdateEventDateTime(ctx, db.UpdateEventDateTimeParams{
		Status: pgtype.Text{
			String: req.Status,
			Valid:  true,
		},
		ID:   eventDateTimeID,
		Type: req.Type,
	})
	if err != nil {
		log.Printf("Error at UpdateEventDateStatus in UpdateEventDateTime err: %v, user: %v\n", err, user.ID)
		err = fmt.Errorf("an error occurred while updating your event date details")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	res := UpdateEventDateStatusRes{
		Status: eventDate.Status,
	}
	log.Println(res)
	if isCoHost {
		HandleCoHostUpdateMsg(ctx, server, userCoHost, user, option, "UpdateEventDateStatus", "event date status", "update event date status")
	}
	ctx.JSON(http.StatusOK, res)
}

func HandleListEventDates(ctx context.Context, server *Server, optionID uuid.UUID, user db.User, funcName string, offset int) ([]EventDateItem, GetUHMDataOptionRes, error) {
	empty := []EventDateItem{
		{
			ID:               "none",
			Name:             "none",
			StartTime:        "",
			EndTime:          "",
			StartDate:        "",
			Status:           "",
			EndDate:          "",
			Tickets:          0,
			Note:             "",
			TimeZone:         "",
			Type:             "",
			EventDates:       []string{"none"},
			NeedBands:        false,
			NeedTickets:      false,
			AbsorbBandCharge: false,
		},
	}
	var data GetUHMDataOptionRes
	if offset < 1 {
		uhmData, err := server.store.GetOptionEventUHMData(ctx, optionID)
		if err != nil {
			log.Printf("There an FuncName: %v error at GetUHMOptionData at GetOptionInfoUHMData: %v, optionID: %v, userID: %v \n", funcName, err.Error(), optionID, user.ID)
			err = fmt.Errorf("could not get the data")
			return empty, GetUHMDataOptionRes{}, err
		}
		data = GetUHMDataOptionRes{
			HostNameOption: uhmData.HostNameOption,
			SpaceAreas:     []string{"none"},
			Price:          "",
			SpaceType:      "",
			Category:       uhmData.Category,
			CategoryTwo:    uhmData.CategoryTwo,
			CategoryThree:  uhmData.CategoryThree,
			CategoryFour:   uhmData.CategoryFour,
			NumOfGuest:     0,
			MainImage:      uhmData.MainImage,
			Images:         uhmData.Images,
			OptionUserID:   tools.UuidToString(uhmData.OptionUserID),
			Street:         "",
			State:          "",
			City:           "",
			Country:        "",
			Postcode:       "",
			CheckInMethod:  "",
			EventType:      uhmData.EventType,
			EventSubType:   uhmData.SubCategoryType,
			Currency:       uhmData.Currency,
			Status:         tools.HandleOptionStatus(uhmData.OptionStatus),
		}
	} else {
		data = GetUHMDataOptionRes{
			HostNameOption: "",
			SpaceAreas:     []string{"none"},
			Price:          "",
			SpaceType:      "",
			Category:       "",
			CategoryTwo:    "",
			CategoryThree:  "",
			CategoryFour:   "",
			NumOfGuest:     0,
			MainImage:     "",
			Images:         []string{"none"},
			OptionUserID:   "",
			Street:         "",
			State:          "",
			City:           "",
			Country:        "",
			Postcode:       "",
			CheckInMethod:  "",
			EventType:      "",
			EventSubType:   "",
			Currency:       "",
			Status:         "",
		}
	}
	return empty, data, nil
}

func (server *Server) UpdatePublishEventCheckInStep(ctx *gin.Context) {
	var req UpdatePublishEventCheckInStepParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at UpdateEventDateStatus in ShouldBindJSON: %v, EventDateTimeID: %v \n", err.Error(), req.EventDateTimeID)
		err = fmt.Errorf("please make sure yuo select at least a start date")
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
		log.Printf("There an error at UpdatePublishEventCheckInStep at GetEventDateTimeByOption: %v, EventDateTimeID: %v, userID: %v \n", err.Error(), eventDateTimeID, user.ID)
		err = fmt.Errorf("could not set event date status, please try again using the format on the app")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	published, err := server.store.UpdateEventPublishCheckInStep(ctx, eventDateTimeID)
	if err != nil {
		log.Printf("There an error at UpdatePublishEventCheckInStep at UpdateEventPublishCheckInStep: %v, EventDateTimeID: %v, userID: %v \n", err.Error(), eventDateTimeID, user.ID)
		err = fmt.Errorf("could not publish your check in step")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	res := UpdatePublishEventCheckInStepRes{
		Published: published,
	}
	log.Println(res)
	if isCoHost {
		HandleCoHostUpdateMsg(ctx, server, userCoHost, user, option, "UpdateEventDateStatus", "event date status", "update event date status")
	}
	ctx.JSON(http.StatusOK, res)
}
