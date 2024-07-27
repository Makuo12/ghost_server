package payment

import "time"

type ReferencePayment struct {
	Reference string `json:"reference"`
	ChargeID  string `json:"charge_id"`
	Type      string `json:"type"`
	AddCard   bool   `json:"add_card"`
	Message   string `json:"message"`
}

type ReferencePaymentResponse struct {
	Status    string          `json:"status"`
	Bank      string          `json:"bank"`
	Verified  bool            `json:"verified"`
	Card      CardAddResponse `json:"card"`
	Email     string          `json:"email"`
	StartTime string          `json:"start_time"`
	Channel   string          `json:"channel"`
	Currency  string          `json:"currency"`
	Amount    string          `json:"amount"`
}

type PaystackVerifyResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    Data   `json:"data"`
}
type History struct {
	Type    string `json:"type"`
	Message string `json:"message"`
	Time    int    `json:"time"`
}
type Log struct {
	StartTime int       `json:"start_time"`
	TimeSpent int       `json:"time_spent"`
	Attempts  int       `json:"attempts"`
	Errors    int       `json:"errors"`
	Success   bool      `json:"success"`
	Mobile    bool      `json:"mobile"`
	Input     []any     `json:"input"`
	History   []History `json:"history"`
}
type Params struct {
	Bearer            string `json:"bearer"`
	TransactionCharge string `json:"transaction_charge"`
	PercentageCharge  string `json:"percentage_charge"`
}
type FeesSplit struct {
	Paystack    int    `json:"paystack"`
	Integration int    `json:"integration"`
	Subaccount  int    `json:"subaccount"`
	Params      Params `json:"params"`
}
type Authorization struct {
	AuthorizationCode string `json:"authorization_code"`
	Bin               string `json:"bin"`
	Last4             string `json:"last4"`
	ExpMonth          string `json:"exp_month"`
	ExpYear           string `json:"exp_year"`
	Channel           string `json:"channel"`
	CardType          string `json:"card_type"`
	Bank              string `json:"bank"`
	CountryCode       string `json:"country_code"`
	Brand             string `json:"brand"`
	Reusable          bool   `json:"reusable"`
	Signature         string `json:"signature"`
	AccountName       any    `json:"account_name"`
}
type Customer struct {
	ID           int    `json:"id"`
	FirstName    any    `json:"first_name"`
	LastName     any    `json:"last_name"`
	Email        string `json:"email"`
	CustomerCode string `json:"customer_code"`
	Phone        any    `json:"phone"`
	Metadata     any    `json:"metadata"`
	RiskAction   string `json:"risk_action"`
}
type PlanObject struct {
}
type Subaccount struct {
	ID                  int     `json:"id"`
	SubaccountCode      string  `json:"subaccount_code"`
	BusinessName        string  `json:"business_name"`
	Description         string  `json:"description"`
	PrimaryContactName  any     `json:"primary_contact_name"`
	PrimaryContactEmail any     `json:"primary_contact_email"`
	PrimaryContactPhone any     `json:"primary_contact_phone"`
	Metadata            any     `json:"metadata"`
	PercentageCharge    float64 `json:"percentage_charge"`
	SettlementBank      string  `json:"settlement_bank"`
	AccountNumber       string  `json:"account_number"`
}
type Data struct {
	ID              int           `json:"id"`
	Domain          string        `json:"domain"`
	Status          string        `json:"status"`
	Reference       string        `json:"reference"`
	Amount          int           `json:"amount"`
	Message         any           `json:"message"`
	GatewayResponse string        `json:"gateway_response"`
	PaidAt          time.Time     `json:"paid_at"`
	CreatedAt       time.Time     `json:"created_at"`
	Channel         string        `json:"channel"`
	Currency        string        `json:"currency"`
	IPAddress       string        `json:"ip_address"`
	Metadata        any           `json:"metadata"`
	Log             Log           `json:"log"`
	Fees            int           `json:"fees"`
	FeesSplit       FeesSplit     `json:"fees_split"`
	Authorization   Authorization `json:"authorization"`
	Customer        Customer      `json:"customer"`
	Plan            any           `json:"plan"`
	OrderID         any           `json:"order_id"`
	PaidAt0         time.Time     `json:"paidAt"`
	CreatedAt0      time.Time     `json:"createdAt"`
	RequestedAmount int           `json:"requested_amount"`
	TransactionDate time.Time     `json:"transaction_date"`
	PlanObject      PlanObject    `json:"plan_object"`
	Subaccount      Subaccount    `json:"subaccount"`
}

type InitCardChargeParams struct {
	Currency        string `json:"currency"`
	CardLast4       string `json:"card_last4"`
	ObjectReference string `json:"object_reference"`
}

type InitRemoveCardParams struct {
	ID string `json:"id"`
}

type InitRemoveCardRes struct {
	Success bool   `json:"success"`
	ID      string `json:"id"`
}

type SetDefaultCardParams struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

type SetDefaultCardRes struct {
	Success bool   `json:"success"`
	ID      string `json:"id"`
	Type    string `json:"type"`
}

type InitCardChargeRes struct {
	Reference string `json:"reference"`
	Reason    string `json:"reason"`
	Charge    int    `json:"charge"`
	Currency  string `json:"currency"`
	Email     string `json:"email"`
}

type CardDetailResponse struct {
	CardID    string `json:"card_id"`
	CardLast4 string `json:"card_last4"`
	CardType  string `json:"card_type"`
	ExpMonth  string `json:"exp_month"`
	ExpYear   string `json:"exp_year"`
	Currency  string `json:"currency"`
}

type CardAddResponse struct {
	CardDetail      CardDetailResponse `json:"card_detail"`
	Account         string             `json:"account"`
	AccountCurrency string             `json:"account_currency"`
	DefaultID       string             `json:"default_id"`
}

type GetWalletResponse struct {
	Cards           []CardDetailResponse `json:"cards"`
	USDAccount      string               `json:"usd_account"`
	NGNAccount      string               `json:"ngn_account"`
	DefaultID       string               `json:"default_id"`
	DefaultPayoutID string               `json:"default_payout_id"`
	HasCard         bool                 `json:"has_card"`
}

type PaystackAuthorization struct {
	AuthorizationCode string                  `json:"authorization_code"`
	Email             string                  `json:"email"`
	Amount            int                     `json:"amount"`
	CallbackUrl       string                  `json:"callback_url"`
	MetaData          PaystackPaymentMetaData `json:"meta_data"`
}

type PaystackPaymentMetaData struct {
	CancelAction string `json:"cancel_action"`
}

type PaystackPaymentResponse struct {
	Status  bool       `json:"status"`
	Message string     `json:"message"`
	Data    DataCharge `json:"data"`
}

type DataChallengeResponse struct {
	AuthorizationUrl string `json:"authorization_url"`
	Reference        string `json:"reference"`
	AccessCode       string `json:"access_code"`
	Paused           bool   `json:"paused"`
}

type PaymentPaymentChallengeResponse struct {
	Status  bool                  `json:"status"`
	Message string                `json:"message"`
	Data    DataChallengeResponse `json:"data"`
}

type DataCharge struct {
	Amount          int                 `json:"amount"`
	Currency        string              `json:"currency"`
	TransactionDate string              `json:"transaction_date"`
	Status          string              `json:"status"`
	Reference       string              `json:"reference"`
	Domain          string              `json:"domain"`
	Metadata        any                 `json:"metadata"` // Use any for handling null
	GatewayResponse string              `json:"gateway_response"`
	Message         any                 `json:"message"` // Use any for handling null
	Channel         string              `json:"channel"`
	IPAddress       any                 `json:"ip_address"` // Use any for handling null
	Log             any                 `json:"log"`        // Use any for handling null
	Fees            int                 `json:"fees"`
	Authorization   AuthorizationCharge `json:"authorization"`
	Customer        CustomerCharge      `json:"customer"`
	Plan            any                 `json:"plan"` // Use any for handling null
	ID              int                 `json:"id"`
}

type AuthorizationCharge struct {
	AuthorizationCode string `json:"authorization_code"`
	BIN               string `json:"bin"`
	Last4             string `json:"last4"`
	ExpMonth          string `json:"exp_month"`
	ExpYear           string `json:"exp_year"`
	Channel           string `json:"channel"`
	CardType          string `json:"card_type"`
	Bank              string `json:"bank"`
	CountryCode       string `json:"country_code"`
	Brand             string `json:"brand"`
	Reusable          bool   `json:"reusable"`
	Signature         string `json:"signature"`
	AccountName       any    `json:"account_name"` // Use any for handling null
}

type CustomerCharge struct {
	ID           int    `json:"id"`
	FirstName    any    `json:"first_name"` // Use any for handling null
	LastName     any    `json:"last_name"`  // Use any for handling null
	Email        string `json:"email"`
	CustomerCode string `json:"customer_code"`
	Phone        any    `json:"phone"`    // Use any for handling null
	Metadata     any    `json:"metadata"` // Use any for handling null
	RiskAction   string `json:"risk_action"`
}
