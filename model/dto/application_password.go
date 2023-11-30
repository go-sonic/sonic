package dto

type ApplicationPasswordDTO struct {
	Name             string `json:"name"`
	Password         string `json:"password"`
	LastActivateTime int64  `json:"last_activate_time"`
	LastActiveIP     string `json:"last_activate_ip"`
	CreateTime       int64  `json:"create_time"`
}
