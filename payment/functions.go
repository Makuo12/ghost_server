package payment

func EmptyPaystackBankCharge() PaystackBankAccountMainRes {
	return PaystackBankAccountMainRes{
		Reference:   "none",
		DisplayText: "none",
	}
}

func EmptyPaystackPWT() PaystackPWTMainRes {
	return PaystackPWTMainRes{
		Reference:     "none",
		Slug:          "none",
		AccountName:   "none",
		AccountNumber: "none",
		ExpiresAt:     "none",
	}
}

func EmptyPaystackUSSD() PaystackUSSDRes {
	return PaystackUSSDRes{
		Reference:   "none",
		DisplayText: "none",
		USSDCode:    "none",
	}
}

func EmptyPaystackCard() InitCardChargeRes {
	return InitCardChargeRes{
		Reference: "none",
		Reason:    "none",
		Charge:    0,
		Currency:  "none",
		Email:     "none",
	}
}
