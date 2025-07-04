package main

import (
	"log"
	"media_transcoder/handlers"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	port := "localhost:8080"
	router := mux.NewRouter()

	// Endpoints for queuing a new task to the backend
	router.HandleFunc("/upload", handlers.UploadHandler).Methods("POST")
	router.HandleFunc("/status", handlers.StatusHandler).Methods("GET")

	log.Fatal(http.ListenAndServe(port, router))
}
