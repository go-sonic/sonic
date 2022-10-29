package file_storage_impl

import "github.com/go-sonic/sonic/injection"

func init() {
	injection.Provide(
		NewMinIO,
		NewLocalFileStorage,
		NewAliyun,
	)
}
