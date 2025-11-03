package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func main() {
	http.HandleFunc("/video", streamVideo)
	fmt.Println("🚀 Server started on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// streamVideo handles HTTP range requests for video streaming
func streamVideo(w http.ResponseWriter, r *http.Request) {
	videoPath := "videos/earth.mp4"

	file, err := os.Open(videoPath)
	if err != nil {
		http.Error(w, "Unable to open video", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	fileStat, err := file.Stat()
	if err != nil {
		http.Error(w, "Unable to get file info", http.StatusInternalServerError)
		return
	}

	fileSize := fileStat.Size()
	rangeHeader := r.Header.Get("Range")

	if rangeHeader == "" {
		// Serve entire file if no Range header is present
		w.Header().Set("Content-Type", "video/mp4")
		w.Header().Set("Content-Length", strconv.FormatInt(fileSize, 10))
		http.ServeContent(w, r, videoPath, fileStat.ModTime(), file)
		return
	}

	// Handle range request (e.g., bytes=2000-)
	parts := strings.Split(strings.Replace(rangeHeader, "bytes=", "", 1), "-")
	start, _ := strconv.ParseInt(parts[0], 10, 64)
	end := fileSize - 1
	if len(parts) == 2 && parts[1] != "" {
		end, _ = strconv.ParseInt(parts[1], 10, 64)
	}

	if start > end || start < 0 || end >= fileSize {
		http.Error(w, "Invalid range", http.StatusRequestedRangeNotSatisfiable)
		return
	}

	chunkSize := end - start + 1
	w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, fileSize))
	w.Header().Set("Accept-Ranges", "bytes")
	w.Header().Set("Content-Length", strconv.FormatInt(chunkSize, 10))
	w.Header().Set("Content-Type", "video/mp4")
	w.WriteHeader(http.StatusPartialContent)

	file.Seek(start, 0)
	buf := make([]byte, 1024*32)
	var sent int64
	for {
		if sent >= chunkSize {
			break
		}
		remaining := chunkSize - sent
		if int64(len(buf)) > remaining {
			buf = buf[:remaining]
		}
		n, err := file.Read(buf)
		if err != nil {
			break
		}
		w.Write(buf[:n])
		sent += int64(n)
	}
}
