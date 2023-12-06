package param

type ScrapPage struct {
	Title     string `json:"title"`
	Summary   string `json:"summary"`
	URL       string `json:"url"`
	OriginURL string `json:"origin_url"`
	AddAt     int64  `json:"add_at"`
	Md5       string `json:"md_5"`
	Content   string `json:"content"`
}
