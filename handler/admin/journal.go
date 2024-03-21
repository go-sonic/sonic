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

type JournalHandler struct {
	JournalService service.JournalService
}

func NewJournalHandler(journalService service.JournalService) *JournalHandler {
	return &JournalHandler{
		JournalService: journalService,
	}
}

func (j *JournalHandler) ListJournal(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	var journalQueryNoEnum param.JournalQueryNoEnum
	err := ctx.BindAndValidate(&journalQueryNoEnum)
	if err != nil {
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("Parameter error")
	}
	journalQueryNoEnum.Sort = &param.Sort{
		Fields: []string{"createTime,desc"},
	}
	journals, totalCount, err := j.JournalService.ListJournal(_ctx, param.AssertJournalQuery(journalQueryNoEnum))
	if err != nil {
		return nil, err
	}
	journalDTOs, err := j.JournalService.ConvertToWithCommentDTOList(_ctx, journals)
	if err != nil {
		return nil, err
	}
	return dto.NewPage(journalDTOs, totalCount, param.AssertJournalQuery(journalQueryNoEnum).Page), nil
}

func (j *JournalHandler) ListLatestJournal(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	top, err := util.MustGetQueryInt(_ctx, ctx, "top")
	if err != nil {
		top = 10
	}
	journalQuery := param.JournalQuery{
		Sort: &param.Sort{Fields: []string{"createTime,desc"}},
		Page: param.Page{PageNum: 0, PageSize: top},
	}
	journals, _, err := j.JournalService.ListJournal(_ctx, journalQuery)
	if err != nil {
		return nil, err
	}
	return j.JournalService.ConvertToWithCommentDTOList(_ctx, journals)
}

func (j *JournalHandler) CreateJournal(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	var journalParam param.Journal
	err := ctx.BindAndValidate(&journalParam)
	if err != nil {
		e := validator.ValidationErrors{}
		if errors.As(err, &e) {
			return nil, xerr.WithStatus(e, xerr.StatusBadRequest).WithMsg(trans.Translate(e))
		}
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("parameter error")
	}
	if journalParam.Content == "" {
		journalParam.Content = journalParam.SourceContent
	}
	journal, err := j.JournalService.Create(_ctx, &journalParam)
	if err != nil {
		return nil, err
	}
	return j.JournalService.ConvertToDTO(journal), nil
}

func (j *JournalHandler) UpdateJournal(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	var journalParam param.Journal
	err := ctx.BindAndValidate(&journalParam)
	if err != nil {
		e := validator.ValidationErrors{}
		if errors.As(err, &e) {
			return nil, xerr.WithStatus(e, xerr.StatusBadRequest).WithMsg(trans.Translate(e))
		}
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("parameter error")
	}

	journalID, err := util.ParamInt32(_ctx, ctx, "journalID")
	if err != nil {
		return nil, err
	}
	return j.JournalService.Update(_ctx, journalID, &journalParam)
}

func (j *JournalHandler) DeleteJournal(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	journalID, err := util.ParamInt32(_ctx, ctx, "journalID")
	if err != nil {
		return nil, err
	}
	return nil, j.JournalService.Delete(_ctx, journalID)
}
