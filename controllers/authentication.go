package controllers

import (
	"fmt"
	"net/http"
	"time"

	"golang.org/x/crypto/scrypt"

	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/jcarley/harbor/models"
	"github.com/jcarley/harbor/utils"
	"github.com/jcarley/harbor/view_models"
	"github.com/jcarley/harbor/web"
)

type AuthenticationController struct {
	Context web.Context
}

func NewAuthenticationController(ctx web.Context) *AuthenticationController {
	return &AuthenticationController{
		Context: ctx,
	}
}

func (this *AuthenticationController) Register(router *mux.Router) {
	router.HandleFunc("/login", this.login).Methods("POST")
	router.HandleFunc("/register", this.register).Methods("POST")
}

func (this *AuthenticationController) login(w http.ResponseWriter, req *http.Request) {

	auth_request := view_models.AuthRequest{}
	if err := utils.Decode(req.Body, &auth_request); err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user := new(models.User)
	if err := user.FindByUsername(auth_request.Email); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if this.authorize(auth_request, user) {
		user_session := models.NewUserSession(user)
		user_session.Save()

		auth_response := view_models.AuthResponse{
			Token:   user_session.SessionKey,
			Expires: time.Now().Add(time.Hour * 24 * 7),
		}

		if err := utils.Encode(w, &auth_response); err != nil {
			fmt.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	} else {
		error_message := `{"message": "%s"}`
		http.Error(w, fmt.Sprintf(error_message, "Username or password incorrect"), http.StatusForbidden)
	}

}

func (this *AuthenticationController) register(w http.ResponseWriter, req *http.Request) {

	auth_request := view_models.AuthRequest{}
	if err := utils.Decode(req.Body, &auth_request); err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// TODO: find a better way of sending a json error message
	if err := auth_request.Valid(); err != nil {
		error_message := `{"message": "%s"}`
		http.Error(w, fmt.Sprintf(error_message, err.Error()), http.StatusBadRequest)
		return
	}

	salt := securecookie.GenerateRandomKey(32)
	dk := this.hash_password([]byte(auth_request.Password), salt)

	password_hash := fmt.Sprintf("%x", dk)
	password_salt := fmt.Sprintf("%x", salt)

	user := models.NewUser(auth_request.Email, password_hash, password_salt, false)

	if err := user.Save(); err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.Encode(w, user)
}

func (this *AuthenticationController) authorize(auth_request view_models.AuthRequest, user *models.User) bool {
	password_hash := user.PasswordHash
	salt := user.PasswordSalt
	hash := this.hash_password([]byte(auth_request.Password), []byte(salt))
	return password_hash == fmt.Sprintf("%x", hash)
}

func (this *AuthenticationController) hash_password(password []byte, salt []byte) []byte {
	dk, _ := scrypt.Key(password, salt, 16384, 8, 1, 32)
	return dk
}
