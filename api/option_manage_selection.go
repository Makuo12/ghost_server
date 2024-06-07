package api

import (
	"fmt"
	"log"
	"net/http"

	db "github.com/makuo12/ghost_server/db/sqlc"
	"github.com/makuo12/ghost_server/tools"
	"github.com/makuo12/ghost_server/utils"

	"github.com/gin-gonic/gin"
)

func (server *Server) ListUHMOptionSelection(ctx *gin.Context) {
	var req OptionSelectionOffsetParams

	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at  GetCalenderOptionItems in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionOffset)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Println(req)
	user, err := HandleGetUser(ctx, server)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var res ListUHMOptionSelectionRes
	var hasData bool
	switch req.Type {
	case "complete":
		res, err, hasData = HandleListOptionSelectComplete(ctx, server, user, req)

	case "in_progress":
		res, err, hasData = HandleListOptionSelectInProgress(ctx, server, user, req)
	case "in_active":
		res, err, hasData = HandleListOptionSelectInActive(ctx, server, user, req)
	}
	if hasData && err == nil {
		ctx.JSON(http.StatusOK, res)
		return

	} else if !hasData && err == nil {
		ctx.JSON(http.StatusNoContent, res)
		return
	} else {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

}

func (server *Server) GetOptionInfoIncomplete(ctx *gin.Context) {
	optionID := ctx.Param("option_id")
	requestID, err := tools.StringToUuid(optionID)
	if err != nil {
		log.Printf("Error at tools.StringToUuid: %v, optionID: %v \n", err.Error(), optionID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := HandleGetUser(ctx, server)
	if err != nil {
		err = fmt.Errorf("you cannot access this resource")
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}
	optionInfo, err := server.store.GetOptionInfoData(ctx, db.GetOptionInfoDataParams{
		HostID:     user.ID,
		IsComplete: false,
		IsActive:   true,
		ID:         requestID,
	})
	if err != nil {
		log.Printf("Error at GetOptionInfoIncomplete  in GetOptionInfoData err: %v, user: %v", err, user.ID)
		err = fmt.Errorf("an error occurred while getting your option info")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	extraInfo := "none"
	switch optionInfo.CurrentState {
	case utils.ShortletInfo:
		shortlet, err := server.store.GetShortlet(ctx, optionInfo.ID)
		if err != nil {
			log.Printf("Error at GetOptionInfoIncomplete  in for ShortletInfo GetShortlet err: %v, user: %v", err, user.ID)
		} else {
			extraInfo = shortlet.SpaceType
		}
	case utils.EventSubType:
		event, err := server.store.GetEventInfo(ctx, optionInfo.ID)
		if err != nil {
			log.Printf("Error at GetOptionInfoIncomplete  in GetEventInfo err: %v, user: %v", err, user.ID)
		} else {
			extraInfo = event.EventType
		}
	}
	res := GetOptionInfoNotComplete{
		ID:                 tools.UuidToString(optionInfo.ID),
		IsComplete:         optionInfo.IsComplete,
		MainOptionType:     optionInfo.MainOptionType,
		Currency:           optionInfo.Currency,
		OptionType:         optionInfo.OptionType,
		CurrentServerView:  optionInfo.CurrentState,
		PreviousServerView: optionInfo.PreviousState,
		HostNameOption:     optionInfo.HostNameOption,
		ExtraInfo:          extraInfo,
		CreatedAt:          tools.ConvertTimeToStringDateOnly(optionInfo.CreatedAt),
	}
	ctx.JSON(http.StatusOK, res)
}
