package server

import (
	"log"
	"net/http"

	"github.com/thedevflex/kubi8p/internal/handler"
)

type Server struct {
	mux *http.ServeMux
	api *http.ServeMux
}

func NewServer() *Server {
	mux := http.NewServeMux()
	api := http.NewServeMux()
	mux.Handle("/api/", http.StripPrefix("/api", api))
	return &Server{mux: mux, api: api}
}

func (s *Server) Register(handler *handler.Handler) {

	// API routes
	s.addApiRoutes(handler)
	// Static file serving
	s.addStaticRoutes()
}

func (s *Server) Start(port string) {
	log.Printf("Starting server on http://localhost:%s", port)
	if err := http.ListenAndServe(":"+port, s.mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
