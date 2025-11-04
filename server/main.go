package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/ParikshitShetty/peercast/server/api/streaming"
	"github.com/ParikshitShetty/peercast/server/internal/configs"
	"github.com/ParikshitShetty/peercast/server/internal/database"

	_ "github.com/ParikshitShetty/peercast/server/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Go Video Streaming API
// @version 1.0
// @description Simple Go video streaming service with Swagger docs.
// @host localhost:8080
// @BasePath /
func main() {
	ctx := context.Background()
	cfg := configs.Load()

	// Initialize DB connection on server startup
	if err := database.Init(ctx, cfg); err != nil {
		log.Fatalf("db init failed: %v", err)
	}
	defer database.Close()

	http.HandleFunc("/video", streaming.StreamVideo)

	// Swagger endpoint
	http.Handle("/docs/", httpSwagger.WrapHandler)

	fmt.Println("🚀 Server running at http://localhost:8080")
	fmt.Println("📘 Swagger docs at http://localhost:8080/docs/index.html")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
