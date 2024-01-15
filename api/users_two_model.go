package api

type UpdateIdentityParams struct {
	Country     string `json:"country"`
	Type        string `json:"type"`
	IDPhoto     string `json:"id_photo"`
	IDBackPhoto string `json:"id_back_photo"`
	FacialPhoto string `json:"facial_photo"`
}

type UpdateIdentityRes struct {
	Status   string `json:"status"`
	Verified bool   `json:"verified"`
}

type CreateEmContactParams struct {
	Name         string `json:"name"`
	Relationship string `json:"relationship"`
	Email        string `json:"email"`
	// Code eg US, NG
	Code        string `json:"code"`
	DialCountry string `json:"dial_country"`
	PhoneNumber string `json:"phone_number"`
	Language    string `json:"language"`
}

type CreateEmContactRes struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type RemoveEmContactParams struct {
	ID string `json:"id"`
}

type EmContactDetail struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

type ProfileUserRes struct {
	Email       string            `json:"email"`
	PhoneNumber string            `json:"phone_number"`
	Status      string            `json:"status"`
	Verified    bool              `json:"verified"`
	PhoneCode   string            `json:"phone_code"`
	DateOfBirth string            `json:"date_of_birth"`
	FirstName   string            `json:"first_name"`
	LastName    string            `json:"last_name"`
	EmContacts  []EmContactDetail `json:"em_contacts"`
	Currency    string            `json:"currency"`
	UserTwoID   string            `json:"user_two_id"`
}

type UserProfileParams struct {
	Work      string   `json:"work"`
	Languages []string `json:"languages"`
	Bio       string   `json:"bio"`
	Type      string   `json:"type" binding:"required,user_profile_type"`
}

type UpdateProfilePhotoParams struct {
	ProfilePhoto string `json:"profile_photo" binding:"required"`
}

type GetUserCurrencyRes struct {
	Currency string `json:"currency"`
}

type GetUserIsHostRes struct {
	IsHost              bool   `json:"is_host"`
	HasIncomplete       bool   `json:"has_incomplete"`
	UnreadMessages      int    `json:"unread_messages"`
	UnreadNotifications int    `json:"unread_notifications"`
	ProfileImage        string `json:"profile_image"`
}

type GetUserProfilePhotoRes struct {
	ProfilePhoto string `json:"profile_photo"`
}

type CreateUpdateUserLocationParams struct {
	Street   string `json:"street"`
	City     string `json:"city"`
	State    string `json:"state" binding:"required"`
	Country  string `json:"country" binding:"required"`
	Postcode string `json:"postcode"`
	Lat      string `json:"lat" binding:"required"`
	Lng      string `json:"lng" binding:"required"`
}

type CreateUpdateUserLocationRes struct {
	State   string `json:"state"`
	Country string `json:"country"`
}

type ProfileDetailOption struct {
	OptionUserID string `json:"option_user_id"`
	Name         string `json:"name"`
	CoverImage   string `json:"cover_image"`
	Type         string `json:"type"`
	MainOption   string `json:"main_option"`
}

type GetUserProfileDetailRes struct {
	Status         string                `json:"status"`
	Verified       bool                  `json:"verified"`
	EmailConfirmed bool                  `json:"email"`
	PhoneConfirmed bool                  `json:"phone"`
	Bio            string                `json:"bio"`
	Work           string                `json:"work"`
	Languages      []string              `json:"languages"`
	State          string                `json:"state"`
	Country        string                `json:"country"`
	YearJoined     string                `json:"year_joined"`
	UserTwoId      string                `json:"user_two_id"`
	Options        []ProfileDetailOption `json:"options"`
}

type CreateFeedbackParams struct {
	Subject string `json:"subject" binding:"required"`
	Detail  string `json:"detail" binding:"required"`
}

type CreateHelpUserParams struct {
	Subject string `json:"subject" binding:"required"`
	Detail  string `json:"detail" binding:"required"`
}

type UpdateCurrencyUserParams struct {
	Currency string `json:"currency" binding:"required"`
}

type CreateHelpParams struct {
	Email   string `json:"email" binding:"required,email"`
	Subject string `json:"subject" binding:"required"`
	Detail  string `json:"detail" binding:"required"`
}

type UserResponseMsg struct {
	Success bool `json:"success"`
}
