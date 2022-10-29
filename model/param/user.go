package param

type User struct {
	Username    string `json:"username" binding:"required,lte=50"`
	Nickname    string `json:"nickname" binding:"required,lte=255"`
	Email       string `json:"email" binding:"required,email,lte=127"`
	Password    string `json:"password"`
	Avatar      string `json:"avatar" binding:"lte=1023"`
	Description string `json:"description" binding:"lte=1023"`
}
