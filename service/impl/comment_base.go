package impl

import (
	"context"

	"gorm.io/gen/field"

	"github.com/go-sonic/sonic/consts"
	"github.com/go-sonic/sonic/dal"
	"github.com/go-sonic/sonic/event"
	"github.com/go-sonic/sonic/model/entity"
	"github.com/go-sonic/sonic/model/param"
	"github.com/go-sonic/sonic/model/property"
	"github.com/go-sonic/sonic/service"
	"github.com/go-sonic/sonic/util"
	"github.com/go-sonic/sonic/util/xerr"
)

type baseCommentServiceImpl struct {
	UserService   service.UserService
	OptionService service.OptionService
	Event         event.Bus
}

func (b baseCommentServiceImpl) LGetByIDs(ctx context.Context, commentIDs []int32) ([]*entity.Comment, error) {
	commentDAL := dal.GetQueryByCtx(ctx).Comment
	comments, err := commentDAL.WithContext(ctx).Where(commentDAL.ID.In(commentIDs...)).Find()
	return comments, WrapDBErr(err)
}

func NewBaseCommentService(userService service.UserService, optionService service.OptionService, event event.Bus) service.BaseCommentService {
	return &baseCommentServiceImpl{
		UserService:   userService,
		OptionService: optionService,
		Event:         event,
	}
}

func (b baseCommentServiceImpl) Page(ctx context.Context, commentQuery param.CommentQuery, commentType consts.CommentType) ([]*entity.Comment, int64, error) {
	commentDAL := dal.GetQueryByCtx(ctx).Comment
	commentDO := commentDAL.WithContext(ctx).Where(commentDAL.Type.Eq(commentType))
	err := BuildSort(commentQuery.Sort, &commentDAL, &commentDO)
	if err != nil {
		return nil, 0, err
	}
	if commentQuery.Keyword != nil && *commentQuery.Keyword != "" {
		commentDO.Where(commentDAL.Content.Like("%" + *commentQuery.Keyword + "%"))
	}
	if commentQuery.CommentStatus != nil {
		commentDO.Where(commentDAL.Status.Eq(*commentQuery.CommentStatus))
	}
	if commentQuery.ParentID != nil {
		commentDO.Where(commentDAL.ParentID.Eq(*commentQuery.ParentID))
	}
	if commentQuery.ContentID != nil {
		commentDO.Where(commentDAL.PostID.Eq(*commentQuery.ContentID))
	}
	comments, totalCount, err := commentDO.FindByPage(commentQuery.PageNum*commentQuery.PageSize, commentQuery.PageSize)
	if err != nil {
		return nil, 0, WrapDBErr(err)
	}
	return comments, totalCount, nil
}

func (b baseCommentServiceImpl) Update(ctx context.Context, comment *entity.Comment) (*entity.Comment, error) {
	if comment.ID == 0 {
		return nil, nil
	}
	commentDAL := dal.GetQueryByCtx(ctx).Comment
	updateResult, err := commentDAL.WithContext(ctx).Where(commentDAL.ID.Eq(comment.ID)).Updates(comment)
	if err != nil {
		return nil, WrapDBErr(err)
	}
	if updateResult.RowsAffected != 1 {
		return nil, xerr.NoType.New("").WithMsg("update comment failed")
	}
	updatedComment, err := commentDAL.WithContext(ctx).Where(commentDAL.ID.Eq(comment.ID)).First()
	if err != nil {
		return nil, WrapDBErr(err)
	}
	return updatedComment, nil
}

func (b baseCommentServiceImpl) GetByID(ctx context.Context, commentID int32) (*entity.Comment, error) {
	commentDAL := dal.GetQueryByCtx(ctx).Comment
	comment, err := commentDAL.WithContext(ctx).Where(commentDAL.ID.Eq(commentID)).First()
	if err != nil {
		return nil, WrapDBErr(err)
	}
	return comment, nil
}

func (b baseCommentServiceImpl) GetByContentID(ctx context.Context, contentID int32, commentType consts.CommentType, sort *param.Sort) ([]*entity.Comment, error) {
	commentDAL := dal.GetQueryByCtx(ctx).Comment
	commentDO := commentDAL.WithContext(ctx).Where(commentDAL.PostID.Eq(contentID), commentDAL.Type.Eq(commentType))

	err := BuildSort(sort, &commentDAL, &commentDO)
	if err != nil {
		return nil, err
	}
	comments, err := commentDO.Find()
	return comments, WrapDBErr(err)
}

func (b baseCommentServiceImpl) DeleteBatch(ctx context.Context, commentIDs []int32) error {
	commentDAL := dal.GetQueryByCtx(ctx).Comment
	deleteResult, err := commentDAL.WithContext(ctx).Where(commentDAL.ID.In(commentIDs...)).Delete()
	if err != nil {
		return WrapDBErr(err)
	}
	if deleteResult.RowsAffected != int64(len(commentIDs)) {
		return xerr.NoType.New("").WithMsg("delete comment failed")
	}
	return nil
}

func (b baseCommentServiceImpl) Delete(ctx context.Context, commentID int32) error {
	commentDAL := dal.GetQueryByCtx(ctx).Comment
	deleteResult, err := commentDAL.WithContext(ctx).Where(commentDAL.ID.Eq(commentID)).Delete()
	if err != nil {
		return WrapDBErr(err)
	}
	if deleteResult.RowsAffected != 1 {
		return xerr.NoType.New("").WithMsg("delete comment failed")
	}
	return nil
}

func (b baseCommentServiceImpl) UpdateStatusBatch(ctx context.Context, commentIDs []int32, commentStatus consts.CommentStatus) ([]*entity.Comment, error) {
	commentDAL := dal.GetQueryByCtx(ctx).Comment
	_, err := commentDAL.WithContext(ctx).Where(commentDAL.ID.In(commentIDs...)).UpdateSimple(commentDAL.Status.Value(commentStatus))
	if err != nil {
		return nil, WrapDBErr(err)
	}
	comments, err := commentDAL.WithContext(ctx).Where(commentDAL.ID.In(commentIDs...)).Find()
	if err != nil {
		return nil, WrapDBErr(err)
	}
	return comments, nil
}

func (b baseCommentServiceImpl) UpdateStatus(ctx context.Context, commentID int32, commentStatus consts.CommentStatus) (*entity.Comment, error) {
	commentDAL := dal.GetQueryByCtx(ctx).Comment
	comment, err := commentDAL.WithContext(ctx).Where(commentDAL.ID.Eq(commentID)).First()
	if err != nil {
		return nil, err
	}
	updateResult, err := commentDAL.WithContext(ctx).Where(commentDAL.ID.Eq(commentID)).UpdateSimple(commentDAL.Status.Value(commentStatus))
	if err != nil {
		return nil, WrapDBErr(err)
	}
	if updateResult.RowsAffected != 1 {
		return nil, xerr.NoType.New("").WithMsg("update comment status failed")
	}
	comment.Status = commentStatus
	if comment.ParentID != 0 {
		go func() {
			b.Event.Publish(context.TODO(), &event.CommentReplyEvent{
				Comment: comment,
			})
		}()
	}
	return comment, nil
}

func (b baseCommentServiceImpl) Create(ctx context.Context, comment *entity.Comment) (*entity.Comment, error) {
	if comment == nil {
		return nil, xerr.BadParam.New("comment can not be empty")
	}
	commentDAL := dal.GetQueryByCtx(ctx).Comment
	if comment.ParentID != 0 {
		_, err := commentDAL.WithContext(ctx).Where(commentDAL.ID.Eq(comment.ParentID)).First()
		err = WrapDBErr(err)
		if xerr.GetType(err) == xerr.NoRecord {
			return nil, xerr.WithMsg(err, "parent comment does not exist").WithStatus(xerr.StatusNotFound)
		}
		if err != nil {
			return nil, err
		}
	}

	comment.IPAddress = util.GetClientIP(ctx)
	comment.UserAgent = util.GetUserAgent(ctx)
	if comment.GravatarMd5 == "" {
		comment.GravatarMd5 = util.Md5Hex(comment.Email)
	}

	authentication, _ := GetAuthorizedUser(ctx)
	if authentication == nil {
		comment.IsAdmin = false
		needCheck, err := b.OptionService.GetOrByDefaultWithErr(ctx, property.CommentNewNeedCheck, true)
		if err != nil {
			return nil, err
		}
		if needCheck.(bool) {
			comment.Status = consts.CommentStatusAuditing
		} else {
			comment.Status = consts.CommentStatusPublished
		}
	} else {
		comment.Email = authentication.Email
		comment.Author = util.IfElse(authentication.Nickname == "", authentication.Username, authentication.Nickname).(string)
		blogURL, err := b.OptionService.GetBlogBaseURL(ctx)
		if err != nil {
			return nil, err
		}
		comment.AuthorURL = blogURL
		comment.IsAdmin = true
		comment.Status = consts.CommentStatusPublished
	}
	comment.GravatarMd5 = util.Md5Hex(comment.Email)
	err := commentDAL.WithContext(ctx).Select(field.Star).Omit(commentDAL.UpdateTime).Create(comment)
	if err != nil {
		return nil, WrapDBErr(err)
	}
	if comment.ParentID != 0 {
		go func() {
			b.Event.Publish(context.TODO(), &event.CommentReplyEvent{
				Comment: comment,
			})
		}()
	} else {
		go func() {
			b.Event.Publish(context.TODO(), &event.CommentNewEvent{
				Comment: comment,
			})
		}()
	}
	return comment, nil
}

func (b baseCommentServiceImpl) BuildAvatarURL(ctx context.Context, gravatarMD5 string, gravatarSource, gravatarDefault *string) (string, error) {
	var gs, gd string
	if gravatarSource == nil {
		temp, err := b.OptionService.GetOrByDefaultWithErr(ctx, property.CommentGravatarSource, property.CommentGravatarSource.DefaultValue)
		if err != nil {
			return "", err
		}
		gs = temp.(string)
	} else {
		gs = *gravatarSource
	}
	if gravatarDefault == nil {
		temp, err := b.OptionService.GetOrByDefaultWithErr(ctx, property.CommentGravatarDefault, property.CommentGravatarDefault.DefaultValue)
		if err != nil {
			return "", err
		}
		gd = temp.(string)
	} else {
		gd = *gravatarDefault
	}

	return gs + gravatarMD5 + "?s=256&d=" + gd, nil
}

func (b baseCommentServiceImpl) ConvertParam(commentParam *param.Comment) *entity.Comment {
	comment := &entity.Comment{
		Author:            commentParam.Author,
		Email:             commentParam.Email,
		AuthorURL:         commentParam.AuthorURL,
		Content:           commentParam.Content,
		PostID:            commentParam.PostID,
		ParentID:          commentParam.ParentID,
		Type:              commentParam.CommentType,
		AllowNotification: commentParam.AllowNotification,
	}

	return comment
}

func (b baseCommentServiceImpl) CountByContentID(ctx context.Context, contentID int32, commentType consts.CommentType, status consts.CommentStatus) (int64, error) {
	postCommentDAL := dal.GetQueryByCtx(ctx).Comment
	count, err := postCommentDAL.WithContext(ctx).Where(postCommentDAL.PostID.Eq(contentID), postCommentDAL.Type.Eq(commentType), postCommentDAL.Status.Eq(status)).Count()
	if err != nil {
		return 0, WrapDBErr(err)
	}
	return count, nil
}

func (b baseCommentServiceImpl) CountByStatusAndContentIDs(ctx context.Context, status consts.CommentStatus, contentIDs []int32) (map[int32]int64, error) {
	var projections []struct {
		PostID       int32
		CommentCount int64 `gorm:"column:comment_count"`
	}

	postCommentDAL := dal.GetQueryByCtx(ctx).Comment
	err := postCommentDAL.WithContext(ctx).Select(postCommentDAL.PostID, postCommentDAL.ID.Count().As("comment_count")).Where(postCommentDAL.Status.Eq(status), postCommentDAL.PostID.In(contentIDs...)).Group(postCommentDAL.PostID).Scan(&projections)
	if err != nil {
		return nil, WrapDBErr(err)
	}
	result := make(map[int32]int64)
	for _, projection := range projections {
		result[projection.PostID] = projection.CommentCount
	}
	return result, nil
}

func (b baseCommentServiceImpl) CreateBy(ctx context.Context, commentParam *param.Comment) (*entity.Comment, error) {
	comment := b.ConvertParam(commentParam)
	return b.Create(ctx, comment)
}

func (*baseCommentServiceImpl) CountChildren(ctx context.Context, parentCommentIDs []int32) (map[int32]int64, error) {
	var projections []struct {
		ParentID     int32
		CommentCount int64
	}

	commentDAL := dal.GetQueryByCtx(ctx).Comment
	err := commentDAL.WithContext(ctx).Select(commentDAL.ParentID, commentDAL.ID.Count().As("comment_count")).Where(commentDAL.Status.Eq(consts.CommentStatusPublished), commentDAL.ParentID.In(parentCommentIDs...)).Group(commentDAL.ParentID).Scan(&projections)
	if err != nil {
		return nil, WrapDBErr(err)
	}
	result := make(map[int32]int64)
	for _, projection := range projections {
		result[projection.ParentID] = projection.CommentCount
	}
	return result, nil
}

func (b *baseCommentServiceImpl) GetChildren(ctx context.Context, parentCommentID int32, contentID int32, commentType consts.CommentType) ([]*entity.Comment, error) {
	allComments, err := b.GetByContentID(ctx, contentID, commentType, nil)
	if err != nil {
		return nil, err
	}
	children := make([]*entity.Comment, 0)
	parentIDMap := make(map[int32][]*entity.Comment, 0)
	for _, comment := range allComments {
		if comment.Status == consts.CommentStatusPublished {
			parentIDMap[comment.ParentID] = append(parentIDMap[comment.ParentID], comment)
		}
	}
	queue := util.NewQueue[int32]()
	queue.Push(parentCommentID)
	for !queue.IsEmpty() {
		parentID := queue.Next()
		children = append(children, parentIDMap[parentID]...)
		for _, child := range parentIDMap[parentID] {
			queue.Push(child.ID)
		}
	}
	return children, nil
}

func (b baseCommentServiceImpl) IncreaseLike(ctx context.Context, commentID int32) error {
	commentDAL := dal.GetQueryByCtx(ctx).Comment
	info, err := commentDAL.WithContext(ctx).Where(commentDAL.ID.Eq(commentID)).UpdateSimple(commentDAL.Likes.Add(1))
	if err != nil {
		return WrapDBErr(err)
	}
	if info.RowsAffected != 1 {
		return xerr.NoType.New("increase comment like failed postID=%v", commentID).WithStatus(xerr.StatusBadRequest).WithMsg("failed to like comment")
	}
	return nil
}
