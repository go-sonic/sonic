package wp

import "github.com/go-sonic/sonic/injection"

func init() {
	injection.Provide(
		NewPostHandler,
		NewUserHandler,
		NewCategoryHandler,
		NewTagHandler,
	)
}
