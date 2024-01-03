package token

import "time"

//maker is an interface for managing tokens

type Maker interface {
	//CreateTokens creates a new token for a specific username and duration
	CreateTokens(username string, duration time.Duration) (string, *Payload, error)

	//VerifyToken checks if the input token is valid or not
	VerifyToken(token string) (*Payload, error)
}