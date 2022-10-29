package param

type LoginParam struct {
	Username string `json:"username" binding:"gte=1,lte=20"`
	Password string `json:"password" binding:"gte=6"`
	AuthCode string `json:"authcode" `
}
