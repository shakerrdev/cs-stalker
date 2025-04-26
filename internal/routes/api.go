package routes

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/finnbechinka/cs-stalker/internal/api"
)

func apiRootHandler(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	resp := make(map[string]string)
	resp["hello"] = "world"

	json, err := json.Marshal(resp)
	if err != nil {
		log.Printf("apiRootHandler: error marshaling json; err %s", err)
		http.Error(w, "D'oh, something went wrong!", http.StatusInternalServerError)
		return
	}

	w.Write(json)
}

func apiResolveUrlHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	url := r.FormValue("url")
	if url == "" {
		log.Println("apiResolveUrlHandler: no url query param, FormValue returned empty string")
		http.Error(w, "D'oh, something went wrong!", http.StatusInternalServerError)
		return

	}

	steam64id, err := api.ResolveUrl(url)
	if err != nil {
		log.Println(fmt.Errorf("apiResolveUrlHandler: %w", err))
		http.Error(w, "D'oh, something went wrong!", http.StatusInternalServerError)
		return
	}

	resp := make(map[string]string)
	resp["steam64id"] = steam64id

	json, err := json.Marshal(resp)
	if err != nil {
		log.Println(fmt.Errorf("apiResolveUrlHandler: error marshaling json; err %w", err))
		http.Error(w, "D'oh, something went wrong!", http.StatusInternalServerError)
		return
	}

	w.Write(json)
}
