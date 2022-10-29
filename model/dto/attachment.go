package dto

import (
	"github.com/go-sonic/sonic/consts"
)

type AttachmentDTO struct {
	ID             int32                 `json:"id"`
	Name           string                `json:"name"`
	Path           string                `json:"path"`
	FileKey        string                `json:"fileKey"`
	ThumbPath      string                `json:"thumbPath"`
	MediaType      string                `json:"mediaType"`
	Suffix         string                `json:"suffix"`
	Width          int32                 `json:"width"`
	Height         int32                 `json:"height"`
	Size           int64                 `json:"size"`
	AttachmentType consts.AttachmentType `json:"type"`
}
