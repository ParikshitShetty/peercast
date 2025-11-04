package configs

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type DBConfig struct {
	Host            string
	Port            int
	User            string
	Password        string
	DBName          string
	SSLMode         string
	MaxConns        int
	MinConns        int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

func Load() *DBConfig {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	getInt := func(key string, def int) int {
		if v, err := strconv.Atoi(os.Getenv(key)); err == nil && v > 0 {
			return v
		}
		return def
	}

	getDuration := func(key string, def time.Duration) time.Duration {
		if v := os.Getenv(key); v != "" {
			if d, err := time.ParseDuration(v); err == nil {
				return d
			}
		}
		return def
	}

	port := getInt("PG_PORT", 5432)

	cfg := &DBConfig{
		Host:            os.Getenv("PG_HOST"),
		Port:            port,
		User:            os.Getenv("PG_USER"),
		Password:        os.Getenv("PG_PASSWORD"),
		DBName:          os.Getenv("PG_DB"),
		SSLMode:         os.Getenv("PG_SSLMODE"),
		MaxConns:        getInt("PG_MAX_CONNS", 5),
		MinConns:        getInt("PG_MIN_CONNS", 1),
		ConnMaxLifetime: getDuration("PG_CONN_MAX_LIFETIME", 30*time.Minute),
		ConnMaxIdleTime: getDuration("PG_CONN_MAX_IDLE_TIME", 5*time.Minute),
	}

	log.Printf("Loaded DB config: host=%s port=%d db=%s max=%d min=%d\n",
		cfg.Host, cfg.Port, cfg.DBName, cfg.MaxConns, cfg.MinConns)

	return cfg
}
