package api

import "github.com/google/uuid"

type GetTicketByIDAndOptionIDRow struct {
	StartDate    string    `json:"start_date"`
	TicketID     uuid.UUID `json:"ticket_id"`
	Level        string    `json:"level"`
	Price        int64     `json:"price"`
	TicketType   string    `json:"ticket_type"`
	TicketName   string    `json:"ticket_name"`
	EventDateID  uuid.UUID `json:"event_date_id"`
	OptionUserID uuid.UUID `json:"option_user_id"`
	Currency     string    `json:"currency"`
	StartTime    string    `json:"start_time"`
	EndTime      string    `json:"end_time"`
	TimeZone     string    `json:"time_zone"`
	PayType      string    `json:"pay_type"`
	AbsorbFees   bool      `json:"absorb_fees"`
}

type TicketReserveItemDB struct {
	ID         string  `json:"id"`
	Grade      string  `json:"grade"`
	Price      float64 `json:"price"`
	ServiceFee float64 `json:"service_fee"`
	PayType    string  `json:"pay_type"`
	AbsorbFees float64 `json:"absorb_fees"`
	Type       string  `json:"type"`
	GroupPrice float64 `json:"group_price"`
}

type DateReserveItemDB struct {
	ID                  string                `json:"id"`
	StartDate           string                `json:"start_date"`
	EndDate             string                `json:"end_date"`
	StartTime           string                `json:"start_time"`
	EndTime             string                `json:"end_time"`
	TimeZone            string                `json:"time_zone"`
	TotalDateServiceFee float64               `json:"total_date_service_fee"`
	TotalDateAbsorbFee  float64               `json:"total_date_absorb_fee"`
	TotalDateFee        float64               `json:"total_date_fee"`
	Tickets             []TicketReserveItemDB `json:"tickets"`
}

type EventDateReserveDB struct {
	ID              string              `json:"id"`
	DateTimes       []DateReserveItemDB `json:"date_times"`
	Currency        string              `json:"currency"`
	TotalFee        float64             `json:"total_fee"`
	TotalServiceFee float64             `json:"total_service_fee"`
	TotalAbsorbFee  float64             `json:"total_absorb_fee"`
}

type TicketReserveItem struct {
	ID         string `json:"id"`
	Grade      string `json:"grade"`
	Price      string `json:"price"`
	ServiceFee string `json:"service_fee"`
	AbsorbFee  string `json:"absorb_fee"`
	Type       string `json:"type"`
	PayType    string `json:"pay_type"`
	GroupPrice string `json:"group_price"`
}

type DateReserveItem struct {
	ID                  string              `json:"id"`
	StartDate           string              `json:"start_date"`
	EndDate             string              `json:"end_date"`
	StartTime           string              `json:"start_time"`
	EndTime             string              `json:"end_time"`
	TimeZone            string              `json:"time_zone"`
	TotalDateFee        string              `json:"total_date_fee"`
	TotalDateServiceFee string              `json:"total_date_service_fee"`
	TotalDateAbsorbFee  string              `json:"total_date_absorb_fee"`
	Tickets             []TicketReserveItem `json:"tickets"`
}

type EventDateReserve struct {
	ID              string            `json:"id"`
	DateTimes       []DateReserveItem `json:"date_times"`
	Currency        string            `json:"currency"`
	TotalFee        string            `json:"total_fee"`
	TotalServiceFee string            `json:"total_service_fee"`
	TotalAbsorbFee  string            `json:"total_absorb_fee"`
}

type CreateEventReserveDetailRes struct {
	ReserveData    EventDateReserve   `json:"reserve_data"`
	EventReference string             `json:"event_reference"`
	DefaultCardID  string             `json:"default_card_id"`
	HasCard        bool               `json:"has_card"`
	CardDetail     CardDetailResponse `json:"card_detail"`
}

type ReserveEventParams struct {
	OptionUserID string       `json:"option_user_id" binding:"required"`
	UserCurrency string       `json:"user_currency" binding:"required,currency"`
	Tickets      []TicketItem `json:"ticket_ids" binding:"required"`
}

type TicketItem struct {
	EventDateID string `json:"event_date_id" binding:"required"`
	StartDate   string `json:"start_date" binding:"required"`
	EndDate     string `json:"end_date" binding:"required"`
	TicketID    string `json:"ticket_id" binding:"required"`
}
