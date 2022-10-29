package dto

type EnvironmentDTO struct {
	Database  string `json:"database"`
	StartTime int64  `json:"startTime"`
	Version   string `json:"version"`
	Mode      string `json:"mode"`
}
