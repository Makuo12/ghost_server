package api

type AccountChangeParams struct {
	Type string `json:"type"`
}

type AccountChangeRes struct {
	Username string `json:"username"`
}

type VerifyAccountChangeParams struct {
	Username string `json:"username"`
	Code     string `json:"code"`
	Type string `json:"type"`
}


