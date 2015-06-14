package handlers

import (
	"fmt"
	"net/http"

	"golang.org/x/crypto/scrypt"

	"github.com/gorilla/securecookie"
	"github.com/jcarley/harbor/models"
)

func signIn(ctx context, w http.ResponseWriter, r *http.Request) {

	email := r.PostFormValue("inputEmail")
	password := r.PostFormValue("inputPassword")

	// salt := securecookie.GenerateRandomKey(32)
	salt := "d61162e555f68c3151133351fc908d688aa2bb1e5bab958859290c443eeec0bc"
	dk, _ := scrypt.Key([]byte(password), []byte(salt), 16384, 8, 1, 32)

	password_hash := fmt.Sprintf("%x", dk)
	password_salt := salt

	user := models.User{
		Username:     email,
		PasswordHash: password_hash,
		PasswordSalt: password_salt,
	}

	fmt.Printf("%+v\n", user)

	if email == "jeff.carley@gmail.com" && password == "password" {
		session, _ := ctx.SessionStore().Get(r, "login")

		session.Values["username"] = email
		session.Values["sessionKey"] = string(securecookie.GenerateRandomKey(16))
		session.Save(r, w)

		w.Write([]byte("success"))
	} else {
		http.Redirect(w, r, "index.html", 301)
	}

}
