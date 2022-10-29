package file_storage

import (
	"context"
	"mime/multipart"

	"github.com/go-sonic/sonic/consts"
	"github.com/go-sonic/sonic/model/dto"
	"github.com/go-sonic/sonic/service/file_storage/file_storage_impl"
)

type FileStorage interface {
	Upload(ctx context.Context, file *multipart.FileHeader) (*dto.AttachmentDTO, error)
	Delete(ctx context.Context, fileKey string) error
	GetAttachmentType() consts.AttachmentType
	GetFilePath(ctx context.Context, relativePath string) (string, error)
}

type FileStorageComposite interface {
	GetFileStorage(storageType consts.AttachmentType) FileStorage
}
type fileStorageComposite struct {
	localStorage *file_storage_impl.LocalFileStorage
	minio        *file_storage_impl.MinIO
	aliyunOSS    *file_storage_impl.Aliyun
}

func NewFileStorageComposite(localStorage *file_storage_impl.LocalFileStorage, minio *file_storage_impl.MinIO, aliyun *file_storage_impl.Aliyun) FileStorageComposite {
	return &fileStorageComposite{
		localStorage: localStorage,
		minio:        minio,
		aliyunOSS:    aliyun,
	}
}

func (f *fileStorageComposite) GetFileStorage(storageType consts.AttachmentType) FileStorage {
	switch storageType {
	case consts.AttachmentTypeLocal:
		return f.localStorage
	case consts.AttachmentTypeMinIO:
		return f.minio
	case consts.AttachmentTypeAliOSS:
		return f.aliyunOSS
	default:
		panic("Unsupported file storage")
	}
}
