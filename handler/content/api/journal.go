package api

import (
	"html/template"

	"github.com/gin-gonic/gin"

	"github.com/go-sonic/sonic/consts"
	"github.com/go-sonic/sonic/handler/binding"
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

func (j *JournalHandler) ListJournal(ctx *gin.Context) (interface{}, error) {
	var journalQuery param.JournalQuery
	err := ctx.ShouldBindWith(&journalQuery, binding.CustomFormBinding)
	if err != nil {
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("Parameter error")
	}
	journalQuery.Sort = &param.Sort{
		Fields: []string{"createTime,desc"},
	}
	journalQuery.JournalType = consts.JournalTypePublic.Ptr()
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

func (j *JournalHandler) GetJournal(ctx *gin.Context) (interface{}, error) {
	journalID, err := util.ParamInt32(ctx, "journalID")
	if err != nil {
		return nil, err
	}
	journals, err := j.JournalService.GetByJournalIDs(ctx, []int32{journalID})
	if err != nil {
		return nil, err
	}
	if len(journals) == 0 {
		return nil, xerr.WithStatus(nil, xerr.StatusBadRequest)
	}
	journalDTOs, err := j.JournalService.ConvertToWithCommentDTOList(ctx, []*entity.Journal{journals[journalID]})
	if err != nil {
		return nil, err
	}
	return journalDTOs[0], nil
}

func (j *JournalHandler) ListTopComment(ctx *gin.Context) (interface{}, error) {
	journalID, err := util.ParamInt32(ctx, "journalID")
	if err != nil {
		return nil, err
	}
	pageSize := j.OptionService.GetOrByDefault(ctx, property.CommentPageSize).(int)

	commentQuery := param.CommentQuery{}
	err = ctx.ShouldBindWith(&commentQuery, binding.CustomFormBinding)
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

	comments, totalCount, err := j.JournalCommentService.Page(ctx, commentQuery, consts.CommentTypeJournal)
	if err != nil {
		return nil, err
	}
	_ = j.JournalCommentAssembler.ClearSensitiveField(ctx, comments)
	commenVOs, err := j.JournalCommentAssembler.ConvertToWithHasChildren(ctx, comments)
	if err != nil {
		return nil, err
	}
	return dto.NewPage(commenVOs, totalCount, commentQuery.Page), nil
}

func (j *JournalHandler) ListChildren(ctx *gin.Context) (interface{}, error) {
	journalID, err := util.ParamInt32(ctx, "journalID")
	if err != nil {
		return nil, err
	}
	parentID, err := util.ParamInt32(ctx, "parentID")
	if err != nil {
		return nil, err
	}
	children, err := j.JournalCommentService.GetChildren(ctx, parentID, journalID, consts.CommentTypeJournal)
	if err != nil {
		return nil, err
	}
	_ = j.JournalCommentAssembler.ClearSensitiveField(ctx, children)
	return j.JournalCommentAssembler.ConvertToDTOList(ctx, children)
}

func (j *JournalHandler) ListCommentTree(ctx *gin.Context) (interface{}, error) {
	journalID, err := util.ParamInt32(ctx, "journalID")
	if err != nil {
		return nil, err
	}
	pageSize := j.OptionService.GetOrByDefault(ctx, property.CommentPageSize).(int)

	commentQuery := param.CommentQuery{}
	err = ctx.ShouldBindWith(&commentQuery, binding.CustomFormBinding)
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

	allComments, err := j.JournalCommentService.GetByContentID(ctx, journalID, consts.CommentTypeJournal, commentQuery.Sort)
	if err != nil {
		return nil, err
	}
	_ = j.JournalCommentAssembler.ClearSensitiveField(ctx, allComments)
	commentVOs, total, err := j.JournalCommentAssembler.PageConvertToVOs(ctx, allComments, commentQuery.Page)
	if err != nil {
		return nil, err
	}
	return dto.NewPage(commentVOs, total, commentQuery.Page), nil
}

func (j *JournalHandler) ListComment(ctx *gin.Context) (interface{}, error) {
	journalID, err := util.ParamInt32(ctx, "journalID")
	if err != nil {
		return nil, err
	}
	pageSize := j.OptionService.GetOrByDefault(ctx, property.CommentPageSize).(int)

	commentQuery := param.CommentQuery{}
	err = ctx.ShouldBindWith(&commentQuery, binding.CustomFormBinding)
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

	comments, total, err := j.JournalCommentService.Page(ctx, commentQuery, consts.CommentTypeJournal)
	if err != nil {
		return nil, err
	}
	_ = j.JournalCommentAssembler.ClearSensitiveField(ctx, comments)
	result, err := j.JournalCommentAssembler.ConvertToWithParentVO(ctx, comments)
	if err != nil {
		return nil, err
	}
	return dto.NewPage(result, total, commentQuery.Page), nil
}

func (j *JournalHandler) CreateComment(ctx *gin.Context) (interface{}, error) {
	p := param.Comment{}
	err := ctx.ShouldBindJSON(&p)
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
	result, err := j.JournalCommentService.CreateBy(ctx, &p)
	if err != nil {
		return nil, err
	}
	return j.JournalCommentAssembler.ConvertToDTO(ctx, result)
}

func (j *JournalHandler) Like(ctx *gin.Context) (interface{}, error) {
	journalID, err := util.ParamInt32(ctx, "journalID")
	if err != nil {
		return nil, err
	}
	err = j.JournalService.IncreaseLike(ctx, journalID)
	if err != nil {
		return nil, err
	}
	return nil, err
}
