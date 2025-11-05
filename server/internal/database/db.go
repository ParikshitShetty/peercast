package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/ParikshitShetty/peercast/server/internal/configs"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

var DB *bun.DB

// Init initializes a global Bun DB connection.
func Init(ctx context.Context, cfg *configs.DBConfig) error {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode,
	)

	// Create sql.DB
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))

	// Tuning
	sqldb.SetMaxOpenConns(cfg.MaxConns)
	sqldb.SetMaxIdleConns(cfg.MinConns)
	sqldb.SetConnMaxLifetime(30 * time.Minute)
	sqldb.SetConnMaxIdleTime(5 * time.Minute)

	// Verify connection
	if err := sqldb.PingContext(ctx); err != nil {
		return fmt.Errorf("ping: %w", err)
	}

	// ✅ Wrap with Bun using PostgreSQL dialect
	DB = bun.NewDB(sqldb, pgdialect.New())

	log.Println("✅ Bun DB connected successfully!")
	return nil
}

func Get() *bun.DB {
	if DB == nil {
		log.Fatal("database not initialized — call database.Init() first")
	}
	return DB
}

func Close() {
	if DB != nil {
		log.Println("🧹 Closing Bun DB connection...")
		if err := DB.DB.Close(); err != nil {
			log.Printf("error closing DB: %v", err)
		}
	}
}
