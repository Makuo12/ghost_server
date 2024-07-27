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

func CreateChargeReference(ctx context.Context, server *Server, userID uuid.UUID, reference string, paymentReference string, objectReference uuid.UUID, hasObjectReference bool, reason string, currency string, mainOptionType string, fee string, funcName string) (charge db.ChargeReference, err error) {
	// Delete the any charge reference that uses the reference
	err = server.store.RemoveChargeReferenceComplete(ctx, reference)
	if err != nil {
		log.Printf("Error at %v at CreateChargeReference, RemoveChargeReferenceComplete err %v\n", funcName, err.Error())
		err = nil
	}
	// userID of type user.UserID
	charge, err = server.store.CreateChargeReference(ctx, db.CreateChargeReferenceParams{
		UserID:             userID,
		Reference:          reference,
		ObjectReference:    objectReference,
		HasObjectReference: hasObjectReference,
		PaymentMedium:      constants.PAYSTACK,
		PaymentChannel:     constants.PAYSTACK_CARD,
		Reason:             reason,
		PaymentReference:   paymentReference,
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

func ObjectOptionPaymentReference(ctx context.Context, server *Server, user db.User, reference string, paymentReference string, objectReference uuid.UUID, amount int, currency string, message string) (success bool, err error) {
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
			HandleInitRefund(ctx, server, user, reference, paymentReference, objectReference, false, "options", constants.USER_OPTION_INVALID_PAYMENT_AMOUNT, tools.IntToMoneyString(int64(amount)), charge.Currency, "ObjectOptionPaymentReference")
			err = fmt.Errorf("amount payed is not the total amount for the listing")
			return
		}
	} else {
		fromCharge = true
		referenceCharge = charge.Reference
		if amount != int(charge.TotalFee) {
			// If amount is not equal we send back a refund and an error
			HandleInitRefund(ctx, server, user, reference, paymentReference, objectReference, true, "options", constants.USER_OPTION_INVALID_PAYMENT_AMOUNT, tools.IntToMoneyString(int64(amount)), charge.Currency, "ObjectOptionPaymentReference")
			err = fmt.Errorf("amount payed is not the total amount for the listing")
			return
		}
	}
	if err == nil {
		host, err := server.store.GetOptionInfoUserIDByUserID(ctx, charge.OptionUserID)
		if err != nil {
			log.Printf("Error at HandleOptionReserveRequest in GetOptionInfoByUserID: %v optionID: %v referenceID: %v, pay_method_reference: %v\n", err.Error(), charge.OptionUserID, charge.Reference, charge.PaymentReference)
			err = nil
		} else {
			header := "Payment and Reservation confirmed"
			msg := "Thank you for using Flizzup"
			checkIn := tools.HandleReadableDate(charge.StartDate, tools.DateDMMYyyy)
			checkout := tools.HandleReadableDate(charge.EndDate, tools.DateDMMYyyy)
			BrevoOptionPaymentSuccess(ctx, server, header, msg, "FinalOptionReserveDetail", charge.ID, host.Email, host.FirstName, host.LastName, tools.UuidToString(charge.ID), tools.UuidToString(host.UserID), user.Email, user.FirstName, user.LastName, tools.UuidToString(user.UserID), host.HostNameOption, checkIn, checkout)
			// Notification for Guest
			msg = "Payment received, and reservation confirmed! Check your email for further information about your booking, including scheduling an inspection and our 100% refund policy if the property does not match the app's description."
			header = fmt.Sprintf("Hey %v, booking confirmed", user.FirstName)
			CreateTypeNotification(ctx, server, charge.ID, user.UserID, constants.OPTION_PAYMENT_SUCCESSFUL, msg, false, header)
			// Notification for Host
			msg = fmt.Sprintf("You have a new booking at %v! Check your email for further details.", host.HostNameOption)
			header = fmt.Sprintf("Hey %v", host.FirstName)
			CreateTypeNotification(ctx, server, charge.ID, host.UserID, constants.HOST_OPTION_PAYMENT_SUCCESSFUL, msg, false, header)
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





func ObjectEventPaymentReference(ctx context.Context, server *Server, user db.User, reference string, paymentReference string, objectReference uuid.UUID, amount int, currency string, message string) (success bool, err error) {
	// A reservation request can never be created for a ticket so we don't need to check for a chargeEventReference
	reserveEventData, err := HandleEventReserveRedisData(tools.UuidToString(user.UserID), reference)
	if err != nil {
		log.Printf("Error at HandleInitPaymentEvent in HandleOptionReserveRedisData: %v, reference: %v, user.UserID: %v \n", err.Error(), reference, user.UserID)
		err = fmt.Errorf("reference has expired, please try again")
		return
	}
	if amount != tools.ConvertToPaystackCharge(reserveEventData.TotalFee) {
		// If amount is not equal we send back a refund and an error
		HandleInitRefund(ctx, server, user, reference, paymentReference, objectReference, false, "events", constants.USER_OPTION_INVALID_PAYMENT_AMOUNT, tools.IntToMoneyString(int64(amount)), reserveEventData.Currency, "ObjectEventPaymentReference")
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


func CheckPaymentProgress(ctx context.Context, server *Server, chargeID uuid.UUID) error {
	chargeData, err := server.store.GetChargeOptionReferencePayment(ctx, chargeID)
	if err != nil {
		return err
	}
	if HandleSqlNullString(chargeData.MainPayoutStatus) == "not_started" || HandleSqlNullString(chargeData.MainRefundStatus) == "not_started" || HandleSqlNullString(chargeData.RefundPayoutStatus) == "not_started" {
		return nil
	} else {
		return fmt.Errorf("This current booking cannot be changed because it is current in use for either refund or payout")
	}
}

