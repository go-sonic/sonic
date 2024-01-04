package filestorageimpl

import (
	"context"
	"image"
	"io"
	"mime/multipart"
	"net/url"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/go-sonic/sonic/consts"
	"github.com/go-sonic/sonic/model/dto"
	"github.com/go-sonic/sonic/model/property"
	"github.com/go-sonic/sonic/service"
	"github.com/go-sonic/sonic/util/xerr"
)

type MinIO struct {
	OptionService service.OptionService
}

func NewMinIO(optionService service.OptionService) *MinIO {
	return &MinIO{
		OptionService: optionService,
	}
}

type minioClient struct {
	*minio.Client
	BucketName string
	Source     string
	EndPoint   string
}

func (m *MinIO) Upload(ctx context.Context, fileHeader *multipart.FileHeader) (*dto.AttachmentDTO, error) {
	minioClientInstance, err := m.getMinioClient(ctx)
	if err != nil {
		return nil, err
	}

	fd, err := newURLFileDescriptor(
		withBaseURL(minioClientInstance.EndPoint+"/"+minioClientInstance.BucketName),
		withSubURLPath(minioClientInstance.Source),
		withShouldRenameURLOption(commonRenamePredicateFunc(ctx, consts.AttachmentTypeMinIO)),
		withOriginalNameURLOption(fileHeader.Filename),
	)
	if err != nil {
		return nil, xerr.WithStatus(err, xerr.StatusInternalServerError)
	}
	file, err := fileHeader.Open()
	if err != nil {
		return nil, xerr.WithStatus(err, xerr.StatusInternalServerError).WithMsg("open upload file error")
	}
	defer file.Close()
	_, err = minioClientInstance.PutObject(ctx, minioClientInstance.BucketName, fd.getRelativePath(), file, fileHeader.Size, minio.PutObjectOptions{})
	if err != nil {
		return nil, xerr.WithMsg(err, "upload to minio error").WithStatus(xerr.StatusInternalServerError).WithErrMsgf("err=%v", err)
	}

	mediaType, _ := getFileContentType(file)
	result := &dto.AttachmentDTO{
		Name:           fd.getFileName(),
		Path:           fd.getRelativePath(),
		FileKey:        fd.getRelativePath(),
		Suffix:         fd.getExtensionName(),
		MediaType:      mediaType,
		AttachmentType: consts.AttachmentTypeMinIO,
		Size:           fileHeader.Size,
	}
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return nil, xerr.WithStatus(err, xerr.StatusInternalServerError)
	}
	err = handleImageMeta(file, result, func(srcImage image.Image) (string, error) {
		return fd.getRelativePath(), nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (m *MinIO) Delete(ctx context.Context, fileKey string) error {
	minioClientInstance, err := m.getMinioClient(ctx)
	if err != nil {
		return err
	}
	err = minioClientInstance.RemoveObject(ctx, minioClientInstance.BucketName, fileKey, minio.RemoveObjectOptions{})
	if err != nil {
		return xerr.WithStatus(err, xerr.StatusInternalServerError).WithErrMsgf("err=%v", err)
	}
	return nil
}

func (m *MinIO) GetAttachmentType() consts.AttachmentType {
	return consts.AttachmentTypeMinIO
}

func (m *MinIO) GetFilePath(ctx context.Context, relativePath string) (string, error) {
	minioClientInstance, err := m.getMinioClient(ctx)
	if err != nil {
		return "", err
	}
	base := minioClientInstance.EndPoint + "/" + minioClientInstance.BucketName
	fullPath, _ := url.JoinPath(base, relativePath)
	fullPath, _ = url.PathUnescape(fullPath)
	return fullPath, nil
}

func (m *MinIO) getMinioClient(ctx context.Context) (*minioClient, error) {
	getClientProperty := func(propertyValue *string, property property.Property, e error) error {
		if e != nil {
			return e
		}
		value, err := m.OptionService.GetOrByDefaultWithErr(ctx, property, property.DefaultValue)
		if err != nil {
			return err
		}
		strValue, ok := value.(string)
		if !ok {
			return xerr.WithStatus(nil, xerr.StatusBadRequest).WithErrMsgf("wrong property type")
		}
		if strValue == "" {
			return xerr.WithStatus(nil, xerr.StatusInternalServerError).WithMsg("property not found: " + property.KeyValue)
		}
		*propertyValue = strValue
		return nil
	}
	var endPoint, bucketName, accessKey, accessSecret, source, region string
	err := getClientProperty(&endPoint, property.MinioEndpoint, nil)
	err = getClientProperty(&bucketName, property.MinioBucketName, err)
	err = getClientProperty(&accessKey, property.MinioAccessKey, err)
	err = getClientProperty(&accessSecret, property.MinioAccessSecret, err)
	err = getClientProperty(&source, property.MinioSource, err)
	err = getClientProperty(&region, property.MinioRegion, err)
	if err != nil {
		return nil, err
	}

	u, err := url.ParseRequestURI(endPoint)
	if err != nil {
		return nil, err
	}

	useSSL := u.Scheme == "https"
	// u.Host like 127.0.0.1:9000 or hostname
	client, err := minio.New(u.Host, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, accessSecret, ""),
		Secure: useSSL,
		Region: region,
	})
	if err != nil {
		return nil, xerr.WithStatus(err, xerr.StatusInternalServerError).WithMsg("failed to initialize minio: " + err.Error())
	}

	minioClientInstance := &minioClient{}

	minioClientInstance.Client = client
	minioClientInstance.BucketName = bucketName
	minioClientInstance.Source = source
	minioClientInstance.EndPoint = endPoint
	return minioClientInstance, nil
}
