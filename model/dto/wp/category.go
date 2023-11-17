package wp

type CategoryDTO struct {
	ID          int32                  `json:"id"`
	Count       int32                  `json:"count"`
	Description string                 `json:"description"`
	Link        string                 `json:"link"`
	Name        string                 `json:"name"`
	Slug        string                 `json:"slug"`
	Taxonomy    string                 `json:"taxonomy"`
	Parent      int32                  `json:"parent"`
	Meta        map[string]interface{} `json:"meta"`
}
