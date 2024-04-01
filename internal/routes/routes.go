package routes

import (
	"html/template"
	"log"
	"net/http"

	"github.com/finnbechinka/cs-stalker/internal/api"
)

func NewRouter() http.Handler {
	mux := http.NewServeMux()

	// resolve requests to staticfiles (css, imgs, ...)
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	mux.HandleFunc("GET /", rootHandler)

	mux.HandleFunc("GET /minion", minionHandler)

	mux.HandleFunc("GET /api", apiRootHandler)

	mux.HandleFunc("GET /api/resolveurl", apiResolveUrlHandler)

	mux.HandleFunc("POST /profile", profilePostHandler)
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

func profilePostHandler(w http.ResponseWriter, r *http.Request) {
	url := r.PostFormValue("url")
	steam64id, _ := api.ResolveUrl(url)
	summary, _ := api.UserSummary(steam64id)

	tmpl := template.Must(template.ParseFiles("./templates/fragments/profile.html"))
	err := tmpl.Execute(w, summary)
	if err != nil {
		log.Printf("profilePostHandler: error executing template; err: %s", err)
		http.Error(w, "D'oh, something went wrong!", http.StatusInternalServerError)
		return
	}
}
