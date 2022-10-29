package param

import "github.com/go-sonic/sonic/consts"

type Sheet struct {
	Title           string             `json:"title" form:"title" binding:"required,lte=100"`
	Status          consts.PostStatus  `json:"status" form:"status"`
	Slug            string             `json:"slug" form:"slug" binding:"lte=255"`
	EditorType      *consts.EditorType `json:"editorType" form:"editorType"`
	Content         string             `json:"content" form:"content"`
	OriginalContent string             `json:"originalContent" form:"originalContent"`
	Summary         string             `json:"summary" form:"summary" `
	Thumbnail       string             `json:"thumbnail" form:"thumbnail" binding:"lte=255"`
	DisallowComment bool               `json:"disallowComment" form:"disallowComment"`
	Password        string             `json:"password" form:"password"`
	Template        string             `json:"template"`
	TopPriority     int32              `json:"topPriority" form:"topPriority" binding:"gte=0"`
	CreateTime      *int64             `json:"createTime" form:"createTime"`
	MetaKeywords    string             `json:"metaKeywords" form:"metaKeywords"`
	MetaDescription string             `json:"metaDescription" form:"metaDescription"`
	Metas           []Meta             `json:"metas" form:"metas"`
}
