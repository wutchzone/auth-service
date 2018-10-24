package main

import (
	"fmt"
	"net/http"

	"github.com/wutchzone/auth-service/pkg/user"

	"github.com/wutchzone/auth-service/pkg/session"

	"github.com/wutchzone/auth-service/pkg/userdb"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

// InitRoutes for api
func InitRoutes() *chi.Mux {
	r := chi.NewRouter()

	r.Route("/api", func(r chi.Router) {
		// Make authorozation for correct privilegies
		//r.Use(Authorize)

		// Attach routes from config files
		// for _, e := range Config.API {
		// 	r.Get(e.From, func(w http.ResponseWriter, r *http.Request) {
		// 		Redirect(w, r)
		// 	})
		// 	r.Post(e.From, func(w http.ResponseWriter, r *http.Request) {
		// 		Redirect(w, r)
		// 	})
		// 	r.Delete(e.From, func(w http.ResponseWriter, r *http.Request) {
		// 		Redirect(w, r)
		// 	})
		// 	r.Put(e.From, func(w http.ResponseWriter, r *http.Request) {
		// 		Redirect(w, r)
		// 	})
		// }

		r.Route("/user", func(r chi.Router) {
			// Make authorozation for correct privilegies
			r.Use(CheckIfAdminOrUser)

			r.Get("/{name}", HandleUserGet) // Display user
			r.Put("/{name}", nil)           // Update user settings
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
		fmt.Println(r.Header)

		uid := r.Header.Get("X-UUID")
		if uid != "test" {
			sendError(w, "Invalid token or you must log in first", http.StatusUnauthorized)
			return
		}
		//ctx := context.WithValue(r.Context(), "UUID", uuid)

		next.ServeHTTP(w, r) //r.WithContext(ctx))
	})
}

func Redirect(w http.ResponseWriter, r *http.Request) {

}
