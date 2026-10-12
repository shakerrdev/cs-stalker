package routes

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/shakerrdev/cs-stalker/internal/api"
)

type site struct {
	templates *template.Template
}

func NewRouter(root string) (http.Handler, error) {
	templates, err := template.ParseFiles(
		filepath.Join(root, "templates", "index.html"),
		filepath.Join(root, "templates", "fragments", "minion.html"),
		filepath.Join(root, "templates", "fragments", "profile.html"),
	)
	if err != nil {
		return nil, fmt.Errorf("parse templates: %w", err)
	}
	s := &site{templates: templates}

	mux := http.NewServeMux()
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir(filepath.Join(root, "static")))))
	mux.HandleFunc("GET /{$}", s.index)
	mux.HandleFunc("GET /minion", s.minion)
	mux.HandleFunc("GET /api", apiRootHandler)
	mux.HandleFunc("GET /api/resolveurl", apiResolveUrlHandler)
	mux.HandleFunc("POST /profile", s.profile)
	return mux, nil
}

func (s *site) render(w http.ResponseWriter, name string, data any) {
	var buf bytes.Buffer
	if err := s.templates.ExecuteTemplate(&buf, name, data); err != nil {
		log.Printf("render %s: %v", name, err)
		http.Error(w, "D'oh, something went wrong!", http.StatusInternalServerError)
		return
	}
	buf.WriteTo(w)
}

func (s *site) index(w http.ResponseWriter, r *http.Request) {
	s.render(w, "index.html", nil)
}

func (s *site) minion(w http.ResponseWriter, r *http.Request) {
	s.render(w, "minion.html", nil)
}

func (s *site) profile(w http.ResponseWriter, r *http.Request) {
	steam64id, _ := api.ResolveUrl(r.PostFormValue("url"))
	summary, _ := api.UserSummary(steam64id)
	s.render(w, "profile.html", summary)
}
