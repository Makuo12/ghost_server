package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"flex_server/constants"
	db "flex_server/db/sqlc"
	"flex_server/tools"
	"flex_server/utils"
	"fmt"
	"log"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

func ChargeIDInPay(dateIDs []string, chargeID uuid.UUID) bool {
	cID := tools.UuidToString(chargeID)
	// If charge is in payout but 4 days have pasted we want to log it our to be aware

	for _, date_id := range dateIDs {
		dateIDsSplit := strings.Split(date_id, "&")
		if len(dateIDsSplit) == 2 {
			date, err := tools.ConvertStringToTime(dateIDsSplit[0])
			if err != nil {
				log.Printf("chargeIDInPayout tools.ConvertStringToTime err:%v\n", err.Error())
				continue
			}
			id := dateIDsSplit[1]
			if id == cID {
				if time.Now().After(date.Add(time.Hour * 96)) {
					// We log it out and return true
					log.Printf("for chargeID %v check if payment was made because it has passed for days and it is still in redis. redis id: %v\n", chargeID, date_id)
				}
				return true
			}
		} else {
			if date_id == cID {
				return true
			}
		}

	}
	return false
}

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
	}
	return
}

func HandleOptionPayoutData(server *Server, dollarToNaira string, dollarToCAD string, data db.ListOptionMainPayoutWithChargeRow, payout PayoutData) (res PayoutData, err error) {
	payoutIDs := append(payout.ChargeIDs, data.ChargeID)
	newAmounts := append(payout.Amounts, tools.ConvertStringToFloat(tools.IntToMoneyString(data.Amount)))
	res = PayoutData{
		ChargeIDs:            payoutIDs,
		HostID:               data.HostID,
		HostDefaultAccountID: data.DefaultAccountID,
		HostUserID:           data.HostUserID,
		Amounts:              newAmounts,
		HostName:             data.FirstName,
	}
	return
}

func HandleEventPayoutData(server *Server, dollarToNaira string, dollarToCAD string, data db.ListTicketMainPayoutWithChargeRow, payout PayoutData) (res PayoutData) {
	payoutIDs := append(payout.ChargeIDs, data.ChargeID)
	newAmounts := append(payout.Amounts, tools.ConvertStringToFloat(tools.IntToMoneyString(data.Amount)))
	res = PayoutData{
		ChargeIDs:            payoutIDs,
		HostID:               data.HostID,
		HostDefaultAccountID: data.DefaultAccountID,
		HostUserID:           data.HostUserID,
		Amounts:              newAmounts,
		HostName:             data.FirstName,
	}
	return
}

func GetEventPayout(ctx context.Context, server *Server, dollarToNaira string, dollarToCAD string, redisDateIDs []string) (map[uuid.UUID]PayoutData, error) {
	var payouts = make(map[uuid.UUID]PayoutData)
	// Lets handle Event payouts
	eventCharges, err := server.store.ListTicketMainPayoutWithCharge(ctx, db.ListTicketMainPayoutWithChargeParams{
		PayoutComplete:        false,
		ChargePaymentComplete: true,
		ChargeCancelled:       false,
	})
	if err != nil || len(eventCharges) == 0 {
		if err != nil {
			log.Printf("DailyHandlePayouts in ListEventMainPayoutWithCharge err:%v\n", err.Error())
		}
	} else {
		// If no error we start making payouts
		for v := 0; v < len(eventCharges); v++ {
			if ChargeIDInPay(redisDateIDs, eventCharges[v].ChargeID) {
				continue
			}
			data := HandleEventPayoutData(server, dollarToNaira, dollarToCAD, eventCharges[v], payouts[eventCharges[v].HostID])

			payouts[eventCharges[v].HostID] = data
		}
	}
	return payouts, nil
}

func GetOptionPayout(ctx context.Context, server *Server, dollarToNaira string, dollarToCAD string, redisDateIDs []string) (map[uuid.UUID]PayoutData, error) {
	var payouts = make(map[uuid.UUID]PayoutData)
	// Lets handle Option payouts
	optionCharges, err := server.store.ListOptionMainPayoutWithCharge(ctx, db.ListOptionMainPayoutWithChargeParams{
		PayoutComplete:        false,
		ChargePaymentComplete: true,
		ChargeCancelled:       false,
	})
	if err != nil || len(optionCharges) == 0 {
		if err != nil {
			log.Printf("GetOptionPayout in ListOptionMainPayoutWithCharge err:%v\n", err.Error())
		}
	} else {
		// If no error we start making payouts
		for v := 0; v < len(optionCharges); v++ {
			if ChargeIDInPay(redisDateIDs, optionCharges[v].ChargeID) {
				continue
			}
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

func AddPayoutChargeIDsToRedis(chargeIDs []uuid.UUID) (err error) {
	var data = []string{}
	for _, id := range chargeIDs {
		date := tools.ConvertTimeToString(time.Now())
		date_id := fmt.Sprintf("%v&%v", date, id)
		data = append(data, date_id)
	}
	// We store it in redis
	err = RedisClient.SAdd(RedisContext, constants.PAYOUT_CHARGE_DATE_IDS, data).Err()
	if err != nil {
		log.Printf("GetOptionPayout in RedisClient.SAdd err:%v\n", err.Error())
	}
	return
}
func RemovePayoutChargeIDsFromRedis(chargeIDs []uuid.UUID) {
	var data = []string{}
	for _, id := range chargeIDs {
		date := tools.ConvertTimeToString(time.Now())
		date_id := fmt.Sprintf("%v&%v", date, id)
		data = append(data, date_id)
	}
	// We store it in redis
	err := RedisClient.SRem(RedisContext, constants.PAYOUT_CHARGE_DATE_IDS, data).Err()
	if err != nil {
		log.Printf("GetOptionPayout in RedisClient.SAdd err:%v\n", err.Error())
	}
}

func HandleBulkTransfer(ctx context.Context, server *Server, payouts map[uuid.UUID]PayoutData, payoutType string) (BulkTransferRequest, map[uuid.UUID]PayoutData) {
	var transfers []Transfer
	var resPayouts = make(map[uuid.UUID]PayoutData)
	for k, v := range payouts {
		accountNum, recipient, err := TransferRecipient(ctx, server, v.HostID, v.HostUserID, v.HostDefaultAccountID, v.HostName)
		if err != nil {
			log.Printf("HandleBulkTransfer in TransferRecipient err:%v, hostID: %v\n", err.Error(), v.HostID)
			continue
		}
		err = AddPayoutChargeIDsToRedis(v.ChargeIDs)
		if err != nil {
			log.Printf("HandleBulkTransfer in AddPayoutChargeIDsToRedis err:%v, hostID: %v\n", err.Error(), v.HostID)
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
		payout, err := server.store.CreatePayout(ctx, db.CreatePayoutParams{
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
			RemovePayoutChargeIDsFromRedis(v.ChargeIDs)
			continue
		}
		reference := tools.UuidToString(payout.ID)

		reason := fmt.Sprintf("%v payout", v.HostName)
		data := Transfer{
			Amount:    paystackAmount,
			Reference: reference,
			Reason:    reason,
			Recipient: recipient,
		}
		transfers = append(transfers, data)
		newData := PayoutData{v.ChargeIDs, v.HostID, v.HostUserID, v.HostDefaultAccountID, v.HostName, v.Amounts, reference}
		resPayouts[k] = newData

	}
	bulkTransfer := BulkTransferRequest{
		Currency:  utils.NGN,
		Source:    "balance",
		Transfers: transfers,
	}
	return bulkTransfer, resPayouts
}

func TransferByPaystack(ctx context.Context, server *Server, bulkTransfer BulkTransferRequest) (resItem TransferQueueResponse, err error) {
	url := "https://api.paystack.co/transfer/bulk"
	var bearer = "Bearer " + server.config.PaystackSecretLiveKey
	var resData = &TransferQueueResponse{}
	buf := new(bytes.Buffer)
	err = json.NewEncoder(buf).Encode(bulkTransfer)
	if err != nil {
		return
	}
	request, err := http.NewRequest("POST", url, buf)
	if err != nil {
		log.Printf("error at %v in http.NewRequest %v \n", "TransferByPaystack", err.Error())
		err = fmt.Errorf("there was an internal server error while verifying your card. If this error continues contact help center")
		return
	}
	request.Header.Add("Authorization", bearer)
	// Send req using http Client
	clientSide := &http.Client{}
	res, err := clientSide.Do(request)
	if err != nil {
		log.Printf("error at %v in clientSide.Do %v \n", "TransferByPaystack", err.Error())
		err = fmt.Errorf("there was an internal server error while verifying your card. Please do not try the again as we'll take a look at the problem then email you later")
		return
	}
	defer res.Body.Close()
	if res.StatusCode >= 400 {
		err = fmt.Errorf("user payout method could not go through")
		return
	}
	if res == nil {
		err = fmt.Errorf("there was an internal server error while verifying your card. Please do not try the again as we'll take a look at the problem then email you later")
		return
	}
	err = json.NewDecoder(res.Body).Decode(&resData)
	if err != nil {
		log.Printf("error at TransferByPaystack in json.NewDecoder %v \n", err.Error())
		return
	}
	resItem = *resData
	return

}

func GetPayoutByReference(resPayouts map[uuid.UUID]PayoutData, reference string) (PayoutData, bool) {
	for _, v := range resPayouts {
		if v.PayoutID == reference {
			return v, true
		}
	}
	return PayoutData{}, false
}

func GetListRedisChargePayoutIDToRemove(redisDateIDs []string, chargeIDs []uuid.UUID) (redisRemoveDateIDs []string) {
	for _, chargeID := range chargeIDs {
		for _, r := range redisDateIDs {
			rSplit := strings.Split(r, "&")
			if len(rSplit) == 2 {
				id := rSplit[1]
				if id == tools.UuidToString(chargeID) {
					redisRemoveDateIDs = append(redisRemoveDateIDs, r)
					break
				}
			}

		}
	}
	return
}

func GetRedisChargeID(redisDateIDs []string, chargeID uuid.UUID) (redisChargeDateID string, exist bool) {
	for _, r := range redisDateIDs {
		rSplit := strings.Split(r, "&")
		if len(rSplit) == 2 {
			id := rSplit[1]
			if id == tools.UuidToString(chargeID) {
				redisChargeDateID = r
				exist = true
				return
			}
		}

	}
	return
}

// HandlePayoutRes this function remove redis date_ids that payout was not successful
func HandlePayoutRes(res TransferQueueResponse, redisDateIDs []string, resPayouts map[uuid.UUID]PayoutData) {
	for _, t := range res.Data {
		if t.Status != "received" {
			data, exist := GetPayoutByReference(resPayouts, t.Reference)
			if exist {
				// Get Redis ids to remove
				removeIDs := GetListRedisChargePayoutIDToRemove(redisDateIDs, data.ChargeIDs)
				if len(removeIDs) > 0 {
					err := RedisClient.SRem(RedisContext, constants.PAYOUT_CHARGE_DATE_IDS, removeIDs).Err()
					if err != nil {
						log.Println("This redis items were not removed from redis even though payout was unsuccessful", removeIDs)
					}
				}
			}
		}
	}
}
