package mail

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type BoolField struct {
	Value bool   `json:"value"`
	Text  string `json:"text"`
}

type AbstractEmailVerifyResult struct {
	Email             string    `json:"email"`
	Autocorrect       string    `json:"autocorrect"`
	Deliverability    string    `json:"deliverability"`
	QualityScore      string    `json:"quality_score"`
	IsValidFormat     BoolField `json:"is_valid_format"`
	IsFreeEmail       BoolField `json:"is_free_email"`
	IsDisposableEmail BoolField `json:"is_disposable_email"`
	IsRoleEmail       BoolField `json:"is_role_email"`
	IsCatchAllEmail   BoolField `json:"is_catchall_email"`
	IsMxFound         BoolField `json:"is_mx_found"`
	IsSmtpValid       BoolField `json:"is_smtp_valid"`
}

func ValidateEmailAddress(email, apiKey string) (*AbstractEmailVerifyResult, error) {
	url := fmt.Sprintf("https://emailvalidation.abstractapi.com/v1/?api_key=%s&email=%s", apiKey, email)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to call abstractapi: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("abstractapi returned status: %s", resp.Status)
	}

	var result AbstractEmailVerifyResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}
