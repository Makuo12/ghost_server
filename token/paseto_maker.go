package token

import (
	"fmt"
	"time"

	"github.com/o1egl/paseto"
	"golang.org/x/crypto/chacha20poly1305"
)

//PasetoMaker is a Paseto token maker

type PasetoMaker struct {
	paseto *paseto.V2
	//we would use symmetric encryption here to encrypt the payload
	symmetricKey []byte
}

//This function would create a new paseto maker instance
func NewPasetoMaker(symmetricKey string) (Maker, error){
	//we check to know if the key size if valid
	if len(symmetricKey) != chacha20poly1305.KeySize{
		return nil, fmt.Errorf("invalid key size: must be exactly %d characters", chacha20poly1305.KeySize) 
	}

	maker := &PasetoMaker{
		paseto: paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
	}
	return maker, nil
}

//CreateTokens creates a new token for a specific email and duration
func (maker *PasetoMaker) CreateTokens(username string, duration time.Duration) (string, *Payload, error){
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", payload, err 
	}
	token, err := maker.paseto.Encrypt(maker.symmetricKey, payload, nil)
	return token, payload, err
}  

//VerifyToken checks if the input token is valid or not
func (maker *PasetoMaker) VerifyToken(token string) (*Payload, error){
	payload := &Payload{}

	err := maker.paseto.Decrypt(token, maker.symmetricKey, payload, nil)
	if err != nil {
		return nil, ErrInvalidToken
	}
	err = payload.Valid()
	if err != nil {
		return nil, err
	}
	return payload, nil
}