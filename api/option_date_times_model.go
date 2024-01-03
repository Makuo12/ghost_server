package api

import (
	"time"

	"github.com/google/uuid"
)

type OptionDateItem struct {
	ID        string `json:"id"`
	OptionID  string `json:"option_id"`
	Date      string `json:"date"`
	Available bool   `json:"available"`
	Price     string `json:"price"`
	IsEmpty   bool   `json:"is_empty"`
}

type ListOptionDateItem struct {
	List         []OptionDateItem       `json:"list"`
	ListBooked   []DateOptionBookedItem `json:"list_booked"`
	BasePrice    string                 `json:"base_price"`
	WeekendPrice string                 `json:"weekend_price"`
}

type CreateUpdateOptionDateTimeParams struct {
	Dates     []string `json:"dates"`
	OptionID  string   `json:"option_id"`
	Available bool     `json:"available"`
	Price     string   `json:"price"`
	Note      string   `json:"note"`
}

type OptionDateParams struct {
	OptionID string `json:"option_id" binding:"required"`
	Date     string `json:"date" binding:"required"`
}

type OptionSelectedDateParams struct {
	OptionDateTimeID string `json:"option_date_time_id" binding:"required"`
	OptionID         string `json:"option_id" binding:"required"`
}

type GetOptionDateNoteRes struct {
	OptionDateTimeID string `json:"option_date_time_id" binding:"required"`
	Note             string `json:"note" binding:"required"`
}

// OptionDateTime is important for reservation
type OptionDateTime struct {
	ID        uuid.UUID `json:"id"`
	OptionID  uuid.UUID `json:"option_id"`
	Date      string    `json:"date"`
	Available bool      `json:"available"`
	Price     string    `json:"price"`
	Note      string    `json:"note"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
