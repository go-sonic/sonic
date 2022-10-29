package dto

type AuthTokenDTO struct {
	AccessToken  string `json:"access_token"`
	ExpiredIn    int    `json:"expired_in"`
	RefreshToken string `json:"refresh_token"`
}
