package impl

import (
	"context"
	"database/sql/driver"
	"net/url"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"

	"github.com/go-sonic/sonic/cache"
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

type postServiceImpl struct {
	service.BasePostService
	CategoryService service.CategoryService
	OptionService   service.OptionService
	Event           event.Bus
	Cache           cache.Cache
}

func NewPostService(basePostService service.BasePostService,
	categoryService service.CategoryService,
	optionService service.OptionService,
	event event.Bus,
	cache cache.Cache,
) service.PostService {
	return &postServiceImpl{
		BasePostService: basePostService,
		CategoryService: categoryService,
		OptionService:   optionService,
		Event:           event,
		Cache:           cache,
	}
}

func (p postServiceImpl) Page(ctx context.Context, postQuery param.PostQuery) ([]*entity.Post, int64, error) {
	if postQuery.PageNum < 0 || postQuery.PageSize <= 0 {
		return nil, 0, xerr.BadParam.New("").WithStatus(xerr.StatusBadRequest).WithMsg("Paging parameter error")
	}
	postDAL := dal.GetQueryByCtx(ctx).Post
	postCategoryDAL := dal.GetQueryByCtx(ctx).PostCategory
	postDo := postDAL.WithContext(ctx).Where(postDAL.Type.Eq(consts.PostTypePost))
	err := BuildSort(postQuery.Sort, &postDAL, &postDo)
	if err != nil {
		return nil, 0, err
	}

	if postQuery.Keyword != nil {
		postDo.Where(postDAL.Title.Like("%" + *postQuery.Keyword + "%")).Or(postDAL.OriginalContent.Like("%" + *postQuery.Keyword + "%"))
	}

	if postQuery.WithPassword != nil && !*postQuery.WithPassword {
		postDo.Where(postDAL.Password.Neq(""))
	}
	if len(postQuery.Statuses) > 0 {
		statuesValue := make([]driver.Valuer, len(postQuery.Statuses))
		for i, status := range postQuery.Statuses {
			statuesValue[i] = driver.Valuer(status)
		}
		postDo = postDo.Where(postDAL.Status.In(statuesValue...))
	}
	if postQuery.CategoryID != nil {
		postDo.Join(&entity.PostCategory{}, postDAL.ID.EqCol(postCategoryDAL.PostID)).Where(postCategoryDAL.CategoryID.Eq(*postQuery.CategoryID))
	}

	posts, totalCount, err := postDo.FindByPage(postQuery.PageNum*postQuery.PageSize, postQuery.PageSize)
	if err != nil {
		return nil, 0, WrapDBErr(err)
	}
	return posts, totalCount, nil
}

func (p postServiceImpl) IncreaseLike(ctx context.Context, postID int32) error {
	postDAL := dal.GetQueryByCtx(ctx).Post
	info, err := postDAL.WithContext(ctx).Where(postDAL.ID.Eq(postID)).UpdateSimple(postDAL.Likes.Add(1))
	if err != nil {
		return WrapDBErr(err)
	}
	if info.RowsAffected != 1 {
		return xerr.NoType.New("increase post like failed postID=%v", postID).WithStatus(xerr.StatusBadRequest).WithMsg("failed to like post")
	}
	return nil
}

func (p postServiceImpl) Create(ctx context.Context, postParam *param.Post) (*entity.Post, error) {
	post, err := p.ConvertParam(ctx, postParam)
	if err != nil {
		return nil, err
	}
	needEncrypt, err := p.CategoryService.IsCategoriesEncrypt(ctx, postParam.CategoryIDs...)
	if err != nil {
		return nil, nil
	}
	if post.Status != consts.PostStatusDraft && (post.Password != "" || needEncrypt) {
		post.Status = consts.PostStatusIntimate
	}

	post, err = p.CreateOrUpdate(ctx, post, postParam.CategoryIDs, postParam.TagIDs, postParam.MetaParam)
	if err != nil {
		return nil, err
	}
	// Todo delete authorization
	p.Event.Publish(ctx, &event.LogEvent{
		LogKey:    strconv.Itoa(int(post.ID)),
		LogType:   consts.LogTypePostPublished,
		Content:   post.Title,
		IpAddress: util.GetClientIP(ctx),
	})
	return post, nil
}

func (p postServiceImpl) Update(ctx context.Context, postID int32, postParam *param.Post) (*entity.Post, error) {
	postDAL := dal.GetQueryByCtx(ctx).Post
	post, err := postDAL.WithContext(ctx).Where(postDAL.ID.Eq(postID)).First()
	if err != nil {
		return nil, WrapDBErr(err)
	}
	postToUpdate, err := p.ConvertParam(ctx, postParam)
	if err != nil {
		return nil, err
	}
	if postToUpdate.CreateTime == (time.Time{}) {
		postToUpdate.CreateTime = post.CreateTime
	}
	postToUpdate.ID = post.ID
	post, err = p.CreateOrUpdate(ctx, postToUpdate, postParam.CategoryIDs, postParam.TagIDs, postParam.MetaParam)
	if err != nil {
		return nil, err
	}
	// TODO should use transcation
	p.Event.Publish(ctx, &event.PostUpdateEvent{
		PostID: post.ID,
	})
	p.Event.Publish(ctx, &event.LogEvent{
		LogKey:    strconv.Itoa(int(post.ID)),
		LogType:   consts.LogTypePostEdited,
		Content:   post.Title,
		IpAddress: util.GetClientIP(ctx),
	})
	return post, nil
}

func (p postServiceImpl) ConvertParam(ctx context.Context, postParam *param.Post) (*entity.Post, error) {
	post := &entity.Post{
		Type:            consts.PostTypePost,
		DisallowComment: postParam.DisallowComment,
		OriginalContent: postParam.OriginalContent,
		Password:        postParam.Password,
		MetaDescription: postParam.MetaDescription,
		MetaKeywords:    postParam.MetaKeywords,
		Template:        postParam.Template,
		Thumbnail:       postParam.Thumbnail,
		Title:           postParam.Title,
		TopPriority:     postParam.TopPriority,
		Status:          postParam.Status,
		EditTime:        util.TimePtr(time.Now()),
		Summary:         postParam.Summary,
		FormatContent:   postParam.Content,
	}
	if postParam.EditorType != nil {
		post.EditorType = *postParam.EditorType
	} else {
		post.EditorType = consts.EditorTypeMarkdown
	}
	if postParam.EditTime != nil {
		post.EditTime = util.TimePtr(time.UnixMilli(*postParam.EditTime))
	}
	if postParam.UpdateTime != nil {
		post.UpdateTime = util.TimePtr(time.UnixMilli(*postParam.UpdateTime))
	}

	post.WordCount = util.HtmlFormatWordCount(post.FormatContent)
	if postParam.Slug == "" {
		post.Slug = util.Slug(postParam.Title)
	} else {
		post.Slug = util.Slug(postParam.Slug)
	}
	if postParam.CreateTime != nil {
		post.CreateTime = time.UnixMilli(*postParam.CreateTime)
	} else {
		post.CreateTime = time.Now()
	}
	return post, nil
}

func (p postServiceImpl) CountByStatus(ctx context.Context, status consts.PostStatus) (int64, error) {
	postDAL := dal.GetQueryByCtx(ctx).Post
	count, err := postDAL.WithContext(ctx).Where(postDAL.Type.Eq(consts.PostTypePost), postDAL.Status.Eq(status)).Count()
	if err != nil {
		return 0, WrapDBErr(err)
	}
	return count, nil
}

func (p postServiceImpl) Preview(ctx context.Context, postID int32) (string, error) {
	post, err := p.GetByPostID(ctx, postID)
	if err != nil {
		return "", err
	}
	token := util.GenUUIDWithOutDash()
	p.Cache.Set(token, token, time.Minute*10)

	previewURL := strings.Builder{}

	isEnabledAbsolutePath, err := p.OptionService.IsEnabledAbsolutePath(ctx)
	if err != nil {
		return "", err
	}
	if isEnabledAbsolutePath {
		blogBaseURL, err := p.OptionService.GetBlogBaseURL(ctx)
		if err != nil {
			return "", err
		}
		previewURL.WriteString(blogBaseURL)
	}
	post.Slug = url.QueryEscape(post.Slug)
	fullPath, err := p.BuildFullPath(ctx, post)
	if err != nil {
		return "", err
	}
	previewURL.WriteString("/")
	previewURL.WriteString(fullPath)
	previewURL.WriteString("?token=")
	previewURL.WriteString(token)
	return previewURL.String(), nil
}

func (p postServiceImpl) CountVisit(ctx context.Context) (int64, error) {
	var count float64
	postDAL := dal.GetQueryByCtx(ctx).Post
	err := postDAL.WithContext(ctx).Select(postDAL.Visits.Sum().IfNull(0)).Where(postDAL.Type.Eq(consts.PostTypePost), postDAL.Status.Eq(consts.PostStatusPublished)).Scan(&count)
	if err != nil {
		return 0, WrapDBErr(err)
	}
	return int64(count), nil
}

func (p postServiceImpl) CountLike(ctx context.Context) (int64, error) {
	var count float64
	postDAL := dal.GetQueryByCtx(ctx).Post
	err := postDAL.WithContext(ctx).Select(postDAL.Likes.Sum().IfNull(0)).Where(postDAL.Type.Eq(consts.PostTypePost), postDAL.Status.Eq(consts.PostStatusPublished)).Scan(&count)
	if err != nil {
		return 0, WrapDBErr(err)
	}
	return int64(count), nil
}

func (p postServiceImpl) GetPrevPosts(ctx context.Context, post *entity.Post, size int) ([]*entity.Post, error) {
	postSort := p.OptionService.GetOrByDefault(ctx, property.IndexSort)
	postDAL := dal.GetQueryByCtx(ctx).Post
	postDO := postDAL.WithContext(ctx).Where(postDAL.Status.Eq(consts.PostStatusPublished))

	if postSort == "createTime" {
		postDO = postDO.Where(postDAL.CreateTime.Gt(post.CreateTime)).Order(postDAL.CreateTime)
	} else if postSort == "editTime" {
		var editTime time.Time
		if post.EditTime == nil {
			editTime = post.CreateTime
		} else {
			editTime = *post.EditTime
		}
		postDO = postDO.Where(postDAL.EditTime.Gt(editTime)).Order(postDAL.EditTime)
	} else if postSort == "visits" {
		postDO = postDO.Where(postDAL.Visits.Gt(post.Visits)).Order(postDAL.EditTime)
	} else {
		return nil, nil
	}

	posts, err := postDO.Find()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, WrapDBErr(err)
	}
	return posts, nil
}

func (p postServiceImpl) GetNextPosts(ctx context.Context, post *entity.Post, size int) ([]*entity.Post, error) {
	postSort := p.OptionService.GetOrByDefault(ctx, property.IndexSort)
	postDAL := dal.GetQueryByCtx(ctx).Post
	postDO := postDAL.WithContext(ctx).Where(postDAL.Status.Eq(consts.PostStatusPublished))

	if postSort == "createTime" {
		postDO = postDO.Where(postDAL.CreateTime.Lt(post.CreateTime)).Order(postDAL.CreateTime.Desc())
	} else if postSort == "editTime" {
		var editTime time.Time
		if post.EditTime == nil {
			editTime = post.CreateTime
		} else {
			editTime = *post.EditTime
		}
		postDO = postDO.Where(postDAL.EditTime.Lt(editTime)).Order(postDAL.EditTime.Desc())
	} else if postSort == "visits" {
		postDO = postDO.Where(postDAL.Visits.Lt(post.Visits)).Order(postDAL.EditTime.Desc())
	} else {
		return nil, nil
	}

	posts, err := postDO.Find()
	if err != nil {
		return nil, WrapDBErr(err)
	}
	return posts, nil
}
