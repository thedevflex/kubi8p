package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	log.Printf("Starting server on :80")
	if err := http.ListenAndServe(":80", nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
