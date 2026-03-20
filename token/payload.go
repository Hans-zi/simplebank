package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrExpiredToken = errors.New("Expired-Token")
	ErrInvalidToken = errors.New("Invalid-Token")
)

type Payload struct {
	ID uuid.UUID `json:"id"`

	Username string `json:"username"`

	IssuedAt  time.Time `json:"issued_at"`
	ExpiresAt time.Time `json:"expires_at"`

	jwt.RegisteredClaims
}

func NewPayload(username string, duration time.Duration) (*Payload, error) {
	tokenId, err := uuid.NewRandom()
	if err != nil {
		return &Payload{}, err
	}
	return &Payload{
		ID:        tokenId,
		Username:  username,
		IssuedAt:  time.Now(),
		ExpiresAt: time.Now().Add(duration),
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "Inorin",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		},
	}, nil
}

func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiresAt) {
		return ErrExpiredToken
	}
	return nil
}

//func NewJWTPayload(username string, duration time.Duration) (*JWTPayload, error) {
//	tokenId, err := uuid.NewRandom()
//	if err != nil {
//		return &JWTPayload{}, err
//	}
//	return &JWTPayload{
//		ID:       tokenId,
//		Username: username,
//		RegisteredClaims: jwt.RegisteredClaims{
//			Issuer:    "Inorin",
//			IssuedAt:  jwt.NewNumericDate(time.Now()),
//			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
//		},
//	}, nil
//}
