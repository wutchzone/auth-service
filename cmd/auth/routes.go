package main

import (
	"fmt"
	"github.com/go-chi/chi"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
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
			r.Get("/", handleGetAllServices)
			r.Post("/", handleRegisterRoute)
			r.Get("/{id}", handleGetOneService)
			r.Delete("/{id}", handleDeleteOneService)
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
	w.Write([]byte(fmt.Sprintf("Server is listening on port: %v", Config.DefaultPort)))
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
		log.Printf("%v => %v", r.URL.Path, "TODO")
		next.ServeHTTP(w, r)
	})
}

// PrepareURL for parsing URL
func PrepareURL(w http.ResponseWriter, r *http.Request) {
	reg := regexp.MustCompile(`{.[^}]*}`)
	matches := reg.FindAllString(i.From, -1)

	final := i.To

	for _, i := range matches {
		final = strings.Replace(final, i, chi.URLParam(r, i[1:len(i)-1]), -1)
	}

	fmt.Println(final)
	Redirect(w, r, final)
}

// Redirect user to specific service
func Redirect(w http.ResponseWriter, r *http.Request, to string) {
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
