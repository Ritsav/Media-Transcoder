package tasks

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// Handles Uploading of Videos
func UploadVideoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only Post allowed at this endpoint", http.StatusMethodNotAllowed)
		return
	}

	// Parsing video File upto 100 MB in memory and remaining memory is stored in disk
	r.ParseMultipartForm(100 << 20)           // << does leftshift 20 times
	file, handler, err := r.FormFile("video") // frontend should send formField name as video while sending vid
	if err != nil {
		log.Println("error retrieving file", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Upload directory
	uploadDir := "./uploads/videos"
	os.MkdirAll(uploadDir, os.ModePerm)
	dstPath := filepath.Join(uploadDir, handler.Filename)

	dst, err := os.Create(dstPath)
	if err != nil {
		log.Println("error creating file", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Uploaded File with filename: %s", handler.Filename)
}

func UploadAudioHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only post method allowed at this endpoint", http.StatusMethodNotAllowed)
		return
	}

	// Parsing audio File upto 100 MB in memory and remaining memory is stored in disk
	r.ParseMultipartForm(100 << 20)
	file, handler, err := r.FormFile("audio")
	if err != nil {
		log.Println("error retrieving audio file", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Upload directory
	uploadDir := "./uploads/audios"
	os.MkdirAll(uploadDir, os.ModePerm)
	dstPath := filepath.Join(uploadDir, handler.Filename)

	dst, err := os.Create(dstPath)
	if err != nil {
		log.Println("error creating audio file", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Uploaded File with filename: %s", handler.Filename)
}
