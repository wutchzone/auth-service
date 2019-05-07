package main

import (
	"encoding/json"
	"github.com/wutchzone/api-response"
	"github.com/wutchzone/auth-service/pkg/accountdb"
	"github.com/wutchzone/auth-service/pkg/decoder"
	"github.com/wutchzone/auth-service/pkg/sessiondb"
	"net/http"
	"time"
)

// HandleLogin route
func HandleLogin(w http.ResponseWriter, r *http.Request) {
	// Decode user
	u, errDecode := decodeUser(r)
	if errDecode != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ := response.CreateResponse(w, response.ResponseError, nil, "")
		return
	}

	// Check if user exists in userDB
	cu, errSearch := UserDB.GetAccount(u.User)
	if errSearch != nil || cu.ComparePswdAndHash((*u).Password) != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_ := response.CreateResponse(w, response.ResponseError, nil, "password or username is not correct")
		return
	}

	// Save user to sessionDB
	si := sessiondb.NewSessionItem(cu.Role)
	uid := sessiondb.NewSessionKey()
	dur := time.Minute * 20
	if err := SessionDB.SetRecord(uid, *si, dur); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ := response.CreateResponse(w, response.ResponseError, nil, "unable to save token into session storage")
		return
	}

	// Return UUID
	resp, _ := json.Marshal(decoder.Authorize{
		Expiration: time.Now().Add(dur),
		Token:      uid,
	})

	w.Write([]byte(resp))
}

// HandleRegister route
func HandleRegister(w http.ResponseWriter, r *http.Request) {
	// Decode user
	u, errDecode := decodeUser(r)
	if errDecode != nil {
		w.WriteHeader(http.StatusBadRequest)
		response.CreateResponse(w, response.ResponseError, nil, "")
	}

	// Create new user
	nu, errCreate := accountdb.NewUser(u.User, u.Password, u.Email, accountdb.DefaultUser)
	if errCreate != nil {
		w.WriteHeader(http.StatusBadRequest)
		response.CreateResponse(w, response.ResponseError, nil, errCreate)
		return
	}

	// Create new user
	cu, errSave := UserDB.SaveUser(u)
	if errCreate != nil {
		w.WriteHeader(http.StatusBadRequest)
		response.CreateResponse(w, response.ResponseError, nil, errSave)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// decodeUser to the user.User form
func decodeUser(r *http.Request) (*decoder.User, error) {
	var u decoder.User

	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		return nil, err
	}

	return &u, nil
}
