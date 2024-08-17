package payment

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/makuo12/ghost_server/tools"
)

func HandlePaystackCard(ctx context.Context, paystackKey string, charge string, currency string, email string, reason string) (resItem InitCardChargeRes, err error) {
	paystackCharge := ConvertToPaystackCharge(charge)
	resItem = InitCardChargeRes{
		Reference: uuid.New().String(),
		Reason:    reason,
		Charge:    paystackCharge,
		Currency:  currency,
		Email:     email,
	}
	return
}

func HandlePaystackBankAccount(ctx context.Context, paystackKey string, charge string, arg PaystackGetBankAccountParams, email string) (resItem PaystackBankChargeAttempt, err error) {
	var resData = &PaystackBankChargeAttempt{}
	var resValidationError = &PaystackValidationError{}
	var resApiError = &PaystackApiError{}
	amount := tools.ConvertToPaystackChargeString(charge)
	bankData := PaystackBankAccountObject{
		Email:         email,
		Amount:        amount,
		Code:          arg.Code,
		Phone:         arg.Phone,
		AccountNumber: arg.AccountNumber,
		Token:         arg.Token,
		//Reference:     reference,
	}

	url := "https://api.paystack.co/charge"
	var bearer = "Bearer " + paystackKey
	buf := new(bytes.Buffer)
	if arg.Code == "50211" {
		err = json.NewEncoder(buf).Encode(bankData.GetPhone())

	} else {
		err = json.NewEncoder(buf).Encode(bankData.GetBank())
	}
	if err != nil {
		return
	}
	request, err := http.NewRequest("POST", url, buf)
	if err != nil {
		log.Printf("error at %v in http.NewRequest %v \n", "HandlePaystackBankAccount", err.Error())
		err = fmt.Errorf("an error occurred while making the payment, try again. %v", err.Error())
		return
	}
	request.Header.Add("Authorization", bearer)
	// Send req using http Client
	clientSide := &http.Client{}
	res, err := clientSide.Do(request)
	if err != nil {
		log.Printf("error at %v in clientSide.Do %v \n", "HandlePaystackBankAccount", err.Error())
		err = fmt.Errorf("an error occurred while making the payment, try again. %v", err.Error())
		return
	}
	if res == nil {
		err = fmt.Errorf("your payment method could not go through, try using another payment channel")
		return
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		// First we try to parse the data
		err = json.NewDecoder(res.Body).Decode(&resValidationError)
		if err != nil {
			err = fmt.Errorf(resApiError.Data.Message)
			if err != nil {
				err = fmt.Errorf("your payment method could not go through, try using another payment channel")
				return
			} else {
				err = fmt.Errorf(resApiError.Data.Message)
				return
			}
		} else {
			err = fmt.Errorf(resValidationError.Message)
		}
	}

	err = json.NewDecoder(res.Body).Decode(&resData)
	if err != nil {
		err = fmt.Errorf("your payment method could not go through, try using another payment channel")
		return
	} else {
		resItem = *resData
	}
	return
}

func HandlePaystackPWT(ctx context.Context, paystackKey string, charge string, email string) (resItem PaystackPWTRes, err error) {
	var resData = &PaystackPWTRes{}
	var resValidationError = &PaystackValidationError{}
	amount := tools.ConvertToPaystackChargeString(charge)
	transferTime := PaystackBankTransfer{time.Now().Add(time.Minute * 30).Format("2006-01-02T15:04:05Z")}
	bankData := PaystackPaymentWithTransferParams{
		Email:  email,
		Amount: amount,
		//Reference:    reference,
		BankTransfer: transferTime,
	}
	url := "https://api.paystack.co/charge"
	var bearer = "Bearer " + paystackKey
	buf := new(bytes.Buffer)
	err = json.NewEncoder(buf).Encode(bankData)
	if err != nil {
		return
	}
	request, err := http.NewRequest("POST", url, buf)
	if err != nil {
		log.Printf("error at %v in http.NewRequest %v \n", "HandlePaystackPWT", err.Error())
		err = fmt.Errorf("an error occurred while making the payment, try again. %v", err.Error())
		return
	}
	request.Header.Add("Authorization", bearer)
	// Send req using http Client
	clientSide := &http.Client{}
	res, err := clientSide.Do(request)
	if err != nil {
		log.Printf("error at %v in clientSide.Do %v \n", "HandlePaystackPWT", err.Error())
		err = fmt.Errorf("an error occurred while making the payment, try again. %v", err.Error())
		return
	}
	if res == nil {
		err = fmt.Errorf("your payment method could not go through, try using another payment channel")
		return
	}
	defer res.Body.Close()
	if res.StatusCode >= 400 {
		// First we try to parse the data
		err = json.NewDecoder(res.Body).Decode(&resValidationError)
		if err != nil {
			err = fmt.Errorf("your payment method could not go through, try using another payment channel")
		} else {
			err = fmt.Errorf(resValidationError.Message)
			return
		}
	}
	err = json.NewDecoder(res.Body).Decode(&resData)
	if err != nil {
		err = fmt.Errorf("your payment method could not go through, try using another payment channel")
		return
	} else {
		resItem = *resData
	}
	return
}

func HandlePaystackUSSD(ctx context.Context, paystackKey string, charge string, ussdCode string, email string, firstName string) (resItem PaystackUSSDChargeAttempt, err error) {
	var resData = &PaystackUSSDChargeAttempt{}
	var resApiError = &PaystackApiError{}
	amount := tools.ConvertToPaystackChargeString(charge)
	metaData := PaystackUSSDMetaData{
		CustomFields: []PaystackUSSDFieldItem{{firstName, "Flizzup USSD Payment", "flizzup_ussd_payment"}},
	}
	ussdType := PaystackUSSDType{ussdCode}
	bankData := PaystackUSSDParams{
		Email:  email,
		Amount: amount,
		//Reference: reference,
		USSD:     ussdType,
		MetaData: metaData,
	}
	url := "https://api.paystack.co/charge"
	var bearer = "Bearer " + paystackKey
	buf := new(bytes.Buffer)
	err = json.NewEncoder(buf).Encode(bankData)
	if err != nil {
		return
	}
	request, err := http.NewRequest("POST", url, buf)
	if err != nil {
		log.Printf("error at %v in http.NewRequest %v \n", "HandlePaystackUSSD", err.Error())
		err = fmt.Errorf("an error occurred while making the payment, try again. %v", err.Error())
		return
	}
	request.Header.Add("Authorization", bearer)
	// Send req using http Client
	clientSide := &http.Client{}
	res, err := clientSide.Do(request)
	if err != nil {
		log.Printf("error at %v in clientSide.Do %v \n", "HandlePaystackUSSD", err.Error())
		err = fmt.Errorf("an error occurred while making the payment, try again. %v", err.Error())
		return
	}
	if res == nil {
		err = fmt.Errorf("your payment method could not go through, try using another payment channel")
		return
	}
	defer res.Body.Close()
	if res.StatusCode >= 400 {
		// First we try to parse the data
		err = json.NewDecoder(res.Body).Decode(&resApiError)
		if err != nil {
			err = fmt.Errorf("your payment method could not go through, try using another payment channel")
		} else {
			err = fmt.Errorf(resApiError.Data.Message)
			return
		}
	}
	err = json.NewDecoder(res.Body).Decode(&resData)
	if err != nil {
		err = fmt.Errorf("your payment method could not go through, try using another payment channel")
		return
	} else {
		resItem = *resData
	}
	return
}
