package ws

type ImageData struct {
	Url string `json:"url"`
	Id  string `json:"id"`
}

type DownloadImageRequest struct {
	Command string      `json:"command"`
	Data    []ImageData `json:"data"`
}
