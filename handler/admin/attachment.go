package admin

import (
	"github.com/gin-gonic/gin"

	"github.com/go-sonic/sonic/handler/binding"
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

func (a *AttachmentHandler) QueryAttachment(ctx *gin.Context) (interface{}, error) {
	queryParam := &param.AttachmentQuery{}
	err := ctx.ShouldBindWith(queryParam, binding.CustomFormBinding)
	if err != nil {
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("param error ")
	}
	attachments, totalCount, err := a.AttachmentService.Page(ctx, queryParam)
	if err != nil {
		return nil, err
	}
	attachmentDTOs, err := a.AttachmentService.ConvertToDTOs(ctx, attachments)
	if err != nil {
		return nil, err
	}
	return dto.NewPage(attachmentDTOs, totalCount, queryParam.Page), nil
}

func (a *AttachmentHandler) GetAttachmentByID(ctx *gin.Context) (interface{}, error) {
	id, err := util.ParamInt32(ctx, "id")
	if err != nil {
		return nil, err
	}
	if id < 0 {
		return nil, xerr.BadParam.New("id < 0").WithStatus(xerr.StatusBadRequest).WithMsg("param error")
	}
	return a.AttachmentService.GetAttachment(ctx, id)
}

func (a *AttachmentHandler) UploadAttachment(ctx *gin.Context) (interface{}, error) {
	fileHeader, err := ctx.FormFile("file")
	if err != nil {
		return nil, xerr.WithMsg(err, "上传文件错误").WithStatus(xerr.StatusBadRequest)
	}
	return a.AttachmentService.Upload(ctx, fileHeader)
}

func (a *AttachmentHandler) UploadAttachments(ctx *gin.Context) (interface{}, error) {
	form, _ := ctx.MultipartForm()
	if len(form.File) == 0 {
		return nil, xerr.BadParam.New("empty files").WithStatus(xerr.StatusBadRequest).WithMsg("empty files")
	}
	files := form.File["files"]
	attachmentDTOs := make([]*dto.AttachmentDTO, 0)
	for _, file := range files {
		attachment, err := a.AttachmentService.Upload(ctx, file)
		if err != nil {
			return nil, err
		}
		attachmentDTOs = append(attachmentDTOs, attachment)
	}
	return attachmentDTOs, nil
}

func (a *AttachmentHandler) UpdateAttachment(ctx *gin.Context) (interface{}, error) {
	id, err := util.ParamInt32(ctx, "id")
	if err != nil {
		return nil, err
	}

	updateParam := &param.AttachmentUpdate{}
	err = ctx.ShouldBind(updateParam)
	if err != nil {
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("param error ")
	}
	return a.AttachmentService.Update(ctx, id, updateParam)
}

func (a *AttachmentHandler) DeleteAttachment(ctx *gin.Context) (interface{}, error) {
	id, err := util.ParamInt32(ctx, "id")
	if err != nil {
		return nil, err
	}
	return a.AttachmentService.Delete(ctx, id)
}

func (a *AttachmentHandler) DeleteAttachmentInBatch(ctx *gin.Context) (interface{}, error) {
	ids := make([]int32, 0)
	err := ctx.ShouldBind(&ids)
	if err != nil {
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("parameter error")
	}
	return a.AttachmentService.DeleteBatch(ctx, ids)
}

func (a *AttachmentHandler) GetAllMediaType(ctx *gin.Context) (interface{}, error) {
	return a.AttachmentService.GetAllMediaTypes(ctx)
}

func (a *AttachmentHandler) GetAllTypes(ctx *gin.Context) (interface{}, error) {
	attachmentTypes, err := a.AttachmentService.GetAllTypes(ctx)
	if err != nil {
		return nil, err
	}
	return attachmentTypes, nil
}
