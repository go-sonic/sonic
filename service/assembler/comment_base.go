package assembler

import (
	"context"

	"go.uber.org/zap"

	"github.com/go-sonic/sonic/log"
	"github.com/go-sonic/sonic/model/dto"
	"github.com/go-sonic/sonic/model/entity"
	"github.com/go-sonic/sonic/model/param"
	"github.com/go-sonic/sonic/model/property"
	"github.com/go-sonic/sonic/model/vo"
	"github.com/go-sonic/sonic/service"
	"github.com/go-sonic/sonic/util"
)

type BaseCommentAssembler interface {
	ConvertToDTO(ctx context.Context, comment *entity.Comment) (*dto.Comment, error)
	ConvertToDTOList(ctx context.Context, comments []*entity.Comment) ([]*dto.Comment, error)
	PageConvertToVOs(ctx context.Context, comments []*entity.Comment, page param.Page) ([]*vo.Comment, int64, error)
	ConvertToWithParentVO(ctx context.Context, comments []*entity.Comment) ([]*vo.CommentWithParent, error)
	ConvertToWithHasChildren(ctx context.Context, comments []*entity.Comment) ([]*vo.CommentWithHasChildren, error)
	ClearSensitiveField(ctx context.Context, comments []*entity.Comment) []*entity.Comment
}

func NewBaseCommentAssembler(
	optionService service.OptionService,
	baseCommentService service.BaseCommentService,
) BaseCommentAssembler {
	return &baseCommentAssembler{
		OptionService:      optionService,
		BaseCommentService: baseCommentService,
	}
}

type baseCommentAssembler struct {
	OptionService      service.OptionService
	BaseCommentService service.BaseCommentService
}

func (*baseCommentAssembler) ClearSensitiveField(ctx context.Context, comments []*entity.Comment) []*entity.Comment {
	for _, comment := range comments {
		comment.Email = ""
		comment.IPAddress = ""
	}
	return comments
}

func (b *baseCommentAssembler) ConvertToWithHasChildren(ctx context.Context, comments []*entity.Comment) ([]*vo.CommentWithHasChildren, error) {
	parentIDs := make([]int32, 0)
	for _, comment := range comments {
		parentIDs = append(parentIDs, comment.ID)
	}
	dtos, err := b.ConvertToDTOList(ctx, comments)
	if err != nil {
		return nil, err
	}
	countMap, err := b.BaseCommentService.CountChildren(ctx, parentIDs)
	if err != nil {
		return nil, err
	}
	result := make([]*vo.CommentWithHasChildren, 0, len(comments))
	for _, commentDTO := range dtos {
		commentWithHasChildren := &vo.CommentWithHasChildren{
			Comment: commentDTO,
		}
		if count, ok := countMap[commentDTO.ID]; ok && count > 0 {
			commentWithHasChildren.HasChildren = true
			commentWithHasChildren.ChildrenCount = count
		} else {
			commentWithHasChildren.HasChildren = false
			commentWithHasChildren.ChildrenCount = 0
		}
		result = append(result, commentWithHasChildren)
	}
	return result, nil
}

func (b *baseCommentAssembler) ConvertToDTO(ctx context.Context, comment *entity.Comment) (*dto.Comment, error) {
	commentDTO := &dto.Comment{
		ID:                comment.ID,
		Author:            comment.Author,
		Email:             comment.Email,
		IPAddress:         comment.IPAddress,
		AuthorURL:         comment.AuthorURL,
		GravatarMD5:       comment.GravatarMd5,
		Content:           comment.Content,
		Status:            comment.Status,
		UserAgent:         comment.UserAgent,
		ParentID:          comment.ParentID,
		IsAdmin:           comment.IsAdmin,
		AllowNotification: comment.AllowNotification,
		CreateTime:        comment.CreateTime.UnixMilli(),
		Likes:             comment.Likes,
	}
	avatarURL, err := b.BaseCommentService.BuildAvatarURL(ctx, comment.GravatarMd5, nil, nil)
	if err != nil {
		return nil, err
	}
	commentDTO.Avatar = avatarURL
	return commentDTO, nil
}

func (b *baseCommentAssembler) ConvertToDTOList(ctx context.Context, comments []*entity.Comment) ([]*dto.Comment, error) {
	gravatarSource, err := b.OptionService.GetOrByDefaultWithErr(ctx, property.CommentGravatarSource, property.CommentGravatarSource.DefaultValue)
	if err != nil {
		return nil, err
	}
	gravatarDefault, err := b.OptionService.GetOrByDefaultWithErr(ctx, property.CommentGravatarDefault, property.CommentGravatarDefault.DefaultValue)
	if err != nil {
		return nil, err
	}
	result := make([]*dto.Comment, 0, len(comments))
	for _, comment := range comments {
		commentDTO := &dto.Comment{
			ID:                comment.ID,
			Author:            comment.Author,
			Email:             comment.Email,
			IPAddress:         comment.IPAddress,
			AuthorURL:         comment.AuthorURL,
			GravatarMD5:       comment.GravatarMd5,
			Content:           comment.Content,
			Status:            comment.Status,
			UserAgent:         comment.UserAgent,
			ParentID:          comment.ParentID,
			IsAdmin:           comment.IsAdmin,
			AllowNotification: comment.AllowNotification,
			CreateTime:        comment.CreateTime.UnixMilli(),
			Likes:             comment.Likes,
		}
		avatarURL, err := b.BaseCommentService.BuildAvatarURL(ctx, comment.GravatarMd5, util.StringPtr(gravatarSource.(string)), util.StringPtr(gravatarDefault.(string)))
		if err != nil {
			return nil, err
		}
		commentDTO.Avatar = avatarURL
		result = append(result, commentDTO)
	}
	return result, nil
}

func (b *baseCommentAssembler) buildCommentTree(ctx context.Context, comments []*entity.Comment) ([]*vo.Comment, error) {
	commentIDMap := make(map[int32]*vo.Comment)
	commentDTOs, err := b.ConvertToDTOList(ctx, comments)
	if err != nil {
		return nil, err
	}
	for _, commentDTO := range commentDTOs {
		commentVO := &vo.Comment{
			Comment:  commentDTO,
			Children: make([]*vo.Comment, 0),
		}
		commentIDMap[commentDTO.ID] = commentVO
	}
	topComments := make([]*vo.Comment, 0)
	for _, comment := range comments {
		if comment.ParentID != 0 {
			parentComment, ok := commentIDMap[comment.ParentID]
			if !ok {
				log.CtxWarn(ctx, "parent comment does not exist", zap.Int32("postID", comment.PostID), zap.Int32("parentID", comment.ParentID))
				continue
			}
			parentComment.Children = append(parentComment.Children, commentIDMap[comment.ID])
		} else {
			topComments = append(topComments, commentIDMap[comment.ID])
		}
	}
	return topComments, nil
}

func (b *baseCommentAssembler) PageConvertToVOs(ctx context.Context, allComments []*entity.Comment, page param.Page) ([]*vo.Comment, int64, error) {
	topComments, err := b.buildCommentTree(ctx, allComments)
	if err != nil {
		return nil, 0, err
	}
	startIndex := page.PageNum * page.PageSize
	endIndex := startIndex + page.PageSize
	if startIndex > len(topComments) || startIndex < 0 {
		return make([]*vo.Comment, 0), 0, nil
	}
	if endIndex > len(topComments) {
		endIndex = len(topComments)
	}
	return topComments[startIndex:endIndex], int64(len(topComments)), nil
}

func (b *baseCommentAssembler) ConvertToWithParentVO(ctx context.Context, comments []*entity.Comment) ([]*vo.CommentWithParent, error) {
	parentIDs := make([]int32, 0)
	for _, comment := range comments {
		if comment.ParentID != 0 {
			parentIDs = append(parentIDs, comment.ParentID)
		}
	}
	dtos, err := b.ConvertToDTOList(ctx, comments)
	if err != nil {
		return nil, err
	}
	parentDTOMap := make(map[int32]*dto.Comment)
	if len(parentIDs) > 0 {
		parents, err := b.BaseCommentService.LGetByIDs(ctx, parentIDs)
		if err != nil {
			return nil, err
		}
		parentDTOs, err := b.ConvertToDTOList(ctx, parents)
		if err != nil {
			return nil, err
		}
		for _, parentDTO := range parentDTOs {
			parentDTOMap[parentDTO.ID] = parentDTO
		}
	}
	result := make([]*vo.CommentWithParent, 0, len(comments))
	for _, commentDTO := range dtos {
		commentWithParent := &vo.CommentWithParent{
			Comment: commentDTO,
		}
		if parent, ok := parentDTOMap[commentDTO.ParentID]; ok {
			commentWithParent.Parent = parent
		}
		result = append(result, commentWithParent)
	}
	return result, nil
}
