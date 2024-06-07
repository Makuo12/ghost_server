package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/makuo12/ghost_server/constants"
	"github.com/makuo12/ghost_server/tools"

	"github.com/gin-gonic/gin"
)

func (server *Server) GetOptionUserCancelDetail(ctx *gin.Context) {
	var req CancelUserDetailParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at GetOptionUserCancelDetail OptionDateParams in ShouldBindJSON: %v, chargeID: %v \n", err.Error(), req.ChargeID)
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
		log.Printf("Error at GetOptionUserCancelDetail tools.StringToUuid: %v, ChargeID: %v \n", err.Error(), req.ChargeID)
		err = fmt.Errorf("this listing reservation cannot be found")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	payoutRedisIDs, err := RedisClient.SMembers(RedisContext, constants.PAYOUT_CHARGE_DATE_IDS).Result()
	if err != nil {
		log.Printf("Error at GetOptionUserCancelDetail Payouts RedisClient.SMembers: %v, ChargeID: %v \n", err.Error(), req.ChargeID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	refundRedisIDs, err := RedisClient.SMembers(RedisContext, constants.REFUND_CHARGE_DATE_IDS).Result()
	if err != nil {
		log.Printf("Error at GetOptionUserCancelDetail Refunds RedisClient.SMembers: %v, ChargeID: %v \n", err.Error(), req.ChargeID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	charge, refund, _, err := GetOptionUserCancelChargeRefund(ctx, server, user, "GetOptionUserCancelDetail", chargeID, payoutRedisIDs, refundRedisIDs)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var amount string
	if refund == 100 {
		amount = tools.IntToMoneyString(charge.TotalFee)
	} else if refund == 0 {
		amount = "0"
	} else {
		amount = tools.ConvertFloatToString(float64(refund/100) * tools.ConvertStringToFloat(tools.IntToMoneyString(charge.TotalFee)))
	}
	res := CancelUserOptionDetailRes{
		MainPrice:     tools.IntToMoneyString(charge.MainPrice),
		CleaningFee:   tools.IntToMoneyString(charge.CleanFee),
		ServiceFee:    tools.IntToMoneyString(charge.ServiceFee),
		TotalFee:      tools.IntToMoneyString(charge.TotalFee),
		GuestFee:      tools.IntToMoneyString(charge.GuestFee),
		PetFee:        tools.IntToMoneyString(charge.PetFee),
		RefundPercent: refund,
		Refund:        amount,
		RefundType:    charge.CancelPolicyOne,
		DateBooked:    tools.ConvertDateOnlyToString(charge.DateBooked),
		Currency:      charge.Currency,
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) GetEventUserCancelDetail(ctx *gin.Context) {
	var req CancelUserDetailParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at GetEventUserCancelDetail OptionDateParams in ShouldBindJSON: %v, chargeID: %v \n", err.Error(), req.ChargeID)
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
		log.Printf("Error at GetEventUserCancelDetail tools.StringToUuid: %v, ChargeID: %v \n", err.Error(), req.ChargeID)
		err = fmt.Errorf("this listing reservation cannot be found")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	payoutRedisIDs, err := RedisClient.SMembers(RedisContext, constants.PAYOUT_CHARGE_DATE_IDS).Result()
	if err != nil {
		log.Printf("Error at GetEventUserCancelDetail Payouts RedisClient.SMembers: %v, ChargeID: %v \n", err.Error(), req.ChargeID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	refundRedisIDs, err := RedisClient.SMembers(RedisContext, constants.REFUND_CHARGE_DATE_IDS).Result()
	if err != nil {
		log.Printf("Error at GetEventUserCancelDetail Refunds RedisClient.SMembers: %v, ChargeID: %v \n", err.Error(), req.ChargeID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	charge, refund, _, err := GetEventUserCancelChargeRefund(ctx, server, user, "GetEventUserCancelDetail", chargeID, payoutRedisIDs, refundRedisIDs)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var amount string
	if refund == 100 {
		amount = tools.IntToMoneyString(charge.TotalFee)
	} else if refund == 0 {
		amount = "0"
	} else {
		amount = tools.ConvertFloatToString(float64(refund/100) * tools.ConvertStringToFloat(tools.IntToMoneyString(charge.TotalFee)))
	}
	res := CancelUserTicketDetailRes{
		ServiceFee:    tools.IntToMoneyString(charge.ServiceFee),
		TicketPrice:   tools.IntToMoneyString(charge.TotalFee),
		RefundPercent: refund,
		Refund:        amount,
		Currency:      charge.Currency,
		RefundType:    charge.CancelPolicyOne,
		DateBooked:    tools.ConvertDateOnlyToString(charge.DateBooked),
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) GetOptionHostCancelDetail(ctx *gin.Context) {
	var req CancelHostOptionDetailParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at GetOptionHostCancelDetail OptionDateParams in ShouldBindJSON: %v, chargeID: %v \n", err.Error(), req.ChargeID)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	optionID, err := tools.StringToUuid(req.OptionID)
	if err != nil {
		log.Printf("FuncName: %v. There an error at CreateOptionHostCancel at tools.StringToUuid(req.OptionID): %v, optionID: %v \n", "CreateOptionHostCancel", err.Error(), req.OptionID)
		err = fmt.Errorf("this listing date cannot be found")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	userID, err := tools.StringToUuid(req.UserID)
	if err != nil {
		log.Printf("FuncName: %v. There an error at CreateOptionHostCancel at tools.StringToUuid(req.UserID): %v, userID: %v \n", "CreateOptionHostCancel", err.Error(), req.UserID)
		err = fmt.Errorf("this listing date cannot be found")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, _, _, _, _, err := HandleGetCompleteOptionReservation(optionID, ctx, server, true)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	chargeID, err := tools.StringToUuid(req.ChargeID)
	if err != nil {
		log.Printf("Error at GetOptionHostCancelDetail tools.StringToUuid: %v, ChargeID: %v \n", err.Error(), req.ChargeID)
		err = fmt.Errorf("this listing reservation cannot be found")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	payoutRedisIDs, err := RedisClient.SMembers(RedisContext, constants.PAYOUT_CHARGE_DATE_IDS).Result()
	if err != nil {
		log.Printf("Error at GetOptionHostCancelDetail Payouts RedisClient.SMembers: %v, ChargeID: %v \n", err.Error(), req.ChargeID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	refundRedisIDs, err := RedisClient.SMembers(RedisContext, constants.REFUND_CHARGE_DATE_IDS).Result()
	if err != nil {
		log.Printf("Error at GetOptionHostCancelDetail Refunds RedisClient.SMembers: %v, ChargeID: %v \n", err.Error(), req.ChargeID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	_, _, hostPayout, err := GetOptionHostCancelChargeRefund(ctx, server, user, "GetOptionHostCancelDetail", chargeID, payoutRedisIDs, refundRedisIDs, userID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	res := CancelHostDetailRes{
		CanCancel:  true,
		HostPayout: hostPayout,
		Amount:     "0.00",
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) GetEventHostCancelDetail(ctx *gin.Context) {
	var req CancelHostEventDetailParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at GetEventHostCancelDetail OptionDateParams in ShouldBindJSON: %v, eventDateTimeID: %v \n", err.Error(), req.EventDateID)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	optionID, err := tools.StringToUuid(req.EventID)
	if err != nil {
		log.Printf("FuncName: %v. There an error at CreateOptionHostCancel at tools.StringToUuid(req.EventID): %v, eventID: %v \n", "CreateOptionHostCancel", err.Error(), req.EventID)
		err = fmt.Errorf("this event date cannot be found")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	eventDateID, err := tools.StringToUuid(req.EventDateID)
	if err != nil {
		log.Printf("FuncName: %v. There an error at CreateOptionHostCancel at tools.StringToUuid(req.EventDateID): %v, eventDateID: %v \n", "CreateOptionHostCancel", err.Error(), req.EventDateID)
		err = fmt.Errorf("this event date cannot be found")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, _, _, _, _, err := HandleGetCompleteOptionReservation(optionID, ctx, server, true)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	_, _, err = GetEventHostChargeCancel(ctx, server, user, "HandleEventUserCancel", req.StartDate, req.EndDate, eventDateID, optionID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	res := CancelHostDetailRes{
		CanCancel:  true,
		HostPayout: 0,
		Amount:     "0.00",
	}
	ctx.JSON(http.StatusOK, res)
}
