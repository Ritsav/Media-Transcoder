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
	AUDIO_ACC
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
	*cmd = CommandChoose(dstPath, outFile, conversionType)
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
// 1. Bitrate wise: 96k(OPUS: Great for low bitrate), 128k, 192k, 320k
// 2. Variable quality from 0-10 (5 being ~128-160kbps)
func CommandChoose(dstFile string, outFile string, conversionType ConversionType) exec.Cmd {
	switch conversionType {
	case VIDEO_WAV:
		return *exec.Command("ffmpeg", "-y", "-i", dstFile, "-vn", "-c:a", "pcm_s16le", outFile)

	case VIDEO_MP3:
		return *exec.Command("ffmpeg", "-y", "-i", dstFile, "-vn", "-c:a", "libmp3lame", "-b:a", "192k", outFile)

	case VIDEO_AAC:
		return *exec.Command("ffmpeg", "-y", "-i", dstFile, "-vn", "-c:a", "aac", "-b:a", "192k", outFile)

	case VIDEO_OGG:
		return *exec.Command("ffmpeg", "-y", "-i", dstFile, "-vn", "-c:a", "libvorbis", "-q:a", "5", outFile)

	case VIDEO_OPUS:
		return *exec.Command("ffmpeg", "-y", "-i", dstFile, "-vn", "-c:a", "libopus", "-b:a", "96k", outFile)

	case VIDEO_FLAC:
		return *exec.Command("ffmpeg", "-y", "-i", dstFile, "-vn", "-c:a", "flac", outFile)

	case VIDEO_ALAC:
		return *exec.Command("ffmpeg", "-y", "-i", dstFile, "-vn", "-c:a", "alac", outFile)

	case AUDIO_WAV:
		return *exec.Command("ffmpeg", "-y", "-i", dstFile, "-c:a", "pcm_s16le", outFile)

	case AUDIO_MP3:
		return *exec.Command("ffmpeg", "-y", "-i", dstFile, "-c:a", "libmp3lame", "-b:a", "192k", outFile)

	case AUDIO_ACC:
		return *exec.Command("ffmpeg", "-y", "-i", dstFile, "-c:a", "aac", "-b:a", "192k", outFile)

	case AUDIO_OGG:
		return *exec.Command("ffmpeg", "-y", "-i", dstFile, "-c:a", "libvorbis", "-q:a", "5", outFile)

	case AUDIO_OPUS:
		return *exec.Command("ffmpeg", "-y", "-i", dstFile, "-c:a", "libopus", "-b:a", "96k", outFile)

	case AUDIO_FLAC:
		return *exec.Command("ffmpeg", "-y", "-i", dstFile, "-c:a", "flac", outFile)

	case AUDIO_ALAC:
		return *exec.Command("ffmpeg", "-y", "-i", dstFile, "-c:a", "alac", outFile)

	default:
		return *exec.Command("ffmpeg", "-i", dstFile, outFile, "-y")
	}
}
