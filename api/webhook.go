package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/makuo12/ghost_server/constants"
	"github.com/makuo12/ghost_server/tools"

	"github.com/gin-gonic/gin"
)

func (server *Server) PaystackWebhook(ctx *gin.Context) {
	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Printf("Error at PaystackWebhook ioutil.ReadAll: %v for clientIP %v", err, ctx.ClientIP())
	}
	paystackToken := server.config.PaystackSecretLiveKey

	reqSignature := ctx.GetHeader("x-paystack-signature")
	hash, err := tools.HandleHMAC(body, paystackToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid signature"})
		return
	}
	if hash != reqSignature {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid signature"})
		return
	}
	var webhookEvent = WebhookEvent{}
	// Unmarshal the JSON data into the WebhookEvent struct
	if err := json.Unmarshal(body, &webhookEvent); err != nil {
		log.Println("Error PaystackWebhook at HandlePayoutWebhook parsing JSON:", err)
		log.Println("err: ", err)
	}
	if IsRefund(webhookEvent.Event) {
		var refund = RefundEvent{}
		// Unmarshal the JSON data into the RefundEvent struct
		if err := json.Unmarshal(body, &refund); err != nil {
			log.Println("Error PaystackWebhook at HandleRefundWebhook parsing JSON:", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "wrong data format"})
			return
		}
		log.Println("r data: ", refund)
		// Lets make sure the uuid is right
		_, err := tools.StringToUuid(refund.Data.TransactionReference)
		if err != nil {
			log.Println("Error PaystackWebhook at HandleRefundWebhook uuid:", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		err = HandleRefundWebhook(ctx, server, body)
		if err != nil {
			log.Println("Error PaystackWebhook at HandleRefundWebhook result:", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

	} else if IsTransfer(webhookEvent.Event) {
		var transferResponse = TransferResponse{}
		// Unmarshal the JSON data into the TransferResponse struct
		if err := json.Unmarshal(body, &transferResponse); err != nil {
			log.Println("Error PaystackWebhook at HandlePayoutWebhook parsing JSON:", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "wrong data format"})
			return
		}
		log.Println("t data: ", transferResponse)
		// Lets make sure the uuid is right
		_, err := tools.StringToUuid(transferResponse.Data.Reference)
		if err != nil {
			log.Println("Error PaystackWebhook at HandleRefundWebhook uuid:", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		}
		err = HandlePayoutWebhook(ctx, server, body)
		if err != nil {
			log.Println("Error PaystackWebhook at HandlePayoutWebhook result:", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
	}
	ctx.JSON(http.StatusOK, "data collected")
}

func HandlePayoutWebhook(ctx context.Context, server *Server, body []byte) error {
	// Create an instance of TransferResponse
	var transferResponse = TransferResponse{}
	// Unmarshal the JSON data into the TransferResponse struct
	if err := json.Unmarshal(body, &transferResponse); err != nil {
		log.Println("Error at HandlePayoutWebhook parsing JSON:", err)
		return err
	}
	// This is were we want to handle updates to payout
	// We are just gonna store the data in redis
	data := []string{
		constants.TRANSFER_ACTION,
		transferResponse.Event,
		constants.TRANSFER_CODE,
		transferResponse.Data.TransferCode,
		constants.PAYOUT_REFERENCE,
		transferResponse.Data.Reference,
		constants.AMOUNT,
		tools.PaystackMoneyToDB(transferResponse.Data.Amount),
	}
	reference := fmt.Sprintf("%v&%v", constants.TRANSFER_ACTION, transferResponse.Data.Reference)
	err := RedisClient.HSet(RedisContext, reference, data).Err()
	if err != nil {
		log.Printf("Error at HandlePayoutWebhook RedisClient.HSet: %v for data %v", err, data)
		return err
	}
	err = RedisClient.SAdd(RedisContext, constants.WEBHOOK_PAYSTACK_TRANSFER_REFERENCE, reference).Err()
	if err != nil {
		log.Printf("Error at HandlePayoutWebhook RedisClient.SAdd: %v for reference %v", err, reference)
		return err
	}
	return nil
}

func HandleRefundWebhook(ctx context.Context, server *Server, body []byte) error {
	// Create an instance of RefundEvent
	var refund = RefundEvent{}
	// Unmarshal the JSON data into the RefundEvent struct
	if err := json.Unmarshal(body, &refund); err != nil {
		log.Println("Error at HandleRefundWebhook parsing JSON:", err)
		return err
	}
	// This is were we want to handle updates to payout
	// We are just gonna store the data in redis
	data := []string{
		constants.REFUND_ACTION,
		refund.Event,
		constants.AMOUNT,
		tools.PaystackMoneyToDB(refund.Data.Amount),
		constants.REFERENCE,
		refund.Data.TransactionReference,
	}
	reference := fmt.Sprintf("%v&%v", constants.REFUND_ACTION, refund.Data.TransactionReference)
	err := RedisClient.HSet(RedisContext, reference, data).Err()
	if err != nil {
		log.Printf("Error at HandleRefundWebhook RedisClient.HSet: %v for data %v", err, data)
		return err
	}
	err = RedisClient.SAdd(RedisContext, constants.WEBHOOK_PAYSTACK_REFUND_REFERENCE, reference).Err()
	if err != nil {
		log.Printf("Error at HandleRefundWebhook RedisClient.SAdd: %v for reference %v", err, refund.Data.TransactionReference)
		return err
	}
	return nil
}

func DailyHandleTransferWebhookData(ctx context.Context, server *Server) func() {
	// All the references are stored in WEBHOOK_PAYSTACK_REFERENCE
	return func() {
		redisChargeDateIDs, err := RedisClient.SMembers(RedisContext, constants.PAYOUT_CHARGE_DATE_IDS).Result()
		if err != nil {
			log.Printf("Error at DailyHandlePayouts in SMembers err:%v\n", err.Error())
		}
		references, err := RedisClient.SMembers(RedisContext, constants.WEBHOOK_PAYSTACK_TRANSFER_REFERENCE).Result()
		if err != nil {
			log.Printf("Error at DailyHandleWebhookData in SMembers err:%v\n", err.Error())
			return
		}

		// We want to remove the reference from PAYOUT_CHARGE_DATE_IDS because we are no more awaiting a response
		for _, r := range references {
			// Remember that this reference is the same with payout.id

			data, err := RedisClient.HGetAll(ctx, r).Result()
			if err != nil {
				log.Printf("Error at DailyHandleWebhookData in RedisClient.HGetAll err:%v, reference: %v\n", err.Error(), r)
				continue
			}

			if data[constants.TRANSFER_ACTION] == "transfer.success" {
				HandleTransferSuccess(ctx, server, data, r, "DailyHandleWebhookData", redisChargeDateIDs)
			} else {
				HandleTransferNotSuccess(ctx, server, data, r, "DailyHandleWebhookData", redisChargeDateIDs)
			}

		}
	}
}

func DailyHandleRefundWebhookData(ctx context.Context, server *Server) func() {
	// All the references are stored in WEBHOOK_PAYSTACK_REFERENCE
	return func() {
		redisRefundDateIDs, err := RedisClient.SMembers(RedisContext, constants.REFUND_CHARGE_DATE_IDS).Result()
		if err != nil {
			log.Printf("Error at DailyHandlePayouts in SMembers err:%v\n", err.Error())
		}
		references, err := RedisClient.SMembers(RedisContext, constants.WEBHOOK_PAYSTACK_REFUND_REFERENCE).Result()
		if err != nil {
			log.Printf("Error at DailyHandleWebhookData in SMembers err:%v\n", err.Error())
			return
		}
		// We want to remove the reference from REFUND_CHARGE_DATE_IDS because we are no more awaiting a response
		for _, r := range references {

			// Remember that this reference is the same with payout.id
			data, err := RedisClient.HGetAll(ctx, r).Result()
			if err != nil {
				log.Printf("Error at DailyHandleWebhookData in RedisClient.HGetAll err:%v, reference: %v\n", err.Error(), r)
				continue
			}
			if data[constants.TRANSFER_ACTION] == "refund.processed" {
				HandleRefundSuccess(ctx, server, data, r, "DailyHandleWebhookData", redisRefundDateIDs)
			}

		}
	}
}
