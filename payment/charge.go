package payment

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	db "github.com/makuo12/ghost_server/db/sqlc"
	"github.com/makuo12/ghost_server/tools"
)

// This charges already store card on the db. This function would handle the paystack conversion
func HandlePaystackChargeAuthorization(ctx context.Context, paymentSuccessUrl string, paymentFailUrl string, paystackKey string, card db.Card, charge string) (resItem PaystackPaymentResponse, resChallengeItem PaymentPaymentChallengeResponse, resChallenged bool, err error) {
	var resData = &PaystackPaymentResponse{}
	var resChallengeData = &PaymentPaymentChallengeResponse{}
	amount := tools.ConvertToPaystackCharge(charge)
	callbackUrl := paymentSuccessUrl
	callbackFailUrl := paymentFailUrl
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
	var bearer = "Bearer " + paystackKey
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
