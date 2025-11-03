package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ParikshitShetty/peercast/server/api"

	_ "github.com/ParikshitShetty/peercast/server/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Go Video Streaming API
// @version 1.0
// @description Simple Go video streaming service with Swagger docs.
// @host localhost:8080
// @BasePath /
func main() {
	http.HandleFunc("/video", api.StreamVideo)

	// Swagger endpoint
	http.Handle("/docs/", httpSwagger.WrapHandler)

	fmt.Println("🚀 Server running at http://localhost:8080")
	fmt.Println("📘 Swagger docs at http://localhost:8080/docs/index.html")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
