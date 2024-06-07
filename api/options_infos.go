package api

import (
	"fmt"
	"log"
	"net/http"

	db "github.com/makuo12/ghost_server/db/sqlc"
	"github.com/makuo12/ghost_server/tools"
	"github.com/makuo12/ghost_server/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (server *Server) ListIncompleteOptionInfos(ctx *gin.Context) {
	var req ListIncompleteOptionInfosParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at ListIncompleteOptionInfos in ShouldBindJSON: %v, Offset: %v \n", err.Error(), req.Offset)
		err = fmt.Errorf("please enter all the details for this ticket")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := HandleGetUser(ctx, server)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	count, err := server.store.CountOptionInfo(ctx, db.CountOptionInfoParams{
		HostID:          user.ID,
		CoUserID:        tools.UuidToString(user.UserID),
		IsComplete:      false,
		IsActive:        true,
		OptionStatusOne: "unlist",
		OptionStatusTwo: "unlist",
	})
	if err != nil {
		log.Printf("Error at  ListIncompleteOptionInfos in CountOptionInfo err: %v, user: %v\n", err, user.ID)
		res := "none"
		ctx.JSON(http.StatusNoContent, res)
		return
	}
	if count <= int64(req.Offset) || count == 0 {
		res := "none"
		ctx.JSON(http.StatusNoContent, res)
		return
	}
	limit := 10
	if req.IsStart {
		limit = 3
	}
	optionInfos, err := server.store.ListOptionInfo(ctx, db.ListOptionInfoParams{
		HostID:          user.ID,
		CoUserID:        tools.UuidToString(user.UserID),
		IsComplete:      false,
		IsActive:        true,
		OptionStatusOne: "unlist",
		OptionStatusTwo: "unlist",
		Limit:           int32(limit),
		Offset:          int32(req.Offset),
	})
	if err != nil {
		log.Printf("Error at  ListIncompleteOptionInfos in ListNotification err: %v, user: %v\n", err, user.ID)
		res := "none"
		ctx.JSON(http.StatusNoContent, res)
	}
	var res ListOptionInfoNotCompleteRowParams
	var resData []ListOptionInfoNotCompleteRow
	for i := 0; i < len(optionInfos); i++ {
		extraInfo := "none"
		var dOptionID uuid.UUID
		switch optionInfos[i].CurrentState {
		case utils.ShortletInfo:
			extraInfo = HandleSqlNullString(optionInfos[i].SpaceType)
		case utils.EventSubType:
			extraInfo = HandleSqlNullString(optionInfos[i].EventType)
		}
		if optionInfos[i].HostType == "co_host" {
			dOptionID = optionInfos[i].CoHostID
		} else {
			dOptionID = optionInfos[i].OptionID
		}
		data := ListOptionInfoNotCompleteRow{
			ID:                 tools.UuidToString(dOptionID),
			IsComplete:         optionInfos[i].IsComplete,
			MainOptionType:     optionInfos[i].MainOptionType,
			Currency:           optionInfos[i].Currency,
			OptionType:         optionInfos[i].OptionType,
			CurrentServerView:  optionInfos[i].CurrentState,
			PreviousServerView: optionInfos[i].PreviousState,
			HostNameOption:     optionInfos[i].HostNameOption,
			ExtraInfo:          extraInfo,
			CreatedAt:          tools.ConvertTimeToStringDateOnly(optionInfos[i].CreatedAt),
		}
		resData = append(resData, data)
	}
	res = ListOptionInfoNotCompleteRowParams{
		OptionInfos: resData,
	}
	log.Println("res ", res)
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) ListCalenderOptionItems(ctx *gin.Context) {
	var req OptionOffsetParams
	var onLastIndex = false
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
	count, err := server.store.CountOptionInfo(ctx, db.CountOptionInfoParams{
		HostID:          user.ID,
		CoUserID:        tools.UuidToString(user.UserID),
		IsComplete:      true,
		IsActive:        true,
		OptionStatusOne: "list",
		OptionStatusTwo: "staged",
	})
	if err != nil {
		log.Printf("Error at  GetCalenderOptionItems in CountOptionInfo err: %v, user: %v", err, user.ID)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if count <= int64(req.OptionOffset) {
		res := "none"
		ctx.JSON(http.StatusNoContent, res)
		return
	}
	// we want to get the options that are complete
	optionInfos, err := server.store.ListOptionInfo(ctx, db.ListOptionInfoParams{
		HostID:          user.ID,
		CoUserID:        tools.UuidToString(user.UserID),
		IsComplete:      true,
		Limit:           5,
		Offset:          int32(req.OptionOffset),
		IsActive:        true,
		OptionStatusOne: "list",
		OptionStatusTwo: "staged",
	})
	if err != nil {
		if err == db.ErrorRecordNotFound {
			res := "none"
			ctx.JSON(http.StatusNoContent, res)
			return
		} else {
			log.Printf("Error at  GetCalenderOptionItems in ListOptionInfoComplete err: %v, user: %v", err, user.ID)
			err = fmt.Errorf("an error occurred while getting your experiences")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
	}
	var res CalenderOptionList
	var resData []CalenderOptionItem
	for i := 0; i < len(optionInfos); i++ {
		var dOptionID uuid.UUID
		if optionInfos[i].HostType == "co_host" {
			dOptionID = optionInfos[i].CoHostID
		} else {
			dOptionID = optionInfos[i].OptionID
		}
		data := CalenderOptionItem{
			HostNameOption: optionInfos[i].HostNameOption,
			OptionID:       tools.UuidToString(dOptionID),
			MainOptionType: optionInfos[i].MainOptionType,
			Currency:       optionInfos[i].Currency,
		}
		resData = append(resData, data)
	}
	if count <= int64(req.OptionOffset+len(optionInfos)) {
		onLastIndex = true
	}
	res = CalenderOptionList{
		List:         resData,
		OptionOffset: req.OptionOffset + len(optionInfos),
		OnLastIndex:  onLastIndex,
	}
	log.Println("calender data", res)
	ctx.JSON(http.StatusOK, res)
}
