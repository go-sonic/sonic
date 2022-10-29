package param

type TestEmail struct {
	To      string `json:"to" form:"to" binding:"email"`
	Subject string `json:"subject" form:"subject" binding:"gte=1"`
	Content string `json:"content" form:"content" binding:"gte=1"`
}
