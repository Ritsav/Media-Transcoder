package services

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"media_transcoder/dto"
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

	// Invalid TYPE
	INVALID_TYPE
)

// TODO: Look into unmarshaling enum type (video/audio ONLY)
func UnMarshelFormData(jsonStr string) (dto.Format, error) {
	var format dto.Format
	err := json.Unmarshal([]byte(jsonStr), &format)
	log.Println("Format: \nFile Type:", format.MediaType, "\nRequired File Type:", format.RequiredFileType)
	return format, err
}

func FileFormatConversion(dstPath string, outFile string, format dto.Format) error {
	// Check Conversion Type
	conversionType := determineConversionType(format)
	if conversionType == INVALID_TYPE {
		return errors.New("invalid conversion type")
	}

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

	// Do Re-encoding if remuxing fails
	log.Println("Remuxing failed, trying to re-encode:", err)
	*cmd = chooseCommand(dstPath, outFile, conversionType, format.MediaType)

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
// 3. Improve function to not let empty string return as command
func chooseCommand(dstFile string, outFile string, conversionType ConversionType, mediaType string) exec.Cmd {
	var args []string
	// Append necessary args
	args = append(args, "-y", "-i", dstFile)

	// Check if we should add media or not
	if mediaType == "video" {
		args = append(args, "-vn")
	}
	args = append(args, "-c:a")

	switch conversionType {
	case VIDEO_WAV, AUDIO_WAV:
		args = append(args, "pcm_s16le")

	case VIDEO_FLAC, AUDIO_FLAC:
		args = append(args, "flac")

	case VIDEO_ALAC, AUDIO_ALAC:
		args = append(args, "alac")

	case VIDEO_OGG, AUDIO_OGG:
		args = append(args, "libvorbis", "-q:a", "5")

	case VIDEO_MP3, AUDIO_MP3:
		args = append(args, "libmp3lame", "-b:a", "192k")

	case VIDEO_AAC, AUDIO_AAC:
		args = append(args, "aac", "-b:a", "192k")

	case VIDEO_OPUS, AUDIO_OPUS:
		args = append(args, "libopus", "-b:a", "96k")

	default:
		args = []string{"-i", "-y", dstFile, outFile}
		return *exec.Command("ffmpeg", args...)
	}

	args = append(args, outFile)
	return *exec.Command("ffmpeg", args...)
}

func determineConversionType(format dto.Format) ConversionType {
	var conversionType ConversionType

	if format.MediaType == "video" {
		switch format.RequiredFileType {
		case "wav":
			conversionType = VIDEO_WAV

		case "mp3":
			conversionType = VIDEO_MP3

		case "ogg":
			conversionType = VIDEO_OGG

		case "opus":
			conversionType = VIDEO_OPUS

		case "flac":
			conversionType = VIDEO_FLAC

		case "aac":
			conversionType = VIDEO_AAC

		case "alac":
			conversionType = VIDEO_ALAC

		default:
			conversionType = INVALID_TYPE
		}
	} else {
		switch format.RequiredFileType {
		case "wav":
			conversionType = AUDIO_WAV

		case "mp3":
			conversionType = AUDIO_MP3

		case "ogg":
			conversionType = AUDIO_OGG

		case "opus":
			conversionType = AUDIO_OPUS

		case "flac":
			conversionType = AUDIO_FLAC

		case "aac":
			conversionType = AUDIO_AAC

		case "alac":
			conversionType = AUDIO_ALAC

		default:
			conversionType = INVALID_TYPE
		}
	}

	return conversionType
}
