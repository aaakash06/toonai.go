package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

func StreamVideo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	videoID := vars["videoID"]

	// Get video path (in practice, you might get this from a database)
	videoPath := filepath.Join("data", "videos", videoID)

	// Check if the video exists
	videoFile, err := os.Open(videoPath)
	if err != nil {
		http.Error(w, "Video not found", http.StatusNotFound)
		return
	}
	defer videoFile.Close()

	// Get video file information
	fileInfo, err := videoFile.Stat()
	if err != nil {
		http.Error(w, "Error reading video", http.StatusInternalServerError)
		return
	}

	// Get the file size
	fileSize := fileInfo.Size()

	// Get the range header
	rangeHeader := r.Header.Get("Range")
	if rangeHeader == "" {
		// Serve the full file if no range is specified
		w.Header().Set("Content-Type", "video/mp4")
		w.Header().Set("Content-Length", strconv.FormatInt(fileSize, 10))
		w.WriteHeader(http.StatusOK)
		io.Copy(w, videoFile)
		return
	}

	// Parse the range header
	ranges := strings.Split(strings.TrimPrefix(rangeHeader, "bytes="), "-")
	if len(ranges) != 2 {
		http.Error(w, "Invalid range header", http.StatusBadRequest)
		return
	}

	// Parse start and end positions
	start, err := strconv.ParseInt(ranges[0], 10, 64)
	if err != nil {
		http.Error(w, "Invalid range header", http.StatusBadRequest)
		return
	}

	var end int64
	if ranges[1] == "" {
		end = fileSize - 1
	} else {
		end, err = strconv.ParseInt(ranges[1], 10, 64)
		if err != nil {
			http.Error(w, "Invalid range header", http.StatusBadRequest)
			return
		}
	}

	// Validate range
	if start >= fileSize || end >= fileSize {
		w.Header().Set("Content-Range", fmt.Sprintf("bytes */%d", fileSize))
		http.Error(w, "Requested range not satisfiable", http.StatusRequestedRangeNotSatisfiable)
		return
	}

	// Set headers for partial content
	w.Header().Set("Content-Type", "video/mp4")
	w.Header().Set("Accept-Ranges", "bytes")
	w.Header().Set("Content-Length", strconv.FormatInt(end-start+1, 10))
	w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, fileSize))
	w.WriteHeader(http.StatusPartialContent)

	// Seek to start position
	_, err = videoFile.Seek(start, 0)
	if err != nil {
		http.Error(w, "Error reading video", http.StatusInternalServerError)
		return
	}

	// Stream the video chunk
	_, err = io.CopyN(w, videoFile, end-start+1)
	if err != nil {
		log.Printf("Error streaming video: %v", err)
	}
} 