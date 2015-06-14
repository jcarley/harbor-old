package models

import "time"

type User struct {
	ID           string
	Username     string
	Fullname     string
	PasswordHash string
	PasswordSalt string
	IsDisabled   bool
}

type UserSession struct {
	SessionKey   string
	UserID       string
	LoginTime    time.Time
	LastSeenTime time.Time
}
