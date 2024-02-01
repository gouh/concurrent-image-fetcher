package ws

// DownloadProgress represents the data structure for download progress.
type DownloadProgress struct {
	Id        string  `json:"id"`
	Progress  float64 `json:"progress"`
	Completed bool    `json:"completed"`
	Error     string  `json:"error,omitempty"`
}
