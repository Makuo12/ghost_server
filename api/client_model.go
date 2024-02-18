package api

type SearchText struct {
	Text string `json:"text"`
}

type EventSearchText struct {
	Currency string `json:"currency"`
	Text     string `json:"text"`
}

type EventSearchTextRes struct {
	List []ExperienceEventData `json:"list"`
}

type MapExperienceLocationParams struct {
	Lat string `json:"lat"`
	Lng string `json:"lng"`
}

type MapExperienceLocationItem struct {
	Lat             string `json:"lat"`
	Lng             string `json:"lng"`
	Name            string `json:"name"`
	StartDate       string `json:"start_date"`
	EndDate         string `json:"end_date"`
	Country         string `json:"country"`
	State           string `json:"state"`
	EventDateTimeID string `json:"event_date_time_id"`
	OptionUserID    string `json:"option_user_id"`
	MainOption      string `json:"main_option"`
	CoverImage      string `json:"cover_image"`
	Price           string `json:"price"`
}

type MapExperienceLocationRes struct {
	List []MapExperienceLocationItem `json:"list"`
}

type SearchTextRes struct {
	List []UHMOptionSelectionRes `json:"list"`
}

type SearchTextCalRes struct {
	List []CalenderOptionItem `json:"list"`
}

// EDT event_date_time
type SearchTextEDT struct {
	Text        string `json:"text"`
	EventInfoID string `json:"event_info_id"`
}

// EDT event_date_time
type SearchTextEDTRes struct {
	List []EventDateItem `json:"list"`
}

type SearchEventByNameRes struct {
	List []UserEventSearchItem `json:"list"`
}

type CurrentTime struct {
	// The Time is the last time data was updated from the front end
	Time string `json:"time"`
}

type CreateMessageParams struct {
	SenderID          string `json:"sender_id"`
	ReceiverID        string `json:"receiver_id"`
	Message           string `json:"message"`
	Type              string `json:"type"`
	Photo             string `json:"photo"`
	ParentID          string `json:"parent_id"`
	Reference         string `json:"reference"`
	SelectedContactID string `json:"selected_contact_id"`
}

type CreateMessageRes struct {
	MsgID     string `json:"msg_id"`
	CreatedAt string `json:"created_at"`
}

type UnreadMessageParams struct {
	SelectedContactID string `json:"selected_contact_id"`
	LatestTime        string `json:"latest_time"`
	UserID            string `json:"user_id"`
}

type GetMessageParams struct {
	SelectedContactID string `json:"selected_contact_id"`
	LatestTime        string `json:"latest_time"`
	UserID            string `json:"user_id"`
}

type GetMessageRes struct {
	MsgList           []MessageMainItem `json:"msg_list"`
	MsgEmpty          bool              `json:"msg_empty"`
	SelectedContactID string            `json:"selected_contact_id"`
	UserID            string            `json:"user_id"`
}

type UnreadMessageRes struct {
	List              []string `json:"list"`
	SelectedContactID string   `json:"selected_contact_id"`
	UserID            string   `json:"user_id"`
}
