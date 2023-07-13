package theme

import (
	"go.uber.org/fx"

	"github.com/go-sonic/sonic/injection"
)

func init() {
	injection.Provide(
		NewFileScanner,
		NewPropertyScanner,
		fx.Annotated{Target: NewMultipartZipThemeFetcher, Name: "multipartZipThemeFetcher"},
		fx.Annotated{Target: NewGitThemeFetcher, Name: "gitRepoThemeFetcher"},
	)
}
