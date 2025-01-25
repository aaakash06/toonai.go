package handlers

import (
	"archive/zip"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"
)

func GetComicChapter(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	chapterID := vars["chapterID"]

	// Get chapter path (in practice, you might get this from a database)
	chapterPath := filepath.Join("data", "comics", chapterID)

	// Check if the chapter exists
	if _, err := os.Stat(chapterPath); os.IsNotExist(err) {
		http.Error(w, "Chapter not found", http.StatusNotFound)
		return
	}

	// Create a temporary zip file
	tempFile, err := os.CreateTemp("", "chapter-*.zip")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer os.Remove(tempFile.Name())

	// Create zip writer
	zipWriter := zip.NewWriter(tempFile)
	defer zipWriter.Close()

	// Walk through the chapter directory and add files to zip
	err = filepath.Walk(chapterPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Only include image files
		if !isImageFile(path) {
			return nil
		}

		// Create zip file entry
		zipFile, err := zipWriter.Create(filepath.Base(path))
		if err != nil {
			return err
		}

		// Open and copy the file content
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(zipFile, file)
		return err
	})

	if err != nil {
		http.Error(w, "Error processing chapter", http.StatusInternalServerError)
		return
	}

	// Close the zip writer
	zipWriter.Close()

	// Serve the zip file
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", "attachment; filename=chapter.zip")
	http.ServeFile(w, r, tempFile.Name())
}

func isImageFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif" || ext == ".webp"
} 