package helper

import (
	"fmt"
	"mime/multipart"
	"strings"
)

func validateFile(maxFileSizeRequest int, file *multipart.FileHeader) error {
	// Batas ukuran file
	var maxFileSize = maxFileSizeRequest << 20

	if file.Size > int64(maxFileSize) {
		return fmt.Errorf("file size is too large, maxium is %vMB", maxFileSizeRequest)
	}

	// Validasi MIME type
	validMIMETypes := []string{"image/jpeg", "image/png", "image/gif"}
	fileType := file.Header.Get("Content-Type")
	isValidType := false
	for _, validType := range validMIMETypes {
		if strings.EqualFold(fileType, validType) {
			isValidType = true
			break
		}
	}

	if !isValidType {
		return fmt.Errorf("file type isn't valid, just only JPEG, PNG, and GIF are allowed")
	}

	return nil
}