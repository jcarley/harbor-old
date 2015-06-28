package view_models

import "errors"

var (
	ErrMissingEmailAddress = errors.New("Email address is required")
	ErrMissingPassword     = errors.New("Password is required")
)

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (this *AuthRequest) Valid() error {
	if this.Email == "" {
		return ErrMissingEmailAddress
	}

	if this.Password == "" {
		return ErrMissingPassword
	}

	return nil
}
