package api

import "github.com/makuo12/ghost_server/payment"

type DatePrice struct {
	Price      string `json:"price"`
	Date       string `json:"date"`
	GroupPrice string `json:"group_price"`
}

type DatePriceFloat struct {
	Price      float64 `json:"price"`
	Date       string  `json:"date"`
	GroupPrice float64 `json:"group_price"`
}

type ReserveOptionParams struct {
	StartDate    string   `json:"start_date" binding:"required,date_only"`
	EndDate      string   `json:"end_date" binding:"required,date_only"`
	Guests       []string `json:"guests" binding:"required,guest_option"`
	OptionUserID string   `json:"option_user_id" binding:"required"`
	UserCurrency string   `json:"user_currency" binding:"required,currency"`
}

type ReDiscount struct {
	Price string `json:"price"`
	Type  string `json:"type"`
}

type ExperienceReserveOModel struct {
	Discount        ReDiscount  `json:"discount"`
	MainPrice       string      `json:"main_price"`
	ServiceFee      string      `json:"service_fee"`
	TotalFee        string      `json:"total_fee"`
	DatePrice       []DatePrice `json:"date_price"`
	Currency        string      `json:"currency"`
	Guests          []string    `json:"guests"`
	GuestFee        string      `json:"guest_fee"`
	PetFee          string      `json:"pet_fee"`
	CleaningFee     string      `json:"cleaning_fee"`
	NightlyPetFee   string      `json:"nightly_pet_fee"`
	NightlyGuestFee string      `json:"nightly_guest_fee"`
	CanInstantBook  bool        `json:"can_instant_book"`
	RequireRequest  bool        `json:"require_request"`
	RequestType     string      `json:"request_type"`
	Reference       string      `json:"reference"`
	OptionUserID    string      `json:"option_user_id"`
	StartDate       string      `json:"start_date"`
	EndDate         string      `json:"end_date"`
}

type CreateOptionReserveDetailRes struct {
	ReserveData   ExperienceReserveOModel    `json:"reserve_data"`
	DefaultCardID string                     `json:"default_card_id"`
	HasCard       bool                       `json:"has_card"`
	CardDetail    payment.CardDetailResponse `json:"card_detail"`
}

type FinalOptionReserveDetailParams struct {
	Reference string `json:"reference" binding:"required"`
	// ID  is the card ID we would use to get the card
	ID      string `json:"id"`
	Message string `json:"message"`
}

type FinalOptionReserveDetailVerificationParams struct {
	Reference        string `json:"reference" binding:"required"`
	Successful       bool   `json:"successful"`
	PaymentReference string `json:"payment_reference" binding:"required"`
	Message          string `json:"message"`
}

type FinalOptionReserveDetailRes struct {
	Reference         string `json:"reference"`
	AuthorizationUrl  string `json:"authorization_url"`
	AccessCode        string `json:"access_code"`
	PaymentReference  string `json:"payment_reference"`
	Paused            bool   `json:"paused"`
	PaymentSuccess    bool   `json:"payment_success"`
	PaymentChallenged bool   `json:"payment_challenged"`
	SuccessUrl        string `json:"success_url"`
	FailureUrl        string `json:"failure_url"`
}

type FinalOptionReserveRequestDetailRes struct {
	Reference   string `json:"reference"`
	Message     string `json:"message"`
	RequestSent bool   `json:"request_sent"`
}

type MsgRequestResponseParams struct {
	MsgID    string `json:"msg_id"`
	Approved bool   `json:"approved"`
	Message  string `json:"message"`
}

type InitMethodPaymentParams struct {
	Reference           string                               `json:"reference"`
	PaymentType         string                               `json:"payment_type" binding:"required"`
	MainOptionType      string                               `json:"main_option_type" binding:"required"`
	PaymentMethod       string                               `json:"payment_method"`
	PaymentChannel      string                               `json:"payment_channel"`
	Message             string                               `json:"message"`
	Type                string                               `json:"type"`
	Currency            string                               `json:"currency"`
	CardLast4           string                               `json:"card_last4"`
	PaystackBankAccount payment.PaystackGetBankAccountParams `json:"paystack_bank_account"`
	PaystackUSSD        string                               `json:"paystack_ussd"`
}

type InitMethodPaymentRes struct {
	Reference          string                             `json:"reference" binding:"required"`
	PaymentType        string                             `json:"payment_type" binding:"required"`
	MainOptionType     string                             `json:"main_option_type" binding:"required"`
	PaymentMethod      string                             `json:"payment_method"`
	PaymentReference   string                             `json:"payment_reference"`
	PaymentChannel     string                             `json:"payment_channel"`
	PaystackBankCharge payment.PaystackBankAccountMainRes `json:"payment_bank_charge"`
	PaystackPWT        payment.PaystackPWTMainRes         `json:"payment_pwt"`
	PaystackCard       payment.InitCardChargeRes          `json:"payment_card"`
	PaystackUSSD       payment.PaystackUSSDRes            `json:"payment_ussd"`
}
