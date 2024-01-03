package api

type InitPaymentParams struct {
	Reference      string `json:"reference" binding:"required"`
	MainOptionType string `json:"main_option_type" binding:"required"`
	Message        string `json:"message"`
}


type InitPaymentRes struct {
	Reference string `json:"reference" binding:"required"`
	Charge    int    `json:"charge" binding:"required"`
	Reason    string `json:"reason" binding:"required"`
	Email     string `json:"email" binding:"required"`
}




type VerifyPaymentReferenceParams struct {
	Reference      string `json:"reference" binding:"required"`
	MainOptionType string `json:"main_option_type" binding:"required"`
	Message        string `json:"message"`
}

type VerifyPaymentReferenceRes struct {
	Reference      string `json:"reference" binding:"required"`
	MainOptionType string `json:"main_option_type" binding:"required"`
	Success        bool `json:"success" binding:"required"`
}


