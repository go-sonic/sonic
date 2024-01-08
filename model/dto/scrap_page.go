package dto

type ScrapPageDTO struct {
	ID           int32   `json:"id"`
	MD5          string  `json:"md_5"`
	Content      string  `json:"content"`
	Title        string  `json:"title"`
	Summary      string  `json:"summary"`
	URL          string  `json:"url"`
	OriginURL    *string `json:"origin_url"`
	Domain       string  `json:"domain"`
	Resource     *string `json:"resource"`
	AttachmentID int32   `json:"attachment_id"`
	FileKey      string  `json:"file_key"`
	Path         string  `json:"path"`
	Size         int64   `json:"size"`
	Suffix       string  `json:"suffix"`
	ThumbPath    string  `json:"thumb_path"`
}
