package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"media_transcoder/dto"
	"media_transcoder/pkg/global"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

// Single File Processing
// Handles Uploading of Files(Audio/Video)
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only Post allowed at this endpoint", http.StatusMethodNotAllowed)
		return
	}

	// Unmarsheling format from formData
	jsonStr := r.FormValue("format")
	format, err := unMarshelFormData(jsonStr)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Parsing video File upto 100 MB in memory and remaining memory is stored in disk
	if err := r.ParseMultipartForm(100 << 20); err != nil { // << does leftshift 20 times
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}
	file, handler, err := r.FormFile("media") // frontend should send formField name as media while sending file
	if err != nil {
		log.Println("error retrieving file", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Upload directory
	uploadDir := "./uploads/media"
	err = uploadFileToDir(uploadDir, handler.Filename, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Uploaded File with filename: %s\n", handler.Filename)

	// Upload handler queues the file conversion process
	go global.TaskQueue.Enqueue(handler.Filename, format)

	// TODO: Refactor this file conversion functionality elsewhere
	// outFile := fmt.Sprintf("%s/output.%s", uploadDir, format.RequiredFileType)
	// err = services.FileFormatConversion(dstPath, outFile, format)
	// if err != nil {
	// 	http.Error(w, "error running command", http.StatusInternalServerError)
	// 	log.Println(err)
	// 	return
	// }
	fmt.Fprintln(w, "File upload complete")
}

// TODO: Look into unmarshaling enum type (video/audio ONLY)
func unMarshelFormData(jsonStr string) (dto.Format, error) {
	var format dto.Format
	err := json.Unmarshal([]byte(jsonStr), &format)
	log.Println("Format: \nFile Type:", format.MediaType, "\nRequired File Type:", format.RequiredFileType)
	return format, err
}

func uploadFileToDir(uploadDir string, filename string, file multipart.File) error {
	os.MkdirAll(uploadDir, os.ModePerm)
	dstPath := filepath.Join(uploadDir, filename)

	dst, err := os.Create(dstPath)
	if err != nil {
		log.Println("error creating file", err)
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		return err
	}
	return nil
}
