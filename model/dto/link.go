package dto

type Link struct {
	ID          int32  `json:"id"`
	Name        string `json:"name"`
	URL         string `json:"url"`
	Logo        string `json:"logo"`
	Description string `json:"description"`
	Team        string `json:"team"`
	Priority    int32  `json:"priority"`
}
