package payment


type WebhookEvent struct {
	Event string `json:"event"`
	Data  any    `json:"data"`
}



type RefundEvent struct {
	Event string `json:"event"`
	Data  struct {
		Status               string `json:"status"`
		TransactionReference string `json:"transaction_reference"`
		RefundReference      any    `json:"refund_reference"`
		Amount               int    `json:"amount"`
		Currency             string `json:"currency"`
		Processor            string `json:"processor"`
		Customer             struct {
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
			Email     string `json:"email"`
		} `json:"customer"`
		Integration int    `json:"integration"`
		Domain      string `json:"domain"`
	} `json:"data"`
}


type PaystackChargeEventWebhookEvent struct {
	Event string                  `json:"event"`
	Data  PaystackChargeEventData `json:"data"`
}

type PaystackChargeEventData struct {
	ID              int                              `json:"id"`
	Domain          string                           `json:"domain"`
	Status          string                           `json:"status"`
	Reference       string                           `json:"reference"`
	Amount          int                              `json:"amount"`
	Message         any                              `json:"message"`
	GatewayResponse string                           `json:"gateway_response"`
	PaidAt          any                              `json:"paid_at"`
	CreatedAt       any                              `json:"created_at"`
	Channel         string                           `json:"channel"`
	Currency        string                           `json:"currency"`
	IPAddress       string                           `json:"ip_address"`
	Metadata        any                              `json:"metadata"`
	Log             PaystackChargeEventLog           `json:"log"`
	Fees            any                              `json:"fees"`
	Customer        PaystackChargeEventCustomer      `json:"customer"`
	Authorization   PaystackChargeEventAuthorization `json:"authorization"`
	Plan            any                              `json:"plan"`
}

type PaystackChargeEventLog struct {
	TimeSpent      int                          `json:"time_spent"`
	Attempts       int                          `json:"attempts"`
	Authentication string                       `json:"authentication"`
	Errors         int                          `json:"errors"`
	Success        bool                         `json:"success"`
	Mobile         bool                         `json:"mobile"`
	Input          []string                     `json:"input"`
	Channel        any                          `json:"channel"`
	History        []PaystackChargeEventHistory `json:"history"`
}

type PaystackChargeEventHistory struct {
	Type    string `json:"type"`
	Message string `json:"message"`
	Time    int    `json:"time"`
}

type PaystackChargeEventCustomer struct {
	ID           int    `json:"id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Email        string `json:"email"`
	CustomerCode string `json:"customer_code"`
	Phone        any    `json:"phone"`
	Metadata     any    `json:"metadata"`
	RiskAction   string `json:"risk_action"`
}

type PaystackChargeEventAuthorization struct {
	AuthorizationCode string `json:"authorization_code"`
	Bin               string `json:"bin"`
	Last4             string `json:"last4"`
	ExpMonth          string `json:"exp_month"`
	ExpYear           string `json:"exp_year"`
	CardType          string `json:"card_type"`
	Bank              string `json:"bank"`
	CountryCode       string `json:"country_code"`
	Brand             string `json:"brand"`
	AccountName       string `json:"account_name"`
}
