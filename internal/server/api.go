package server

import (
	"net/http"

	"github.com/thedevflex/kubi8p/internal/handler"
)

func (s *Server) addApiRoutes(handler *handler.Handler) {

	// Health check endpoint
	s.mux.HandleFunc("/api/_health_", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	// Add your API sub-routes here
	s.api.HandleFunc("/initiate-db", handler.InitiateDB)
	s.api.HandleFunc("/initiate-webhook", handler.InitiateWebhook)
	s.api.HandleFunc("/initiate-dns", handler.InitiateDNS)
	s.api.HandleFunc("/get-dns-status", handler.GetDNSStatus)
	s.api.HandleFunc("/configure-dns", handler.ConfigureDNS)
	s.api.HandleFunc("/verify-dns", handler.VerifyDNS)
	// for rest return not found
	s.api.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not Found"))
	})

}
