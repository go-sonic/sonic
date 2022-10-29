package dto

type CategoryDTO struct {
	ID          int32  `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	Thumbnail   string `json:"thumbnail"`
	ParentID    int32  `json:"parentId"`
	Password    string `json:"password"`
	CreateTime  int64  `json:"createTime"`
	FullPath    string `json:"fullPath"`
	Priority    int32  `json:"priority"`
	Type        int32  `json:"type"`
}

type CategoryWithPostCount struct {
	*CategoryDTO
	PostCount int32 `json:"postCount"`
}
