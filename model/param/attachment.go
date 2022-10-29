package param

import "github.com/go-sonic/sonic/consts"

type AttachmentQuery struct {
	Page
	Keyword        string                 `json:"keyword" form:"keyword"`
	MediaType      string                 `json:"mediaType" form:"mediaType"`
	AttachmentType *consts.AttachmentType `json:"attachmentType" form:"attachmentType"`
}

type AttachmentUpdate struct {
	Name string `json:"name" binding:"gte=1,lte=255"`
}
