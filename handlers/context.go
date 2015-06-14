package handlers

import (
	"github.com/gorilla/sessions"
	"github.com/jcarley/harbor/service"
)

type context interface {
	Hub() *service.Hub
	SessionStore() sessions.Store
}

type ctx struct {
	hub          *service.Hub
	sessionStore sessions.Store
}

func (this ctx) Hub() *service.Hub {
	return this.hub
}

func (this ctx) SessionStore() sessions.Store {
	return this.sessionStore
}
