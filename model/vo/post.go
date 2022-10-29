package vo

import "github.com/go-sonic/sonic/model/dto"

type Post struct {
	dto.Post
	CommentCount int64                  `json:"commentCount"`
	Tags         []*dto.Tag             `json:"tags"`
	Categories   []*dto.CategoryDTO     `json:"categories"`
	Metas        map[string]interface{} `json:"metas"`
}

type PostDetailVO struct {
	dto.PostDetail
	TagIDs      []int32            `json:"tagIds"`
	Tags        []*dto.Tag         `json:"tags"`
	CategoryIDs []int32            `json:"categoryIds"`
	Categories  []*dto.CategoryDTO `json:"categories"`
	MetaIDs     []int64            `json:"metaIds"`
	Metas       []*dto.Meta        `json:"metas"`
}
