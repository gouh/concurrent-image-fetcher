package handlers

import (
	"concurrent-image-fetcher/internal/DAO"
	"concurrent-image-fetcher/internal/DTO"
	"concurrent-image-fetcher/internal/container"
	"concurrent-image-fetcher/internal/requests"
	"concurrent-image-fetcher/internal/responses"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type (
	ImageHandlerInterface interface {
		GetImages(c *gin.Context)
		GetImage(c *gin.Context)
		PostImage(c *gin.Context)
		DeleteImage(c *gin.Context)
	}
	ImageHandler struct {
		AppFileDao DAO.AppFileDaoInterface
		ImagePath  string
	}
)

func isImage(file *multipart.FileHeader) bool {
	mimeType := file.Header.Get("Content-Type")
	extensions := []string{"jpg", "jpeg", "png", "gif", "bmp", "svg"}
	for _, ext := range extensions {
		if mimeType == mime.TypeByExtension("."+ext) {
			return true
		}
	}
	return false
}

func (handler *ImageHandler) GetImages(c *gin.Context) {
	var params requests.PaginationRequest
	params = params.GetDefaultPaginationParams(c)
	totalPages, items, err := handler.AppFileDao.GetAll(params.GetDefaultPaginationParams(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.GetErrorResponse("Error on get app files"))
	}

	c.JSON(http.StatusOK, responses.GetImagesResponse(params, totalPages, items))
}

func (handler *ImageHandler) GetImage(c *gin.Context) {
	appFileId := c.Param("id")

	appFile, errAppFile := handler.AppFileDao.Get(appFileId)
	if errAppFile != nil {
		c.JSON(http.StatusInternalServerError, responses.GetErrorResponse("Error on get app files"))
	}

	if appFile == nil {
		c.JSON(http.StatusNotFound, responses.GetErrorResponse("AppFile not found"))
		return
	}

	c.JSON(http.StatusOK, responses.CommonResponse{
		Data: appFile,
	})
}

func (handler *ImageHandler) PostImage(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 10<<20)

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.GetErrorResponse("Error on get file"))
		return
	}

	if !isImage(file) {
		c.JSON(http.StatusBadRequest, responses.GetErrorResponse("File is not image"))
		return
	}

	fileUUID := uuid.New().String()
	fileName := fileUUID + filepath.Ext(file.Filename)
	filePath := filepath.Join(handler.ImagePath, fileName)
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, responses.GetErrorResponse("Error on save file"))
		return
	}

	appFile := DTO.GetAppFileFromFileHeader(file)
	appFile.ID = fileName
	appFile.Path = strings.ReplaceAll(handler.ImagePath, ".", "")
	if err := handler.AppFileDao.Save(*appFile); err != nil {
		c.JSON(http.StatusInternalServerError, responses.GetErrorResponse("Error on save metadata"))
		return
	}

	appFileStored, errAppFile := handler.AppFileDao.Get(appFile.ID)
	if errAppFile != nil {
		c.JSON(http.StatusInternalServerError, responses.GetErrorResponse("Error on get app files"))
	}

	c.JSON(http.StatusCreated, responses.CommonResponse{
		Data: appFileStored,
	})
}

func (handler *ImageHandler) DeleteImage(c *gin.Context) {
	appFileId := c.Param("id")
	imagePath := fmt.Sprintf(filepath.Join(handler.ImagePath, "%s.jpg"), appFileId)
	thumbnailPath := fmt.Sprintf(filepath.Join(handler.ImagePath, "%s_thumb.jpg"), appFileId)

	appFile, errAppFile := handler.AppFileDao.Get(appFileId)
	if errAppFile != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error on get app file"})
	}

	if appFile == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "AppFile not found"})
		return
	}

	deleteIfExists := func(path string) {
		if _, err := os.Stat(path); err == nil || !os.IsNotExist(err) {
			if err := os.Remove(path); err != nil {
				fmt.Println("Error al eliminar el archivo:", err)
			}
		}
	}

	deleteIfExists(imagePath)
	deleteIfExists(thumbnailPath)

	if err := handler.AppFileDao.Delete(appFileId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar la imagen de la base de datos"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func NewImageHandler(container *container.Container) ImageHandlerInterface {
	return &ImageHandler{
		AppFileDao: DAO.NewAppFileDao(container.Db),
		ImagePath:  container.Config.Image.ImagePath,
	}
}
