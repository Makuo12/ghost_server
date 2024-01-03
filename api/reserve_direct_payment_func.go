package api

import (
	"context"
	"flex_server/constants"
	db "flex_server/db/sqlc"
	"flex_server/tools"
	"fmt"
	"log"
)

func HandleInitPaymentEvent(user db.User, reference string) (res InitPaymentRes, err error) {
	reserveData, err := HandleEventReserveRedisData(tools.UuidToString(user.ID), reference)
	if err != nil {
		log.Printf("Error at HandleInitPaymentEvent in HandleOptionReserveRedisData: %v, reference: %v, userID: %v \n", err.Error(), reference, user.ID)
		err = fmt.Errorf("reference has expired, please try again")
		return
	}
	res = InitPaymentRes{
		Reference: reference,
		Charge:    tools.ConvertToPaystackCharge(reserveData.TotalFee),
		Reason:    constants.USER_EVENT_PAYMENT,
		Email:     user.Email,
	}
	return
}

func HandleInitRefund(ctx context.Context, server *Server, user db.User, referenceString string, reason string, amount int, currency string, funcName string) {
	reference, err := tools.StringToUuid(referenceString)
	if err != nil {
		log.Printf("Error at FuncName: %v HandleInitRefund in tools.StringToUuid: %v, reference: %v, userID: %v \n", funcName, err.Error(), reference, user.ID)
		return
	}
	charge, err := server.store.CreateChargeReference(ctx, db.CreateChargeReferenceParams{
		UserID:     user.UserID,
		Reason:     reason,
		Charge:     int64(amount),
		Currency:   currency,
		IsComplete: true,
		Reference:  reference,
	})
	if err != nil {
		log.Printf("Error at FuncName: %v HandleInitRefund in CreateChargeReference: %v, reference: %v, userID: %v \n", funcName, err.Error(), reference, user.ID)
		return
	}
	_, err = server.store.CreateMainRefund(ctx, db.CreateMainRefundParams{
		ChargeID: charge.ID,
		UserPercent: 100,
		HostPercent: 0,
		ChargeType: constants.CHARGE_REFERENCE,
		Type: constants.USER_PAYMENT_INVALID,
	})
	if err != nil {
		log.Printf("Error at FuncName: %v HandleInitRefund in CreateMainRefund: %v, reference: %v, userID: %v \n", funcName, err.Error(), reference, user.ID)
		
	}
	return
}
