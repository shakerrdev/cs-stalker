package routes

import (
	"encoding/json"
	"log"
	"net/http"
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

