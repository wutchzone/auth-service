package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/wutchzone/auth-service/api"
)

const apiRoot string = "http://localhost:7080/auth/"

var c = http.Client{
	Timeout: 10 * time.Second,
}

var u, _ = json.Marshal(parser.UserJSON{
	Name:     "test",
	Password: "testtest",
	Email:    "test@test.test",
})

func TestHandleRegister(t *testing.T) {
	// Make POST for user registration
	r, errRegister := c.Post(apiRoot+"register", "application/json", strings.NewReader(string(u)))
	if errRegister != nil {
		t.Errorf("Failed post to server for registration, got: %v", errRegister)
		return
	}
	if r.StatusCode == http.StatusOK {
		fmt.Printf("User registered succesfully\n")
	} else {
		t.Errorf("Registration error expected 200, but got %v\n", r.StatusCode)
	}
}

func TestHandleLogin(t *testing.T) {
	r, err := c.Post(apiRoot+"login", "application/json", strings.NewReader(string(u)))

	if err != nil {
		t.Errorf("Failed post to server for login, got: %v", err)
		return
	}

	// Decode POST response for user login
	var uid parser.UUIDJSON
	errDecode := json.NewDecoder(r.Body).Decode(&uid)
	if errDecode != nil {
		t.Errorf("Failed to decode registration response got: %v", errDecode)
		return
	}
	fmt.Printf("User registered, with token: %v\n", uid.UUID)
}
