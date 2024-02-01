package DTO

import (
	"concurrent-image-fetcher/internal/models"
	"mime/multipart"
)

func GetAppFileFromFileHeader(file *multipart.FileHeader) *models.AppFile {
	return &models.AppFile{
		Name:     file.Filename,
		MimeType: file.Header.Get("Content-Type"),
		Size:     file.Size,
	}
}
