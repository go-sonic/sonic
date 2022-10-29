package param

type Tag struct {
	Name      string `json:"name" form:"name" binding:"gte=1,lte=255"`
	Slug      string `json:"slug" form:"slug" binding:"lte=255"`
	Thumbnail string `json:"thumbnail" form:"thumbnail" binding:"lte=1023"`
	Color     string `json:"color" form:"color" biding:"lte=24"`
}
