package param

type ResetPasswordParam struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Code     string `json:"code"`
	Password string `json:"Password"`
}
