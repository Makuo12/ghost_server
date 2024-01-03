package api

import (
	"bytes"
	"context"
	"encoding/json"
	"flex_server/constants"
	db "flex_server/db/sqlc"
	"flex_server/tools"
	"fmt"
	"strings"

	//"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

func HandlePaystackReference(server *Server, ctx *gin.Context, reference string, user db.User) (resItem PaystackVerifyResponse, err error) {
	log.Println("reference is", reference)
	var resData = &PaystackVerifyResponse{}
	clientSide := &http.Client{}
	payStackToken := server.config.PaystackSecretLiveKey

	url := "https://api.paystack.co/transaction/verify/" + reference
	var bearer = "Bearer " + payStackToken
	log.Println("Bearer", bearer)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("error at %v in http.NewRequest %v \n", "HandlePaystackReference", err.Error())
		err = fmt.Errorf("there was an internal server error while getting the list of banks in your country. If this error continues contact help center")
		return
	}
	request.Close = true

	request.Header.Add("Authorization", bearer)
	// Send req using http Client

	res, err := clientSide.Do(request)
	if err != nil {
		log.Printf("error at %v in clientSide.Do %v \n", "HandlePaystackReference", err.Error())
		err = fmt.Errorf("there was an internal server error while getting the list of banks in your country. Please do not try the again as we'll take a look at the problem then email you later")
		return
	}
	defer res.Body.Close()
	if res.StatusCode >= 400 {
		log.Printf("error at %v in http.Get", "HandlePaystackReference")
		err = fmt.Errorf("there was an error while getting the data")
		return
	}
	if res == nil {
		err = fmt.Errorf("there was an internal server error while verifying your card. Please do not try the again as we'll take a look at the problem then email you later")
		return
	}
	rd := json.NewDecoder(res.Body)
	err = rd.Decode(&resData)
	if err != nil {
		log.Printf("error at HandlePaystackReference in json.NewDecoder %v \n", err.Error())
		err = fmt.Errorf("there was an internal server error while verifying your card. Please do not try the again as we'll take a look at the problem then email you later")
		return
	}

	// We check to see if the transaction failed
	if resData.Data.Status != "success" {
		log.Printf("Error at HandlePaystackReference payment did not go through")
		err = fmt.Errorf("there was an error while verifying your card. Please try again if the error continues try using another card or help center")
		return
	}
	log.Printf("reference %v\n", resData.Data.Reference)
	if err != nil {
		log.Printf("Error at StringToUuid payment was successful, but reference ID was not able to convert %v", err.Error())
		err = fmt.Errorf("please contact us with error code 404-203-ID, something went wrong. reference ID could not match")
		return
	}
	resItem = *resData
	return
}

// This charges already store card on the db. This function would handle the paystack conversion
func HandlePaystackChargeAuthorization(server *Server, ctx context.Context, card db.Card, charge string) (resItem PaystackPaymentResponse, resChallengeItem PaymentPaymentChallengeResponse, resChallenged bool, err error) {
	var resData = &PaystackPaymentResponse{}
	var resChallengeData = &PaymentPaymentChallengeResponse{}
	amount := tools.ConvertToPaystackCharge(charge)
	callbackUrl := server.config.PaymentSuccessUrl
	callbackFailUrl := server.config.PaymentFailUrl
	metaData := PaystackPaymentMetaData{
		CancelAction: callbackFailUrl,
	}

	data := PaystackAuthorization{
		AuthorizationCode: card.AuthorizationCode,
		Email:             card.Email,
		Amount:            amount,
		CallbackUrl:       callbackUrl,
		MetaData:          metaData,
	}
	url := "https://api.paystack.co/transaction/charge_authorization"
	var bearer = "Bearer " + server.config.PaystackSecretLiveKey
	buf := new(bytes.Buffer)
	err = json.NewEncoder(buf).Encode(data)
	if err != nil {
		return
	}
	request, err := http.NewRequest("POST", url, buf)
	if err != nil {
		log.Printf("error at %v in http.NewRequest %v \n", "HandlePaystackChargeAuthorization", err.Error())
		err = fmt.Errorf("an error occurred while making the payment, try again. %v", err.Error())
		return
	}
	request.Header.Add("Authorization", bearer)
	// Send req using http Client
	clientSide := &http.Client{}
	res, err := clientSide.Do(request)
	if err != nil {
		log.Printf("error at %v in clientSide.Do %v \n", "HandlePaystackChargeAuthorization", err.Error())
		err = fmt.Errorf("an error occurred while making the payment, try again. %v", err.Error())
		return
	}
	defer res.Body.Close()
	if res.StatusCode >= 400 {
		err = fmt.Errorf("your payment method could not go through, ensure you have the required funds before making the purchase")
		return
	}
	if res == nil {
		err = fmt.Errorf("your payment method could not go through, ensure you have the required funds before making the purchase")
		return
	}

	err = json.NewDecoder(res.Body).Decode(&resData)
	if err != nil {
		// This error might occur because the card was challenge so we want to use the resChallengeData to hold it
		log.Printf("error at HandlePaystackChargeAuthorization in json.NewDecoder %v \n", err.Error())
		// We want to set the err back to nil
		err = json.NewDecoder(res.Body).Decode(&resChallengeData)
		if err != nil {
			log.Printf("error at HandlePaystackChargeAuthorization in json.NewDecoder %v \n", err.Error())
			// If there is an error here we know that the transaction failed
			err = fmt.Errorf("payment method failed error %v, try again", err.Error())
			return
		} else {
			err = nil
			resChallenged = true
			resChallengeItem = *resChallengeData
			return
		}
	} else {
		resItem = *resData
	}
	return
}

func (server *Server) VerifyAddCardChargeReference(ctx *gin.Context) {
	var req ReferencePayment
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at VerifyAddCardChargeReference in ShouldBindJSON: %v, currency: %v \n", err.Error(), req.Reference)
		err = fmt.Errorf("reference invalid")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	reference, err := tools.StringToUuid(req.Reference)
	if err != nil {
		log.Printf("Error at VerifyAddCardChargeReference in tools.StringToUuid %v, currency: %v \n", err.Error(), req.Reference)
		err = fmt.Errorf("reference expired")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := HandleGetUser(ctx, server)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	chargeData, err := RedisClient.Get(RedisContext, req.Reference).Result()
	if err != nil {
		log.Printf("Error at VerifyAddCardChargeReference in RedisClient.Get(ctx %v, currency: %v \n", err.Error(), req.Reference)
		err = fmt.Errorf("reference expired")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	split := strings.Split(chargeData, "&")
	if len(split) != 4 {
		err = fmt.Errorf("reference expired")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	userIDString, currency, charge, reason := split[0], split[1], split[2], split[3]

	userID, err := tools.StringToUuid(userIDString)
	if err != nil || userID != user.UserID {
		err = fmt.Errorf("user do not match the one with charge")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	resData, err := HandlePaystackReference(server, ctx, req.Reference, user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	chargeRef, err := server.store.CreateChargeReference(ctx, db.CreateChargeReferenceParams{
		UserID:     userID,
		Reason:     reason,
		Charge:     tools.MoneyStringToInt(charge),
		Currency:   currency,
		IsComplete: true,
		Reference:  reference,
	})
	if err != nil {
		log.Printf("Error at VerifyAddCardChargeReference in CreateChargeReference: %v, currency: %v \n", err.Error(), req.Reference)
		err = fmt.Errorf("verification of your card was made successful, however an error occurred contacts us. error code E4-GO")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	accountName := fmt.Sprintf("%v", resData.Data.Authorization.AccountName)
	argsCard := db.CreateCardParams{
		UserID:            user.ID,
		Email:             resData.Data.Customer.Email,
		AuthorizationCode: resData.Data.Authorization.AuthorizationCode,
		CardType:          resData.Data.Authorization.CardType,
		Last4:             resData.Data.Authorization.Last4,
		ExpMonth:          resData.Data.Authorization.ExpMonth,
		ExpYear:           resData.Data.Authorization.ExpYear,
		Bank:              resData.Data.Authorization.Bank,
		CountryCode:       resData.Data.Authorization.CountryCode,
		Reusable:          resData.Data.Authorization.Reusable,
		Channel:           resData.Data.Authorization.Channel,
		AccountName:       accountName,
		Bin:               resData.Data.Authorization.Bin,
		CardSignature:     resData.Data.Authorization.Signature,
		Currency:          resData.Data.Currency,
	}
	card, err := server.store.CreateCard(ctx, argsCard)
	if err != nil {
		log.Printf("Error at VerifyAddCardChargeReference in CreateCard. Error is %v for AuthorizationCode: %v and user: %v \n", err, resData.Data.Authorization.AuthorizationCode, user.ID)
		err := fmt.Errorf("there was an error while verifying your card. Please don't try just again we are working on it")

		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	// we delete it from paystack
	err = RedisClient.Del(RedisContext, req.Reference).Err()
	if err != nil {
		log.Printf("Error at VerifyAddCardChargeReference in RedisClient.Del(ctx %v, currency: %v \n", err.Error(), req.Reference)
		err = nil
	}

	// update user default card
	switch req.Type {
	case "payout":
		_, err = server.store.UpdateUser(ctx, db.UpdateUserParams{
			DefaultPayoutCard: pgtype.Text{
				String: tools.UuidToString(card.ID),
				Valid:  true,
			},
			ID: user.ID,
		})
		if err != nil {
			log.Printf("Error at VerifyAddCardChargeReference payout in UpdateUser. Error is %v for DefaultID: %v and user: %v \n", err, card.ID, user.ID)
		}

	case "payment":
		_, err = server.store.UpdateUser(ctx, db.UpdateUserParams{
			DefaultCard: pgtype.Text{
				String: tools.UuidToString(card.ID),
				Valid:  true,
			},
			ID: user.ID,
		})
		if err != nil {
			log.Printf("Error at VerifyAddCardChargeReference payment in UpdateUser. Error is %v for DefaultID: %v and user: %v \n", err, card.ID, user.ID)
		}

	}

	// We create a refund to send all the money back
	_, err = server.store.CreateMainRefund(ctx, db.CreateMainRefundParams{
		ChargeID:    chargeRef.ID,
		UserPercent: 100,
		HostPercent: 0,
		ChargeType:  constants.CHARGE_REFERENCE,
		Type:        constants.ADD_CARD,
	})

	if err != nil {
		log.Printf("Error at VerifyAddCardChargeReference in CreateMainRefund. Error is %v for AuthorizationCode: %v and user: %v \n", err, resData.Data.Authorization.AuthorizationCode, user.ID)
	}
	// If card is not reusable then we send an error
	if !resData.Data.Authorization.Reusable {
		err := fmt.Errorf("there was an error while verifying your card. Please note that this card is not reusable so it cannot be used as a payment method. You would be sent your charge refund shortly")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	resCard := CardDetailResponse{
		CardID:    tools.UuidToString(card.ID),
		CardLast4: card.Last4,
		CardType:  card.CardType,
		ExpMonth:  card.ExpMonth,
		ExpYear:   card.ExpYear,
		Currency:  card.Currency,
	}

	res := CardAddResponse{
		CardDetail:      resCard,
		Account:         "0.00",
		AccountCurrency: chargeRef.Currency,
		DefaultID:       user.DefaultCard,
	}
	ctx.JSON(http.StatusOK, res)
}
