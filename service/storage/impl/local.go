package filestorageimpl

import (
	"context"
	"image"
	"io"
	"mime/multipart"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"time"

	"github.com/disintegration/imaging"
	"go.uber.org/zap"
	_ "golang.org/x/image/webp"

	"github.com/go-sonic/sonic/config"
	"github.com/go-sonic/sonic/consts"
	"github.com/go-sonic/sonic/log"
	"github.com/go-sonic/sonic/model/dto"
	"github.com/go-sonic/sonic/service"
	"github.com/go-sonic/sonic/util"
	"github.com/go-sonic/sonic/util/xerr"
)

var (
	thumbWidth      = 256
	thumbHeight     = 256
	thumbnailSuffix = "-thumbnail"
)

type LocalFileStorage struct {
	Config        *config.Config
	OptionService service.OptionService
}

func NewLocalFileStorage(config *config.Config, optionService service.OptionService) *LocalFileStorage {
	return &LocalFileStorage{
		Config:        config,
		OptionService: optionService,
	}
}

func (l *LocalFileStorage) Upload(ctx context.Context, fileHeader *multipart.FileHeader) (*dto.AttachmentDTO, error) {
	now := time.Now()
	year, month, _ := now.Date()

	fd, err := newLocalFileDescriptor(
		withOriginalName(fileHeader.Filename),
		withBasePath(l.Config.Sonic.WorkDir),
		withSubPath(l.getUploadPath(year, int(month))),
		withShouldRename(commonRenamePredicateFunc(ctx, consts.AttachmentTypeLocal)),
	)
	if err != nil {
		return nil, err
	}

	log.CtxDebug(ctx, "Local FileStorage ",
		zap.String("originalFilename", fileHeader.Filename),
		zap.String("absPath", fd.getFullPath()))

	err = os.MkdirAll(fd.getFullDirPath(), os.ModePerm)
	if isExist := os.IsExist(err); err != nil && !isExist {
		return nil, xerr.WithMsg(err, "Upload file error")
	}

	srcFile, err := fileHeader.Open()
	if err != nil {
		return nil, xerr.WithMsg(err, "Upload file error")
	}
	defer srcFile.Close()

	out, err := os.Create(fd.getFullPath())
	if os.IsExist(err) {
		return nil, xerr.WithMsg(err, "The file already exists!")
	} else if err != nil {
		return nil, xerr.WithMsg(err, "Upload file error")
	}
	defer out.Close()
	defer func() {
		if err != nil {
			_ = l.Delete(ctx, fd.getFullPath())
		}
	}()

	_, err = io.Copy(out, srcFile)
	if err != nil {
		return nil, xerr.WithMsg(err, "Error writing file")
	}

	mediaType, _ := getFileContentType(srcFile)
	attachment := dto.AttachmentDTO{
		Name:           fd.getFileName(),
		Path:           fd.getRelativePath(),
		FileKey:        fd.getRelativePath(),
		Suffix:         fd.getExtensionName(),
		MediaType:      mediaType,
		AttachmentType: consts.AttachmentTypeLocal,
		Size:           fileHeader.Size,
	}

	thumbnailFn := func(srcImage image.Image) (string, error) {
		if fd.getExtensionName() == "webp" {
			return fd.getRelativePath(), nil
		}
		thumbnailFd, err := newLocalFileDescriptor(
			withOriginalName(fileHeader.Filename),
			withBasePath(l.Config.Sonic.WorkDir),
			withSubPath(l.getUploadPath(year, int(month))),
			withShouldRename(func(r string) (bool, error) {
				return false, nil
			}),
			withSuffix("-thumbnail"),
		)
		if err != nil {
			return "", err
		}
		bounds := srcImage.Bounds()
		srcWidth := bounds.Dx()
		srcHeight := bounds.Dy()
		width, height := keepAspectRatio(srcWidth, srcHeight, thumbWidth, thumbHeight)

		dstImage := imaging.Thumbnail(srcImage, width, height, imaging.Box)
		err = imaging.Save(dstImage, thumbnailFd.getFullPath())
		if err != nil {
			return "", xerr.NoType.Wrap(err).WithMsg("Save thumb srcImage err")
		}
		return thumbnailFd.getRelativePath(), nil
	}
	_, err = srcFile.Seek(0, io.SeekStart)
	if err != nil {
		return nil, xerr.WithStatus(err, xerr.StatusInternalServerError)
	}
	err = handleImageMeta(srcFile, &attachment, thumbnailFn)
	if err != nil {
		return nil, err
	}
	return &attachment, nil
}

func (l *LocalFileStorage) Delete(ctx context.Context, fileKey string) error {
	filePath, err := filepath.Abs(fileKey)
	if err != nil {
		return xerr.NoType.Wrap(err)
	}
	extName := filepath.Ext(fileKey)
	fullName := path.Base(fileKey)
	originalFilename := fullName[0 : len(fullName)-len(extName)]

	thumbFileName := originalFilename + thumbnailSuffix + extName
	thumbFilePath := filepath.Join(filepath.Dir(fileKey), thumbFileName)

	err = os.Remove(filePath)
	if err != nil && !os.IsNotExist(err) {
		return xerr.NoType.Wrap(err).WithMsg("delete file failed")
	}
	err = os.Remove(thumbFilePath)
	if err != nil && !os.IsNotExist(err) {
		return xerr.NoType.Wrap(err).WithMsg("delete file failed")
	}
	return nil
}

func (l *LocalFileStorage) GetAttachmentType() consts.AttachmentType {
	return consts.AttachmentTypeLocal
}

func (l *LocalFileStorage) GetFilePath(ctx context.Context, relativePath string) (string, error) {
	isEnabled, _ := l.OptionService.IsEnabledAbsolutePath(ctx)
	var blogBaseURL string
	if isEnabled {
		blogBaseURL, _ = l.OptionService.GetBlogBaseURL(ctx)
	}
	fullPath, _ := url.JoinPath(blogBaseURL, relativePath)
	if blogBaseURL == "" {
		fullPath, _ = url.JoinPath("/", relativePath)
	}
	fullPath, _ = url.PathUnescape(fullPath)
	return fullPath, nil
}

func (l *LocalFileStorage) getUploadPath(year, month int) string {
	uploadPath := filepath.Join(
		consts.SonicUploadDir,
		strconv.Itoa(year),
		util.IfElse(month >= 10, strconv.Itoa(month), "0"+strconv.Itoa(month)).(string),
	)
	return filepath.Clean(uploadPath)
}
