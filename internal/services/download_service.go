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
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
)

// DownloadService provides functionalities to download images and publish progress.
type (
	DownloadServiceInterface interface {
		DownloadImage(context.Context, requestsWs.ImageData, string)
		ProcessImages([]requestsWs.ImageData, string)
	}
	DownloadService struct {
		RedisPubSub utils.PubSubInterface
		AppFileDao  DAO.AppFileDaoInterface
		ImagePath   string
	}
)

func getFilenameFromURL(imageURL string) string {
	parsedURL, err := url.Parse(imageURL)
	if err != nil {
		return ""
	}
	filename := path.Base(parsedURL.Path)
	decodedFilename, err := url.QueryUnescape(filename)
	if err != nil {
		return filename
	}

	return decodedFilename
}

func getFilenameFromContentDisposition(cd string) string {
	re := regexp.MustCompile(`filename\*?="?([^";]+)"?;?`)
	matches := re.FindStringSubmatch(cd)
	for i, name := range re.SubexpNames() {
		if name == "filename" && i < len(matches) {
			return matches[i]
		}
	}
	return ""
}

func getFilename(cd string, imageURL string) string {
	filename := getFilenameFromContentDisposition(cd)
	if filename == "" {
		filename = getFilenameFromURL(imageURL)
	}
	return filename
}

// DownloadImage downloads an image from the specified URL and saves it locally.
// It also publishes download progress updates using Redis Pub/Sub.
func (service *DownloadService) DownloadImage(ctx context.Context, url requestsWs.ImageData, roomId string) {
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
	filename := getFilename(cd, url.Url)

	contentType := response.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		errMsg := fmt.Sprintf("File isnt a image: %s, Content-Type: %s", url.Url, contentType)
		errPub := service.RedisPubSub.PublishProgress(ctx, roomId, responsesWs.DownloadProgress{Id: url.Id, Error: errMsg, Event: "progress"})
		if errPub != nil {
			fmt.Println(errPub)
		}
		return
	}

	totalSize := response.ContentLength
	receivedSize := 0

	newFilePath := filepath.Join(service.ImagePath, fmt.Sprintf("%s"+filepath.Ext(filename), url.Id))
	file, err := os.Create(newFilePath)
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
			errPub := service.RedisPubSub.PublishProgress(ctx, roomId, responsesWs.DownloadProgress{Id: url.Id, Progress: progress, Completed: false, Event: "progress"})
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

	errPub := service.RedisPubSub.PublishProgress(ctx, roomId, responsesWs.DownloadProgress{Id: url.Id, Progress: 100, Completed: true, Event: "progress"})
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
	} else {
		fileMini := url.Id + "_mini" + filepath.Ext(filename)
		fileMiniPath := filepath.Join(service.ImagePath, fileMini)
		go func() {
			errResize := utils.ResizeImage(newFilePath, fileMiniPath)
			if errResize != nil {
				fmt.Println(errResize)
			}
		}()
	}

}

func (service *DownloadService) ProcessImages(images []requestsWs.ImageData, roomId string) {
	var wg sync.WaitGroup
	ctx := context.Background()
	counter := 0
	for i, image := range images {
		wg.Add(1)
		go func(img requestsWs.ImageData) {
			defer wg.Done()
			service.DownloadImage(ctx, img, roomId)
		}(image)

		counter++

		if counter%5 == 0 && i > 0 {
			time.Sleep(10 * time.Second)
		}
	}
	wg.Wait()

	errPub := service.RedisPubSub.PublishProgress(ctx, roomId, responsesWs.DownloadProgress{Id: "none", Progress: 100, Completed: true, Event: "final"})
	if errPub != nil {
		fmt.Println(errPub)
	}
}

// NewDownloadService creates a new instance of DownloadService.
func NewDownloadService(container *container.Container) DownloadServiceInterface {
	return &DownloadService{
		RedisPubSub: utils.NewPubSub(container.Redis),
		AppFileDao:  DAO.NewAppFileDao(container.Db),
		ImagePath:   container.Config.Image.ImagePath,
	}
}
