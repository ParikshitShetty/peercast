package api

import (
	"net/http"

	"github.com/ParikshitShetty/peercast/server/api/files"
	"github.com/ParikshitShetty/peercast/server/api/streaming"
)

func RegisterRoutes(mux *http.ServeMux) {
	// File CRUD
	// mux.HandleFunc("/api/files/create", files.)
	mux.HandleFunc("/api/files/list", files.GetAll)
	// mux.HandleFunc("/api/files/get", repository.GetFileByID)
	// mux.HandleFunc("/api/files/update", repository.UpdateFile())
	// mux.HandleFunc("/api/files/delete", repository.DeleteFileByID())
	//
	// Video streaming
	mux.HandleFunc("/api/streaming/video", streaming.StreamVideo)
}
