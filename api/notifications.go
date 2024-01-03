package api

import (
	"context"
	"errors"
	db "flex_server/db/sqlc"
	"flex_server/tools"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateTypeNotification(ctx context.Context, server *Server, itemID uuid.UUID, userID uuid.UUID, notifyType string, msg string, itemIDFake bool, header string) {
	err := server.store.CreateNotification(ctx, db.CreateNotificationParams{
		ItemID:     itemID,
		ItemIDFake: itemIDFake,
		UserID:     userID,
		Type:       notifyType,
		Message:    msg,
		Header:     header,
	})
	if err != nil {
		log.Printf("Error at CreateTypeNotification in CreateNotification: %v, itemID: %v \n", err.Error(), itemID)
	}
	// We want to send an apn
	HandleUserIdApn(ctx, server, userID, tools.CapitalizeFirstCharacter(header), msg)
}

func (server *Server) ListNotification(ctx *gin.Context) {
	var req ListNotificationParams
	var onLastIndex = false
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at  ListNotification in ShouldBindJSON: %v, offset: %v \n", err.Error(), req.Offset)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := HandleGetUser(ctx, server)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	count, err := server.store.CountNotification(ctx, user.UserID)
	if err != nil {
		log.Printf("Error at  ListNotification in CountNotification err: %v, user: %v\n", err, user.ID)
		res := "none"
		ctx.JSON(http.StatusNoContent, res)
		return
	}
	if count <= int64(req.Offset) || count == 0 {
		res := "none"
		ctx.JSON(http.StatusNoContent, res)
		return
	}
	notifications, err := server.store.ListNotification(ctx, db.ListNotificationParams{
		UserID: user.UserID,
		Limit:  30,
		Offset: int32(req.Offset),
	})
	if err != nil {
		log.Printf("Error at  ListNotification in ListNotification err: %v, user: %v\n", err, user.ID)
		res := "none"
		ctx.JSON(http.StatusNoContent, res)
	}
	var res ListNotificationRes
	var resData []NotificationItem
	for _, n := range notifications {
		data := NotificationItem{
			ID:        tools.UuidToString(n.ID),
			Type:      n.Type,
			Header:    n.Header,
			Message:   n.Message,
			CreatedAt: tools.ConvertTimeToString(n.CreatedAt),
		}
		resData = append(resData, data)
	}
	log.Println("at 7")
	if count <= int64(req.Offset+len(notifications)) {
		onLastIndex = true
	}
	var timeString string = "none"
	if len(resData) > 0 {
		timeString = resData[0].CreatedAt
	}
	res = ListNotificationRes{
		List:        resData,
		Offset:      req.Offset + len(notifications),
		OnLastIndex: onLastIndex,
		UserID:      tools.UuidToString(user.UserID),
		Time:        timeString,
	}

	ctx.JSON(http.StatusOK, res)
}

func (server *Server) NotificationOptionReserveDetail(ctx *gin.Context) {
	var req NotificationOptionReserveDetailParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("error at NotificationOptionReserveDetail in ShouldBindJSON: %v \n", err)
		err = fmt.Errorf("there was an error while processing your inputs please try again later")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := HandleGetUser(ctx, server)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	nID, err := tools.StringToUuid(req.ID)
	if err != nil {
		log.Printf("error at NotificationOptionReserveDetail at tools.StringToUuid nID %v,: %v, userID: %v\n", nID, err.Error(), user.ID)
		err = errors.New("this listing does not exist")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	nData, err := server.store.GetNotificationUserRequest(ctx, db.GetNotificationUserRequestParams{
		NotificationID: nID,
		UserID:         user.UserID,
	})
	if err != nil {
		log.Printf("error at NotificationOptionReserveDetail at GetNotificationUserRequest nID %v,: %v, userID: %v\n", nID, err.Error(), user.ID)
		err = errors.New("this listing does not exist")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	charge, err := server.store.GetChargeOptionReference(ctx, db.GetChargeOptionReferenceParams{
		ID:               nData.ChargeID,
		UserID:           user.UserID,
		PaymentCompleted: false,
		ChargeCancelled:  false,
		RequestApproved:  true,
	})
	if err != nil {
		log.Printf("error at NotificationOptionReserveDetail at store.GetChargeOptionReference chargeID %v, err := tools.StringToUuid: %v, userID: %v\n", nData.ChargeID, err.Error(), user.ID)
		err = errors.New("this reservation you made was either not accepted or not found")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if time.Now().Add(time.Hour).After(charge.StartDate) || time.Now().Add(time.Hour).After(nData.CreatedAt.Add(time.Hour*48)) {
		err = errors.New("this reservation is no more available for booking because current date has passed your booking start date")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	defaultCardID, cardDetail, hasCard := HandleReserveCard(ctx, server, user, "NotificationOptionReserveDetail")
	reserveData, err := HandleOptionChargeToReserve(server, charge, req.Currency)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	// Experience Option Info Detail
	reqDetail := ExperienceDetailParams{
		OptionUserID:   reserveData.OptionUserID,
		MainOptionType: "options",
		Currency:       req.Currency,
	}
	resOption, hasData, err := HandleDetailOptionExperience(ctx, server, reqDetail)
	if err != nil || !hasData {
		if !hasData && err == nil {
			err = errors.New("listing details was not found")
		}
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	mainOptionData, err := HandleChargeToOptionData(ctx, server, charge, "NotificationOptionReserveDetail", req.Currency)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	// We update charge option reference cause the user is using a new currency
	err = UpdateChargeOptionReferencePrice(ctx, server, reserveData, user, "NotificationOptionReserveDetail")
	res := NotificationOptionReserveDetailRes{
		ReserveData:   reserveData,
		DefaultCardID: defaultCardID,
		HasCard:       hasCard,
		CardDetail:    cardDetail,
		OptionDetail:  resOption,
		Option:        mainOptionData,
	}
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Printf("NotificationOptionReserveDetail successfully (%v) id: %v\n", user.Email, user.ID)
	ctx.JSON(http.StatusOK, res)
}
