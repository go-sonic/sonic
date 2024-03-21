package admin

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/go-sonic/sonic/model/dto"
	"github.com/go-sonic/sonic/model/param"
	"github.com/go-sonic/sonic/service"
	"github.com/go-sonic/sonic/util"
	"github.com/go-sonic/sonic/util/xerr"
)

type AttachmentHandler struct {
	AttachmentService service.AttachmentService
}

func NewAttachmentHandler(attachmentService service.AttachmentService) *AttachmentHandler {
	return &AttachmentHandler{
		AttachmentService: attachmentService,
	}
}

func (a *AttachmentHandler) QueryAttachment(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	queryParamNoEnum := param.AttachmentQueryNoEnum{}
	err := ctx.BindAndValidate(&queryParamNoEnum)
	if err != nil {
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("param error ")
	}
	attachmentQuery := param.AssertAttachmentQuery(queryParamNoEnum)
	attachments, totalCount, err := a.AttachmentService.Page(_ctx, &attachmentQuery)
	if err != nil {
		return nil, err
	}
	attachmentDTOs, err := a.AttachmentService.ConvertToDTOs(_ctx, attachments)
	if err != nil {
		return nil, err
	}
	return dto.NewPage(attachmentDTOs, totalCount, attachmentQuery.Page), nil
}

func (a *AttachmentHandler) GetAttachmentByID(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	id, err := util.ParamInt32(_ctx, ctx, "id")
	if err != nil {
		return nil, err
	}
	if id < 0 {
		return nil, xerr.BadParam.New("id < 0").WithStatus(xerr.StatusBadRequest).WithMsg("param error")
	}
	return a.AttachmentService.GetAttachment(_ctx, id)
}

func (a *AttachmentHandler) UploadAttachment(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	fileHeader, err := ctx.FormFile("file")
	if err != nil {
		return nil, xerr.WithMsg(err, "上传文件错误").WithStatus(xerr.StatusBadRequest)
	}
	return a.AttachmentService.Upload(_ctx, fileHeader)
}

func (a *AttachmentHandler) UploadAttachments(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	form, _ := ctx.MultipartForm()
	if len(form.File) == 0 {
		return nil, xerr.BadParam.New("empty files").WithStatus(xerr.StatusBadRequest).WithMsg("empty files")
	}
	files := form.File["files"]
	attachmentDTOs := make([]*dto.AttachmentDTO, 0)
	for _, file := range files {
		attachment, err := a.AttachmentService.Upload(_ctx, file)
		if err != nil {
			return nil, err
		}
		attachmentDTOs = append(attachmentDTOs, attachment)
	}
	return attachmentDTOs, nil
}

func (a *AttachmentHandler) UpdateAttachment(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	id, err := util.ParamInt32(_ctx, ctx, "id")
	if err != nil {
		return nil, err
	}

	updateParam := &param.AttachmentUpdate{}
	err = ctx.BindAndValidate(updateParam)
	if err != nil {
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("param error ")
	}
	return a.AttachmentService.Update(_ctx, id, updateParam)
}

func (a *AttachmentHandler) DeleteAttachment(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	id, err := util.ParamInt32(_ctx, ctx, "id")
	if err != nil {
		return nil, err
	}
	return a.AttachmentService.Delete(_ctx, id)
}

func (a *AttachmentHandler) DeleteAttachmentInBatch(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	ids := make([]int32, 0)
	err := ctx.BindAndValidate(&ids)
	if err != nil {
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("parameter error")
	}
	return a.AttachmentService.DeleteBatch(_ctx, ids)
}

func (a *AttachmentHandler) GetAllMediaType(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	return a.AttachmentService.GetAllMediaTypes(_ctx)
}

func (a *AttachmentHandler) GetAllTypes(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	attachmentTypes, err := a.AttachmentService.GetAllTypes(_ctx)
	if err != nil {
		return nil, err
	}
	return attachmentTypes, nil
}
