package param

type Meta struct {
	Key   string `json:"key" form:"key" binding:"required"`
	Value string `json:"value" form:"value" binding:"required"`
}
