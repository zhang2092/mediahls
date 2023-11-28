package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Different types of error returned by the VerifyToken function
var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")
)

// Payload contains the payload data of the token
type Payload struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// NewPayload creates a new token payload with a specific username and duration
func NewPayload(id string, username string, duration time.Duration) *Payload {
	payload := &Payload{
		ID:       id,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)), // 过期时间
			IssuedAt:  jwt.NewNumericDate(time.Now()),               // 签发时间
			NotBefore: jwt.NewNumericDate(time.Now()),               // 生效时间
		},
	}
	return payload
}

// Valid checks if the token payload is valid or not
func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiresAt.Time) {
		return ErrExpiredToken
	}
	return nil
}
