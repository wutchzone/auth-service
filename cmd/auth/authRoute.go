package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/wutchzone/auth-service/pkg/sessiondb"

	"github.com/google/uuid"
	"github.com/wutchzone/auth-service/api"
	"github.com/wutchzone/auth-service/pkg/user"
	"github.com/wutchzone/auth-service/pkg/accountdb"
)

// HandleLogin route
func HandleLogin(w http.ResponseWriter, r *http.Request) {
	// Decode user
	u, errDecode := decodeUser(r)
	if errDecode != nil {
		sendError(w, errDecode.Error(), http.StatusUnprocessableEntity)
		return
	}

	// Check if user exists in userDB
	cu, errSearch := accountdb.GetUser((*u).Name)
	if errSearch != nil || cu.ComparePswdAndHash((*u).Password) != nil {
		sendError(w, "User name or password is invalid", http.StatusUnauthorized)
		return
	}

	// Save user to sessionDB
	uid, errSession := sessiondb.GetRecord(u.Name)
	if errSession != nil {
		uuid, _ := uuid.NewUUID()
		uid = uuid.String()
		sessiondb.SetRecord(uid, cu.Username, 20*time.Minute)
	}
	// Return UUID
	json, _ := json.Marshal(parser.UUIDJSON{
		UUID: uid,
	})
	w.Write([]byte(json))
}

// HandleRegister route
func HandleRegister(w http.ResponseWriter, r *http.Request) {
	// Decode user
	u, errDecode := decodeUser(r)
	if errDecode != nil {
		sendError(w, errDecode.Error(), http.StatusUnprocessableEntity)
		return
	}

	// Create new user
	cu, errCreate := createNewUser(u)
	if errCreate != nil {
		sendError(w, errCreate.Error(), http.StatusBadRequest)
		return
	}

	// Handler saving user to DB
	errSave := accountdb.SaveUser(*cu)
	if errSave != nil {
		sendError(w, errCreate.Error(), http.StatusInternalServerError)
		return
	}
}

func sendError(w http.ResponseWriter, message string, hs int) {
	w.WriteHeader(hs)
	json, _ := json.Marshal(parser.ErrorJSON{
		Error: message,
	})
	w.Write([]byte(json))
}

// decodeUser to the user.User form
func decodeUser(r *http.Request) (*parser.UserJSON, error) {
	var u parser.UserJSON
	errDecode := json.NewDecoder(r.Body).Decode(&u)
	if errDecode != nil {
		return nil, errDecode
	}
	return &u, nil
}

// createNewUser from parsed JSON
func createNewUser(u *parser.UserJSON) (*user.User, error) {
	errUsrEx := checkIfUserExists(u)
	if errUsrEx != nil {
		return nil, errUsrEx
	}

	cu, errCreate := user.NewUser(u.Name, u.Password, u.Email, user.DefaultUser)
	if errCreate != nil {
		return nil, errCreate
	}

	return cu, nil
}

// checkIfUserExists in User DB
func checkIfUserExists(usr *parser.UserJSON) error {
	u, _ := accountdb.GetUser(usr.Name)

	if u != nil {
		return errors.New("User already exists")
	}
	return nil
}
