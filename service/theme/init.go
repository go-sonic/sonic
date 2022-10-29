package theme

import "github.com/go-sonic/sonic/injection"

func init() {
	injection.Provide(
		NewFileScanner,
		NewPropertyScanner,
		NewMultipartZipThemeFetcher,
	)
}
