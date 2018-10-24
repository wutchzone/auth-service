package main

import (
	"encoding/json"
	"net/http"

	"github.com/wutchzone/auth-service/pkg/user"

	"github.com/go-chi/chi"
	"github.com/wutchzone/auth-service/pkg/userdb"
)

// HandleUserGet
func HandleUserGet(w http.ResponseWriter, r *http.Request) {
	u, err := userdb.GetUser(chi.URLParam(r, "name"))
	if err != nil {
		sendError(w, "User does not exist", http.StatusBadRequest)
		return
	}

	u.Password = ""
	su, _ := json.Marshal(u)
	w.Write(su)
}

// HandleUserPut can edit users preferences like email, password, role
func HandleUserPut(w http.ResponseWriter, r *http.Request) {
	du, _ := decodeUser(r)

	u, err := userdb.GetUser(chi.URLParam(r, "name"))
	if err != nil {
		sendError(w, "User does not exist", http.StatusBadRequest)
		return
	}

	if du.Password != "" {
		u.Password = du.Password
	}
	if du.Email != "" {
		u.Email = du.Email
	}

	if du.Role == user.SuperAdmin {
		u.Role = du.Role
	}

	userdb.SaveUser(*u)
}
