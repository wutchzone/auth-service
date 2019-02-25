package main

import (
	"fmt"
	"github.com/go-chi/chi"
	"io/ioutil"
	"net/http"
	"time"
)

var client = http.Client{
	Timeout: 10 * time.Second,
}

// InitRoutes for api
func InitRoutes() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/", handleHello)

	r.Route("/api", func(r chi.Router) {
		// Check if user has enough privilegies
		r.Use(Authorize)

		// Redirect
		r.Use(Redirect)
		r.Get("/", handleHello)
	})

	r.Route("/user", func(r chi.Router) {
		// Make authorization for correct privilegies
		// Because normal user can see only his profile
		// r.Use(CheckIfAdminOrUser)

		// r.Get("/{name}", nil) // Display user
		// r.Put("/{name}", nil) // Update user settings
	})

	r.Route("/service", func(r chi.Router) {
		// Check if user has enough privilegies
		r.Use(Authorize)

		// Register new route
		r.Get("/", handleGetAllServices)
		r.Post("/", handleRegisterRoute)
		r.Get("/{id}", handleGetOneService)
		r.Delete("/{id}", handleDeleteOneService)
	})

	r.Route("/auth", func(r chi.Router) {
		//r.Post("/register", nil)
		//r.Post("/login", nil)
		//r.Post("/logout", nil)
		//r.Post("/reset", nil)

		// Send available roles
		// r.Get("/roles", nil)
	})

	return r
}

func handleHello(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte(fmt.Sprintf("Server is listening on port: %v", Config.DefaultPort)))
}

// CheckIfAdminOrUser checks for correct ptivilegies
//func CheckIfAdminOrUser(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		url := strings.Split(r.URL.String(), "/")
//		n := url[len(url)-1]
//		uid := r.Header.Get("X-UUID")
//
//		// Check if UUID is in sessiondb DB
//		sn, errSession := sessiondb.GetRecord(uid)
//		if errSession != nil {
//			sendError(w, "Invalid UUID", http.StatusBadRequest)
//			return
//		}
//
//		// Check if user is current user
//		if sn != n {
//			w.WriteHeader(http.StatusUnauthorized)
//			return
//		}
//
//		// Check if user is admin
//		u, errUserdb := accountdb.GetUser(n)
//		if errUserdb != nil {
//			w.WriteHeader(http.StatusInternalServerError)
//			return
//		}
//		if sn == n {
//			next.ServeHTTP(w, r)
//		}
//		if u.Role < user.Admin {
//			w.WriteHeader(http.StatusUnauthorized)
//			return
//		}
//
//		next.ServeHTTP(w, r)
//	})
//}
//
// Authorize is a middleware for authorization
func Authorize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check authorization
		uid := r.Header.Get("X-UUID")

		_, err := SessionDB.GetRecord(uid)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

// Redirect user to specific service
func Redirect(_ http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		originalReq := r.URL.Path
		to := ""
		didFind := false

		for _, i := range Config.Routes {
			if originalReq == i.From {
				to = i.To
				didFind = true
				break
			}
		}

		if didFind == false {
			w.WriteHeader(http.StatusNotFound)
		} else {
			nr, _ := http.NewRequest(r.Method, to, r.Body)
			nr.Header = r.Header

			mr, err := client.Do(nr)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			} else {
				body, _ := ioutil.ReadAll(mr.Body)
				for k, v := range mr.Header {
					w.Header().Set(k, v[0])
				}

				w.Write(body)
			}
		}
	})
}
