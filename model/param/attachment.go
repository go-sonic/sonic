package param

import "github.com/go-sonic/sonic/consts"

type AttachmentQuery struct {
	Page
	Keyword        string                 `json:"keyword" form:"keyword"`
	MediaType      string                 `json:"mediaType" form:"mediaType"`
	AttachmentType *consts.AttachmentType `json:"attachmentType" form:"attachmentType"`
}

type AttachmentQueryNoEnum struct {
	Page
	Keyword        string `json:"keyword" form:"keyword"`
	MediaType      string `json:"mediaType" form:"mediaType"`
	AttachmentType string `json:"attachmentType" form:"attachmentType"`
}

type AttachmentUpdate struct {
	Name string `json:"name" binding:"gte=1,lte=255"`
}

func AssertAttachmentQuery(t AttachmentQueryNoEnum) AttachmentQuery {
	res := AttachmentQuery{
		Page:      t.Page,
		Keyword:   t.Keyword,
		MediaType: t.MediaType,
	}
	str := t.AttachmentType
	switch str {
	case `"LOCAL"`:
		*res.AttachmentType = consts.AttachmentTypeLocal
	case `"UPOSS"`:
		*res.AttachmentType = consts.AttachmentTypeUpOSS
	case `"QINIUOSS"`:
		*res.AttachmentType = consts.AttachmentTypeQiNiuOSS
	case `"AttachmentTypeSMMS"`:
		*res.AttachmentType = consts.AttachmentTypeSMMS
	case `"ALIOSS"`:
		*res.AttachmentType = consts.AttachmentTypeAliOSS
	case `"BAIDUBOS"`:
		*res.AttachmentType = consts.AttachmentTypeBaiDuOSS
	case `"TENCENTCOS"`:
		*res.AttachmentType = consts.AttachmentTypeTencentCOS
	case `"HUAWEIOBS"`:
		*res.AttachmentType = consts.AttachmentTypeHuaweiOBS
	case `"MINIO"`:
		*res.AttachmentType = consts.AttachmentTypeMinIO
	}
	return res
}
