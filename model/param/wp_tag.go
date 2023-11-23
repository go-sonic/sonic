package param

type TagListParam struct {
	Search string `form:"search" json:"search"`
}

type TagCreateParam struct {
	Description string                 `json:"description"`
	Name        string                 `json:"name"`
	Slug        string                 `json:"slug"`
	Meta        map[string]interface{} `json:"meta"`
}
