package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	r "github.com/dancannon/gorethink"
	"github.com/jcarley/harbor/models"
	"github.com/jcarley/harbor/service"
	. "github.com/onsi/gomega"
)

const (
	checkMark = "\u2713"
	ballotX   = "\u2717"
)

// salt := "d61162e555f68c3151133351fc908d688aa2bb1e5bab958859290c443eeec0bc"

func CreateTestUser() (models.User, error) {
	user := models.User{
		Username:     "jeff.carley@gmail.com",
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

func UserCount() (int, error) {
	res, err := r.Db("harbor").Table("users").Count().Run(service.Session())
	if err != nil {
		return 0, err
	}
	defer res.Close()

	var count int
	res.One(&count)

	return count, nil
}

func DeleteUser(user models.User) {
	_, err := r.Db("harbor").Table("users").Get(user.Id).Delete().Run(service.Session())
	if err != nil {
		fmt.Println(err)
		return
	}
}

func DeleteAll() {
	_, err := r.Db("harbor").Table("users").Delete().Run(service.Session())
	if err != nil {
		fmt.Println(err)
		return
	}
}

func TestRegister(t *testing.T) {
	RegisterTestingT(t)
	defer DeleteAll()

	// Test Setup
	current_count, err := UserCount()
	if err != nil {
		t.Fatal(err)
	}

	// Send register request
	body :=
		`{
			"email":"jcarley@gmail.com",
			"password":"secret"
		}`

	req, err := http.NewRequest("POST", "/register", strings.NewReader(body))
	Expect(err).ShouldNot(HaveOccurred(), "Should be able to create a request")

	req.Header.Add("Content-Type", "application/json")

	// Run test
	w := httptest.NewRecorder()
	router := NewRouter()
	router.ServeHTTP(w, req)

	Expect(w.Code).To(Equal(200), "Should receive 200 status")

	// Assert if a user was added
	actual_count, err := UserCount()
	if err != nil {
		t.Fatal(err)
	}

	Expect(actual_count > current_count).To(BeTrue(), "Expected user count to be greater than %d, but was %d", current_count, actual_count)
}

func TestUserLogin(t *testing.T) {
	RegisterTestingT(t)

	user, _ := CreateTestUser()
	defer DeleteUser(user)

	// Send register request
	body := fmt.Sprintf(`{"email":"%s", "password":"%s"}`, user.Username, "password")

	req, err := http.NewRequest("POST", "/login", strings.NewReader(body))
	Expect(err).ShouldNot(HaveOccurred(), "Should be able to create a request")

	req.Header.Add("Content-Type", "application/json")

	w := httptest.NewRecorder()

	router := NewRouter()
	router.ServeHTTP(w, req)

	Expect(w.Code).To(Equal(200), "Should receive 200 status")

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)

	Expect(err).ShouldNot(HaveOccurred())
	Expect(response["token"]).ToNot(BeNil())
	Expect(response["expires"]).ToNot(BeNil())
}
