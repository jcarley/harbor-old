package handlers

import "github.com/gorilla/sessions"

type context interface {
	SessionStore() sessions.Store
}

type ctx struct {
	sessionStore sessions.Store
}

func (this ctx) SessionStore() sessions.Store {
	return this.sessionStore
}
