package api

//import (
//	"database/sql"
//	db "flex_server/db/sqlc"
//	"flex_server/tools"
//	"flex_server/utils"
//	"fmt"
//	"log"
//	"net/http"

//	"github.com/gin-gonic/gin"
//)

//func (server *Server) RemoveLodgeInfo(ctx *gin.Context) {
//	var req OptionInfoRemoveRequest
//	if err := ctx.ShouldBindJSON(&req); err != nil {
//		log.Printf("Error at RemoveLodgeInfoType in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
//		err = fmt.Errorf("error occurred while taking you back")
//		ctx.JSON(http.StatusBadRequest, errorResponse(err))
//		return
//	}
//	log.Println(req)
//	requestID, err := tools.StringToUuid(req.OptionID)
//	if err != nil {
//		log.Printf("Error at tools.StringToUuid: %v, optionID: %v \n", err.Error(), req.OptionID)
//		err = fmt.Errorf("error occurred while processing your request")
//		ctx.JSON(http.StatusBadRequest, errorResponse(err))
//		return
//	}
//	user, option, err := HandleGetOptionIncomplete(requestID, ctx, server, false)
//	if err != nil {
//		ctx.JSON(http.StatusBadRequest, errorResponse(err))
//		return
//	}

//	err = server.store.RemoveLodgeInfo(ctx, option.ID)
//	if err != nil {
//		if err == db.ErrorRecordNotFound {
//			log.Printf("Hotel doesn't exist %v \n", err.Error())
//		} else {
//			log.Printf("Error at RemoveLodgeInfoType in RemoveLodgeInfo: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
//			err = fmt.Errorf("error occurred while taking you back, please try again")
//			ctx.JSON(http.StatusNotFound, errorResponse(err))
//			return
//		}

//	}

//	completeOption, err := HandleCompleteOption(utils.LodgeRoomType, utils.OptionType, server, ctx, option, user.ID)
//	if err != nil {
//		log.Printf("Error at HandleCompleteOption in RemoveLodgeInfoType: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
//	}
//	res := CreateOptionInfoParams{
//		OptionType:         option.OptionType,
//		Currency:           option.Currency,
//		MainOptionType:     option.MainOptionType,
//		OptionImg:          option.OptionImg,
//		OptionID:           tools.UuidToString(option.ID),
//		CurrentServerView:  completeOption.CurrentState,
//		PreviousServerView: completeOption.PreviousState,
//	}
//	ctx.JSON(http.StatusOK, res)
//}

//func (server *Server) RemoveRecreationInfo(ctx *gin.Context) {
//	var req OptionInfoRemoveRequest
//	if err := ctx.ShouldBindJSON(&req); err != nil {
//		log.Printf("Error at RemoveRecreationInfo in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
//		err = fmt.Errorf("error occurred while taking you back")
//		ctx.JSON(http.StatusBadRequest, errorResponse(err))
//		return
//	}
//	log.Println(req)
//	requestID, err := tools.StringToUuid(req.OptionID)
//	if err != nil {
//		log.Printf("Error at tools.StringToUuid: %v, optionID: %v \n", err.Error(), req.OptionID)
//		err = fmt.Errorf("error occurred while processing your request")
//		ctx.JSON(http.StatusBadRequest, errorResponse(err))
//		return
//	}
//	user, option, err := HandleGetOptionIncomplete(requestID, ctx, server, false)
//	if err != nil {
//		ctx.JSON(http.StatusBadRequest, errorResponse(err))
//		return
//	}

//	err = server.store.RemoveRecreationInfo(ctx, option.ID)
//	if err != nil {
//		if err == db.ErrorRecordNotFound {
//			log.Printf("Recreation experience doesn't exist %v \n", err.Error())
//		} else {
//			log.Printf("Error at RemoveRecreationInfo in RemoveRecreationInfo: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
//			err = fmt.Errorf("error occurred while taking you back, please try again")
//			ctx.JSON(http.StatusNotFound, errorResponse(err))
//			return
//		}

//	}

//	completeOption, err := HandleCompleteOption(utils.RecreationType, utils.OptionType, server, ctx, option, user.ID)
//	if err != nil {
//		log.Printf("Error at HandleCompleteOption in RemoveRecreationInfo: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
//	}
//	res := CreateOptionInfoParams{
//		OptionType:         option.OptionType,
//		Currency:           option.Currency,
//		MainOptionType:     option.MainOptionType,
//		OptionImg:          option.OptionImg,
//		OptionID:           tools.UuidToString(option.ID),
//		CurrentServerView:  completeOption.CurrentState,
//		PreviousServerView: completeOption.PreviousState,
//	}
//	ctx.JSON(http.StatusOK, res)
//}

//func (server *Server) RemoveChillInfo(ctx *gin.Context) {
//	var req OptionInfoRemoveRequest
//	if err := ctx.ShouldBindJSON(&req); err != nil {
//		log.Printf("Error at RemoveChillInfo in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
//		err = fmt.Errorf("error occurred while taking you back")
//		ctx.JSON(http.StatusBadRequest, errorResponse(err))
//		return
//	}
//	log.Println(req)
//	requestID, err := tools.StringToUuid(req.OptionID)
//	if err != nil {
//		log.Printf("Error at tools.StringToUuid: %v, optionID: %v \n", err.Error(), req.OptionID)
//		err = fmt.Errorf("error occurred while processing your request")
//		ctx.JSON(http.StatusBadRequest, errorResponse(err))
//		return
//	}
//	user, option, err := HandleGetOptionIncomplete(requestID, ctx, server, false)
//	if err != nil {
//		ctx.JSON(http.StatusBadRequest, errorResponse(err))
//		return
//	}

//	err = server.store.RemoveChillInfo(ctx, option.ID)
//	if err != nil {
//		if err == db.ErrorRecordNotFound {
//			log.Printf("Chill experience doesn't exist %v \n", err.Error())
//		} else {
//			log.Printf("Error at RemoveChillInfo in RemoveChillInfo: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
//			err = fmt.Errorf("error occurred while taking you back, please try again")
//			ctx.JSON(http.StatusNotFound, errorResponse(err))
//			return
//		}
//	}
//	completeOption, err := HandleCompleteOption(utils.ChillInfo, utils.OptionType, server, ctx, option, user.ID)
//	if err != nil {
//		log.Printf("Error at HandleCompleteOption in RemoveChillInfo: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
//	}
//	res := CreateOptionInfoParams{
//		OptionType:         option.OptionType,
//		Currency:           option.Currency,
//		MainOptionType:     option.MainOptionType,
//		OptionImg:          option.OptionImg,
//		OptionID:           tools.UuidToString(option.ID),
//		CurrentServerView:  completeOption.CurrentState,
//		PreviousServerView: completeOption.PreviousState,
//	}
//	ctx.JSON(http.StatusOK, res)
//}

//func (server *Server) RemoveOption(ctx *gin.Context) {
//	var req OptionInfoRemoveRequest
//	if err := ctx.ShouldBindJSON(&req); err != nil {
//		log.Printf("Error at RemoveOption in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
//		err = fmt.Errorf("error occurred while taking you back, try again")
//		ctx.JSON(http.StatusBadRequest, errorResponse(err))
//		return
//	}
//	log.Println(req)
//	requestID, err := tools.StringToUuid(req.OptionID)
//	if err != nil {
//		log.Printf("Error at tools.StringToUuid: %v, optionID: %v \n", err.Error(), req.OptionID)
//		err = fmt.Errorf("error occurred while processing your request")
//		ctx.JSON(http.StatusBadRequest, errorResponse(err))
//		return
//	}

//	user, option, err := HandleGetOptionIncomplete(requestID, ctx, server, false)
//	if err != nil {
//		ctx.JSON(http.StatusBadRequest, errorResponse(err))
//		return
//	}
//	// if option is an event we want to remove the event info
//	if option.MainOptionType == "events" {
//		err = HandleRemoveEventInfo(server, ctx, option, user.ID)
//		if err != nil {
//			ctx.JSON(http.StatusBadRequest, errorResponse(err))
//			return
//		}
//	}
//	err = server.store.RemoveCompleteOptionInfo(ctx, option.ID)
//	if err != nil {
//		log.Printf("error at RemoveOption at RemoveCompleteOptionInfo: %v, optionID: %v, userID: %vv\n", err.Error(), option.ID, user.ID)
//		err = fmt.Errorf("error occurred while taking you back, please try again")
//		ctx.JSON(http.StatusNotFound, errorResponse(err))
//		return
//	}
//	err = server.store.RemoveOptionInfoDetail(ctx, option.ID)
//	if err != nil {
//		log.Printf("error at RemoveOption at RemoveOptionInfoDetail: %v, optionID: %v, userID: %vv\n", err.Error(), option.ID, user.ID)
//		err = fmt.Errorf("error occurred while taking you back, please try again")
//		ctx.JSON(http.StatusNotFound, errorResponse(err))
//		return
//	}
//	arg := db.RemoveOptionInfoParams{
//		ID:     requestID,
//		HostID: user.ID,
//	}

//	err = server.store.RemoveOptionInfo(ctx, arg)
//	if err != nil {
//		if err == db.ErrorRecordNotFound {
//			log.Printf("Option doesn't exist %v \n", err.Error())
//		} else {
//			log.Printf("Error at RemoveOption in RemoveOptionInfo: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
//			err = fmt.Errorf("error occurred while taking you back, please try again")
//			ctx.JSON(http.StatusNotFound, errorResponse(err))
//			return
//		}

//	}

//	res := OptionInfoRemoveResponse{
//		Success: true,
//	}
//	ctx.JSON(http.StatusOK, res)
//}

//func (server *Server) RemoveShortletType(ctx *gin.Context) {
//	var req OptionInfoRemoveRequest
//	if err := ctx.ShouldBindJSON(&req); err != nil {
//		log.Printf("Error at RemoveShortletType in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
//		err = fmt.Errorf("error occurred while taking you back")
//		ctx.JSON(http.StatusBadRequest, errorResponse(err))
//		return
//	}
//	log.Println(req)
//	requestID, err := tools.StringToUuid(req.OptionID)
//	if err != nil {
//		log.Printf("Error at tools.StringToUuid: %v, optionID: %v \n", err.Error(), req.OptionID)
//		err = fmt.Errorf("error occurred while processing your request")
//		ctx.JSON(http.StatusBadRequest, errorResponse(err))
//		return
//	}
//	user, option, err := HandleGetOptionIncomplete(requestID, ctx, server, false)
//	if err != nil {
//		ctx.JSON(http.StatusBadRequest, errorResponse(err))
//		return
//	}
//	err = server.store.RemoveShortletType(ctx, option.ID)
//	if err != nil {
//		if err == db.ErrorRecordNotFound {
//			log.Printf("Shortlet doesn't exist %v \n", err.Error())
//		} else {
//			log.Printf("Error at RemoveShortletType in RemoveShortletType: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
//			err = fmt.Errorf("error occurred while taking you back, please try again")
//			ctx.JSON(http.StatusNotFound, errorResponse(err))
//			return
//		}

//	}
//	completeOption, err := HandleCompleteOption(utils.ShortletType, utils.OptionType, server, ctx, option, user.ID)
//	if err != nil {
//		log.Printf("Error at RemoveShortletType in HandleCompleteOption: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
//	}
//	res := CreateOptionInfoParams{
//		OptionType:         option.OptionType,
//		Currency:           option.Currency,
//		MainOptionType:     option.MainOptionType,
//		OptionImg:          option.OptionImg,
//		OptionID:           tools.UuidToString(option.ID),
//		CurrentServerView:  completeOption.CurrentState,
//		PreviousServerView: completeOption.PreviousState,
//	}
//	ctx.JSON(http.StatusOK, res)
//}

//func (server *Server) RemoveHallInfo(ctx *gin.Context) {
//	var req OptionInfoRemoveRequest
//	if err := ctx.ShouldBindJSON(&req); err != nil {
//		log.Printf("Error at RemoveHallInfo in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
//		err = fmt.Errorf("error occurred while taking you back")
//		ctx.JSON(http.StatusBadRequest, errorResponse(err))
//		return
//	}
//	log.Println(req)
//	requestID, err := tools.StringToUuid(req.OptionID)
//	if err != nil {
//		log.Printf("Error at tools.StringToUuid: %v, optionID: %v \n", err.Error(), req.OptionID)
//		err = fmt.Errorf("error occurred while processing your request")
//		ctx.JSON(http.StatusBadRequest, errorResponse(err))
//		return
//	}
//	user, option, err := HandleGetOptionIncomplete(requestID, ctx, server, false)
//	if err != nil {
//		ctx.JSON(http.StatusBadRequest, errorResponse(err))
//		return
//	}
//	err = server.store.RemoveHall(ctx, option.ID)
//	if err != nil {
//		if err == db.ErrorRecordNotFound {
//			log.Printf("Hall doesn't exist %v \n", err.Error())
//		} else {
//			log.Printf("Error at RemoveHallInfo in RemoveHall: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
//			err = fmt.Errorf("error occurred while taking you back, please try again")
//			ctx.JSON(http.StatusNotFound, errorResponse(err))
//			return
//		}
//	}
//	completeOption, err := HandleCompleteOption(utils.ShortletType, utils.OptionType, server, ctx, option, user.ID)
//	if err != nil {
//		log.Printf("Error at RemoveHallInfo in HandleCompleteOption: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
//	}
//	res := CreateOptionInfoParams{
//		OptionType:         option.OptionType,
//		Currency:           option.Currency,
//		MainOptionType:     option.MainOptionType,
//		OptionImg:          option.OptionImg,
//		OptionID:           tools.UuidToString(option.ID),
//		CurrentServerView:  completeOption.CurrentState,
//		PreviousServerView: completeOption.PreviousState,
//	}
//	ctx.JSON(http.StatusOK, res)
//}

//func (server *Server) RemoveLearnType(ctx *gin.Context) {
//	var req OptionInfoRemoveRequest
//	if err := ctx.ShouldBindJSON(&req); err != nil {
//		log.Printf("Error at RemoveLearnType in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
//		err = fmt.Errorf("error occurred while taking you back")
//		ctx.JSON(http.StatusBadRequest, errorResponse(err))
//		return
//	}
//	log.Println(req)
//	requestID, err := tools.StringToUuid(req.OptionID)
//	if err != nil {
//		log.Printf("Error at tools.StringToUuid: %v, optionID: %v \n", err.Error(), req.OptionID)
//		err = fmt.Errorf("error occurred while processing your request")
//		ctx.JSON(http.StatusBadRequest, errorResponse(err))
//		return
//	}
//	user, option, err := HandleGetOptionIncomplete(requestID, ctx, server, false)
//	if err != nil {
//		ctx.JSON(http.StatusBadRequest, errorResponse(err))
//		return
//	}
//	err = server.store.RemoveLearn(ctx, option.ID)
//	if err != nil {
//		if err == db.ErrorRecordNotFound {
//			log.Printf("Learn doesn't exist %v \n", err.Error())
//		} else {
//			log.Printf("Error at RemoveLearnType in RemoveLearnType: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
//			err = fmt.Errorf("error occurred while taking you back, please try again")
//			ctx.JSON(http.StatusNotFound, errorResponse(err))
//			return
//		}

//	}
//	completeOption, err := HandleCompleteOption(utils.LearnType, utils.OptionType, server, ctx, option, user.ID)
//	if err != nil {
//		log.Printf("Error at RemoveLearnType in HandleCompleteOption: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
//	}
//	res := CreateOptionInfoParams{
//		OptionType:         option.OptionType,
//		Currency:           option.Currency,
//		MainOptionType:     option.MainOptionType,
//		OptionImg:          option.OptionImg,
//		OptionID:           tools.UuidToString(option.ID),
//		CurrentServerView:  completeOption.CurrentState,
//		PreviousServerView: completeOption.PreviousState,
//	}
//	ctx.JSON(http.StatusOK, res)
//}

//func (server *Server) RemoveYatchInfo(ctx *gin.Context) {
//	var req OptionInfoRemoveRequest
//	if err := ctx.ShouldBindJSON(&req); err != nil {
//		log.Printf("Error at RemoveYatchInfo in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
//		err = fmt.Errorf("error occurred while taking you back")
//		ctx.JSON(http.StatusBadRequest, errorResponse(err))
//		return
//	}
//	log.Println(req)
//	requestID, err := tools.StringToUuid(req.OptionID)
//	if err != nil {
//		log.Printf("Error at tools.StringToUuid: %v, optionID: %v \n", err.Error(), req.OptionID)
//		err = fmt.Errorf("error occurred while processing your request")
//		ctx.JSON(http.StatusBadRequest, errorResponse(err))
//		return
//	}
//	user, option, err := HandleGetOptionIncomplete(requestID, ctx, server, false)
//	if err != nil {
//		ctx.JSON(http.StatusBadRequest, errorResponse(err))
//		return
//	}
//	err = server.store.RemoveYatch(ctx, option.ID)
//	if err != nil {
//		if err == db.ErrorRecordNotFound {
//			log.Printf("Learn doesn't exist %v \n", err.Error())
//		} else {
//			log.Printf("Error at RemoveYatchInfo in RemoveYatch: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
//			err = fmt.Errorf("error occurred while taking you back, please try again")
//			ctx.JSON(http.StatusNotFound, errorResponse(err))
//			return
//		}

//	}
//	completeOption, err := HandleCompleteOption(utils.YatchInfo, utils.OptionType, server, ctx, option, user.ID)
//	if err != nil {
//		log.Printf("Error at RemoveYatchInfo in HandleCompleteOption: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
//	}
//	res := CreateOptionInfoParams{
//		OptionType:         option.OptionType,
//		Currency:           option.Currency,
//		MainOptionType:     option.MainOptionType,
//		OptionImg:          option.OptionImg,
//		OptionID:           tools.UuidToString(option.ID),
//		CurrentServerView:  completeOption.CurrentState,
//		PreviousServerView: completeOption.PreviousState,
//	}
//	ctx.JSON(http.StatusOK, res)
//}

//func (server *Server) RemoveSportInfo(ctx *gin.Context) {
//	var req OptionInfoRemoveRequest
//	if err := ctx.ShouldBindJSON(&req); err != nil {
//		log.Printf("Error at RemoveSportInfo in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
//		err = fmt.Errorf("error occurred while taking you back")
//		ctx.JSON(http.StatusBadRequest, errorResponse(err))
//		return
//	}
//	log.Println(req)
//	requestID, err := tools.StringToUuid(req.OptionID)
//	if err != nil {
//		log.Printf("Error at tools.StringToUuid: %v, optionID: %v \n", err.Error(), req.OptionID)
//		err = fmt.Errorf("error occurred while processing your request")
//		ctx.JSON(http.StatusBadRequest, errorResponse(err))
//		return
//	}
//	user, option, err := HandleGetOptionIncomplete(requestID, ctx, server, false)
//	if err != nil {
//		ctx.JSON(http.StatusBadRequest, errorResponse(err))
//		return
//	}

//	err = server.store.RemoveSportInfo(ctx, option.ID)
//	if err != nil {
//		if err == db.ErrorRecordNotFound {
//			log.Printf("Sport experience doesn't exist %v \n", err.Error())
//		} else {
//			log.Printf("Error at RemoveSportInfo in RemoveSportInfo: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
//			err = fmt.Errorf("error occurred while taking you back, please try again")
//			ctx.JSON(http.StatusNotFound, errorResponse(err))
//			return
//		}

//	}

//	completeOption, err := HandleCompleteOption(utils.SportType, utils.OptionType, server, ctx, option, user.ID)
//	if err != nil {
//		log.Printf("Error at HandleCompleteOption in RemoveSportInfo: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
//	}
//	res := CreateOptionInfoParams{
//		OptionType: option.OptionType,
//		Currency:   option.Currency,

//		MainOptionType:     option.MainOptionType,
//		OptionImg:          option.OptionImg,
//		OptionID:           tools.UuidToString(option.ID),
//		CurrentServerView:  completeOption.CurrentState,
//		PreviousServerView: completeOption.PreviousState,
//	}
//	ctx.JSON(http.StatusOK, res)
//}
