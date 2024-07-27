package payout

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func HandlePaystackFetchRefund(ctx context.Context, paystackKey, refundID string) (resItem RefundResponse, err error) {
	var resData = &RefundResponse{}
	clientSide := &http.Client{}
	payStackToken := paystackKey
	url := "https://api.paystack.co/refund/" + refundID
	var bearer = "Bearer " + payStackToken
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("error at %v in http.NewRequest %v \n", "HandlePaystackFetchRefund", err.Error())
		err = fmt.Errorf("internal server error occurred while verifying your transaction %v", err.Error())
		return
	}
	//request.Close = true
	request.Header.Add("Authorization", bearer)
	// Send req using http Client

	res, err := clientSide.Do(request)
	if err != nil {
		log.Printf("error at %v in clientSide.Do %v \n", "HandlePaystackFetchRefund", err.Error())
		err = fmt.Errorf("internal server error occurred while verifying your transaction %v", err.Error())
		return
	}
	defer res.Body.Close()
	if res.StatusCode >= 400 {
		log.Printf("error at %v in http.Get", "HandlePaystackFetchRefund")
		err = fmt.Errorf("an error %v occurred so your transaction could not be verified,", err.Error())
		return
	}
	if res == nil {
		err = fmt.Errorf("no data ")
		return
	}
	rd := json.NewDecoder(res.Body)
	err = rd.Decode(&resData)
	if err != nil {
		log.Printf("error at HandlePaystackFetchRefund in json.NewDecoder %v \n", err.Error())
		err = fmt.Errorf("there was an internal server error while verifying your card. Please do not try the again as we'll take a look at the problem then email you later")
		return
	}

	// We check to see if the transaction failed
	if resData.Data.Status != "success" {
		log.Printf("Error at HandlePaystackFetchRefund payment did not go through")
		err = fmt.Errorf(resData.Data.Reason)
		return
	}
	log.Printf("reference %v\n", resData.Data.Amount)
	if err != nil {
		log.Printf("Error at StringToUuid payment was successful, but reference ID was not able to convert %v", err.Error())
		err = fmt.Errorf("please contact us with error code 404-203-ID, something went wrong. reference ID could not match")
		return
	}
	resItem = *resData
	return
}

func HandlePaystackListRefund(paystackKey string) (resItem PaystackRefundResponse, err error) {
	var resData = &PaystackRefundResponse{}
	url := "https://api.paystack.co/refund"
	var bearer = "Bearer " + paystackKey
	request, err := http.NewRequest("GET", url, nil)
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
		log.Printf("error at HandlePaystackChargeAuthorization in json.NewDecoder %v \n", err.Error())
	} else {
		resItem = *resData
	}
	return
}
