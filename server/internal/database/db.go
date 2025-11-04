package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ParikshitShetty/peercast/server/internal/configs"
	"github.com/jackc/pgx/v5/pgxpool"
)

var pool *pgxpool.Pool

// Init initializes the global database connection pool.
func Init(ctx context.Context, cfg *configs.DBConfig) error {
	dsn := fmt.Sprintf(
		"user=%s password=%s host=%s port=%d dbname=%s sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode,
	)

	poolCfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return fmt.Errorf("parse config: %w", err)
	}

	// Apply tuning / limits
	poolCfg.MaxConns = int32(cfg.MaxConns)
	poolCfg.MinConns = int32(cfg.MinConns)
	poolCfg.MaxConnLifetime = 30 * time.Minute
	poolCfg.MaxConnIdleTime = 5 * time.Minute

	// Connect
	dbpool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return fmt.Errorf("connect: %w", err)
	}

	// Test connection
	if err := dbpool.Ping(ctx); err != nil {
		dbpool.Close()
		return fmt.Errorf("ping: %w", err)
	}

	pool = dbpool
	log.Println("✅ Database connected successfully!")
	return nil
}

// GetPool returns the global pool instance
func GetPool() *pgxpool.Pool {
	if pool == nil {
		log.Fatal("database pool not initialized — call database.Init() first")
	}
	return pool
}

// Close terminates the pool connection
func Close() {
	if pool != nil {
		log.Println("🧹 Closing database pool...")
		pool.Close()
	}
}
