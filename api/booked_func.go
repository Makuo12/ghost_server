package api

import (
	"context"
	"fmt"
	"log"

	"github.com/makuo12/ghost_server/constants"
	db "github.com/makuo12/ghost_server/db/sqlc"
	"github.com/makuo12/ghost_server/tools"

	"github.com/google/uuid"
)

func GetEventDatesBooked(ctx context.Context, server *Server, eventDateTimeID uuid.UUID, user db.User, funcName string) (dates []db.ListChargeDateReferenceDatesRow) {
	dates, err := server.store.ListChargeDateReferenceDates(ctx, db.ListChargeDateReferenceDatesParams{
		EventDateID: eventDateTimeID,
		IsComplete:  true,
		Cancelled:   false,
	})
	if err != nil || len(dates) == 0 {
		if err != nil {
			log.Printf("Error at GetEventDatesBooked .ListChargeDateReferenceDates in funcName: %v err: %v, user: %v\n", funcName, err, user.ID)
		}
		dates = []db.ListChargeDateReferenceDatesRow{}
	}
	return
}

func IsUpdatedEventDatesBooked(ctx context.Context, server *Server, dates []string, user db.User, eventDateTimeID uuid.UUID, funcName string) (valid bool) {
	datesData := GetEventDatesBooked(ctx, server, eventDateTimeID, user, funcName)
	if len(datesData) != 0 {
		for _, d := range dates {
			for _, booked := range datesData {
				startDate := tools.ConvertDateOnlyToString(booked.StartDate)
				endDate := tools.ConvertDateOnlyToString(booked.EndDate)
				if d == startDate || d == endDate {
					valid = false
					return
				}
			}
		}
	}
	valid = true
	return
}

func EventDateIsBooked(ctx context.Context, server *Server, eventDateTimeID uuid.UUID, funcName string, id string) (isBooked bool) {
	count, err := server.store.CountChargeTicketReferenceByEventDateID(ctx, db.CountChargeTicketReferenceByEventDateIDParams{
		EventDateID: eventDateTimeID,
		Cancelled:   false,
		IsComplete:  true,
	})
	if err != nil {
		log.Printf("Error at EventDateIsBooked in CountChargeTicketReferenceByEventDateID funcName: %v, eventDateTimeID: %v, err: %v, id: %v\n", funcName, eventDateTimeID, err, id)

		return
	}
	isBooked = count > 0
	return
}

func EventDateIsBookedAny(ctx context.Context, server *Server, eventDateTimeID uuid.UUID, funcName string, id string) (isBooked bool) {
	count, err := server.store.CountChargeTicketReferenceByEventDateIDAny(ctx, eventDateTimeID)
	if err != nil {
		log.Printf("Error at EventDateIsBookedAny in CountChargeTicketReferenceByEventDateID funcName: %v, eventDateTimeID: %v, err: %v, id: %v\n", funcName, eventDateTimeID, err, id)
		return
	}
	isBooked = count > 0
	return
}

func EventDateTicketIsBooked(ctx context.Context, server *Server, ticketID uuid.UUID, funcName string, id string) (isBooked bool) {
	ticketCount, err := server.store.CountChargeTicketReference(ctx, db.CountChargeTicketReferenceParams{
		TicketID:   ticketID,
		Cancelled:  false,
		IsComplete: true,
	})
	if err != nil {
		log.Printf("Error at EventDateIsBooked in CountChargeTicketReferenceByEventDateID funcName: %v, ticketID: %v, err: %v, id: %v\n", funcName, ticketID, err, id)

		return
	}
	isBooked = ticketCount > 0
	return
}

func EventDateTicketIsBookedAny(ctx context.Context, server *Server, ticketID uuid.UUID, funcName string, id string) (isBooked bool) {
	ticketCount, err := server.store.CountChargeTicketReferenceAny(ctx, ticketID)
	if err != nil {
		log.Printf("Error at EventDateIsBookedAny in CountChargeTicketReferenceByEventDateID funcName: %v, ticketID: %v, err: %v, id: %v\n", funcName, ticketID, err, id)

		return
	}
	isBooked = ticketCount > 0
	return
}

func OptionDateIsBooked(ctx context.Context, server *Server, optionUserID uuid.UUID, funcName string, id string) (isBooked bool) {
	count, err := server.store.CountChargeOptionReferenceBook(ctx, db.CountChargeOptionReferenceBookParams{
		OptionUserID: optionUserID,
		IsComplete:   true,
		Cancelled:    false,
	})
	if err != nil {
		log.Printf("Error at OptionDateIsBooked in CountChargeOptionReferenceBook funcName: %v, eventDateTimeID: %v, err: %v, id: %v\n", funcName, optionUserID, err, id)

		return
	}
	isBooked = count > 0
	return
}

func GetOptionDatesBooked(ctx context.Context, server *Server, optionUserID uuid.UUID, funcName string, id string) (dates []db.ListChargeOptionReferenceBookRow) {
	dates, err := server.store.ListChargeOptionReferenceBook(ctx, db.ListChargeOptionReferenceBookParams{
		OptionUserID: optionUserID,
		IsComplete:   true,
		Cancelled:    false,
	})
	if err != nil {
		log.Printf("Error at GetOptionDatesBooked in CountChargeOptionReferenceBook funcName: %v, eventDateTimeID: %v, err: %v, id: %v\n", funcName, optionUserID, err, id)
		dates = []db.ListChargeOptionReferenceBookRow{}
	}
	return
}

// This makes checks if the given capacity is ok for the ticket
// We need to make sure the quantity given isn't less than the ticket sold
func TicketQuantity(server *Server, ctx context.Context, capacity int, ticketID uuid.UUID, eventDateTimeID uuid.UUID) (quantityOk bool, err error) {
	eventDate, err := server.store.GetEventDateTime(ctx, eventDateTimeID)
	if err != nil {
		log.Printf("Error at TicketQuantity in GetEventDateTime: %v eventDateIDString: %v ticketIDString: %v\n", err.Error(), eventDateTimeID, ticketID)
		return
	}
	//if ticket.Type == "free" || ticket.Type == "donation_allowed" {
	//	err = nil
	//	return
	//}
	if eventDate.Type == "single" {
		ticketCount, errCharge := server.store.CountChargeTicketReference(ctx, db.CountChargeTicketReferenceParams{
			TicketID:   ticketID,
			Cancelled:  false,
			IsComplete: true,
		})
		if errCharge != nil {
			log.Printf("Error at TicketQuantity in CountChargeTicketReference %v eventDateIDString: %v ticketIDString: %v\n", errCharge.Error(), eventDateTimeID, ticketID)
			if err == db.ErrorRecordNotFound {
				err = nil
				quantityOk = true
				return
			} else {
				return
			}

		}
		if capacity < int(ticketCount) {
			err = fmt.Errorf("there are %v tickets sold, the capacity cannot be less than the tickets sold", ticketCount)
			quantityOk = false
		} else {
			err = nil
			quantityOk = true
		}

	} else {
		// We loop through the dates in ticket and try to get the one that has sold the most tickets
		var maxTicketCount int
		for _, d := range eventDate.EventDates {
			date, err := tools.ConvertDateOnlyStringToDate(d)
			if err != nil {
				log.Printf("Error at TicketQuantity in tools.ConvertDateOnlyStringToDate %v eventDateIDString: %v ticketIDString: %v\n", err.Error(), eventDateTimeID, ticketID)
				if err == db.ErrorRecordNotFound {
					err = nil
					continue
				} else {
					break
				}
			}
			ticketCount, errCharge := server.store.CountChargeTicketReferenceByStartDate(ctx, db.CountChargeTicketReferenceByStartDateParams{
				Date:       date,
				TicketID:   ticketID,
				Cancelled:  false,
				IsComplete: true,
			})
			if errCharge != nil {
				log.Printf("Error at TicketQuantity in store.CountChargeTicketReferenceByStartDate %v eventDateIDString: %v ticketIDString: %v\n", errCharge.Error(), eventDateTimeID, ticketID)
				if errCharge == db.ErrorRecordNotFound {
					err = nil
					continue
				} else {
					break
				}
			}
			if maxTicketCount < int(ticketCount) {
				maxTicketCount = int(ticketCount)
			}
		}
		if capacity < int(maxTicketCount) {
			err = fmt.Errorf("there are %v tickets sold, the capacity cannot be less than the tickets sold", maxTicketCount)
			quantityOk = false
		} else {
			err = nil
			quantityOk = true
		}
	}
	return
}

func DailyChangeDateEventHostUpdate(ctx context.Context, server *Server) func() {

	return func() {
		result, err := RedisClient.SMembers(RedisContext, constants.CHARGE_TICKET_ID_UPDATE).Result()
		if err != nil || len(result) == 0 {
			if err != nil {
				log.Printf("There an error at DailyChangeDateEventHostUpdate at RedisClient.SMembers: %v, type: %v \n", err.Error(), constants.CHARGE_TICKET_ID_UPDATE)
			}
			return
		}
		for _, uniqueID := range result {
			ProcessEventHostUpdateUniqueID(ctx, server, uniqueID, "DailyChangeDateEventHostUpdate")
		}
	}
}

func ProcessEventHostUpdateUniqueID(ctx context.Context, server *Server, uniqueID string, funcName string) {
	data, err := RedisClient.HGetAll(RedisContext, uniqueID).Result()
	if err != nil {
		log.Printf("FuncName: %v. There an error at ProcessEventHostUpdateUniqueID at RedisClient.HGetAll: %v, uniqueID: %v \n", funcName, err.Error(), uniqueID)
		return
	}
	startDate, err := tools.ConvertDateOnlyStringToDate(data[constants.START_DATE])
	if err != nil {
		log.Printf("FuncName: %v. At startDate There an error at ProcessEventHostUpdateUniqueID at tools.ConvertDateOnlyStringToDate: %v, uniqueID: %v \n", funcName, err.Error(), uniqueID)
		return
	}

	endDate, err := tools.ConvertDateOnlyStringToDate(data[constants.END_DATE])
	if err != nil {
		log.Printf("FuncName: %v. At endDate There an error at ProcessEventHostUpdateUniqueID at tools.ConvertDateOnlyStringToDate: %v, uniqueID: %v \n", funcName, err.Error(), uniqueID)
		endDate = startDate
	}

	result, err := RedisClient.SMembers(RedisContext, data[constants.REFERENCE]).Result()
	if err != nil {
		log.Printf("FuncName: %v. There an error at ProcessEventHostUpdateUniqueID at RedisClient.SMembers: %v, data[constants.REFERENCE]: %v \n", funcName, err.Error(), data[constants.REFERENCE])
		return
	}
	var redisRemovedCount int
	for _, id := range result {
		chargeID, err := tools.StringToUuid(id)
		if err != nil {
			log.Printf("FuncName: %v. There an error at ProcessEventHostUpdateUniqueID at tools.StringToUuid: %v, data[constants.REFERENCE]: %v, chargeID: %v \n", funcName, err.Error(), data[constants.REFERENCE], id)
			continue
		}
		charge, err := server.store.GetChargeTicketReferenceByChargeID(ctx, db.GetChargeTicketReferenceByChargeIDParams{
			Cancelled:  false,
			IsComplete: true,
			ID:         chargeID,
		})
		if err != nil {
			log.Printf("FuncName: %v. There an error at ProcessEventHostUpdateUniqueID at GetChargeTicketReferenceByChargeID: %v, data[constants.REFERENCE]: %v, chargeID: %v \n", funcName, err.Error(), data[constants.REFERENCE], chargeID)
			continue
		}
		// We want to update
		// We update the charge date to the current date
		_, err = server.store.UpdateChargeDateReferenceDates(ctx, db.UpdateChargeDateReferenceDatesParams{
			StartDate: startDate,
			EndDate:   endDate,
			ID:        charge.DateTimeID,
		})
		if err != nil {
			log.Printf("FuncName: %v. There an error at ProcessEventHostUpdateUniqueID at UpdateChargeDateReferenceDate: %v, data[constants.REFERENCE]: %v, chargeID: %v \n", funcName, err.Error(), data[constants.REFERENCE], chargeID)
			continue
		}
		header := fmt.Sprintf("%v's, changed dates for %v event", tools.CapitalizeFirstCharacter(charge.HostFirstName), tools.HandleReadableDates(charge.StartDate, charge.EndDate, tools.DateDMM))
		CreateTypeNotification(ctx, server, charge.ChargeID, charge.HostUserID, constants.USER_CANCEL, data[constants.MESSAGE], false, header)

		err = RedisClient.SRem(RedisContext, data[constants.REFERENCE], id).Err()
		if err != nil {
			log.Printf("FuncName: %v. There an error at ProcessEventHostUpdateUniqueID at .SRem: %v, uniqueID: %v, chargeID: %v \n", funcName, err.Error(), uniqueID, chargeID)
			continue
		} else {
			redisRemovedCount++
		}
	}
	if redisRemovedCount == len(result) {
		err = RedisClient.Del(RedisContext, uniqueID).Err()
		if err != nil {
			log.Printf("FuncName: %v. There an error at ProcessEventHostUpdateUniqueID at RedisClient.Del: %v, uniqueID: %v \n", funcName, err.Error(), uniqueID)
			return
		}
		err = RedisClient.SRem(RedisContext, constants.CHARGE_TICKET_ID_UPDATE, uniqueID).Err()
		if err != nil {
			log.Printf("FuncName: %v. There an error at ProcessEventHostUpdateUniqueID at RedisClient.Del: %v, uniqueID: %v \n", funcName, err.Error(), uniqueID)
			return
		}
	}
}
