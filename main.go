package main

import (
	"log"
	"media_transcoder/tasks"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	port := "localhost:8080"
	router := mux.NewRouter()

	// Endpoints for queuing a new task to the backend
	router.HandleFunc("/upload-video", tasks.UploadVideoHandler)
	router.HandleFunc("/upload-audio", tasks.UploadAudioHandler)

	log.Fatal(http.ListenAndServe(port, router))
}
