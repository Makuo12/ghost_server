package payout

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func HandlePaystackVerifyPayout(ctx context.Context, paystackKey string, reference string) (resItem PaystackTransferResponse, err error) {
	var resData = &PaystackTransferResponse{}
	clientSide := &http.Client{}
	payStackToken := paystackKey
	url := "https://api.paystack.co/transfer/verify/" + reference
	var bearer = "Bearer " + payStackToken
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("error at %v in http.NewRequest %v \n", "HandlePaystackVerifyPayment", err.Error())
		err = fmt.Errorf("internal server error occurred while verifying your transaction %v", err.Error())
		return
	}
	//request.Close = true
	request.Header.Add("Authorization", bearer)
	// Send req using http Client

	res, err := clientSide.Do(request)
	if err != nil {
		log.Printf("error at %v in clientSide.Do %v \n", "HandlePaystackVerifyPayment", err.Error())
		err = fmt.Errorf("internal server error occurred while verifying your transaction %v", err.Error())
		return
	}
	defer res.Body.Close()
	if res.StatusCode >= 400 {
		log.Printf("error at %v in http.Get", "HandlePaystackVerifyPayment")
		err = fmt.Errorf("an error occurred so your transaction could not be verified")
		return
	}
	if res == nil {
		err = fmt.Errorf("no data ")
		return
	}
	rd := json.NewDecoder(res.Body)
	err = rd.Decode(&resData)
	if err != nil {
		log.Printf("error at HandlePaystackVerifyPayment in json.NewDecoder %v \n", err.Error())
		err = fmt.Errorf("there was an internal server error while verifying your card. Please do not try the again as we'll take a look at the problem then email you later")
		return
	}

	// We check to see if the transaction failed
	if resData.Data.Status != "success" {
		log.Printf("Error at HandlePaystackVerifyPayment payment did not go through")
		err = fmt.Errorf(resData.Data.Reason)
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

func TransferByPaystack(ctx context.Context, paystackKey string, bulkTransfer BulkTransferRequest) (resItem TransferQueueResponse, err error) {
	url := "https://api.paystack.co/transfer/bulk"
	var bearer = "Bearer " + paystackKey
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
