package api

import (
	"fmt"
	"log"
	"net/http"

	db "github.com/makuo12/ghost_server/db/sqlc"
	"github.com/makuo12/ghost_server/tools"
	"github.com/makuo12/ghost_server/val"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

func (server *Server) CreateUpdateEventDateDetail(ctx *gin.Context) {
	var req CreateUpdateEventDateDetailReq
	var dataExist bool
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at CreateUpdateEventDateDetailParams in ShouldBindJSON: %v, EventDateTimeID: %v \n", err.Error(), req.EventDateTimeID)
		err = fmt.Errorf("please make sure yuo select at least a start date")
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
	requestID, err := tools.StringToUuid(req.EventDateTimeID)
	if err != nil {
		log.Printf("Error at tools.StringToUuid: %v, EventDateTimeID: %v \n", err.Error(), req.EventDateTimeID)
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
		ID:         requestID,
		ID_2:       user.ID,
		ID_3:       option.ID,
		IsComplete: true,
	})
	if err != nil {
		log.Printf("There an error at CreateUpdateEventDateDetail at GetEventDateTimeByOption: %v, EventDateTimeID: %v, userID: %v \n", err.Error(), requestID, user.ID)
		err = fmt.Errorf("could not set event date details, please try again using the format on the app")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var startTimeString string
	var endTimeString string
	_, err = tools.ConvertStringToTimeOnly(req.StartTime)
	if err != nil {
		startTimeString = "none"
	} else {
		startTimeString = req.StartTime
	}
	_, err = tools.ConvertStringToTimeOnly(req.EndTime)
	if err != nil {
		endTimeString = "none"
	} else {
		endTimeString = req.EndTime
	}
	_, err = server.store.GetEventDateDetail(ctx, requestID)
	if err != nil {
		log.Printf("Error at  CreateUpdateEventDateDetail in GetEventDateDetail err: %v, user: %v\n", err, user.ID)
		dataExist = false
	} else {
		dataExist = true
	}
	var res CreateUpdateEventDateDetailRes
	var eventDateDetail db.EventDateDetail
	var name string
	if dataExist {
		// If data exist we want to update
		eventDateDetail, err = server.store.UpdateEventDateDetail(ctx, db.UpdateEventDateDetailParams{
			StartTime: pgtype.Text{
				String: startTimeString,
				Valid:  true,
			},
			EndTime: pgtype.Text{
				String: endTimeString,
				Valid:  true,
			},
			TimeZone: pgtype.Text{
				String: req.TimeZone,
				Valid:  true,
			},
			EventDateTimeID: requestID,
		})
		if err != nil {
			log.Printf("Error at  CreateUpdateEventDateDetail in UpdateEventDateDetail err: %v, user: %v\n", err, user.ID)
			err = fmt.Errorf("an error occurred while updating your event date details")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
	} else {
		// If data does not exist we want to update
		eventDateDetail, err = server.store.CreateEventDateDetail(ctx, db.CreateEventDateDetailParams{
			EventDateTimeID: requestID,
			StartTime:       startTimeString,
			EndTime:         endTimeString,
			TimeZone:        req.TimeZone,
		})
		if err != nil {
			log.Printf("Error at  CreateUpdateEventDateDetail in .CreateEventDateDetail err: %v, user: %v\n", err, user.ID)
			err = fmt.Errorf("an error occurred while creating your event date details")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
	}
	eventDateTime, err := server.store.UpdateEventDateTime(ctx, db.UpdateEventDateTimeParams{
		Name: pgtype.Text{
			String: req.Name,
			Valid:  true,
		},
		ID:   requestID,
		Type: req.Type,
	})
	if err != nil {
		log.Printf("Error at  CreateUpdateEventDateDetail in UpdateEventDateTime err: %v, user: %v\n", err, user.ID)
		name = ""
	} else {
		name = eventDateTime.Name
	}
	res = CreateUpdateEventDateDetailRes{
		EventDateTimeID: tools.UuidToString(eventDateDetail.EventDateTimeID),
		Name:            name,
		StartTime:       eventDateDetail.StartTime,
		EndTime:         eventDateDetail.EndTime,
		TimeZone:        eventDateDetail.TimeZone,
	}
	log.Println(res)
	if isCoHost {
		HandleCoHostUpdateMsg(ctx, server, userCoHost, user, option, "CreateUpdateEventDateDetail", "event date detail", "create-update event date detail")
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) GetEventDateDetail(ctx *gin.Context) {
	var req GetEventDateDetailParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at GetEventDateDetailParams in ShouldBindJSON: %v, EventDateTimeID: %v \n", err.Error(), req.EventDateTimeID)
		err = fmt.Errorf("please make sure yuo select at least a start date")
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
	requestID, err := tools.StringToUuid(req.EventDateTimeID)
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
	eventDT, err := server.store.GetEventDateTimeByOption(ctx, db.GetEventDateTimeByOptionParams{
		ID:         requestID,
		ID_2:       user.ID,
		ID_3:       option.ID,
		IsComplete: true,
	})
	if err != nil {
		log.Printf("There an error at GetEventDateDetail at GetEventDateTimeByOption: %v, EventDateTimeID: %v, userID: %v \n", err.Error(), requestID, user.ID)
		err = fmt.Errorf("could not set location for this event date, please try again using the format on the app")
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}
	eventDateDetail, err := server.store.GetEventDateDetail(ctx, requestID)
	if err != nil {
		if err == db.ErrorRecordNotFound {
			log.Printf("Error at GetEventDateDetail in GetEventDateDetail err: %v, user: %v\n", err, user.ID)
			res := "none"
			ctx.JSON(http.StatusNoContent, res)
			return
		}
		log.Printf("Error at GetEventDateDetail in GetEventDateDetail err: %v, user: %v\n", err, user.ID)
		err = fmt.Errorf("an error occurred while setting your event date details")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	res := CreateUpdateEventDateDetailRes{
		EventDateTimeID: tools.UuidToString(eventDateDetail.EventDateTimeID),
		Name:            eventDT.Name,
		StartTime:       eventDateDetail.StartTime,
		EndTime:         eventDateDetail.EndTime,
		TimeZone:        eventDateDetail.TimeZone,
	}
	log.Println(res)
	ctx.JSON(http.StatusOK, res)
}

// Publishes

func (server *Server) GetEventDatePublish(ctx *gin.Context) {
	var req GetEventDatePublishParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at GetEventDatePublishParams in ShouldBindJSON: %v, EventDateTimeID: %v \n", err.Error(), req.EventDateTimeID)
		err = fmt.Errorf("please make sure yuo select at least a start date")
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
	requestID, err := tools.StringToUuid(req.EventDateTimeID)
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

	_, err = server.store.GetEventDateTimeByOption(ctx, db.GetEventDateTimeByOptionParams{
		ID:         requestID,
		ID_2:       user.ID,
		ID_3:       option.ID,
		IsComplete: true,
	})
	if err != nil {
		log.Printf("There an error at GetEventDatePublish at GetEventDateTimeByOption: %v, EventDateTimeID: %v, userID: %v \n", err.Error(), requestID, user.ID)
		err = fmt.Errorf("could not set event date publish details, please try again using the format on the app")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	publish, err := server.store.GetEventDatePublish(ctx, requestID)
	if err != nil {
		log.Printf("There an error at GetEventDatePublish at GetEventDatePublish: %v, EventDateTimeID: %v, userID: %v \n", err.Error(), requestID, user.ID)
		err = fmt.Errorf("could not set event date publish details, please try again using the format on the app")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var resData []PrivateAudienceItem
	var hasAudience bool
	newList, err := server.store.ListEventDatePrivateAudience(ctx, requestID)
	if err != nil {
		log.Printf("There an error at GetEventDatePublish at ListEventDatePrivateAudience: %v, EventDateTimeID: %v, userID: %v \n", err.Error(), requestID, user.ID)
	}
	if len(newList) == 0 {
		hasAudience = false
		data := PrivateAudienceItem{Name: "none", Type: "none", Email: "none", Number: "none", Sent: false}
		resData = append(resData, data)
	} else {
		hasAudience = true
		for _, l := range newList {
			data := PrivateAudienceItem{ID: tools.UuidToString(l.ID), Name: l.Name, Email: l.Email, Number: l.Number, Sent: l.Sent, Exist: true, Type: l.Type}
			log.Println("data", data)
			resData = append(resData, data)
		}
	}
	res := UpdateEventDatePublishParams{
		EventDateTimeID:      tools.UuidToString(publish.EventDateTimeID),
		EventGoingPublicDate: tools.ConvertDateOnlyToString(publish.EventGoingPublicDate),
		EventGoingPublicTime: tools.ConvertTimeOnlyToString(publish.EventGoingPublicTime),
		EventPublic:          publish.EventPublic,
		EventGoingPublic:     publish.EventGoingPublic,
		HasAudience:          hasAudience,
		Audiences:            resData,
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) UpdatePrivateAudience(ctx *gin.Context) {
	var req UpdatePrivateAudienceParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at UpdatePrivateAudienceParams in ShouldBindJSON: %v, EventDateTimeID: %v \n", err.Error(), req.EventDateTimeID)
		err = fmt.Errorf("please make sure yuo select at least a start date")
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
	requestID, err := tools.StringToUuid(req.EventDateTimeID)
	if err != nil {
		log.Printf("Error at tools.StringToUuid: %v, EventDateTimeID: %v \n", err.Error(), req.EventDateTimeID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	id, err := tools.StringToUuid(req.ID)
	if err != nil {
		log.Printf("Error at tools.StringToUuid: %v, ID: %v \n", err.Error(), req.ID)
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
		ID:         requestID,
		ID_2:       user.ID,
		ID_3:       option.ID,
		IsComplete: true,
	})
	if err != nil {
		log.Printf("There an error at UpdatePrivateAudience at GetEventDateTimeByOption: %v, EventDateTimeID: %v, userID: %v \n", err.Error(), requestID, user.ID)
		err = fmt.Errorf("could not set event date publish details, please try again using the format on the app")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	// We want to first get the aud
	audData, err := server.store.GetEventDatePrivateAudience(ctx, db.GetEventDatePrivateAudienceParams{
		EventDateTimeID: requestID,
		ID:              id,
	})
	if err != nil {
		log.Printf("There an error at UpdatePrivateAudience at ListEventDatePrivateAudience: %v, EventDateTimeID: %v, userID: %v \n", err.Error(), requestID, user.ID)
		err = fmt.Errorf("you cannot update a user info that the message has already been sent to")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var aud db.EventDatePrivateAudience
	// We want to update the audience based on the type
	switch audData.Type {
	case "email":
		if val.ValidateEmail(req.Email) {
			aud, err = server.store.UpdateEventDatePrivateAudienceTwo(ctx, db.UpdateEventDatePrivateAudienceTwoParams{
				EventDateTimeID: audData.EventDateTimeID,
				ID:              audData.ID,
				Email: pgtype.Text{
					String: req.Email,
					Valid:  true,
				},
			})
		}
	case "phone":
		if val.ValidatePhoneNumber(req.Number) {
			aud, err = server.store.UpdateEventDatePrivateAudienceTwo(ctx, db.UpdateEventDatePrivateAudienceTwoParams{
				EventDateTimeID: audData.EventDateTimeID,
				ID:              audData.ID,
				Number: pgtype.Text{
					String: req.Number,
					Valid:  true,
				},
			})
		}

	}

	if err != nil {
		log.Printf("There an error at UpdatePrivateAudience at ListEventDatePrivateAudience: %v, EventDateTimeID: %v, userID: %v \n", err.Error(), requestID, user.ID)
		err = fmt.Errorf("you cannot update a user info that the message has already been sent to")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	res := PrivateAudienceItem{
		ID:     tools.UuidToString(aud.ID),
		Name:   aud.Name,
		Type:   aud.Type,
		Email:  aud.Email,
		Number: aud.Number,
		Sent:   aud.Sent,
		Exist:  true,
	}
	if isCoHost {
		HandleCoHostUpdateMsg(ctx, server, userCoHost, user, option, "UpdatePrivateAudience", "event private audience", "update event private audience")
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) RemovePrivateAudience(ctx *gin.Context) {
	var req RemovePrivateAudienceParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at RemovePrivateAudienceParams in ShouldBindJSON: %v, EventDateTimeID: %v \n", err.Error(), req.EventDateTimeID)
		err = fmt.Errorf("please make sure yuo select at least a start date")
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
	requestID, err := tools.StringToUuid(req.EventDateTimeID)
	if err != nil {
		log.Printf("Error at tools.StringToUuid: %v, EventDateTimeID: %v \n", err.Error(), req.EventDateTimeID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	id, err := tools.StringToUuid(req.ID)
	if err != nil {
		log.Printf("Error at tools.StringToUuid: %v, ID: %v \n", err.Error(), req.ID)
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
		ID:         requestID,
		ID_2:       user.ID,
		ID_3:       option.ID,
		IsComplete: true,
	})
	if err != nil {
		log.Printf("There an error at RemovePrivateAudience at GetEventDateTimeByOption: %v, EventDateTimeID: %v, userID: %v \n", err.Error(), requestID, user.ID)
		err = fmt.Errorf("could not set event date publish details, please try again using the format on the app")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	// We want to update the audience
	// We set Sent to false because we only want tp remove if the msg has not been sent
	err = server.store.RemoveEventDatePrivateAudienceBySent(ctx, db.RemoveEventDatePrivateAudienceBySentParams{
		EventDateTimeID: requestID,
		ID:              id,
		Sent:            false,
	})
	if err != nil {
		log.Printf("There an error at RemovePrivateAudience at RemoveEventDatePrivateAudience: %v, EventDateTimeID: %v, userID: %v \n", err.Error(), requestID, user.ID)
		err = fmt.Errorf("you cannot remove a user info that the message has already been sent to")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	res := "none"
	if isCoHost {
		HandleCoHostUpdateMsg(ctx, server, userCoHost, user, option, "UpdatePrivateAudience", "event private audience", "remove event private audience")
	}
	ctx.JSON(http.StatusNoContent, res)
}

func (server *Server) UpdateEventDatePublish(ctx *gin.Context) {
	var req UpdateEventDatePublishParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at UpdateEventDatePublishParams in ShouldBindJSON: %v, EventDateTimeID: %v \n", err.Error(), req.EventDateTimeID)
		err = fmt.Errorf("please make sure yuo select at least a start date")
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
	requestID, err := tools.StringToUuid(req.EventDateTimeID)
	if err != nil {
		log.Printf("Error at tools.StringToUuid: %v, EventDateTimeID: %v \n", err.Error(), req.EventDateTimeID)
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
		ID:         requestID,
		ID_2:       user.ID,
		ID_3:       option.ID,
		IsComplete: true,
	})
	if err != nil {
		log.Printf("There an error at UpdateEventDatePublish at GetEventDateTimeByOption: %v, EventDateTimeID: %v, userID: %v \n", err.Error(), requestID, user.ID)
		err = fmt.Errorf("could not set event date publish details, please try again using the format on the app")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var reqTime string
	if len(req.EventGoingPublicTime) == 0 || req.EventGoingPublic == "no" || req.EventPublic == "public" {
		reqTime = "00:00"
	} else {
		reqTime = req.EventGoingPublicTime
	}
	publishTime, err := tools.ConvertStringToTimeOnly(reqTime)
	log.Println("Show start time", publishTime)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var reqDate string
	if len(req.EventGoingPublicDate) == 0 || req.EventGoingPublic == "no" || req.EventPublic == "public" {
		reqDate = "1777-12-07"
	} else {
		reqDate = req.EventGoingPublicDate
	}
	publishDate, err := tools.ConvertDateOnlyStringToDate(reqDate)
	log.Println("Show end time", publishDate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	// We want to update the publish
	publish, err := server.store.UpdateEventDatePublish(ctx, db.UpdateEventDatePublishParams{
		EventPublic:          req.EventPublic,
		EventGoingPublic:     req.EventGoingPublic,
		EventGoingPublicDate: publishDate,
		EventGoingPublicTime: publishTime,
		EventDateTimeID:      requestID,
	})
	if err != nil {
		log.Printf("There an error at UpdateEventDatePublish at UpdateEventDatePublish: %v, EventDateTimeID: %v, userID: %v \n", err.Error(), requestID, user.ID)
		err = fmt.Errorf("could not set event date publish details, please try again using the format on the app")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	list, err := server.store.ListEventDatePrivateAudience(ctx, requestID)
	if err != nil {
		log.Printf("There an error at UpdateEventDatePublish at ListEventDatePrivateAudience: %v, EventDateTimeID: %v, userID: %v \n", err.Error(), requestID, user.ID)
		list = []db.ListEventDatePrivateAudienceRow{}
	}
	if len(req.Audiences) > 0 {
		for _, a := range req.Audiences {
			exist := false
			if !a.Exist {
				for _, l := range list {
					// we check if there is any match
					if a.Type == l.Type && (a.Email == l.Email || a.Number == l.Number) {
						exist = true
						break
					}
				}
				// If it doesn't exist then we create it
				if !exist {
					_, err = server.store.CreateEventDatePrivateAudience(ctx, db.CreateEventDatePrivateAudienceParams{
						EventDateTimeID: publish.EventDateTimeID,
						Name:            a.Name,
						Email:           a.Email,
						Number:          a.Number,
						Type:            a.Type,
					})
					if err != nil {
						log.Printf("There an error at UpdateEventDatePublish at CreateEventDatePrivateAudience: %v, EventDateTimeID: %v, userID: %v \n", err.Error(), requestID, user.ID)
					}
				}
			}
		}
	}
	var resData []PrivateAudienceItem
	var hasAudience bool
	newList, err := server.store.ListEventDatePrivateAudience(ctx, requestID)
	if err != nil {
		log.Printf("There an error at UpdateEventDatePublish at ListEventDatePrivateAudience: %v, EventDateTimeID: %v, userID: %v \n", err.Error(), requestID, user.ID)
	}
	if len(newList) == 0 {
		hasAudience = false
		data := PrivateAudienceItem{ID: "", Name: "none", Type: "none", Email: "none", Number: "none", Sent: false, Exist: false}
		resData = append(resData, data)
	} else {
		hasAudience = true
		for _, l := range newList {
			data := PrivateAudienceItem{ID: tools.UuidToString(l.ID), Name: l.Name, Email: l.Email, Number: l.Number, Sent: l.Sent, Exist: true, Type: l.Type}
			log.Println("data", data)
			resData = append(resData, data)
		}
	}
	res := UpdateEventDatePublishParams{
		EventDateTimeID:      tools.UuidToString(publish.EventDateTimeID),
		EventGoingPublicDate: tools.ConvertDateOnlyToString(publish.EventGoingPublicDate),
		EventGoingPublicTime: tools.ConvertTimeOnlyToString(publish.EventGoingPublicTime),
		EventPublic:          publish.EventPublic,
		EventGoingPublic:     publish.EventGoingPublic,
		HasAudience:          hasAudience,
		Audiences:            resData,
	}

	if isCoHost {
		HandleCoHostUpdateMsg(ctx, server, userCoHost, user, option, "UpdatePrivateAudience", "event private audience", "update event private audience")
	}

	ctx.JSON(http.StatusOK, res)
}
