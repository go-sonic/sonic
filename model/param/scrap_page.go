package param

type ScrapPage struct {
	Title     string  `json:"title" form:"title"`
	Summary   *string `json:"summary" form:"summary"`
	URL       string  `json:"url" form:"url"`
	OriginURL string  `json:"origin_url" form:"origin_url"`
	AddAt     *int64  `json:"add_at" form:"add_at"`
	Md5       string  `json:"md_5" form:"md_5"`
	Content   *string `json:"content" form:"content"`
	Resource  *string `json:"resource" form:"resource"`
}

type ScrapPageQuery struct {
	Page
	*Sort
	KeyWord string `json:"key_word" form:"key_word""`
}
