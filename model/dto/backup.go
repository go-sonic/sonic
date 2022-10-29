package dto

type BackupDTO struct {
	DownloadLink string `json:"downloadLink"`
	Filename     string `json:"filename"`
	UpdateTime   int64  `json:"updateTime"`
	FileSize     int64  `json:"fileSize"`
}
