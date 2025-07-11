package token

import "time"

// Maker is an interface for managing tokens
type Maker interface {
	// CreateToken create a new token for a specific user_id and duration
	CreateToken(user_id string, role string, duration time.Duration, tokenType TokenType) (string, *Payload, error)

	// VerifyToken check if the token is valid or not
	VerifyToken(token string, tokenType TokenType) (*Payload, error)
}
