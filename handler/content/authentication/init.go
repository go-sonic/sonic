package authentication

import "github.com/go-sonic/sonic/injection"

func init() {
	injection.Provide(
		NewCategoryAuthentication,
		NewPostAuthentication,
	)
}
