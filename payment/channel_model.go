package payment

import (
	"log"
	"math"
	"strconv"
)

// Channels refer to the various payment channels like bank accounts, cards, etc.
// / PAYSTACK ///
// Bank Account Channel
type PaystackBankAccountBankData struct {
	Code          string `json:"code"`
	AccountNumber string `json:"account_number"`
}
type PaystackBankAccountPhoneData struct {
	Code  string `json:"code"`
	Phone string `json:"phone"`
	Token string `json:"token"`
}

type PaystackGetBankAccountParams struct {
	Code          string `json:"code"`
	Phone         string `json:"phone"`
	Token         string `json:"token"`
	AccountNumber string `json:"account_number"`
}

type PaystackBankAccountBankParams struct {
	Email     string                      `json:"email"`
	Amount    string                      `json:"amount"`
	Reference string                      `json:"reference"`
	Bank      PaystackBankAccountBankData `json:"bank"`
}

type PaystackBankAccountPhoneParams struct {
	Email     string                       `json:"email"`
	Amount    string                       `json:"amount"`
	Reference string                       `json:"reference"`
	Bank      PaystackBankAccountPhoneData `json:"bank"`
}

type PaystackBankAccountObject struct {
	Email         string
	Amount        string
	Code          string
	Phone         string
	AccountNumber string
	Token         string
	Reference     string
}
type PaystackValidationError struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Meta    struct {
		NextStep string `json:"nextStep"`
	} `json:"meta"`
	Type string `json:"type"`
	Code string `json:"code"`
}

type PaystackApiError struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Reference string `json:"reference"`
		Status    string `json:"status"`
		Message   string `json:"message"`
	} `json:"data"`
	Meta struct {
		NextStep string `json:"nextStep"`
	} `json:"meta"`
	Type string `json:"type"`
	Code string `json:"code"`
}

func (p *PaystackBankAccountObject) GetPhone() PaystackBankAccountPhoneParams {
	bank := PaystackBankAccountPhoneData{
		Code:  p.Code,
		Phone: p.Phone,
		Token: p.Token,
	}
	res := PaystackBankAccountPhoneParams{
		Email:     p.Email,
		Amount:    p.Amount,
		Reference: p.Reference,
		Bank:      bank,
	}
	return res
}

func (p *PaystackBankAccountObject) GetBank() PaystackBankAccountBankParams {
	bank := PaystackBankAccountBankData{
		Code:          p.Code,
		AccountNumber: p.AccountNumber,
	}
	res := PaystackBankAccountBankParams{
		Email:     p.Email,
		Amount:    p.Amount,
		Reference: p.Reference,
		Bank:      bank,
	}
	return res
}

type PaystackBankChargeAttempt struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Reference   string `json:"reference"`
		Status      string `json:"status"`
		DisplayText string `json:"display_text"`
	} `json:"data"`
}

type PaystackBankAccountMainRes struct {
	Reference   string `json:"reference"`
	DisplayText string `json:"display_text"`
}

// End Bank Account Channel

// Pay With Transfer
type PaystackBankTransfer struct {
	AccountExpiresAt string `json:"account_expires_at"`
}

type PaystackPaymentWithTransferParams struct {
	Email        string               `json:"email"`
	Amount       string               `json:"amount"`
	Reference    string               `json:"reference"`
	BankTransfer PaystackBankTransfer `json:"bank_transfer"`
}

type PaystackPWTBank struct {
	Slug string `json:"slug"`
	Name string `json:"name"`
	ID   int    `json:"id"`
}

type PaystackPWTRes struct {
	Status  bool               `json:"status"`
	Message string             `json:"message"`
	Data    PaystackPWTResData `json:"data"`
}

type PaystackPWTResData struct {
	Reference        string             `json:"reference"`
	Status           string             `json:"status"`
	DisplayText      string             `json:"display_text"`
	AccountName      string             `json:"account_name"`
	AccountNumber    string             `json:"account_number"`
	Bank             PaystackPWTResBank `json:"bank"`
	AccountExpiresAt string             `json:"account_expires_at"`
}

type PaystackPWTResBank struct {
	Slug string `json:"slug"`
	Name string `json:"name"`
	ID   int    `json:"id"`
}

type PaystackPWTMainRes struct {
	Reference     string `json:"reference"`
	Slug          string `json:"slug"`
	AccountName   string `json:"account_name"`
	BankName      string `json:"bank_name"`
	AccountNumber string `json:"account_number"`
	ExpiresAt     string `json:"expires_at"`
}

// End Pay With Transfer

//Pay With USSD

type PaystackUSSDType struct {
	Type string `json:"type"`
}

type PaystackUSSDFieldItem struct {
	Value        string `json:"value"`
	DisplayName  string `json:"display_name"`
	VariableName string `json:"variable_name"`
}

type PaystackUSSDMetaData struct {
	CustomFields []PaystackUSSDFieldItem `json:"custom_fields"`
}

type PaystackUSSDParams struct {
	Email     string               `json:"email"`
	Amount    string               `json:"amount"`
	Reference string               `json:"reference"`
	USSD      PaystackUSSDType     `json:"ussd"`
	MetaData  PaystackUSSDMetaData `json:"meta_data"`
}

type PaystackUSSDChargeAttempt struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Reference   string `json:"reference"`
		Status      string `json:"status"`
		DisplayText string `json:"display_text"`
		UssdCode    string `json:"ussd_code"`
	} `json:"data"`
}
type PaystackUSSDRes struct {
	Reference   string `json:"reference"`
	DisplayText string `json:"display_text"`
	USSDCode    string `json:"ussd_code"`
}

// End of USSD

// PAYSTACK CARD PAYMENT
type PaystackCardPaymentRes struct {
	Reference string `json:"reference" binding:"required"`
	Charge    int    `json:"charge" binding:"required"`
	Reason    string `json:"reason" binding:"required"`
	Email     string `json:"email" binding:"required"`
}

func ConvertStringToFloat(num string) float64 {
	result, err := strconv.ParseFloat(num, 64)
	if err != nil {
		log.Printf("Error at ConvertStringToFloat %v\n", err)
		return 0.00
	}
	return result
}

func ConvertToPaystackCharge(charge string) int {
	paystackCharge := ConvertStringToFloat(charge)
	return int(math.Ceil(paystackCharge * 100))
}

// END PAYSTACK CARD PAYMENT

// Flutterwave

type FlutterwaveMetaData struct {
	ConsumerID  int    `json:"consumer_id"`
	ConsumerMac string `json:"consumer_mac"`
}

type FlutterwaveCustomer struct {
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Name        string `json:"name"`
}

type FlutterwaveCustomization struct {
	Title string `json:"title"`
	Logo  string `json:"logo"`
}

type FlutterwaveDataLink struct {
	Link string `json:"link"`
}

type FlutterwaveRes struct {
	Status  string              `json:"status"`
	Message string              `json:"message"`
	Data    FlutterwaveDataLink `json:"data"`
}

type FlutterwaveChargeParams struct {
	TxRef         string                   `json:"tx_ref"`
	Amount        string                   `json:"amount"`
	Currency      string                   `json:"currency"`
	RedirectUrl   string                   `json:"redirect_url"`
	Meta          FlutterwaveMetaData      `json:"meta"`
	Customer      FlutterwaveCustomer      `json:"customer"`
	Customization FlutterwaveCustomization `json:"customization"`
}

// End Flutterwave
