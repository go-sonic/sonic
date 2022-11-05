package dto

type Meta struct {
	ID         int32  `json:"id"`
	PostID     int32  `json:"postId"`
	Key        string `json:"key"`
	Value      string `json:"value"`
	CreateTime int64  `json:"createTime"`
}
