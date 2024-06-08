package api

type ListFirebasePhoto struct {
	Photos []string `json:"photos"`
}

type UpdateFirebasePhotoParams struct {
	Actual string `json:"actual"`
	Update string `json:"update"`
}

type UpdateFirebasePhoto struct {
	Photos []UpdateFirebasePhotoParams `json:"photos"`
}


