package api

type CreateEventDateTimeParams struct {
	EventInfoID string   `json:"event_info_id" binding:"required"`
	StartDate   string   `json:"start_date"`
	EndDate     string   `json:"end_date"`
	Type        string   `json:"type"`
	EventDates  []string `json:"event_dates"`
	TimeZone    string   `json:"time_zone"`
}

type UpdateEventDateTimeNoteParams struct {
	EventInfoID     string `json:"event_info_id" binding:"required"`
	EventDateTimeID string `json:"event_date_time_id" binding:"required"`
	Note            string `json:"note" binding:"required"`
}

type CreateEventDateTimeRes struct {
	EventDateTimeID  string   `json:"event_date_time_id"`
	StartDate        string   `json:"start_date" `
	EndDate          string   `json:"end_date"`
	Status           string   `json:"status"`
	Type             string   `json:"type"`
	EventDates       []string `json:"event_dates"`
	NeedBands        bool     `json:"need_bands"`
	NeedTickets      bool     `json:"need_tickets"`
	AbsorbBandCharge bool     `json:"absorb_band_charge"`
}

type UpdateEventDateTimeParams struct {
	EventInfoID     string   `json:"event_info_id" binding:"required"`
	EventDateTimeID string   `json:"event_date_time_id" binding:"required"`
	StartDate       string   `json:"start_date"`
	EndDate         string   `json:"end_date"`
	Type            string   `json:"type"`
	EventDates      []string `json:"event_dates"`
}

type UpdateEventDateTimeControlParams struct {
	EventInfoID      string `json:"event_info_id" binding:"required"`
	EventDateTimeID  string `json:"event_date_time_id" binding:"required"`
	NeedTickets      bool   `json:"need_Tickets"`
	NeedBands        bool   `json:"need_bands"`
	AbsorbBandCharge bool   `json:"absorb_band_charge"`
	Type             string `json:"type"`
}

type UpdateEventDateTimeControlRes struct {
	EventDateTimeID  string `json:"event_date_time_id"`
	NeedTickets      bool   `json:"need_Tickets"`
	NeedBands        bool   `json:"need_bands"`
	AbsorbBandCharge bool   `json:"absorb_band_charge"`
}

type EventDateItem struct {
	ID               string   `json:"id"`
	Name             string   `json:"name"`
	StartTime        string   `json:"start_time"`
	EndTime          string   `json:"end_time"`
	Status           string   `json:"status"`
	StartDate        string   `json:"start_date"`
	EndDate          string   `json:"end_date"`
	Tickets          int      `json:"tickets"`
	Note             string   `json:"note"`
	TimeZone         string   `json:"time_zone"`
	EventDates       []string `json:"event_dates"`
	Type             string   `json:"type"`
	NeedBands        bool     `json:"need_bands"`
	NeedTickets      bool     `json:"need_tickets"`
	AbsorbBandCharge bool     `json:"absorb_band_charge"`
}

type EventDateNormalItem struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	StartDate  string   `json:"start_date"`
	EndDate    string   `json:"end_date"`
	Type       string   `json:"type"`
	EventDates []string `json:"event_dates"`
}

type EventDateEDTItem struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Status    string `json:"status"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	Note      string `json:"note"`
	TimeZone  string `json:"time_zone"`
}

type ListEventDateNormalItem struct {
	List        []EventDateNormalItem `json:"list" `
	ItemOffset  int                   `json:"item_offset"`
	OnLastIndex bool                  `json:"on_last_index" binding:"required"`
}

type ListEventDateItem struct {
	List        []EventDateItem     `json:"list"`
	IsEmpty     bool                `json:"is_empty" binding:"required"`
	OptionData  GetUHMDataOptionRes `json:"option_data"`
	ItemOffset  int                 `json:"item_offset"`
	OnLastIndex bool                `json:"on_last_index" binding:"required"`
}

type EventDateParams struct {
	ItemOffset  int    `json:"item_offset"`
	EventInfoID string `json:"event_info_id" binding:"required"`
}

type CreateUpdateEventDateDetailReq struct {
	EventInfoID     string `json:"event_info_id"`
	EventDateTimeID string `json:"event_date_time_id" binding:"required"`
	StartTime       string `json:"start_time" binding:"required,time_only"`
	EndTime         string `json:"end_time" binding:"required,time_only"`
	Name            string `json:"name" binding:"required"`
	TimeZone        string `json:"time_zone" binding:"required"`
	Type            string `json:"type" binding:"required"`
}

type UpdateEventDateStatusReq struct {
	EventInfoID     string `json:"event_info_id"`
	EventDateTimeID string `json:"event_date_time_id" binding:"required"`
	Status          string `json:"status" binding:"required"`
	Type            string `json:"type" binding:"required"`
}

type UpdatePublishEventCheckInStepParams struct {
	EventInfoID     string `json:"event_info_id"`
	EventDateTimeID string `json:"event_date_time_id" binding:"required"`
}

type UpdatePublishEventCheckInStepRes struct {
	Published bool `json:"published" binding:"required"`
}

type UpdateEventDateStatusRes struct {
	Status string `json:"status" binding:"required"`
}

type CreateUpdateEventDateDetailRes struct {
	EventInfoID     string `json:"event_info_id"`
	EventDateTimeID string `json:"event_date_time_id" binding:"required"`
	StartTime       string `json:"start_time" binding:"required,time_only"`
	EndTime         string `json:"end_time" binding:"required,time_only"`
	Name            string `json:"name" binding:"required"`
	TimeZone        string `json:"time_zone" binding:"required"`
}

type CreateUpdateEventDateLocationReq struct {
	EventInfoID     string `json:"event_info_id"`
	EventDateTimeID string `json:"event_date_time_id" binding:"required"`
	Street          string `json:"street" binding:"required"`
	City            string `json:"city" binding:"required"`
	State           string `json:"state" binding:"required"`
	Country         string `json:"country" binding:"required"`
	Postcode        string `json:"postcode" binding:"required"`
	Lat             string `json:"lat" binding:"required"`
	Lng             string `json:"lng" binding:"required"`
}

type CreateUpdateEventDateLocationRes struct {
	EventDateTimeID string `json:"event_date_time_id" binding:"required"`
	Street          string `json:"street" binding:"required"`
	City            string `json:"city" binding:"required"`
	State           string `json:"state" binding:"required"`
	Country         string `json:"country" binding:"required"`
	Postcode        string `json:"postcode" binding:"required"`
	Lat             string `json:"lat" binding:"required"`
	Lng             string `json:"lng" binding:"required"`
}

type CreateEventDateTicketParams struct {
	EventInfoID     string `json:"event_info_id" binding:"required"`
	EventDateTimeID string `json:"event_date_time_id" binding:"required"`
	StartDate       string `json:"start_date" binding:"required,date_only"`
	EndDate         string `json:"end_date" binding:"required,date_only"`
	StartTime       string `json:"start_time" binding:"required,time_only"`
	EndTime         string `json:"end_time" binding:"required,time_only"`
	Name            string `json:"name" binding:"required"`
	Capacity        int    `json:"capacity"`
	Price           string `json:"price" binding:"required"`
	AbsorbFees      bool   `json:"absorb_fees"`
	Description     string `json:"description" binding:"required"`
	Type            string `json:"type" binding:"required,event_ticket_type"`
	Level           string `json:"level" binding:"required,event_ticket_level"`
	TicketType      string `json:"ticket_type" binding:"required,event_ticket_main_type"`
	NumOfSeats      int    `json:"num_of_seats"`
	FreeRefreshment bool   `json:"free_refreshment"`
}

type RemoveEventDateTicketParams struct {
	EventInfoID     string `json:"event_info_id" binding:"required"`
	EventDateTimeID string `json:"event_date_time_id" binding:"required"`
	TicketID        string `json:"ticket_id" binding:"required"`
}

type RemoveEventDateTicketRes struct {
	EventDateTimeID string `json:"event_date_time_id" binding:"required"`
	TicketID        string `json:"ticket_id" binding:"required"`
	TicketCount     int    `json:"ticket_count" binding:"required"`
}

type RemoveEventDateTimeParams struct {
	EventInfoID     string `json:"event_info_id" binding:"required"`
	EventDateTimeID string `json:"event_date_time_id" binding:"required"`
}

type RemoveEventDateTimeRes struct {
	EventDateTimeID string `json:"event_date_time_id" binding:"required"`
}

type UpdateEventDateTicketParams struct {
	EventInfoID     string `json:"event_info_id" binding:"required"`
	EventDateTimeID string `json:"event_date_time_id" binding:"required"`
	TicketID        string `json:"ticket_id" binding:"required"`
	StartDate       string `json:"start_date" binding:"required,date_only"`
	EndDate         string `json:"end_date" binding:"required,date_only"`
	StartTime       string `json:"start_time" binding:"required,time_only"`
	EndTime         string `json:"end_time" binding:"required,time_only"`
	Name            string `json:"name" binding:"required"`
	Capacity        int    `json:"capacity"`
	Price           string `json:"price" binding:"required"`
	AbsorbFees      bool   `json:"absorb_fees"`
	Description     string `json:"description" binding:"required"`
	Type            string `json:"type" binding:"required,event_ticket_type"`
	Level           string `json:"level" binding:"required,event_ticket_level"`
	TicketType      string `json:"ticket_type" binding:"required,event_ticket_main_type"`
	NumOfSeats      int    `json:"num_of_seats"`
	FreeRefreshment bool   `json:"free_refreshment"`
}

type CreateEventDateTicketRes struct {
	EventDateTimeID string `json:"event_date_time_id" binding:"required"`
	ID              string `json:"id" binding:"required"`
	StartDate       string `json:"start_date" binding:"required,date_only"`
	EndDate         string `json:"end_date" binding:"required,date_only"`
	StartTime       string `json:"start_time" binding:"required,time_only"`
	EndTime         string `json:"end_time" binding:"required,time_only"`
	Name            string `json:"name" binding:"required"`
	Capacity        int    `json:"capacity"`
	Price           string `json:"price" binding:"required"`
	AbsorbFees      bool   `json:"absorb_fees" binding:"required"`
	Description     string `json:"description" binding:"required"`
	Type            string `json:"type" binding:"required,event_ticket_type"`
	Level           string `json:"level" binding:"required,event_ticket_level"`
	TicketType      string `json:"ticket_type" binding:"required,event_ticket_main_type"`
	NumOfSeats      int    `json:"num_of_seats"`
	FreeRefreshment bool   `json:"free_refreshment"`
	TicketCount     int    `json:"ticket_count"`
}

type ListEventDateTicketParams struct {
	EventInfoID     string `json:"event_info_id" binding:"required"`
	EventDateTimeID string `json:"event_date_time_id" binding:"required"`
	OptionOffset    int    `json:"option_offset"`
}

type ListEventDateTicketRes struct {
	List         []CreateEventDateTicketRes `json:"list" binding:"required"`
	Grades       []string                   `json:"grades" binding:"required"`
	OptionOffset int                        `json:"option_offset"`
	OnLastIndex  bool                       `json:"on_last_index" binding:"required"`
}

type GetEventDateLocationReq struct {
	EventInfoID     string `json:"event_info_id" binding:"required"`
	EventDateTimeID string `json:"event_date_time_id" binding:"required"`
}

type GetEventDateLocationRes struct {
	EventDateTimeID string `json:"event_date_time_id" binding:"required"`
}

type GetEventDateDetailParams struct {
	EventInfoID     string `json:"event_info_id" binding:"required"`
	EventDateTimeID string `json:"event_date_time_id" binding:"required"`
}

type GetOptionDateIsBookedRes struct {
	IsBooked bool `json:"is_booked" binding:"required"`
}

type GetEventDateIsBookedRes struct {
	IsBooked bool `json:"is_booked" binding:"required"`
}

type PrivateAudienceItem struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Type   string `json:"type"`
	Email  string `json:"email"`
	Number string `json:"number"`
	Sent   bool   `json:"sent"`
	Exist  bool   `json:"exist"`
}

type UpdatePrivateAudienceParams struct {
	EventInfoID     string `json:"event_info_id" binding:"required"`
	EventDateTimeID string `json:"event_date_time_id" binding:"required"`
	ID              string `json:"id" binding:"required"`
	Email           string `json:"email"`
	Number          string `json:"number"`
}

type RemovePrivateAudienceParams struct {
	EventInfoID     string `json:"event_info_id"`
	EventDateTimeID string `json:"event_date_time_id" binding:"required"`
	ID              string `json:"id" binding:"required"`
}

type UpdateEventDatePublishParams struct {
	EventInfoID          string                `json:"event_info_id"`
	EventDateTimeID      string                `json:"event_date_time_id" binding:"required"`
	EventGoingPublicDate string                `json:"event_going_public_date"`
	EventGoingPublicTime string                `json:"event_going_public_time"`
	EventPublic          string                `json:"event_public"`
	EventGoingPublic     string                `json:"event_going_public"`
	Audiences            []PrivateAudienceItem `json:"audiences"`
	HasAudience          bool                  `json:"has_audience"`
}

type GetEventDatePublishParams struct {
	EventInfoID     string `json:"event_info_id" binding:"required"`
	EventDateTimeID string `json:"event_date_time_id" binding:"required"`
}

type GetEventDateTimeParams struct {
	EventInfoID     string `json:"event_info_id" binding:"required"`
	EventDateTimeID string `json:"event_date_time_id" binding:"required"`
}

type UpdateEventCheckInStepParams struct {
	StepID          string `json:"step_id" binding:"required"`
	EventInfoID     string `json:"event_info_id" binding:"required"`
	EventDateTimeID string `json:"event_date_time_id" binding:"required"`
	Des             string `json:"des"`
	Image           string `json:"image"`
	// This would tell us whether to use des or photo to update
	Type string `json:"type" binding:"required"`
}

type CreateEventCheckInStepParams struct {
	EventInfoID     string `json:"event_info_id" binding:"required"`
	EventDateTimeID string `json:"event_date_time_id" binding:"required"`
	Des             string `json:"des"`
	Image           string `json:"image"`
	// This would tell us whether to use des or photo to update
	Type string `json:"type" binding:"required"`
}

type RemoveEventCheckInStepParams struct {
	EventInfoID     string `json:"event_info_id" binding:"required"`
	EventDateTimeID string `json:"event_date_time_id" binding:"required"`
	StepID          string `json:"step_id" binding:"required"`
}

type ListEventCheckInStepRes struct {
	List        []CheckInStepRes `json:"list"`
	Street      string           `json:"street"`
	City        string           `json:"city"`
	State       string           `json:"state"`
	Country     string           `json:"country"`
	Postcode    string           `json:"postcode"`
	HasLocation bool             `json:"has_location"`
	Published   bool             `json:"published"`
}
