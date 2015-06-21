package models

import (
	"errors"
	"time"

	"github.com/gorilla/securecookie"
)

var (
	ErrUserSessionNotFound = errors.New("User session not found")
)

type UserSessionStore interface {
	Get(sessionKey string) (UserSession, error)
	Save(session UserSession) error
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

type MemoryUserSessionStore struct {
	db map[string]UserSession
}

func NewMemoryUserSessionStore() MemoryUserSessionStore {
	return MemoryUserSessionStore{
		db: make(map[string]UserSession),
	}
}

func (this *MemoryUserSessionStore) Get(sessionKey string) (UserSession, error) {
	user_session, exists := this.db[sessionKey]
	if exists {
		return user_session, nil
	} else {
		return UserSession{}, ErrUserSessionNotFound
	}
}

func (this *MemoryUserSessionStore) Save(session UserSession) error {
	this.db[session.SessionKey] = session
	return nil
}
