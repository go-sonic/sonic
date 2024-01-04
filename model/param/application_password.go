package param

type ApplicationPasswordParam struct {
	Name string `json:"name" binding:"required"`
}
