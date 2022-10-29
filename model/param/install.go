package param

type Install struct {
	User
	Locale string `json:"locale"`
	Title  string `json:"title" binding:"required"`
	Url    string `json:"url"`
}
