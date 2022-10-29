package param

type Category struct {
	ID          int32  `json:"id" binding:"gte=0"`
	Name        string `json:"name" binding:"gte=1,lte=255"`
	Slug        string `json:"slug" binding:"gte=0,lte=255"`
	Description string `json:"description" binding:"gte=0,lte=100"`
	Thumbnail   string `json:"thumbnail" binding:"gte=0,lte=1023"`
	Password    string `json:"password" binding:"gte=0,lte=255"`
	ParentID    int32  `json:"parentId" binding:"gte=0"`
	Priority    int32  `json:"priority" binding:"gte=0"`
}
