package impl

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"strconv"
	"strings"
	"time"

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

type installServiceImpl struct {
	OptionService      service.OptionService
	Event              event.Bus
	UerService         service.UserService
	CategoryService    service.CategoryService
	PostService        service.PostService
	PostCommentService service.PostCommentService
	SheetService       service.SheetService
	MenuService        service.MenuService
}

func NewInstallService(
	optionService service.OptionService,
	event event.Bus,
	uerService service.UserService,
	categoryService service.CategoryService,
	postService service.PostService,
	postCommentService service.PostCommentService,
	sheetService service.SheetService,
	menuService service.MenuService,
) service.InstallService {
	return &installServiceImpl{
		OptionService:      optionService,
		Event:              event,
		UerService:         uerService,
		CategoryService:    categoryService,
		PostService:        postService,
		PostCommentService: postCommentService,
		SheetService:       sheetService,
		MenuService:        menuService,
	}
}

func (i installServiceImpl) InstallBlog(ctx context.Context, installParam param.Install) error {
	isInstalled, err := i.OptionService.GetOrByDefaultWithErr(ctx, property.IsInstalled, false)
	if err != nil {
		return nil
	}
	if isInstalled.(bool) {
		return xerr.BadParam.New("").WithStatus(xerr.StatusBadRequest).WithMsg("Blog has been installed")
	}
	err = dal.Transaction(ctx, func(txCtx context.Context) error {
		if err := i.createJWTSecret(txCtx); err != nil {
			return err
		}
		if err := i.createDefaultSetting(txCtx, installParam); err != nil {
			return err
		}
		user, err := i.createUser(txCtx, installParam.User)
		if err != nil {
			return err
		}
		category, err := i.createDefaultCategory(txCtx)
		if err != nil {
			return err
		}
		post, err := i.createDefaultPost(txCtx, category)
		if err != nil {
			return err
		}
		_, err = i.createDefaultSheet(txCtx)
		if err != nil {
			return err
		}
		_, err = i.createDefaultComment(txCtx, post)
		if err != nil {
			return err
		}
		err = i.createDefaultMenu(txCtx)
		if err != nil {
			return err
		}
		i.Event.Publish(txCtx, &event.LogEvent{
			LogKey:  strconv.Itoa(int(user.ID)),
			LogType: consts.LogTypeBlogInitialized,
			Content: "博客已成功初始化",
		})
		return nil
	})

	return err
}

func (i installServiceImpl) createDefaultSetting(ctx context.Context, installParam param.Install) error {
	optionMap := make(map[string]string)
	optionMap[property.IsInstalled.KeyValue] = "true"
	optionMap[property.GlobalAbsolutePathEnabled.KeyValue] = "false"
	optionMap[property.BlogTitle.KeyValue] = installParam.Title
	if installParam.Url == "" {
		blogURL, err := i.OptionService.GetBlogBaseURL(ctx)
		if err != nil {
			return err
		}
		optionMap[property.BlogUrl.KeyValue] = blogURL
	} else {
		optionMap[property.BlogUrl.KeyValue] = installParam.Url
	}
	if installParam.Locale == "" {
		optionMap[property.BlogLocale.KeyValue] = property.BlogLocale.DefaultValue.(string)
	}
	optionMap[property.BirthDay.KeyValue] = strconv.FormatInt(time.Now().UnixMilli(), 10)
	err := i.OptionService.Save(ctx, optionMap)
	return err
}

func (i installServiceImpl) createUser(ctx context.Context, user param.User) (*entity.User, error) {
	emailMd5 := md5.Sum([]byte(user.Email))
	avatar := "//cn.gravatar.com/avatar/" + hex.EncodeToString(emailMd5[:]) + "?s=256&d=mm"
	user.Avatar = avatar
	userEntity, err := i.UerService.CreateByParam(ctx, user)
	return userEntity, err
}

func (i installServiceImpl) createDefaultCategory(ctx context.Context) (*entity.Category, error) {
	categoryDal := dal.Use(dal.GetDBByCtx(ctx)).Category
	count, err := categoryDal.WithContext(ctx).Count()
	if err != nil {
		return nil, WrapDBErr(err)
	}
	if count > 0 {
		return nil, nil
	}
	categoryParam := param.Category{
		Name:        "默认分类",
		Slug:        "default",
		Description: "这是你的默认分类，如不需要，删除即可",
	}
	category, err := i.CategoryService.Create(ctx, &categoryParam)
	if err != nil {
		return nil, err
	}
	return category, nil
}

func (i installServiceImpl) createDefaultPost(ctx context.Context, category *entity.Category) (*entity.Post, error) {
	if category == nil {
		return nil, nil
	}
	postDAL := dal.Use(dal.GetDBByCtx(ctx)).Post
	count, err := postDAL.WithContext(ctx).Where(postDAL.Status.Eq(consts.PostStatusPublished)).Count()
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, nil
	}
	content := `
## Hello Sonic

如果你看到了这一篇文章，那么证明你已经安装成功了，感谢使用 [Sonic](https://go-sonic.org) 进行创作，希望能够使用愉快。

## 相关链接

- 官网：[https://github.com/go-sonic](https://github.com/go-sonic)
- 主题仓库：[https://github.com/go-sonic/default-theme-anatole](https://github.com/go-sonic/default-theme-anatole)
- 开源地址：[https://github.com/go-sonic/sonic](https://github.com/go-sonic/sonic)

在使用过程中，有任何问题都可以通过以上链接找寻答案，或者联系我们。

> 这是一篇自动生成的文章，请删除这篇文章之后开始你的创作吧！
`
	formatContent := `<h2 id="hello-sonic" tabindex="-1">Hello Sonic</h2>
	<p>如果你看到了这一篇文章，那么证明你已经安装成功了，感谢使用 <a href="https://go-sonic.org" target="_blank">Sonic</a> 进行创作，希望能够使用愉快。</p>
	<h2 id="%E7%9B%B8%E5%85%B3%E9%93%BE%E6%8E%A5" tabindex="-1">相关链接</h2>
	<ul>
	<li>官网：<a href="https://github.com/go-sonic" target="_blank">https://github.com/go-sonic</a></li>
	<li>主题仓库：<a href="https://github.com/go-sonic/default-theme-anatole" target="_blank">https://github.com/go-sonic/default-theme-anatole</a></li>
	<li>开源地址：<a href="https://github.com/go-sonic/sonic" target="_blank">https://github.com/go-sonic/sonic</a></li>
	</ul>
	<p>在使用过程中，有任何问题都可以通过以上链接找寻答案，或者联系我们。</p>
	<blockquote>
	<p>这是一篇自动生成的文章，请删除这篇文章之后开始你的创作吧！</p>
	</blockquote>
	`
	postParam := param.Post{
		Title:           "Hello Sonic",
		Status:          consts.PostStatusPublished,
		Slug:            "hello-sonic",
		OriginalContent: content,
		Content:         formatContent,
		CategoryIDs:     []int32{category.ID},
	}
	return i.PostService.Create(ctx, &postParam)
}

func (i installServiceImpl) createDefaultSheet(ctx context.Context) (*entity.Post, error) {
	postDAL := dal.Use(dal.GetDBByCtx(ctx)).Post
	count, err := postDAL.WithContext(ctx).Where(postDAL.Status.Eq(consts.PostStatusPublished)).Count()
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, nil
	}
	content := "## 关于页面 \n\n" +
		" 这是一个自定义页面，你可以在后台的 `页面` -> `所有页面` -> `自定义页面` 找到它，" +
		"你可以用于新建关于页面、留言板页面等等。发挥你自己的想象力！\n\n" +
		"> 这是一篇自动生成的页面，你可以在后台删除它。"
	sheetParam := param.Sheet{
		Title:           "关于页面",
		Status:          consts.PostStatusPublished,
		Slug:            "about",
		OriginalContent: content,
	}
	return i.SheetService.Create(ctx, &sheetParam)
}

func (i installServiceImpl) createDefaultComment(ctx context.Context, post *entity.Post) (*entity.Comment, error) {
	if post == nil {
		return nil, nil
	}

	count, err := i.PostCommentService.CountByStatus(ctx, consts.CommentStatusPublished)
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, nil
	}
	content := "欢迎使用 Sonic，这是你的第一条评论，头像来自 [Gravatar](https://cn.gravatar.com)，" +
		"你也可以通过注册 [Gravatar]" +
		"(https://cn.gravatar.com) 来显示自己的头像。"
	comment := &entity.Comment{
		Type:              consts.CommentTypePost,
		AllowNotification: true,
		Author:            "Sonic",
		AuthorURL:         "https://sonic.run",
		Content:           content,
		Email:             "hi@sonic.run",
		ParentID:          0,
		PostID:            post.ID,
		Status:            consts.CommentStatusPublished,
	}
	return i.PostCommentService.Create(ctx, comment)
}

func (i installServiceImpl) createDefaultMenu(ctx context.Context) error {
	menuIndex := &param.Menu{
		Name:     "首页",
		URL:      "/",
		Priority: 1,
	}
	menuArchive := &param.Menu{
		Name:     "文章归档",
		URL:      "/archives",
		Priority: 2,
	}
	menuCategory := &param.Menu{
		Name:     "默认分类",
		URL:      "/categories/default",
		Priority: 3,
	}
	menuSheet := &param.Menu{
		Name:     "关于页面",
		URL:      "/s/about",
		Priority: 4,
	}
	createMenu := func(menu *param.Menu, err error) error {
		if err != nil {
			return err
		}
		_, err = i.MenuService.Create(ctx, menu)
		return err
	}
	err := createMenu(menuIndex, nil)
	err = createMenu(menuArchive, err)
	err = createMenu(menuCategory, err)
	err = createMenu(menuSheet, err)
	return err
}

func (i installServiceImpl) createJWTSecret(ctx context.Context) error {
	secret := &strings.Builder{}
	secret.Grow(256)
	for i := 0; i < 8; i++ {
		secret.WriteString(util.GenUUIDWithOutDash())
	}
	m := map[string]string{property.JWTSecret.KeyValue: secret.String()}
	return i.OptionService.Save(ctx, m)
}
