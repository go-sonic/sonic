package content

import (
	"github.com/gin-gonic/gin"

	"github.com/go-sonic/sonic/handler/content/model"
	"github.com/go-sonic/sonic/service"
	"github.com/go-sonic/sonic/template"
	"github.com/go-sonic/sonic/util"
)

type PhotoHandler struct {
	OptionService service.OptionService
	PhotoService  service.PhotoService
	PhotoModel    *model.PhotoModel
}

func NewPhotoHandler(
	optionService service.OptionService,
	photoService service.PhotoService,
	photoModel *model.PhotoModel,

) *PhotoHandler {
	return &PhotoHandler{
		OptionService: optionService,
		PhotoService:  photoService,
		PhotoModel:    photoModel,
	}
}

func (p *PhotoHandler) PhotosPage(ctx *gin.Context, model template.Model) (string, error) {
	page, err := util.ParamInt32(ctx, "page")
	if err != nil {
		return "", err
	}
	return p.PhotoModel.Photos(ctx, model, int(page-1))
}

func (p *PhotoHandler) Phtotos(ctx *gin.Context, model template.Model) (string, error) {
	return p.PhotoModel.Photos(ctx, model, 0)
}
