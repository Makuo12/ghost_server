package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/makuo12/ghost_server/constants"
	db "github.com/makuo12/ghost_server/db/sqlc"
	"github.com/makuo12/ghost_server/tools"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

func (server *Server) ValidateOptionCoHost(ctx *gin.Context) {
	var req ValidateOptionCOHostParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at  ValidateOptionCOHostParams in ShouldBindJSON: %v, Code: %v \n", err.Error(), req.Code)
		err = fmt.Errorf("an error occurred while getting up this amenity detail, try again")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Println(req)

	user, err := HandleGetUser(ctx, server)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	result, err := RedisClient.Get(RedisContext, req.Code).Result()
	if err != nil {
		log.Printf("There an error at ValidateOptionCoHost at RedisClient.Get: %v, code: %v, userID: %v \n", err.Error(), req.Code, user.ID)
		err = fmt.Errorf("this code must have expired, try asking the host to resend the code")
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}

	coHostID, err := tools.StringToUuid(result)
	if err != nil {
		log.Printf("There an error at ValidateOptionCoHost at tools.StringToUuid: %v, code: %v, userID: %v \n", err.Error(), req.Code, user.ID)
		err = fmt.Errorf("this code must have expired, try asking the host to resend the code")
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}
	// Lets get the OPTION_INFO
	optionData, err := server.store.GetOptionCOHostByID(ctx, coHostID)
	if err != nil {
		log.Printf("There an error at ValidateOptionCoHost at store.GetOptionCOHostByID: %v, code: %v, userID: %v \n", err.Error(), req.Code, user.ID)
		err = fmt.Errorf("this code must have expired, try asking the host to resend the code")
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}
	if optionData.UserID == user.UserID {
		err = fmt.Errorf("you cannot cohost, when your the main host")
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}
	coHost, err := server.store.UpdateOptionCOHostTwo(ctx, db.UpdateOptionCOHostTwoParams{
		ID: coHostID,
		CoUserID: pgtype.Text{
			String: tools.UuidToString(user.UserID),
			Valid:  true,
		},
		Accepted: pgtype.Bool{
			Bool:  true,
			Valid: true,
		},
	})
	if err != nil {
		log.Printf("There an error at ValidateOptionCoHost at server.store.UpdateOptionCOHostTwo: %v, code: %v, userID: %v \n", err.Error(), req.Code, user.ID)
		err = fmt.Errorf("this code must have expired, try asking the host to resend the code")
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}
	err = RedisClient.Del(RedisContext, req.Code).Err()
	if err != nil {
		log.Printf("There an error at ValidateOptionCoHost at RedisClient.Del(RedisContext,: %v, code: %v, userID: %v \n", err.Error(), req.Code, user.ID)
	}
	option, err := server.store.GetOptionAndUser(ctx, coHost.OptionID)
	if err != nil {
		log.Printf("There an error at ValidateOptionCoHost at GetOptionAndUser: %v, code: %v, userID: %v \n", err.Error(), req.Code, user.ID)
		err = fmt.Errorf("this code must have expired, try asking the host to resend the code")
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}
	res := OptionCoHostItem{
		ID:             tools.UuidToString(coHost.ID),
		MainHostName:   option.FirstName,
		MainOptionName: option.HostNameOption,
		OptionCoHostID: tools.UuidToString(option.CoHostID),
		CoverImage:     option.CoverImage,
		MainOption:     option.MainOptionType,
		IsPrimaryHost:  false,
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) ListOptionCoHostItem(ctx *gin.Context) {
	var req ListOptionCoHostItemParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at  ValidateOptionCOHostParams in ShouldBindJSON: %v, Offset: %v \n", err.Error(), req.Offset)
		err = fmt.Errorf("an error occurred while getting up this amenity detail, try again")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Println(req)

	user, err := HandleGetUser(ctx, server)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	count, err := server.store.CountOptionCoHostByCoHost(ctx, tools.UuidToString(user.UserID))
	if err != nil {
		log.Printf("Error at  ListOptionCoHostItem in store.CountOptionCoHostByCoHost err: %v, user: %v\n", err, user.ID)
		res := "none"
		ctx.JSON(http.StatusNoContent, res)
		return
	}
	if count <= int64(req.Offset) || count == 0 {
		res := "none"
		ctx.JSON(http.StatusNoContent, res)
		return
	}
	options, err := server.store.ListOptionCoHostByCoHost(ctx, db.ListOptionCoHostByCoHostParams{
		Limit:    15,
		Offset:   int32(req.Offset),
		CoUserID: tools.UuidToString(user.UserID),
	})
	if err != nil {
		log.Printf("Error at  ListOptionCoHostItem in store.ListOptionCoHostByCoHost err: %v, user: %v\n", err, user.ID)
		res := "none"
		ctx.JSON(http.StatusNoContent, res)
	}

	var resData []OptionCoHostItem
	for _, o := range options {
		var isPrimaryHost bool
		if o.PrimaryUserID == user.UserID {
			isPrimaryHost = true
		}
		data := OptionCoHostItem{
			ID:             tools.UuidToString(o.ID),
			MainHostName:   o.FirstName,
			MainOptionName: o.HostNameOption,
			OptionCoHostID: tools.UuidToString(o.CoHostID),
			CoverImage:     o.CoverImage,
			MainOption:     o.MainOptionType,
			IsPrimaryHost:  isPrimaryHost,
		}
		resData = append(resData, data)
	}
	res := ListOptionCoHostItemRes{
		List: resData,
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) GetOptionCoHostItemDetail(ctx *gin.Context) {
	var req OptionCoHostItemDetailParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at  OptionCoHostItemDetailParams in ShouldBindJSON: %v, ID: %v \n", err.Error(), req.ID)
		err = fmt.Errorf("an error occurred while getting up this amenity detail, try again")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Println(req)

	user, err := HandleGetUser(ctx, server)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	id, err := tools.StringToUuid(req.ID)
	if err != nil {
		log.Printf("There an error at DeactivateOptionCoHost at tools.StringToUuid: %v, ID: %v, userID: %v \n", err.Error(), req.ID, user.ID)
		err = fmt.Errorf("you don't have access to this resource")
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}
	coHost, err := server.store.GetOptionCOHostByCoHost(ctx, db.GetOptionCOHostByCoHostParams{
		ID:       id,
		CoUserID: tools.UuidToString(user.UserID),
	})
	if err != nil {
		log.Printf("There an error at GetOptionCOHostByCoHost at store.DeactivateCoHost: %v, ID: %v, userID: %v \n", err.Error(), req.ID, user.ID)
		err = fmt.Errorf("you don't have access to this resource")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	res := OptionCoHostItemDetailRes{
		Reservations:        coHost.Reservations,
		Post:                coHost.Post,
		ScanCode:            coHost.ScanCode,
		Calender:            coHost.Calender,
		EditOptionInfo:      coHost.EditOptionInfo,
		EditEventDatesTimes: coHost.EditEventDatesTimes,
		EditCoHosts:         coHost.EditCoHosts,
		Insights:            coHost.Insights,
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) DeactivateOptionCoHost(ctx *gin.Context) {
	var req DeactivateOptionCOHostParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at  DeactivateOptionCOHostParams in ShouldBindJSON: %v, ID: %v \n", err.Error(), req.ID)
		err = fmt.Errorf("an error occurred while getting up this amenity detail, try again")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Println(req)

	user, err := HandleGetUser(ctx, server)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	id, err := tools.StringToUuid(req.ID)
	if err != nil {
		log.Printf("There an error at DeactivateOptionCoHost at tools.StringToUuid: %v, ID: %v, userID: %v \n", err.Error(), req.ID, user.ID)
		err = fmt.Errorf("you don't have access to this resource")
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}
	coHostID, err := server.store.DeactivateCoHost(ctx, db.DeactivateCoHostParams{
		CoUserID: tools.UuidToString(user.UserID),
		ID:       id,
	})
	if err != nil {
		log.Printf("There an error at DeactivateOptionCoHost at store.DeactivateCoHost: %v, ID: %v, userID: %v \n", err.Error(), req.ID, user.ID)
		err = fmt.Errorf("you don't have access to this resource")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	err = RedisClient.SAdd(RedisContext, constants.DEACTIVATE_CO_HOST_IDS, tools.UuidToString(coHostID)).Err()
	if err != nil {
		log.Printf("There an error at DeactivateOptionCoHost at RedisClient.SAdd: %v, ID: %v, userID: %v \n", err.Error(), req.ID, user.ID)
		err = nil
	}

	res := DeactivateOptionCOHostRes{
		ID:      req.ID,
		Success: true,
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) UpdatePublishOptionCheckInStep(ctx *gin.Context) {
	var req UpdatePublishCheckInStepParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at UpdatePublishOptionCheckInStep in ShouldBindJSON: %v, OptionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("please make sure yuo select at least a start date")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Println(req)
	optionID, err := tools.StringToUuid(req.OptionID)
	if err != nil {
		log.Printf("Error at tools.StringToUuid: %v, OptionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, _, option, isCoHost, userCoHost, err := HandleGetCompleteOptionEditOptionInfo(optionID, ctx, server, true)
	if err != nil {
		err = fmt.Errorf("you cannot access this resource")
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}
	published, err := server.store.UpdateShortletPublishCheckInStep(ctx, option.ID)
	if err != nil {
		log.Printf("There an error at UpdatePublishOptionCheckInStep at UpdateShortletPublishCheckInStep: %v, option.ID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		err = fmt.Errorf("could not publish your check in step")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	res := UpdatePublishCheckInStepRes{
		Published: published,
	}
	log.Println(res)
	if isCoHost {
		HandleCoHostUpdateMsg(ctx, server, userCoHost, user, option, "UpdatePublishOptionCheckInStep", "listing check in step", "update listing check in step")
	}
	ctx.JSON(http.StatusOK, res)
}
