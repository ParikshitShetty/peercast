package files

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ParikshitShetty/peercast/server/internal/repository"
)

func GetAll(w http.ResponseWriter, r *http.Request) {
	files, err := repository.GetAllFiles(context.Background())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(files)
	fmt.Println("files:", files)
}
