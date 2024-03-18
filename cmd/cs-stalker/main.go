package main

import (
	"fmt"
	"github.com/finnbechinka/cs-stalker/internal/api"
	"github.com/finnbechinka/cs-stalker/internal/routes"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

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

	router := routes.NewRouter()
	port := ":8085"

	s := &http.Server{
		Addr:    port,
		Handler: router,
	}

	log.Println(fmt.Sprintf("server listening on http://localhost%s", port))

	id, _ := api.ResolveVanityUrl("cantremovethis")
	log.Printf(id)
	id, _ = api.ResolveVanityUrl("https://steamcommunity.com/id/cantremovethis/")
	log.Printf(id)
	id, _ = api.ResolveVanityUrl("76561198056395137")
	log.Printf(id)

	log.Fatal(s.ListenAndServe())
}
