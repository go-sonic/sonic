package dto

type Tag struct {
	ID         int32  `json:"id"`
	Name       string `json:"name"`
	Slug       string `json:"slug"`
	Thumbnail  string `json:"thumbnail"`
	CreateTime int64  `json:"createTime"`
	FullPath   string `json:"fullPath"`
	Color      string `json:"color"`
}

type TagWithPostCount struct {
	*Tag
	PostCount int64 `json:"postCount"`
}
