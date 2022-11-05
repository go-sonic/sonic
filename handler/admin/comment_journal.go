package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/go-sonic/sonic/consts"
	"github.com/go-sonic/sonic/handler/binding"
	"github.com/go-sonic/sonic/handler/trans"
	"github.com/go-sonic/sonic/model/dto"
	"github.com/go-sonic/sonic/model/param"
	"github.com/go-sonic/sonic/model/property"
	"github.com/go-sonic/sonic/service"
	"github.com/go-sonic/sonic/service/assembler"
	"github.com/go-sonic/sonic/service/impl"
	"github.com/go-sonic/sonic/util"
	"github.com/go-sonic/sonic/util/xerr"
)

type JournalCommentHandler struct {
	JournalCommentService   service.JournalCommentService
	OptionService           service.OptionService
	JournalService          service.JournalService
	JournalCommentAssembler assembler.JournalCommentAssembler
}

func NewJournalCommentHandler(journalCommentService service.JournalCommentService, optionService service.OptionService, journalService service.JournalService, journalCommentAssembler assembler.JournalCommentAssembler) *JournalCommentHandler {
	return &JournalCommentHandler{
		JournalCommentService:   journalCommentService,
		OptionService:           optionService,
		JournalService:          journalService,
		JournalCommentAssembler: journalCommentAssembler,
	}
}

func (j *JournalCommentHandler) ListJournalComment(ctx *gin.Context) (interface{}, error) {
	var commentQuery param.CommentQuery
	err := ctx.ShouldBindWith(&commentQuery, binding.CustomFormBinding)
	if err != nil {
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("Parameter error")
	}
	commentQuery.Sort = &param.Sort{
		Fields: []string{"createTime,desc"},
	}
	comments, totalCount, err := j.JournalCommentService.Page(ctx, commentQuery, consts.CommentTypeJournal)
	if err != nil {
		return nil, err
	}
	commentDTOs, err := j.JournalCommentAssembler.ConvertToWithJournal(ctx, comments)
	if err != nil {
		return nil, err
	}
	return dto.NewPage(commentDTOs, totalCount, commentQuery.Page), nil
}

func (j *JournalCommentHandler) ListJournalCommentLatest(ctx *gin.Context) (interface{}, error) {
	top, err := util.MustGetQueryInt32(ctx, "top")
	if err != nil {
		return nil, err
	}
	commentQuery := param.CommentQuery{
		Sort: &param.Sort{Fields: []string{"createTime,desc"}},
		Page: param.Page{PageNum: 0, PageSize: int(top)},
	}
	comments, _, err := j.JournalCommentService.Page(ctx, commentQuery, consts.CommentTypeSheet)
	if err != nil {
		return nil, err
	}
	return j.JournalCommentAssembler.ConvertToWithJournal(ctx, comments)
}

func (j *JournalCommentHandler) ListJournalCommentAsTree(ctx *gin.Context) (interface{}, error) {
	journalID, err := util.ParamInt32(ctx, "journalID")
	if err != nil {
		return nil, err
	}
	pageNum, err := util.MustGetQueryInt32(ctx, "page")
	if err != nil {
		return nil, err
	}
	pageSize, err := j.OptionService.GetOrByDefaultWithErr(ctx, property.CommentPageSize, property.CommentPageSize.DefaultValue)
	if err != nil {
		return nil, err
	}
	page := param.Page{PageSize: pageSize.(int), PageNum: int(pageNum)}

	allComments, err := j.JournalCommentService.GetByContentID(ctx, journalID, consts.CommentTypeJournal, &param.Sort{Fields: []string{"createTime,desc"}})
	if err != nil {
		return nil, err
	}

	commentVOs, totalCount, err := j.JournalCommentAssembler.PageConvertToVOs(ctx, allComments, page)
	if err != nil {
		return nil, err
	}
	return dto.NewPage(commentVOs, totalCount, page), nil
}

func (j *JournalCommentHandler) ListJournalCommentWithParent(ctx *gin.Context) (interface{}, error) {
	journalID, err := util.ParamInt32(ctx, "journalID")
	if err != nil {
		return nil, err
	}
	pageNum, err := util.MustGetQueryInt32(ctx, "page")
	if err != nil {
		return nil, err
	}

	pageSize, err := j.OptionService.GetOrByDefaultWithErr(ctx, property.CommentPageSize, property.CommentPageSize.DefaultValue)
	if err != nil {
		return nil, err
	}

	page := param.Page{PageSize: pageSize.(int), PageNum: int(pageNum)}

	comments, totalCount, err := j.JournalCommentService.Page(ctx, param.CommentQuery{
		ContentID: &journalID,
		Page:      page,
		Sort:      &param.Sort{Fields: []string{"createTime,desc"}},
	}, consts.CommentTypePost)
	if err != nil {
		return nil, err
	}

	commentsWithParent, err := j.JournalCommentAssembler.ConvertToWithParentVO(ctx, comments)
	if err != nil {
		return nil, err
	}
	return dto.NewPage(commentsWithParent, totalCount, page), nil
}

func (j *JournalCommentHandler) CreateJournalComment(ctx *gin.Context) (interface{}, error) {
	var commentParam *param.AdminComment
	err := ctx.ShouldBindJSON(&commentParam)
	if err != nil {
		if e, ok := err.(validator.ValidationErrors); ok {
			return nil, xerr.WithStatus(e, xerr.StatusBadRequest).WithMsg(trans.Translate(e))
		}
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("parameter error")
	}
	user, err := impl.MustGetAuthorizedUser(ctx)
	if err != nil || user == nil {
		return nil, err
	}
	blogURL, err := j.OptionService.GetBlogBaseURL(ctx)
	if err != nil {
		return nil, err
	}
	commonParam := param.Comment{
		Author:            user.Username,
		Email:             user.Email,
		AuthorURL:         blogURL,
		Content:           commentParam.Content,
		PostID:            commentParam.PostID,
		ParentID:          commentParam.ParentID,
		AllowNotification: true,
		CommentType:       consts.CommentTypeJournal,
	}
	comment, err := j.JournalCommentService.CreateBy(ctx, &commonParam)
	if err != nil {
		return nil, err
	}
	return j.JournalCommentAssembler.ConvertToDTO(ctx, comment)
}

func (j *JournalCommentHandler) UpdateJournalCommentStatus(ctx *gin.Context) (interface{}, error) {
	commentID, err := util.ParamInt32(ctx, "commentID")
	if err != nil {
		return nil, err
	}
	strStatus, err := util.ParamString(ctx, "status")
	if err != nil {
		return nil, err
	}
	status, err := consts.CommentStatusFromString(strStatus)
	if err != nil {
		return nil, err
	}
	return j.JournalCommentService.UpdateStatus(ctx, commentID, status)
}

func (j *JournalCommentHandler) UpdateJournalComment(ctx *gin.Context) (interface{}, error) {
	commentID, err := util.ParamInt32(ctx, "commentID")
	if err != nil {
		return nil, err
	}
	var commentParam *param.Comment
	err = ctx.ShouldBindJSON(&commentParam)
	if err != nil {
		if e, ok := err.(validator.ValidationErrors); ok {
			return nil, xerr.WithStatus(e, xerr.StatusBadRequest).WithMsg(trans.Translate(e))
		}
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("parameter error")
	}
	if commentParam.AuthorURL != "" {
		err = util.Validate.Var(commentParam.AuthorURL, "url")
		if err != nil {
			return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("url is not available")
		}
	}
	comment, err := j.JournalCommentService.UpdateBy(ctx, commentID, commentParam)
	if err != nil {
		return nil, err
	}
	return j.JournalCommentAssembler.ConvertToDTO(ctx, comment)
}

func (j *JournalCommentHandler) UpdateJournalStatusBatch(ctx *gin.Context) (interface{}, error) {
	status, err := util.ParamInt32(ctx, "status")
	if err != nil {
		return nil, err
	}

	ids := make([]int32, 0)
	err = ctx.ShouldBindJSON(&ids)
	if err != nil {
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("post ids error")
	}
	comments, err := j.JournalCommentService.UpdateStatusBatch(ctx, ids, consts.CommentStatus(status))
	if err != nil {
		return nil, err
	}
	return j.JournalCommentAssembler.ConvertToDTOList(ctx, comments)
}

func (j *JournalCommentHandler) DeleteJournalComment(ctx *gin.Context) (interface{}, error) {
	commentID, err := util.ParamInt32(ctx, "commentID")
	if err != nil {
		return nil, err
	}
	return nil, j.JournalCommentService.Delete(ctx, commentID)
}

func (j *JournalCommentHandler) DeleteJournalCommentBatch(ctx *gin.Context) (interface{}, error) {
	ids := make([]int32, 0)
	err := ctx.ShouldBindJSON(&ids)
	if err != nil {
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("post ids error")
	}
	return nil, j.JournalCommentService.DeleteBatch(ctx, ids)
}
