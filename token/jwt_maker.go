package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const minSecretKeySize = 32

//JWTMaker is json web Token maker
type JWTMaker struct{
	secretKey string
}

//NewJWTMaker creates a new JWTMaker 

func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid key size: must be at least %d characters", minSecretKeySize)
	}
	return &JWTMaker{secretKey: secretKey}, nil
}

//CreateTokens creates a new token for a specific username and duration
func (maker *JWTMaker) CreateTokens(username string, duration time.Duration) (string, *Payload, error){
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", payload, err
	}
	//we create a new jwttoken usinf jwt.NewWithClaims function
	//first argument is the signing so we use HS256
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256,payload)
	token, err := jwtToken.SignedString([]byte(maker.secretKey))
	return token, payload, err
}

//VerifyToken checks if the input token is valid or not
func (maker *JWTMaker) VerifyToken(token string) (*Payload, error){
	// keyFunc is our key function which receive a parse but unverfied token. So we should verfiy its header 
	//to make sure that the signing algorithm matches what you normally use to sign the token
	//if it matches, you return the key
	keyFunc := func(token *jwt.Token) (interface{}, error){
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok{
			return nil, ErrInvalidToken
		}
		return []byte(maker.secretKey), nil
	}
	//first we have to parse the token 
	jwtToken, err := jwt.ParseWithClaims(token, &Payload{},keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ErrExpiredToken){
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}
	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}
	return payload, nil
}
