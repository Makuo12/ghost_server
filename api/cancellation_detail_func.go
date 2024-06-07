package api

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/makuo12/ghost_server/constants"
	db "github.com/makuo12/ghost_server/db/sqlc"
	"github.com/makuo12/ghost_server/tools"
	"github.com/makuo12/ghost_server/val"

	"github.com/google/uuid"
)

func GetOptionUserCancelChargeRefund(ctx context.Context, server *Server, user db.User, funcName string, chargeID uuid.UUID, payoutRedisIDs []string, refundRedisIDs []string) (charge db.GetChargeOptionReferenceByUserIDRow, refund int, hostPayout int, err error) {
	charge, err = server.store.GetChargeOptionReferenceByUserID(ctx, db.GetChargeOptionReferenceByUserIDParams{
		UserID:     user.UserID,
		Cancelled:  false,
		IsComplete: true,
		ID:         chargeID,
	})
	if err != nil {
		log.Printf("FuncName: %v. There an error at GetOptionUserCancelRefund at GetChargeOptionReferenceByUserID: %v, chargeID: %v, userID: %v \n", funcName, err.Error(), chargeID, user.ID)
		err = fmt.Errorf("could not find your booking. Please try again later or try contacting us")
		return
	}
	//// Once a host is paid you cannot perform a cancel
	switch true {
	case time.Now().After(charge.EndDate.Add(time.Hour * -10)):
		err = fmt.Errorf("your cancellation cannot be made because your stay is already coming to an end")
		return
	}
	// Lets check if the guest has refund available
	startDate := tools.ConvertToTimeZoneTwo(charge.StartDate, charge.TimeZone)
	cancelPolicy, exist := val.GetCancelPolicy(charge.CancelPolicyOne)
	if ChargeIDInPay(refundRedisIDs, chargeID) {
		err = fmt.Errorf("your cancellation cannot be made because the booking's refund is currently processing")
		return
	}
	if charge.MainPayoutComplete || ChargeIDInPay(payoutRedisIDs, chargeID) {
		refund = 0
		hostPayout = 0
	} else {
		if exist {
			for _, p := range cancelPolicy.Items {
				if p.Type == "hard" {

					timeWithZone := tools.ConvertToTimeZoneTwo(time.Now(), charge.TimeZone)
					startTime := startDate.Add(time.Hour).Add(time.Hour * -time.Duration(p.HoursTwo))
					timeOfBooking := timeWithZone.Add(time.Hour).Before(charge.DateBooked.Add(time.Hour * time.Duration(p.Hours)))
					if time.Now().Before(startTime) && timeOfBooking {
						// We need to check the dateBooked against Hours
						{
							refund = p.Percent
							hostPayout = 100 - refund
							break
						}
					}
				} else {

					startTime := startDate.Add(time.Hour).Add(time.Hour * -time.Duration(p.Hours))
					if time.Now().Before(startTime) {
						refund = p.Percent
						hostPayout = 100 - refund
						break
					}
				}
			}
		}
	}
	return
}

func GetOptionHostCancelChargeRefund(ctx context.Context, server *Server, user db.User, funcName string, chargeID uuid.UUID, payoutRedisIDs []string, refundRedisIDs []string, userID uuid.UUID) (charge db.GetChargeOptionReferenceByHostIDRow, refund int, hostPayout int, err error) {
	charge, err = server.store.GetChargeOptionReferenceByHostID(ctx, db.GetChargeOptionReferenceByHostIDParams{
		UserID:     userID,
		Cancelled:  false,
		IsComplete: true,
		ID:         chargeID,
		HostUserID: user.UserID,
	})
	if err != nil {
		log.Printf("FuncName: %v. There an error at GetOptionEventCancelChargeRefund at GetChargeOptionReferenceByHostID: %v, chargeID: %v, userID: %v \n", funcName, err.Error(), chargeID, user.ID)
		err = fmt.Errorf("could not find this booking. Please try again later or try contacting us")
		return
	}
	// Once a host is paid you cannot perform a cancel
	switch true {
	case charge.MainPayoutComplete:
		err = fmt.Errorf("your cancellation cannot be made because you have been paid")
		return
	case time.Now().After(charge.StartDate):
		err = fmt.Errorf("your cancellation cannot be made because your guest has already started staying")
		return
	case ChargeIDInPay(payoutRedisIDs, chargeID):
		err = fmt.Errorf("your cancellation cannot be made because your payment is currently being processed")
		return
	case ChargeIDInPay(refundRedisIDs, chargeID):
		err = fmt.Errorf("your cancellation cannot be made because a refund is currently being processed")
		return
	case time.Now().After(charge.EndDate.Add(time.Hour * -10)):
		err = fmt.Errorf("your cancellation cannot be made because your guest stay is already coming to an end")
		return
	}
	// Because is the host cancelling refund will always be 100%
	refund = 100
	hostPayout = 0
	return

}

func GetEventUserCancelChargeRefund(ctx context.Context, server *Server, user db.User, funcName string, chargeID uuid.UUID, payoutRedisIDs []string, refundRedisIDs []string) (charge db.GetChargeTicketReferenceByUserIDRow, refund int, hostPayout int, err error) {
	charge, err = server.store.GetChargeTicketReferenceByUserID(ctx, db.GetChargeTicketReferenceByUserIDParams{
		UserID:     user.UserID,
		Cancelled:  false,
		IsComplete: true,
		ID:         chargeID,
	})
	if err != nil {
		log.Printf("FuncName: %v. There an error at HandleEventUserCancel at GetChargeTicketReferenceByUserID: %v, chargeID: %v, userID: %v \n", funcName, err.Error(), chargeID, user.ID)
		err = fmt.Errorf("could not find your booking. Please try again later or try contacting us")
		return
	}
	// Once a host is paid you cannot perform a cancel
	switch true {
	case time.Now().After(charge.EndDate.Add(time.Hour * -2)):
		err = fmt.Errorf("your cancellation cannot be made because your event has ended or is about to end")
		return
	}
	// Lets check if the guest has refund available
	startDate := tools.ConvertToTimeZoneTwo(charge.StartDate, charge.TimeZone)
	cancelPolicy, exist := val.GetCancelPolicy(charge.CancelPolicyOne)
	log.Println("cancel policy", cancelPolicy)
	if ChargeIDInPay(refundRedisIDs, chargeID) {
		err = fmt.Errorf("your cancellation cannot be made because the booking's refund is currently processing")
		return
	}
	if charge.MainPayoutComplete || ChargeIDInPay(payoutRedisIDs, chargeID) {
		log.Println("here one ", charge.MainPayoutComplete)
		refund = 0
		hostPayout = 0
	} else {
		log.Println("here two ", charge.MainPayoutComplete)
		if exist {
			log.Println("here three ", charge.MainPayoutComplete)
			for _, p := range cancelPolicy.Items {
				log.Println("p ", p)
				if p.Type == "hard" {
					log.Println("hard ", p.Type)
					timeWithZone := tools.ConvertToTimeZoneTwo(time.Now(), charge.TimeZone)
					startTime := startDate.Add(time.Hour).Add(time.Hour * -time.Duration(p.HoursTwo))
					timeOfBooking := timeWithZone.Add(time.Hour).Before(charge.DateBooked.Add(time.Hour * time.Duration(p.Hours)))
					log.Println("p hard startTime ", startTime, " ", timeOfBooking)
					if time.Now().Before(startTime) && timeOfBooking {
						// We need to check the dateBooked against Hours
						refund = p.Percent
						hostPayout = 100 - refund
						break
					}
				} else {

					startTime := startDate.Add(time.Hour).Add(time.Hour * -time.Duration(p.Hours))
					log.Println("p startTime ", startTime)
					if time.Now().Before(startTime) {
						refund = p.Percent
						hostPayout = 100 - refund
						break
					}
				}
			}
		}
	}
	return
}

func GetEventHostChargeCancel(ctx context.Context, server *Server, user db.User, funcName string, startDateString string, endDateString string, eventDateID uuid.UUID, eventID uuid.UUID) (chargeIDs []uuid.UUID, eventDate db.GetEventDateTimeByUIDRow, err error) {
	startDate, err := tools.ConvertDateOnlyStringToDate(startDateString)
	if err != nil {
		err = fmt.Errorf("this date for the event is not in the right format")
		return
	}
	endDate, err := tools.ConvertDateOnlyStringToDate(endDateString)
	if err != nil {
		err = fmt.Errorf("this date for the event is not in the right format")
		return
	}
	// If an event is about to end you cannot also cancel it
	switch true {
	case time.Now().After(endDate.Add(time.Hour * -5)):
		err = fmt.Errorf("your cancellation cannot be made because your event will soon end")
		return
	}
	eventDate, err = server.store.GetEventDateTimeByUID(ctx, db.GetEventDateTimeByUIDParams{
		EventDateTimeID: eventDateID,
		UID:             user.ID,
		EventInfoID:     eventID,
	})
	if err != nil {
		log.Printf("FuncName: %v. There an error at HandleEventUserCancel at .CreateRefund: %v, startDate: %v, userID: %v \n", funcName, err.Error(), startDateString, user.ID)
		err = fmt.Errorf("this event date was not found")
		return
	}
	// Lets check if the particular event has data in redis that needs to be cancelled or updated
	keyCancel := fmt.Sprintf("%v&%v&%v&%v", tools.UuidToString(eventDate.EventDateTimeID), eventDate.Type, startDateString, "cancel")
	result, err := RedisClient.HExists(RedisContext, keyCancel, constants.REFERENCE).Result()
	if err != nil {
		log.Printf("FuncName: %v. Cancel There an error at RedisClient.SMembers: %v, startDate: %v, userID: %v \n", funcName, err.Error(), startDateString, user.ID)
		err = fmt.Errorf("an error occurred while checking details")
		return
	}
	if result {
		err = fmt.Errorf("there is already a cancellation process in progress for this event, please try again is the cancellation fails")
		return
	}
	keyUpdate := fmt.Sprintf("%v&%v&%v&%v", tools.UuidToString(eventDate.EventDateTimeID), eventDate.Type, startDateString, "update")
	result, err = RedisClient.HExists(RedisContext, keyUpdate, constants.REFERENCE).Result()
	if err != nil {
		log.Printf("FuncName: %v. Update There an error at RedisClient.SMembers RedisClient.SMembers(RedisContext, keyUpdate: %v, startDate: %v, userID: %v \n", "UpdateEventDatesBooked", err.Error(), startDateString, user.ID)
		err = fmt.Errorf("an error occurred while checking details")
		return
	}
	if result {
		err = fmt.Errorf("there is already an update process in progress for this event, please try again if the process fails")
		return
	}
	// Lets get all the current charge ticket ids for this event
	chargeIDs, err = server.store.ListChargeTicketReferenceIDByStartDate(ctx, db.ListChargeTicketReferenceIDByStartDateParams{
		Date:        startDate,
		Cancelled:   false,
		IsComplete:  true,
		EventDateID: eventDate.EventDateTimeID,
	})
	if err != nil {
		log.Printf("FuncName: %v. There an error at HandleEventUserCancel at .ListChargeTicketReferenceIDByStartDate: %v, startDate: %v, userID: %v \n", funcName, err.Error(), startDateString, user.ID)
		err = fmt.Errorf("unable to find booked tickets for this event date")
		return
	}
	return
}
