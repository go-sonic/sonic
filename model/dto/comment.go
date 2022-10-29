package dto

import "github.com/go-sonic/sonic/consts"

type Comment struct {
	ID                int64                `json:"id"`
	Author            string               `json:"author"`
	Email             string               `json:"email"`
	IpAddress         string               `json:"ipAddress"`
	AuthorURL         string               `json:"authorUrl"`
	GravatarMD5       string               `json:"gravatarMd5"`
	Content           string               `json:"content"`
	Status            consts.CommentStatus `json:"status"`
	UserAgent         string               `json:"userAgent"`
	ParentID          int64                `json:"parentId"`
	IsAdmin           bool                 `json:"isAdmin"`
	AllowNotification bool                 `json:"allowNotification"`
	CreateTime        int64                `json:"createTime"`
	Avatar            string               `json:"avatar"`
}
