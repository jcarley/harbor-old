package models

import (
	"time"

	r "github.com/dancannon/gorethink"
	"github.com/jcarley/harbor/service"
)

type User struct {
	Id           string    `gorethink:"id,omitempty", json:"id,omitempty"`
	Username     string    `gorethink:"username",     json:"username"`
	PasswordHash string    `gorethink:"passwordhash", json:"password"`
	PasswordSalt string    `gorethink:"passwordsalt", json:"salt"`
	IsDisabled   bool      `gorethink:"is_disabled",  json:is_disabled`
	Created      time.Time `gorethink:"created",      json:"created"`
	Updated      time.Time `gorethink:"updated",      json:"updated"`
}

func NewUser(username, password_hash, password_salt string, is_disabled bool) *User {
	user := User{
		Username:     username,
		PasswordHash: password_hash,
		PasswordSalt: password_salt,
		IsDisabled:   is_disabled,
	}
	user.Created = time.Now()
	user.Updated = time.Now()
	return &user
}

func (this *User) Save() error {
	res, err := r.Db("harbor").Table("users").Insert(this).RunWrite(service.Session())
	if err != nil {
		return err
	}
	this.Id = res.GeneratedKeys[0]

	return nil
}

func (this *User) FindByUsername(username string) error {
	res, err := r.Db("harbor").Table("users").Filter(r.Row.Field("username").Eq(username)).Run(service.Session())
	if err != nil {
		return err
	}
	defer res.Close()
	return res.One(this)
}
