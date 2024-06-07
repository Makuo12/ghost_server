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

func (server *Server) RemoveEventSubCategory(ctx *gin.Context) {
	var req OptionInfoRemoveRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at RemoveEventSubCategory in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
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
	event, err := server.store.UpdateEventInfo(ctx, db.UpdateEventInfoParams{
		OptionID: option.ID,
		SubCategoryType: pgtype.Text{
			String: "none",
			Valid:  true,
		},
	})
	if err != nil {
		log.Printf("Error at RemoveEventSubCategory in UpdateEventInfo(: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		err = fmt.Errorf("error occurred while controlling your state, please try again")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return

	}
	completeOption, err := HandleCompleteOption(utils.EventSubType, utils.OptionType, server, ctx, option, user.ID)
	if err != nil {
		log.Printf("Error at RemoveEventSubCategory in HandleCompleteOption: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
	}
	res := CreateOptionInfoParams{
		OptionItemType:     event.EventType,
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

//func (server *Server) RemoveEventLocation(ctx *gin.Context) {
//	var req OptionInfoRemoveRequest
//	if err := ctx.ShouldBindJSON(&req); err != nil {
//		log.Printf("Error at RemoveEventLocation in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
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
//	eventInfo, err := server.store.GetEventInfo(ctx, option.ID)
//	if err != nil {
//		log.Printf("Error at RemoveEventLocation in RemoveEventLocation: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
//		err = fmt.Errorf("error occurred while updating the state of your event")
//		ctx.JSON(http.StatusBadRequest, errorResponse(err))
//		return
//	}
//	locationType := eventInfo.LocationType
//	_, err = server.store.UpdateEventInfo(ctx, db.UpdateEventInfoParams{
//		OptionID: option.ID,
//		LocationType: pgtype.Text{
//			String: "none",
//			Valid:  true,
//		},
//	})
//	if err != nil {
//		log.Printf("Error at RemoveEventLocation in RemoveEventLocation: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
//		err = fmt.Errorf("error occurred while updating the state of your event")
//		ctx.JSON(http.StatusBadRequest, errorResponse(err))
//		return
//	}
//	err = server.store.RemoveLocation(ctx, option.ID)
//	if err != nil {
//		if err ==  {
//			// locationType it is important because if it RemoveLocation fails we can still get the past data
//			log.Printf("Location doesn't exist %v \n", err.Error())
//		} else {
//			log.Printf("Error at RemoveEventLocation in RemoveLocation: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
//			_, err = server.store.UpdateEventInfo(ctx, db.UpdateEventInfoParams{
//				OptionID: option.ID,
//				LocationType: pgtype.Text{
//					String: locationType,
//					Valid:  true,
//				},
//			})
//			if err != nil {
//				log.Printf("Error at RemoveEventLocation in UpdateEventInfo: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
//			}
//			err = fmt.Errorf("error occurred while control your state, please try again")
//			ctx.JSON(http.StatusNotFound, errorResponse(err))
//			return
//		}

//	}
//	currentState, previousState := utils.LocationReverseViewState(option.OptionType)
//	completeOption, err := HandleCompleteOption(currentState, previousState, server, ctx, option, user.ID)
//	if err != nil {
//		log.Printf("Error at RemoveLocation in HandleCompleteOption: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
//	}
//	res := LocationHandleEventParams(server, ctx, option, user.ID, completeOption)
//	ctx.JSON(http.StatusOK, res)
//}

func (server *Server) RemoveOptionPhoto(ctx *gin.Context) {
	var req OptionInfoRemoveRequest
	var photoExist bool = true
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at RemoveOption in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("error occurred while taking you back, try again")
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

	completeOption, err := HandleCompleteOption(utils.Photo, utils.Highlight, server, ctx, option, user.ID)
	if err != nil {
		log.Printf("Error at RemoveOptionName in HandleCompleteOption: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
	}
	_, err = server.store.GetOptionInfoPhoto(ctx, option.ID)
	if err != nil {
		if err == db.ErrorRecordNotFound {
			photoExist = false
		} else {
			err = fmt.Errorf("something went wrong while uploading your photos, try again")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
	}
	if photoExist {
		// we want to remove the photo from firebase
		err = RemoveAllPhoto(server, ctx, option)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
	}
	var res CreateOptionInfoHighlight
	optionDetail, err := server.store.GetOptionInfoDetail(ctx, option.ID)
	if err != nil {
		res = CreateOptionInfoHighlight{
			OptionID:           tools.UuidToString(option.ID),
			Highlight:          []string{""},
			UserOptionID:       tools.UuidToString(option.OptionUserID),
			CurrentServerView:  completeOption.CurrentState,
			PreviousServerView: completeOption.PreviousState,
			MainOptionType:     option.MainOptionType,
			OptionType:         option.OptionType,
			Currency:           option.Currency,
		}
	}
	res = CreateOptionInfoHighlight{
		OptionID:           tools.UuidToString(option.ID),
		Highlight:          optionDetail.OptionHighlight,
		UserOptionID:       tools.UuidToString(option.OptionUserID),
		CurrentServerView:  completeOption.CurrentState,
		PreviousServerView: completeOption.PreviousState,
		MainOptionType:     option.MainOptionType,
		OptionType:         option.OptionType,
		Currency:           option.Currency,
	}
	ctx.JSON(http.StatusOK, res)
}
