package impl

import (
	"context"
	"errors"
	"mime/multipart"
	"os"
	"strings"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/go-sonic/sonic/consts"
	"github.com/go-sonic/sonic/dal"
	"github.com/go-sonic/sonic/log"
	"github.com/go-sonic/sonic/model/dto"
	"github.com/go-sonic/sonic/model/entity"
	"github.com/go-sonic/sonic/model/param"
	"github.com/go-sonic/sonic/service"
	"github.com/go-sonic/sonic/service/storage"
	"github.com/go-sonic/sonic/util/xerr"
)

type attachmentServiceImpl struct {
	OptionService        service.OptionService
	FileStorageComposite storage.FileStorageComposite
}

func (a *attachmentServiceImpl) ConvertToDTOs(ctx context.Context, attachments []*entity.Attachment) ([]*dto.AttachmentDTO, error) {
	dtos := make([]*dto.AttachmentDTO, 0, len(attachments))
	for _, attachment := range attachments {
		dto := &dto.AttachmentDTO{
			ID:             attachment.ID,
			Name:           attachment.Name,
			Path:           attachment.Path,
			FileKey:        attachment.FileKey,
			ThumbPath:      attachment.ThumbPath,
			MediaType:      attachment.MediaType,
			Suffix:         attachment.Suffix,
			Width:          attachment.Width,
			Height:         attachment.Height,
			Size:           attachment.Size,
			AttachmentType: attachment.Type,
		}
		dtos = append(dtos, dto)
		path, err := a.FileStorageComposite.GetFileStorage(attachment.Type).GetFilePath(ctx, attachment.Path)
		if err != nil {
			log.CtxError(ctx, "GetFilePath err", zap.Error(err))
		}
		thumbPath, err := a.FileStorageComposite.GetFileStorage(attachment.Type).GetFilePath(ctx, attachment.ThumbPath)
		if err != nil {
			log.CtxError(ctx, "GetFilePath err", zap.Error(err))
		}
		dto.Path = path
		dto.ThumbPath = thumbPath
	}
	return dtos, nil
}

func NewAttachmentService(optionService service.OptionService, fileStorageComposite storage.FileStorageComposite) service.AttachmentService {
	return &attachmentServiceImpl{
		FileStorageComposite: fileStorageComposite,
		OptionService:        optionService,
	}
}

func (a *attachmentServiceImpl) ConvertToDTO(ctx context.Context, attachment *entity.Attachment) (*dto.AttachmentDTO, error) {
	dto := &dto.AttachmentDTO{
		ID:             attachment.ID,
		Name:           attachment.Name,
		Path:           attachment.Path,
		FileKey:        attachment.FileKey,
		ThumbPath:      attachment.ThumbPath,
		MediaType:      attachment.MediaType,
		Suffix:         attachment.Suffix,
		Width:          attachment.Width,
		Height:         attachment.Height,
		Size:           attachment.Size,
		AttachmentType: attachment.Type,
	}
	path, err := a.FileStorageComposite.GetFileStorage(attachment.Type).GetFilePath(ctx, attachment.Path)
	if err != nil {
		return nil, err
	}
	thumbPath, err := a.FileStorageComposite.GetFileStorage(attachment.Type).GetFilePath(ctx, attachment.ThumbPath)
	if err != nil {
		return nil, err
	}
	dto.Path = path
	dto.ThumbPath = thumbPath
	return dto, nil
}

func (a *attachmentServiceImpl) Page(ctx context.Context, queryParam *param.AttachmentQuery) ([]*entity.Attachment, int64, error) {
	attachmentDAL := dal.GetQueryByCtx(ctx).Attachment
	attachmentDo := attachmentDAL.WithContext(ctx)
	if queryParam.Keyword != "" {
		attachmentDo = attachmentDo.Where(attachmentDAL.Name.Like(queryParam.Keyword))
	}
	if queryParam.PageNum >= 0 && queryParam.PageSize >= 0 {
		if queryParam.PageSize > 50 {
			queryParam.PageSize = 50
		}
	} else {
		queryParam.PageSize = 10
		queryParam.PageNum = 0
	}
	if queryParam.Keyword != "" {
		attachmentDo = attachmentDo.Where(attachmentDAL.Name.Like(queryParam.Keyword))
	}
	if queryParam.AttachmentType != nil {
		attachmentDo = attachmentDo.Where(attachmentDAL.Type.Eq(*queryParam.AttachmentType))
	}
	if queryParam.MediaType != "" {
		attachmentDo = attachmentDo.Where(attachmentDAL.MediaType.Eq(queryParam.MediaType))
	}
	attachments, totalCount, err := attachmentDo.FindByPage(queryParam.PageNum*queryParam.PageSize, queryParam.PageSize)
	if err != nil {
		return nil, 0, WrapDBErr(err)
	}
	return attachments, totalCount, nil
}

func (a *attachmentServiceImpl) GetAttachment(ctx context.Context, attachmentID int32) (*entity.Attachment, error) {
	attachmentDAL := dal.GetQueryByCtx(ctx).Attachment
	attachment, err := attachmentDAL.WithContext(ctx).Where(attachmentDAL.ID.Eq(attachmentID)).Take()
	if err != nil {
		return nil, WrapDBErr(err)
	}
	return attachment, nil
}

func (a *attachmentServiceImpl) Upload(ctx context.Context, fileHeader *multipart.FileHeader) (attachmentDTO *dto.AttachmentDTO, err error) {
	attachmentType := a.OptionService.GetAttachmentType(ctx)

	fileStorage := a.FileStorageComposite.GetFileStorage(attachmentType)
	attachmentDAL := dal.GetQueryByCtx(ctx).Attachment

	attachmentDTO, err = fileStorage.Upload(ctx, fileHeader)
	if err != nil {
		return nil, err
	}

	record, err := attachmentDAL.WithContext(ctx).Where(attachmentDAL.Path.Eq(attachmentDTO.Path)).Take()
	if record != nil && err == nil {
		return nil, xerr.BadParam.New("附件路径为 " + attachmentDTO.Path + " 已经存在").
			WithStatus(xerr.StatusBadRequest).
			WithMsg("附件路径为 " + attachmentDTO.Path + " 已经存在")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, WrapDBErr(err)
	}
	attachmentEntity := &entity.Attachment{
		FileKey:   attachmentDTO.FileKey,
		Height:    attachmentDTO.Height,
		MediaType: attachmentDTO.MediaType,
		Name:      attachmentDTO.Name,
		Path:      strings.ReplaceAll(attachmentDTO.Path, string(os.PathSeparator), "/"),
		Size:      attachmentDTO.Size,
		Suffix:    attachmentDTO.Suffix,
		ThumbPath: strings.ReplaceAll(attachmentDTO.Path, string(os.PathSeparator), "/"),
		Type:      attachmentDTO.AttachmentType,
		Width:     attachmentDTO.Width,
	}

	// change file path separator to url separator
	if err = attachmentDAL.WithContext(ctx).Create(attachmentEntity); err != nil {
		return nil, WrapDBErr(err)
	}
	attachmentDTO.ID = attachmentEntity.ID
	attachmentDTO.Path, err = fileStorage.GetFilePath(ctx, attachmentEntity.Path)
	if err != nil {
		return nil, err
	}
	attachmentDTO.ThumbPath, err = fileStorage.GetFilePath(ctx, attachmentEntity.ThumbPath)
	if err != nil {
		return nil, err
	}

	return attachmentDTO, nil
}

func (a *attachmentServiceImpl) Delete(ctx context.Context, attachmentID int32) (*entity.Attachment, error) {
	attachmentDAL := dal.GetQueryByCtx(ctx).Attachment
	attachment, err := attachmentDAL.WithContext(ctx).Where(attachmentDAL.ID.Eq(attachmentID)).Take()
	if err != nil {
		return nil, WrapDBErr(err)
	}
	result, err := attachmentDAL.WithContext(ctx).Where(attachmentDAL.ID.Eq(attachmentID)).Delete()
	if err != nil || result.RowsAffected != 1 {
		return nil, xerr.WithMsg(err, "delete file failed")
	}
	fileStorage := a.FileStorageComposite.GetFileStorage(attachment.Type)
	err = fileStorage.Delete(ctx, attachment.FileKey)
	if err != nil {
		return nil, xerr.WithMsg(err, "delete file failed")
	}
	return attachment, nil
}

func (a *attachmentServiceImpl) DeleteBatch(ctx context.Context, ids []int32) (attachments []*entity.Attachment, err error) {
	attachments = make([]*entity.Attachment, 0)
	var globalErr error
	for _, id := range ids {
		attachment, err := a.Delete(ctx, id)
		if err != nil {
			globalErr = err
		} else {
			attachments = append(attachments, attachment)
		}
	}
	if globalErr != nil {
		return attachments, xerr.NoType.Wrap(err).WithMsg("Failed to delete some files")
	}
	return attachments, nil
}

func (a *attachmentServiceImpl) Update(ctx context.Context, id int32, updateParam *param.AttachmentUpdate) (*entity.Attachment, error) {
	attachmentDAL := dal.GetQueryByCtx(ctx).Attachment
	attachment, err := attachmentDAL.WithContext(ctx).Where(attachmentDAL.ID.Eq(id)).Take()
	if err != nil {
		return nil, WrapDBErr(err)
	}
	result, err := attachmentDAL.WithContext(ctx).Where(attachmentDAL.ID.Eq(id)).Updates(entity.Attachment{Name: updateParam.Name})
	if err != nil {
		return nil, WrapDBErr(err)
	}
	if result.RowsAffected != 1 {
		return nil, xerr.NoType.New("update failed")
	}
	attachment.Name = updateParam.Name
	return attachment, nil
}

func (a *attachmentServiceImpl) GetAllMediaTypes(ctx context.Context) ([]string, error) {
	attachmentDAL := dal.GetQueryByCtx(ctx).Attachment
	var allMediaTypes []string
	err := attachmentDAL.WithContext(ctx).Distinct(attachmentDAL.MediaType).Select(attachmentDAL.MediaType).Scan(&allMediaTypes)
	return allMediaTypes, WrapDBErr(err)
}

func (a *attachmentServiceImpl) GetAllTypes(ctx context.Context) ([]consts.AttachmentType, error) {
	attachmentDAL := dal.GetQueryByCtx(ctx).Attachment
	var allTypes []consts.AttachmentType
	err := attachmentDAL.WithContext(ctx).Distinct(attachmentDAL.Type).Select(attachmentDAL.Type).Scan(&allTypes)
	return allTypes, err
}
