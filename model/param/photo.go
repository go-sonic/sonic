package param

type Photo struct {
	Name        string `json:"name" form:"name" binding:"gte=1"`
	Thumbnail   string `json:"thumbnail" form:"thumbnail" binding:"gte=1"`
	TakeTime    *int64 `json:"takeTime" form:"takeTime"`
	URL         string `json:"url" form:"url" binding:"gte=1"`
	Team        string `json:"team" form:"team"`
	Location    string `json:"location" form:"location"`
	Description string `json:"description" form:"description"`
}
