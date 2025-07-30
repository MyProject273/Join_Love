package apperr

import "errors"

var (
	// ==== User ====
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrUserBlocked        = errors.New("user is blocked")
	ErrUserInactive       = errors.New("user is inactive")
	ErrCreateUser         = errors.New("failed to create user")

	// ==== Password ====
	ErrHashPassword     = errors.New("failed to hash password")
	ErrMismatchPassword = errors.New("password is mismatch")

	// ==== Token ====
	ErrCreateAccessToken  = errors.New("failed to create access token")
	ErrCreateRefreshToken = errors.New("failed to create refresh token")
	ErrInvalidToken       = errors.New("invalid token")
	ErrExpiredToken       = errors.New("token has expired")
	ErrUnauthorized       = errors.New("unauthorized access")

	// ==== Session ====
	ErrCreateSession   = errors.New("failed to create session")
	ErrSessionNotFound = errors.New("session not found")
	ErrSessionBlocked  = errors.New("session is blocked")
	ErrSessionExpired  = errors.New("session expired")

	// ==== Auth ====
	ErrLoginFailed       = errors.New("login failed")
	ErrPermissionDenied  = errors.New("permission denied")
	ErrInvalidAuthHeader = errors.New("invalid authorization header")
	ErrMissingAuthToken  = errors.New("missing authorization token")

	// ==== Validation ====
	ErrValidationFailed     = errors.New("validation failed")
	ErrMissingRequiredField = errors.New("missing required field")

	// ==== Database ====
	ErrDBQueryFailed   = errors.New("database query failed")
	ErrDBTransaction   = errors.New("database transaction failed")
	ErrRecordNotFound  = errors.New("record not found")
	ErrDuplicateRecord = errors.New("duplicate record")

	// ==== File ====
	ErrFileTooLarge     = errors.New("file size exceeds limit")
	ErrFileUnsupported  = errors.New("unsupported file format")
	ErrFileUploadFailed = errors.New("failed to upload file")

	// ==== Mail ====
	ErrSendEmailFailed    = errors.New("failed to send email")
	ErrInvalidEmailFormat = errors.New("invalid email format")

	// ==== Internal ====
	ErrInternalServer = errors.New("internal server error")
)
