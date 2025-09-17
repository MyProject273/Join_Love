package auth_dto

import (
	"time"
)

type LoginReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type Token struct {
	AccessToken           string    `json:"access_token" binding:"required"`
	RefreshToken          string    `json:"refresh_token" binding:"required"`
	AccessTokenExpiresIn  time.Time `json:"access_token_expires_in"`
	RefreshTokenExpiresIn time.Time `json:"refresh_token_expires_in"`
}

type LoginRes struct {
	UserID    string `json:"user_id"`
	Email     string `json:"email"`
	UserName  string `json:"user_name"`
	SessionID string `json:"session_id"`
	Token     Token  `json:"token"`
}

type SignupReq struct {
	UserName string `json:"user_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type SignupRes struct {
	UserID   string `json:"user_id"`
	UserName string `json:"user_name"`
	Email    string `json:"email"`
}

type VerifyEmailReq struct {
	VerifyEmailID string `form:"email_id" binding:"required"`
	SecretCode    string `form:"secret_code" binding:"required"`
}

type VerifyEmailRes struct {
	IsVerified string `json:"is_verified"`
}
