package api

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/makuo12/ghost_server/constants"
	db "github.com/makuo12/ghost_server/db/sqlc"
	"github.com/makuo12/ghost_server/tools"
)

func CreateChargeReference(ctx context.Context, server *Server, userID uuid.UUID, paymentReference string, objectReference uuid.UUID, hasObjectReference bool, reason string, currency string, mainOptionType string, fee string, funcName string) (charge db.ChargeReference, err error) {
	// userID of type user.UserID
	charge, err = server.store.CreateChargeReference(ctx, db.CreateChargeReferenceParams{
		UserID:             userID,
		Reference:          paymentReference,
		ObjectReference:    objectReference,
		HasObjectReference: hasObjectReference,
		PaymentMedium:      constants.PAYSTACK,
		PaymentChannel:     constants.PAYSTACK_CARD,
		Reason:             reason,
		MainObjectType:     mainOptionType,
		Charge:             tools.MoneyStringToInt(fee),
		Currency:           currency,
		IsComplete:         false,
	})
	if err != nil {
		log.Printf("Error at %v at CreateChargeReference, err %v\n", funcName, err.Error())
	}
	return
}

func ObjectOptionPaymentReference(ctx context.Context, server *Server, user db.User, reference string, objectReference uuid.UUID, amount int, currency string, message string) (success bool, err error) {
	var fromCharge bool
	var reserveOptionData ExperienceReserveOModel
	var referenceCharge string
	// If there is no ChargeOptionReference instead we get the charge from redis
	// This usually occurs when a user makes a reservation request so there is already a charge created.
	charge, err := server.store.GetChargeOptionReferenceByID(ctx, db.GetChargeOptionReferenceByIDParams{
		ChargeID: objectReference,
		UserID:   user.UserID,
	})
	if err != nil {
		log.Printf("Error at ObjectOptionPaymentReference in GetChargeOptionReferenceByID: %v, reference: %v, userID: %v \n", err.Error(), objectReference, user.UserID)
		err = nil
		reserveOptionData, err = HandleOptionReserveRedisData(tools.UuidToString(user.UserID), reference)
		if err != nil {
			log.Printf("Error at ObjectOptionPaymentReference in HandleOptionReserveRedisData: %v, reference: %v, userID: %v \n", err.Error(), reference, user.UserID)
			err = fmt.Errorf("req.Reference has expired, please try again")
			return
		}
		if amount != tools.ConvertToPaystackCharge(reserveOptionData.TotalFee) {
			// If amount is not equal we send back a refund and an error
			HandleInitRefund(ctx, server, user, reference, objectReference, false, "options", constants.USER_OPTION_INVALID_PAYMENT_AMOUNT, tools.IntToMoneyString(int64(amount)), currency, "ObjectOptionPaymentReference")
			err = fmt.Errorf("amount payed is not the total amount for the listing")
			return
		}
	} else {
		fromCharge = true
		referenceCharge = charge.Reference
		if amount != int(charge.TotalFee) {
			// If amount is not equal we send back a refund and an error
			HandleInitRefund(ctx, server, user, reference, objectReference, true, "options", constants.USER_OPTION_INVALID_PAYMENT_AMOUNT, tools.IntToMoneyString(int64(amount)), currency, "ObjectOptionPaymentReference")
			err = fmt.Errorf("amount payed is not the total amount for the listing")
			return
		}
	}

	// We Store the receipt and snap shot
	err = HandleOptionReserveComplete(server, ctx, reserveOptionData, referenceCharge, reference, user, message, fromCharge)
	if err != nil {
		return
	}
	success = true
	return

}

func ObjectEventPaymentReference(ctx context.Context, server *Server, user db.User, reference string, objectReference uuid.UUID, amount int, currency string, message string) (success bool, err error) {
	// A reservation request can never be created for a ticket so we don't need to check for a chargeEventReference
	reserveEventData, err := HandleEventReserveRedisData(tools.UuidToString(user.UserID), reference)
	if err != nil {
		log.Printf("Error at HandleInitPaymentEvent in HandleOptionReserveRedisData: %v, reference: %v, user.UserID: %v \n", err.Error(), reference, user.UserID)
		err = fmt.Errorf("reference has expired, please try again")
		return
	}
	if amount != tools.ConvertToPaystackCharge(reserveEventData.TotalFee) {
		// If amount is not equal we send back a refund and an error
		HandleInitRefund(ctx, server, user, reference, objectReference, false, "events", constants.USER_OPTION_INVALID_PAYMENT_AMOUNT, tools.IntToMoneyString(int64(amount)), currency, "ObjectEventPaymentReference")
		err = fmt.Errorf("amount payed is not the total amount for the event")
		return
	}
	// Creating snapshot
	err = HandleEventReserveComplete(server, ctx, reserveEventData, reference, user, message, reference)
	if err != nil {
		return
	}
	success = true
	return
}
