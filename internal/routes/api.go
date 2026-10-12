package routes

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/shakerrdev/cs-stalker/internal/api"
)

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		log.Printf("write json response: %v", err)
	}
}

func apiRootHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"hello": "world"})
}

func apiResolveUrlHandler(w http.ResponseWriter, r *http.Request) {
	input := r.FormValue("url")
	if input == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing url parameter"})
		return
	}

	steam64id, err := api.ResolveUrl(input)
	if err != nil {
		log.Printf("apiResolveUrlHandler: %v", err)
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "could not resolve a steam profile"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"steam64id": steam64id})
}
