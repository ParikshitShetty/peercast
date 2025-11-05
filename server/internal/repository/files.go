package repository

import (
	"context"

	"github.com/ParikshitShetty/peercast/server/internal/database"
	"github.com/ParikshitShetty/peercast/server/internal/models"
)

func CreateFile(ctx context.Context, file *models.File) error {
	db := database.Get()
	_, err := db.NewInsert().Model(file).Exec(ctx)
	return err
}

func GetFileByID(ctx context.Context, id int64) (*models.File, error) {
	db := database.Get()
	file := new(models.File)
	err := db.NewSelect().Model(file).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func DeleteFileByID(ctx context.Context, id int64) error {
	db := database.Get()
	_, err := db.NewDelete().Model((*models.File)(nil)).Where("id = ?", id).Exec(ctx)
	return err
}

func GetAllFiles(ctx context.Context) ([]*models.File, error) {
	db := database.Get()
	var files []*models.File
	err := db.NewSelect().Model(&files).Order("created_at DESC").Scan(ctx)
	if err != nil {
		return nil, err
	}
	return files, nil
}

func UpdateFile(ctx context.Context, file *models.File) error {
	db := database.Get()
	_, err := db.NewUpdate().Model(file).WherePK().Exec(ctx)
	return err
}
