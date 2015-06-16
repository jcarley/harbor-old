package models

import (
	"time"

	"github.com/gorilla/securecookie"
)

type User struct {
	Id           string    `gorethink:"id,omitempty"`
	Username     string    `gorethink:"username"`
	Fullname     string    `gorethink:"fullname"`
	PasswordHash string    `gorethink:"passwordhash"`
	PasswordSalt string    `gorethink:"passwordsalt"`
	IsDisabled   bool      `gorethink:"is_disabled"`
	Created      time.Time `gorethink:"created"`
}

type UserSession struct {
	Id           string    `gorethink:"id,omitempty"`
	SessionKey   string    `gorethink:"session_key"`
	UserId       string    `gorethink:"user_id"`
	LoginTime    time.Time `gorethink:"login_time"`
	LastSeenTime time.Time `gorethink:"last_seen_time"`
}

func NewUserSession(user User) UserSession {
	return UserSession{
		SessionKey:   string(securecookie.GenerateRandomKey(16)),
		UserId:       user.Id,
		LoginTime:    time.Now(),
		LastSeenTime: time.Now(),
	}
}
