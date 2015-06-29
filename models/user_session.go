package models

import (
	"errors"
	"fmt"
	"time"

	r "github.com/dancannon/gorethink"
	"github.com/gorilla/securecookie"
	"github.com/jcarley/harbor/service"
)

var (
	ErrUserSessionNotFound = errors.New("User session not found")
)

type UserSession struct {
	Id           string    `gorethink:"id,omitempty"`
	SessionKey   string    `gorethink:"session_key"`
	UserId       string    `gorethink:"user_id"`
	LoginTime    time.Time `gorethink:"login_time"`
	LastSeenTime time.Time `gorethink:"last_seen_time"`
	Created      time.Time `gorethink:"created"`
	Updated      time.Time `gorethink:"updated"`
}

func NewUserSession(user *User) *UserSession {
	session_key := fmt.Sprintf("%x", securecookie.GenerateRandomKey(16))
	session := UserSession{
		SessionKey:   session_key,
		UserId:       user.Id,
		LoginTime:    time.Now(),
		LastSeenTime: time.Now(),
	}
	session.Created = time.Now()
	session.Updated = time.Now()
	return &session
}

func (this *UserSession) Save() error {
	res, err := this.coll().Insert(this).RunWrite(service.Session())
	if err != nil {
		return err
	}
	this.Id = res.GeneratedKeys[0]

	return nil
}

func (this *UserSession) FindBySessionKey(session_key string) error {
	res, err := this.coll().Filter(r.Row.Field("session_key").Eq(session_key)).Run(service.Session())
	if err != nil {
		return ErrUserSessionNotFound
	}
	defer res.Close()
	return res.One(this)
}

func (this *UserSession) coll() r.Term {
	return r.Db("harbor").Table("user_sessions")
}
