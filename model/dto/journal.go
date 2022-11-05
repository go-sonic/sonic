package dto

import "github.com/go-sonic/sonic/consts"

type Journal struct {
	ID            int32              `json:"id"`
	SourceContent string             `json:"sourceContent"`
	Content       string             `json:"content"`
	Likes         int64              `json:"likes"`
	CreateTime    int64              `json:"createTime"`
	JournalType   consts.JournalType `json:"type"`
}

type JournalWithComment struct {
	Journal
	CommentCount int64 `json:"commentCount"`
}
