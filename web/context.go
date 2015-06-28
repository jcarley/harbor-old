package web

import (
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

type Context interface {
	SessionStore() sessions.Store
}

type ContextImpl struct {
	store sessions.Store
}

func NewContext(sessionStore sessions.Store) *ContextImpl {
	return &ContextImpl{
		store: sessionStore,
	}
}

func (this *ContextImpl) SessionStore() sessions.Store {
	if this.store == nil {
		return nil
	}
	return this.store
}

func NewCookieStore() *sessions.CookieStore {
	authKey := securecookie.GenerateRandomKey(64)
	encryptionKey := securecookie.GenerateRandomKey(32)
	return sessions.NewCookieStore(
		authKey,
		encryptionKey,
	)
}
