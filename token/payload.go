package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
)
//DIfferent types of error returned by the verifyToken function
var (
	ErrExpiredToken = errors.New("token has expired")
	ErrInvalidToken = errors.New("token is invalid")
)
//Payload contains the payload data of the token

type Payload struct {
	ID       uuid.UUID `json:"id"`
	Username string `json:"username"`
	IssuedAt time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

//NewPayload creates a new token payload with a specific username and duration
func NewPayload(username string, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	payload := &Payload{
		ID: tokenID,
		Username: username,
		IssuedAt: time.Now(),
		ExpiredAt: time.Now().Add(duration),
		
	}
	return payload, nil
}
func (payload *Payload) Valid() error{
	//first we check the time of the time of the token 
	if time.Now().After(payload.ExpiredAt){
		return ErrExpiredToken
	}
	return nil
}