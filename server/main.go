package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/ParikshitShetty/peercast/server/api"
	"github.com/ParikshitShetty/peercast/server/internal/configs"
	"github.com/ParikshitShetty/peercast/server/internal/database"
	mw_api "github.com/ParikshitShetty/peercast/server/internal/middleware/api"
	mw_log "github.com/ParikshitShetty/peercast/server/internal/middleware/log"

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

	// 🗄️ Initialize DB
	if err := database.Init(ctx, cfg); err != nil {
		log.Fatalf("❌ DB init failed: %v", err)
	}
	defer database.Close()

	if err := database.EnsureSchema(ctx); err != nil {
		log.Fatalf("❌ Schema setup failed: %v", err)
	}

	mux := http.NewServeMux()
	api.RegisterRoutes(mux)

	mux.Handle("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))

	apiHandler := mw_api.RecoverPanic(
		mw_log.Logger(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			mux.ServeHTTP(w, r)
		})),
	)

	fmt.Println("🚀 Server running at http://localhost:8080")
	fmt.Println("📘 Swagger docs at http://localhost:8080/swagger/index.html")

	if err := http.ListenAndServe(":8080", apiHandler); err != nil {
		log.Fatalf("❌ Server failed: %v", err)
	}
}
