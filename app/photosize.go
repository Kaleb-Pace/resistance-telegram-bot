package main

// PhotoSize contains information about photos.
type PhotoSize struct {
	FileID   string `json:"file_id"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
	FileSize int    `json:"file_size"` // optional
}

// GetImageResponse
type GetFileResponse struct {
	Ok     bool                `json:"ok"`
	Result GetFileResponseData `json:"result"`
}

// GetImageResponseData
type GetFileResponseData struct {
	FilePath string `json:"file_path"`
	FileSize int    `json:"file_size"`
	ID       string `json:"file_id"`
}
