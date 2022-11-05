package param

import "github.com/go-sonic/sonic/consts"

type CommentQuery struct {
	Page
	*Sort
	ContentID     *int32
	Keyword       *string               `json:"keyword"`
	CommentStatus *consts.CommentStatus `json:"status"`
	ParentID      *int32                `json:"parentID"`
}

type Comment struct {
	Author            string             `json:"author" form:"author" binding:"gte=1,lte=50"`
	Email             string             `json:"email" form:"email" binding:"email,lte=255"`
	AuthorURL         string             `json:"authorUrl" form:"authorUrl" binding:"lte=255"`
	Content           string             `json:"content" form:"content" binding:"gte=1,lte=1023"`
	PostID            int32              `json:"postId" form:"postId" binding:"gte=1"`
	ParentID          int32              `json:"parentId" form:"parentId" binding:"gte=0"`
	AllowNotification bool               `json:"allowNotification" form:"allowNotification"`
	CommentType       consts.CommentType `json:"-"`
}

type AdminComment struct {
	Author            string             `json:"author" form:"author"`
	Email             string             `json:"email" form:"email"`
	AuthorURL         string             `json:"authorUrl" form:"authorUrl"`
	Content           string             `json:"content" form:"content"`
	PostID            int32              `json:"postId" form:"postId"`
	ParentID          int32              `json:"parentId" form:"parentId"`
	AllowNotification bool               `json:"allowNotification"`
	CommentType       consts.CommentType `json:"-"`
}
