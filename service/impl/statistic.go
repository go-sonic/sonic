package impl

import (
	"context"
	"time"

	"github.com/go-sonic/sonic/consts"
	"github.com/go-sonic/sonic/model/dto"
	"github.com/go-sonic/sonic/model/property"
	"github.com/go-sonic/sonic/service"
)

type statisticServiceImpl struct {
	PostService           service.PostService
	PostCommentService    service.PostCommentService
	JournalCommentService service.JournalCommentService
	JournalService        service.JournalService
	SheetCommentService   service.SheetCommentService
	TagService            service.TagService
	OptionService         service.OptionService
	CategoryService       service.CategoryService
	LinkService           service.LinkService
	SheetService          service.SheetService
	UserService           service.UserService
}

func NewStatisticService(postService service.PostService,
	postCommentService service.PostCommentService,
	journalCommentService service.JournalCommentService,
	sheetCommentService service.SheetCommentService,
	tagService service.TagService,
	optionService service.OptionService,
	categoryService service.CategoryService,
	linkService service.LinkService,
	sheetService service.SheetService,
	userService service.UserService,
	journalService service.JournalService,
) service.StatisticService {
	return &statisticServiceImpl{
		PostService:           postService,
		PostCommentService:    postCommentService,
		JournalCommentService: journalCommentService,
		SheetCommentService:   sheetCommentService,
		TagService:            tagService,
		OptionService:         optionService,
		CategoryService:       categoryService,
		LinkService:           linkService,
		SheetService:          sheetService,
		UserService:           userService,
		JournalService:        journalService,
	}
}

func (s statisticServiceImpl) Statistic(ctx context.Context) (*dto.Statistic, error) {
	var statistic dto.Statistic
	postCount, err := s.PostService.CountByStatus(ctx, consts.PostStatusPublished)
	if err != nil {
		return nil, err
	}
	postCommentCount, err := s.PostCommentService.CountByStatus(ctx, consts.CommentStatusPublished)
	if err != nil {
		return nil, err
	}
	sheetCommentCount, err := s.SheetCommentService.CountByStatus(ctx, consts.CommentStatusPublished)
	if err != nil {
		return nil, err
	}
	journalCommentCount, err := s.JournalCommentService.CountByStatus(ctx, consts.CommentStatusPublished)
	if err != nil {
		return nil, err
	}
	tagCount, err := s.TagService.CountAllTag(ctx)
	if err != nil {
		return nil, err
	}
	journalCount, err := s.JournalService.Count(ctx)
	if err != nil {
		return nil, err
	}
	categoryCount, err := s.CategoryService.Count(ctx)
	if err != nil {
		return nil, err
	}
	linkCount, err := s.LinkService.Count(ctx)
	if err != nil {
		return nil, err
	}
	postVisitCount, err := s.PostService.CountVisit(ctx)
	if err != nil {
		return nil, err
	}
	sheetVisitCount, err := s.SheetService.CountVisit(ctx)
	if err != nil {
		return nil, err
	}
	postLikeCount, err := s.PostService.CountVisit(ctx)
	if err != nil {
		return nil, err
	}
	sheetLikeCount, err := s.SheetService.CountLike(ctx)
	if err != nil {
		return nil, err
	}
	birthday, err := s.OptionService.GetOrByDefaultWithErr(ctx, property.BirthDay, property.BirthDay.DefaultValue)
	if err != nil {
		return nil, err
	}
	statistic.PostCount = postCount
	statistic.CommentCount = postCommentCount + journalCommentCount + sheetCommentCount
	statistic.TagCount = tagCount
	statistic.JournalCount = journalCount
	statistic.CategoryCount = categoryCount
	statistic.LinkCount = linkCount
	statistic.VisitCount = postVisitCount + sheetVisitCount
	statistic.LikeCount = postLikeCount + sheetLikeCount
	statistic.Birthday = birthday.(int64)
	statistic.EstablishDays = (time.Now().UnixMilli() - birthday.(int64)) / (1000 * 24 * 3600)
	return &statistic, nil
}

func (s statisticServiceImpl) StatisticWithUser(ctx context.Context) (*dto.StatisticWithUser, error) {
	statisticDTO, err := s.Statistic(ctx)
	if err != nil {
		return nil, err
	}
	user, err := MustGetAuthorizedUser(ctx)
	if err != nil {
		return nil, err
	}
	statisticDTOWithUser := &dto.StatisticWithUser{
		Statistic: *statisticDTO,
		User:      *s.UserService.ConvertToDTO(ctx, user),
	}
	return statisticDTOWithUser, nil
}
