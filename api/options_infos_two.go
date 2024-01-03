package api

import (
	db "flex_server/db/sqlc"
	"flex_server/tools"
	"flex_server/utils"

	//"flex_server/val"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

func (server *Server) CreateEventSubCategory(ctx *gin.Context) {
	var req CreateEventSubCategoryParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at CreateEventSubCategoryParams in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("an error occurred while processing the details you sent please make sure you select one of the provided options")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Println(req)
	requestID, err := tools.StringToUuid(req.OptionID)
	if err != nil {
		log.Printf("Error at tools.StringToUuid: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, option, err := HandleGetOptionIncomplete(requestID, ctx, server, false)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if option.OptionType != "events" {
		err = fmt.Errorf("this request is not allowed")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	arg := db.UpdateEventInfoParams{
		SubCategoryType: pgtype.Text{
			String: req.SubCategoryType,
			Valid:  true,
		},
		OptionID: option.ID,
	}
	_, err = server.store.UpdateEventInfo(ctx, arg)
	if err != nil {
		log.Printf("error at CreateEventSubCategory in UpdateEventInfo err is %v, optionID: %v, userID: %v \n", err, option.ID, user.ID)
		err = fmt.Errorf("error occurred while setting the sub category of your event")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	completeOption, err := HandleCompleteOption(utils.Description, utils.EventSubType, server, ctx, option, user.ID)
	if err != nil {
		log.Printf("Error at CreateEventSubCategory in HandleCompleteOption: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
	}
	res := OptionInfoResponse{
		Success:            true,
		OptionID:           tools.UuidToString(option.ID),
		UserOptionID:       tools.UuidToString(option.OptionUserID),
		CurrentServerView:  completeOption.CurrentState,
		PreviousServerView: completeOption.PreviousState,
		MainOptionType:     option.MainOptionType,
		OptionType:         option.OptionType,
		Currency:           option.Currency,
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) PublishOption(ctx *gin.Context) {
	var req PublishOption
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at PublishOption in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("an error occurred while adding your price")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Println(req)
	requestID, err := tools.StringToUuid(req.OptionID)
	if err != nil {
		log.Printf("Error at tools.StringToUuid: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, option, err := HandleGetOptionIncomplete(requestID, ctx, server, false)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	_, err = server.store.UpdateOptionInfoComplete(ctx, db.UpdateOptionInfoCompleteParams{
		ID:         option.ID,
		IsComplete: true,
	})
	if err != nil {
		log.Printf("There an error at PublishOption at UpdateOptionInfoComplete: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	completeOption, err := HandleCompleteOption(utils.Publish, utils.Publish, server, ctx, option, user.ID)
	if err != nil {
		log.Printf("There an error at PublishOption at HandleCompleteOption: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
	}
	// We want to update option info status
	_, err = server.store.UpdateOptionInfoStartStatus(ctx, db.UpdateOptionInfoStartStatusParams{
		Status:   "list",
		OptionID: option.ID,
	})
	if err != nil {
		log.Printf("There an error at PublishOption at store.UpdateOptionInfoStartStatus: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
	}
	switch option.MainOptionType {
	case "options":
		CreateOptionAlgo(ctx, server, option, user)
	case "events":
		CreateEventAlgo(ctx, server, option, user)
	}
	res := OptionInfoResponse{
		OptionItemType:     "",
		Success:            true,
		OptionID:           tools.UuidToString(option.ID),
		UserOptionID:       tools.UuidToString(option.OptionUserID),
		CurrentServerView:  completeOption.CurrentState,
		PreviousServerView: completeOption.PreviousState,
		MainOptionType:     option.MainOptionType,
		OptionType:         option.OptionType,
		Currency:           option.Currency,
	}
	ctx.JSON(http.StatusOK, res)
}
