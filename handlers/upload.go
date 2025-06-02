package handlers

import (
	"fmt"
	"log"
	"media_transcoder/services"
	"net/http"
)

// Single File Processing
// Handles Uploading of Videos
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only Post allowed at this endpoint", http.StatusMethodNotAllowed)
		return
	}

	// Unmarsheling format from formData
	jsonStr := r.FormValue("format")
	format, err := services.UnMarshelFormData(jsonStr)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Parsing video File upto 100 MB in memory and remaining memory is stored in disk
	if err := r.ParseMultipartForm(100 << 20); err != nil { // << does leftshift 20 times
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}
	file, handler, err := r.FormFile("media") // frontend should send formField name as video while sending vid
	if err != nil {
		log.Println("error retrieving file", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Upload directory
	uploadDir := "./uploads/media"
	dstPath, err := services.UploadFileToDir(uploadDir, handler.Filename, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Uploaded File with filename: %s\n", handler.Filename)

	// Running the appropriate command for conversion of file format
	outFile := fmt.Sprintf("%s/output.%s", uploadDir, format.RequiredFileType)
	err = services.FileFormatConversion(dstPath, outFile, format)
	if err != nil {
		http.Error(w, "error running command", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	fmt.Fprintf(w, "Conversion of file complete: %s\n", outFile)
}
