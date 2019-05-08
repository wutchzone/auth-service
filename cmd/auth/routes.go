package main

import (
	"fmt"
	"github.com/go-chi/chi"
	"log"
	"net/http"
)

// InitRoutes for api
func InitRoutes() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/", handleHello)

	r.Route("/", func(r chi.Router) {
		// Log all access to STDIN
		r.Use(Log)

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
			//r.Get("/", handleGetAllServices)
			//r.Post("/", handleRegisterRoute)
			//r.Get("/{id}", handleGetOneService)
			//r.Delete("/{id}", handleDeleteOneService)
		})

		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", HandleRegister)
			r.Post("/login", HandleLogin)
			//r.Post("/logout", nil) // logout
			//r.Post("/reset", nil) // Reset password
			//r.Get //validate

			// Send available roles
			// r.Get("/roles", nil)
		})
	})

	return r
}

func handleHello(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte(fmt.Sprintf("Server is listening on port: %v\n", Config.DefaultPort)))
}

// Authorize is a middleware for authorization
func Authorize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check authorization
		uid := r.Header.Get("Authorization")

		_, err := SessionDB.GetRecord(uid)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

// Log input to STDIN
func Log(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%v", r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
