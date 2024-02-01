package services

import (
	"concurrent-image-fetcher/internal/DAO"
	"concurrent-image-fetcher/internal/container"
	"concurrent-image-fetcher/internal/models"
	requestsWs "concurrent-image-fetcher/internal/requests/ws"
	responsesWs "concurrent-image-fetcher/internal/responses/ws"
	"concurrent-image-fetcher/internal/utils"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

// DownloadService provides functionalities to download images and publish progress.
type (
	DownloadServiceInterface interface {
		DownloadImage(context.Context, requestsWs.ImageData, *sync.WaitGroup, string)
		ProcessImages([]requestsWs.ImageData, string)
	}
	DownloadService struct {
		RedisPubSub utils.PubSubInterface
		AppFileDao  DAO.AppFileDaoInterface
		ImagePath   string
	}
)

// getFilenameFromContentDisposition extrae el nombre del archivo del encabezado Content-Disposition.
func getFilenameFromContentDisposition(cd string) string {
	re := regexp.MustCompile(`filename="?(?P<filename>[^"]+)"?`)
	matches := re.FindStringSubmatch(cd)
	for i, name := range re.SubexpNames() {
		if name == "filename" && i < len(matches) {
			return matches[i]
		}
	}
	return ""
}

// DownloadImage downloads an image from the specified URL and saves it locally.
// It also publishes download progress updates using Redis Pub/Sub.
func (service *DownloadService) DownloadImage(ctx context.Context, url requestsWs.ImageData, wg *sync.WaitGroup, roomId string) {
	defer wg.Done()

	response, err := http.Get(url.Url)
	if err != nil {
		fmt.Println(url)
		fmt.Println("Error al descargar la imagen:", err)
		return
	}
	defer response.Body.Close()

	errFolder := os.MkdirAll(service.ImagePath, os.ModePerm)
	if errFolder != nil {
		fmt.Println("Error on folder creation:", errFolder)
		return
	}

	cd := response.Header.Get("Content-Disposition")
	filename := getFilenameFromContentDisposition(cd)

	contentType := response.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		errMsg := fmt.Sprintf("File isnt a image: %s, Content-Type: %s", url.Url, contentType)
		errPub := service.RedisPubSub.PublishProgress(ctx, roomId, responsesWs.DownloadProgress{Id: url.Id, Error: errMsg})
		if errPub != nil {
			fmt.Println(errPub)
		}
		return
	}

	totalSize := response.ContentLength
	receivedSize := 0

	file, err := os.Create(filepath.Join(service.ImagePath, fmt.Sprintf("%s.jpg", url.Id)))
	if err != nil {
		fmt.Println("Error al crear el archivo:", err)
		return
	}
	defer file.Close()

	buffer := make([]byte, 1024)
	for {
		n, err := response.Body.Read(buffer)
		if n > 0 {
			receivedSize += n
			_, err = file.Write(buffer[:n])
			if err != nil {
				fmt.Println("Error al escribir la imagen:", err)
				return
			}
			progress := float64(receivedSize) / float64(totalSize) * 100
			errPub := service.RedisPubSub.PublishProgress(ctx, roomId, responsesWs.DownloadProgress{Id: url.Id, Progress: progress, Completed: false})
			if errPub != nil {
				fmt.Println(errPub)
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("Error al leer la imagen:", err)
			return
		}
	}

	errPub := service.RedisPubSub.PublishProgress(ctx, roomId, responsesWs.DownloadProgress{Id: url.Id, Progress: 100, Completed: true})
	if errPub != nil {
		fmt.Println(errPub)
	}

	errOnSave := service.AppFileDao.Save(models.AppFile{
		ID:       url.Id,
		Name:     filename,
		MimeType: contentType,
		Size:     totalSize,
		Path:     strings.ReplaceAll(service.ImagePath, ".", ""),
	})
	if errOnSave != nil {
		fmt.Println(errOnSave)
	}
}

func (service *DownloadService) ProcessImages(images []requestsWs.ImageData, roomId string) {
	var wg sync.WaitGroup
	ctx := context.Background()
	for _, image := range images {
		wg.Add(1)
		go service.DownloadImage(ctx, image, &wg, roomId)
	}
	wg.Wait()
}

// NewDownloadService creates a new instance of DownloadService.
func NewDownloadService(container *container.Container) DownloadServiceInterface {
	return &DownloadService{
		RedisPubSub: utils.NewPubSub(container.Redis),
		AppFileDao:  DAO.NewAppFileDao(container.Db),
		ImagePath:   container.Config.Image.ImagePath,
	}
}
