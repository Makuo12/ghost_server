package api

type ListMessageContactParams struct {
	HasData bool `json:"has_data"`
	Offset  int  `json:"offset"`
}

type MessageContactItem struct {
	MsgID                      string `json:"msg_id"`
	ConnectedUserID            string `json:"connected_user_id"`
	FirstName                  string `json:"first_name"`
	Photo                      string `json:"photo"`
	LastMessage                string `json:"last_message"`
	LastMessageTime            string `json:"last_message_time"`
	UnreadMessageCount         int    `json:"unread_message_count"`
	UnreadUserRequestCount     int    `json:"unread_user_request_count"`
	UnreadUserCancelCount      int    `json:"unread_user_cancel_count"`
	UnreadHostCancelCount      int    `json:"unread_host_cancel_count"`
	UnreadHostChangeDatesCount int    `json:"unread_host_change_dates_count"`
	RoomID                     string `json:"room_id"`
}

type ListMessageContactRes struct {
	List        []MessageContactItem `json:"list"`
	Offset      int                  `json:"offset"`
	OnLastIndex bool                 `json:"on_last_index"`
	UserID      string               `json:"user_id"`
	Time        string               `json:"time"`
}

type ListMessageParams struct {
	HasData   bool   `json:"has_data"`
	Offset    int    `json:"offset"`
	ContactID string `json:"contact_id" binding:"required"`
}

type ListRequestNotifyParams struct {
	Offset    int    `json:"offset"`
	ContactID string `json:"contact_id" binding:"required"`
}

type RequestNotifyItem struct {
	MID            string `json:"mid" binding:"required"`
	MsgID          string `json:"msg_id" binding:"required"`
	HostNameOption string `json:"host_name_option" binding:"required"`
	Text           string `json:"text" binding:"required"`
	StartDate      string `json:"start_date" binding:"required"`
	EndDate        string `json:"end_date" binding:"required"`
	MainOptionType string `json:"main_option_type" binding:"required"`
	Category       string `json:"category" binding:"required"`
	SpecialType    string `json:"special_type" binding:"required"`
	Reference      string `json:"reference"`
}

type RequestNotifyDetailParams struct {
	Reference   string `json:"reference" binding:"required"`
	ContactID   string `json:"contact_id" binding:"required"`
	SpecialType string `json:"special_type" binding:"required"`
	MID         string `json:"mid" binding:"required"`
	Currency    string `json:"currency" binding:"required"`
}

type OUserRequestNotifyDetailRes struct {
	StartDate        string   `json:"start_date" binding:"required"`
	EndDate          string   `json:"end_date" binding:"required"`
	CoverImage       string   `json:"cover_image" binding:"required"`
	Guests           []string `json:"guests" binding:"required"`
	Price            string   `json:"price" binding:"required"`
	EmailVerified    bool     `json:"email_verified" binding:"required"`
	PhoneVerified    bool     `json:"phone_verified" binding:"required"`
	IdentityVerified bool     `json:"identity_verified" binding:"required"`
	Text             string   `json:"text" binding:"required"`
}

type ListRequestNotifyRes struct {
	List              []RequestNotifyItem `json:"list"`
	Offset            int                 `json:"offset"`
	OnLastIndex       bool                `json:"on_last_index" binding:"required"`
	SelectedContactID string              `json:"selected_contact_id"`
}

type MessageMainItem struct {
	MainMsg        MessageItem `json:"main_msg"`
	MainMsgEmpty   bool        `json:"main_msg_empty"`
	ParentMsg      MessageItem `json:"parent_msg"`
	ParentMsgEmpty bool        `json:"parent_msg_empty"`
}

type MessageItem struct {
	ID         string `json:"id"`
	SenderID   string `json:"sender_id"`
	ReceiverID string `json:"receiver_id"`
	Message    string `json:"message"`
	Type       string `json:"type"`
	Read       bool   `json:"read"`
	Photo      string `json:"photo"`
	ParentID   string `json:"parent_id"`
	Reference  string `json:"reference"`
	CreatedAt  string `json:"created_at"`
}

type ListMessageRes struct {
	List              []MessageMainItem `json:"list"`
	Offset            int               `json:"offset"`
	OnLastIndex       bool              `json:"on_last_index" binding:"required"`
	SelectedContactID string            `json:"selected_contact_id"`
	RoomID            string            `json:"room_id" binding:"required"`
	IsEmpty           bool              `json:"is_empty"`
}

type ListMessageContactListenRes struct {
	List        []MessageContactItem `json:"list"`
	CurrentTime string               `json:"current_time"`
}
