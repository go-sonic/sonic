package service

type OneTimeTokenService interface {
	Get(oneTimeToken string) (string, bool)
	Create(value string) string
}
