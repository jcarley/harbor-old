package models

import "time"

type User struct {
	Id           string    `gorethink:"id,omitempty", json:"id,omitempty"`
	Username     string    `gorethink:"username",     json:"username"`
	PasswordHash string    `gorethink:"passwordhash", json:"password"`
	PasswordSalt string    `gorethink:"passwordsalt", json:"salt"`
	IsDisabled   bool      `gorethink:"is_disabled",  json:is_disabled`
	Created      time.Time `gorethink:"created",      json:"created"`
}
