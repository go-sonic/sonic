package admin

import (
	"context"
	"errors"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/go-playground/validator/v10"

	"github.com/go-sonic/sonic/handler/trans"
	"github.com/go-sonic/sonic/model/dto"
	"github.com/go-sonic/sonic/model/param"
	"github.com/go-sonic/sonic/service"
	"github.com/go-sonic/sonic/util"
	"github.com/go-sonic/sonic/util/xerr"
)

type PhotoHandler struct {
	PhotoService service.PhotoService
}

func NewPhotoHandler(photoService service.PhotoService) *PhotoHandler {
	return &PhotoHandler{
		PhotoService: photoService,
	}
}

func (p *PhotoHandler) ListPhoto(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	sort := param.Sort{}
	err := ctx.BindAndValidate(&sort)
	if err != nil {
		return nil, xerr.WithMsg(err, "sort parameter error").WithStatus(xerr.StatusBadRequest)
	}
	if len(sort.Fields) == 0 {
		sort.Fields = append(sort.Fields, "createTime,desc")
	}
	photos, err := p.PhotoService.List(_ctx, &sort)
	if err != nil {
		return nil, err
	}
	return p.PhotoService.ConvertToDTOs(_ctx, photos), nil
}

func (p *PhotoHandler) PagePhotos(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	type Param struct {
		param.Page
		param.Sort
	}
	param := Param{}
	err := ctx.BindAndValidate(&param)
	if err != nil {
		return nil, xerr.WithMsg(err, "parameter error").WithStatus(xerr.StatusBadRequest)
	}
	if len(param.Fields) == 0 {
		param.Fields = append(param.Fields, "createTime,desc")
	}
	photos, totalCount, err := p.PhotoService.Page(_ctx, param.Page, &param.Sort)
	if err != nil {
		return nil, err
	}
	return dto.NewPage(p.PhotoService.ConvertToDTOs(_ctx, photos), totalCount, param.Page), nil
}

func (p *PhotoHandler) GetPhotoByID(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	id, err := util.ParamInt32(_ctx, ctx, "id")
	if err != nil {
		return nil, err
	}
	photo, err := p.PhotoService.GetByID(_ctx, id)
	if err != nil {
		return nil, err
	}
	return p.PhotoService.ConvertToDTO(_ctx, photo), nil
}

func (p *PhotoHandler) CreatePhoto(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	photoParam := &param.Photo{}
	err := ctx.BindAndValidate(photoParam)
	if err != nil {
		e := validator.ValidationErrors{}
		if errors.As(err, &e) {
			return nil, xerr.WithStatus(e, xerr.StatusBadRequest).WithMsg(trans.Translate(e))
		}
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("parameter error")
	}
	photo, err := p.PhotoService.Create(_ctx, photoParam)
	if err != nil {
		return nil, err
	}
	return p.PhotoService.ConvertToDTO(_ctx, photo), nil
}

func (p *PhotoHandler) CreatePhotoBatch(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	photosParam := make([]*param.Photo, 0)
	err := ctx.BindAndValidate(&photosParam)
	if err != nil {
		e := validator.ValidationErrors{}
		if errors.As(err, &e) {
			return nil, xerr.WithStatus(e, xerr.StatusBadRequest).WithMsg(trans.Translate(e))
		}
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("parameter error")
	}
	photos, err := p.PhotoService.CreateBatch(_ctx, photosParam)
	if err != nil {
		return nil, err
	}
	return p.PhotoService.ConvertToDTOs(_ctx, photos), nil
}

func (p *PhotoHandler) UpdatePhoto(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	id, err := util.ParamInt32(_ctx, ctx, "id")
	if err != nil {
		return nil, err
	}
	photoParam := &param.Photo{}
	err = ctx.BindAndValidate(photoParam)
	if err != nil {
		e := validator.ValidationErrors{}
		if errors.As(err, &e) {
			return nil, xerr.WithStatus(e, xerr.StatusBadRequest).WithMsg(trans.Translate(e))
		}
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("parameter error")
	}
	photo, err := p.PhotoService.Update(_ctx, id, photoParam)
	if err != nil {
		return nil, err
	}
	return p.PhotoService.ConvertToDTO(_ctx, photo), nil
}

func (p *PhotoHandler) DeletePhoto(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	id, err := util.ParamInt32(_ctx, ctx, "id")
	if err != nil {
		return nil, err
	}
	return nil, p.PhotoService.Delete(_ctx, id)
}

func (p *PhotoHandler) DeletePhotoBatch(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	photosParam := make([]int32, 0)
	err := ctx.BindAndValidate(&photosParam)
	if err != nil {
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("parameter error")
	}
	for _, id := range photosParam {
		err := p.PhotoService.Delete(_ctx, id)
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func (p *PhotoHandler) ListPhotoTeams(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	return p.PhotoService.ListTeams(_ctx)
}
