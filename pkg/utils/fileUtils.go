package utils

import (
	"image"
	"log"
	"mime/multipart"
)

func IsImageFile(file multipart.File) bool {
	_, format, err := image.Decode(file)
	if err != nil {
		log.Printf("failed to decode image: %v", err)
		return false
	}
	return format == "jpeg" || format == "png" || format == "gif" || format == "webp"
}
