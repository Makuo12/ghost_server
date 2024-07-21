package payment

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	db "github.com/makuo12/ghost_server/db/sqlc"
)

// This function help verify any payment using paystack that has a reference
func HandlePaystackVerifyPayment(ctx context.Context, paystackKey, reference string, user db.User) (resItem PaystackVerifyResponse, err error) {
	log.Println("reference is", reference)
	var resData = &PaystackVerifyResponse{}
	clientSide := &http.Client{}
	payStackToken := paystackKey

	url := "https://api.paystack.co/transaction/verify/" + reference
	var bearer = "Bearer " + payStackToken
	log.Println("Bearer", bearer)
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
		log.Printf("error at HandlePaystackVerifyPayment in json.NewDecoder %v \n", err.Error())
		err = fmt.Errorf("there was an internal server error while verifying your card. Please do not try the again as we'll take a look at the problem then email you later")
		return
	}

	// We check to see if the transaction failed
	if resData.Data.Status != "success" {
		log.Printf("Error at HandlePaystackVerifyPayment payment did not go through")
		err = fmt.Errorf(resData.Data.GatewayResponse)
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


func GetFakeCardRes() CardAddResponse {
	detail := CardDetailResponse{"none", "none", "none", "none", "none", "none"}
	card := CardAddResponse {detail, "none", "none", "none"}
	return card
}