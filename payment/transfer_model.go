package payment

type ChargeReferenceItem struct {
	ChargeID       string `json:"charge_id"`
	Reason         string `json:"reason"`
	IsComplete     bool   `json:"is_complete"`
	Charge         string `json:"charge"`
	Currency       string `json:"currency"`
	PaymentChannel string `json:"payment_channel"`
	CreatedAt      string `json:"created_at"`
}

func GetChargeReferenceReason(reason string) string {
	switch reason {
	case "USER_OPTION_PAYMENT":
		return "Payment for a stay"
	case "add_card_reason":
		return "Add a card"
	}
	return ""
}

func GetChargeReferencePaymentChannel(reason string) string {
	switch reason {
	case "paystack_card":
		return "Card"
	}
	return "Transfer"
}

type ChargeReferenceRes struct {
	List []ChargeReferenceItem `json:"list"`
}

type ChargeReferenceParams struct {
	Offset int `json:"offset"`
}
