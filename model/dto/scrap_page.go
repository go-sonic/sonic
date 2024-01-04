package dto

type ScrapPageDTO struct {
	ID      int32  `json:"id"`
	Content string `json:"content"`
	Title   string `json:"title"`
	Summary string `json:"summary"`
	URL     string `json:"url"`
}
