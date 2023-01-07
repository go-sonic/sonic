package param

type Install struct {
	User
	Locale string `json:"locale"`
	Title  string `json:"title" binding:"required"`
	URL    string `json:"url"`
}
