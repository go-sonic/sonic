package filestorageimpl

import "github.com/go-sonic/sonic/injection"

func init() {
	injection.Provide(
		NewMinIO,
		NewLocalFileStorage,
		NewAliyun,
	)
}
