package routes

import (
	"html/template"
	"log"
	"net/http"
)

func NewRouter() http.Handler {
	mux := http.NewServeMux()

	// resolve requests to staticfiles (css, imgs, ...)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	mux.HandleFunc("/", rootHandler)

	mux.HandleFunc("/minion", minionHandler)

	mux.HandleFunc("/api", apiRootHandler)

	return mux
}

func minionHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("./templates/fragments/minion.html"))
	err := tmpl.Execute(w, nil)
	if err != nil {
		log.Printf("minionHandler: error executing templatel; err: %s", err)
		http.Error(w, "D'oh, something went wrong!", http.StatusInternalServerError)
		return
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("./templates/index.html"))
	err := tmpl.Execute(w, nil)
	if err != nil {
		log.Printf("rootHandler: error executing template; err: %s", err)
		http.Error(w, "D'oh, something went wrong!", http.StatusInternalServerError)
		return
	}
}
