package param

type ThemeContent struct {
	Path    string `json:"path" form:"path" binding:"gte=1"`
	Content string `json:"content" form:"path"`
}
