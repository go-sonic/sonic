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
	EditTime        *int64             `json:"editTime" form:"editTime"`
	UpdateTime      *int64             `json:"updateTime" form:"updateTime"`
}

type PostContent struct {
	Content         string `json:"content" form:"content"`
	OriginalContent string `json:"originalContent" form:"orginalContent"`
}

type PostQueryNoEnum struct {
	Page
	*Sort
	Keyword      *string  `json:"keyword" form:"keyword"`
	Statuses     []string `json:"statuses" form:"statuses"`
	CategoryID   *int32   `json:"categoryId" form:"categoryId"`
	More         *bool    `json:"more" form:"more"`
	TagID        *int32   `json:"tagId" form:"tagId"`
	WithPassword *bool    `json:"-" form:"-"`
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

func AssertPostQuery(t PostQueryNoEnum) PostQuery {
	res := PostQuery{
		Page:         t.Page,
		Sort:         t.Sort,
		Keyword:      t.Keyword,
		CategoryID:   t.CategoryID,
		More:         t.More,
		TagID:        t.TagID,
		WithPassword: t.WithPassword,
	}
	var statues []*consts.PostStatus
	for _, str := range t.Statuses {
		var status consts.PostStatus
		switch str {
		case "PUBLISHED":
			status = consts.PostStatusPublished
		case "DRAFT":
			status = consts.PostStatusDraft
		case "INTIMATE":
			status = consts.PostStatusIntimate
		}
		statues = append(statues, &status)
	}
	res.Statuses = statues
	return res
}
