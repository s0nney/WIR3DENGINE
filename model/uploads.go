package model

import (
	"crypto/md5"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

func FnMixer(originalFilename string) string {
	fileExtension := filepath.Ext(originalFilename)

	randomString := uuid.New().String()

	newFilename := strings.TrimSuffix(randomString, "-") + fileExtension

	return newFilename
}

func TrimFn(filename string, maxLength int) string {
	fileExt := filepath.Ext(filename)
	if len(filename) > maxLength {
		return filename[:maxLength] + fileExt
	}
	return filename
}

func ComputeChecksum(filePath string) string {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return ""
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		fmt.Println("Error copying file content:", err)
		return ""
	}

	return fmt.Sprintf("%x", hash.Sum(nil))
}

func ReadChecksum(file *multipart.FileHeader) (string, error) {
	f, err := file.Open()
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
