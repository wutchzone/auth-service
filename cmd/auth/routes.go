package main

import (
	"io/ioutil"
	"net/http"
	"time"

	"github.com/wutchzone/auth-service/pkg/user"

	"github.com/wutchzone/auth-service/pkg/session"

	"github.com/wutchzone/auth-service/pkg/userdb"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

var client = http.Client{
	Timeout: 10 * time.Second,
}

// InitRoutes for api
func InitRoutes() *chi.Mux {
	r := chi.NewRouter()

	r.Route("/api", func(r chi.Router) {
		// Make authorozation for correct privilegies
		r.Use(Authorize)

		// Attach routes from config files
		for _, e := range Config.API {
			r.Get(e.From, func(w http.ResponseWriter, r *http.Request) {
				Redirect(w, r, e.To)
			})
			r.Post(e.From, func(w http.ResponseWriter, r *http.Request) {
				Redirect(w, r, e.To)
			})
			r.Delete(e.From, func(w http.ResponseWriter, r *http.Request) {
				Redirect(w, r, e.To)
			})
			r.Put(e.From, func(w http.ResponseWriter, r *http.Request) {
				Redirect(w, r, e.To)
			})
		}

		r.Route("/user", func(r chi.Router) {
			// Make authorozation for correct privilegies
			r.Use(CheckIfAdminOrUser)

			r.Get("/{name}", HandleUserGet) // Display user
			r.Put("/{name}", HandleUserPut) // Update user settings
		})
	})

	r.Route("/auth", func(r chi.Router) {
		r.Use(middleware.AllowContentType("application/json"))
		r.Post("/register", HandleRegister)
		r.Post("/login", HandleLogin)
	})

	return r
}

// CheckIfAdminOrUser checks for correct ptivilegies
func CheckIfAdminOrUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := chi.URLParam(r, "name")
		uid := r.Header.Get("X-UUID")

		// Check if UUID is in session DB
		sn, errSession := session.GetRecord(uid)
		if errSession != nil {
			sendError(w, "Invalid UUID", http.StatusBadRequest)
			return
		}

		// Check if user is current user
		if sn != n {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Check if user is admin
		u, errUserdb := userdb.GetUser(n)
		if errUserdb != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if u.Role < user.Admin {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Authorize is a middleware for authorization
func Authorize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check authorization

		uid := r.Header.Get("X-UUID")
		id := r.Header.Get("X-USERNAME")

		st, _ := session.GetRecord(id)

		if st != uid {
			sendError(w, "Invalid token or you must log in first", http.StatusUnauthorized)
			return
		}
		//ctx := context.WithValue(r.Context(), "UUID", uuid)

		next.ServeHTTP(w, r) //r.WithContext(ctx))
	})
}

// Redirect user to specific service
func Redirect(w http.ResponseWriter, r *http.Request, location string) {
	nr, _ := http.NewRequest(r.Method, location, r.Body)
	nr.Header = r.Header

	mr, err := client.Do(nr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	body, _ := ioutil.ReadAll(mr.Body)
	// TODO implement Header to response
	w.Write(body)
}
