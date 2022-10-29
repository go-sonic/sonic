package impl

import (
	"context"
	"database/sql/driver"
	"regexp"
	"strconv"
	"strings"
	"time"

	"gorm.io/gen/field"

	"github.com/go-sonic/sonic/consts"
	"github.com/go-sonic/sonic/dal"
	"github.com/go-sonic/sonic/log"
	"github.com/go-sonic/sonic/model/entity"
	"github.com/go-sonic/sonic/model/param"
	"github.com/go-sonic/sonic/model/property"
	"github.com/go-sonic/sonic/service"
	"github.com/go-sonic/sonic/util"
	"github.com/go-sonic/sonic/util/xerr"
)

type basePostServiceImpl struct {
	OptionService      service.OptionService
	BaseCommentService service.BaseCommentService
	CounterCache       *util.CounterCache[int32]
}

func NewBasePostService(optionService service.OptionService, baseCommentService service.BaseCommentService) service.BasePostService {
	counterCache := util.NewCounterCache(time.Second*5, nil, func(postID int32, count int64) {
		ctx := context.Background()
		postDAL := dal.Use(dal.GetDBByCtx(ctx)).Post
		_, err := postDAL.WithContext(ctx).Where(postDAL.ID.Eq(postID)).UpdateSimple(postDAL.Visits.Add(count))
		if err != nil {
			log.CtxErrorf(ctx, "increase visit err postID=%v", postID)
		}
	})
	b := &basePostServiceImpl{
		CounterCache:       counterCache,
		OptionService:      optionService,
		BaseCommentService: baseCommentService,
	}
	return b
}

func (b basePostServiceImpl) GetByStatus(ctx context.Context, status []consts.PostStatus, postType consts.PostType, sort *param.Sort) ([]*entity.Post, error) {
	postDAL := dal.Use(dal.GetDBByCtx(ctx)).Post
	postDo := postDAL.WithContext(ctx)
	if err := BuildSort(sort, &postDAL, &postDo); err != nil {
		return nil, err
	}
	statusAdapt := make([]driver.Valuer, len(status))
	for i, status := range status {
		statusAdapt[i] = status
	}
	posts, err := postDAL.WithContext(ctx).Where(postDAL.Status.In(statusAdapt...), postDAL.Type.Eq(postType)).Find()
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (b basePostServiceImpl) BuildFullPath(ctx context.Context, post *entity.Post) (string, error) {
	if post.Type == consts.PostTypePost {
		return b.buildPostFullPath(ctx, post)
	}
	return b.buildSheetFullPath(ctx, post)
}

func (b basePostServiceImpl) GetByPostIDs(ctx context.Context, postIDs []int32) (map[int32]*entity.Post, error) {
	postDAL := dal.Use(dal.GetDBByCtx(ctx)).Post
	posts, err := postDAL.WithContext(ctx).Where(postDAL.ID.In(postIDs...)).Find()
	if err != nil {
		return nil, WrapDBErr(err)
	}
	result := make(map[int32]*entity.Post)
	for _, post := range posts {
		result[post.ID] = post
	}
	return result, nil
}

func (b basePostServiceImpl) GetBySlug(ctx context.Context, slug string) (*entity.Post, error) {
	postDAL := dal.Use(dal.GetDBByCtx(ctx)).Post
	post, err := postDAL.WithContext(ctx).Where(postDAL.Slug.Eq(slug)).Take()
	if err != nil {
		return nil, WrapDBErr(err)
	}
	return post, nil
}

func (b basePostServiceImpl) buildPostFullPath(ctx context.Context, post *entity.Post) (string, error) {
	postPermaLinkType, err := b.OptionService.GetOrByDefaultWithErr(ctx, property.PostPermalinkType, property.PostPermalinkType.DefaultValue)
	if err != nil {
		return "", err
	}
	pathSuffix, err := b.OptionService.GetPathSuffix(ctx)
	if err != nil {
		return "", err
	}
	archivePrefix, err := b.OptionService.GetArchivePrefix(ctx)
	if err != nil {
		return "", err
	}
	month := post.CreateTime.Month()
	monthStr := util.IfElse(month < 10, "0"+strconv.Itoa(int(month)), strconv.Itoa(int(month))).(string)
	day := post.CreateTime.Day()
	dayStr := util.IfElse(month < 10, "0"+strconv.Itoa(day), strconv.Itoa(day)).(string)

	fullPath := strings.Builder{}
	isEnabled, err := b.OptionService.IsEnabledAbsolutePath(ctx)
	if err != nil {
		return "", err
	}
	if isEnabled {
		blogBaseUrl, err := b.OptionService.GetBlogBaseURL(ctx)
		if err != nil {
			return "", err
		}
		fullPath.WriteString(blogBaseUrl)
	}
	fullPath.WriteString("/")
	switch consts.PostPermalinkType(postPermaLinkType.(string)) {
	case consts.PostPermalinkTypeDefault:
		fullPath.WriteString(archivePrefix)
		fullPath.WriteString("/")
		fullPath.WriteString(post.Slug)
		fullPath.WriteString(pathSuffix)
	case consts.PostPermalinkTypeDate:
		fullPath.WriteString(strconv.Itoa(post.CreateTime.Year()))
		fullPath.WriteString("/")
		fullPath.WriteString(monthStr)
		fullPath.WriteString("/")
		fullPath.WriteString(post.Slug)
		fullPath.WriteString(pathSuffix)
	case consts.PostPermalinkTypeDay:
		fullPath.WriteString(strconv.Itoa(post.CreateTime.Year()))
		fullPath.WriteString("/")
		fullPath.WriteString(monthStr)
		fullPath.WriteString("/")
		fullPath.WriteString(dayStr)
		fullPath.WriteString("/")
		fullPath.WriteString(post.Slug)
		fullPath.WriteString(pathSuffix)
	case consts.PostPermalinkTypeYear:
		fullPath.WriteString(strconv.Itoa(post.CreateTime.Year()))
		fullPath.WriteString("/")
		fullPath.WriteString(post.Slug)
		fullPath.WriteString(pathSuffix)
	case consts.PostPermalinkTypeIDSlug:
		fullPath.WriteString(archivePrefix)
		fullPath.WriteString("/")
		fullPath.WriteString(strconv.Itoa(int(post.ID)))
		fullPath.WriteString(post.Slug)
	case consts.PostPermalinkTypeID:
		fullPath.WriteString("?p=")
		fullPath.WriteString(strconv.Itoa(int(post.ID)))
	}
	return fullPath.String(), nil
}

func (b basePostServiceImpl) buildSheetFullPath(ctx context.Context, sheet *entity.Post) (string, error) {
	sheetPermaLinkType, err := b.OptionService.GetOrByDefaultWithErr(ctx, property.SheetPermalinkType, property.SheetPermalinkType.DefaultValue)
	if err != nil {
		return "", err
	}
	pathSuffix, err := b.OptionService.GetPathSuffix(ctx)
	if err != nil {
		return "", err
	}
	sheetPrefix, err := b.OptionService.GetOrByDefaultWithErr(ctx, property.SheetPrefix, property.SheetPrefix.DefaultValue)
	if err != nil {
		return "", err
	}

	fullPath := strings.Builder{}
	isEnabled, err := b.OptionService.IsEnabledAbsolutePath(ctx)
	if err != nil {
		return "", err
	}
	if isEnabled {
		blogBaseUrl, err := b.OptionService.GetBlogBaseURL(ctx)
		if err != nil {
			return "", err
		}
		fullPath.WriteString(blogBaseUrl)
	}
	fullPath.WriteString("/")
	switch consts.SheetPermaLinkType(sheetPermaLinkType.(string)) {
	case consts.SheetPermaLinkTypeSecondary:
		fullPath.WriteString(sheetPrefix.(string))
		fullPath.WriteString("/")
		fullPath.WriteString(sheet.Slug)
		fullPath.WriteString(pathSuffix)
	case consts.SheetPermaLinkTypeRoot:
		fullPath.WriteString(sheet.Slug)
		fullPath.WriteString(pathSuffix)
	}
	return fullPath.String(), nil
}

func (b basePostServiceImpl) GetByPostID(ctx context.Context, postID int32) (*entity.Post, error) {
	postDAL := dal.Use(dal.GetDBByCtx(ctx)).Post
	post, err := postDAL.WithContext(ctx).Where(postDAL.ID.Eq(postID)).First()
	if err != nil {
		return nil, WrapDBErr(err)
	}
	return post, nil
}

var summaryPattern = regexp.MustCompile(`[\t\r\n]`)

func (b basePostServiceImpl) GenerateSummary(ctx context.Context, htmlContent string) string {
	text := util.CleanHtmlTag(htmlContent)
	text = summaryPattern.ReplaceAllString(text, "")
	summaryLength := b.OptionService.GetPostSummaryLength(ctx)
	end := summaryLength
	textRune := []rune(text)
	if len(textRune) < end {
		end = len(textRune)
	}
	return string(textRune[:end])
}

func (b basePostServiceImpl) Delete(ctx context.Context, postID int32) error {
	err := dal.Use(dal.GetDBByCtx(ctx)).Transaction(func(tx *dal.Query) error {
		postDAL := tx.Post
		postTagDAL := tx.PostTag
		postCategoryDAL := tx.PostCategory
		postMetaDAL := tx.Meta
		postCommentDAL := tx.Comment

		deleteResult, err := postDAL.WithContext(ctx).Where(postDAL.ID.Eq(postID)).Delete()
		if err != nil {
			return WrapDBErr(err)
		}
		if deleteResult.RowsAffected != 1 {
			return xerr.NoType.New("").WithMsg("delete post failed")
		}
		_, err = postTagDAL.WithContext(ctx).Where(postTagDAL.PostID.Eq(postID)).Delete()
		if err != nil {
			return WrapDBErr(err)
		}
		_, err = postCategoryDAL.WithContext(ctx).Where(postCategoryDAL.PostID.Eq(postID)).Delete()
		if err != nil {
			return WrapDBErr(err)
		}
		_, err = postMetaDAL.WithContext(ctx).Where(postMetaDAL.PostID.Eq(postID)).Delete()
		if err != nil {
			return WrapDBErr(err)
		}
		_, err = postCommentDAL.WithContext(ctx).Where(postCommentDAL.PostID.Eq(postID)).Delete()
		if err != nil {
			return WrapDBErr(err)
		}
		return nil
	})
	return err
}

func (b basePostServiceImpl) UpdateStatus(ctx context.Context, postID int32, status consts.PostStatus) (*entity.Post, error) {
	if postID < 0 || status < consts.PostStatusPublished || status > consts.PostStatusIntimate {
		return nil, xerr.BadParam.New("").WithMsg("postID or status parameter error").WithStatus(xerr.StatusBadRequest)
	}

	postDAL := dal.Use(dal.GetDBByCtx(ctx)).Post
	post, err := postDAL.WithContext(ctx).Where(postDAL.ID.Eq(postID)).First()
	if err != nil {
		return nil, WrapDBErr(err)
	}
	updateResult, err := postDAL.WithContext(ctx).Where(postDAL.ID.Eq(postID)).UpdateColumnSimple(postDAL.Status.Value(status))
	if err != nil {
		return nil, WrapDBErr(err)
	}
	if updateResult.RowsAffected != 1 {
		return nil, xerr.NoType.New("update post status failed postID=%v", postID).WithMsg("update post status failed")
	}
	post.Status = status
	return post, nil
}

func (b basePostServiceImpl) DeleteBatch(ctx context.Context, postIDs []int32) error {
	err := dal.Use(dal.GetDBByCtx(ctx)).Transaction(func(tx *dal.Query) error {
		postDAL := tx.Post
		postTagDAL := tx.PostTag
		postCategoryDAL := tx.PostCategory
		postMetaDAL := tx.Meta
		postCommentDAL := tx.Comment

		deleteResult, err := postDAL.WithContext(ctx).Where(postDAL.ID.In(postIDs...)).Delete()
		if err != nil {
			return WrapDBErr(err)
		}
		if deleteResult.RowsAffected != 1 {
			return xerr.NoType.New("").WithMsg("delete post failed")
		}
		_, err = postTagDAL.WithContext(ctx).Where(postTagDAL.PostID.In(postIDs...)).Delete()
		if err != nil {
			return WrapDBErr(err)
		}
		_, err = postCategoryDAL.WithContext(ctx).Where(postCategoryDAL.PostID.In(postIDs...)).Delete()
		if err != nil {
			return WrapDBErr(err)
		}
		_, err = postMetaDAL.WithContext(ctx).Where(postMetaDAL.PostID.In(postIDs...)).Delete()
		if err != nil {
			return WrapDBErr(err)
		}
		_, err = postCommentDAL.WithContext(ctx).Where(postCommentDAL.PostID.In(postIDs...)).Delete()
		if err != nil {
			return WrapDBErr(err)
		}
		return nil
	})
	return err
}

func (b basePostServiceImpl) CreateOrUpdate(ctx context.Context, post *entity.Post, categoryIDs, tagIDs []int32, metas []param.Meta) (*entity.Post, error) {
	err := dal.Use(dal.GetDBByCtx(ctx)).Transaction(func(tx *dal.Query) error {
		postDAL := tx.Post
		postCategoryDAL := tx.PostCategory
		postTagDAL := tx.PostTag
		categoryDAL := tx.Category
		tagDAL := tx.Tag
		postMetaDAL := tx.Meta

		// create post
		if post.ID == 0 {
			postCount, err := postDAL.WithContext(ctx).Select(field.Star).Omit(postDAL.UpdateTime).Where(postDAL.Slug.Eq(post.Slug)).Count()
			if err != nil {
				return WrapDBErr(err)
			}
			if postCount > 0 {
				return xerr.BadParam.New("").WithMsg("文章别名已存在(Article alias already exists)").WithStatus(xerr.StatusBadRequest)
			}
			err = postDAL.WithContext(ctx).Create(post)
			if err != nil {
				return WrapDBErr(err)
			}
		} else {
			// update post
			slugCount, err := postDAL.WithContext(ctx).Where(postDAL.Slug.Eq(post.Slug), postDAL.ID.Neq(post.ID)).Count()
			if err != nil {
				return WrapDBErr(err)
			}
			if slugCount > 0 {
				return xerr.BadParam.New("").WithMsg("文章别名已存在(Article alias already exists)").WithStatus(xerr.StatusBadRequest)
			}
			updateResult, err := postDAL.WithContext(ctx).Select(field.Star).Omit(postDAL.Likes, postDAL.Visits).Where(postDAL.ID.Eq(post.ID)).Updates(post)
			if err != nil {
				return WrapDBErr(err)
			}
			if updateResult.RowsAffected != 1 {
				return xerr.NoType.New("").WithMsg("更新文章失败(update post failed)")
			}
		}

		// create post_category
		if post.ID > 0 {
			_, err := postCategoryDAL.WithContext(ctx).Where(postCategoryDAL.PostID.Eq(post.ID)).Delete()
			if err != nil {
				return WrapDBErr(err)
			}
		}
		if len(categoryIDs) > 0 {
			categoryCount, err := categoryDAL.WithContext(ctx).Where(categoryDAL.ID.In(categoryIDs...)).Count()
			if err != nil {
				return WrapDBErr(err)
			}
			if int(categoryCount) != len(categoryIDs) {
				return xerr.BadParam.New("").WithMsg("category not exist").WithStatus(xerr.StatusBadRequest)
			}
			pcs := make([]*entity.PostCategory, 0, len(categoryIDs))
			for _, categoryID := range categoryIDs {
				pc := &entity.PostCategory{
					CategoryID: categoryID,
					PostID:     post.ID,
				}
				pcs = append(pcs, pc)
			}
			err = postCategoryDAL.WithContext(ctx).Create(pcs...)
			if err != nil {
				return WrapDBErr(err)
			}
		}
		// create post_tag
		if post.ID > 0 {
			_, err := postTagDAL.WithContext(ctx).Where(postTagDAL.PostID.Eq(post.ID)).Delete()
			if err != nil {
				return WrapDBErr(err)
			}
		}
		if len(tagIDs) > 0 {
			tagCount, err := tagDAL.WithContext(ctx).Where(tagDAL.ID.In(tagIDs...)).Count()
			if err != nil {
				return WrapDBErr(err)
			}
			if int(tagCount) != len(tagIDs) {
				return xerr.BadParam.New("").WithMsg("tag not exist").WithStatus(xerr.StatusBadRequest)
			}
			pts := make([]*entity.PostTag, 0, len(tagIDs))
			for _, tagID := range tagIDs {
				pts = append(pts, &entity.PostTag{
					PostID: post.ID,
					TagID:  tagID,
				})
			}
			err = postTagDAL.WithContext(ctx).Create(pts...)
			if err != nil {
				return err
			}
		}

		// create metas
		if post.ID > 0 {
			_, err := postMetaDAL.WithContext(ctx).Where(postMetaDAL.PostID.Eq(post.ID)).Delete()
			if err != nil {
				return WrapDBErr(err)
			}
		}
		if len(metas) > 0 {
			pms := make([]*entity.Meta, 0, len(metas))
			for _, meta := range metas {
				pms = append(pms, &entity.Meta{
					PostID:    post.ID,
					MetaValue: meta.Value,
					MetaKey:   meta.Key,
				})
			}
			err := postMetaDAL.WithContext(ctx).Create(pms...)
			if err != nil {
				return WrapDBErr(err)
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (b basePostServiceImpl) UpdateStatusBatch(ctx context.Context, status consts.PostStatus, postIDs []int32) ([]*entity.Post, error) {
	if status < consts.PostStatusPublished || status > consts.PostStatusIntimate {
		return nil, xerr.BadParam.New("").WithMsg("postID or status parameter error").WithStatus(xerr.StatusBadRequest)
	}

	uniquePostIDMap := make(map[int32]struct{})
	for _, postID := range postIDs {
		uniquePostIDMap[postID] = struct{}{}
	}
	uniqueIDs := make([]int32, 0)
	for postID := range uniquePostIDMap {
		uniqueIDs = append(uniqueIDs, postID)
	}
	err := dal.Use(dal.GetDBByCtx(ctx)).Transaction(func(tx *dal.Query) error {
		postDAL := tx.Post
		updateResult, err := postDAL.WithContext(ctx).Where(postDAL.ID.In(uniqueIDs...)).UpdateColumnSimple(postDAL.Status.Value(status))
		if err != nil {
			return WrapDBErr(err)
		}
		if updateResult.RowsAffected != int64(len(uniqueIDs)) {
			return xerr.NoType.New("").WithMsg("update post status failed")
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	postDAL := dal.Use(dal.GetDBByCtx(ctx)).Post
	posts, err := postDAL.WithContext(ctx).Where(postDAL.ID.In(uniqueIDs...)).Find()
	if err != nil {
		return nil, WrapDBErr(err)
	}
	return posts, nil
}

func (b basePostServiceImpl) UpdateDraftContent(ctx context.Context, postID int32, content string) (*entity.Post, error) {
	postDAL := dal.Use(dal.GetDBByCtx(ctx)).Post
	post, err := postDAL.WithContext(ctx).Where(postDAL.ID.Eq(postID)).First()
	if err != nil {
		return nil, WrapDBErr(err)
	}
	if post.OriginalContent != content {
		updateResult, err := postDAL.WithContext(ctx).Where(postDAL.ID.Eq(postID)).UpdateColumnSimple(postDAL.OriginalContent.Value(content))
		if err != nil {
			return nil, WrapDBErr(err)
		}
		if updateResult.RowsAffected != 1 {
			return nil, xerr.NoType.New("").WithMsg("update post content failed")
		}
	}
	post.OriginalContent = content
	return post, nil
}

func (b basePostServiceImpl) IncreaseVisit(ctx context.Context, postID int32) {
	b.CounterCache.IncrBy(postID, 1)
}
