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

	for _, key := range []string{"STEAMAPIKEY", "LEETIFY_AUTH_TOKEN"} {
		if os.Getenv(key) == "" {
			log.Fatalf("missing environment variable: %s", key)
		}
	}

	dsn, err := db.DSNFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Connect(dsn); err != nil {
		log.Fatal(err)
	}
	log.Println("connected to database")

	router, err := routes.NewRouter(".")
	if err != nil {
		log.Fatal(err)
	}

	s := &http.Server{
		Addr:    ":8085",
		Handler: logRequests(trimTrailingSlash(router)),
	}

	log.Printf("server listening on http://localhost%s", s.Addr)
	log.Fatal(s.ListenAndServe())
}
