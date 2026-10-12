package main

import (
	"errors"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/shakerrdev/cs-stalker/internal/db"
	"github.com/shakerrdev/cs-stalker/internal/routes"
)

func trimTrailingSlash(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
		}
		next.ServeHTTP(w, r)
	})
}

func logRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func main() {
	log.Println("cs-stalker started")

	if err := godotenv.Load(); err != nil && !errors.Is(err, fs.ErrNotExist) {
		log.Fatalf("load .env: %v", err)
	}

	if os.Getenv("STEAMAPIKEY") == "" {
		log.Fatal("missing environment variable: STEAMAPIKEY")
	}

	dsn, err := db.DSNFromEnv()
	switch {
	case errors.Is(err, db.ErrNotConfigured):
		log.Println("no database configured, running without one")
	case err != nil:
		log.Fatal(err)
	default:
		if err := db.Connect(dsn); err != nil {
			log.Fatal(err)
		}
		log.Println("connected to database")
	}

	router, err := routes.NewRouter(".")
	if err != nil {
		log.Fatal(err)
	}
	if password := os.Getenv("SITE_PASSWORD"); password != "" {
		router = routes.RequirePassword(password, router)
	} else {
		log.Println("SITE_PASSWORD not set, the site is open to everyone")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8085"
	}

	s := &http.Server{
		Addr:    ":" + port,
		Handler: logRequests(trimTrailingSlash(router)),
	}

	log.Printf("server listening on http://localhost:%s", port)
	log.Fatal(s.ListenAndServe())
}
