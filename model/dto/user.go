package dto

import "github.com/go-sonic/sonic/consts"

type User struct {
	ID          int32          `json:"id"`
	Username    string         `json:"username"`
	Nickname    string         `json:"nickname"`
	Email       string         `json:"email"`
	Avatar      string         `json:"avatar"`
	Description string         `json:"description"`
	MFAType     consts.MFAType `json:"mfaType"`
	CreateTime  int64          `json:"createTime"`
	UpdateTime  int64          `json:"updateTime"`
}
