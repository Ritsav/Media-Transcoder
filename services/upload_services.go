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

// ENUM
type ConversionType int

const (
	// Video --> Audio Conversions
	VIDEO_WAV  ConversionType = iota // Lossy, compressed -- .wav
	VIDEO_MP3                        // .mp3
	VIDEO_OGG                        // Open Source Lossy -- .ogg
	VIDEO_OPUS                       // Low Bitrate/Modern lossy -- .opus
	VIDEO_FLAC                       // Lossless compressed -- .flac
	VIDEO_AAC                        // Apple Lossy - Apple Devices -- .m4a
	VIDEO_ALAC                       // Apple Lossless -- .m4a

	// Audio --> Audio Conversions
	AUDIO_WAV
	AUDIO_MP3
	AUDIO_OGG
	AUDIO_OPUS
	AUDIO_FLAC
	AUDIO_AAC
	AUDIO_ALAC
)

func determineConversionType(format models.Format) ConversionType {
	switch {
	case format.MediaType == "video" && format.RequiredFileType == "wav":
		return ConversionType(VIDEO_WAV)

	default:
		return ConversionType(VIDEO_MP3)
	}
}

func UnMarshelFormData(jsonStr string) (models.Format, error) {
	var format models.Format
	err := json.Unmarshal([]byte(jsonStr), &format)
	log.Println("Format: \nFile Type:", format.MediaType, "\nRequired File Type:", format.RequiredFileType)
	return format, err
}

func FileFormatConversion(dstPath string, outFile string, format models.Format) error {
	conversionType := determineConversionType(format)

	// Try Remuxing first
	cmd := exec.Command("ffmpeg", "-y", "-i", dstPath, "-c", "copy", outFile)
	err := cmd.Start()
	if err != nil {
		return err
	}
	log.Println("Waiting for command to finish...")
	err = cmd.Wait()
	if err == nil {
		log.Println("Command finished")
		return nil
	}

	log.Println("Remuxing failed, trying to re-encode:", err)

	// Do Re-encoding if remuxing fails
	// cmd = exec.Command("ffmpeg", "-i", dstPath, "-c:v", "libx264", "-c:a", "aac", outFile, "-y")
	*cmd = ChooseCommand(dstPath, outFile, conversionType, format.MediaType)
	err = cmd.Start()
	if err != nil {
		return err
	}
	log.Println("Waiting for command to finish...")
	err = cmd.Wait()
	if err != nil {
		return err
	}
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

// TODO: AUDIO QUALITY level set:
// Add a ENUM BitrateLevels
// 1. Bitrate wise: 96k(OPUS: Great for low bitrate), 128k, 192k, 320k
// 2. Variable quality from 0-10 (5 being ~128-160kbps)
func ChooseCommand(dstFile string, outFile string, conversionType ConversionType, mediaType string) exec.Cmd {
	var media, codec, setBitrate, bitrate string
	switch mediaType {
	case "video":
		media = "-vn"

	default:
		media = ""
	}

	switch conversionType {
	case VIDEO_WAV, AUDIO_WAV:
		codec = "pcm_s16le"

	case VIDEO_FLAC, AUDIO_FLAC:
		codec = "flac"

	case VIDEO_ALAC, AUDIO_ALAC:
		codec = "alac"

	case VIDEO_OGG, AUDIO_OGG:
		codec = "libvorbis"
		setBitrate = "-q:a"
		bitrate = "5"

	case VIDEO_MP3, AUDIO_MP3:
		codec = "libmp3lame"
		setBitrate = "-b:a"
		bitrate = "192k"

	case VIDEO_AAC, AUDIO_AAC:
		codec = "aac"
		setBitrate = "-b:a"
		bitrate = "192k"

	case VIDEO_OPUS, AUDIO_OPUS:
		codec = "libopus"
		setBitrate = "-b:a"
		bitrate = "96k"

	default:
		return *exec.Command("ffmpeg", "-i", dstFile, outFile, "-y")
	}

	return *exec.Command("ffmpeg", "-y", "-i", dstFile, media, "-c:a", codec, setBitrate, bitrate, outFile)
}
