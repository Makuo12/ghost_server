package api

import (
	"context"
	"flex_server/constants"
	db "flex_server/db/sqlc"
	"flex_server/tools"
	"fmt"
	"log"
	"time"
)

// This function remove the reserve option stored in redis
func DailyRemoveOptionReserveUser() {
	// All the references are stored in constants.REMOVE_OPTION_RESERVE_USER
	refs, err := RedisClient.SMembers(RedisContext, constants.REMOVE_OPTION_RESERVE_USER).Result()
	if err != nil {
		log.Printf("DailyRemoveOptionReserveUser in SMembers err:%v\n", err.Error())
		return
	}
	for _, ref := range refs {
		timeData, err := RedisClient.HGetAll(RedisContext, ref).Result()
		if err != nil {
			log.Printf("DailyRemoveOptionReserveUser in HGetAll reference: %v, err:%v\n", ref, err.Error())
			continue
		}
		t, err := tools.ConvertStringToTime(timeData[constants.TIME])
		if err != nil {
			log.Printf("DailyRemoveOptionReserveUser in ConvertStringToTime reference: %v, %v err:%v, time: \n", ref, err.Error(), timeData[constants.TIME])
			continue
		}
		// We run this after an hour
		// If the time in redis after the current time then we run the function
		if time.Now().Add(time.Hour * 1).After(t) {
			err = TimerRemoveOptionReserveUser(timeData[constants.MAIN_REFERENCE])
			// Lets remove timeData
			if err == nil {
				err = RedisClient.Del(RedisContext, ref).Err()
				if err != nil {
					log.Printf("DailyRemoveOptionReserveUser in RedisClient.Del reference: %v, err:%v\n", ref, err.Error())
				}
				err = RedisClient.SRem(RedisContext, constants.REMOVE_OPTION_RESERVE_USER, ref).Err()
				if err != nil {
					log.Printf("DailyRemoveOptionReserveUser in RedisClient.SRem reference: %v, err:%v\n", ref, err.Error())
				}
			}

		}
	}
}

func DailyRemoveEventReserveUser() {
	// All the references are stored in constants.REMOVE_EVENT_RESERVE_USER
	refs, err := RedisClient.SMembers(RedisContext, constants.REMOVE_EVENT_RESERVE_USER).Result()
	if err != nil {
		log.Printf("DailyHandleRedisMsgToDB in SMembers err:%v\n", err.Error())
		return
	}
	for _, ref := range refs {
		timeData, err := RedisClient.HGetAll(RedisContext, ref).Result()
		if err != nil {
			log.Printf("DailyHandleRedisMsgToDB in HGetAll reference: %v, err:%v\n", ref, err.Error())
			continue
		}
		t, err := tools.ConvertStringToTime(timeData[constants.TIME])
		if err != nil {
			log.Printf("DailyHandleRedisMsgToDB in ConvertStringToTime reference: %v, %v err:%v, time: \n", ref, err.Error(), timeData[constants.TIME])
			continue
		}
		// We run this after an hour
		// If the time in redis after the current time then we run the function
		if time.Now().Add(time.Hour * 1).After(t) {
			err = TimerRemoveEventReserveUser(timeData[constants.MAIN_REFERENCE])
			if err == nil {
				err = RedisClient.Del(RedisContext, ref).Err()
				if err != nil {
					log.Printf("DailyRemoveEventReserveUser in RedisClient.Del reference: %v, err:%v\n", ref, err.Error())
				}
				err = RedisClient.SRem(RedisContext, constants.REMOVE_EVENT_RESERVE_USER, ref).Err()
				if err != nil {
					log.Printf("DailyRemoveEventReserveUser in RedisClient.SRem reference: %v, err:%v\n", ref, err.Error())
				}
			}
		}
	}
}



func DailyHandleUserRequest(ctx context.Context, server *Server) func() {
	// All the ids are stored in constants.USER_REQUEST_APPROVE
	return func() {
		ids, err := RedisClient.SMembers(RedisContext, constants.USER_REQUEST_APPROVE).Result()
		if err != nil {
			log.Printf("DailyHandleUserRequest in SMembers err:%v\n", err.Error())
			return
		}
		for _, id := range ids {
			timeData, err := RedisClient.HGetAll(RedisContext, id).Result()
			if err != nil {
				log.Printf("DailyHandleUserRequest in HGetAll id: %v, err:%v\n", id, err.Error())
				return
			}
			t, err := tools.ConvertStringToTime(timeData[constants.TIME])
			if err != nil {
				log.Printf("DailyHandleUserRequest in ConvertStringToTime id: %v, %v err:%v, time: \n", id, err.Error(), timeData[constants.TIME])
				continue
			}
			mID, err := tools.StringToUuid(timeData[constants.MID])
			if err != nil {
				log.Printf("DailyHandleUserRequest in tools.StringToUuid id: %v, %v err:%v, time: \n", id, err.Error(), timeData[constants.TIME])
				continue
			}
			// We run this after an hour
			// If the time in redis after the current time then we run the function
			if time.Now().Add(time.Hour * 1).After(t) {
				HandleURApproved(ctx, server, mID, timeData[constants.SENDER_ID], timeData[constants.RECEIVER_ID], timeData[constants.FIRSTNAME], timeData[constants.REFERENCE])
				err = RedisClient.Del(RedisContext, id).Err()
				if err != nil {
					log.Printf("DailyHandleUserRequest in RedisClient.Del id: %v, err:%v\n", id, err.Error())
				}
			
			}
		}
	}
}

func DailyHandleSnooze(ctx context.Context, server *Server) func() {
	// All the ids are stored in constants.USER_REQUEST_APPROVE
	return func() {
		err := server.store.UpdateUnSnoozeStatus(ctx)
		if err != nil {
			log.Printf("DailyHandleSnooze in server.store.UpdateUnSnoozeStatus, err:%v\n", err.Error())
		}

		err = server.store.UpdateSnoozeStatus(ctx)
		if err != nil {
			log.Printf("DailyHandleSnooze in server.store.UpdateSnoozeStatus, err:%v\n", err.Error())
		}
	}
}

func DailyDeactivateCoHost(ctx context.Context, server *Server) func() {
	return func() {
		// First lets get the list of deactivated accounts
		deactivatesID, err := RedisClient.SMembers(RedisContext, constants.DEACTIVATE_CO_HOST_IDS).Result()
		if err != nil || len(deactivatesID) == 0 {
			if err != nil {
				log.Printf("Error at DailyDeactivateCoHost in CreateNotification RedisClient.SMembers: %v \n", err.Error())
			}
			return
		}
		for _, id := range deactivatesID {
			coHostID, err := tools.StringToUuid(id)
			if err != nil {
				log.Printf("DailyDeactivateCoHost in tools.StringToUuid id: %v, err:%v, time: \n", id, err.Error())
				continue
			}
			coHost, err := server.store.GetDeactivateOptionCOHostByID(ctx, coHostID)
			if err != nil {
				log.Printf("DailyDeactivateCoHost in GetDeactivateOptionCOHostByID: %v, err:%v, time: \n", id, err.Error())
				continue
			}
			topHeader := "End to co-hosting session."
			header := fmt.Sprintf("%v, has ending the co-hosting session for %v.", coHost.CoUserFirstName, coHost.HostNameOption)
			msg := fmt.Sprintf("Hey %v, %v has ending the co-hosting session for %v so all access you gave to %v would be removed.", coHost.MainHostName, coHost.CoUserFirstName, coHost.HostNameOption, coHost.CoUserFirstName)
			// We would send an email
			err = SendCustomEmail(server, coHost.MainUserEmail, coHost.MainHostName, header, topHeader, msg, "DailyDeactivateCoHost")
			if err != nil {
				log.Printf("DailyDeactivateCoHost in SendCustomEmail: %v, err:%v, time: \n", id, err.Error())
			}
			// We would send a notification
			CreateTypeNotification(ctx, server, coHostID, coHost.MainUserID, constants.CO_HOST_DEACTIVATE, msg, false, topHeader)

			err = RedisClient.SRem(RedisContext, constants.DEACTIVATE_CO_HOST_IDS, id).Err()
			if err != nil {
				log.Printf("DailyDeactivateCoHost in RedisClient.SRem: %v, err:%v, time: \n", id, err.Error())
			}

		}

	}
}

func DailyValidatedChargeTicket(ctx context.Context, server *Server) func() {
	return func() {
		// First lets get the list of tickets
		ticketChargeIDs, err := RedisClient.SMembers(RedisContext, constants.SCANNED_CHARGE_TICKET_ID).Result()
		if err != nil {
			log.Printf("Error at DailyValidatedChargeTicket in RedisClient.SMembers: %v \n", err.Error())
		}
		if len(ticketChargeIDs) != 0 {
			// First we would remove it from Redis just to make sure the next cron job doesn't run it
			err = RedisClient.SRem(RedisContext, constants.SCANNED_CHARGE_TICKET_ID, ticketChargeIDs).Err()
			if err != nil {
				log.Printf("Error at DailyValidatedChargeTicket in RedisClient.SRem(: %v ids: %v \n", err.Error(), ticketChargeIDs)
			} else {
				// We want to all send notifications if the data was removed
				for _, id := range ticketChargeIDs {
					chargeID, err := tools.StringToUuid(id)
					if err != nil {
						log.Printf("Error at DailyValidatedChargeTicket in tools.StringToUuid: %v chargeID: %v \n", err.Error(), id)
						continue
					}
					charge, err := server.store.GetScannedChargeTicketByID(ctx, db.GetScannedChargeTicketByIDParams{
						ChargeID:         chargeID,
						PaymentCompleted: true,
						Cancelled:        false,
						ChargeScanned:    true,
					})
					if err != nil {
						log.Printf("Error at DailyValidatedChargeTicket in GetScannedChargeTicketByID: %v chargeID: %v \n", err.Error(), id)
						continue
					}
					// We create the notification
					dateString := tools.HandleReadableDates(charge.StartDate, charge.EndDate, tools.DateDMM)
					header := "Check in successful"
					msg := fmt.Sprintf("Hey %v,\nYour ticket for %v was successfully scanned by %v.\nEvent: %v\nEvent date: %v", charge.UserName, charge.HostNameOption, charge.ScannedByName, charge.HostNameOption, dateString)
					CreateTypeNotification(ctx, server, charge.ChargeID, charge.UserID, constants.SCANNED_CHARGE_TICKET, msg, false, header)
				}
			}
		}
	}
}

func DailyValidatedChargeOption(ctx context.Context, server *Server) func() {
	return func() {
		// First lets get the list of tickets
		optionChargeIDs, err := RedisClient.SMembers(RedisContext, constants.SCANNED_CHARGE_OPTION_ID).Result()
		if err != nil {
			log.Printf("Error at DailyValidatedChargeOption in RedisClient.SMembers: %v \n", err.Error())
		}
		if len(optionChargeIDs) != 0 {
			// First we would remove it from Redis just to make sure the next cron job doesn't run it
			err = RedisClient.SRem(RedisContext, constants.SCANNED_CHARGE_OPTION_ID, optionChargeIDs).Err()
			if err != nil {
				log.Printf("Error at DailyValidatedChargeOption in RedisClient.SRem(: %v ids: %v \n", err.Error(), optionChargeIDs)
			} else {
				// We want to all send notifications if the data was removed
				for _, id := range optionChargeIDs {
					chargeID, err := tools.StringToUuid(id)
					if err != nil {
						log.Printf("Error at DailyValidatedChargeOption in tools.StringToUuid: %v chargeID: %v \n", err.Error(), id)
						continue
					}
					// We send the notification for options immediately
					charge, err := server.store.GetScannedChargeTicketByID(ctx, db.GetScannedChargeTicketByIDParams{
						ChargeID:         chargeID,
						PaymentCompleted: true,
						Cancelled:        false,
						ChargeScanned:    true,
					})
					if err != nil {
						log.Printf("Error at DailyValidatedChargeOption in GetScannedChargeTicketByID: %v chargeID: %v \n", err.Error(), id)
						continue
					}
					// We create the notification
					dateString := tools.HandleReadableDates(charge.StartDate, charge.EndDate, tools.DateDMM)
					header := "Check in successful"
					msg := fmt.Sprintf("Hey %v,\nYour code for %v was successfully scanned by %v.\nListing: %v\nStaying: %v", charge.UserName, charge.HostNameOption, charge.ScannedByName, charge.HostNameOption, dateString)
					CreateTypeNotification(ctx, server, chargeID, charge.UserID, constants.SCANNED_CHARGE_OPTION, msg, false, header)
				}
			}
		}
	}
}
