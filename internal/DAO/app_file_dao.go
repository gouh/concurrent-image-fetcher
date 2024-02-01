package DAO

import (
	"concurrent-image-fetcher/internal/models"
	"concurrent-image-fetcher/internal/requests"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type (
	AppFileDaoInterface interface {
		GetAll(params requests.PaginationRequest) (*int, *[]models.AppFile, error)
		Save(models.AppFile) error
		Delete(imageID string) error
		Get(fileId string) (*models.AppFile, error)
	}
	AppFileDao struct {
		db *sql.DB
	}
)

func (dao *AppFileDao) GetAll(params requests.PaginationRequest) (*int, *[]models.AppFile, error) {
	var totalItems int
	countQuery := "SELECT COUNT(*) FROM AppFile"
	err := dao.db.QueryRow(countQuery).Scan(&totalItems)
	if err != nil {
		return nil, nil, err
	}

	totalPages := totalItems / *params.PageSize
	if totalItems%*params.PageSize != 0 {
		totalPages++
	}

	offset := (*params.Page - 1) * *params.PageSize
	query := `SELECT * FROM AppFile ORDER BY created_at DESC LIMIT ? OFFSET ?`
	rows, err := dao.db.Query(query, *params.PageSize, offset)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	var appFiles []models.AppFile
	for rows.Next() {
		var appFile models.AppFile
		if err := rows.Scan(&appFile.ID, &appFile.Name, &appFile.MimeType, &appFile.Size, &appFile.Path, &appFile.CreatedAt, &appFile.UpdatedAt); err != nil {
			return nil, nil, err
		}
		appFiles = append(appFiles, appFile)
	}
	if err := rows.Err(); err != nil {
		return nil, nil, err
	}

	return &totalPages, &appFiles, nil
}

func (dao *AppFileDao) Get(fileId string) (*models.AppFile, error) {
	query := `SELECT * FROM AppFile WHERE id = ?`
	var appFile models.AppFile

	err := dao.db.QueryRow(query, fileId).Scan(&appFile.ID, &appFile.Name, &appFile.MimeType, &appFile.Size, &appFile.Path, &appFile.CreatedAt, &appFile.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("query row: %w", err)
	}

	return &appFile, nil
}

func (dao *AppFileDao) Delete(fileId string) error {
	_, err := dao.db.Exec("DELETE FROM AppFile WHERE id = ?", fileId)
	return err
}

func (dao *AppFileDao) Save(af models.AppFile) error {
	query := `INSERT INTO AppFile (id, name, mime_type, size, path) VALUES (?, ?, ?, ?, ?)`
	_, err := dao.db.Exec(query, af.ID, af.Name, af.MimeType, af.Size, af.Path)
	return err
}

func NewAppFileDao(db *sql.DB) AppFileDaoInterface {
	return &AppFileDao{
		db: db,
	}
}
