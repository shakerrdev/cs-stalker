package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/finnbechinka/cs-stalker/internal/api"
	"github.com/finnbechinka/cs-stalker/internal/routes"
	"github.com/joho/godotenv"
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

	err := godotenv.Load()
	if err != nil {
		log.Panicf("main: error loading .env; err: %s", err)
	}

	evar, exists := os.LookupEnv("STEAMAPIKEY")
	if exists {
		log.Printf("steam api key: %s...", evar[:5])
	} else {
		log.Panicf("no steam api key env var set")
	}

	router := loggingMiddleware(routes.NewRouter())
	port := ":8085"

	s := &http.Server{
		Addr:    port,
		Handler: router,
	}

	log.Println(fmt.Sprintf("server listening on http://localhost%s", port))

	id, _ := api.UserSummary("76561198056395137")
	log.Printf("%+v", id)

	time, _ := api.UserPlaytime("76561198056395137")
	log.Printf("%d", time)

	profile, _ := api.LeetifyProfile("76561198056395137")
	log.Printf("%+v", profile)

	_, err = api.LeetifyProfile("76561198801755202")
	if err != nil {
		log.Println(err)
	}
	log.Fatal(s.ListenAndServe())
}
