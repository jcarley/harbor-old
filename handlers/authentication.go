package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"golang.org/x/crypto/scrypt"

	r "github.com/dancannon/gorethink"
	"github.com/gorilla/securecookie"
	"github.com/jcarley/harbor/models"
	"github.com/jcarley/harbor/service"
)

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token   string    `json:"token,omitempty"`
	Expires time.Time `json:"expires,omitempty"`
}

func Encode(writer io.Writer, data interface{}) error {
	encoder := json.NewEncoder(writer)
	return encoder.Encode(data)
}

func Decode(reader io.Reader, data interface{}) error {
	decoder := json.NewDecoder(reader)
	return decoder.Decode(data)
}

func login(ctx context, w http.ResponseWriter, req *http.Request) {

	auth_request := AuthRequest{}
	if err := Decode(req.Body, &auth_request); err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res, err := r.Db("harbor").Table("users").Filter(r.Row.Field("username").Eq(auth_request.Email)).Run(service.Session())
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

	password_hash := user.PasswordHash
	salt := user.PasswordSalt

	dk, _ := scrypt.Key([]byte(auth_request.Password), []byte(salt), 16384, 8, 1, 32)

	if password_hash == fmt.Sprintf("%x", dk) {
		user_session := models.NewUserSession(user)

		session, _ := ctx.SessionStore().Get(req, "login")
		session.Values["username"] = user.Username
		session.Values["sessionKey"] = user_session.SessionKey
		session.Save(req, w)

		auth_response := AuthResponse{
			Token:   fmt.Sprintf("%x", user_session.SessionKey),
			Expires: time.Now().Add(time.Hour * 24 * 7),
		}

		if err := Encode(w, &auth_response); err != nil {
			fmt.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	} else {
		http.Redirect(w, req, "index.html", 301)
	}

}

func register(ctx context, w http.ResponseWriter, req *http.Request) {

	auth_request := AuthRequest{}

	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&auth_request)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	salt := securecookie.GenerateRandomKey(32)
	dk, _ := scrypt.Key([]byte(auth_request.Password), salt, 16384, 8, 1, 32)

	password_hash := fmt.Sprintf("%x", dk)
	password_salt := fmt.Sprintf("%x", salt)

	user := models.User{
		Username:     auth_request.Email,
		PasswordHash: password_hash,
		PasswordSalt: password_salt,
		IsDisabled:   false,
	}
	user.Created = time.Now()

	res, err := r.Db("harbor").Table("users").Insert(user).RunWrite(service.Session())
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user.Id = res.GeneratedKeys[0]

	encoder := json.NewEncoder(w)
	encoder.Encode(user)
}
