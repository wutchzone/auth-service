package main

import (
	"encoding/json"
	"github.com/wutchzone/api-response"
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
		_ := response.CreateResponse(w, response.ResponseError, nil)
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

//
//// HandleRegister route
//func HandleRegister(w http.ResponseWriter, r *http.Request) {
//	// Decode user
//	u, errDecode := decodeUser(r)
//	if errDecode != nil {
//		sendError(w, errDecode.Error(), http.StatusUnprocessableEntity)
//		return
//	}
//
//	// Create new user
//	cu, errCreate := createNewUser(u)
//	if errCreate != nil {
//		sendError(w, errCreate.Error(), http.StatusBadRequest)
//		return
//	}
//
//	// Handler saving user to DB
//	errSave := accountdb.SaveUser(*cu)
//	if errSave != nil {
//		sendError(w, errCreate.Error(), http.StatusInternalServerError)
//		return
//	}
//}
//
//func sendError(w http.ResponseWriter, message string, hs int) {
//	w.WriteHeader(hs)
//	json, _ := json.Marshal(parser.ErrorJSON{
//		Error: message,
//	})
//	w.Write([]byte(json))
//}

// decodeUser to the user.User form
func decodeUser(r *http.Request) (*decoder.User, error) {
	var u decoder.User

	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

//// createNewUser from parsed JSON
//func createNewUser(u *parser.UserJSON) (*user.User, error) {
//	errUsrEx := checkIfUserExists(u)
//	if errUsrEx != nil {
//		return nil, errUsrEx
//	}
//
//	cu, errCreate := user.NewUser(u.Name, u.Password, u.Email, user.DefaultUser)
//	if errCreate != nil {
//		return nil, errCreate
//	}
//
//	return cu, nil
//}
//
//// checkIfUserExists in User DB
//func checkIfUserExists(usr *parser.UserJSON) error {
//	u, _ := accountdb.GetUser(usr.Name)
//
//	if u != nil {
//		return errors.New("User already exists")
//	}
//	return nil
//}
