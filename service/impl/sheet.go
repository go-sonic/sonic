package impl

import (
	"context"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/go-sonic/sonic/cache"
	"github.com/go-sonic/sonic/consts"
	"github.com/go-sonic/sonic/dal"
	"github.com/go-sonic/sonic/event"
	"github.com/go-sonic/sonic/model/dto"
	"github.com/go-sonic/sonic/model/entity"
	"github.com/go-sonic/sonic/model/param"
	"github.com/go-sonic/sonic/service"
	"github.com/go-sonic/sonic/util"
)

type sheetServiceImpl struct {
	service.BasePostService
	MetaService         service.MetaService
	OptionService       service.OptionService
	SheetCommentService service.SheetCommentService
	Event               event.Bus
	Cache               cache.Cache
	ThemeService        service.ThemeService
}

func NewSheetService(basePostService service.BasePostService,
	metaService service.MetaService,
	optionService service.OptionService,
	sheetCommentService service.SheetCommentService,
	event event.Bus,
	cache cache.Cache,
	themeService service.ThemeService,
) service.SheetService {
	return &sheetServiceImpl{
		BasePostService:     basePostService,
		MetaService:         metaService,
		OptionService:       optionService,
		SheetCommentService: sheetCommentService,
		Event:               event,
		Cache:               cache,
		ThemeService:        themeService,
	}
}

func (s sheetServiceImpl) Page(ctx context.Context, page param.Page, sort *param.Sort) ([]*entity.Post, int64, error) {
	sheetDAL := dal.GetQueryByCtx(ctx).Post
	sheetDO := sheetDAL.WithContext(ctx).Where(sheetDAL.Type.Eq(consts.PostTypeSheet))
	err := BuildSort(sort, &sheetDAL, &sheetDO)
	if err != nil {
		return nil, 0, err
	}
	sheets, totalCount, err := sheetDO.FindByPage(page.PageSize*page.PageNum, page.PageSize)
	if err != nil {
		return nil, 0, err
	}
	return sheets, totalCount, nil
}

func (s sheetServiceImpl) Create(ctx context.Context, sheetParam *param.Sheet) (*entity.Post, error) {
	sheet, err := s.ConvertParam(ctx, sheetParam)
	if err != nil {
		return nil, err
	}
	sheet, err = s.CreateOrUpdate(ctx, sheet, nil, nil, sheetParam.Metas)
	if err != nil {
		return nil, err
	}
	return sheet, nil
}

func (s sheetServiceImpl) ConvertParam(ctx context.Context, sheetParam *param.Sheet) (*entity.Post, error) {
	sheet := &entity.Post{
		Type:            consts.PostTypeSheet,
		DisallowComment: sheetParam.DisallowComment,
		OriginalContent: sheetParam.OriginalContent,
		FormatContent:   sheetParam.Content,
		Password:        sheetParam.Password,
		MetaDescription: sheetParam.MetaDescription,
		MetaKeywords:    sheetParam.MetaKeywords,
		Template:        sheetParam.Template,
		Thumbnail:       sheetParam.Thumbnail,
		Title:           sheetParam.Title,
		TopPriority:     sheetParam.TopPriority,
		Status:          sheetParam.Status,
		Summary:         sheetParam.Summary,
		EditTime:        util.TimePtr(time.Now()),
	}
	if sheetParam.EditorType != nil {
		sheet.EditorType = *sheetParam.EditorType
	} else {
		sheet.EditorType = consts.EditorTypeMarkdown
	}

	sheet.WordCount = util.HTMLFormatWordCount(sheet.FormatContent)
	if sheetParam.Slug == "" {
		sheet.Slug = util.Slug(sheetParam.Title)
	} else {
		sheet.Slug = util.Slug(sheetParam.Slug)
	}
	if sheetParam.CreateTime != nil {
		sheet.CreateTime = time.Unix(*sheetParam.CreateTime, 0)
	}
	return sheet, nil
}

func (s sheetServiceImpl) Update(ctx context.Context, sheetID int32, sheetParam *param.Sheet) (*entity.Post, error) {
	sheetDAL := dal.GetQueryByCtx(ctx).Post
	sheet, err := sheetDAL.WithContext(ctx).Where(sheetDAL.ID.Eq(sheetID)).First()
	if err != nil {
		return nil, WrapDBErr(err)
	}
	sheetToUpdate, err := s.ConvertParam(ctx, sheetParam)
	if err != nil {
		return nil, err
	}
	sheetToUpdate.ID = sheet.ID
	sheetToUpdate.Type = sheet.Type
	if sheetToUpdate.CreateTime == (time.Time{}) {
		sheetToUpdate.CreateTime = sheet.CreateTime
	}
	sheetToUpdate.CreateTime = sheet.CreateTime
	sheetToUpdate.Likes = sheet.Likes
	sheetToUpdate.Visits = sheet.Visits

	sheet, err = s.CreateOrUpdate(ctx, sheetToUpdate, nil, nil, sheetParam.Metas)
	if err != nil {
		return nil, err
	}
	s.Event.Publish(ctx, &event.LogEvent{
		LogKey:    strconv.Itoa(int(sheet.ID)),
		LogType:   consts.LogTypeSheetEdited,
		Content:   sheet.Title,
		IPAddress: util.GetClientIP(ctx),
	})
	return sheet, nil
}

func (s sheetServiceImpl) Preview(ctx context.Context, sheetID int32) (string, error) {
	post, err := s.GetByPostID(ctx, sheetID)
	if err != nil {
		return "", err
	}
	token := util.GenUUIDWithOutDash()
	s.Cache.Set(token, token, time.Minute*10)

	previewURL := strings.Builder{}

	isEnabledAbsolutePath, err := s.OptionService.IsEnabledAbsolutePath(ctx)
	if err != nil {
		return "", err
	}
	if isEnabledAbsolutePath {
		blogBaseURL, err := s.OptionService.GetBlogBaseURL(ctx)
		if err != nil {
			return "", err
		}
		previewURL.WriteString(blogBaseURL)
	}
	post.Slug = url.QueryEscape(post.Slug)
	fullPath, err := s.BuildFullPath(ctx, post)
	if err != nil {
		return "", err
	}
	previewURL.WriteString("/")
	previewURL.WriteString(fullPath)
	previewURL.WriteString("?token=")
	previewURL.WriteString(token)
	return previewURL.String(), nil
}

func (s sheetServiceImpl) CountVisit(ctx context.Context) (int64, error) {
	var count float64
	sheetDAL := dal.GetQueryByCtx(ctx).Post
	err := sheetDAL.WithContext(ctx).Select(sheetDAL.Visits.Sum().IfNull(0)).Where(sheetDAL.Type.Eq(consts.PostTypeSheet), sheetDAL.Status.Eq(consts.PostStatusPublished)).Scan(&count)
	if err != nil {
		return 0, WrapDBErr(err)
	}
	return int64(count), nil
}

func (s sheetServiceImpl) CountLike(ctx context.Context) (int64, error) {
	var count float64
	sheetDAL := dal.GetQueryByCtx(ctx).Post
	err := sheetDAL.WithContext(ctx).Select(sheetDAL.Likes.Sum().IfNull(0)).Where(sheetDAL.Type.Eq(consts.PostTypeSheet), sheetDAL.Status.Eq(consts.PostStatusPublished)).Scan(&count)
	if err != nil {
		return 0, WrapDBErr(err)
	}
	return int64(count), nil
}

func (s sheetServiceImpl) ListIndependentSheets(ctx context.Context) ([]*dto.IndependentSheet, error) {
	isEnableAbsolutePath, err := s.OptionService.IsEnabledAbsolutePath(ctx)
	if err != nil {
		return nil, err
	}
	var sheetBaseURL string
	if isEnableAbsolutePath {
		sheetBaseURL, err = s.OptionService.GetBlogBaseURL(ctx)
		if err != nil {
			return nil, err
		}
		sheetBaseURL += "/"
	} else {
		sheetBaseURL = "/"
	}
	linkPrefix, err := s.OptionService.GetLinksPrefix(ctx)
	if err != nil {
		return nil, err
	}
	linkSheet := dto.IndependentSheet{
		ID:        1,
		Title:     "友情链接",
		FullPath:  sheetBaseURL + linkPrefix,
		RouteName: "LinkList",
		Available: true,
	}
	linkSheet.Available, _ = s.ThemeService.TemplateExist(ctx, "links.tmpl")

	photoPrefix, err := s.OptionService.GetPhotoPrefix(ctx)
	if err != nil {
		return nil, err
	}
	photoSheet := dto.IndependentSheet{
		ID:        2,
		Title:     "图库页面",
		FullPath:  sheetBaseURL + photoPrefix,
		RouteName: "PhotoList",
		Available: true,
	}
	photoSheet.Available, _ = s.ThemeService.TemplateExist(ctx, "photos.tmpl")

	journalPrefix, err := s.OptionService.GetJournalPrefix(ctx)
	if err != nil {
		return nil, err
	}
	journalSheet := dto.IndependentSheet{
		ID:        3,
		Title:     "日志页面",
		FullPath:  sheetBaseURL + journalPrefix,
		RouteName: "JournalList",
		Available: true,
	}
	journalSheet.Available, _ = s.ThemeService.TemplateExist(ctx, "journals.tmpl")
	return []*dto.IndependentSheet{&linkSheet, &photoSheet, &journalSheet}, nil
}
