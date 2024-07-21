package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/makuo12/ghost_server/constants"
	db "github.com/makuo12/ghost_server/db/sqlc"
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
	}
	return
}

func DailyHandleRefund(ctx context.Context, server *Server) func() {
	return func() {
		redisRefundDateIDs, err := RedisClient.SMembers(RedisContext, constants.REFUND_CHARGE_DATE_IDS).Result()
		if err != nil {
			log.Printf("Error at DailyHandlePayouts in SMembers err:%v\n", err.Error())
		}
		HandleOptionMainRefund(ctx, server, redisRefundDateIDs)
		HandleTicketMainRefund(ctx, server, redisRefundDateIDs)
		HandleAddCardRefund(ctx, server, redisRefundDateIDs)
	}
}

func HandleAddCardRefund(ctx context.Context, server *Server, redisRefunds []string) {
	refunds, err := server.store.ListMainRefundWithCharge(ctx, false)
	if err != nil || len(refunds) == 0 {
		if err != nil {
			log.Printf("HandleAddCardRefund in ListOptionMainRefundWithCharge err:%v\n", err.Error())
		}
	} else {
		url := "https://api.paystack.co/refund"
		bearer := "Bearer " + server.config.PaystackSecretLiveKey
		for _, r := range refunds {
			err := HandleMainRefund(ctx, server, r.ChargeID, uuid.New(), r.UID, "0", "0", r.Reference, url, bearer, int(r.UserPercent), int(r.HostPercent), r.Currency, redisRefunds)
			if err != nil {
				log.Printf("HandleAddCardRefund in HandleMainRefund err:%v\n", err.Error())
			}
		}
	}

}

func HandleOptionMainRefund(ctx context.Context, server *Server, redisRefunds []string) {
	refunds, err := server.store.ListOptionMainRefundWithCharge(ctx, false)
	if err != nil || len(refunds) == 0 {
		if err != nil {
			log.Printf("HandleOptionMainRefund in ListOptionMainRefundWithCharge err:%v\n", err.Error())
		}
	} else {
		url := "https://api.paystack.co/refund"
		bearer := "Bearer " + server.config.PaystackSecretLiveKey
		for _, r := range refunds {
			err := HandleMainRefund(ctx, server, r.ChargeID, r.HostID, r.UID, tools.IntToMoneyString(r.ServiceFee), tools.IntToMoneyString(r.TotalFee), r.PaymentReference, url, bearer, int(r.UserPercent), int(r.HostPercent), r.Currency, redisRefunds)
			if err != nil {
				log.Printf("HandleOptionMainRefund in HandleMainRefund err:%v\n", err.Error())
			}
		}
	}
}

func HandleTicketMainRefund(ctx context.Context, server *Server, redisRefunds []string) {
	refunds, err := server.store.ListTicketMainRefundWithCharge(ctx, false)
	if err != nil || len(refunds) == 0 {
		if err != nil {
			log.Printf("HandleTicketMainRefund in ListTicketMainRefundWithCharge err:%v\n", err.Error())
		}
	} else {
		url := "https://api.paystack.co/refund"
		bearer := "Bearer " + server.config.PaystackSecretLiveKey
		for _, r := range refunds {
			err := HandleMainRefund(ctx, server, r.ChargeID, r.HostID, r.UID, tools.IntToMoneyString(r.ServiceFee), tools.IntToMoneyString(r.TotalFee), r.PaymentReference, url, bearer, int(r.UserPercent), int(r.HostPercent), r.Currency, redisRefunds)
			if err != nil {
				log.Printf("HandleTicketMainRefund in HandleMainRefund err:%v\n", err.Error())
			}
		}
	}
}

func HandleMainRefund(ctx context.Context, server *Server, chargeID uuid.UUID, hostID uuid.UUID, uID uuid.UUID, serviceFee string, totalFee string, paymentReference, url string, bearer string, userPercent int, hostPercent int, currency string, redisRefunds []string) error {
	_, exist := GetRedisChargeID(redisRefunds, chargeID)
	if exist {
		log.Printf("ChargeID:%v is currently in redis so cannot resend the refund\n", chargeID)
		return nil
	}
	buf := new(bytes.Buffer)
	var resData = RefundData{}
	if userPercent == 0 {
		HandleRefundPayoutDB(ctx, server, chargeID, hostID, int(hostPercent), serviceFee, totalFee, "HandleMainRefund", currency)
		return nil
	}
	if userPercent == 100 {
		err := json.NewEncoder(buf).Encode(FullRefundParams{
			Transaction: paymentReference,
		})
		if err != nil {
			return err
		}
	} else {
		amountAfterPercent := tools.ConvertFloatToString(float64(userPercent/100) * tools.ConvertStringToFloat(totalFee))

		err := json.NewEncoder(buf).Encode(PartialRefundParams{
			Transaction: paymentReference,
			Amount:      tools.ConvertToPaystackPayout(amountAfterPercent),
		})
		if err != nil {
			return err
		}
	}
	request, err := http.NewRequest("POST", url, buf)
	if err != nil {
		log.Printf("error at %v in http.NewRequest %v \n", "HandleMainRefund", err.Error())
		err = fmt.Errorf("there was an internal server error while verifying your card. If this error continues contact help center")
		return err
	}
	request.Header.Add("Authorization", bearer)
	// Send req using http Client
	clientSide := &http.Client{}
	res, err := clientSide.Do(request)
	if err != nil {
		log.Printf("error at %v in clientSide.Do %v \n", "HandleMainRefund", err.Error())
		err = fmt.Errorf("there was an internal server error while verifying your card. Please do not try the again as we'll take a look at the problem then email you later")
		return err
	}
	defer res.Body.Close()
	log.Println("resCode", res.StatusCode)
	if res.StatusCode == 400 {
		err = fmt.Errorf("user payment method could not go through")
		return err
	}
	if res == nil {
		err = fmt.Errorf("there was an internal server error while verifying your card. Please do not try the again as we'll take a look at the problem then email you later")
		return err
	}
	err = json.NewDecoder(res.Body).Decode(&resData)
	if err != nil {
		log.Printf("error at HandleMainRefund in json.NewDecoder %v \n", err.Error())
		return err
	}
	log.Printf("main refund response %v\n", resData.Data)
	date := tools.ConvertTimeToString(time.Now())
	date_id := fmt.Sprintf("%v&%v", date, chargeID)
	// We store it in redis
	// We want to store the charge id in redis for refunds
	err = RedisClient.SAdd(RedisContext, constants.REFUND_CHARGE_DATE_IDS, date_id).Err()
	if err != nil {
		log.Printf("error at HandleMainRefund in RedisClient.SAdd %v \n", err.Error())
		err = nil
	}
	HandleRefundDB(ctx, server, chargeID, uID, "paystack", resData.Data.Amount, paymentReference, "HandleMainRefund")
	HandleRefundPayoutDB(ctx, server, chargeID, hostID, int(hostPercent), serviceFee, totalFee, "HandleMainRefund", currency)
	return nil

}

func HandleRefundDB(ctx context.Context, server *Server, chargeID uuid.UUID, uID uuid.UUID, sendMedium string, amount int, paymentReference string, funcName string) {
	_, err := server.store.CreateRefund(ctx, db.CreateRefundParams{
		ChargeID:   chargeID,
		Reference:  paymentReference,
		SendMedium: sendMedium,
		UserID:     uID,
		Amount:     int64(amount),
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

func HandleRefundSuccess(ctx context.Context, server *Server, data map[string]string, reference string, funcName string, redisIDs []string) {
	split := strings.Split(reference, "&")
	if len(split) != 2 {
		return
	}
	mainReference := split[1]
	refund, err := server.store.UpdateRefund(ctx, db.UpdateRefundParams{
		IsComplete: true,
		TimePaid:   time.Now(),
		Reference:  mainReference,
	})
	if err != nil {
		log.Printf("Error at funcName: %v, HandleRefundSuccess in store.UpdateRefund err:%v, mainReference: %v\n", funcName, err.Error(), mainReference)
	}
	// We want to update main_payout chargeIDs
	_, err = server.store.UpdateMainRefund(ctx, db.UpdateMainRefundParams{
		IsPayed:  true,
		ChargeID: refund.ChargeID,
	})
	if err != nil {
		log.Printf("Error at funcName: %v, HandleRefundSuccess in UpdateMainRefund err:%v, mainReference: %v\n", funcName, err.Error(), mainReference)
	}
	// Get Redis ids to remove
	rCID, exist := GetRedisChargeID(redisIDs, refund.ChargeID)
	if exist {
		err := RedisClient.SRem(RedisContext, constants.REFUND_CHARGE_DATE_IDS, rCID).Err()
		if err != nil {
			log.Println("This redis  HandleRefundSuccess items were not removed from redis even though payout was unsuccessful", refund.ChargeID)
		}
	}
	// We want to remove the data and mainReference from redis
	err = RedisClient.Del(RedisContext, reference).Err()
	if err != nil {
		log.Printf("Error at funcName: %v, HandleRefundSuccess in RedisClient.Del err:%v, reference: %v\n", funcName, err.Error(), reference)
	}
	err = RedisClient.SRem(RedisContext, constants.WEBHOOK_PAYSTACK_REFUND_REFERENCE, reference).Err()
	if err != nil {
		log.Printf("Error at funcName: %v, HandleRefundSuccess in RedisClient.SRem err:%v, reference: %v\n", funcName, err.Error(), reference)
	}

}
