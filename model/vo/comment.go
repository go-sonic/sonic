package vo

import "github.com/go-sonic/sonic/model/dto"

type Comment struct {
	*dto.Comment
	Children []*Comment
}

type CommentWithParent struct {
	*dto.Comment
	Parent *dto.Comment
}

type PostCommentWithPost struct {
	dto.Comment
	Post *dto.PostMinimal `json:"post"`
}

type SheetCommentWithSheet struct {
	dto.Comment
	*dto.PostMinimal
}
type JournalCommentWithJournal struct {
	dto.Comment
	Journal *dto.Journal `json:"journal"`
}

type CommentWithHasChildren struct {
	*dto.Comment
	HasChildren   bool  `json:"hasChildren"`
	ChildrenCount int64 `json:"childrenCount"`
}
