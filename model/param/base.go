package param

type Page struct {
	PageNum  int `json:"page" form:"page"`
	PageSize int `json:"size" form:"size"`
}

type Sort struct {
	Fields []string `json:"sort" form:"sort"`
}
