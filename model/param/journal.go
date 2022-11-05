package param

import "github.com/go-sonic/sonic/consts"

type JournalQuery struct {
	Page
	*Sort
	Keyword     *string             `json:"keyword" form:"keyword"`
	JournalType *consts.JournalType `json:"journalType" form:"journalType"`
}

type Journal struct {
	SourceContent string             `json:"sourceContent" form:"sourceContent" binding:"gte=1"`
	Content       string             `json:"content" form:"content"`
	Type          consts.JournalType `json:"type" form:"type"`
}
