package api

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/makuo12/ghost_server/db/sqlc"
	"github.com/makuo12/ghost_server/payout"
	"github.com/makuo12/ghost_server/tools"
)

func HandleOptionPayouts(ctx context.Context, server *Server) {
	dollarToNaira := server.config.DollarToNaira
	dollarToCAD := server.config.DollarToCAD
	paystackKey := server.config.PaystackSecretLiveKey
	payouts, err := GetOptionPayout(ctx, server, dollarToNaira, dollarToCAD)
	if err != nil || len(payouts) == 0 {
		if err != nil {
			log.Printf("Error at HandleOptionPayouts in GetOptionPayout err:%v\n", err.Error())
		}
		return
	}
	var transferData = make(map[uuid.UUID]payout.PayoutData)
	var i = 1
	for key, value := range payouts {
		transferData[key] = value
		i++
		if i == len(payouts) {
			bulkTransfer, resPayouts := HandleBulkTransfer(ctx, server, transferData, "main_payout")
			res, err := payout.TransferByPaystack(ctx, paystackKey, bulkTransfer)
			if err != nil {
				log.Printf("Error at HandleOptionPayouts in TransferByPaystack err:%v\n", err.Error())
			} else {
				log.Println("Transfer response from HandleOptionPayouts", res)
			}
			HandlePayoutRes(ctx, server, res, resPayouts)
		} else if i%100 == 0 {
			bulkTransfer, resPayouts := HandleBulkTransfer(ctx, server, transferData, "main_payout")
			res, err := payout.TransferByPaystack(ctx, paystackKey, bulkTransfer)
			if err != nil {
				log.Printf("Error at HandleOptionPayouts in TransferByPaystack err:%v\n", err.Error())
			} else {
				log.Println("Transfer response from HandleOptionPayouts", res)
			}
			HandlePayoutRes(ctx, server, res, resPayouts)
			duration := 7 * time.Second
			time.Sleep(duration)
			transferData = make(map[uuid.UUID]payout.PayoutData)
		}
	}

}

func HandleEventPayouts(ctx context.Context, server *Server) {
	dollarToNaira := server.config.DollarToNaira
	dollarToCAD := server.config.DollarToCAD
	paystackKey := server.config.PaystackSecretLiveKey
	payouts, err := GetEventPayout(ctx, server, dollarToNaira, dollarToCAD)
	if err != nil || len(payouts) == 0 {
		if err != nil {
			log.Printf("Error at HandleEventPayouts in GetEventPayout err:%v\n", err.Error())
		}
		return
	}
	var transferData = make(map[uuid.UUID]payout.PayoutData)
	var i = 1
	for key, value := range payouts {
		transferData[key] = value

		if i == len(payouts) {
			bulkTransfer, resPayouts := HandleBulkTransfer(ctx, server, transferData, "main_payout")
			log.Println("bulk transfer: ", bulkTransfer)
			res, err := payout.TransferByPaystack(ctx, paystackKey, bulkTransfer)
			if err != nil {
				log.Printf("Error at HandleEventPayouts in TransferByPaystack err:%v\n", err.Error())
			} else {
				log.Println("Transfer response from HandleEventPayouts", res)
			}
			HandlePayoutRes(ctx, server, res, resPayouts)
		} else if i%100 == 0 {
			bulkTransfer, resPayouts := HandleBulkTransfer(ctx, server, transferData, "main_payout")
			res, err := payout.TransferByPaystack(ctx, paystackKey, bulkTransfer)
			if err != nil {
				log.Printf("Error at HandleEventPayouts in TransferByPaystack err:%v\n", err.Error())
			} else {
				log.Println("Transfer response from HandleEventPayouts", res)
			}
			HandlePayoutRes(ctx, server, res, resPayouts)
			duration := 7 * time.Second
			time.Sleep(duration)
			transferData = make(map[uuid.UUID]payout.PayoutData)
		}
		i++
	}

}

func DailyHandlePayouts(ctx context.Context, server *Server) func() {
	return func() {
		// First we start by handling options
		HandleOptionPayouts(ctx, server)
		// Next we handle events payouts
		HandleEventPayouts(ctx, server)
	}
}

func DailyHandleStatusPayouts(ctx context.Context, server *Server) func() {
	return func() {
		paystackKey := server.config.PaystackSecretLiveKey
		payouts, err := server.store.ListPayout(ctx)
		if err != nil {
			log.Printf("Error at DailyHandleStatusPayouts ListPayout err:%v\n", err.Error())
			return
		}
		for _, p := range payouts {
			transferData, err := payout.HandlePaystackVerifyPayout(ctx, paystackKey, tools.UuidToString(p.ID))
			if err != nil {
				log.Printf("Error at DailyHandleStatusPayouts HandlePaystackVerifyPayout err:%v\n", err.Error())
				continue
			}
			if transferData.Data.Status == "success" {
				timePaid, err := tools.ConvertStringToTime(transferData.Data.TransferredAt)
				if err != nil {
					timePaid = time.Now().UTC()
					err = nil
				}
				payoutData, err := server.store.UpdatePayout(ctx, db.UpdatePayoutParams{
					AmountPayed: pgtype.Int8{
						Int64: int64(transferData.Data.Amount),
						Valid: true,
					},
					IsComplete: pgtype.Bool{
						Bool:  true,
						Valid: true,
					},
					TimePaid: pgtype.Timestamptz{
						Time:  timePaid,
						Valid: true,
					},
					ID: p.ID,
				})
				if err != nil {
					log.Printf("Error at DailyHandleStatusPayouts UpdatePayout err:%v\n", err.Error())
					continue
				}
				// We want to update the main_payout
				for _, chargeID := range payoutData.PayoutIds {
					err = server.store.UpdateMainPayout(ctx, db.UpdateMainPayoutParams{
						Status: pgtype.Text{
							String: "completed",
							Valid:  true,
						},
						TimePaid: pgtype.Timestamptz{
							Time:  timePaid,
							Valid: true,
						},
						ChargeID: chargeID,
					})
					if err != nil {
						log.Printf("Error at DailyHandleStatusPayouts UpdateMainPayout err:%v\n", err.Error())
						continue
					}
				}
			}
		}
	}
}

//func HandleTransferSuccess(ctx context.Context, server *Server, data map[string]string, reference string, funcName string, redisIDs []string) {
//	split := strings.Split(reference, "&")
//	if len(split) != 2 {
//		return
//	}
//	mainReference, err := tools.StringToUuid(split[1])
//	if err != nil {
//		log.Printf("Error at DailyHandleWebhookData in tools.StringToUuid(r) err:%v, reference: %v\n", err.Error(), reference)
//		return
//	}
//	payout, err := server.store.UpdatePayout(ctx, db.UpdatePayoutParams{
//		TransferCode: pgtype.Text{
//			String: data[constants.TRANSFER_CODE],
//			Valid:  true,
//		},
//		AmountPayed: pgtype.Int8{
//			Int64: tools.MoneyStringToInt(data[constants.AMOUNT]),
//			Valid: true,
//		},
//		IsComplete: pgtype.Bool{
//			Bool:  true,
//			Valid: true,
//		},
//		TimePaid: pgtype.Timestamptz{
//			Time:  time.Now(),
//			Valid: true,
//		},
//		ID: mainReference,
//	})
//	if err != nil {
//		log.Printf("Error at funcName: %v, HandleTransferSuccess in .UpdatePayout err:%v, reference: %v\n", funcName, err.Error(), reference)
//	}
//	// We want to update main_payout chargeIDs
//	switch payout.ParentType {
//	case "main_payout":
//		for _, chargeID := range payout.PayoutIds {
//			err = server.store.UpdateMainPayout(ctx, db.UpdateMainPayoutParams{
//				IsComplete:    true,
//				ChargeID:      chargeID,
//				AccountNumber: payout.AccountNumber,
//				TimePaid:      payout.TimePaid,
//			})
//			if err != nil {
//				log.Printf("Error at funcName: %v, HandleTransferSuccess in UpdateMainPayout err:%v, reference: %v\n", funcName, err.Error(), reference)
//			}
//			// Get Redis ids to remove
//			rCID, exist := GetRedisChargeID(redisIDs, chargeID)
//			if exist {
//				err := RedisClient.SRem(RedisContext, constants.PAYOUT_CHARGE_DATE_IDS, rCID).Err()
//				if err != nil {
//					log.Println("This redis  HandleTransferSuccess items were not removed from redis even though payout was unsuccessful", chargeID)
//				}
//			}

//		}
//	case "refund_payout":
//		for _, chargeID := range payout.PayoutIds {
//			err = server.store.UpdateRefundPayout(ctx, db.UpdateRefundPayoutParams{
//				IsComplete:    true,
//				ChargeID:      chargeID,
//				AccountNumber: payout.AccountNumber,
//				TimePaid:      payout.TimePaid,
//			})
//			if err != nil {
//				log.Printf("Error at funcName: %v, HandleTransferSuccess in UpdateRefundPayout err:%v, reference: %v\n", funcName, err.Error(), reference)
//			}
//			// Get Redis ids to remove
//			rCID, exist := GetRedisChargeID(redisIDs, chargeID)
//			if exist {
//				err := RedisClient.SRem(RedisContext, constants.PAYOUT_CHARGE_DATE_IDS, rCID).Err()
//				if err != nil {
//					log.Println("This redis  HandleTransferSuccess items were not removed from redis even though payout was unsuccessful", chargeID)
//				}
//			}

//		}
//	}
//	err = RedisClient.Del(RedisContext, reference).Err()
//	if err != nil {
//		log.Printf("Error at funcName: %v, HandleTransferSuccess in RedisClient.Del err:%v, reference: %v\n", funcName, err.Error(), reference)
//	}
//	err = RedisClient.SRem(RedisContext, constants.WEBHOOK_PAYSTACK_TRANSFER_REFERENCE, reference).Err()
//	if err != nil {
//		log.Printf("Error at funcName: %v, HandleTransferSuccess in RedisClient.SRem err:%v, reference: %v\n", funcName, err.Error(), reference)
//	}

//}

//func HandleTransferNotSuccess(ctx context.Context, server *Server, data map[string]string, reference string, funcName string, redisIDs []string) {
//	split := strings.Split(reference, "&")
//	if len(split) != 2 {
//		return
//	}
//	mainReference, err := tools.StringToUuid(split[1])
//	if err != nil {
//		log.Printf("Error at DailyHandleWebhookData in tools.StringToUuid(r) err:%v, reference: %v\n", err.Error(), reference)
//		return
//	}
//	// We only update the TransferCode because a transfer was still made even though is was not successful
//	payout, err := server.store.UpdatePayout(ctx, db.UpdatePayoutParams{
//		TransferCode: pgtype.Text{
//			String: data[constants.TRANSFER_CODE],
//			Valid:  true,
//		},
//		ID: mainReference,
//	})
//	if err != nil {
//		log.Printf("Error at funcName: %v, HandleTransferNotSuccess in .UpdatePayout err:%v, mainReference: %v\n", funcName, err.Error(), mainReference)
//	}
//	// We want to update main_payout chargeIDs
//	for _, chargeID := range payout.PayoutIds {
//		// Get Redis ids to remove
//		rCID, exist := GetRedisChargeID(redisIDs, chargeID)
//		if exist {
//			err := RedisClient.SRem(RedisContext, constants.PAYOUT_CHARGE_DATE_IDS, rCID).Err()
//			if err != nil {
//				log.Println("This redis  HandleTransferSuccess items were not removed from redis even though payout was unsuccessful", chargeID)
//			}
//		}
//	}
//	err = RedisClient.Del(RedisContext, reference).Err()
//	if err != nil {
//		log.Printf("Error at funcName: %v, HandleTransferNotSuccess in RedisClient.Del err:%v, reference: %v\n", funcName, err.Error(), reference)
//	}
//	err = RedisClient.SRem(RedisContext, constants.WEBHOOK_PAYSTACK_TRANSFER_REFERENCE, reference).Err()
//	if err != nil {
//		log.Printf("Error at funcName: %v, HandleTransferNotSuccess in RedisClient.SRem err:%v, reference: %v\n", funcName, err.Error(), reference)
//	}

//}
