package api

// "time"

// "github.com/google/uuid"

type CreateUserRequest struct {
	Password    string `json:"password" binding:"required,password"`
	Email       string `json:"email" binding:"required,email"`
	Username    string `json:"username" binding:"required"`
	FirstName   string `json:"first_name" binding:"required,person_name"`
	Currency    string `json:"currency" binding:"required,currency"`
	LastName    string `json:"last_name" binding:"required,person_name"`
	DateOfBirth string `json:"date_of_birth" binding:"required,date_only"`
}

type VerifyPhoneNumberRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
	Code        string `json:"code" binding:"required,len=6"`
	CountryName string `json:"country_name" binding:"required"`
}

type VerifyPhoneNumberResponse struct {
	CodeSent    bool   `json:"code_sent"`
	Username    string `json:"username"`
	PhoneNumber string `json:"phone_number" binding:"required"`
}

type VerifyEmailResponse struct {
	CodeSent bool   `json:"code_sent"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type DetailResponse struct {
	Success bool `json:"success"`
}

type ConfirmPhoneNumberRequest struct {
	Code     string `json:"code" binding:"required,len=6"`
	Username string `json:"username" binding:"required"`
}

type ConfirmPhoneNumberResponse struct {
	Confirmed bool   `json:"confirmed"`
	Username  string `json:"username"`
}
type LoginUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,password"`
}
type LoginUserResponse struct {
	Email                string `json:"email"`
	FireFight            string `json:"fire_fight"`
	FirstName            string `json:"first_name"`
	LastName             string `json:"last_name"`
	AccessToken          string `json:"access_token"`
	Currency             string `json:"currency"`
	ProfilePhoto         string `json:"profile_photo"`
	PublicID             string `json:"public_id"`
	AccessTokenExpiresAt string `json:"access_token_expires_at"`
	RefreshToken         string `json:"refresh_token"`
}

type UpdateUserResponse struct {
	Email         string `json:"email"`
	FireFight     string `json:"fire_fight"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	ProfilePhoto  string `json:"profile_photo"`
	Currency      string `json:"currency"`
	IsHost        bool   `json:"is_host"`
	HasIncomplete bool   `json:"has_incomplete"`
	DateOfBirth   string `json:"date_of_birth"`
}

type GetFireEmailAndPasswordRes struct {
	Email        string `json:"email"`
	FireFight    string `json:"fire_fight"`
	ProfilePhoto string `json:"profile_photo"`
}

type ProfileUserResponse struct {
	Email        string `json:"email"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	ProfilePhoto string `json:"profile_photo"`
	DateOfBirth  string `json:"date_of_birth"`
	IsHost       bool   `json:"is_host"`
	Rating       string `json:"rating"`
	DialCode     string `json:"dial_code"`
	DialCountry  string `json:"dial_country"`
	PhoneNumber  string `json:"phone_number"`
}

type ProfileHostResponse struct {
	HostLevel      string `json:"host_level"`
	Rating         string `json:"rating"`
	IsHostVerified bool   `json:"is_host_verified"`
	Bio            string `json:"bio"`
}

type ProfileAllResponse struct {
	ProfileUser ProfileUserResponse `json:"profile_user"`
	ProfileHost ProfileHostResponse `json:"profile_host"`
}

type UpdateUserInfoParams struct {
	DateOfBirth string `json:"date_of_birth" binding:"required,date_only"`
	FirstName   string `json:"first_name" binding:"required,person_name"`
	LastName    string `json:"last_name" binding:"required,person_name"`
}

// type LoginHostResponse struct {
// 	HostLevel      string `json:"host_level"`
// 	Rating         string `json:"rating"`
// 	IsActive       bool   `json:"is_active"`
// 	IsHostVerified bool   `json:"is_host_verified"`
// 	Bio            string `json:"bio"`
// }

// type LoginResponse struct {
// 	User LoginUserResponse `json:"user"`
// 	Host LoginHostResponse `json:"host"`
// }

type UpdateCurrencyParams struct {
	Currency string `json:"currency" binding:"required,currency"`
}

type UserFullDataResponse struct {
	UserUpdate UpdateUserResponse `json:"user_update"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required,password"`
	Password        string `json:"password" binding:"required,password"`
	PasswordTwo     string `json:"password_two" binding:"required,password"`
}

type GetAppPolicyParams struct {
	Type string `json:"type" binding:"required"`
}

type GetAppPolicyRes struct {
	Type string `json:"type" binding:"required"`
	Link string `json:"link" binding:"required"`
}

type ChangePasswordResponse struct {
	Updated              bool   `json:"updated"`
	AccessToken          string `json:"access_token"`
	AccessTokenExpiresAt string `json:"access_token_expires_at"`
	RefreshToken         string `json:"refresh_token"`
}

type ChangeEmailRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ChangeEmailResponse struct {
	Updated bool   `json:"updated"`
	Email   string `json:"email"`
}

type ChangePhoneRequest struct {
	Username string `json:"username" binding:"required"`
}
type ChangePhoneResponse struct {
	Updated     bool   `json:"updated"`
	PhoneNumber string `json:"phone_number"`
	Code        string `json:"code"`
	CountryName string `json:"country_name"`
}

type ForgotPasswordNotLoggedRequest struct {
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Code        string `json:"code"`
	CountryName string `json:"country_name"`
}

type JoinWithPhoneParams struct {
	PhoneNumber string `json:"phone_number"`
	Code        string `json:"code"`
	CountryName string `json:"country_name"`
}

type ConfirmCodeJoinParams struct {
	Currency  string `json:"currency"`
	Type      string `json:"type"`
	Code      string `json:"code" binding:"required,len=6"`
	Username  string `json:"username" binding:"required"`
}





type ConfirmJoinSignUpParams struct {
	Type      string `json:"type"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username" binding:"required"`
}



type ConfirmJoinSignUpRes struct {
	Type        string `json:"type"`
	CodeSent    bool   `json:"code_sent"`
	Username    string `json:"username"`
	PhoneNumber string `json:"phone_number" binding:"required"`
}



type JoinWithPhoneRes struct {
	Type        string `json:"type"`
	CountryName string `json:"country_name"`
	CodeSent    bool   `json:"code_sent"`
	Username    string `json:"username"`
	PhoneNumber string `json:"phone_number" binding:"required"`
}

type UpdateVerifyEmailPhoneRequest struct {
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Code        string `json:"code"`
	CountryName string `json:"country_name"`
}

type ConfirmCodeRequest struct {
	Code     string `json:"code" binding:"required,len=6"`
	Username string `json:"username" binding:"required"`
}

type ConfirmCodeResponse struct {
	Confirmed bool   `json:"confirmed"`
	Username  string `json:"username"`
}

type NewPasswordRequest struct {
	Username    string `json:"username" binding:"required,min=5"`
	PasswordOne string `json:"password_one" binding:"required,password"`
	PasswordTwo string `json:"password_two" binding:"required,password"`
}

type NewPasswordResponse struct {
	Updated bool `json:"updated"`
}

type UpdatePhoneNumberRequest struct {
	Username string `json:"username" binding:"required,min=5"`
}

type UpdatePhoneNumberResponse struct {
	Updated     bool   `json:"updated"`
	PhoneNumber string `json:"phone_number"`
}

type JoinUserVerifyEmail struct {
	Email string `json:"email" binding:"required"`
}

type UpdateCodePhoneResponse struct {
	PhoneNumber          string `json:"phone_number"`
	AccessToken          string `json:"access_token"`
	AccessTokenExpiresAt string `json:"access_token_expires_at"`
	RefreshToken         string `json:"refresh_token"`
}

type UpdateCodeEmailResponse struct {
	Email                string `json:"email"`
	AccessToken          string `json:"access_token"`
	AccessTokenExpiresAt string `json:"access_token_expires_at"`
	RefreshToken         string `json:"refresh_token"`
}

type GetUserParams struct {
	Email         string          `json:"email"`
	FireFight     string          `json:"fire_fight"`
	FirstName     string          `json:"first_name"`
	LastName      string          `json:"last_name"`
	Currency      string          `json:"currency"`
	ProfilePhoto  string          `json:"profile_photo"`
	IsHost        bool            `json:"is_host"`
	WishlistList  ListWishlistRes `json:"wishlist_list"`
	HasIncomplete bool            `json:"has_incomplete"`
}

type CreateUserAPNDetailParams struct {
	PublicID            string `json:"public_id"`
	Name                string `json:"name"`
	Model               string `json:"model"`
	IdentifierForVendor string `json:"identifier_for_vendor"`
	Token               string `json:"token"`
}
