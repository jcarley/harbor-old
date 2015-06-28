package view_models

import "time"

type AuthResponse struct {
	Token   string    `json:"token,omitempty"`
	Expires time.Time `json:"expires,omitempty"`
}
