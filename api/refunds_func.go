package api

import (
	"context"
	"flex_server/constants"
	db "flex_server/db/sqlc"
	"flex_server/tools"
	"flex_server/utils"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)



func HandleRefundPayoutData(server *Server, dollarToNaira string, dollarToCAD string, data db.ListRefundPayoutWithUserRow, payout PayoutData) (res PayoutData, err error) {
	amount, err := tools.ConvertPrice(tools.IntToMoneyString(data.Amount), data.Currency, utils.NGN, dollarToNaira, dollarToCAD, data.HostID)
	if err != nil {
		log.Printf("HandleRefundPayoutData tools.ConvertPrice err:%v\n", err.Error())
		return
	}
	payoutIDs := append(payout.ChargeIDs, data.ChargeID)
	newAmounts := append(payout.Amounts, amount)
	res = PayoutData{
		ChargeIDs:            payoutIDs,
		HostID:               data.HostID,
		HostDefaultAccountID: data.HostDefaultAccountID,
		HostUserID:           data.HostUserID,
		Amounts:              newAmounts,
		HostName:             data.HostName,
	}
	return
}

func GetRefundPayout(ctx context.Context, server *Server, dollarToNaira string, dollarToCAD string, redisDateIDs []string) (map[uuid.UUID]PayoutData, error) {
	var payouts = make(map[uuid.UUID]PayoutData)
	refundPayouts, err := server.store.ListRefundPayoutWithUser(ctx, false)
	if err != nil || len(refundPayouts) == 0 {
		if err != nil {
			log.Printf("GetRefundPayout in ListRefundPayoutWithUser err:%v\n", err.Error())
		}
	} else {
		for v := 0; v < len(refundPayouts); v++ {
			if ChargeIDInPay(redisDateIDs, refundPayouts[v].ChargeID) {
				continue
			}
			data, err := HandleRefundPayoutData(server, dollarToNaira, dollarToCAD, refundPayouts[v], payouts[refundPayouts[v].HostID])
			if err != nil {
				log.Printf("GetRefundPayout in HandleRefundPayoutData err:%v\n", err.Error())
			} else {
				chargeID := refundPayouts[v].ChargeID
				date := time.Now()
				date_id := fmt.Sprintf("%v&%v", date, chargeID)
				// We store it in redis
				err := RedisClient.SAdd(RedisContext, constants.PAYOUT_CHARGE_DATE_IDS, date_id).Err()
				if err != nil {
					log.Printf("GetRefundPayout in RedisClient.SAdd err:%v\n", err.Error())
				}
				payouts[refundPayouts[v].HostID] = data
			}
		}
	}
	return payouts, nil
}


func HandleRefundPayouts(ctx context.Context, server *Server, redisDateIDs []string, refundIDs []string) {
	dollarToNaira := server.config.DollarToNaira
	dollarToCAD := server.config.DollarToCAD
	payouts, err := GetRefundPayout(ctx, server, dollarToNaira, dollarToCAD, redisDateIDs)
	if err != nil || len(payouts) == 0 {
		if err != nil {
			log.Printf("Error at HandleRefundPayouts in GetRefundPayout err:%v\n", err.Error())
		}
		return
	}
	var transferData = make(map[uuid.UUID]PayoutData)
	var i = 1
	for key, value := range payouts {
		transferData[key] = value
		i++
		if i == len(payouts) {
			bulkTransfer, resPayouts := HandleBulkTransfer(ctx, server, transferData, "refund_payout")
			res, err := TransferByPaystack(ctx, server, bulkTransfer)
			if err != nil {
				log.Printf("Error at HandleRefundPayouts in TransferByPaystack err:%v\n", err.Error())
			} else {
				log.Println("Transfer response from HandleRefundPayouts", res)
			}
			HandlePayoutRes(res, redisDateIDs, resPayouts)
		} else if i%100 == 0 {
			bulkTransfer, resPayouts := HandleBulkTransfer(ctx, server, transferData, "refund_payout")
			res, err := TransferByPaystack(ctx, server, bulkTransfer)
			if err != nil {
				log.Printf("Error at HandleRefundPayouts in TransferByPaystack err:%v\n", err.Error())
			} else {
				log.Println("Transfer response from HandleRefundPayouts", res)
			}
			HandlePayoutRes(res, redisDateIDs, resPayouts)
			duration := 7 * time.Second
			time.Sleep(duration)
			transferData = make(map[uuid.UUID]PayoutData)
		}
	}

}

func DailyHandleRefundPayouts(ctx context.Context, server *Server) func() {
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
		// First we start by handling refund payouts
		HandleRefundPayouts(ctx, server, redisChargeDateIDs, []string{})
	}
}