package middleware

import (
	"net/http"

	"gopkg.in/mgo.v2/bson"

	"github.com/codegangsta/negroni"
	"github.com/jcarley/harbor/models"
)

type Authenticator struct {
	currentUser utils.CurrentUserAccessor
}

func NewAuthenticator(currentUser utils.CurrentUserAccessor) *Authenticator {
	return &Authenticator{currentUser}
}

func (a *Authenticator) Middleware() negroni.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

		session_key := r.Header.Get("Authorization")
		user_session := new(models.UserSession)

		if err := user_session.FindBySessionKey(session_key); err != nil {

		}

		user := new(models.User)
		user.FindById(user_session.UserId)

		if bson.IsObjectIdHex(userID) && user.FindByID(bson.ObjectIdHex(userID), a.database.Get(r)) == nil {
			a.currentUser.Set(r, user)
		} else {
			a.currentUser.Clear(r)
		}
		next(rw, r)
	}
}
