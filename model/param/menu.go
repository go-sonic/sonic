package param

type Menu struct {
	ID       int32  `json:"id" form:"id"`
	Name     string `json:"name" form:"name" binding:"gte=1,lte=50"`
	URL      string `json:"URL" form:"url" binding:"gte=1,lte=1023"`
	Priority int32  `json:"priority" form:"priority" binding:"gte=0"`
	Target   string `json:"target" form:"target" binding:"lte=50"`
	Icon     string `json:"icon" form:"icon" binding:"lte=50"`
	ParentID int32  `json:"parentId" form:"parentId" binding:"gte=0"`
	Team     string `json:"team" form:"team" binding:"lte=255"`
}
