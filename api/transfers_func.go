package api

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/makuo12/ghost_server/constants"
	db "github.com/makuo12/ghost_server/db/sqlc"
	"github.com/makuo12/ghost_server/payout"
	"github.com/makuo12/ghost_server/tools"
	"github.com/makuo12/ghost_server/utils"
)

func TransferAmount(amounts []float64) float64 {
	var totalAmount float64
	for _, a := range amounts {
		totalAmount += a
	}
	return math.Floor(totalAmount)
}

func TransferRecipient(ctx context.Context, server *Server, hostID uuid.UUID, hostUserID uuid.UUID, hostDefaultAccountID string, hostName string) (accountNum string, recipient string, err error) {
	if !tools.ServerStringEmpty(hostDefaultAccountID) {
		accountID, err := tools.StringToUuid(hostDefaultAccountID)
		if err != nil {

			log.Printf("TransferRecipient tools.StringToUuid err:%v hostDefaultAccountID:%v\n", err.Error(), hostDefaultAccountID)
			msg := fmt.Sprintf("Hey %v,\nPlease try updating your default payout method as it is not valid.", hostName)
			CreateTypeNotification(ctx, server, uuid.New(), hostUserID, constants.INVALID_DEFAULT_ACCOUNT_NUMBER, msg, true, "Invalid payout method")
		} else {
			account, err := server.store.GetDefaultAccountNumber(ctx, db.GetDefaultAccountNumberParams{
				UserID: hostID,
				ID:     accountID,
			})
			if err != nil || tools.ServerStringEmpty(account.AccountNumber) || tools.ServerStringEmpty(account.RecipientCode) {

				if err != nil {
					log.Printf("TransferRecipient in .GetDefaultAccountNumber err:%v accountID:%v\n", err.Error(), accountID)
				}
				msg := fmt.Sprintf("Hey %v, Please try updating your default payout method as it is not valid.", hostName)
				CreateTypeNotification(ctx, server, accountID, hostUserID, constants.INVALID_DEFAULT_ACCOUNT_NUMBER, msg, false, "Invalid payout method")
				err = fmt.Errorf("account number not found")
			}
			log.Printf("data in account number %v, %v ", account.RecipientCode, account.AccountNumber)
			recipient = account.RecipientCode
			accountNum = account.AccountNumber
			log.Printf("recipient %v\n", recipient)
			log.Printf("accountNum %v\n", accountNum)
		}
	} else {
		err = errors.New("no valid account number")
		msg := fmt.Sprintf("Hey %v, Please try adding a default payout method as we could not make payout to your account.", hostName)
		CreateTypeNotification(ctx, server, uuid.New(), hostUserID, constants.NO_DEFAULT_ACCOUNT_NUMBER, msg, true, "No payout method")
		// Need to send an email
	}
	return
}

func HandleOptionPayoutData(server *Server, dollarToNaira string, dollarToCAD string, data db.ListOptionMainPayoutWithChargeRow, payoutData payout.PayoutData) (res payout.PayoutData, err error) {
	payoutIDs := append(payoutData.ChargeIDs, data.ChargeID)
	newAmounts := append(payoutData.Amounts, tools.ConvertStringToFloat(tools.IntToMoneyString(data.Amount)))
	res = payout.PayoutData{
		ChargeIDs:            payoutIDs,
		HostID:               data.HostID,
		HostDefaultAccountID: data.DefaultAccountID,
		HostUserID:           data.HostUserID,
		Amounts:              newAmounts,
		HostName:             data.FirstName,
	}
	return
}

func HandleEventPayoutData(server *Server, dollarToNaira string, dollarToCAD string, data db.ListTicketMainPayoutWithChargeRow, payoutData payout.PayoutData) (res payout.PayoutData) {
	payoutIDs := append(payoutData.ChargeIDs, data.ChargeID)
	newAmounts := append(payoutData.Amounts, tools.ConvertStringToFloat(tools.IntToMoneyString(data.Amount)))
	res = payout.PayoutData{
		ChargeIDs:            payoutIDs,
		HostID:               data.HostID,
		HostDefaultAccountID: data.DefaultAccountID,
		HostUserID:           data.HostUserID,
		Amounts:              newAmounts,
		HostName:             data.FirstName,
	}
	return
}

func GetEventPayout(ctx context.Context, server *Server, dollarToNaira string, dollarToCAD string) (map[uuid.UUID]payout.PayoutData, error) {
	var payouts = make(map[uuid.UUID]payout.PayoutData)
	// Lets handle Event payouts
	eventCharges, err := server.store.ListTicketMainPayoutWithCharge(ctx, db.ListTicketMainPayoutWithChargeParams{
		PayoutComplete:        false,
		ChargePaymentComplete: true,
		ChargeCancelled:       false,
		PayoutStatus:          "not_started",
	})
	if err != nil || len(eventCharges) == 0 {
		if err != nil {
			log.Printf("DailyHandlePayouts in ListEventMainPayoutWithCharge err:%v\n", err.Error())
		}
	} else {
		// If no error we start making payouts
		for v := 0; v < len(eventCharges); v++ {
			data := HandleEventPayoutData(server, dollarToNaira, dollarToCAD, eventCharges[v], payouts[eventCharges[v].HostID])
			payouts[eventCharges[v].HostID] = data
		}
	}
	return payouts, nil
}

func GetOptionPayout(ctx context.Context, server *Server, dollarToNaira string, dollarToCAD string) (map[uuid.UUID]payout.PayoutData, error) {
	var payouts = make(map[uuid.UUID]payout.PayoutData)
	// Lets handle Option payouts
	optionCharges, err := server.store.ListOptionMainPayoutWithCharge(ctx, db.ListOptionMainPayoutWithChargeParams{
		PayoutComplete:        false,
		ChargePaymentComplete: true,
		ChargeCancelled:       false,
		PayoutStatus:          "not_started",
	})
	if err != nil || len(optionCharges) == 0 {
		if err != nil {
			log.Printf("GetOptionPayout in ListOptionMainPayoutWithCharge err:%v\n", err.Error())
		}
	} else {
		// If no error we start making payouts
		for v := 0; v < len(optionCharges); v++ {
			// for each user we have a payout.PayoutData
			data, err := HandleOptionPayoutData(server, dollarToNaira, dollarToCAD, optionCharges[v], payouts[optionCharges[v].HostID])
			if err != nil {
				log.Printf("GetOptionPayout in HandleOptionPayoutData err:%v\n", err.Error())
			} else {
				payouts[optionCharges[v].HostID] = data
			}
		}
	}
	return payouts, nil
}

func HandleBulkTransfer(ctx context.Context, server *Server, payouts map[uuid.UUID]payout.PayoutData, payoutType string) (payout.BulkTransferRequest, map[uuid.UUID]payout.PayoutData) {
	var transfers []payout.Transfer
	var resPayouts = make(map[uuid.UUID]payout.PayoutData)
	for k, v := range payouts {
		accountNum, recipient, err := TransferRecipient(ctx, server, v.HostID, v.HostUserID, v.HostDefaultAccountID, v.HostName)
		if err != nil {
			log.Printf("HandleBulkTransfer in TransferRecipient err:%v, hostID: %v\n", err.Error(), v.HostID)
			continue
		}
		amount := TransferAmount(v.Amounts)
		// We add from our money to pay for the charges
		if amount <= 5000.0 {
			amount += 10
		} else if amount >= 5001.0 && amount <= 50000 {
			amount += 25
		} else {
			amount += 50
		}
		paystackAmount := tools.ConvertToPaystackPayout(tools.ConvertFloatToString(amount))
		payoutData, err := server.store.CreatePayout(ctx, db.CreatePayoutParams{
			PayoutIds:     v.ChargeIDs,
			SendMedium:    "paystack",
			UserID:        v.HostID,
			Amount:        int64(paystackAmount),
			AmountPayed:   int64(paystackAmount),
			ParentType:    payoutType,
			AccountNumber: accountNum,
		})
		if err != nil {
			log.Printf("HandleBulkTransfer in TransferRecipient err:%v, hostID: %v\n", err.Error(), v.HostID)
			continue
		}
		reference := tools.UuidToString(payoutData.ID)

		reason := fmt.Sprintf("%v payout", v.HostName)
		data := payout.Transfer{
			Amount:    paystackAmount,
			Reference: reference,
			Reason:    reason,
			Recipient: recipient,
		}
		transfers = append(transfers, data)
		newData := payout.PayoutData{ChargeIDs: v.ChargeIDs, HostID: v.HostID, HostUserID: v.HostUserID, HostDefaultAccountID: v.HostDefaultAccountID, HostName: v.HostName, Amounts: v.Amounts, PayoutID: reference}
		resPayouts[k] = newData
	}
	bulkTransfer := payout.BulkTransferRequest{
		Currency:  utils.NGN,
		Source:    "balance",
		Transfers: transfers,
	}
	return bulkTransfer, resPayouts
}

// HandlePayoutRes this function remove redis date_ids that payout was not successful
func HandlePayoutRes(ctx context.Context, server *Server, res payout.TransferQueueResponse, resPayouts map[uuid.UUID]payout.PayoutData) {
	// Setup the transfer codes
	for _, transfer := range res.Data {
		id, err := tools.StringToUuid(transfer.Reference)
		if err != nil {
			log.Printf("error at HandlePayoutRes in StringToUuid %v \n", err.Error())
			continue
		}
		_, err = server.store.UpdatePayout(ctx, db.UpdatePayoutParams{
			TransferCode: pgtype.Text{
				String: transfer.TransferCode,
				Valid:  true,
			},
			ID: id,
		})
		if err != nil {
			log.Printf("error at HandlePayoutRes in UpdatePayout %v \n", err.Error())
			continue
		}
	}
	for _, value := range resPayouts {
		// We want to update all the chargeIDs to be processing
		for _, chargeID := range value.ChargeIDs {
			err := server.store.UpdateMainPayout(ctx, db.UpdateMainPayoutParams{
				Status: pgtype.Text{
					String: "processing",
					Valid:  true,
				},
				ChargeID: chargeID,
			})
			if err != nil {
				log.Printf("error at HandlePayoutRes in UpdateMainPayout %v \n", err.Error())
			}
		}
	}
}


func HandleRefundPayoutRes(ctx context.Context, server *Server, res payout.TransferQueueResponse, resPayouts map[uuid.UUID]payout.PayoutData) {
	// Setup the transfer codes
	for _, transfer := range res.Data {
		id, err := tools.StringToUuid(transfer.Reference)
		if err != nil {
			log.Printf("error at HandlePayoutRes in StringToUuid %v \n", err.Error())
			continue
		}
		_, err = server.store.UpdatePayout(ctx, db.UpdatePayoutParams{
			TransferCode: pgtype.Text{
				String: transfer.TransferCode,
				Valid:  true,
			},
			ID: id,
		})
		if err != nil {
			log.Printf("error at HandlePayoutRes in UpdatePayout %v \n", err.Error())
			continue
		}
	}
	for _, value := range resPayouts {
		// We want to update all the chargeIDs to be processing
		for _, chargeID := range value.ChargeIDs {
			err := server.store.UpdateRefundPayout(ctx, db.UpdateRefundPayoutParams{
				Status: pgtype.Text{
					String: "processing",
					Valid:  true,
				},
				ChargeID: chargeID,
			})
			if err != nil {
				log.Printf("error at HandlePayoutRes in UpdateRefundPayout %v \n", err.Error())
			}
		}
	}
}