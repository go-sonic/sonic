package api

import (
	"context"
	"html/template"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/go-sonic/sonic/consts"
	"github.com/go-sonic/sonic/model/dto"
	"github.com/go-sonic/sonic/model/entity"
	"github.com/go-sonic/sonic/model/param"
	"github.com/go-sonic/sonic/model/property"
	"github.com/go-sonic/sonic/service"
	"github.com/go-sonic/sonic/service/assembler"
	"github.com/go-sonic/sonic/util"
	"github.com/go-sonic/sonic/util/xerr"
)

type JournalHandler struct {
	JournalService          service.JournalService
	JournalCommentService   service.JournalCommentService
	OptionService           service.ClientOptionService
	JournalCommentAssembler assembler.JournalCommentAssembler
}

func NewJournalHandler(
	journalService service.JournalService,
	journalCommentService service.JournalCommentService,
	optionService service.ClientOptionService,
	journalCommentAssembler assembler.JournalCommentAssembler,
) *JournalHandler {
	return &JournalHandler{
		JournalService:          journalService,
		JournalCommentService:   journalCommentService,
		OptionService:           optionService,
		JournalCommentAssembler: journalCommentAssembler,
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
	journalQuery := param.AssertJournalQuery(journalQueryNoEnum)
	journalQuery.JournalType = consts.JournalTypePublic.Ptr()
	journals, totalCount, err := j.JournalService.ListJournal(_ctx, journalQuery)
	if err != nil {
		return nil, err
	}
	journalDTOs, err := j.JournalService.ConvertToWithCommentDTOList(_ctx, journals)
	if err != nil {
		return nil, err
	}
	return dto.NewPage(journalDTOs, totalCount, journalQuery.Page), nil
}

func (j *JournalHandler) GetJournal(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	journalID, err := util.ParamInt32(_ctx, ctx, "journalID")
	if err != nil {
		return nil, err
	}
	journals, err := j.JournalService.GetByJournalIDs(_ctx, []int32{journalID})
	if err != nil {
		return nil, err
	}
	if len(journals) == 0 {
		return nil, xerr.WithStatus(nil, xerr.StatusBadRequest)
	}
	journalDTOs, err := j.JournalService.ConvertToWithCommentDTOList(_ctx, []*entity.Journal{journals[journalID]})
	if err != nil {
		return nil, err
	}
	return journalDTOs[0], nil
}

func (j *JournalHandler) ListTopComment(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	journalID, err := util.ParamInt32(_ctx, ctx, "journalID")
	if err != nil {
		return nil, err
	}
	pageSize := j.OptionService.GetOrByDefault(_ctx, property.CommentPageSize).(int)

	commentQuery := param.CommentQuery{}
	err = ctx.BindAndValidate(&commentQuery)
	if err != nil {
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("Parameter error")
	}
	if commentQuery.Sort != nil && len(commentQuery.Fields) > 0 {
		commentQuery.Sort = &param.Sort{
			Fields: []string{"createTime,desc"},
		}
	}
	commentQuery.ContentID = &journalID
	commentQuery.Keyword = nil
	commentQuery.CommentStatus = consts.CommentStatusPublished.Ptr()
	commentQuery.PageSize = pageSize
	commentQuery.ParentID = util.Int32Ptr(0)

	comments, totalCount, err := j.JournalCommentService.Page(_ctx, commentQuery, consts.CommentTypeJournal)
	if err != nil {
		return nil, err
	}
	_ = j.JournalCommentAssembler.ClearSensitiveField(_ctx, comments)
	commenVOs, err := j.JournalCommentAssembler.ConvertToWithHasChildren(_ctx, comments)
	if err != nil {
		return nil, err
	}
	return dto.NewPage(commenVOs, totalCount, commentQuery.Page), nil
}

func (j *JournalHandler) ListChildren(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	journalID, err := util.ParamInt32(_ctx, ctx, "journalID")
	if err != nil {
		return nil, err
	}
	parentID, err := util.ParamInt32(_ctx, ctx, "parentID")
	if err != nil {
		return nil, err
	}
	children, err := j.JournalCommentService.GetChildren(_ctx, parentID, journalID, consts.CommentTypeJournal)
	if err != nil {
		return nil, err
	}
	_ = j.JournalCommentAssembler.ClearSensitiveField(_ctx, children)
	return j.JournalCommentAssembler.ConvertToDTOList(_ctx, children)
}

func (j *JournalHandler) ListCommentTree(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	journalID, err := util.ParamInt32(_ctx, ctx, "journalID")
	if err != nil {
		return nil, err
	}
	pageSize := j.OptionService.GetOrByDefault(_ctx, property.CommentPageSize).(int)

	commentQueryNoEnum := param.CommentQueryNoEnum{}
	err = ctx.BindAndValidate(&commentQueryNoEnum)
	if err != nil {
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("Parameter error")
	}
	if commentQueryNoEnum.Sort != nil && len(commentQueryNoEnum.Fields) > 0 {
		commentQueryNoEnum.Sort = &param.Sort{
			Fields: []string{"createTime,desc"},
		}
	}
	commentQuery := param.AssertCommentQuery(commentQueryNoEnum)
	commentQuery.ContentID = &journalID
	commentQuery.Keyword = nil
	commentQuery.CommentStatus = consts.CommentStatusPublished.Ptr()
	commentQuery.PageSize = pageSize
	commentQuery.ParentID = util.Int32Ptr(0)

	allComments, err := j.JournalCommentService.GetByContentID(_ctx, journalID, consts.CommentTypeJournal, commentQuery.Sort)
	if err != nil {
		return nil, err
	}
	_ = j.JournalCommentAssembler.ClearSensitiveField(_ctx, allComments)
	commentVOs, total, err := j.JournalCommentAssembler.PageConvertToVOs(_ctx, allComments, commentQuery.Page)
	if err != nil {
		return nil, err
	}
	return dto.NewPage(commentVOs, total, commentQuery.Page), nil
}

func (j *JournalHandler) ListComment(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	journalID, err := util.ParamInt32(_ctx, ctx, "journalID")
	if err != nil {
		return nil, err
	}
	pageSize := j.OptionService.GetOrByDefault(_ctx, property.CommentPageSize).(int)

	commentQuery := param.CommentQuery{}
	err = ctx.BindAndValidate(&commentQuery)
	if err != nil {
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("Parameter error")
	}
	if commentQuery.Sort != nil && len(commentQuery.Fields) > 0 {
		commentQuery.Sort = &param.Sort{
			Fields: []string{"createTime,desc"},
		}
	}
	commentQuery.ContentID = &journalID
	commentQuery.Keyword = nil
	commentQuery.CommentStatus = consts.CommentStatusPublished.Ptr()
	commentQuery.PageSize = pageSize
	commentQuery.ParentID = util.Int32Ptr(0)

	comments, total, err := j.JournalCommentService.Page(_ctx, commentQuery, consts.CommentTypeJournal)
	if err != nil {
		return nil, err
	}
	_ = j.JournalCommentAssembler.ClearSensitiveField(_ctx, comments)
	result, err := j.JournalCommentAssembler.ConvertToWithParentVO(_ctx, comments)
	if err != nil {
		return nil, err
	}
	return dto.NewPage(result, total, commentQuery.Page), nil
}

func (j *JournalHandler) CreateComment(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	p := param.Comment{}
	err := ctx.BindAndValidate(&p)
	if err != nil {
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("Parameter error")
	}
	if p.AuthorURL != "" {
		err = util.Validate.Var(p.AuthorURL, "http_url")
		if err != nil {
			return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("Parameter error")
		}
	}
	p.Author = template.HTMLEscapeString(p.Author)
	p.AuthorURL = template.HTMLEscapeString(p.AuthorURL)
	p.Content = template.HTMLEscapeString(p.Content)
	p.Email = template.HTMLEscapeString(p.Email)
	p.CommentType = consts.CommentTypeJournal
	result, err := j.JournalCommentService.CreateBy(_ctx, &p)
	if err != nil {
		return nil, err
	}
	return j.JournalCommentAssembler.ConvertToDTO(_ctx, result)
}

func (j *JournalHandler) Like(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	journalID, err := util.ParamInt32(_ctx, ctx, "journalID")
	if err != nil {
		return nil, err
	}
	err = j.JournalService.IncreaseLike(_ctx, journalID)
	if err != nil {
		return nil, err
	}
	return nil, err
}
