package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	r "github.com/dancannon/gorethink"
	"github.com/jcarley/harbor/models"
	"github.com/jcarley/harbor/service"
)

const (
	checkMark = "\u2713"
	ballotX   = "\u2717"
)

// salt := "d61162e555f68c3151133351fc908d688aa2bb1e5bab958859290c443eeec0bc"

func CreateTestUser() (models.User, error) {
	user := models.User{
		Username:     "jeff.carley@gmail.com",
		Fullname:     "Jefferson Carley",
		PasswordHash: "0b2f219acb4b0cd9c5181f77ed41484fc286d0c11878005be2d4e7695255e2dc",
		PasswordSalt: "d61162e555f68c3151133351fc908d688aa2bb1e5bab958859290c443eeec0bc",
		IsDisabled:   false,
	}
	user.Created = time.Now()

	res, err := r.Db("harbor").Table("users").Insert(user).RunWrite(service.Session())
	if err != nil {
		return models.User{}, err
	}
	user.Id = res.GeneratedKeys[0]

	return user, nil
}

func DeleteUser(user models.User) {
	_, err := r.Db("harbor").Table("users").Get(user.Id).Delete().Run(service.Session())
	if err != nil {
		fmt.Println(err)
		return
	}
}

func TestUserSignIn(t *testing.T) {

	user, _ := CreateTestUser()
	defer DeleteUser(user)

	formData := url.Values{}
	formData.Add("inputEmail", user.Username)
	formData.Add("inputPassword", "password")
	body := strings.NewReader(formData.Encode())

	req, err := http.NewRequest("POST", "/signin", body)
	if err != nil {
		t.Fatal("Should be able to create a request", ballotX, err)
	}
	t.Log("Should be able to create a request", checkMark)

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()

	router := NewRouter()
	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatal("Should receive 200 status", ballotX, w.Code)
	}
	t.Log("Should receive 200 status", checkMark)

	cookie := w.Header()["Set-Cookie"]
	if len(cookie) == 0 {
		t.Fatal("Should have auth cookies", ballotX)
	}
	t.Log("Should have auth cookies", checkMark)

}
