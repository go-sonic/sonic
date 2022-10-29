package dto

type Photo struct {
	ID          int32  `json:"id"`
	Name        string `json:"name"`
	Thumbnail   string `json:"thumbnail"`
	TakeTime    int64  `json:"takeTime"`
	URL         string `json:"url"`
	Team        string `json:"team"`
	Location    string `json:"location"`
	Description string `json:"description"`
	Likes       int64  `json:"likes"`
}
