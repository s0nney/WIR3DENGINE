package model

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"WIR3DENGINE/utils"

	"github.com/jinzhu/gorm"
)

type Paginate struct {
	gorm.Model
}

type PageVariables struct {
	Title  string
	Images []string
}

type Image struct {
	Filename string
}

type ImageFile struct {
	Name         string
	CreationTime time.Time
}

func GetImageFiles(directory string) ([]ImageFile, error) {
	var imageFiles []ImageFile

	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && isImageFile(info.Name()) {
			imageFiles = append(imageFiles, ImageFile{
				Name:         info.Name(),
				CreationTime: info.ModTime(),
			})
		}
		return nil
	})

	return imageFiles, err
}

func SortImageFilesByCreationTime(imageFiles []ImageFile) {
	sort.Slice(imageFiles, func(i, j int) bool {
		return imageFiles[i].CreationTime.Before(imageFiles[j].CreationTime)
	})
}

func GetRecentImagesFromDB(maxCount int) ([]string, error) {
	rows, err := utils.Db.Query("SELECT filename FROM images ORDER BY id DESC LIMIT $1", maxCount)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var recentImages []string

	for rows.Next() {
		var image string
		err := rows.Scan(&image)
		if err != nil {
			return nil, err
		}
		recentImages = append(recentImages, image)
	}

	return recentImages, nil
}

func GetRandomImagesFromDB(maxCount int) ([]string, error) {
	rows, err := utils.Db.Query("SELECT filename FROM images ORDER BY RANDOM() DESC LIMIT $1", maxCount)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var recentImages []string

	for rows.Next() {
		var image string
		err := rows.Scan(&image)
		if err != nil {
			return nil, err
		}
		recentImages = append(recentImages, image)
	}

	return recentImages, nil
}

func isImageFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif" || ext == ".bmp" || ext == ".webp"
}

func ManipulateImage(inputFile, outputFile string) error {
	cmd := exec.Command("convert", inputFile, "-resize", "50%", outputFile)

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("ImageMagick command failed: %v", err)
	}

	return nil
}
