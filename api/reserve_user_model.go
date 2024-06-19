package api

type ReserveUserItem struct {
	ID                     string   `json:"id"`
	HostUserID             string   `json:"host_user_id"`
	MainOption             string   `json:"main_option"`
	HostNameOption         string   `json:"host_name_option"`
	StartDate              string   `json:"start_date"`
	EndDate                string   `json:"end_date"`
	HostName               string   `json:"host_name"`
	ProfilePhoto           string   `json:"profile_photo"`
	OptionCoverImage       string   `json:"option_cover_image"`
	PublicOptionCoverImage string   `json:"public_option_cover_image"`
	OptionPhotos           []string `json:"option_photos"`
	PublicOptionPhotos     []string `json:"public_option_photos"`
	PublicProfilePhoto     string   `json:"public_profile_photo"`
	StartTime              string   `json:"start_time"`
	EndTime                string   `json:"end_time"`
	CheckInMethod          string   `json:"check_in_method"`
	Type                   string   `json:"type"`
	Grade                  string   `json:"grade"`
	OptionType             string   `json:"option_type"`
	SpaceType              string   `json:"space_type"`
	State                  string   `json:"state"`
	Country                string   `json:"country"`
	Street                 string   `json:"street"`
	Timezone               string   `json:"timezone"`
	City                   string   `json:"city"`
	ReviewStatus           string   `json:"review_status"`
	RoomID                 string   `json:"room_id"`
}

type ListReserveUserItemRes struct {
	List        []ReserveUserItem `json:"list"`
	Offset      int               `json:"offset"`
	OnLastIndex bool              `json:"on_last_index"`
	MainOption  string            `json:"main_option"`
	UserID      string            `json:"user_id"`
}

type ListReserveUserItemParams struct {
	Offset     int    `json:"offset"`
	MainOption string `json:"main_option"`
	Type       string `json:"type"`
}

type ReserveUserInfoParams struct {
	ID         string `json:"id"`
	MainOption string `json:"main_option"`
}

type ReserveUserDirectionRes struct {
	Street    string `json:"street"`
	City      string `json:"city"`
	State     string `json:"state"`
	Country   string `json:"country"`
	Postcode  string `json:"postcode"`
	Lat       string `json:"lat"`
	Lng       string `json:"lng"`
	Direction string `json:"direction"`
}

type ReserveUserWifiRes struct {
	NetworkName string `json:"network_name"`
	Password    string `json:"password"`
}

// RU -> ReserveUser
type RUCheckInStepRes struct {
	ID    string `json:"id"`
	Des   string `json:"des"`
	Photo string `json:"photo"`
}

// RU -> ReserveUser
type RUListCheckInStepRes struct {
	List []RUCheckInStepRes `json:"list"`
}

// RU -> ReserveUser
type RUHelpManualRes struct {
	Help string `json:"help"`
}

type RUCheckInMethodRes struct {
	Des  string `json:"des"`
	Type string `json:"type"`
}

type ReserveEventReceiptRes struct {
	Grade          string `json:"grade"`
	Price          string `json:"price"`
	Currency       string `json:"currency"`
	Type           string `json:"type"`
	TicketType     string `json:"ticket_type"`
	HostNameOption string `json:"host_name_option"`
}

type ReserveOptionReceiptRes struct {
	Discount        string   `json:"discount"`
	MainPrice       string   `json:"main_price"`
	ServiceFee      string   `json:"service_fee"`
	TotalFee        string   `json:"total_fee"`
	DatePrice       []string `json:"date_price"`
	Currency        string   `json:"currency"`
	GuestFee        string   `json:"guest_fee"`
	PetFee          string   `json:"pet_fee"`
	CleanFee        string   `json:"clean_fee"`
	NightlyPetFee   string   `json:"nightly_pet_fee"`
	NightlyGuestFee string   `json:"nightly_guest_fee"`
	HostNameOption  string   `json:"host_name_option"`
}
