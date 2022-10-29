package dto

import "github.com/go-sonic/sonic/consts"

type Post struct {
	PostMinimal
	Summary         string `json:"summary"`
	Thumbnail       string `json:"thumbnail"`
	Visits          int64  `json:"visits"`
	DisallowComment bool   `json:"disallowComment"`
	Password        string `json:"password"`
	Template        string `json:"template"`
	TopPriority     int32  `json:"topPriority"`
	Likes           int64  `json:"likes"`
	WordCount       int64  `json:"wordCount"`
	Topped          bool   `json:"topped"`
}

type PostMinimal struct {
	ID              int32             `json:"id"`
	Title           string            `json:"title"`
	Status          consts.PostStatus `json:"status"`
	Slug            string            `json:"slug"`
	EditorType      consts.EditorType `json:"editorType"`
	CreateTime      int64             `json:"createTime"`
	EditTime        int64             `json:"editTime"`
	UpdateTime      int64             `json:"updateTime"`
	MetaKeywords    string            `json:"metaKeywords"`
	MetaDescription string            `json:"metaDescription"`
	FullPath        string            `json:"fullPath"`
}

type PostDetail struct {
	Post
	OriginalContent string `json:"originalContent"`
	Content         string `json:"content"`
	CommentCount    int64  `json:"commentCount"`
}
