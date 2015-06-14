package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

const (
	checkMark = "\u2713"
	ballotX   = "\u2717"
)

func TestUserSignIn(t *testing.T) {

	formData := url.Values{}
	formData.Add("inputEmail", "jeff.carley@gmail.com")
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
