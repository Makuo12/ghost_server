package api

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/makuo12/ghost_server/db/sqlc"
	"github.com/makuo12/ghost_server/payout"
	"github.com/makuo12/ghost_server/tools"
	"github.com/makuo12/ghost_server/utils"

	"github.com/google/uuid"
)

func HandleRefundPayoutData(server *Server, dollarToNaira string, dollarToCAD string, data db.ListRefundPayoutWithUserRow, payoutData payout.PayoutData) (res payout.PayoutData, err error) {
	amount, err := tools.ConvertPrice(tools.IntToMoneyString(data.Amount), data.Currency, utils.NGN, dollarToNaira, dollarToCAD, data.HostID)
	if err != nil {
		log.Printf("HandleRefundPayoutData tools.ConvertPrice err:%v\n", err.Error())
		return
	}
	payoutIDs := append(payoutData.ChargeIDs, data.ChargeID)
	newAmounts := append(payoutData.Amounts, amount)
	res = payout.PayoutData{
		ChargeIDs:            payoutIDs,
		HostID:               data.HostID,
		HostDefaultAccountID: data.HostDefaultAccountID,
		HostUserID:           data.HostUserID,
		Amounts:              newAmounts,
		HostName:             data.HostName,
	}
	return
}

func GetRefundPayout(ctx context.Context, server *Server, dollarToNaira string, dollarToCAD string) (map[uuid.UUID]payout.PayoutData, error) {
	var payouts = make(map[uuid.UUID]payout.PayoutData)
	refundPayouts, err := server.store.ListRefundPayoutWithUser(ctx, db.ListRefundPayoutWithUserParams{
		IsComplete: false,
		Status:     "not_started",
	})
	if err != nil || len(refundPayouts) == 0 {
		if err != nil {
			log.Printf("GetRefundPayout in ListRefundPayoutWithUser err:%v\n", err.Error())
		}
	} else {
		for _, refundPayout := range refundPayouts {
			data, err := HandleRefundPayoutData(server, dollarToNaira, dollarToCAD, refundPayout, payouts[refundPayout.HostID])
			if err != nil {
				log.Printf("GetRefundPayout in HandleRefundPayoutData err:%v\n", err.Error())
			} else {
				payouts[refundPayout.HostID] = data
			}
		}
	}
	return payouts, nil
}

func HandleRefundPayouts(ctx context.Context, server *Server) {
	dollarToNaira := server.config.DollarToNaira
	dollarToCAD := server.config.DollarToCAD
	paystackKey := server.config.PaystackSecretLiveKey
	payouts, err := GetRefundPayout(ctx, server, dollarToNaira, dollarToCAD)
	if err != nil || len(payouts) == 0 {
		if err != nil {
			log.Printf("Error at HandleRefundPayouts in GetRefundPayout err:%v\n", err.Error())
		}
		return
	}
	var transferData = make(map[uuid.UUID]payout.PayoutData)
	var i = 1
	for key, value := range payouts {
		transferData[key] = value
		i++
		if i == len(payouts) {
			bulkTransfer, resPayouts := HandleBulkTransfer(ctx, server, transferData, "refund_payout")
			res, err := payout.TransferByPaystack(ctx, paystackKey, bulkTransfer)
			if err != nil {
				log.Printf("Error at HandleRefundPayouts in TransferByPaystack err:%v\n", err.Error())
			} else {
				log.Println("Transfer response from HandleRefundPayouts", res)
			}
			HandleRefundPayoutRes(ctx, server, res, resPayouts)
		} else if i%100 == 0 {
			bulkTransfer, resPayouts := HandleBulkTransfer(ctx, server, transferData, "refund_payout")
			res, err := payout.TransferByPaystack(ctx, paystackKey, bulkTransfer)
			if err != nil {
				log.Printf("Error at HandleRefundPayouts in TransferByPaystack err:%v\n", err.Error())
			} else {
				log.Println("Transfer response from HandleRefundPayouts", res)
			}
			HandleRefundPayoutRes(ctx, server, res, resPayouts)
			duration := 7 * time.Second
			time.Sleep(duration)
			transferData = make(map[uuid.UUID]payout.PayoutData)
		}
	}

}

func DailyHandleRefundPayouts(ctx context.Context, server *Server) func() {
	return func() {
		HandleRefundPayouts(ctx, server)
	}
}

func DailyHandleStatusRefundPayouts(ctx context.Context, server *Server) func() {
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
					err = server.store.UpdateRefundPayout(ctx, db.UpdateRefundPayoutParams{
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
						log.Printf("Error at DailyHandleStatusPayouts UpdateRefundPayout err:%v\n", err.Error())
						continue
					}
				}
			}
		}
	}
}
