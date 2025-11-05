package database

import (
	"context"
	"log"

	"github.com/ParikshitShetty/peercast/server/internal/models"
)

// EnsureSchema ensures that the "files" table exists.
func EnsureSchema(ctx context.Context) error {
	db := Get() // returns *bun.DB

	// Automatically create the files table if it doesn’t exist
	_, err := db.NewCreateTable().
		Model((*models.File)(nil)).
		IfNotExists().
		Exec(ctx)

	if err != nil {
		return err
	}

	log.Println("✅ Schema verified or created successfully!")
	return nil
}
