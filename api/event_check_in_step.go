package api

import (
	db "flex_server/db/sqlc"
	"flex_server/tools"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (server *Server) ListEventCheckInStep(ctx *gin.Context) {
	var req GetEventDateTimeParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at GetEventDateTimeParams in ShouldBindJSON: %v, EventDateTimeID: %v \n", err.Error(), req.EventDateTimeID)
		err = fmt.Errorf("an error occurred while setting up this amenity, try again")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Println(req)
	eventInfoID, err := tools.StringToUuid(req.EventInfoID)
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
	user, _, option, _, _, err := HandleGetCompleteOptionEditEventDateTimes(eventInfoID, ctx, server, true)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	eventDate, err := server.store.GetEventDateTimeByOption(ctx, db.GetEventDateTimeByOptionParams{
		ID:         eventDateTimeID,
		ID_2:       user.ID,
		ID_3:       option.ID,
		IsComplete: true,
	})
	if err != nil {
		log.Printf("There an error at ListEventCheckInStep at GetEventDateTimeByOption: %v, EventDateTimeID: %v, userID: %v \n", err.Error(), eventDateTimeID, user.ID)
		err = fmt.Errorf("could not set event date publish details, please try again using the format on the app")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	steps, err := server.store.ListEventCheckInStepOrdered(ctx, eventDateTimeID)
	if err != nil || len(steps) == 0 {
		if err != nil {
			log.Printf("There an error at ListEventCheckInStep at ListEventCheckInStepOrdered: %v, eventDateTimeID: %v, userID: %v \n", err.Error(), eventDateTimeID, user.ID)
		}
		res := "none"
		ctx.JSON(http.StatusNoContent, res)
		return
	}
	var res ListEventCheckInStepRes
	var resData []CheckInStepRes
	for i := 0; i < len(steps); i++ {
		data := CheckInStepRes{
			ID:    tools.UuidToString(steps[i].ID),
			Des:   steps[i].Des,
			Photo: steps[i].Photo,
		}
		resData = append(resData, data)
	}
	l, err := server.store.GetEventDateLocation(ctx, eventDateTimeID)
	if err != nil {
		res = ListEventCheckInStepRes{List: resData, Street: "", City: "", State: "", Country: "", Postcode: "", HasLocation: false, Published: eventDate.PublishCheckInSteps}

		ctx.JSON(http.StatusOK, res)
	} else {
		res = ListEventCheckInStepRes{List: resData, Street: l.Street, City: l.City, State: l.State, Country: l.Country, Postcode: l.Postcode, HasLocation: true, Published: eventDate.PublishCheckInSteps}
		ctx.JSON(http.StatusOK, res)
	}

}

func (server *Server) RemoveEventCheckInStepPhoto(ctx *gin.Context) {
	var req RemoveEventCheckInStepParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at RemoveCheckInStepPhotoParams in ShouldBindJSON: %v, eventDateTimeID: %v \n", err.Error(), req.EventDateTimeID)
		err = fmt.Errorf("an error occurred while setting up this amenity, try again")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Println(req)
	eventInfoID, err := tools.StringToUuid(req.EventInfoID)
	if err != nil {
		log.Printf("Error at tools.StringToUuid: %v, EventInfoID: %v \n", err.Error(), req.EventInfoID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	eventDateTimeID, err := tools.StringToUuid(req.EventDateTimeID)
	if err != nil {
		log.Printf("Error at tools.StringToUuid: %v, eventDateTimeID: %v \n", err.Error(), req.EventDateTimeID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	stepID, err := tools.StringToUuid(req.StepID)
	if err != nil {
		log.Printf("Error at tools.StringToUuid: %v, StepID: %v \n", err.Error(), req.StepID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, _, option, isCoHost, userCoHost, err := HandleGetCompleteOptionEditEventDateTimes(eventInfoID, ctx, server, true)
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
		log.Printf("There an error at RemoveEventCheckInStepPhoto at GetEventDateTimeByOption: %v, EventDateTimeID: %v, userID: %v \n", err.Error(), eventDateTimeID, user.ID)
		err = fmt.Errorf("could not set event date publish details, please try again using the format on the app")

		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// We want to remove photo from fire base
	stepDetail, err := server.store.GetEventCheckInStep(ctx, db.GetEventCheckInStepParams{
		ID:              stepID,
		EventDateTimeID: eventDateTimeID,
	})
	if err != nil {
		log.Printf("There an error at RemoveEventCheckInStepPhoto at GetCheckInStep: %v, eventDateTimeID: %v, userID: %v, stepID: %v \n", err.Error(), eventDateTimeID, user.ID, stepID)
		err = fmt.Errorf("could not find this step")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	err = RemoveFirebasePhoto(server, ctx, stepDetail.Photo)
	if err != nil {
		log.Printf("There an error at RemoveEventCheckInStepPhoto at RemoveFirebasePhoto: %v, eventDateTimeID: %v, userID: %v, stepID: %v \n", err.Error(), eventDateTimeID, user.ID, stepID)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	step, err := server.store.UpdateEventCheckInStepPhoto(ctx, db.UpdateEventCheckInStepPhotoParams{
		Photo:           "none",
		EventDateTimeID: eventDateTimeID,
		ID:              stepID,
	})
	if err != nil {
		log.Printf("There an error at RemoveEventCheckInStepPhoto at UpdateCheckInStepPhoto: %v, eventDateTimeID: %v, userID: %v \n", err.Error(), eventDateTimeID, user.ID)
		err = fmt.Errorf("your photo was deleted but not updated on the database, please if anything feels wrong just connect us")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if isCoHost {
		// we want to a push notification and store the message in the database
		HandleCoHostUpdateMsg(ctx, server, userCoHost, user, option, "ListEventCheckInStep", "event check in steps", "remove photo")
	}
	res := CheckInStepRes{
		ID:    tools.UuidToString(step.ID),
		Des:   tools.HandleString(step.Des),
		Photo: step.Photo,
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) RemoveEventCheckInStep(ctx *gin.Context) {
	var req RemoveEventCheckInStepParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at RemoveEventCheckInStepParams in ShouldBindJSON: %v, eventDateTimeID: %v \n", err.Error(), req.EventDateTimeID)
		err = fmt.Errorf("an error occurred while setting up this amenity, try again")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Println(req)
	eventInfoID, err := tools.StringToUuid(req.EventInfoID)
	if err != nil {
		log.Printf("Error at tools.StringToUuid: %v, EventInfoID: %v \n", err.Error(), req.EventInfoID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	eventDateTimeID, err := tools.StringToUuid(req.EventDateTimeID)
	if err != nil {
		log.Printf("Error at tools.StringToUuid: %v, eventDateTimeID: %v \n", err.Error(), req.EventDateTimeID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	stepID, err := tools.StringToUuid(req.StepID)
	if err != nil {
		log.Printf("Error at tools.StringToUuid: %v, StepID: %v \n", err.Error(), req.StepID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, _, option, isCoHost, userCoHost, err := HandleGetCompleteOptionEditEventDateTimes(eventInfoID, ctx, server, true)
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
		log.Printf("There an error at RemoveEventCheckInStep at GetEventDateTimeByOption: %v, EventDateTimeID: %v, userID: %v \n", err.Error(), eventDateTimeID, user.ID)
		err = fmt.Errorf("could not set event date publish details, please try again using the format on the app")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	// We want to remove photo from fire base
	stepDetail, err := server.store.GetEventCheckInStep(ctx, db.GetEventCheckInStepParams{
		ID:              stepID,
		EventDateTimeID: eventDateTimeID,
	})
	if err != nil {
		log.Printf("There an error at RemoveEventCheckInStep at GetCheckInStep: %v, eventDateTimeID: %v, userID: %v, stepID: %v \n", err.Error(), eventDateTimeID, user.ID, stepID)
		err = fmt.Errorf("could not find this step")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if len(stepDetail.Photo) != 0 && stepDetail.Photo != "none" {
		err = RemoveFirebasePhoto(server, ctx, stepDetail.Photo)
		if err != nil {
			log.Printf("There an error at RemoveEventCheckInStep at RemoveFirebasePhoto: %v, eventDateTimeID: %v, userID: %v, stepID: %v \n", err.Error(), eventDateTimeID, user.ID, stepID)
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
	}

	err = server.store.RemoveEventCheckInStep(ctx, db.RemoveEventCheckInStepParams{
		EventDateTimeID: eventDateTimeID,
		ID:              stepID,
	})
	if err != nil {
		log.Printf("There an error at RemoveEventCheckInStep at RemoveEventCheckInStep: %v, eventDateTimeID: %v, userID: %v \n", err.Error(), eventDateTimeID, user.ID)
		err = fmt.Errorf("the photo for this step was removed however something went wrong while updating it on the database. please refresh then connect us if anything feels wrong")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	steps, err := server.store.ListEventCheckInStepOrdered(ctx, eventDateTimeID)
	if err != nil || len(steps) == 0 {
		if err != nil {
			log.Printf("There an error at RemoveEventCheckInStep at ListEventCheckInStepOrdered: %v, eventDateTimeID: %v, userID: %v \n", err.Error(), eventDateTimeID, user.ID)
		}
		res := "none"
		ctx.JSON(http.StatusNoContent, res)
		return
	}
	var res ListCheckInStepRes
	var resData []CheckInStepRes
	for i := 0; i < len(steps); i++ {
		data := CheckInStepRes{
			ID:    tools.UuidToString(steps[i].ID),
			Des:   steps[i].Des,
			Photo: steps[i].Photo,
		}
		resData = append(resData, data)
	}
	res = ListCheckInStepRes{
		List: resData,
	}
	if isCoHost {
		HandleCoHostUpdateMsg(ctx, server, userCoHost, user, option, "ListEventCheckInStep", "event check in steps", "remove event check in step")
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) UpdateEventCheckInStep(ctx *gin.Context) {
	var req UpdateEventCheckInStepParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at UpdateEventCheckInStepParams in ShouldBindJSON: %v, eventDateTimeID: %v \n", err.Error(), req.EventDateTimeID)
		err = fmt.Errorf("an error occurred while setting up this amenity, try again")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Println(req)
	eventInfoID, err := tools.StringToUuid(req.EventInfoID)
	if err != nil {
		log.Printf("Error at tools.StringToUuid: %v, EventInfoID: %v \n", err.Error(), req.EventInfoID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	eventDateTimeID, err := tools.StringToUuid(req.EventDateTimeID)
	if err != nil {
		log.Printf("Error at tools.StringToUuid: %v, eventDateTimeID: %v \n", err.Error(), req.EventDateTimeID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	stepID, err := tools.StringToUuid(req.StepID)
	if err != nil {
		log.Printf("Error at tools.StringToUuid: %v, StepID: %v \n", err.Error(), req.StepID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, _, option, isCoHost, userCoHost, err := HandleGetCompleteOptionEditEventDateTimes(eventInfoID, ctx, server, true)
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
		log.Printf("There an error at UpdateEventCheckInStepParams at GetEventDateTimeByOption: %v, EventDateTimeID: %v, userID: %v \n", err.Error(), eventDateTimeID, user.ID)
		err = fmt.Errorf("could not set event date publish details, please try again using the format on the app")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var res CheckInStepRes
	switch req.Type {
	case "photo":
		step, err := server.store.UpdateEventCheckInStepPhoto(ctx, db.UpdateEventCheckInStepPhotoParams{
			Photo:           req.Photo,
			EventDateTimeID: eventDateTimeID,
			ID:              stepID,
		})
		if err != nil {
			log.Printf("There an error at UpdateEventCheckInStep at UpdateEventCheckInStepPhoto: %v, eventDateTimeID: %v, userID: %v \n", err.Error(), eventDateTimeID, user.ID)
			err = fmt.Errorf("could not update your photo in server")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		} else {
			res = CheckInStepRes{
				ID:    tools.UuidToString(step.ID),
				Des:   tools.HandleString(step.Des),
				Photo: step.Photo,
			}
		}
	case "des":
		step, err := server.store.UpdateEventCheckInStepDes(ctx, db.UpdateEventCheckInStepDesParams{
			Des:             req.Des,
			EventDateTimeID: eventDateTimeID,
			ID:              stepID,
		})
		if err != nil {
			log.Printf("There an error at UpdateEventCheckInStep at UpdateEventCheckInStepDes: %v, eventDateTimeID: %v, userID: %v \n", err.Error(), eventDateTimeID, user.ID)
			err = fmt.Errorf("could not update your des in server")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		} else {
			res = CheckInStepRes{
				ID:    tools.UuidToString(step.ID),
				Des:   tools.HandleString(step.Des),
				Photo: step.Photo,
			}
		}
	default:
		err = fmt.Errorf("type not found")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if isCoHost {
		HandleCoHostUpdateMsg(ctx, server, userCoHost, user, option, "ListEventCheckInStep", "event check in steps", "update check in step")
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) CreateEventCheckInStep(ctx *gin.Context) {
	var req CreateEventCheckInStepParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at CreateEventCheckInStepParams in CreateCheckInStep in ShouldBindJSON: %v, eventDateTimeID: %v \n", err.Error(), req.EventDateTimeID)
		err = fmt.Errorf("an error occurred while setting up this amenity, try again")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Println(req)
	eventInfoID, err := tools.StringToUuid(req.EventInfoID)
	if err != nil {
		log.Printf("Error at tools.StringToUuid: %v, EventInfoID: %v \n", err.Error(), req.EventInfoID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	eventDateTimeID, err := tools.StringToUuid(req.EventDateTimeID)
	if err != nil {
		log.Printf("Error at tools.StringToUuid: %v, eventDateTimeID: %v \n", err.Error(), req.EventDateTimeID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, _, option, isCoHost, userCoHost, err := HandleGetCompleteOptionEditEventDateTimes(eventInfoID, ctx, server, true)
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
		log.Printf("There an error at CreateEventCheckInStep at GetEventDateTimeByOption: %v, EventDateTimeID: %v, userID: %v \n", err.Error(), eventDateTimeID, user.ID)
		err = fmt.Errorf("could not set event date publish details, please try again using the format on the app")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var res CheckInStepRes
	switch req.Type {
	case "photo":
		step, err := server.store.CreateEventCheckInStep(ctx, db.CreateEventCheckInStepParams{
			EventDateTimeID: eventDateTimeID,
			Photo:           req.Photo,
			Des:             "none",
		})
		if err != nil {
			log.Printf("There an error at CreateEventCheckInStep at CreateEventCheckInStep for photo: %v, eventDateTimeID: %v, userID: %v \n", err.Error(), eventDateTimeID, user.ID)
			err = fmt.Errorf("could not update your photo in server")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		} else {
			res = CheckInStepRes{
				ID:    tools.UuidToString(step.ID),
				Des:   tools.HandleString(step.Des),
				Photo: step.Photo,
			}
		}
	case "des":
		step, err := server.store.CreateEventCheckInStep(ctx, db.CreateEventCheckInStepParams{
			EventDateTimeID: eventDateTimeID,
			Photo:           "none",
			Des:             req.Des,
		})
		if err != nil {
			log.Printf("There an error at CreateEventCheckInStep at CreateEventCheckInStep for des: %v, eventDateTimeID: %v, userID: %v \n", err.Error(), eventDateTimeID, user.ID)
			err = fmt.Errorf("could not update your photo in server")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		} else {
			res = CheckInStepRes{
				ID:    tools.UuidToString(step.ID),
				Des:   tools.HandleString(step.Des),
				Photo: step.Photo,
			}
		}
	default:
		err = fmt.Errorf("type not found")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if isCoHost {
		HandleCoHostUpdateMsg(ctx, server, userCoHost, user, option, "ListEventCheckInStep", "event check in steps", "create check in step")
	}
	ctx.JSON(http.StatusOK, res)
}
