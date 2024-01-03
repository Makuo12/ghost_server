package api

import (
	"context"
	"flex_server/constants"
	db "flex_server/db/sqlc"
	"flex_server/tools"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func HandleOptionPayouts(ctx context.Context, server *Server, redisDateIDs []string, refundIDs []string) {
	dollarToNaira := server.config.DollarToNaira
	dollarToCAD := server.config.DollarToCAD
	payouts, err := GetOptionPayout(ctx, server, dollarToNaira, dollarToCAD, redisDateIDs)
	if err != nil || len(payouts) == 0 {
		if err != nil {
			log.Printf("Error at HandleOptionPayouts in GetOptionPayout err:%v\n", err.Error())
		}
		return
	}
	var transferData = make(map[uuid.UUID]PayoutData)
	var i = 1
	for key, value := range payouts {
		transferData[key] = value
		i++
		if i == len(payouts) {
			bulkTransfer, resPayouts := HandleBulkTransfer(ctx, server, transferData, "main_payout")
			res, err := TransferByPaystack(ctx, server, bulkTransfer)
			if err != nil {
				log.Printf("Error at HandleOptionPayouts in TransferByPaystack err:%v\n", err.Error())
			} else {
				log.Println("Transfer response from HandleOptionPayouts", res)
			}
			HandlePayoutRes(res, redisDateIDs, resPayouts)
		} else if i%100 == 0 {
			bulkTransfer, resPayouts := HandleBulkTransfer(ctx, server, transferData, "main_payout")
			res, err := TransferByPaystack(ctx, server, bulkTransfer)
			if err != nil {
				log.Printf("Error at HandleOptionPayouts in TransferByPaystack err:%v\n", err.Error())
			} else {
				log.Println("Transfer response from HandleOptionPayouts", res)
			}
			HandlePayoutRes(res, redisDateIDs, resPayouts)
			duration := 7 * time.Second
			time.Sleep(duration)
			transferData = make(map[uuid.UUID]PayoutData)
		}
	}

}

func HandleEventPayouts(ctx context.Context, server *Server, redisDateIDs []string, refundIDs []string) {
	dollarToNaira := server.config.DollarToNaira
	dollarToCAD := server.config.DollarToCAD
	payouts, err := GetEventPayout(ctx, server, dollarToNaira, dollarToCAD, redisDateIDs)
	log.Println("payouts ", payouts)
	if err != nil || len(payouts) == 0 {
		if err != nil {
			log.Printf("Error at HandleEventPayouts in GetEventPayout err:%v\n", err.Error())
		}
		return
	}
	var transferData = make(map[uuid.UUID]PayoutData)
	var i = 1
	for key, value := range payouts {
		transferData[key] = value

		if i == len(payouts) {
			bulkTransfer, resPayouts := HandleBulkTransfer(ctx, server, transferData, "main_payout")
			log.Println("bulk transfer: ", bulkTransfer)
			res, err := TransferByPaystack(ctx, server, bulkTransfer)
			if err != nil {
				log.Printf("Error at HandleEventPayouts in TransferByPaystack err:%v\n", err.Error())
			} else {
				log.Println("Transfer response from HandleEventPayouts", res)
			}
			HandlePayoutRes(res, redisDateIDs, resPayouts)
		} else if i%100 == 0 {
			bulkTransfer, resPayouts := HandleBulkTransfer(ctx, server, transferData, "main_payout")
			res, err := TransferByPaystack(ctx, server, bulkTransfer)
			if err != nil {
				log.Printf("Error at HandleEventPayouts in TransferByPaystack err:%v\n", err.Error())
			} else {
				log.Println("Transfer response from HandleEventPayouts", res)
			}
			HandlePayoutRes(res, redisDateIDs, resPayouts)
			duration := 7 * time.Second
			time.Sleep(duration)
			transferData = make(map[uuid.UUID]PayoutData)
		}
		i++
	}

}

func DailyHandlePayouts(ctx context.Context, server *Server) func() {
	// All the ids are stored in constants.USER_REQUEST_APPROVE
	return func() {
		redisChargeDateIDs, err := RedisClient.SMembers(RedisContext, constants.PAYOUT_CHARGE_DATE_IDS).Result()
		if err != nil {
			log.Printf("Error at DailyHandlePayouts in SMembers err:%v\n", err.Error())
			return
		}
		//redisRefundChargeDateIDs, err := RedisClient.SMembers(RedisContext, constants.REFUND_CHARGE_DATE_IDS).Result()
		//if err != nil {
		//	log.Printf("Error at DailyHandlePayouts in SMembers err:%v\n", err.Error())
		//	return
		//}
		// First we start by handling options
		HandleOptionPayouts(ctx, server, redisChargeDateIDs, []string{})
		// Next we handle events payouts
		HandleEventPayouts(ctx, server, redisChargeDateIDs, []string{})
	}
}

func HandleTransferSuccess(ctx context.Context, server *Server, data map[string]string, reference string, funcName string, redisIDs []string) {
	split := strings.Split(reference, "&")
	if len(split) != 2 {
		return
	}
	mainReference, err := tools.StringToUuid(split[1])
	if err != nil {
		log.Printf("Error at DailyHandleWebhookData in tools.StringToUuid(r) err:%v, reference: %v\n", err.Error(), reference)
		return
	}
	payout, err := server.store.UpdatePayout(ctx, db.UpdatePayoutParams{
		TransferCode: pgtype.Text{
			String: data[constants.TRANSFER_CODE],
			Valid:  true,
		},
		AmountPayed: pgtype.Int8{
			Int64: tools.MoneyStringToInt(data[constants.AMOUNT]),
			Valid: true,
		},
		IsComplete: pgtype.Bool{
			Bool:  true,
			Valid: true,
		},
		TimePaid: pgtype.Timestamptz{
			Time:  time.Now(),
			Valid: true,
		},
		ID: mainReference,
	})
	if err != nil {
		log.Printf("Error at funcName: %v, HandleTransferSuccess in .UpdatePayout err:%v, reference: %v\n", funcName, err.Error(), reference)
	}
	// We want to update main_payout chargeIDs
	switch payout.ParentType {
	case "main_payout":
		for _, chargeID := range payout.PayoutIds {
			err = server.store.UpdateMainPayout(ctx, db.UpdateMainPayoutParams{
				IsComplete:    true,
				ChargeID:      chargeID,
				AccountNumber: payout.AccountNumber,
				TimePaid:      payout.TimePaid,
			})
			if err != nil {
				log.Printf("Error at funcName: %v, HandleTransferSuccess in UpdateMainPayout err:%v, reference: %v\n", funcName, err.Error(), reference)
			}
			// Get Redis ids to remove
			rCID, exist := GetRedisChargeID(redisIDs, chargeID)
			if exist {
				err := RedisClient.SRem(RedisContext, constants.PAYOUT_CHARGE_DATE_IDS, rCID).Err()
				if err != nil {
					log.Println("This redis  HandleTransferSuccess items were not removed from redis even though payout was unsuccessful", chargeID)
				}
			}

		}
	case "refund_payout":
		for _, chargeID := range payout.PayoutIds {
			err = server.store.UpdateRefundPayout(ctx, db.UpdateRefundPayoutParams{
				IsComplete:    true,
				ChargeID:      chargeID,
				AccountNumber: payout.AccountNumber,
				TimePaid:      payout.TimePaid,
			})
			if err != nil {
				log.Printf("Error at funcName: %v, HandleTransferSuccess in UpdateRefundPayout err:%v, reference: %v\n", funcName, err.Error(), reference)
			}
			// Get Redis ids to remove
			rCID, exist := GetRedisChargeID(redisIDs, chargeID)
			if exist {
				err := RedisClient.SRem(RedisContext, constants.PAYOUT_CHARGE_DATE_IDS, rCID).Err()
				if err != nil {
					log.Println("This redis  HandleTransferSuccess items were not removed from redis even though payout was unsuccessful", chargeID)
				}
			}

		}
	}
	err = RedisClient.Del(RedisContext, reference).Err()
	if err != nil {
		log.Printf("Error at funcName: %v, HandleTransferSuccess in RedisClient.Del err:%v, reference: %v\n", funcName, err.Error(), reference)
	}
	err = RedisClient.SRem(RedisContext, constants.WEBHOOK_PAYSTACK_TRANSFER_REFERENCE, reference).Err()
	if err != nil {
		log.Printf("Error at funcName: %v, HandleTransferSuccess in RedisClient.SRem err:%v, reference: %v\n", funcName, err.Error(), reference)
	}

}

func HandleTransferNotSuccess(ctx context.Context, server *Server, data map[string]string, reference string, funcName string, redisIDs []string) {
	split := strings.Split(reference, "&")
	if len(split) != 2 {
		return
	}
	mainReference, err := tools.StringToUuid(split[1])
	if err != nil {
		log.Printf("Error at DailyHandleWebhookData in tools.StringToUuid(r) err:%v, reference: %v\n", err.Error(), reference)
		return
	}
	// We only update the TransferCode because a transfer was still made even though is was not successful
	payout, err := server.store.UpdatePayout(ctx, db.UpdatePayoutParams{
		TransferCode: pgtype.Text{
			String: data[constants.TRANSFER_CODE],
			Valid:  true,
		},
		ID: mainReference,
	})
	if err != nil {
		log.Printf("Error at funcName: %v, HandleTransferNotSuccess in .UpdatePayout err:%v, mainReference: %v\n", funcName, err.Error(), mainReference)
	}
	// We want to update main_payout chargeIDs
	for _, chargeID := range payout.PayoutIds {
		// Get Redis ids to remove
		rCID, exist := GetRedisChargeID(redisIDs, chargeID)
		if exist {
			err := RedisClient.SRem(RedisContext, constants.PAYOUT_CHARGE_DATE_IDS, rCID).Err()
			if err != nil {
				log.Println("This redis  HandleTransferSuccess items were not removed from redis even though payout was unsuccessful", chargeID)
			}
		}
	}
	err = RedisClient.Del(RedisContext, reference).Err()
	if err != nil {
		log.Printf("Error at funcName: %v, HandleTransferNotSuccess in RedisClient.Del err:%v, reference: %v\n", funcName, err.Error(), reference)
	}
	err = RedisClient.SRem(RedisContext, constants.WEBHOOK_PAYSTACK_TRANSFER_REFERENCE, reference).Err()
	if err != nil {
		log.Printf("Error at funcName: %v, HandleTransferNotSuccess in RedisClient.SRem err:%v, reference: %v\n", funcName, err.Error(), reference)
	}

}
