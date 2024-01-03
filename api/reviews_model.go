package api

type CreateGeneralReviewParams struct {
	General        int    `json:"general"`
	ChargeID       string `json:"charge_id"`
	CurrentState   string `json:"current_state"`
	PreviousState  string `json:"previous_state"`
	MainOptionType string `json:"main_option_type"`
}

type GetStateReviewParams struct {
	ChargeID       string `json:"charge_id"`
	MainOptionType string `json:"main_option_type"`
}

type CreateDetailReviewParams struct {
	Environment    int    `json:"environment"`
	Accuracy       int    `json:"accuracy"`
	CheckIn        int    `json:"check_in"`
	Communication  int    `json:"communication"`
	Location       int    `json:"location"`
	ChargeID       string `json:"charge_id"`
	CurrentState   string `json:"current_state"`
	PreviousState  string `json:"previous_state"`
	MainOptionType string `json:"main_option_type"`
}

type CreatePrivateNoteReviewParams struct {
	PrivateNote    string `json:"private_note"`
	ChargeID       string `json:"charge_id"`
	CurrentState   string `json:"current_state"`
	PreviousState  string `json:"previous_state"`
	MainOptionType string `json:"main_option_type"`
}

type CreatePublicNoteReviewParams struct {
	PublicNote     string `json:"public_note"`
	ChargeID       string `json:"charge_id"`
	CurrentState   string `json:"current_state"`
	PreviousState  string `json:"previous_state"`
	MainOptionType string `json:"main_option_type"`
}

type PlaceholderReviewParams struct {
	ChargeID       string `json:"charge_id"`
	CurrentState   string `json:"current_state"`
	PreviousState  string `json:"previous_state"`
	MainOptionType string `json:"main_option_type"`
}

type CreateStayCleanReviewParams struct {
	StayClean      string `json:"stay_clean"`
	ChargeID       string `json:"charge_id"`
	CurrentState   string `json:"current_state"`
	PreviousState  string `json:"previous_state"`
	MainOptionType string `json:"main_option_type"`
}

type CreateComfortReviewParams struct {
	StayComfort    string `json:"stay_comfort"`
	ChargeID       string `json:"charge_id"`
	CurrentState   string `json:"current_state"`
	PreviousState  string `json:"previous_state"`
	MainOptionType string `json:"main_option_type"`
}

type CreateHostReviewParams struct {
	HostReview     string `json:"host_review"`
	ChargeID       string `json:"charge_id"`
	CurrentState   string `json:"current_state"`
	PreviousState  string `json:"previous_state"`
	MainOptionType string `json:"main_option_type"`
}

type CompleteOptionReviewParams struct {
	ChargeID      string `json:"charge_id"`
	CurrentState  string `json:"current_state"`
	PreviousState string `json:"previous_state"`
}

type CreateAmenityReviewItem struct {
	Tag      string `json:"tag"`
	Answer   string `json:"answer"`
	ChargeID string `json:"charge_id"`
}

type ListAmenityReviewItemParams struct {
	ChargeID string `json:"charge_id"`
}

type ListAmenityReviewItemRes struct {
	Amenities []string             `json:"amenities"`
	Selected  ListAmenityReviewRes `json:"selected"`
	IsEmpty   bool                 `json:"is_empty"`
}

type RemoveAmenityReviewItem struct {
	Tag      string `json:"tag"`
	ChargeID string `json:"charge_id"`
}

type ListAmenityReviewRes struct {
	List    []CreateAmenityReviewItem `json:"list"`
	IsEmpty bool                      `json:"is_empty"`
}

type ReviewRes struct {
	ChargeID       string `json:"charge_id"`
	CurrentState   string `json:"current_state"`
	PreviousState  string `json:"previous_state"`
	MainOptionType string `json:"main_option_type"`
}

type RemoveReviewParams struct {
	MainOptionType string `json:"main_option_type"`
	ChargeID       string `json:"charge_id"`
}
