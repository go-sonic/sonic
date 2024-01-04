package wp

type TagDTO struct {
	ID          int32                  `json:"id"`
	Count       int32                  `json:"count"`
	Description string                 `json:"description"`
	Link        string                 `json:"link"`
	Name        string                 `json:"name"`
	Slug        string                 `json:"slug"`
	Taxonomy    string                 `json:"taxonomy"`
	Meta        map[string]interface{} `json:"meta"`
}
