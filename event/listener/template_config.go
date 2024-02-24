package listener

import (
	"context"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/go-sonic/sonic/config"
	"github.com/go-sonic/sonic/consts"
	"github.com/go-sonic/sonic/dal"
	"github.com/go-sonic/sonic/event"
	"github.com/go-sonic/sonic/handler"
	"github.com/go-sonic/sonic/model/property"
	"github.com/go-sonic/sonic/service"
	"github.com/go-sonic/sonic/template"
)

type TemplateConfigListener struct {
	Template      *template.Template
	ThemeService  service.ThemeService
	OptionService service.ClientOptionService
	UserService   service.UserService
	Logger        *zap.Logger
	Config        *config.Config
	Router        *gin.Engine
}

func NewTemplateConfigListener(bus event.Bus,
	template *template.Template,
	themeService service.ThemeService,
	optionService service.ClientOptionService,
	logger *zap.Logger,
	userService service.UserService,
	config *config.Config, server *handler.Server,
) {
	t := &TemplateConfigListener{
		Template:      template,
		ThemeService:  themeService,
		OptionService: optionService,
		Logger:        logger,
		UserService:   userService,
		Config:        config,
		Router:        server.Router,
	}
	bus.Subscribe(event.ThemeUpdateEventName, t.HandleThemeUpdateEvent)
	bus.Subscribe(event.UserUpdateEventName, t.HandleUserUpdateEvent)
	bus.Subscribe(event.OptionUpdateEventName, t.HandleOptionUpdateEvent)
	bus.Subscribe(event.StartEventName, t.HandleStartEvent)
	bus.Subscribe(event.ThemeActivatedEventName, t.HandleThemeUpdateEvent)
	bus.Subscribe(event.ThemeFileUpdatedEventName, t.HandleThemeFileUpdateEvent)
}

func (t *TemplateConfigListener) HandleThemeUpdateEvent(ctx context.Context, themeUpdateEvent event.Event) error {
	err := t.loadThemeConfig(ctx)
	if err != nil {
		return err
	}
	return t.loadThemeTemplate(ctx)
}

func (t *TemplateConfigListener) HandleUserUpdateEvent(ctx context.Context, userUpdateEvent event.Event) error {
	return t.loadUser(ctx)
}

func (t *TemplateConfigListener) HandleOptionUpdateEvent(ctx context.Context, optionUpdateEvent event.Event) error {
	err := t.loadThemeConfig(ctx)
	if err != nil {
		return err
	}
	return t.loadOption(ctx)
}

func (t *TemplateConfigListener) HandleStartEvent(ctx context.Context, startEvent event.Event) error {
	ctx = dal.SetCtxQuery(ctx, dal.GetQueryByCtx(ctx).ReplaceDB(dal.GetDB().Session(
		&gorm.Session{Logger: dal.DB.Logger.LogMode(logger.Warn)},
	)))
	err := t.loadThemeConfig(ctx)
	if err != nil {
		return err
	}
	err = t.loadUser(ctx)
	if err != nil {
		return err
	}
	err = t.loadOption(ctx)
	if err != nil {
		return err
	}
	return t.loadThemeTemplate(ctx)
}

func (t *TemplateConfigListener) HandleThemeFileUpdateEvent(ctx context.Context, themeFileUpdateEvent event.Event) error {
	return t.loadThemeTemplate(ctx)
}

func (t *TemplateConfigListener) loadThemeTemplate(ctx context.Context) error {
	theme, err := t.ThemeService.GetActivateTheme(ctx)
	if err != nil {
		return err
	}
	err = t.Template.Load([]string{filepath.Join(t.Config.Sonic.TemplateDir, "common"), theme.ThemePath})
	return err
}

func (t *TemplateConfigListener) loadThemeConfig(ctx context.Context) error {
	theme, err := t.ThemeService.GetActivateTheme(ctx)
	if err != nil {
		return nil
	}
	isEnabledAbsolutePath, err := t.OptionService.IsEnabledAbsolutePath(ctx)
	if err != nil {
		return err
	}
	blogBaseURL, err := t.OptionService.GetBlogBaseURL(ctx)
	if err != nil {
		return err
	}
	var themeBasePath string
	if isEnabledAbsolutePath {
		themeBasePath = blogBaseURL + "/themes/" + theme.FolderName
	} else {
		themeBasePath = "/themes/" + theme.FolderName
	}
	themeSetting, err := t.ThemeService.GetThemeSettingMap(ctx, theme.ID)
	if err != nil {
		return err
	}

	t.Template.SetSharedVariable("theme_base", themeBasePath)
	t.Template.SetSharedVariable("theme", theme)
	t.Template.SetSharedVariable("settings", themeSetting)
	t.Logger.Debug("load theme success", zap.String("theme", theme.Name))
	return nil
}

func (t *TemplateConfigListener) loadUser(ctx context.Context) error {
	users, err := t.UserService.GetAllUser(ctx)
	if err != nil {
		return err
	}
	if len(users) == 0 {
		return nil
	}
	user := users[0]
	user.Password = ""
	user.MfaKey = ""
	user.MfaType = consts.MFANone
	t.Template.SetSharedVariable("user", user)
	t.Logger.Debug("load user success", zap.Any("user", user))
	return nil
}

func (t *TemplateConfigListener) loadOption(ctx context.Context) error {
	options, err := t.OptionService.ListAllOption(ctx)
	if err != nil {
		return err
	}
	optionMap := make(map[string]interface{})
	for _, option := range options {
		optionMap[option.Key] = option.Value
	}
	blogBaseURL := t.OptionService.GetOrByDefault(ctx, property.BlogURL)
	blogTitle := t.OptionService.GetOrByDefault(ctx, property.BlogTitle)
	blogLogo := t.OptionService.GetOrByDefault(ctx, property.BlogLogo)
	globalAbsolutePathEnabled := t.OptionService.GetOrByDefault(ctx, property.GlobalAbsolutePathEnabled)
	seoKeywords := t.OptionService.GetOrByDefault(ctx, property.SeoKeywords)
	seoDescription := t.OptionService.GetOrByDefault(ctx, property.SeoDescription)
	journalPrefix := t.OptionService.GetOrByDefault(ctx, property.JournalsPrefix)
	archivePrefix := t.OptionService.GetOrByDefault(ctx, property.ArchivesPrefix)
	categoryPrefix := t.OptionService.GetOrByDefault(ctx, property.CategoriesPrefix)
	tagPrefix := t.OptionService.GetOrByDefault(ctx, property.TagsPrefix)
	linkPrefix := t.OptionService.GetOrByDefault(ctx, property.LinksPrefix)
	photoPrefix := t.OptionService.GetOrByDefault(ctx, property.PhotosPrefix)
	urlContext := "/"
	if globalAbsolutePathEnabled.(bool) {
		urlContext = blogBaseURL.(string) + "/"
	}
	t.Template.SetSharedVariable("version", consts.SonicVersion)
	t.Template.SetSharedVariable("options", optionMap)
	t.Template.SetSharedVariable("context", urlContext)
	t.Template.SetSharedVariable("globalAbsolutePathEnabled", globalAbsolutePathEnabled.(bool))
	t.Template.SetSharedVariable("blog_title", blogTitle)
	t.Template.SetSharedVariable("blog_logo", blogLogo)
	t.Template.SetSharedVariable("blog_url", blogBaseURL)
	t.Template.SetSharedVariable("seo_keywords", seoKeywords)
	t.Template.SetSharedVariable("seo_description", seoDescription)
	t.Template.SetSharedVariable("rss_url", blogBaseURL.(string)+"/rss.xml")
	t.Template.SetSharedVariable("atom_url", blogBaseURL.(string)+"/atom.xml")
	t.Template.SetSharedVariable("sitemap_xml_url", blogBaseURL.(string)+"/sitemap.xml")
	t.Template.SetSharedVariable("sitemap_html_url", blogBaseURL.(string)+"/sitemap.html")
	t.Template.SetSharedVariable("links_url", urlContext+linkPrefix.(string))
	t.Template.SetSharedVariable("photos_url", urlContext+photoPrefix.(string))
	t.Template.SetSharedVariable("journals_url", urlContext+journalPrefix.(string))
	t.Template.SetSharedVariable("archives_url", urlContext+archivePrefix.(string))
	t.Template.SetSharedVariable("categories_url", urlContext+categoryPrefix.(string))
	t.Template.SetSharedVariable("tags_url", urlContext+tagPrefix.(string))
	return nil
}
