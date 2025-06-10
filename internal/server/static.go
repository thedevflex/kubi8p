package server

import (
	"net/http"
	"os"
	"path/filepath"
)

func (s *Server) addStaticRoutes() {

	s.mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Get the requested path
		path := r.URL.Path
		if path == "/" {
			path = "/index.html"
		}

		// Construct the file path
		filePath := filepath.Join("public", path)

		// Check if file exists
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			// If file doesn't exist, serve 404 page
			notFoundPath := filepath.Join("public", "404.html")
			if _, err := os.Stat(notFoundPath); os.IsNotExist(err) {
				http.NotFound(w, r)
				return
			}
			http.ServeFile(w, r, notFoundPath)
			return
		}

		// Set content type based on file extension
		ext := filepath.Ext(filePath)
		switch ext {
		case ".html":
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
		case ".css":
			w.Header().Set("Content-Type", "text/css")
		case ".js":
			w.Header().Set("Content-Type", "application/javascript")
		}

		// Set cache control headers to prevent caching
		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")

		// Serve the file
		http.ServeFile(w, r, filePath)
	})
}
