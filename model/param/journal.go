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

type JournalQueryNoEnum struct {
	Page
	Sort        *Sort   `json:"sort" form:"sort"`
	Keyword     *string `json:"keyword" form:"keyword"`
	JournalType string  `json:"journalType" form:"journalType"`
}

func AssertJournalQuery(t JournalQueryNoEnum) JournalQuery {
	str := t.JournalType
	res := JournalQuery{
		Page:    t.Page,
		Sort:    t.Sort,
		Keyword: t.Keyword,
	}
	switch str {
	case `"PUBLIC"`:
		*res.JournalType = consts.JournalTypePublic
	case `"INTIMATE"`:
		*res.JournalType = consts.JournalTypeIntimate
	}
	return res
}
