package api

type CreateUserOptionCancellationParams struct {
	ChargeID  string `json:"charge_id" binding:"required"`
	ReasonTwo string `json:"reason_two"`
	ReasonOne string `json:"reason_one" binding:"required,user_option_cancel"`
	Message   string `json:"message" binding:"required"`
}

type CreateUserEventCancellationParams struct {
	ChargeID  string `json:"charge_id" binding:"required"`
	ReasonTwo string `json:"reason_two"`
	ReasonOne string `json:"reason_one" binding:"required,user_event_cancel"`
	Message   string `json:"message" binding:"required"`
}

type CreateHostOptionCancellationParams struct {
	ChargeID  string `json:"charge_id" binding:"required"`
	ReasonTwo string `json:"reason_two"`
	ReasonOne string `json:"reason_one" binding:"required,host_option_cancel"`
	Message   string `json:"message" binding:"required"`
	UserID    string `json:"user_id" binding:"required"`
	OptionID  string `json:"option_id" binding:"required"`
}

type CreateCancellationRes struct {
	ChargeID string `json:"charge_id" binding:"required"`
	Msg      string `json:"msg" binding:"required"`
}

type CreateHostEventCancellationParams struct {
	StartDate   string `json:"start_date" binding:"required"`
	EndDate     string `json:"end_date" binding:"required"`
	EventID     string `json:"event_id" binding:"required"`
	EventDateID string `json:"event_date_id" binding:"required"`
	ReasonTwo   string `json:"reason_two"`
	ReasonOne   string `json:"reason_one" binding:"required,host_event_cancel"`
	Message     string `json:"message" binding:"required"`
}

type CreateHostEventCancellationRes struct {
	StartDate   string `json:"start_date" binding:"required"`
	EndDate     string `json:"end_date" binding:"required"`
	EventID     string `json:"event_id" binding:"required"`
	EventDateID string `json:"event_date_id" binding:"required"`
	Msg         string `json:"msg" binding:"required"`
}

// Details

type CancelUserDetailParams struct {
	ChargeID string `json:"charge_id" binding:"required"`
}

type CancelUserOptionDetailRes struct {
	MainPrice     string `json:"main_price" binding:"required"`
	CleaningFee   string `json:"cleaning_fee" binding:"required"`
	ServiceFee    string `json:"service_fee" binding:"required"`
	TotalFee      string `json:"total_fee" binding:"required"`
	GuestFee      string `json:"guest_fee" binding:"required"`
	PetFee        string `json:"pet_fee" binding:"required"`
	RefundPercent int    `json:"refund_percent" binding:"required"`
	Refund        string `json:"refund" binding:"required"`
	RefundType    string `json:"refund_type" binding:"required"`
	DateBooked    string `json:"date_booked" binding:"required"`
	Currency      string `json:"currency" binding:"required"`
}

type CancelUserTicketDetailRes struct {
	TicketPrice   string `json:"ticket_price" binding:"required"`
	ServiceFee    string `json:"service_fee" binding:"required"`
	RefundPercent int    `json:"refund_percent" binding:"required"`
	Refund        string `json:"refund" binding:"required"`
	RefundType    string `json:"refund_type" binding:"required"`
	DateBooked    string `json:"date_booked" binding:"required"`
	Currency      string `json:"currency" binding:"required"`
}

type CancelHostOptionDetailParams struct {
	ChargeID string `json:"charge_id" binding:"required"`
	UserID   string `json:"user_id" binding:"required"`
	OptionID string `json:"option_id" binding:"required"`
}

type CancelHostEventDetailParams struct {
	StartDate   string `json:"start_date" binding:"required"`
	EndDate     string `json:"end_date" binding:"required"`
	EventID     string `json:"event_id" binding:"required"`
	EventDateID string `json:"event_date_id" binding:"required"`
}

type CancelHostDetailRes struct {
	CanCancel  bool   `json:"can_cancel" binding:"required"`
	HostPayout int    `json:"host_payout" binding:"required"`
	Amount     string `json:"amount" binding:"required"`
}
