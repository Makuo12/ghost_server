package api

import (
	"context"
	"flex_server/constants"
	db "flex_server/db/sqlc"
	"flex_server/tools"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

func DailyCreateEventHostCancel(ctx context.Context, server *Server) func() {

	return func() {
		result, err := RedisClient.SMembers(RedisContext, constants.CHARGE_TICKET_ID_CANCEL).Result()
		if err != nil || len(result) == 0 {
			if err != nil {
				log.Printf("There an error at DailyCreateEventHostCancel at RedisClient.SMembers: %v, type: %v \n", err.Error(), constants.CHARGE_TICKET_ID_CANCEL)
			}
			return
		}
		for _, uniqueID := range result {
			ProcessEventHostCancelUniqueID(ctx, server, uniqueID, "DailyCreateEventHostCancel")
		}
	}
}

func ProcessEventHostCancelUniqueID(ctx context.Context, server *Server, uniqueID string, funcName string) {
	data, err := RedisClient.HGetAll(RedisContext, uniqueID).Result()
	if err != nil {
		log.Printf("FuncName: %v. There an error at ProcessEventHostCancelUniqueID at RedisClient.HGetAll: %v, uniqueID: %v \n", funcName, err.Error(), uniqueID)
		return
	}

	result, err := RedisClient.SMembers(RedisContext, data[constants.REFERENCE]).Result()
	if err != nil {
		log.Printf("FuncName: %v. There an error at ProcessEventHostCancelUniqueID at RedisClient.SMembers: %v, data[constants.REFERENCE]: %v \n", funcName, err.Error(), data[constants.REFERENCE])
		return
	}
	var redisRemovedCount int
	for _, id := range result {
		chargeID, err := tools.StringToUuid(id)
		if err != nil {
			log.Printf("FuncName: %v. There an error at ProcessEventHostCancelUniqueID at tools.StringToUuid: %v, uniqueID: %v, chargeID: %v \n", funcName, err.Error(), uniqueID, id)
			continue
		}
		charge, err := server.store.GetChargeTicketReferenceByChargeID(ctx, db.GetChargeTicketReferenceByChargeIDParams{
			Cancelled:  false,
			IsComplete: true,
			ID:         chargeID,
		})
		if err != nil {
			log.Printf("FuncName: %v. There an error at ProcessEventHostCancelUniqueID at GetChargeTicketReferenceByChargeID: %v, uniqueID: %v, chargeID: %v \n", funcName, err.Error(), uniqueID, chargeID)
			continue
		}
		// We create the cancellation
		cancel, err := server.store.CreateCancellation(ctx, db.CreateCancellationParams{
			ChargeID:       charge.ChargeID,
			ChargeType:     charge.ChargeType,
			Type:           constants.HOST_CANCEL,
			CancelUserID:   charge.HostUserID,
			ReasonOne:      data[constants.REASON_ONE],
			ReasonTwo:      data[constants.REASON_TWO],
			Message:        data[constants.MESSAGE],
			MainOptionType: "events",
		})
		if err != nil {
			log.Printf("FuncName: %v. There an error at ProcessEventHostCancelUniqueID at .CreateCancellation: %v, uniqueID: %v, chargeID: %v \n", funcName, err.Error(), uniqueID, chargeID)
			continue
		}
		// We want to update
		// We update the charge ticket to cancel
		_, err = server.store.UpdateChargeTicketReferenceByID(ctx, db.UpdateChargeTicketReferenceByIDParams{
			Cancelled: pgtype.Bool{
				Bool:  true,
				Valid: true,
			},
			ID: chargeID,
		})
		if err != nil {
			log.Printf("FuncName: %v. There an error at ProcessEventHostCancelUniqueID at .CreateCancellation: %v, uniqueID: %v, chargeID: %v \n", funcName, err.Error(), uniqueID, chargeID)
			continue
		}

		// When a guest cancels we always want to create a refund to know how much we are giving the host and the guest
		_, err = server.store.CreateMainRefund(ctx, db.CreateMainRefundParams{
			ChargeID:    cancel,
			ChargeType:  constants.CHARGE_TICKET_REFERENCE,
			UserPercent: 100,
			HostPercent: 0,
			Type:        constants.HOST_CANCEL,
		})
		if err != nil {
			log.Printf("FuncName: %v. There an error at ProcessEventHostCancelUniqueID at .CreateMainRefund: %v, uniqueID: %v, chargeID: %v \n", funcName, err.Error(), uniqueID, chargeID)
			continue
		}
		// We want to send a message to the host saying a cancellation was made]
		header := fmt.Sprintf("Cancellation %v event", tools.HandleReadableDates(charge.StartDate, charge.EndDate, tools.DateDMMYyyy))
		CreateTypeNotification(ctx, server, charge.ChargeID, charge.GuestUserID, constants.HOST_CANCEL, data[constants.MESSAGE], false, header)
		err = RedisClient.SRem(RedisContext, data[constants.REFERENCE], id).Err()
		if err != nil {
			log.Printf("FuncName: %v. There an error at ProcessEventHostCancelUniqueID at .SRem: %v, uniqueID: %v, chargeID: %v \n", funcName, err.Error(), uniqueID, chargeID)
			continue
		} else {
			redisRemovedCount++
		}
	}
	if redisRemovedCount == len(result) {
		err = RedisClient.Del(RedisContext, uniqueID).Err()
		if err != nil {
			log.Printf("FuncName: %v. There an error at ProcessEventHostCancelUniqueID at RedisClient.Del: %v, uniqueID: %v \n", funcName, err.Error(), uniqueID)
			return
		}
		err = RedisClient.SRem(RedisContext, constants.CHARGE_TICKET_ID_CANCEL, uniqueID).Err()
		if err != nil {
			log.Printf("FuncName: %v. There an error at ProcessEventHostCancelUniqueID at RedisClient.Del: %v, uniqueID: %v \n", funcName, err.Error(), uniqueID)
			return
		}
	}
}

func (server *Server) CreateOptionUserCancel(ctx *gin.Context) {
	var req CreateUserOptionCancellationParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at CreateOptionUserCancel OptionDateParams in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.ChargeID)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	payoutIDs, err := RedisClient.SMembers(RedisContext, constants.PAYOUT_CHARGE_DATE_IDS).Result()
	if err != nil {
		log.Printf("Error at CreateOptionUserCancel Payouts RedisClient.SMembers: %v, ChargeID: %v \n", err.Error(), req.ChargeID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	refundIDs, err := RedisClient.SMembers(RedisContext, constants.REFUND_CHARGE_DATE_IDS).Result()
	if err != nil {
		log.Printf("Error at CreateOptionUserCancel Refunds RedisClient.SMembers: %v, ChargeID: %v \n", err.Error(), req.ChargeID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := HandleGetUser(ctx, server)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	chargeID, err := tools.StringToUuid(req.ChargeID)
	if err != nil {
		log.Printf("Error at CreateOptionUserCancel at tools.StringToUuid: %v, ChargeID: %v \n", err.Error(), req.ChargeID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	msg, err := HandleOptionUserCancel(ctx, server, user, req, "CreateOptionUserCancel", chargeID, payoutIDs, refundIDs)
	if err != nil {
		log.Printf("Error at CreateOptionUserCancel at HandleOptionUserCancel: %v, ChargeID: %v \n", err.Error(), req.ChargeID)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	res := CreateCancellationRes{
		ChargeID: req.ChargeID,
		Msg:      msg,
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) CreateEventUserCancel(ctx *gin.Context) {
	var req CreateUserEventCancellationParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at CreateEventUserCancel OptionDateParams in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.ChargeID)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	payoutIDs, err := RedisClient.SMembers(RedisContext, constants.PAYOUT_CHARGE_DATE_IDS).Result()
	if err != nil {
		log.Printf("Error at CreateEventUserCancel Payouts RedisClient.SMembers: %v, ChargeID: %v \n", err.Error(), req.ChargeID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	refundIDs, err := RedisClient.SMembers(RedisContext, constants.REFUND_CHARGE_DATE_IDS).Result()
	if err != nil {
		log.Printf("Error at CreateEventUserCancel Refunds RedisClient.SMembers: %v, ChargeID: %v \n", err.Error(), req.ChargeID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := HandleGetUser(ctx, server)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	chargeID, err := tools.StringToUuid(req.ChargeID)
	if err != nil {
		log.Printf("Error at CreateEventUserCancel at tools.StringToUuid: %v, ChargeID: %v \n", err.Error(), req.ChargeID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	msg, err := HandleEventUserCancel(ctx, server, user, req, "CreateEventUserCancel", chargeID, payoutIDs, refundIDs)
	if err != nil {
		log.Printf("Error at CreateEventUserCancel at HandleOptionUserCancel: %v, ChargeID: %v \n", err.Error(), req.ChargeID)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	res := CreateCancellationRes{
		ChargeID: req.ChargeID,
		Msg:      msg,
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) CreateOptionHostCancel(ctx *gin.Context) {
	var req CreateHostOptionCancellationParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at CreateOptionHostCancel OptionDateParams in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.ChargeID)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	optionID, err := tools.StringToUuid(req.OptionID)
	if err != nil {
		log.Printf("FuncName: %v. There an error at CreateOptionHostCancel at tools.StringToUuid(req.OptionID): %v, optionID: %v \n", "CreateOptionHostCancel", err.Error(), req.OptionID)
		err = fmt.Errorf("this event date cannot be found")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, _, option, isCoHost, userCoHost, err := HandleGetCompleteOptionReservation(optionID, ctx, server, true)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payoutIDs, err := RedisClient.SMembers(RedisContext, constants.PAYOUT_CHARGE_DATE_IDS).Result()
	if err != nil {
		log.Printf("Error at CreateOptionHostCancel Payouts RedisClient.SMembers: %v, ChargeID: %v \n", err.Error(), req.ChargeID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	refundIDs, err := RedisClient.SMembers(RedisContext, constants.REFUND_CHARGE_DATE_IDS).Result()
	if err != nil {
		log.Printf("Error at CreateOptionHostCancel Refunds RedisClient.SMembers: %v, ChargeID: %v \n", err.Error(), req.ChargeID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	chargeID, err := tools.StringToUuid(req.ChargeID)
	if err != nil {
		log.Printf("Error at CreateOptionHostCancel at tools.StringToUuid: %v, ChargeID: %v \n", err.Error(), req.ChargeID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	msg, err := HandleOptionHostCancel(ctx, server, user, req, "CreateOptionHostCancel", chargeID, payoutIDs, refundIDs)
	if err != nil {
		log.Printf("Error at CreateOptionHostCancel at HandleOptionUserCancel: %v, ChargeID: %v \n", err.Error(), req.ChargeID)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if isCoHost {
		HandleCoHostUpdateMsg(ctx, server, userCoHost, user, option, "CreateOptionHostCancel", "listing option booking cancellation", "listing option booking cancellation")
	}
	res := CreateCancellationRes{
		ChargeID: req.ChargeID,
		Msg:      msg,
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) CreateEventHostCancel(ctx *gin.Context) {
	var req CreateHostEventCancellationParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at CreateEventHostCancel OptionDateParams in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.EventID)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	eventID, err := tools.StringToUuid(req.EventID)
	if err != nil {
		log.Printf("FuncName: %v. There an error at HandleEventHostCancel at tools.StringToUuid(req.EventDateID): %v, startDate: %v, eventID: %v \n", "UpdateEventDatesBooked", err.Error(), req.StartDate, req.EventID)
		err = fmt.Errorf("this event date cannot be found")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	eventDateID, err := tools.StringToUuid(req.EventDateID)
	if err != nil {
		log.Printf("FuncName: %v. There an error at HandleEventHostCancel at tools.StringToUuid(req.EventDateID): %v, startDate: %v, eventDateID: %v \n", "UpdateEventDatesBooked", err.Error(), req.StartDate, req.EventDateID)
		err = fmt.Errorf("this event date cannot be found")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, _, option, isCoHost, userCoHost, err := HandleGetCompleteOptionReservation(eventID, ctx, server, true)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	msg, err := HandleEventHostCancel(ctx, server, user, req, eventDateID, eventID, "CreateEventHostCancel")
	if err != nil {
		log.Printf("Error at CreateEventHostCancel at HandleOptionUserCancel: %v, EventDateID: %v \n", err.Error(), req.EventDateID)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if isCoHost {
		HandleCoHostUpdateMsg(ctx, server, userCoHost, user, option, "ListEventDateBooked", "event date status", "update event date status")
	}
	res := CreateHostEventCancellationRes{
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		EventID:     req.EventID,
		EventDateID: req.EventDateID,
		Msg:         msg,
	}
	ctx.JSON(http.StatusOK, res)
}
