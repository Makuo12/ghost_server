package api

import (
	"context"
	"flex_server/constants"
	db "flex_server/db/sqlc"
	"flex_server/tools"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func HandleOptionUserCancel(ctx context.Context, server *Server, user db.User, req CreateUserOptionCancellationParams, funcName string, chargeID uuid.UUID, payoutRedisIDs []string, refundRedisIDs []string) (msg string, err error) {
	// We want to make sure charge ID is not in payoutsRedisIDs, refundRedisIDs
	if req.ReasonOne == "other" && tools.ServerStringEmpty(req.ReasonTwo) {
		err = fmt.Errorf("please give a detailed reason even though you selected other")
		return
	} else {
		if req.ReasonOne != "other" {
			req.ReasonTwo = "none"
		}
	}
	charge, refund, hostPayout, err := GetOptionUserCancelChargeRefund(ctx, server, user, funcName, chargeID, payoutRedisIDs, refundRedisIDs)
	if err != nil {
		log.Printf("FuncName: %v. There an error at HandleOptionUserCancel at GetOptionUserCancelChargeRefund: %v, chargeID: %v, userID: %v \n", funcName, err.Error(), chargeID, user.ID)
		err = fmt.Errorf("could not cancel your booking as an error occurred")
		return
	}
	// We create the cancellation
	cancel, err := server.store.CreateCancellation(ctx, db.CreateCancellationParams{
		ChargeID:       charge.ChargeID,
		ChargeType:     charge.ChargeType,
		Type:           constants.USER_CANCEL,
		CancelUserID:   user.UserID,
		ReasonOne:      req.ReasonOne,
		ReasonTwo:      req.ReasonTwo,
		Message:        req.Message,
		MainOptionType: "options",
	})
	if err != nil {
		log.Printf("FuncName: %v. There an error at HandleOptionUserCancel at .server.store.CreateCancellation: %v, chargeID: %v, userID: %v \n", funcName, err.Error(), chargeID, user.ID)
		err = fmt.Errorf("could not cancel your booking as an error occurred")
		return
	}
	// We want to update
	_, err = server.store.UpdateChargeOptionReferenceByID(ctx, db.UpdateChargeOptionReferenceByIDParams{
		Cancelled: pgtype.Bool{
			Bool:  true,
			Valid: true,
		},
		ID: charge.ChargeID,
	})
	if err != nil {
		log.Printf("FuncName: %v. There an error at HandleOptionUserCancel at .UpdateChargeOptionReferenceByID: %v, chargeID: %v, userID: %v \n", funcName, err.Error(), chargeID, user.ID)
		err = fmt.Errorf("could not cancel your booking as an error occurred")
		return
	}

	// When a guest cancels we always want to create a refund to know how much we are giving the host and the guest
	fund, err := server.store.CreateMainRefund(ctx, db.CreateMainRefundParams{
		ChargeID:    cancel,
		ChargeType:  constants.CHARGE_OPTION_REFERENCE,
		UserPercent: int32(refund),
		HostPercent: int32(hostPayout),
		Type:        constants.USER_CANCEL,
	})
	if err != nil {
		log.Printf("FuncName: %v. There an error at HandleOptionUserCancel at server.store.CreateMainRefund: %v, chargeID: %v, userID: %v refund: %v \n", funcName, err.Error(), chargeID, user.ID, refund)
		err = fmt.Errorf("your cancellation was successful, but refund was not created")
		return
	} else {
		if charge.MainPayoutComplete {
			// We send the user a notification saying that host has already been paid so no refund available
			msg = fmt.Sprintf("Hey %v, the host has already been paid so no refund currently available. Please contact us so we can know if we can provide you with a refund", charge.UserFirstName)
		} else {
			msg = fmt.Sprintf("Hey %v, you refund of %v percent was created and would be sent to you shortly. This may take up to 5 to 15 days", charge.UserFirstName, fund)
		}
	}
	// We want to send a message to the host saying a cancellation was made]
	header := fmt.Sprintf("%v's, cancelled for %v", tools.CapitalizeFirstCharacter(charge.UserFirstName), charge.HostNameOption)
	CreateTypeNotification(ctx, server, charge.ChargeID, charge.HostUserID, constants.USER_CANCEL, req.Message, false, header)
	return
}

func HandleOptionHostCancel(ctx context.Context, server *Server, user db.User, req CreateHostOptionCancellationParams, funcName string, chargeID uuid.UUID, payoutRedisIDs []string, refundRedisIDs []string) (msg string, err error) {
	if req.ReasonOne == "other" && tools.ServerStringEmpty(req.ReasonTwo) {
		err = fmt.Errorf("please give a detailed reason even though you selected other")
		return
	} else {
		if req.ReasonOne != "other" {
			req.ReasonTwo = "none"
		}
	}

	userID, err := tools.StringToUuid(req.UserID)
	if err != nil {
		log.Printf("FuncName: %v. There an error at HandleOptionHostCancel at tools.StringToUuid: %v, chargeID: %v, userID: %v \n", funcName, err.Error(), chargeID, user.ID)
		err = fmt.Errorf("could not find the guest for the booking. Please try again later or try contacting us")
		return
	}
	charge, refund, hostPayout, err := GetOptionHostCancelChargeRefund(ctx, server, user, "HandleOptionHostCancel", chargeID, payoutRedisIDs, refundRedisIDs, userID)
	if err != nil {
		return
	}
	// We create the cancellation
	cancel, err := server.store.CreateCancellation(ctx, db.CreateCancellationParams{
		ChargeID:       charge.ChargeID,
		ChargeType:     charge.ChargeType,
		Type:           constants.HOST_CANCEL,
		CancelUserID:   user.UserID,
		ReasonOne:      req.ReasonOne,
		ReasonTwo:      req.ReasonTwo,
		Message:        req.Message,
		MainOptionType: "options",
	})
	if err != nil {
		log.Printf("FuncName: %v. There an error at HandleOptionHostCancel at .server.store.CreateCancellation: %v, chargeID: %v, userID: %v \n", funcName, err.Error(), chargeID, user.ID)
		err = fmt.Errorf("could not cancel this booking an error occurred")
		return
	}
	// We want to update
	_, err = server.store.UpdateChargeOptionReferenceByID(ctx, db.UpdateChargeOptionReferenceByIDParams{
		Cancelled: pgtype.Bool{
			Bool:  true,
			Valid: true,
		},
		ID: charge.ChargeID,
	})
	if err != nil {
		log.Printf("FuncName: %v. There an error at HandleOptionHostCancel at .UpdateChargeOptionReferenceByID: %v, chargeID: %v, userID: %v \n", funcName, err.Error(), chargeID, user.ID)
		err = fmt.Errorf("could not cancel this booking an error occurred")
		return
	}
	// When a guest cancels we always want to create a refund to know how much we are giving the host and the guest
	_, err = server.store.CreateMainRefund(ctx, db.CreateMainRefundParams{
		ChargeID:    cancel,
		ChargeType:  constants.CHARGE_OPTION_REFERENCE,
		UserPercent: int32(refund),
		HostPercent: int32(hostPayout),
		Type:        constants.HOST_CANCEL,
	})
	if err != nil {
		log.Printf("FuncName: %v. There an error at HandleOptionHostCancel at store.CreateMainRefund: %v, chargeID: %v, userID: %v refund: %v \n", funcName, err.Error(), chargeID, user.ID, refund)
		return
	}
	// We want to send a message to the host saying a cancellation was made]
	header := fmt.Sprintf("%v's, cancelled for %v", tools.CapitalizeFirstCharacter(charge.HostFirstName), charge.HostNameOption)
	CreateTypeNotification(ctx, server, charge.ChargeID, userID, constants.HOST_CANCEL, req.Message, false, header)
	return
}

func HandleEventUserCancel(ctx context.Context, server *Server, user db.User, req CreateUserEventCancellationParams, funcName string, chargeID uuid.UUID, payoutRedisIDs []string, refundRedisIDs []string) (msg string, err error) {
	if req.ReasonOne == "other" && tools.ServerStringEmpty(req.ReasonTwo) {
		err = fmt.Errorf("please give a detailed reason even though you selected other")
		return
	} else {
		if req.ReasonOne != "other" {
			req.ReasonTwo = "none"
		}
	}
	charge, refund, hostPayout, err := GetEventUserCancelChargeRefund(ctx, server, user, "HandleEventUserCancel", chargeID, payoutRedisIDs, refundRedisIDs)
	if err != nil {
		return
	}

	// We create the cancellation
	cancel, err := server.store.CreateCancellation(ctx, db.CreateCancellationParams{
		ChargeID:       charge.ChargeID,
		ChargeType:     charge.ChargeType,
		Type:           constants.USER_CANCEL,
		CancelUserID:   user.UserID,
		ReasonOne:      req.ReasonOne,
		ReasonTwo:      req.ReasonTwo,
		Message:        req.Message,
		MainOptionType: "events",
	})
	if err != nil {
		log.Printf("FuncName: %v. There an error at HandleEventUserCancel at .CreateCancellation: %v, chargeID: %v, userID: %v \n", funcName, err.Error(), chargeID, user.ID)
		err = fmt.Errorf("could not cancel your booking as an error occurred")
		return
	}
	// We want to update
	_, err = server.store.UpdateChargeTicketReferenceByID(ctx, db.UpdateChargeTicketReferenceByIDParams{
		Cancelled: pgtype.Bool{
			Bool:  true,
			Valid: true,
		},
		ID: charge.ChargeID,
	})
	if err != nil {
		log.Printf("FuncName: %v. There an error at HandleEventUserCancel at .UpdateChargeTicketReferenceByID: %v, chargeID: %v, userID: %v \n", funcName, err.Error(), chargeID, user.ID)
		err = fmt.Errorf("could not cancel your booking as an error occurred")
		return
	}

	// When a guest cancels we always want to create a refund to know how much we are giving the host and the guest
	fund, err := server.store.CreateMainRefund(ctx, db.CreateMainRefundParams{
		ChargeID:    cancel,
		ChargeType:  constants.CHARGE_TICKET_REFERENCE,
		UserPercent: int32(refund),
		HostPercent: int32(hostPayout),
		Type:        constants.USER_CANCEL,
	})
	if err != nil {
		log.Printf("FuncName: %v. There an error at HandleEventUserCancel at .CreateMainRefund: %v, chargeID: %v, userID: %v refund: %v \n", funcName, err.Error(), chargeID, user.ID, refund)
		err = fmt.Errorf("your cancellation was successful, but refund was not created")
		return
	} else {
		if charge.MainPayoutComplete {
			// We send the user a notification saying that host has already been paid so no refund available
			msg = fmt.Sprintf("Hey %v, the host has already been paid so no refund currently available. Please contact us so we can know if we can provide you with a refund", charge.UserFirstName)
		} else {
			msg = fmt.Sprintf("Hey %v, you refund of %v percent was created and would be sent to you shortly. This may take one to three days", charge.UserFirstName, fund)
		}
	}
	// We want to send a message to the host saying a cancellation was made]
	header := fmt.Sprintf("%v's, cancelled for %v", charge.UserFirstName, tools.HandleReadableDates(charge.StartDate, charge.EndDate, tools.DateDMM))
	CreateTypeNotification(ctx, server, charge.ChargeID, charge.HostUserID, constants.USER_CANCEL, req.Message, false, header)
	return
}

func HandleEventHostCancel(ctx context.Context, server *Server, user db.User, req CreateHostEventCancellationParams, eventDateID uuid.UUID, eventID uuid.UUID, funcName string) (msg string, err error) {
	if req.ReasonOne == "other" && tools.ServerStringEmpty(req.ReasonTwo) {
		err = fmt.Errorf("please give a detailed reason even though you selected other")
		return
	} else {
		if req.ReasonOne != "other" {
			req.ReasonTwo = "none"
		}
	}

	chargeIDs, eventDate, err := GetEventHostChargeCancel(ctx, server, user, "HandleEventUserCancel", req.StartDate, req.EndDate, eventDateID, eventID)
	if err != nil {
		return
	}
	// Lets
	// When cancelling a single event we set that event to in_active
	// When cancelling a recurring event we remove that particular date
	switch eventDate.Type {
	case "single":
		_, err = server.store.UpdateEventDateTimeActive(ctx, db.UpdateEventDateTimeActiveParams{
			IsActive: pgtype.Bool{
				Bool:  false,
				Valid: true,
			},
			ID: eventDate.EventDateTimeID,
		})
		if err != nil {
			log.Printf("FuncName: %v. There an error at HandleEventUserCancel at UpdateEventDateTime: %v, startDate: %v, userID: %v \n", funcName, err.Error(), req.StartDate, user.ID)
			err = fmt.Errorf("unable to cancel this single event")
			return
		}
	case "recurring":
		newDates, errRecur := tools.RemoveRecurDate(req.StartDate, eventDate.EventDates)
		if errRecur != nil {
			err = errRecur
			return
		}
		_, err = server.store.UpdateEventDateTimeDates(ctx, db.UpdateEventDateTimeDatesParams{
			EventDates: newDates,
			ID:         eventDate.EventDateTimeID,
		})
		if err != nil {
			log.Printf("FuncName: %v. There an error at HandleEventUserCancel at UpdateEventDateTimeDates: %v, startDate: %v, userID: %v \n", funcName, err.Error(), req.StartDate, user.ID)
			err = fmt.Errorf("unable to cancel this recurring event")
			return
		}
	}
	if err != nil {
		return
	}
	if len(chargeIDs) == 0 {
		log.Println("chargeID empty ", chargeIDs)
		return
	}
	chargeReference := tools.UuidToString(uuid.New())
	data := []string{
		constants.START_DATE,
		req.StartDate,
		constants.END_DATE,
		req.StartDate,
		constants.MESSAGE,
		req.Message,
		constants.REASON_ONE,
		req.ReasonOne,
		constants.REASON_TWO,
		req.ReasonTwo,
		constants.REFERENCE,
		chargeReference,
	}
	uniqueID := fmt.Sprintf("%v&%v&%v&%v", tools.UuidToString(eventDate.EventDateTimeID), eventDate.Type, req.StartDate, "cancel")
	err = RedisClient.SAdd(RedisContext, constants.CHARGE_TICKET_ID_CANCEL, uniqueID).Err()
	if err != nil {
		log.Printf("FuncName: %v. There an error at HandleEventUserCancel at RedisClient.SAdd: %v, startDate: %v, userID: %v, chargeIDs: %v \n", funcName, err.Error(), req.StartDate, user.ID, chargeIDs)
		err = nil
	}
	err = RedisClient.HSet(RedisContext, uniqueID, data).Err()
	if err != nil {
		log.Printf("FuncName: %v. There an error at HandleEventUserCancel at RedisClient.HSet: %v, startDate: %v, userID: %v, chargeIDs: %v \n", funcName, err.Error(), req.StartDate, user.ID, chargeIDs)
		err = nil
	}
	// We want to store the chargeIDs in redis so that we can later give full refund and set it to cancel
	err = RedisClient.SAdd(RedisContext, chargeReference, tools.ListUuidToString(chargeIDs)).Err()
	if err != nil {
		log.Printf("FuncName: %v. There an error at HandleEventUserCancel at RedisClient.SAdd: %v, startDate: %v, userID: %v, chargeIDs: %v \n", funcName, err.Error(), req.StartDate, user.ID, chargeIDs)
		err = nil
	}

	return
}
