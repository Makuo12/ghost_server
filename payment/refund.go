package payment

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/makuo12/ghost_server/tools"
)

func HandlePaystackCreateRefund(paymentReference string, url string, userPercent int, totalFee string, bearer string) (resItem PaystackCreateRefundData, err error) {
	buf := new(bytes.Buffer)
	var resData = &PaystackCreateRefundData{}
	if userPercent == 100 {
		err = json.NewEncoder(buf).Encode(PaystackFullRefundParams{
			Transaction: paymentReference,
		})
		if err != nil {
			return
		}
	} else {
		amountAfterPercent := tools.ConvertFloatToString(float64(userPercent/100) * tools.ConvertStringToFloat(totalFee))

		err = json.NewEncoder(buf).Encode(PaystackPartialRefundParams{
			Transaction: paymentReference,
			Amount:      tools.ConvertToPaystackPayout(amountAfterPercent),
		})
		if err != nil {
			return
		}
	}
	request, err := http.NewRequest("POST", url, buf)
	if err != nil {
		log.Printf("error at %v in http.NewRequest %v \n", "HandleMainRefund", err.Error())
		err = fmt.Errorf("there was an internal server error while verifying your card. If this error continues contact help center")
		return
	}
	request.Header.Add("Authorization", bearer)
	// Send req using http Client
	clientSide := &http.Client{}
	res, err := clientSide.Do(request)
	if err != nil {
		log.Printf("error at %v in clientSide.Do %v \n", "HandleMainRefund", err.Error())
		err = fmt.Errorf("there was an internal server error while verifying your card. Please do not try the again as we'll take a look at the problem then email you later")
		return
	}
	defer res.Body.Close()
	log.Println("resCode", res.StatusCode)
	if res.StatusCode == 400 {
		err = fmt.Errorf("user payment method could not go through")
		return
	}
	if res == nil {
		err = fmt.Errorf("there was an internal server error while verifying your card. Please do not try the again as we'll take a look at the problem then email you later")
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

