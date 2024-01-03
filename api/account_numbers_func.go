package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetBankList(server *Server, country string) (resBank BankData, err error) {
	var banks = &BankData{}
	payStackToken := server.config.PaystackSecretLiveKey
	url := "https://api.paystack.co/bank?country=" + strings.ToLower(country)
	var bearer = "Bearer " + payStackToken
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("error at %v in http.NewRequest %v \n", "GetBankList", err.Error())
		err = fmt.Errorf("there was an internal server error while getting the list of banks in your country. If this error continues contact help center")
		return
	}
	request.Header.Add("Authorization", bearer)
	// Send req using http Client
	clientSide := &http.Client{}
	res, err := clientSide.Do(request)
	if err != nil {
		log.Printf("error at %v in clientSide.Do %v \n", "GetBankList", err.Error())
		err = fmt.Errorf("there was an internal server error while getting the list of banks in your country. Please do not try the again as we'll take a look at the problem then email you later")
		return
	}
	defer res.Body.Close()
	if res.StatusCode >= 400 {
		log.Printf("error at %v in http.Get", "GetBankList")
		err = fmt.Errorf("there was an error while getting the data")
		return
	}
	if res == nil {
		err = fmt.Errorf("there was an internal server error while verifying your card. Please do not try the again as we'll take a look at the problem then email you later")
		return
	}
	err = json.NewDecoder(res.Body).Decode(&banks)
	if err != nil {
		log.Printf("error at GetBankList in json.NewDecoder %v \n", err.Error())
		err = fmt.Errorf("there was an internal server error while getting the list of banks in your country")
		return
	}
	// We check to see if the transaction failed
	if !banks.Status {
		log.Printf("Error at GetBankList payment did not go through")
		err = fmt.Errorf("there was an error while getting the list of banks in your country")
		return
	}
	resBank = *banks
	return
}

func VerifyAccountNumber(server *Server, ctx *gin.Context, accountNumber string, bankCode string) (resAccountData AccountData, err error) {
	var accountData = &AccountData{}
	url := fmt.Sprintf("https://api.paystack.co/bank/resolve?account_number=%v&bank_code=%v", accountNumber, bankCode)
	payStackToken := server.config.PaystackSecretLiveKey
	var bearer = "Bearer " + payStackToken
	log.Println("Bearer", bearer)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("error at %v in http.NewRequest %v \n", "VerifyAccountNumber", err.Error())
		err = fmt.Errorf("there was an internal server error while getting the list of banks in your country. If this error continues contact help center")
		return
	}
	request.Header.Add("Authorization", bearer)
	// Send req using http Client
	clientSide := &http.Client{}
	res, err := clientSide.Do(request)
	if err != nil {
		log.Printf("error at %v in clientSide.Do %v \n", "VerifyAccountNumber", err.Error())
		err = fmt.Errorf("there was an internal server error while getting the list of banks in your country. Please do not try the again as we'll take a look at the problem then email you later")
		return
	}
	defer res.Body.Close()
	if res.StatusCode >= 400 {
		log.Printf("error at %v in http.Get", "VerifyAccountNumber")
		err = fmt.Errorf("there was an error while getting the data")
		return
	}
	if res == nil {
		err = fmt.Errorf("there was an internal server error while verifying your card. Please do not try the again as we'll take a look at the problem then email you later")
		return
	}
	err = json.NewDecoder(res.Body).Decode(&accountData)
	if err != nil {
		log.Printf("error at VerifyAccountNumber in json.NewDecoder %v \n", err.Error())
		err = fmt.Errorf("there was an internal server error while verifying your account number")
		return
	}
	// We check to see if the transaction failed
	if !accountData.Status {
		log.Printf("Error at VerifyAccountNumber payment did not go through")
		err = fmt.Errorf("there was an error while verifying your account number")
		return
	}
	resAccountData = *accountData
	return
}

func CreateTransferRecipient(server *Server, ctx *gin.Context, accountNumber string, bankCode string, name string, currency string) (resRecipientData TransferRecipientRes, err error) {
	var recipientData = &TransferRecipientRes{}
	url := "https://api.paystack.co/transferrecipient"
	payStackToken := server.config.PaystackSecretLiveKey
	var bearer = "Bearer " + payStackToken
	log.Println("Bearer", bearer)
	var recipient = TransferRecipientParams{
		Type:          "nuban",
		Name:          name,
		AccountNumber: accountNumber,
		BankCode:      bankCode,
		Currency:      currency,
	}
	buf := new(bytes.Buffer)
	err = json.NewEncoder(buf).Encode(recipient)
	if err != nil {
		return
	}
	request, err := http.NewRequest("POST", url, buf)
	if err != nil {
		log.Printf("error at %v in http.NewRequest %v \n", "HandleMainRefund", err.Error())
		err = fmt.Errorf("there was an internal server error while verifying your card. If this error continues contact help center")
		return TransferRecipientRes{}, err
	}
	request.Header.Add("Authorization", bearer)
	// Send req using http Client
	clientSide := &http.Client{}
	res, err := clientSide.Do(request)
	if err != nil {
		log.Printf("error at %v in clientSide.Do %v \n", "HandleMainRefund", err.Error())
		err = fmt.Errorf("there was an internal server error while verifying your card. Please do not try the again as we'll take a look at the problem then email you later")
		return TransferRecipientRes{}, err
	}
	defer res.Body.Close()
	log.Println("resCode", res.StatusCode)
	if res.StatusCode == 400 {
		err = fmt.Errorf("user payment method could not go through")
		return TransferRecipientRes{}, err
	}
	if res == nil {
		err = fmt.Errorf("there was an internal server error while verifying your card. Please do not try the again as we'll take a look at the problem then email you later")
		return TransferRecipientRes{}, err
	}
	err = json.NewDecoder(res.Body).Decode(&recipientData)
	if err != nil {
		log.Printf("error at HandleMainRefund in json.NewDecoder %v \n", err.Error())
		return TransferRecipientRes{}, err
	}
	resRecipientData = *recipientData
	return
}

func HandleListBank(server *Server, ctx *gin.Context, country string) (res ListBankRes, err error) {
	banks, err := GetBankList(server, country)
	if err != nil {
		return
	}
	var resData []BankItem
	for _, b := range banks.Data {
		data := BankItem{
			Name:     b.Name,
			Code:     b.Code,
			Type:     b.Type,
			Currency: b.Currency,
		}
		resData = append(resData, data)
	}
	res = ListBankRes{
		List: resData,
	}
	return
}

func HandleVerifyBankCode(server *Server, ctx *gin.Context, country string, code string) (bank Bank, err error) {
	banks, err := GetBankList(server, country)
	if err != nil {
		return
	}
	for _, b := range banks.Data {
		if b.Code == code {
			bank = b
			return
		}
	}
	err = fmt.Errorf("bank for this account number is not found")
	return
}
