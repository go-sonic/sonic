package file_storage_impl

import (
	"context"
	"image"
	"io"
	"math"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/go-sonic/sonic/consts"
	"github.com/go-sonic/sonic/dal"
	"github.com/go-sonic/sonic/model/dto"
	"github.com/go-sonic/sonic/util"
	"github.com/go-sonic/sonic/util/xerr"
)

func commonRenamePredicateFunc(ctx context.Context, attachmentType consts.AttachmentType) func(relativePath string) (bool, error) {
	return func(relativePath string) (bool, error) {
		attachmentDAL := dal.GetQueryByCtx(ctx).Attachment
		count, err := attachmentDAL.WithContext(ctx).Where(attachmentDAL.Type.Eq(attachmentType), attachmentDAL.FileKey.Eq(relativePath)).Count()
		if err != nil {
			return false, err
		}
		return count > 0, nil
	}
}

func getFileContentType(file multipart.File) (string, error) {
	_, err := file.Seek(0, io.SeekStart)
	if err != nil {
		return "", xerr.WithMsg(err, "seek file error").WithStatus(xerr.StatusInternalServerError)
	}
	buf := make([]byte, 512)
	_, err = file.Read(buf)
	if err != nil {
		return "", err
	}

	return http.DetectContentType(buf), nil
}

func isImageType(mediaType string) bool {
	return strings.Contains(mediaType, "image")
}

func handleImageMeta(file multipart.File, attachment *dto.AttachmentDTO, thumbnailFn func(srcImage image.Image) (string, error)) error {
	if attachment == nil {
		return nil
	}
	if !isImageType(attachment.MediaType) {
		return nil
	}
	srcImage, _, err := image.Decode(file)
	if err != nil {
		return xerr.NoType.Wrap(err).WithMsg("Handle srcImage error")
	}
	bounds := srcImage.Bounds()

	srcWidth := bounds.Dx()
	srcHeight := bounds.Dy()
	attachment.Width = int32(srcWidth)
	attachment.Height = int32(srcHeight)

	if thumbnailFn == nil {
		return nil
	}

	thumbnailPath, err := thumbnailFn(srcImage)
	if err != nil {
		return err
	}
	attachment.ThumbPath = thumbnailPath
	return nil
}

func keepAspectRatio(sourceWidth, sourceHeight, targetWidth, targetHeight int) (width, height int) {
	sourceRatio := float64(sourceWidth) / float64(sourceHeight)
	targetRatio := float64(targetWidth) / float64(targetHeight)
	if sourceRatio > targetRatio {
		width = targetWidth
		height = int(math.Floor(float64(targetWidth)/sourceRatio + 0.5))
	} else {
		width = int(math.Floor(float64(targetHeight)*sourceRatio + 0.5))
		height = targetHeight
	}
	width = util.IfElse(width == 0, 1, width).(int)
	height = util.IfElse(height == 0, 1, height).(int)
	return
}
