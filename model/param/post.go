package param

import "github.com/go-sonic/sonic/consts"

type Post struct {
	Title           string             `json:"title" form:"title" binding:"gte=1,lte=100"`
	Status          consts.PostStatus  `json:"status" form:"status" binding:"gte=0"`
	Slug            string             `json:"slug" form:"slug" binding:"lte=255"`
	EditorType      *consts.EditorType `json:"editorType" form:"editorType"`
	OriginalContent string             `json:"originalContent" form:"originalContent"`
	Summary         string             `json:"summary" form:"summary"`
	Thumbnail       string             `json:"thumbnail" form:"thumbnail"`
	DisallowComment bool               `json:"disallowComment" form:"disallowComment"`
	Password        string             `json:"password" form:"password" binding:"lte=255"`
	Template        string             `json:"template" form:"template" binding:"lte=255"`
	TopPriority     int32              `json:"topPriority" form:"topPriority" binding:"gte=0"`
	CreateTime      *int64             `json:"createTime" form:"createTime" `
	MetaKeywords    string             `json:"metaKeywords" form:"metaKeywords" `
	MetaDescription string             `json:"metaDescription" form:"metaDescription"`
	TagIDs          []int32            `json:"tagIds" form:"tagIds"`
	CategoryIDs     []int32            `json:"categoryIds" form:"categoryIds"`
	MetaParam       []Meta             `json:"metas" form:"metas"`
	Content         string             `json:"content" form:"content"`
}

type PostContent struct {
	Content string `json:"content" form:"content"`
}

type PostQuery struct {
	Page
	*Sort
	Keyword      *string              `json:"keyword" form:"keyword"`
	Statuses     []*consts.PostStatus `json:"statuses" form:"statuses"`
	CategoryID   *int32               `json:"categoryId" form:"categoryId"`
	More         *bool                `json:"more" form:"more"`
	TagID        *int32               `json:"tagId" form:"tagId"`
	WithPassword *bool                `json:"-" form:"-"`
}
