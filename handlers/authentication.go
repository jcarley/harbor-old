package handlers

import (
	"fmt"
	"net/http"

	"golang.org/x/crypto/scrypt"

	r "github.com/dancannon/gorethink"

	"github.com/gorilla/securecookie"
	"github.com/jcarley/harbor/models"
	"github.com/jcarley/harbor/service"
)

func signIn(ctx context, w http.ResponseWriter, req *http.Request) {

	email := req.PostFormValue("inputEmail")
	password := req.PostFormValue("inputPassword")

	res, err := r.Db("harbor").Table("users").Filter(r.Row.Field("username").Eq(email)).Run(service.Session())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer res.Close()

	user := models.User{}
	err = res.One(&user)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// fmt.Printf("User: %+v\n", user)

	password_hash := user.PasswordHash
	salt := user.PasswordSalt

	dk, _ := scrypt.Key([]byte(password), []byte(salt), 16384, 8, 1, 32)

	fmt.Printf("%s\n", password_hash)
	fmt.Printf("%x\n", dk)

	if password_hash == fmt.Sprintf("%x", dk) {
		user_session := models.NewUserSession(user)

		session, _ := ctx.SessionStore().Get(req, "login")
		session.Values["username"] = user.Username
		session.Values["sessionKey"] = user_session.SessionKey
		session.Save(req, w)

		w.Write([]byte("success"))
	} else {
		http.Redirect(w, req, "index.html", 301)
	}

}

func register(ctx context, w http.ResponseWriter, r *http.Request) {

	email := r.PostFormValue("inputEmail")
	password := r.PostFormValue("inputPassword")

	salt := securecookie.GenerateRandomKey(32)
	// salt := "d61162e555f68c3151133351fc908d688aa2bb1e5bab958859290c443eeec0bc"
	dk, _ := scrypt.Key([]byte(password), salt, 16384, 8, 1, 32)

	password_hash := fmt.Sprintf("%x", dk)
	password_salt := fmt.Sprintf("%x", salt)

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
