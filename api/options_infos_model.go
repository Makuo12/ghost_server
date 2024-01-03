package api

type ListIncompleteOptionInfosParams struct {
	Offset  int  `json:"offset"`
	IsStart bool `json:"is_start"`
}

type CreateOptionInfoParams struct {
	Currency           string `json:"currency" binding:"required,currency"`
	OptionItemType     string `json:"option_item_type"`
	UserOptionID       string `json:"user_option_id"`
	OptionID           string `json:"option_id"`
	CurrentServerView  string `json:"current_server_view"`
	PreviousServerView string `json:"previous_server_view"`
	MainOptionType     string `json:"main_option_type" binding:"required"`
	OptionType         string `json:"option_type" binding:"required"`
	OptionImg          string `json:"option_img" binding:"required"`
}
type CreateLocationParams struct {
	UserOptionID         string `json:"user_option_id"`
	OptionID             string `json:"option_id"`
	Street               string `json:"street" binding:"required"`
	City                 string `json:"city" binding:"required"`
	State                string `json:"state" binding:"required"`
	Country              string `json:"country" binding:"required"`
	Postcode             string `json:"postcode" binding:"required"`
	Lat                  string `json:"lat" binding:"required"`
	Lng                  string `json:"lng" binding:"required"`
	MainOptionType       string `json:"main_option_type" binding:"required"`
	ShowSpecificLocation bool   `json:"show_specific_location"`
	CurrentServerView    string `json:"current_server_view"`
	PreviousServerView   string `json:"previous_server_view"`
	Currency             string `json:"currency"`
	OptionType           string `json:"option_type"`
}

type CreateOptionSE struct {
	Currency           string `json:"currency" binding:"required,currency"`
	UserOptionID       string `json:"user_option_id"`
	OptionID           string `json:"option_id"`
	CurrentServerView  string `json:"current_server_view"`
	PreviousServerView string `json:"previous_server_view"`
	MainOptionType     string `json:"main_option_type" binding:"required"`
	OptionType         string `json:"option_type" binding:"required"`
	OptionImg          string `json:"option_img" binding:"required"`
	ShortletType       string `json:"shortlet_type" binding:"required"`
	EventType          string `json:"event_type" binding:"required"`
	TimeZone           string `json:"time_zone" binding:"required"`
}

type OptionInfoResponse struct {
	Success            bool   `json:"success"`
	UserOptionID       string `json:"user_option_id"`
	OptionID           string `json:"option_id"`
	CurrentServerView  string `json:"current_server_view"`
	PreviousServerView string `json:"previous_server_view"`
	MainOptionType     string `json:"main_option_type" binding:"required"`
	OptionType         string `json:"option_type" binding:"required"`
	Currency           string `json:"currency" binding:"required,currency"`
	OptionSubType      string `json:"option_sub_type"`
	OptionItemType     string `json:"option_item_type"`
}

type OptionInfoFirstResponse struct {
	Success            bool   `json:"success"`
	UserOptionID       string `json:"user_option_id"`
	OptionID           string `json:"option_id"`
	CurrentServerView  string `json:"current_server_view"`
	PreviousServerView string `json:"previous_server_view"`
	MainOptionType     string `json:"main_option_type" binding:"required"`
	OptionType         string `json:"option_type" binding:"required"`
	Currency           string `json:"currency" binding:"required,currency"`
	OptionSubType      string `json:"option_sub_type"`
	OptionItemType     string `json:"option_item_type"`
	UserIsHost         bool   `json:"user_is_host"`
	HasIncomplete      bool   `json:"has_incomplete"`
}

type OptionInfoRemoveRequest struct {
	UserOptionID string `json:"user_option_id"`
	OptionID     string `json:"option_id"`
}

type CreateAmenitiesAndSafety struct {
	UserOptionID       string   `json:"user_option_id"`
	OptionID           string   `json:"option_id" binding:"required"`
	PopularAm          []string `json:"popular_am"`
	HomeSafetyAm       []string `json:"home_safety_am"`
	CurrentServerView  string   `json:"current_server_view"`
	PreviousServerView string   `json:"previous_server_view"`
	MainOptionType     string   `json:"main_option_type" binding:"required"`
	OptionType         string   `json:"option_type" binding:"required"`
	Currency           string   `json:"currency" binding:"required,currency"`
}

type OptionInfoRemoveResponse struct {
	Success bool `json:"success"`
}
type RemoveResponse struct {
	Success bool `json:"success"`
}

type OptionInfoRemoveFirstResponse struct {
	Success       bool `json:"success"`
	UserIsHost    bool `json:"user_is_host"`
	HasIncomplete bool `json:"has_incomplete"`
}

type OptionTypeRemoveResponse struct {
	Success      bool   `json:"success"`
	UserOptionID string `json:"user_option_id"`
	OptionID     string `json:"option_id"`
}

type CreateShortletTypeParams struct {
	UserOptionID       string `json:"user_option_id"`
	OptionID           string `json:"option_id"`
	ShortletType       string `json:"shortlet_type" binding:"required,shortlet_type"`
	CurrentServerView  string `json:"current_server_view"`
	PreviousServerView string `json:"previous_server_view"`
	MainOptionType     string `json:"main_option_type" binding:"required"`
	OptionType         string `json:"option_type" binding:"required"`
	Currency           string `json:"currency" binding:"required,currency"`
}

type CreateShortletSpaceParams struct {
	UserOptionID       string `json:"user_option_id"`
	OptionID           string `json:"option_id"`
	SpaceType          string `json:"space_type" binding:"required,shortlet_space"`
	CurrentServerView  string `json:"current_server_view"`
	PreviousServerView string `json:"previous_server_view"`
	MainOptionType     string `json:"main_option_type" binding:"required"`
	OptionType         string `json:"option_type" binding:"required"`
	Currency           string `json:"currency" binding:"required,currency"`
}

type CreateEventSubCategoryParams struct {
	UserOptionID       string `json:"user_option_id"`
	OptionID           string `json:"option_id"`
	SubCategoryType    string `json:"sub_category_type" binding:"required"`
	EventType          string `json:"event_type" binding:"required"`
	CurrentServerView  string `json:"current_server_view"`
	PreviousServerView string `json:"previous_server_view"`
	MainOptionType     string `json:"main_option_type" binding:"required"`
	OptionType         string `json:"option_type" binding:"required"`
	Currency           string `json:"currency" binding:"required,currency"`
}

type CreateEventLocationParams struct {
	UserOptionID         string `json:"user_option_id"`
	OptionID             string `json:"option_id"`
	LocationType         string `json:"location_type" binding:"required,event_location_type"`
	Street               string `json:"street"`
	City                 string `json:"city"`
	State                string `json:"state"`
	Country              string `json:"country"`
	Postcode             string `json:"postcode"`
	Lat                  string `json:"lat"`
	Lng                  string `json:"lng"`
	ShowSpecificLocation bool   `json:"show_specific_location"`
	CurrentServerView    string `json:"current_server_view"`
	PreviousServerView   string `json:"previous_server_view"`
	MainOptionType       string `json:"main_option_type" binding:"required"`
	OptionType           string `json:"option_type" binding:"required"`
	Currency             string `json:"currency" binding:"required,currency"`
}

type CreateEventDateParams struct {
	UserOptionID       string `json:"user_option_id"`
	OptionID           string `json:"option_id"`
	DateType           string `json:"date_type" binding:"required,event_date_type"`
	StartingDate       string `json:"starting_date"`
	EndingDate         string `json:"ending_date"`
	StartTime          string `json:"start_time"`
	EndTime            string `json:"end_time"`
	CurrentServerView  string `json:"current_server_view"`
	PreviousServerView string `json:"previous_server_view"`
	MainOptionType     string `json:"main_option_type" binding:"required"`
	OptionType         string `json:"option_type" binding:"required"`
	Currency           string `json:"currency" binding:"required,currency"`
}

type CreateShortletInfoParams struct {
	UserOptionID       string   `json:"user_option_id"`
	OptionID           string   `json:"option_id"`
	AnySpaceShared     bool     `json:"any_space_shared"`
	Space              []string `json:"space" binding:"required"`
	CurrentServerView  string   `json:"current_server_view"`
	PreviousServerView string   `json:"previous_server_view"`
	MainOptionType     string   `json:"main_option_type" binding:"required"`
	OptionType         string   `json:"option_type" binding:"required"`
	Currency           string   `json:"currency" binding:"required,currency"`
}
type ReverseShortletInfoParams struct {
	UserOptionID       string   `json:"user_option_id"`
	OptionID           string   `json:"option_id"`
	AnySpaceShared     bool     `json:"any_space_shared"`
	Space              []string `json:"space" binding:"required"`
	SpaceType          string   `json:"space_type" binding:"required"`
	CurrentServerView  string   `json:"current_server_view"`
	PreviousServerView string   `json:"previous_server_view"`
	MainOptionType     string   `json:"main_option_type" binding:"required"`
	OptionType         string   `json:"option_type" binding:"required"`
	Currency           string   `json:"currency" binding:"required,currency"`
}

type CreateOptionQuestionParams struct {
	UserOptionID        string `json:"user_option_id"`
	OptionID            string `json:"option_id"`
	HasSecurityCamera   bool   `json:"has_security_camera"`
	HostAsIndividual    bool   `json:"host_as_individual"`
	HasWeapons          bool   `json:"has_weapons"`
	HasDangerousAnimals bool   `json:"has_dangerous_animals"`
	CurrentServerView   string `json:"current_server_view"`
	PreviousServerView  string `json:"previous_server_view"`
	MainOptionType      string `json:"main_option_type" binding:"required"`
	OptionType          string `json:"option_type" binding:"required"`
	Currency            string `json:"currency" binding:"required,currency"`
	OrganizationName    string `json:"organization_name"`
}

type CreateYatchInfoParams struct {
	UserOptionID       string `json:"user_option_id"`
	OptionID           string `json:"option_id"`
	NumFullBathrooms   int    `json:"num_full_bathrooms" binding:"required"`
	NumHalfBathrooms   int    `json:"num_half_bathrooms" binding:"required"`
	NumOfGuest         int    `json:"num_of_guest" binding:"required,gte=0"`
	NumOfBedrooms      int    `json:"num_of_bedrooms" binding:"required,gte=0"`
	CurrentServerView  string `json:"current_server_view"`
	PreviousServerView string `json:"previous_server_view"`
	MainOptionType     string `json:"main_option_type" binding:"required"`

	OptionType string `json:"option_type" binding:"required"`
	Currency   string `json:"currency" binding:"required,currency"`
}

type CreateOptionInfoDescription struct {
	UserOptionID       string `json:"user_option_id"`
	OptionID           string `json:"option_id"`
	Description        string `json:"description" binding:"required"`
	CurrentServerView  string `json:"current_server_view"`
	PreviousServerView string `json:"previous_server_view"`
	MainOptionType     string `json:"main_option_type" binding:"required"`

	OptionType string `json:"option_type" binding:"required"`
	Currency   string `json:"currency" binding:"required,currency"`
}
type PublishOption struct {
	UserOptionID       string `json:"user_option_id"`
	OptionID           string `json:"option_id"`
	CurrentServerView  string `json:"current_server_view"`
	PreviousServerView string `json:"previous_server_view"`
	MainOptionType     string `json:"main_option_type" binding:"required"`
	OptionType         string `json:"option_type" binding:"required"`
	Currency           string `json:"currency" binding:"required,currency"`
}

type PublishShortletOptionParams struct {
	UserOptionID       string `json:"user_option_id"`
	OptionID           string `json:"option_id"`
	CurrentServerView  string `json:"current_server_view"`
	PreviousServerView string `json:"previous_server_view"`
	MainOptionType     string `json:"main_option_type" binding:"required"`
	OptionType         string `json:"option_type" binding:"required"`
	Currency           string `json:"currency" binding:"required,currency"`
	Name               string `json:"name" binding:"required"`
	// option_main_type is would be main_event_type for events and main_shortlet_type for shortlets
	OptionMainType       string   `json:"option_main_type" binding:"required"`
	Space                []string `json:"space" binding:"required"`
	NumOfGuest           int      `json:"num_of_guests" binding:"required"`
	Description          string   `json:"description" binding:"required"`
	PopularAm            []string `json:"popular_am" binding:"required"`
	HomeSafetyAm         []string `json:"home_safety_am" binding:"required"`
	Street               string   `json:"street"`
	CoverImage           string   `json:"cover_image"`
	City                 string   `json:"city"`
	State                string   `json:"state"`
	Country              string   `json:"country"`
	Postcode             string   `json:"postcode"`
	ShowSpecificLocation bool     `json:"show_specific_location"`
	FirstName            string   `json:"first_name"`
}

type HostCurrentOptionParams struct {
	OptionID       string `json:"option_id"`
	MainOptionType string `json:"main_option_type" binding:"required"`
	OptionType     string `json:"option_type" binding:"required"`
	Currency       string `json:"currency" binding:"required,currency"`
	HostNameOption string `json:"host_name_option" binding:"required"`
	CoverImage     string `json:"cover_image"`
	State          string `json:"state"`
	Country        string `json:"country"`
	// option_main_type is would be main_event_type for events and main_shortlet_type for shortlets
	OptionMainType string `json:"option_main_type"`
}

type PublishEventOptionParams struct {
	UserOptionID       string `json:"user_option_id"`
	OptionID           string `json:"option_id"`
	CurrentServerView  string `json:"current_server_view"`
	PreviousServerView string `json:"previous_server_view"`
	MainOptionType     string `json:"main_option_type" binding:"required"`
	OptionType         string `json:"option_type" binding:"required"`
	Currency           string `json:"currency" binding:"required,currency"`
	Name               string `json:"name" binding:"required"`
	// option_main_type is would be main_event_type for events and main_shortlet_type for shortlets
	OptionMainType string `json:"option_main_type" binding:"required"`
	Description    string `json:"description" binding:"required"`
	CoverImage     string `json:"cover_image"`
	FirstName      string `json:"first_name"`
}

type CreateRecreationInfo struct {
	UserOptionID       string   `json:"user_option_id"`
	OptionID           string   `json:"option_id"`
	RecreationTypes    []string `json:"recreation_types"`
	CurrentServerView  string   `json:"current_server_view"`
	PreviousServerView string   `json:"previous_server_view"`
	MainOptionType     string   `json:"main_option_type" binding:"required"`
	OptionType         string   `json:"option_type" binding:"required"`
	Currency           string   `json:"currency" binding:"required,currency"`
}

type CreateOptionInfoName struct {
	UserOptionID       string `json:"user_option_id"`
	OptionID           string `json:"option_id"`
	HostNameOption     string `json:"host_option_name" binding:"required"`
	CurrentServerView  string `json:"current_server_view"`
	PreviousServerView string `json:"previous_server_view"`
	MainOptionType     string `json:"main_option_type" binding:"required"`
	OptionType         string `json:"option_type" binding:"required"`
	Currency           string `json:"currency" binding:"required,currency"`
}

type CreateOptionPrice struct {
	UserOptionID       string `json:"user_option_id"`
	OptionID           string `json:"option_id"`
	Price              string `json:"price" binding:"required"`
	CurrentServerView  string `json:"current_server_view"`
	PreviousServerView string `json:"previous_server_view"`
	MainOptionType     string `json:"main_option_type" binding:"required"`
	OptionType         string `json:"option_type" binding:"required"`
	Currency           string `json:"currency" binding:"required,currency"`
}

type CreateOptionInfoHighlight struct {
	UserOptionID       string   `json:"user_option_id"`
	OptionID           string   `json:"option_id"`
	Highlight          []string `json:"highlight" binding:"required"`
	CurrentServerView  string   `json:"current_server_view"`
	PreviousServerView string   `json:"previous_server_view"`
	MainOptionType     string   `json:"main_option_type" binding:"required"`
	OptionType         string   `json:"option_type" binding:"required"`
	Currency           string   `json:"currency" binding:"required,currency"`
}

type PublishToHostViewParams struct {
	UserOptionID       string `json:"user_option_id"`
	OptionID           string `json:"option_id"`
	CurrentServerView  string `json:"current_server_view"`
	PreviousServerView string `json:"previous_server_view"`
	MainOptionType     string `json:"main_option_type" binding:"required"`
	OptionType         string `json:"option_type" binding:"required"`
	Currency           string `json:"currency" binding:"required,currency"`
}

type OptionInfo struct {
	ID             int    `json:"id"`
	OptionType     string `json:"option_type"`
	HostNameOption string `json:"host_name_option"`
	IsActive       bool   `json:"is_active"`
}

type CreateOptionInfoPhotoParams struct {
	UserOptionID string `json:"user_option_id"`
	OptionID     string `json:"option_id"`

	CoverImage         string   `json:"cover_image"`
	Photo              []string `json:"photo"`
	CurrentServerView  string   `json:"current_server_view"`
	PreviousServerView string   `json:"previous_server_view"`
	MainOptionType     string   `json:"main_option_type" binding:"required"`
	OptionType         string   `json:"option_type" binding:"required"`
	Currency           string   `json:"currency" binding:"required,currency"`
}

type UpdatedOptionInfoFieldResponse struct {
	Updated bool `json:"updated"`
}

type ListOptionInfoNotCompleteRow struct {
	ID                 string `json:"id"`
	IsComplete         bool   `json:"is_complete"`
	MainOptionType     string `json:"main_option_type"`
	OptionType         string `json:"option_type"`
	HostNameOption     string `json:"host_name_option"`
	CreatedAt          string `json:"created_at"`
	CurrentServerView  string `json:"current_server_view"`
	PreviousServerView string `json:"previous_server_view"`
	Currency           string `json:"currency"`
	ExtraInfo          string `json:"extra_info"`
	//// option_main_type is would be main_event_type for events and main_shortlet_type for shortlets
	//OptionMainType string `json:"option_main_type" binding:"required"`
}

type GetOptionInfoNotComplete struct {
	ID                 string `json:"id"`
	IsComplete         bool   `json:"is_complete"`
	MainOptionType     string `json:"main_option_type"`
	OptionType         string `json:"option_type"`
	HostNameOption     string `json:"host_name_option"`
	CreatedAt          string `json:"created_at"`
	CurrentServerView  string `json:"current_server_view"`
	PreviousServerView string `json:"previous_server_view"`
	Currency           string `json:"currency"`
	ExtraInfo          string `json:"extra_info"`
	//// option_main_type is would be main_event_type for events and main_shortlet_type for shortlets
	//OptionMainType string `json:"option_main_type" binding:"required"`
}

type ListOptionInfoNotCompleteRowParams struct {
	OptionInfos []ListOptionInfoNotCompleteRow `json:"option_infos"`
}

type GetSingleOption struct {
	UserOptionID string `json:"user_option_id"`
	OptionID     string `json:"option_id"`
}

type SwitchToHostingParams struct {
	CanSwitchToHosting bool `json:"can_switch_to_hosting"`
}

type OptionOffsetParams struct {
	OptionOffset int `json:"option_offset"`
}

type CalenderOptionList struct {
	List         []CalenderOptionItem `json:"list"`
	OptionOffset int                  `json:"option_offset"`
	OnLastIndex  bool                 `json:"on_last_index"`
}

type CalenderOptionItem struct {
	MainOptionType string `json:"main_option_type"`
	HostNameOption string `json:"host_name_option"`
	OptionID       string `json:"option_id"`
	Currency       string `json:"currency"`
}

type UserEventSearchItem struct {
	HostNameOption string `json:"host_name_option"`
	OptionUserID   string `json:"option_user_id"`
	IsVerified     bool   `json:"is_verified"`
	CoverImage     string `json:"cover_image"`
}

type OptionQuestionNote struct {
	Tag  string `json:"tag"`
	Type string `json:"type"`
}

var optionQuestionNote = []OptionQuestionNote{
	{
		Tag:  "cameras_audio_devices",
		Type: "safety_devices",
	},
	{
		Tag:  "dangerous_animal",
		Type: "safety_consider",
	},
	{
		Tag:  "weapon_on_property",
		Type: "property_info",
	},
}
