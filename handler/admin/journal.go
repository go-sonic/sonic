package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/go-sonic/sonic/handler/binding"
	"github.com/go-sonic/sonic/handler/trans"
	"github.com/go-sonic/sonic/model/dto"
	"github.com/go-sonic/sonic/model/param"
	"github.com/go-sonic/sonic/service"
	"github.com/go-sonic/sonic/util"
	"github.com/go-sonic/sonic/util/xerr"
)

type JournalHandler struct {
	JournalService service.JournalService
}

func NewJournalHandler(journalService service.JournalService) *JournalHandler {
	return &JournalHandler{
		JournalService: journalService,
	}
}

func (j *JournalHandler) ListJournal(ctx *gin.Context) (interface{}, error) {
	var journalQuery param.JournalQuery
	err := ctx.ShouldBindWith(&journalQuery, binding.CustomFormBinding)
	if err != nil {
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("Parameter error")
	}
	journalQuery.Sort = &param.Sort{
		Fields: []string{"createTime,desc"},
	}
	journals, totalCount, err := j.JournalService.ListJournal(ctx, journalQuery)
	if err != nil {
		return nil, err
	}
	journalDTOs, err := j.JournalService.ConvertToWithCommentDTOList(ctx, journals)
	if err != nil {
		return nil, err
	}
	return dto.NewPage(journalDTOs, totalCount, journalQuery.Page), nil
}

func (j *JournalHandler) ListLatestJournal(ctx *gin.Context) (interface{}, error) {
	top, err := util.MustGetQueryInt(ctx, "top")
	if err != nil {
		top = 10
	}
	journalQuery := param.JournalQuery{
		Sort: &param.Sort{Fields: []string{"createTime,desc"}},
		Page: param.Page{PageNum: 0, PageSize: top},
	}
	journals, _, err := j.JournalService.ListJournal(ctx, journalQuery)
	if err != nil {
		return nil, err
	}
	return j.JournalService.ConvertToWithCommentDTOList(ctx, journals)
}

func (j *JournalHandler) CreateJournal(ctx *gin.Context) (interface{}, error) {
	var journalParam param.Journal
	err := ctx.ShouldBindJSON(&journalParam)
	if err != nil {
		if e, ok := err.(validator.ValidationErrors); ok {
			return nil, xerr.WithStatus(e, xerr.StatusBadRequest).WithMsg(trans.Translate(e))
		}
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("parameter error")
	}
	journal, err := j.JournalService.Create(ctx, &journalParam)
	if err != nil {
		return nil, err
	}
	return j.JournalService.ConvertToDTO(journal), nil
}

func (j *JournalHandler) UpdateJournal(ctx *gin.Context) (interface{}, error) {
	var journalParam param.Journal
	err := ctx.ShouldBindJSON(&journalParam)
	if err != nil {
		if e, ok := err.(validator.ValidationErrors); ok {
			return nil, xerr.WithStatus(e, xerr.StatusBadRequest).WithMsg(trans.Translate(e))
		}
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("parameter error")
	}

	journalID, err := util.ParamInt32(ctx, "journalID")
	if err != nil {
		return nil, err
	}
	return j.JournalService.Update(ctx, journalID, &journalParam)
}

func (j *JournalHandler) DeleteJournal(ctx *gin.Context) (interface{}, error) {
	journalID, err := util.ParamInt32(ctx, "journalID")
	if err != nil {
		return nil, err
	}
	return nil, j.JournalService.Delete(ctx, int32(journalID))
}
