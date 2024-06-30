package api

import (
	"fmt"
	"log"
	"net/http"

	db "github.com/makuo12/ghost_server/db/sqlc"
	"github.com/makuo12/ghost_server/tools"
	"github.com/makuo12/ghost_server/utils"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

func (server *Server) RemoveOptionSE(ctx *gin.Context) {
	var req OptionInfoRemoveRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at RemoveOption in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("error occurred while taking you back, try again")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Println(req)

	requestID, err := tools.StringToUuid(req.OptionID)

	if err != nil {
		var isHost bool
		user, err := HandleGetUser(ctx, server)
		userIsHost, hasIncomplete := UserIsHost(ctx, server, user)
		if err != nil {
			isHost = false
		} else {
			isHost = userIsHost
		}
		res := OptionInfoRemoveFirstResponse{
			Success:       true,
			UserIsHost:    isHost,
			HasIncomplete: hasIncomplete,
		}
		ctx.JSON(http.StatusOK, res)
		return
	}

	user, option, err := HandleGetOptionIncomplete(requestID, ctx, server, false)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// if option is an event we want to remove the event info

	err = server.store.RemoveCompleteOptionInfo(ctx, option.ID)
	if err != nil {
		log.Printf("error at RemoveOption at RemoveCompleteOptionInfo: %v, optionID: %v, userID: %vv\n", err.Error(), option.ID, user.ID)
		err = fmt.Errorf("error occurred while taking you back, please try again")
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}
	if option.MainOptionType == "events" {
		err = server.store.RemoveEventInfo(ctx, option.ID)
		if err != nil {
			log.Printf("Error at HandleEventInfo for optionID: %v, userID: %v, err is: %v", option.ID, user.ID, err.Error())
			err = fmt.Errorf("an error occurred while removing your event, try again")
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		HandleRemoveEventDefaultTables(server, ctx, option, user)
	} else {
		// we want to remove the shortlets
		err = server.store.RemoveShortlet(ctx, option.ID)
		if err != nil {
			log.Printf("error at RemoveOption at RemoveShortlet: %v, optionID: %v, userID: %vv\n", err.Error(), option.ID, user.ID)
			err = fmt.Errorf("error occurred while taking you back, please try again")
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		HandleRemoveOptionDefaultTables(server, ctx, option, user)
	}

	arg := db.RemoveOptionInfoParams{
		ID:     requestID,
		HostID: user.ID,
	}

	err = server.store.RemoveOptionInfo(ctx, arg)
	if err != nil {
		if err == db.ErrorRecordNotFound {
			log.Printf("Option doesn't exist %v \n", err.Error())
		} else {
			log.Printf("Error at RemoveOption in RemoveOptionInfo: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
			err = fmt.Errorf("error occurred while taking you back, please try again")
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

	}
	userIsHost, hasIncomplete := UserIsHost(ctx, server, user)
	res := OptionInfoRemoveFirstResponse{
		Success:       true,
		UserIsHost:    userIsHost,
		HasIncomplete: hasIncomplete,
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) RemoveShortletSpace(ctx *gin.Context) {
	var req OptionInfoRemoveRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at RemoveShortletSpace in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("error occurred while taking you back")
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

	arg := db.UpdateShortletInfoParams{
		SpaceType: pgtype.Text{
			String: "none",
			Valid:  true,
		},
		OptionID: option.ID,
	}
	shortlet, err := server.store.UpdateShortletInfo(ctx, arg)
	if err != nil {
		log.Printf("Error at RemoveShortletSpace in UpdateShortletInfo: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		err = fmt.Errorf("error occurred while taking you back, please try again")
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}
	completeOption, err := HandleCompleteOption(utils.ShortletSpace, utils.OptionType, server, ctx, option, user.ID)
	if err != nil {
		log.Printf("Error at RemoveShortletSpace in HandleCompleteOption: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
	}
	res := CreateOptionInfoParams{
		OptionItemType:     shortlet.TypeOfShortlet,
		OptionType:         option.OptionType,
		Currency:           option.Currency,
		MainOptionType:     option.MainOptionType,
		OptionImg:          option.OptionImg,
		OptionID:           tools.UuidToString(option.ID),
		UserOptionID:       tools.UuidToString(option.OptionUserID),
		CurrentServerView:  completeOption.CurrentState,
		PreviousServerView: completeOption.PreviousState,
	}

	ctx.JSON(http.StatusOK, res)
}

func (server *Server) RemoveShortletInfo(ctx *gin.Context) {
	var req OptionInfoRemoveRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at RemoveShortletType in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("error occurred while processing your request. Please ensure you entered the right details")
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
	err = server.store.RemoveSpaceAreaAll(ctx, option.ID)
	if err != nil {
		log.Printf("An error at RemoveShortletInfo in RemoveSpaceAreaAll err: %v for optionID: %q", err.Error(), option.ID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	argUpdate := db.UpdateShortletInfoParams{
		GuestWelcomed: pgtype.Int4{
			Int32: int32(0),
			Valid: true,
		},
		OptionID: option.ID,
	}
	// wev change guest number back to 0

	_, err = server.store.UpdateShortletInfo(ctx, argUpdate)
	if err != nil {
		log.Printf("There an error at RemoveShortletInfo at UpdateShortletInfo: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
	}
	completeOption, err := HandleCompleteOption(utils.ShortletInfo, utils.ShortletSpace, server, ctx, option, user.ID)
	if err != nil {
		log.Printf("There an error at RemoveShortletInfo at HandleCompleteOption: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
	}
	// Get the shortletSpace data
	var res CreateShortletSpaceParams
	shortlet, err := server.store.GetShortlet(ctx, option.ID)
	if err != nil {
		res = CreateShortletSpaceParams{
			OptionID:           tools.UuidToString(option.ID),
			SpaceType:          "none",
			UserOptionID:       tools.UuidToString(option.OptionUserID),
			CurrentServerView:  completeOption.CurrentState,
			PreviousServerView: completeOption.PreviousState,
			MainOptionType:     option.MainOptionType,
			OptionType:         option.OptionType,
			Currency:           option.Currency,
		}
	} else {
		res = CreateShortletSpaceParams{
			OptionID:           tools.UuidToString(option.ID),
			SpaceType:          shortlet.SpaceType,
			UserOptionID:       tools.UuidToString(option.OptionUserID),
			CurrentServerView:  completeOption.CurrentState,
			PreviousServerView: completeOption.PreviousState,
			MainOptionType:     option.MainOptionType,
			OptionType:         option.OptionType,
			Currency:           option.Currency,
		}
	}

	ctx.JSON(http.StatusOK, res)
}

func (server *Server) RemoveLocation(ctx *gin.Context) {
	var req OptionInfoRemoveRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at RemoveLocation in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("error occurred while processing your request. Please ensure you entered the right details")
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

	err = server.store.RemoveLocation(ctx, option.ID)
	if err != nil {
		if err == db.ErrorRecordNotFound {
			// the reason for this is that it is possible that this amenities was never created
			log.Printf("Location doesn't exist %v \n", err.Error())
		} else {
			log.Printf("Error at RemoveLocation in RemoveLocation: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
			err = fmt.Errorf("error occurred while control your state, please try again")
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

	}
	currentState, previousState := utils.LocationReverseViewState(option.OptionType)
	completeOption, err := HandleCompleteOption(currentState, previousState, server, ctx, option, user.ID)
	if err != nil {
		log.Printf("Error at RemoveLocation in HandleCompleteOption: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
	}
	res := LocationHandleShortletParams(server, ctx, option, user.ID, completeOption)
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) RemoveAmenitiesAndSafety(ctx *gin.Context) {
	var req OptionInfoRemoveRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at RemoveLocation in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("error occurred while taking you back")
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

	err = server.store.RemoveAllAmenity(ctx, option.ID)
	if err != nil {
		if err == db.ErrorRecordNotFound {
			// the reason for this is that it is possible that this amenities was never created
			log.Printf("Amenities doesn't exist %v \n", err.Error())
		} else {
			log.Printf("Error at RemoveAmenitiesAndSafety in RemoveAllAmenity: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
			err = fmt.Errorf("error occurred while control your state, please try again")
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
	}
	completeOption, err := HandleCompleteOption(utils.Amenities, utils.LocationView, server, ctx, option, user.ID)
	if err != nil {
		log.Printf("Error at RemoveAmenitiesAndSafety in HandleCompleteOption: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
	}
	res := HandleAmenityLocationData(server, ctx, option, completeOption)
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) RemoveOptionDescription(ctx *gin.Context) {
	var req OptionInfoRemoveRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at RemoveOptionDescription in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("error occurred while processing your request. Please ensure you entered the right details")
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
	_, err = server.store.UpdateOptionInfoDetail(ctx, db.UpdateOptionInfoDetailParams{
		OptionID: option.ID,
		Des: pgtype.Text{
			String: "none",
			Valid:  false,
		},
	})
	if err != nil {
		// We expect it to have already been created
		log.Printf("Error at RemoveOptionDescription in UpdateOptionInfoDetail: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		err = fmt.Errorf("error occurred while control your state, please try again")
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return

	}
	currentState, previousState := utils.DescriptionReserveViewState(option.OptionType)
	completeOption, err := HandleCompleteOption(currentState, previousState, server, ctx, option, user.ID)
	if err != nil {
		log.Printf("Error at RemoveOptionDescription in HandleCompleteOption: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
	}
	switch option.OptionType {
	case "shortlets":
		fmt.Printf("at shortlets for des")
		res := DescriptionHandleAmenitiesParams(server, ctx, option, user.ID, completeOption)
		ctx.JSON(http.StatusOK, res)
		return
	case "events":
		fmt.Printf("at events for des")
		res := DescriptionHandleEventParams(server, ctx, option, user.ID, completeOption)
		ctx.JSON(http.StatusOK, res)
		return
	default:
		err = fmt.Errorf("option type was not found")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return

	}

}

func (server *Server) RemoveOptionName(ctx *gin.Context) {
	var req OptionInfoRemoveRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at RemoveOptionName in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("error occurred while taking you back")
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
	_, err = server.store.UpdateOptionInfoDetail(ctx, db.UpdateOptionInfoDetailParams{
		OptionID: option.ID,
		HostNameOption: pgtype.Text{
			String: "none",
			Valid:  false,
		},
	})
	if err != nil {
		// We expect it to have already been created
		log.Printf("Error at RemoveOptionName in UpdateOptionInfoDetail: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		err = fmt.Errorf("error occurred while control your state, please try again")
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return

	}

	completeOption, err := HandleCompleteOption(utils.Name, utils.Description, server, ctx, option, user.ID)
	if err != nil {
		log.Printf("Error at RemoveOptionName in HandleCompleteOption: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
	}
	var res CreateOptionInfoDescription
	optionDetail, err := server.store.GetOptionInfoDetail(ctx, option.ID)
	if err != nil {
		res = CreateOptionInfoDescription{
			OptionID:           tools.UuidToString(option.ID),
			Description:        "none",
			UserOptionID:       tools.UuidToString(option.OptionUserID),
			CurrentServerView:  completeOption.CurrentState,
			PreviousServerView: completeOption.PreviousState,
			MainOptionType:     option.MainOptionType,
			OptionType:         option.OptionType,
			Currency:           option.Currency,
		}
	} else {
		res = CreateOptionInfoDescription{
			OptionID:           tools.UuidToString(option.ID),
			Description:        optionDetail.Des,
			UserOptionID:       tools.UuidToString(option.OptionUserID),
			CurrentServerView:  completeOption.CurrentState,
			PreviousServerView: completeOption.PreviousState,
			MainOptionType:     option.MainOptionType,
			OptionType:         option.OptionType,
			Currency:           option.Currency,
		}
	}
	ctx.JSON(http.StatusOK, res)

}
func (server *Server) RemoveOptionHighlight(ctx *gin.Context) {
	var req OptionInfoRemoveRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at RemoveOptionHighlight in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("error occurred while processing your request. Please ensure you entered the right details")
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
	_, err = server.store.UpdateOptionInfoDetail(ctx, db.UpdateOptionInfoDetailParams{
		OptionID:        option.ID,
		OptionHighlight: []string{""},
	})
	if err != nil {
		log.Printf("Error at RemoveOptionHighlight in UpdateOptionInfoDetail: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		err = fmt.Errorf("error occurred while controlling your state, please try again")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return

	}
	currentState, previousState := utils.HighlightReserveViewState(option.MainOptionType, option.OptionType)
	completeOption, err := HandleCompleteOption(currentState, previousState, server, ctx, option, user.ID)
	if err != nil {
		log.Printf("Error at RemoveOptionHighlight in HandleCompleteOption: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
	}
	switch option.OptionType {
	case "shortlets":
		res := HighlightHandlePriceParams(server, ctx, option, user.ID, completeOption)
		ctx.JSON(http.StatusOK, res)
		return
	case "events":
		res := HighlightHandleNameParams(server, ctx, option, user.ID, completeOption)
		ctx.JSON(http.StatusOK, res)
		return
	default:
		err = fmt.Errorf("option type was not found")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

}

func (server *Server) RemoveOptionPrice(ctx *gin.Context) {
	var req OptionInfoRemoveRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at RemoveOptionPrice in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("error occurred while processing your request. Please ensure you entered the right details")
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
	err = server.store.RemoveOptionPrice(ctx, option.ID)
	if err != nil {
		if err == db.ErrorRecordNotFound {
			log.Printf("Shortlet doesn't exist %v \n", err.Error())
		} else {
			log.Printf("Error at RemoveOptionPrice in RemoveOptionPrice: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
			err = fmt.Errorf("error occurred while control your state, please try again")
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

	}
	completeOption, err := HandleCompleteOption(utils.Price, utils.Name, server, ctx, option, user.ID)
	if err != nil {
		log.Printf("Error at RemoveOptionPrice in HandleCompleteOption: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
	}
	var res CreateOptionInfoName
	optionDetail, err := server.store.GetOptionInfoDetail(ctx, option.ID)
	if err != nil {
		res = CreateOptionInfoName{
			OptionID:           tools.UuidToString(option.ID),
			HostNameOption:     "none",
			UserOptionID:       tools.UuidToString(option.OptionUserID),
			CurrentServerView:  completeOption.CurrentState,
			PreviousServerView: completeOption.PreviousState,
			MainOptionType:     option.MainOptionType,
			OptionType:         option.OptionType,
			Currency:           option.Currency,
		}
	} else {
		res = CreateOptionInfoName{
			OptionID:           tools.UuidToString(option.ID),
			HostNameOption:     optionDetail.HostNameOption,
			UserOptionID:       tools.UuidToString(option.OptionUserID),
			CurrentServerView:  completeOption.CurrentState,
			PreviousServerView: completeOption.PreviousState,
			MainOptionType:     option.MainOptionType,
			OptionType:         option.OptionType,
			Currency:           option.Currency,
		}
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) RemoveOptionQuestion(ctx *gin.Context) {
	var req OptionInfoRemoveRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at RemoveOptionQuestion in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("error occurred while taking you back")
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

	err = server.store.RemoveOptionQuestion(ctx, option.ID)
	if err != nil {
		if err == db.ErrorRecordNotFound {
			log.Printf("Questions section not accessible %v \n", err.Error())
		} else {
			log.Printf("Error at RemoveOptionQuestion in RemoveOptionQuestion: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
			err = fmt.Errorf("error occurred while taking you back, please try again")
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

	}
	err = server.store.RemoveAllThingToNote(ctx, option.ID)
	if err != nil {
		if err == db.ErrorRecordNotFound {
			log.Printf("RemoveAllThingToNote section not accessible %v \n", err.Error())
		} else {
			log.Printf("Error at RemoveOptionQuestion in RemoveAllThingToNote: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		}
	}

	completeOption, err := HandleCompleteOption(utils.HostQuestion, utils.Photo, server, ctx, option, user.ID)
	if err != nil {
		log.Printf("Error at HandleCompleteOption in RemoveOptionQuestion: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
	}
	photoData, err := server.store.GetOptionInfoPhoto(ctx, option.ID)
	var res CreateOptionInfoPhotoParams

	if err != nil {
		res = CreateOptionInfoPhotoParams{
			OptionID:           tools.UuidToString(option.ID),
			UserOptionID:       tools.UuidToString(option.OptionUserID),
			CurrentServerView:  completeOption.CurrentState,
			PreviousServerView: completeOption.PreviousState,
			MainOptionType:     option.MainOptionType,
			OptionType:         option.OptionType,
			Currency:           option.Currency,
			MainPhoto:          []string{},
			MainCoverImage:     "",
			PublicPhoto:        []string{},
			PublicCoverImage:   "",
		}
	} else {
		res = CreateOptionInfoPhotoParams{
			OptionID:           tools.UuidToString(option.ID),
			UserOptionID:       tools.UuidToString(option.OptionUserID),
			CurrentServerView:  completeOption.CurrentState,
			PreviousServerView: completeOption.PreviousState,
			MainOptionType:     option.MainOptionType,
			OptionType:         option.OptionType,
			Currency:           option.Currency,
			MainPhoto:          photoData.Photo,
			MainCoverImage:     photoData.CoverImage,
			PublicPhoto:        photoData.PublicPhoto,
			PublicCoverImage:   photoData.PublicCoverImage,
		}
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) RemovePublishOption(ctx *gin.Context) {
	var req OptionInfoRemoveRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at RemovePublishOption in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("error occurred while taking you back")
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
	currentState, previousState := utils.PublishReverseViewState(option.MainOptionType, option.OptionType)
	completeOption, err := HandleCompleteOption(currentState, previousState, server, ctx, option, user.ID)
	if err != nil {
		log.Printf("Error at HandleCompleteOption in RemovePublishOption: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
	}
	switch option.OptionType {
	case "shortlets":
		res := PublishHandleHostQuestion(server, ctx, option, user.ID, completeOption)
		ctx.JSON(http.StatusOK, res)
		return
	case "events":
		res := ReversePublishToPhoto(server, ctx, option, user.ID, completeOption)
		ctx.JSON(http.StatusOK, res)
	default:
		err = fmt.Errorf("option type was not found")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
}
