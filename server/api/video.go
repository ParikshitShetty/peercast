package api

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// streamVideo godoc
// @Summary Stream a video file
// @Description Streams MP4 video content with support for range requests
// @Tags Video
// @Produce video/mp4
// @Success 206 {string} string "Partial Content"
// @Failure 404 {string} string "Video not found"
// @Router /video [get]
func StreamVideo(w http.ResponseWriter, r *http.Request) {
	videoPath := "videos/earth.mp4"
	file, err := os.Open(videoPath)
	if err != nil {
		http.Error(w, "Video not found", http.StatusNotFound)
		return
	}
	defer file.Close()

	fileStat, _ := file.Stat()
	fileSize := fileStat.Size()
	rangeHeader := r.Header.Get("Range")

	if rangeHeader == "" {
		w.Header().Set("Content-Type", "video/mp4")
		w.Header().Set("Content-Length", strconv.FormatInt(fileSize, 10))
		http.ServeContent(w, r, videoPath, fileStat.ModTime(), file)
		return
	}

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
	fmt.Println("Video streaming handler initialized")
}
