package models

import (
	"time"

	"github.com/uptrace/bun"
)

type File struct {
	bun.BaseModel `bun:"table:files"`

	ID           int64     `bun:",pk,autoincrement"`
	FileName     string    `bun:"file_name,notnull"`
	Size         int64     `bun:"size"`
	FileLocation string    `bun:"file_location,notnull"`
	CreatedAt    time.Time `bun:"created_at,default:current_timestamp"`
}
