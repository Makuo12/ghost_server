package api

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/makuo12/ghost_server/constants"
	db "github.com/makuo12/ghost_server/db/sqlc"
	"github.com/makuo12/ghost_server/payment"
	"github.com/makuo12/ghost_server/payout"
	"github.com/makuo12/ghost_server/tools"
	"github.com/makuo12/ghost_server/utils"

	"github.com/google/uuid"
)

func HandleInitRefund(ctx context.Context, server *Server, user db.User, referenceString string, paymentReference string, objectReference uuid.UUID, hasObjectReference bool, mainOption string, reason string, amount string, currency string, funcName string) {
	charge, err := CreateChargeReference(ctx, server, user.UserID, referenceString, paymentReference, objectReference, hasObjectReference, reason, currency, mainOption, amount, "HandleInitRefund")
	if err != nil {
		log.Printf("Error at FuncName: %v HandleInitRefund in CreateChargeReference: %v, reference: %v, userID: %v \n", funcName, err.Error(), referenceString, user.ID)
		return
	}
	_, err = server.store.CreateMainRefund(ctx, db.CreateMainRefundParams{
		ChargeID:    charge.ID,
		UserPercent: 100,
		HostPercent: 0,
		ChargeType:  constants.CHARGE_REFERENCE,
		Type:        reason,
	})
	if err != nil {
		log.Printf("Error at FuncName: %v HandleInitRefund in CreateMainRefund: %v, reference: %v, userID: %v \n", funcName, err.Error(), referenceString, user.ID)
		return
	}
}

func DailyHandleSendRefund(ctx context.Context, server *Server) func() {
	return func() {
		HandleOptionMainRefund(ctx, server)
		HandleTicketMainRefund(ctx, server)
		HandleAddCardRefund(ctx, server)
	}
}

func HandleAddCardRefund(ctx context.Context, server *Server) {
	refunds, err := server.store.ListMainRefundWithCharge(ctx, false)
	if err != nil || len(refunds) == 0 {
		if err != nil {
			log.Printf("HandleAddCardRefund in ListOptionMainRefundWithCharge err:%v\n", err.Error())
		}
	} else {
		url := "https://api.paystack.co/refund"
		bearer := "Bearer " + server.config.PaystackSecretLiveKey
		for _, r := range refunds {
			if r.Status == "not_started" {
				err := HandleMainRefund(ctx, server, r.ChargeID, uuid.New(), r.UID, "0", "0", r.Reference, url, bearer, int(r.UserPercent), int(r.HostPercent), r.Currency)
				if err != nil {
					log.Printf("HandleAddCardRefund in HandleMainRefund err:%v\n", err.Error())
				}
			}

		}
	}

}

func HandleOptionMainRefund(ctx context.Context, server *Server) {
	refunds, err := server.store.ListOptionMainRefundWithCharge(ctx, false)
	if err != nil || len(refunds) == 0 {
		if err != nil {
			log.Printf("HandleOptionMainRefund in ListOptionMainRefundWithCharge err:%v\n", err.Error())
		}
	} else {
		url := "https://api.paystack.co/refund"
		bearer := "Bearer " + server.config.PaystackSecretLiveKey
		for _, r := range refunds {
			if r.Status == "not_started" {
				err := HandleMainRefund(ctx, server, r.ChargeID, r.HostID, r.UID, tools.IntToMoneyString(r.ServiceFee), tools.IntToMoneyString(r.TotalFee), r.PaymentReference, url, bearer, int(r.UserPercent), int(r.HostPercent), r.Currency)
				if err != nil {
					log.Printf("HandleOptionMainRefund in HandleMainRefund err:%v\n", err.Error())
				}
			}

		}
	}
}

func HandleTicketMainRefund(ctx context.Context, server *Server) {
	refunds, err := server.store.ListTicketMainRefundWithCharge(ctx, false)
	if err != nil || len(refunds) == 0 {
		if err != nil {
			log.Printf("HandleTicketMainRefund in ListTicketMainRefundWithCharge err:%v\n", err.Error())
		}
	} else {
		url := "https://api.paystack.co/refund"
		bearer := "Bearer " + server.config.PaystackSecretLiveKey
		for _, r := range refunds {
			if r.Status == "not_started" {
				err := HandleMainRefund(ctx, server, r.ChargeID, r.HostID, r.UID, tools.IntToMoneyString(r.ServiceFee), tools.IntToMoneyString(r.TotalFee), r.PaymentReference, url, bearer, int(r.UserPercent), int(r.HostPercent), r.Currency)
				if err != nil {
					log.Printf("HandleTicketMainRefund in HandleMainRefund err:%v\n", err.Error())
				}
			}

		}
	}
}

func HandleMainRefund(ctx context.Context, server *Server, chargeID uuid.UUID, hostID uuid.UUID, uID uuid.UUID, serviceFee string, totalFee string, paymentReference string, url string, bearer string, userPercent int, hostPercent int, currency string) error {
	if userPercent == 0 {
		HandleRefundPayoutDB(ctx, server, chargeID, hostID, int(hostPercent), serviceFee, totalFee, "HandleMainRefund", currency)
		return nil
	}
	resData, err := payment.HandlePaystackCreateRefund(paymentReference, url, userPercent, totalFee, bearer)
	if err != nil {
		return err
	}
	// Has we have made the refund request, we now update the main refund
	_, err = server.store.UpdateMainRefund(ctx, db.UpdateMainRefundParams{
		ChargeID: chargeID,
		Status: pgtype.Text{
			String: "processing",
			Valid:  true,
		},
	})
	if err != nil {
		header := "Error at HandleMainRefund for UpdateMainRefund"
		message := fmt.Sprintf("Error occurred at HandleMainRefund in UpdateMainRefund err: %v, for chargeID: %v", err.Error(), chargeID)
		log.Printf("error at HandleMainRefund in json.NewDecoder %v \n", err.Error())
		BrevoErrorMessage(ctx, server, header, message, "HandleMainRefund")
	}
	refundID := fmt.Sprint("%v", resData.Data.ID)
	HandleRefundDB(ctx, server, chargeID, uID, "paystack", resData.Data.Amount, paymentReference, "HandleMainRefund", refundID)
	HandleRefundPayoutDB(ctx, server, chargeID, hostID, int(hostPercent), serviceFee, totalFee, "HandleMainRefund", currency)
	return nil

}

func HandleRefundDB(ctx context.Context, server *Server, chargeID uuid.UUID, uID uuid.UUID, sendMedium string, amount int, paymentReference string, funcName string, refundID string) {
	_, err := server.store.CreateRefund(ctx, db.CreateRefundParams{
		ChargeID:   chargeID,
		Reference:  paymentReference,
		SendMedium: sendMedium,
		UserID:     uID,
		Amount:     int64(amount),
		RefundID:   refundID,
	})
	if err != nil {
		log.Printf("FuncName: %v. There an error at HandleRefundDB at server.store.CreateRefund: %v, chargeID: %v \n", funcName, err.Error(), chargeID)
	}
}

func HandleRefundPayoutDB(ctx context.Context, server *Server, chargeID uuid.UUID, hostID uuid.UUID, hostPercent int, serviceFee string, totalFee string, funcName string, currency string) {
	if hostPercent != 0 {
		dollarToNaira := server.config.DollarToNaira
		dollarToCAD := server.config.DollarToCAD
		var payoutAmount float64
		var servicePercent float64
		var payoutServiceFee float64
		// We need to get the host percent. How much the host is going to take
		amount := float64(hostPercent/100) * tools.ConvertStringToFloat(totalFee)
		amount = amount - tools.ConvertStringToFloat(serviceFee)
		amount, err := tools.ConvertPrice(tools.ConvertFloatToString(amount), currency, utils.PayoutCurrency, dollarToNaira, dollarToCAD, hostID)
		if err != nil {
			log.Printf("FuncName: %v. There an error at HandleRefundPayoutDB at tools.ConvertPrice: %v, chargeID: %v \n", funcName, err.Error(), chargeID)
			err = nil
			return
		}
		switch utils.PayoutCurrency {
		case utils.NGN:
			servicePercent = tools.ConvertStringToFloat(server.config.LocServiceEventHostPercent)
		default:
			servicePercent = tools.ConvertStringToFloat(server.config.IntServiceEventHostPercent)
		}
		payoutServiceFee = (servicePercent / 100) * amount
		payoutAmount = amount - payoutServiceFee
		_, err = server.store.CreateRefundPayout(ctx, db.CreateRefundPayoutParams{
			ChargeID:   chargeID,
			Amount:     tools.MoneyFloatToInt(payoutAmount),
			UserID:     hostID,
			Currency:   utils.PayoutCurrency,
			ServiceFee: tools.MoneyStringToInt(serviceFee),
		})
		if err != nil {
			log.Printf("FuncName: %v. There an error at ProcessEventHostCancelUniqueID at server.store.CreateRefundPayout: %v, chargeID: %v \n", funcName, err.Error(), chargeID)
		}
	}
}

func HandleRefundSuccess(ctx context.Context, server *Server, funcName string, amountPaid int, paymentReference string, timePaid time.Time) {
	// Let update it directly to the database
	refundData, err := server.store.UpdateRefund(ctx, db.UpdateRefundParams{
		IsComplete: pgtype.Bool{
			Bool:  true,
			Valid: true,
		},
		TimePaid: pgtype.Timestamptz{
			Time:  timePaid,
			Valid: true,
		},
		AmountPayed: pgtype.Int8{
			Int64: int64(amountPaid),
			Valid: true,
		},
		Reference: paymentReference,
	})
	if err != nil {
		header := "Error at HandleRefundWebhook"
		msg := fmt.Sprintf("Error at HandleRefundWebhook at UpdateRefund err: %v", err)
		log.Printf("Error at HandleRefundWebhook in SMembers err:%v\n", err.Error())
		BrevoErrorMessage(ctx, server, header, msg, "HandleRefundWebhook")
	}

	_, err = server.store.UpdateMainRefund(ctx, db.UpdateMainRefundParams{
		HasPaid: pgtype.Bool{
			Valid: true,
			Bool:  true,
		},
		Status: pgtype.Text{
			Valid:  true,
			String: "completed",
		},
		ChargeID: refundData.ChargeID,
	})

}

func HandleDailyRefundUpdate(ctx context.Context, server *Server) {
	dbRefunds, err := server.store.ListMainRefundProcessing(ctx)
	if err != nil || len(dbRefunds) == 0 {
		if err != nil {
			header := "Error at HandleDailyRefundUpdate"
			msg := fmt.Sprintf("Error at HandleDailyRefundUpdate at ListMainRefundProcessing err: %v", err)
			log.Printf("Error at HandleDailyRefundUpdate in SMembers err:%v\n", err.Error())
			BrevoErrorMessage(ctx, server, header, msg, "HandleDailyRefundUpdate")
		}
		return
	}
	for _, dbRefund := range dbRefunds {
		paystackRefund, err := payout.HandlePaystackFetchRefund(ctx, server.config.PaystackSecretLiveKey, HandleSqlNullString(dbRefund.RefundID))
		if paystackRefund.Data.Status != "processed" {
			continue
		}
		if err != nil {
			log.Printf("Error at HandleDailyRefundUpdate in HandlePaystackFetchRefund err:%v\n", err.Error())
			continue
		}
		timePaid, err := tools.ConvertStringToTime(paystackRefund.Data.RefundedAt)
		if err != nil {
			timePaid = time.Now().UTC()
			err = nil
		}
		HandleRefundSuccess(ctx, server, "HandleDailyRefundUpdate", paystackRefund.Data.Amount, paystackRefund.Data.TransactionReference, timePaid)
	}

}


