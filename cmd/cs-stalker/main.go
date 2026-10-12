package main

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"

	// "github.com/shakerrdev/cs-stalker/internal/api"
	"github.com/joho/godotenv"
	"github.com/shakerrdev/cs-stalker/internal/db"
	"github.com/shakerrdev/cs-stalker/internal/routes"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			// make sure **/x is reachable from **/x/
			// ref: https://natedenlinger.com/dealing-with-trailing-slashes-on-requesturi-in-go-with-mux/
			r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
		}
		// Do stuff here
		log.Printf("INCOMING REQUEST: %s %s", r.Method, r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
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

	handler, err := routes.NewRouter(".")
	if err != nil {
		log.Fatal(err)
	}
	router := loggingMiddleware(handler)
	port := ":8085"

	s := &http.Server{
		Addr:    port,
		Handler: router,
	}

	log.Println(fmt.Sprintf("server listening on http://localhost%s", port))

	// id, _ := api.UserSummary("76561198056395137")
	// log.Printf("UserSummary: %+v", id)

	// time, _ := api.UserPlaytime("76561198056395137")
	// log.Printf("UserPlaytime: %d", time)

	// profile, _ := api.LeetifyProfile("76561198056395137")
	// log.Printf("LeetifyProfile: %+v", profile)

	// _, err = api.LeetifyProfile("76561198801755202")
	// if err != nil {
	// 	log.Println(err)
	// }
	log.Fatal(s.ListenAndServe())
}
