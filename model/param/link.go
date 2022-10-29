package param

type Link struct {
	Name        string `json:"name" form:"name" binding:"gte=1,lte=255"`
	URL         string `json:"url" form:"url" binding:"url,lte=255"`
	Logo        string `json:"logo" form:"logo" binding:"lte=1023"`
	Description string `json:"description" form:"description" binding:"lte=255"`
	Team        string `json:"team" form:"team" binding:"lte=255"`
	Priority    int32  `json:"priority" form:"priority" binding:"gte=0"`
}
