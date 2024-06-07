package api

//import (
//	"database/sql"
//	db "github.com/makuo12/ghost_server/db/sqlc"
//	"github.com/makuo12/ghost_server/tools"
//	"github.com/makuo12/ghost_server/utils"
//	"fmt"
//	"log"
//	"net/http"

//	"github.com/gin-gonic/gin"
//)

//func (server *Server) CreateLearnType(ctx *gin.Context) {
//	var req CreateLearnTypeParams
//	if err := ctx.ShouldBindJSON(&req); err != nil {
//		log.Printf("Error at CreateLearnType in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
//		err = fmt.Errorf("select one of the provided options")
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
//	if option.OptionType != "learn" {
//		err = fmt.Errorf("error occurred while performing your request. This experience type must be for a learn activity")
//		ctx.JSON(http.StatusForbidden, errorResponse(err))
//		return
//	}
//	err = server.store.RemoveLearn(ctx, option.ID)
//	if err != nil {
//		if err == db.ErrorRecordNotFound {
//			log.Printf("learn doesn't exist %v \n", err.Error())
//		} else {
//			log.Println(err.Error())
//			err = fmt.Errorf("error occurred while performing your request. Try again")
//			ctx.JSON(http.StatusForbidden, errorResponse(err))
//			return
//		}
//	}
//	argLearn := db.CreateLearnParams{
//		OptionID:  option.ID,
//		LearnType: req.LearnType,
//	}

//	Learn, err := server.store.CreateLearn(ctx, argLearn)
//	if err != nil {
//		log.Printf("Error at CreateLearnType in CreateLearn: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
//		err = fmt.Errorf("an error occurred while creating your Learn, try again or contact help")
//		ctx.JSON(http.StatusBadRequest, errorResponse(err))
//		return
//	}
//	completeOption, err := HandleCompleteOption(utils.LocationView, utils.LearnType, server, ctx, option, user.ID)
//	if err != nil {
//		log.Printf("Error at CreateLearnType in HandleCompleteOption: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
//	}
//	res := OptionInfoResponse{
//		OptionItemType:     "",
//		Success:            true,
//		OptionID:           tools.UuidToString(Learn.OptionID),
//		CurrentServerView:  completeOption.CurrentState,
//		PreviousServerView: completeOption.PreviousState,
//		MainOptionType:     option.MainOptionType,
//		OptionType:         option.OptionType,
//		Currency:           option.Currency,
//	}
//	ctx.JSON(http.StatusOK, res)
//}

//func (server *Server) CreateShortletType(ctx *gin.Context) {
//	var req CreateShortletTypeParams
//	if err := ctx.ShouldBindJSON(&req); err != nil {
//		log.Printf("Error at CreateShortletType in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
//		err = fmt.Errorf("an error occurred while processing the details you sent please make sure you select one of the provided options")
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
//	if option.OptionType != "shortlets" {
//		err = fmt.Errorf("error occurred while performing your request. This experience type must be for a shortlet")
//		ctx.JSON(http.StatusForbidden, errorResponse(err))
//		return
//	}
//	err = server.store.RemoveShortletType(ctx, option.ID)
//	if err != nil {
//		if err == db.ErrorRecordNotFound {
//			log.Printf("Shortlet doesn't exist %v \n", err.Error())
//		} else {
//			log.Println(err.Error())
//			err = fmt.Errorf("error occurred while performing your request. Try again")
//			ctx.JSON(http.StatusForbidden, errorResponse(err))
//			return
//		}
//	}
//	argShortlet := db.CreateShortletParams{
//		OptionID:       option.ID,
//		TypeOfShortlet: req.ShortletType,
//	}

//	shortlet, err := server.store.CreateShortlet(ctx, argShortlet)
//	if err != nil {
//		log.Printf("Error at CreateShortletType in CreateShortlet: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
//		err = fmt.Errorf("an error occurred while creating your shortlet, try again or contact help")
//		ctx.JSON(http.StatusBadRequest, errorResponse(err))
//		return
//	}
//	completeOption, err := HandleCompleteOption(utils.ShortletSpace, utils.ShortletType, server, ctx, option, user.ID)
//	if err != nil {
//		log.Printf("Error at CreateShortletType in HandleCompleteOption: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
//	}
//	res := OptionInfoResponse{
//		Success:            true,
//		OptionID:           tools.UuidToString(shortlet.OptionID),
//		CurrentServerView:  completeOption.CurrentState,
//		PreviousServerView: completeOption.PreviousState,
//		MainOptionType:     option.MainOptionType,
//		OptionType:         option.OptionType,
//		Currency:           option.Currency,
//	}
//	ctx.JSON(http.StatusOK, res)
//}

//func (server *Server) CreateYatchInfo(ctx *gin.Context) {
//	var req CreateYatchInfoParams
//	var numFullBathrooms = 0
//	var numHalfBathrooms = 0
//	if err := ctx.ShouldBindJSON(&req); err != nil {
//		log.Printf("Error at CreateYatchInfo in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
//		err = fmt.Errorf("you must be able to welcome at least one guest and have a bedroom")
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

//	if req.NumFullBathrooms > 0 {
//		numFullBathrooms = req.NumFullBathrooms
//	}
//	if req.NumHalfBathrooms > 0 {
//		numHalfBathrooms = req.NumHalfBathrooms
//	}
//	arg := db.CreateYatchParams{
//		NumFullBathrooms: int32(numFullBathrooms),
//		NumHalfBathrooms: int32(numHalfBathrooms),
//		NumOfGuest:       int32(req.NumOfGuest),
//		NumOfBedrooms:    int32(req.NumOfBedrooms),
//		OptionID:         option.ID,
//	}
//	yatch, err := server.store.CreateYatch(ctx, arg)
//	if err != nil {
//		log.Printf("There an error at CreateYatchInfo at CreateYatch: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
//		err = fmt.Errorf("please make sure you have at least a bedroom and can welcome a guest")
//		ctx.JSON(http.StatusBadRequest, errorResponse(err))
//		return
//	}
//	completeOption, err := HandleCompleteOption(utils.LocationView, utils.YatchInfo, server, ctx, option, user.ID)
//	if err != nil {
//		log.Printf("There an error at CreateYatchInfo at HandleCompleteOption: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
//	}
//	res := OptionInfoResponse{
//		OptionItemType:     "",
//		Success:            true,
//		OptionID:           tools.UuidToString(yatch.OptionID),
//		CurrentServerView:  completeOption.CurrentState,
//		PreviousServerView: completeOption.PreviousState,
//		MainOptionType:     option.MainOptionType,
//		OptionType:         option.OptionType,
//		Currency:           option.Currency,
//	}
//	ctx.JSON(http.StatusOK, res)
//}

//func (server *Server) CreateHallInfo(ctx *gin.Context) {
//	var req CreateHallParams

//	if err := ctx.ShouldBindJSON(&req); err != nil {
//		log.Printf("Error at CreateHallInfo in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
//		err = fmt.Errorf("you must be able to welcome at least one guest and have a bedroom")
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
//	if option.OptionType != "hall" {
//		err = fmt.Errorf("error occurred while performing your request. This experience type must be for a hall")
//		ctx.JSON(http.StatusForbidden, errorResponse(err))
//		return
//	}
//	err = server.store.RemoveHall(ctx, option.ID)
//	if err != nil {
//		if err == db.ErrorRecordNotFound {
//			log.Printf("hall doesn't exist %v \n", err.Error())
//		} else {
//			log.Println(err.Error())
//			err = fmt.Errorf("error occurred while performing your request. Try again")
//			ctx.JSON(http.StatusForbidden, errorResponse(err))
//			return

//		}
//	}
//	argHall := db.CreateHallParams{
//		OptionID:    option.ID,
//		HallLength:  req.HallHeight,
//		HallWidth:   req.HallWidth,
//		HallHeight:  req.HallHeight,
//		MaxNumGuest: int32(req.MaxNumGuest),
//	}
//	hall, err := server.store.CreateHall(ctx, argHall)
//	if err != nil {
//		log.Printf("There an error at CreateHallInfo at CreateHall: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
//		err = fmt.Errorf("please make you enter all the detail requested about your hall")
//		ctx.JSON(http.StatusBadRequest, errorResponse(err))
//		return
//	}
//	completeOption, err := HandleCompleteOption(utils.LocationView, utils.HallInfo, server, ctx, option, user.ID)
//	if err != nil {
//		log.Printf("There an error at CreateHallInfo at HandleCompleteOption: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
//	}
//	res := OptionInfoResponse{
//		OptionItemType:     "",
//		Success:            true,
//		OptionID:           tools.UuidToString(hall.OptionID),
//		CurrentServerView:  completeOption.CurrentState,
//		PreviousServerView: completeOption.PreviousState,
//		MainOptionType:     option.MainOptionType,
//		OptionType:         option.OptionType,
//		Currency:           option.Currency,
//	}
//	ctx.JSON(http.StatusOK, res)
//}

//func (server *Server) CreateLodgeInfo(ctx *gin.Context) {
//	var req CreateLodgeInfo
//	if err := ctx.ShouldBindJSON(&req); err != nil {
//		log.Printf("Error at CreateLodgeInfoType in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
//		err = fmt.Errorf("error occurred while processing your request. Please ensure you entered the right details")
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
//	if option.OptionType != "lodge" {
//		err = fmt.Errorf("error occurred while performing your request. This experience type must be for a hotel")
//		ctx.JSON(http.StatusForbidden, errorResponse(err))
//		return
//	}
//	err = server.store.RemoveLodgeInfo(ctx, option.ID)
//	if err != nil {
//		if err == db.ErrorRecordNotFound {
//			log.Printf("hotel doesn't exist %v \n", err.Error())
//		} else {
//			log.Println(err.Error())
//		}
//	}
//	_, err = server.store.CreateLodgeInfo(ctx, db.CreateLodgeInfoParams{
//		OptionID:  option.ID,
//		RoomTypes: req.RoomTypes,
//		LodgeType: option.OptionImg,
//	})
//	if err != nil {
//		log.Printf("Error at CreateLodgeInfoType in CreateLodgeInfo: %v, optionID: %v, userID: %v \n", err.Error(), requestID, user.ID)
//		err = fmt.Errorf("error occurred while creating your experience, try again")
//		ctx.JSON(http.StatusBadRequest, errorResponse(err))
//		return
//	}

//	completeOption, err := HandleCompleteOption(utils.LocationView, utils.LodgeRoomType, server, ctx, option, user.ID)
//	if err != nil {
//		log.Printf("Error at CreateLodgeInfo in HandleCompleteOption: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
//	}
//	res := OptionInfoResponse{
//		OptionItemType:     "",
//		Success:            true,
//		OptionID:           tools.UuidToString(option.ID),
//		CurrentServerView:  completeOption.CurrentState,
//		PreviousServerView: completeOption.PreviousState,
//		MainOptionType:     option.MainOptionType,
//		OptionType:         option.OptionType,
//		Currency:           option.Currency,
//	}
//	ctx.JSON(http.StatusOK, res)
//}

//func (server *Server) CreateRecreationInfo(ctx *gin.Context) {
//	var req CreateRecreationInfo
//	if err := ctx.ShouldBindJSON(&req); err != nil {
//		log.Printf("CreateRecreationInfo in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
//		err = fmt.Errorf("error occurred while processing your request. Please ensure you entered the right details")
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
//	if option.OptionType != "recreation" {
//		err = fmt.Errorf("error occurred while performing your request. This experience type must be for a recreation type")
//		ctx.JSON(http.StatusForbidden, errorResponse(err))
//		return
//	}
//	err = server.store.RemoveRecreationInfo(ctx, option.ID)
//	if err != nil {
//		if err == db.ErrorRecordNotFound {
//			log.Printf("recreation experience doesn't exist %v \n", err.Error())
//		} else {
//			log.Println(err.Error())
//		}
//	}
//	_, err = server.store.CreateRecreationInfo(ctx, db.CreateRecreationInfoParams{
//		OptionID:        option.ID,
//		RecreationTypes: req.RecreationTypes,
//	})
//	if err != nil {
//		log.Printf("Error at CreateRecreationInfo in CreateRecreationInfo: %v, optionID: %v, userID: %v \n", err.Error(), requestID, user.ID)
//		err = fmt.Errorf("error occurred while creating your recreation experience, try again")
//		ctx.JSON(http.StatusBadRequest, errorResponse(err))
//		return
//	}

//	completeOption, err := HandleCompleteOption(utils.LocationView, utils.RecreationType, server, ctx, option, user.ID)
//	if err != nil {
//		log.Printf("Error at CreateRecreationInfo in HandleCompleteOption: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
//	}
//	res := OptionInfoResponse{
//		OptionItemType:     "",
//		Success:            true,
//		OptionID:           tools.UuidToString(option.ID),
//		CurrentServerView:  completeOption.CurrentState,
//		PreviousServerView: completeOption.PreviousState,
//		MainOptionType:     option.MainOptionType,
//		OptionType:         option.OptionType,
//		Currency:           option.Currency,
//	}
//	ctx.JSON(http.StatusOK, res)
//}

//func (server *Server) CreateSportInfo(ctx *gin.Context) {
//	var req CreateSportInfo
//	if err := ctx.ShouldBindJSON(&req); err != nil {
//		log.Printf("CreateSportInfo in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
//		err = fmt.Errorf("error occurred while processing your request. Please ensure you entered the right details")
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
//	if option.OptionType != "Sport" {
//		err = fmt.Errorf("error occurred while performing your request. This experience type must be for a sport type")
//		ctx.JSON(http.StatusForbidden, errorResponse(err))
//		return
//	}
//	err = server.store.RemoveSportInfo(ctx, option.ID)
//	if err != nil {
//		if err == db.ErrorRecordNotFound {
//			log.Printf("Sport experience doesn't exist %v \n", err.Error())
//		} else {
//			log.Println(err.Error())
//			err = fmt.Errorf("error occurred while performing your request. Try again")
//			ctx.JSON(http.StatusForbidden, errorResponse(err))
//			return
//		}
//	}
//	_, err = server.store.CreateSportInfo(ctx, db.CreateSportInfoParams{
//		OptionID:   option.ID,
//		SportTypes: req.SportTypes,
//	})
//	if err != nil {
//		log.Printf("Error at CreateSportInfo in CreateSportInfo: %v, optionID: %v, userID: %v \n", err.Error(), requestID, user.ID)
//		err = fmt.Errorf("error occurred while creating your Sport experience, try again")
//		ctx.JSON(http.StatusBadRequest, errorResponse(err))
//		return
//	}

//	completeOption, err := HandleCompleteOption(utils.LocationView, utils.SportType, server, ctx, option, user.ID)
//	if err != nil {
//		log.Printf("Error at CreateSportInfo in HandleCompleteOption: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
//	}
//	res := OptionInfoResponse{
//		OptionItemType:     "",
//		Success:            true,
//		OptionID:           tools.UuidToString(option.ID),
//		CurrentServerView:  completeOption.CurrentState,
//		PreviousServerView: completeOption.PreviousState,
//		MainOptionType:     option.MainOptionType,
//		OptionType:         option.OptionType,
//		Currency:           option.Currency,
//	}
//	ctx.JSON(http.StatusOK, res)
//}

//func (server *Server) CreateChillInfo(ctx *gin.Context) {
//	var req CreateChillInfoParams
//	if err := ctx.ShouldBindJSON(&req); err != nil {
//		log.Printf("Error at CreateChillInfo in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
//		err = fmt.Errorf("error occurred while processing your request. Please ensure you entered the right details")
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
//	if option.OptionType != "chill" {
//		err = fmt.Errorf("error occurred while performing your request. This experience type must be for a restaurant, lounge, club")
//		ctx.JSON(http.StatusForbidden, errorResponse(err))
//		return
//	}
//	err = server.store.RemoveChillInfo(ctx, option.ID)
//	if err != nil {
//		if err == db.ErrorRecordNotFound {
//			log.Printf("chill info experience doesn't exist %v \n", err.Error())
//		} else {
//			log.Println(err.Error())
//		}
//	}
//	_, err = server.store.CreateChillInfo(ctx, db.CreateChillInfoParams{
//		OptionID:        option.ID,
//		ChillType:       option.OptionImg,
//		ComesWithFood:   req.ComesWithFood,
//		NumReservations: int32(req.NumReservations),
//	})
//	if err != nil {
//		log.Printf("Error at CreateChillInfo in CreateChillInfo: %v, optionID: %v, userID: %v \n", err.Error(), requestID, user.ID)
//		err = fmt.Errorf("error occurred while creating your %v experience", option.OptionImg)
//		ctx.JSON(http.StatusBadRequest, errorResponse(err))
//		return
//	}

//	completeOption, err := HandleCompleteOption(utils.LocationView, utils.ChillInfo, server, ctx, option, user.ID)
//	if err != nil {
//		log.Printf("Error at CreateChillInfo in HandleCompleteOption: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
//	}
//	res := OptionInfoResponse{
//		OptionItemType:     "",
//		Success:            true,
//		OptionID:           tools.UuidToString(option.ID),
//		CurrentServerView:  completeOption.CurrentState,
//		PreviousServerView: completeOption.PreviousState,
//		MainOptionType:     option.MainOptionType,
//		OptionType:         option.OptionType,
//		Currency:           option.Currency,
//	}
//	ctx.JSON(http.StatusOK, res)
//}
