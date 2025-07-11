package token

import (
	"time"

	err "github.com/MyProject273/Join_Love/pkg/error"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenType byte
type Payload struct {
	ID        uuid.UUID `json:"id"`
	Type      TokenType `json:"token_type"`
	UserID    string    `json:"user_id"`
	Role      string    `json:"role"`
	IssueAt   time.Time `json:"issue_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

// GetAudience implements jwt.Claims.
func (payload *Payload) GetAudience() (jwt.ClaimStrings, error) {
	panic("unimplemented")
}

// GetExpirationTime implements jwt.Claims.
func (payload *Payload) GetExpirationTime() (*jwt.NumericDate, error) {
	panic("unimplemented")
}

// GetIssuedAt implements jwt.Claims.
func (payload *Payload) GetIssuedAt() (*jwt.NumericDate, error) {
	panic("unimplemented")
}

// GetIssuer implements jwt.Claims.
func (payload *Payload) GetIssuer() (string, error) {
	panic("unimplemented")
}

// GetNotBefore implements jwt.Claims.
func (payload *Payload) GetNotBefore() (*jwt.NumericDate, error) {
	panic("unimplemented")
}

// GetSubject implements jwt.Claims.
func (payload *Payload) GetSubject() (string, error) {
	panic("unimplemented")
}

func NewPayload(user_id string, role string, duration time.Duration, tokenType TokenType) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := &Payload{
		ID:        tokenID,
		Type:      tokenType,
		UserID:    user_id,
		Role:      role,
		IssueAt:   time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}

	return payload, nil
}

func (payload *Payload) Valid(tokenType TokenType) error {
	if payload.Type != tokenType {
		return err.ErrInvalidToken
	}
	if time.Now().After(payload.ExpiredAt) {
		return err.ErrExpiredToken
	}
	return nil
}
