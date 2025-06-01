package services

import (
	"encoding/json"
	"io"
	"log"
	"media_transcoder/models"
	"mime/multipart"
	"os"
	"os/exec"
	"path/filepath"
)

func UnMarshelFormData(jsonStr string) (models.Format, error) {
	var format models.Format
	err := json.Unmarshal([]byte(jsonStr), &format)
	log.Println("Format: \nFile Type:", format.FileType, "\nRequired File Type:", format.RequiredFileType)
	return format, err
}

// TODO: Add check logic to see if the outFile already exists or not
// ELSE: Add logic to fetch the already existing outFile and change name so as not to overwrite
func FileFormatConversion(dstPath string, outFile string) error {
	cmd := exec.Command("ffmpeg", "-i", dstPath, outFile)
	err := cmd.Start()
	if err != nil {
		return err
	}

	log.Println("Waiting for command to finish...")
	err = cmd.Wait()
	log.Println("Command finished, err:", err)
	return nil
}

func UploadFileToDir(uploadDir string, filename string, file multipart.File) (string, error) {
	os.MkdirAll(uploadDir, os.ModePerm)
	dstPath := filepath.Join(uploadDir, filename)

	dst, err := os.Create(dstPath)
	if err != nil {
		log.Println("error creating file", err)
		return dstPath, err
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		return dstPath, err
	}
	return dstPath, nil
}
