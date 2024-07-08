package api

type GetChargeCodeParams struct {
	ID         string `json:"id"`
	MainOption string `json:"main_option"`
}

type GetChargeCodeRes struct {
	Code       string `json:"code"`
	ID         string `json:"id"`
	TicketType string `json:"ticket_type"`
	Grade      string `json:"grade"`
	MainOption string `json:"main_option"`
}

type GetChargeCodeScannedRes struct {
	Message          string `json:"message"`
	ScannedUserImage string `json:"scanned_user_image"`
	ScannedByName    string `json:"scanned_by_name"`
	ScannedTime      string `json:"scanned_time"`
	ID               string `json:"id"`
}

type DeleteChargeCodeParams struct {
	ID         string `json:"id"`
	MainOption string `json:"main_option"`
}

type DeleteChargeCodeRes struct {
	WasScanned bool `json:"was_scanned"`
	Success    bool `json:"success"`
}

type ValidateChargeCodeParams struct {
	OptionID        string `json:"option_id"`
	Code            string `json:"code"`
	EventDateTimeID string `json:"event_date_time_id"`
	MainOption      string `json:"main_option"`
	StartDate       string `json:"start_date"`
}

type ValidateChargeCodeRes struct {
	FirstName      string `json:"first_name"`
	HostOptionName string `json:"host_option_name"`
	Grade          string `json:"grade"`
	TicketType     string `json:"ticket_type"`
	Success        bool   `json:"success"`
	MainOption     string `json:"main_option"`
	StartDate      string `json:"start_date"`
	Used           bool   `json:"used"`
	EndDate        string `json:"end_date"`
}
