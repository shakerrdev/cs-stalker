package routes

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/shakerrdev/cs-stalker/internal/api"
)

type site struct {
	pages map[string]*template.Template
}

func NewRouter(root string) (http.Handler, error) {
	dir := filepath.Join(root, "templates")
	minion := filepath.Join(dir, "fragments", "minion.html")

	s := &site{pages: map[string]*template.Template{}}
	for _, name := range []string{"home", "player"} {
		page, err := template.New(name).Funcs(funcs).ParseFiles(filepath.Join(dir, "layout.html"), filepath.Join(dir, name+".html"), minion)
		if err != nil {
			return nil, fmt.Errorf("parse %s template: %w", name, err)
		}
		s.pages[name] = page
	}

	mux := http.NewServeMux()
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir(filepath.Join(root, "static")))))
	mux.HandleFunc("GET /{$}", s.home)
	mux.HandleFunc("GET /lookup", s.lookup)
	mux.HandleFunc("GET /player/{id}", s.player)
	mux.HandleFunc("GET /banana", s.banana)
	mux.HandleFunc("GET /healthz", healthz)
	mux.HandleFunc("GET /api", apiRootHandler)
	mux.HandleFunc("GET /api/resolveurl", apiResolveUrlHandler)
	return mux, nil
}

func healthz(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
}

func (s *site) render(w http.ResponseWriter, status int, page string, data any) {
	var buf bytes.Buffer
	if err := s.pages[page].ExecuteTemplate(&buf, "layout", data); err != nil {
		log.Printf("render %s: %v", page, err)
		http.Error(w, "D'oh, something went wrong!", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	buf.WriteTo(w)
}

type homePage struct {
	Query  string
	Error  string
	Banana bool
}

func (s *site) home(w http.ResponseWriter, r *http.Request) {
	s.render(w, http.StatusOK, "home", homePage{})
}

func (s *site) lookup(w http.ResponseWriter, r *http.Request) {
	query := strings.TrimSpace(r.URL.Query().Get("q"))
	switch {
	case query == "":
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	case strings.EqualFold(query, "banana"):
		http.Redirect(w, r, "/banana", http.StatusSeeOther)
		return
	}

	steam64id, err := api.ResolveUrl(query)
	switch {
	case errors.Is(err, api.ErrNotFound):
		s.render(w, http.StatusNotFound, "home", homePage{Query: query, Error: "No Steam profile matches that. Paste the full profile link or the 17-digit Steam64 ID."})
	case err != nil:
		log.Printf("lookup: %v", err)
		s.render(w, http.StatusBadGateway, "home", homePage{Query: query, Error: "Steam didn't answer. Try again in a moment."})
	default:
		http.Redirect(w, r, "/player/"+steam64id, http.StatusSeeOther)
	}
}

func (s *site) player(w http.ResponseWriter, r *http.Request) {
	steam64id := r.PathValue("id")
	if !api.IsSteam64(steam64id) {
		s.render(w, http.StatusNotFound, "home", homePage{Query: steam64id, Error: "That isn't a Steam64 ID. Paste the profile link instead."})
		return
	}

	summary, err := api.UserSummary(steam64id)
	switch {
	case errors.Is(err, api.ErrNotFound):
		s.render(w, http.StatusNotFound, "home", homePage{Query: steam64id, Error: "No Steam account has that ID."})
		return
	case err != nil:
		log.Printf("player: %v", err)
		s.render(w, http.StatusBadGateway, "home", homePage{Query: steam64id, Error: "Steam didn't answer. Try again in a moment."})
		return
	}

	s.render(w, http.StatusOK, "player", gatherPlayer(summary))
}

func (s *site) banana(w http.ResponseWriter, r *http.Request) {
	s.render(w, http.StatusOK, "home", homePage{Banana: true})
}
