package handler

import (
	"context"
	"path/filepath"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/go-sonic/sonic/config"
	"github.com/go-sonic/sonic/consts"
	"github.com/go-sonic/sonic/dal"
	"github.com/go-sonic/sonic/handler/middleware"
)

func (s *Server) RegisterRouters() {
	router := s.Router
	if config.IsDev() {
		router.Use(cors.New(cors.Config{
			AllowAllOrigins:  true,
			AllowOrigins:     []string{},
			AllowMethods:     []string{"PUT", "PATCH", "GET", "DELETE", "POST", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Admin-Authorization", "Content-Type"},
			AllowCredentials: true,
			ExposeHeaders:    []string{"Content-Length"},
		}))
	}

	{
		router.GET("/ping", func(ctx *gin.Context) {
			_, _ = ctx.Writer.Write([]byte("pong"))
		})
		{
			staticRouter := router.Group("/")
			staticRouter.StaticFS(s.Config.Sonic.AdminURLPath, gin.Dir(s.Config.Sonic.AdminResourcesDir, false))
			staticRouter.StaticFS("/css", gin.Dir(filepath.Join(s.Config.Sonic.AdminResourcesDir, "css"), false))
			staticRouter.StaticFS("/js", gin.Dir(filepath.Join(s.Config.Sonic.AdminResourcesDir, "js"), false))
			staticRouter.StaticFS("/images", gin.Dir(filepath.Join(s.Config.Sonic.AdminResourcesDir, "images"), false))
			staticRouter.Use(middleware.NewCacheControlMiddleware(middleware.WithMaxAge(time.Hour*24*7)).CacheControl()).
				StaticFS(consts.SonicUploadDir, gin.Dir(s.Config.Sonic.UploadDir, false))
			staticRouter.StaticFS("/themes/", gin.Dir(s.Config.Sonic.ThemeDir, false))
		}
		{
			adminAPIRouter := router.Group("/api/admin")
			adminAPIRouter.Use(s.LogMiddleware.LoggerWithConfig(middleware.GinLoggerConfig{}), s.RecoveryMiddleware.RecoveryWithLogger(), s.InstallRedirectMiddleware.InstallRedirect())
			adminAPIRouter.GET("/is_installed", s.wrapHandler(s.AdminHandler.IsInstalled))
			adminAPIRouter.POST("/login/precheck", s.wrapHandler(s.AdminHandler.AuthPreCheck))
			adminAPIRouter.POST("/login", s.wrapHandler(s.AdminHandler.Auth))
			adminAPIRouter.POST("/refresh/:refreshToken", s.wrapHandler(s.AdminHandler.RefreshToken))
			adminAPIRouter.POST("/installations", s.wrapHandler(s.InstallHandler.InstallBlog))
			{
				authRouter := adminAPIRouter.Group("")
				authRouter.Use(s.AuthMiddleware.GetWrapHandler())
				authRouter.POST("/logout", s.wrapHandler(s.AdminHandler.LogOut))
				authRouter.POST("/password/code", s.wrapHandler(s.AdminHandler.SendResetCode))
				authRouter.GET("/environments", s.wrapHandler(s.AdminHandler.GetEnvironments))
				authRouter.GET("/sonic/logfile", s.wrapHandler(s.AdminHandler.GetLogFiles))
				{
					attachmentRouter := authRouter.Group("/attachments")
					attachmentRouter.POST("/upload", s.wrapHandler(s.AttachmentHandler.UploadAttachment))
					attachmentRouter.POST("/uploads", s.wrapHandler(s.AttachmentHandler.UploadAttachments))
					attachmentRouter.DELETE("/:id", s.wrapHandler(s.AttachmentHandler.DeleteAttachment))
					attachmentRouter.DELETE("", s.wrapHandler(s.AttachmentHandler.DeleteAttachmentInBatch))
					attachmentRouter.GET("", s.wrapHandler(s.AttachmentHandler.QueryAttachment))
					attachmentRouter.GET("/:id", s.wrapHandler(s.AttachmentHandler.GetAttachmentByID))
					attachmentRouter.PUT("/:id", s.wrapHandler(s.AttachmentHandler.UpdateAttachment))
					attachmentRouter.GET("/media_types", s.wrapHandler(s.AttachmentHandler.GetAllMediaType))
					attachmentRouter.GET("types", s.wrapHandler(s.AttachmentHandler.GetAllTypes))
				}
				{
					backupRouter := authRouter.Group("/backups")
					backupRouter.POST("/work-dir", s.wrapHandler(s.BackupHandler.BackupWholeSite))
					backupRouter.GET("/work-dir", s.wrapHandler(s.BackupHandler.ListBackups))
					backupRouter.GET("/work-dir/*path", s.BackupHandler.HandleWorkDir)
					backupRouter.DELETE("/work-dir", s.wrapHandler(s.BackupHandler.DeleteBackups))
					backupRouter.POST("/data", s.wrapHandler(s.BackupHandler.ExportData))
					backupRouter.DELETE("/data", s.wrapHandler(s.BackupHandler.DeleteDataFile))
					backupRouter.GET("/data/*path", s.BackupHandler.HandleData)
					backupRouter.POST("/markdown/export", s.wrapHandler(s.BackupHandler.ExportMarkdown))
					backupRouter.POST("/markdown/import", s.wrapHandler(s.BackupHandler.ImportMarkdown))
					backupRouter.GET("/markdown/fetch", s.wrapHandler(s.BackupHandler.GetMarkDownBackup))
					backupRouter.GET("/markdown/export", s.wrapHandler(s.BackupHandler.ListMarkdowns))
					backupRouter.DELETE("/markdown/export", s.wrapHandler(s.BackupHandler.DeleteMarkdowns))
					backupRouter.GET("/markdown/export/:filename", s.BackupHandler.DownloadMarkdown)
				}
				{
					categoryRouter := authRouter.Group("/categories")
					categoryRouter.PUT("/batch", s.wrapHandler(s.CategoryHandler.UpdateCategoryBatch))
					categoryRouter.GET("/:categoryID", s.wrapHandler(s.CategoryHandler.GetCategoryByID))
					categoryRouter.GET("", s.wrapHandler(s.CategoryHandler.ListAllCategory))
					categoryRouter.GET("/tree_view", s.wrapHandler(s.CategoryHandler.ListAsTree))
					categoryRouter.POST("", s.wrapHandler(s.CategoryHandler.CreateCategory))
					categoryRouter.PUT("/:categoryID", s.wrapHandler(s.CategoryHandler.UpdateCategory))
					categoryRouter.DELETE("/:categoryID", s.wrapHandler(s.CategoryHandler.DeleteCategory))
				}
				{
					postRouter := authRouter.Group("/posts")
					postRouter.GET("", s.wrapHandler(s.PostHandler.ListPosts))
					postRouter.GET("/latest", s.wrapHandler(s.PostHandler.ListLatestPosts))
					postRouter.GET("/status/:status", s.wrapHandler(s.PostHandler.ListPostsByStatus))
					postRouter.GET("/:postID", s.wrapHandler(s.PostHandler.GetByPostID))
					postRouter.POST("", s.wrapHandler(s.PostHandler.CreatePost))
					postRouter.PUT("/:postID", s.wrapHandler(s.PostHandler.UpdatePost))
					postRouter.PUT("/:postID/status/:status", s.wrapHandler(s.PostHandler.UpdatePostStatus))
					postRouter.PUT("/status/:status", s.wrapHandler(s.PostHandler.UpdatePostStatusBatch))
					postRouter.PUT("/:postID/status/draft/content", s.wrapHandler(s.PostHandler.UpdatePostDraft))
					postRouter.DELETE("/:postID", s.wrapHandler(s.PostHandler.DeletePost))
					postRouter.DELETE("", s.wrapHandler(s.PostHandler.DeletePostBatch))
					postRouter.GET("/:postID/preview", s.PostHandler.PreviewPost)
					{
						postCommentRouter := postRouter.Group("/comments")
						postCommentRouter.GET("", s.wrapHandler(s.PostCommentHandler.ListPostComment))
						postCommentRouter.GET("/latest", s.wrapHandler(s.PostCommentHandler.ListPostCommentLatest))
						postCommentRouter.GET("/:postID/tree_view", s.wrapHandler(s.PostCommentHandler.ListPostCommentAsTree))
						postCommentRouter.GET("/:postID/list_view", s.wrapHandler(s.PostCommentHandler.ListPostCommentWithParent))
						postCommentRouter.POST("", s.wrapHandler(s.PostCommentHandler.CreatePostComment))
						postCommentRouter.PUT("/:commentID", s.wrapHandler(s.PostCommentHandler.UpdatePostComment))
						postCommentRouter.PUT("/:commentID/status/:status", s.wrapHandler(s.PostCommentHandler.UpdatePostCommentStatus))
						postCommentRouter.PUT("/status/:status", s.wrapHandler(s.PostCommentHandler.UpdatePostCommentStatusBatch))
						postCommentRouter.DELETE("/:commentID", s.wrapHandler(s.PostCommentHandler.DeletePostComment))
						postCommentRouter.DELETE("", s.wrapHandler(s.PostCommentHandler.DeletePostCommentBatch))
					}
				}
				{
					optionRouter := authRouter.Group("/options")
					optionRouter.GET("", s.wrapHandler(s.OptionHandler.ListAllOptions))
					optionRouter.GET("/map_view", s.wrapHandler(s.OptionHandler.ListAllOptionsAsMap))
					optionRouter.POST("/map_view/keys", s.wrapHandler(s.OptionHandler.ListAllOptionsAsMapWithKey))
					optionRouter.POST("/saving", s.wrapHandler(s.OptionHandler.SaveOption))
					optionRouter.POST("/map_view/saving", s.wrapHandler(s.OptionHandler.SaveOptionWithMap))
				}
				{
					logRouter := authRouter.Group("/logs")
					logRouter.GET("/latest", s.wrapHandler(s.LogHandler.PageLatestLog))
					logRouter.GET("", s.wrapHandler(s.LogHandler.PageLog))
					logRouter.GET("/clear", s.wrapHandler(s.LogHandler.ClearLog))
				}
				{
					statisticRouter := authRouter.Group("/statistics")
					statisticRouter.GET("", s.wrapHandler(s.StatisticHandler.Statistics))
					statisticRouter.GET("user", s.wrapHandler(s.StatisticHandler.StatisticsWithUser))
				}
				{
					sheetRouter := authRouter.Group("/sheets")
					sheetRouter.GET("/:sheetID", s.wrapHandler(s.SheetHandler.GetSheetByID))
					sheetRouter.GET("", s.wrapHandler(s.SheetHandler.ListSheet))
					sheetRouter.POST("", s.wrapHandler(s.SheetHandler.CreateSheet))
					sheetRouter.PUT("/:sheetID", s.wrapHandler(s.SheetHandler.UpdateSheet))
					sheetRouter.PUT("/:sheetID/:status", s.wrapHandler(s.SheetHandler.UpdateSheetStatus))
					sheetRouter.PUT("/:sheetID/status/draft/content", s.wrapHandler(s.SheetHandler.UpdateSheetDraft))
					sheetRouter.DELETE("/:sheetID", s.wrapHandler(s.SheetHandler.DeleteSheet))
					sheetRouter.GET("/preview/:sheetID", s.SheetHandler.PreviewSheet)
					sheetRouter.GET("/independent", s.wrapHandler(s.SheetHandler.IndependentSheets))
					{
						sheetCommentRouter := sheetRouter.Group("/comments")
						sheetCommentRouter.GET("", s.wrapHandler(s.SheetCommentHandler.ListSheetComment))
						sheetCommentRouter.GET("/latest", s.wrapHandler(s.SheetCommentHandler.ListSheetCommentLatest))
						sheetCommentRouter.GET("/:sheetID/tree_view", s.wrapHandler(s.SheetCommentHandler.ListSheetCommentAsTree))
						sheetCommentRouter.GET("/:sheetID/list_view", s.wrapHandler(s.SheetCommentHandler.ListSheetCommentWithParent))
						sheetCommentRouter.POST("/", s.wrapHandler(s.SheetCommentHandler.CreateSheetComment))
						sheetCommentRouter.PUT("/:commentID/status/:status", s.wrapHandler(s.SheetCommentHandler.UpdateSheetCommentStatus))
						sheetCommentRouter.PUT("/status/:status", s.wrapHandler(s.SheetCommentHandler.UpdateSheetCommentStatusBatch))
						sheetCommentRouter.DELETE("/:commentID", s.wrapHandler(s.SheetCommentHandler.DeleteSheetComment))
						sheetCommentRouter.DELETE("", s.wrapHandler(s.SheetCommentHandler.DeleteSheetCommentBatch))
					}
				}
				{
					journalRouter := authRouter.Group("/journals")
					journalRouter.GET("", s.wrapHandler(s.JournalHandler.ListJournal))
					journalRouter.GET("/latest", s.wrapHandler(s.JournalHandler.ListLatestJournal))
					journalRouter.POST("", s.wrapHandler(s.JournalHandler.CreateJournal))
					journalRouter.PUT("/:journalID", s.wrapHandler(s.JournalHandler.UpdateJournal))
					journalRouter.DELETE("/:journalID", s.wrapHandler(s.JournalHandler.DeleteJournal))
					{
						journalCommentRouter := journalRouter.Group("/comments")
						journalCommentRouter.GET("", s.wrapHandler(s.JournalCommentHandler.ListJournalComment))
						journalCommentRouter.GET("/latest", s.wrapHandler(s.JournalCommentHandler.ListJournalCommentLatest))
						journalCommentRouter.GET("/:journalID/tree_view", s.wrapHandler(s.JournalCommentHandler.ListJournalCommentAsTree))
						journalCommentRouter.GET("/:journalID/list_view", s.wrapHandler(s.JournalCommentHandler.ListJournalCommentWithParent))
						journalCommentRouter.POST("/", s.wrapHandler(s.JournalCommentHandler.CreateJournalComment))
						journalCommentRouter.PUT("/:commentID/status/:status", s.wrapHandler(s.JournalCommentHandler.UpdateJournalCommentStatus))
						journalCommentRouter.PUT("/status/:status", s.wrapHandler(s.JournalCommentHandler.UpdateJournalStatusBatch))
						journalCommentRouter.PUT("/:commentID", s.wrapHandler(s.JournalCommentHandler.UpdateJournalComment))
						journalCommentRouter.DELETE("/:commentID", s.wrapHandler(s.JournalCommentHandler.DeleteJournalComment))
						journalCommentRouter.DELETE("", s.wrapHandler(s.JournalCommentHandler.DeleteJournalCommentBatch))
					}
				}

				{
					linkRouter := authRouter.Group("/links")
					linkRouter.GET("", s.wrapHandler(s.LinkHandler.ListLinks))
					linkRouter.GET("/:id", s.wrapHandler(s.LinkHandler.GetLinkByID))
					linkRouter.POST("", s.wrapHandler(s.LinkHandler.CreateLink))
					linkRouter.PUT("/:id", s.wrapHandler(s.LinkHandler.UpdateLink))
					linkRouter.DELETE("/:id", s.wrapHandler(s.LinkHandler.DeleteLink))
					linkRouter.GET("/teams", s.wrapHandler(s.LinkHandler.ListLinkTeams))
				}
				{
					menuRouter := authRouter.Group("/menus")
					menuRouter.GET("", s.wrapHandler(s.MenuHandler.ListMenus))
					menuRouter.GET("/tree_view", s.wrapHandler(s.MenuHandler.ListMenusAsTree))
					menuRouter.GET("/team/tree_view", s.wrapHandler(s.MenuHandler.ListMenusAsTreeByTeam))
					menuRouter.GET("/:id", s.wrapHandler(s.MenuHandler.GetMenuByID))
					menuRouter.POST("", s.wrapHandler(s.MenuHandler.CreateMenu))
					menuRouter.POST("/batch", s.wrapHandler(s.MenuHandler.CreateMenuBatch))
					menuRouter.PUT("/:id", s.wrapHandler(s.MenuHandler.UpdateMenu))
					menuRouter.PUT("/batch", s.wrapHandler(s.MenuHandler.UpdateMenuBatch))
					menuRouter.DELETE("/:id", s.wrapHandler(s.MenuHandler.DeleteMenu))
					menuRouter.DELETE("/batch", s.wrapHandler(s.MenuHandler.DeleteMenuBatch))
					menuRouter.GET("/teams", s.wrapHandler(s.MenuHandler.ListMenuTeams))
				}
				{
					tagRouter := authRouter.Group("/tags")
					tagRouter.GET("", s.wrapHandler(s.TagHandler.ListTags))
					tagRouter.GET("/:id", s.wrapHandler(s.TagHandler.GetTagByID))
					tagRouter.POST("", s.wrapHandler(s.TagHandler.CreateTag))
					tagRouter.PUT("/:id", s.wrapHandler(s.TagHandler.UpdateTag))
					tagRouter.DELETE("/:id", s.wrapHandler(s.TagHandler.DeleteTag))
				}
				{
					photoRouter := authRouter.Group("/photos")
					photoRouter.GET("/latest", s.wrapHandler(s.PhotoHandler.ListPhoto))
					photoRouter.GET("", s.wrapHandler(s.PhotoHandler.PagePhotos))
					photoRouter.GET("/:id", s.wrapHandler(s.PhotoHandler.GetPhotoByID))
					photoRouter.DELETE("/batch", s.wrapHandler(s.PhotoHandler.DeletePhotoBatch))
					photoRouter.POST("", s.wrapHandler(s.PhotoHandler.CreatePhoto))
					photoRouter.POST("/batch", s.wrapHandler(s.PhotoHandler.CreatePhotoBatch))
					photoRouter.PUT("/:id", s.wrapHandler(s.PhotoHandler.UpdatePhoto))
					photoRouter.GET("/teams", s.wrapHandler(s.PhotoHandler.ListPhotoTeams))
				}
				{
					userRouter := authRouter.Group("/users")
					userRouter.GET("/profiles", s.wrapHandler(s.UserHandler.GetCurrentUserProfile))
					userRouter.PUT("/profiles", s.wrapHandler(s.UserHandler.UpdateUserProfile))
					userRouter.PUT("/profiles/password", s.wrapHandler(s.UserHandler.UpdatePassword))
					userRouter.PUT("/mfa/generate", s.wrapHandler(s.UserHandler.GenerateMFAQRCode))
					userRouter.PUT("/mfa/update", s.wrapHandler(s.UserHandler.UpdateMFA))
				}
				{
					themeRouter := authRouter.Group("themes")
					themeRouter.GET("/activation", s.wrapHandler(s.ThemeHandler.GetActivatedTheme))
					themeRouter.GET("/:themeID", s.wrapHandler(s.ThemeHandler.GetThemeByID))
					themeRouter.GET("", s.wrapHandler(s.ThemeHandler.ListAllThemes))
					themeRouter.GET("/activation/files", s.wrapHandler(s.ThemeHandler.ListActivatedThemeFile))
					themeRouter.GET("/:themeID/files", s.wrapHandler(s.ThemeHandler.ListThemeFileByID))
					themeRouter.GET("files/content", s.wrapHandler(s.ThemeHandler.GetThemeFileContent))
					themeRouter.GET("/:themeID/files/content", s.wrapHandler(s.ThemeHandler.GetThemeFileContentByID))
					themeRouter.PUT("/files/content", s.wrapHandler(s.ThemeHandler.UpdateThemeFile))
					themeRouter.PUT("/:themeID/files/content", s.wrapHandler(s.ThemeHandler.UpdateThemeFileByID))
					themeRouter.GET("activation/template/custom/sheet", s.wrapHandler(s.ThemeHandler.ListCustomSheetTemplate))
					themeRouter.GET("activation/template/custom/post", s.wrapHandler(s.ThemeHandler.ListCustomPostTemplate))
					themeRouter.POST("/:themeID/activation", s.wrapHandler(s.ThemeHandler.ActivateTheme))
					themeRouter.GET("activation/configurations", s.wrapHandler(s.ThemeHandler.GetActivatedThemeConfig))
					themeRouter.GET("/:themeID/configurations", s.wrapHandler(s.ThemeHandler.GetThemeConfigByID))
					themeRouter.GET("/:themeID/configurations/groups/:group", s.wrapHandler(s.ThemeHandler.GetThemeConfigByGroup))
					themeRouter.GET("/:themeID/configurations/groups", s.wrapHandler(s.ThemeHandler.GetThemeConfigGroupNames))
					themeRouter.GET("activation/settings", s.wrapHandler(s.ThemeHandler.GetActivatedThemeSettingMap))
					themeRouter.GET("/:themeID/settings", s.wrapHandler(s.ThemeHandler.GetThemeSettingMapByID))
					themeRouter.GET("/:themeID/groups/:group/settings", s.wrapHandler(s.ThemeHandler.GetThemeSettingMapByGroupAndThemeID))
					themeRouter.POST("activation/settings", s.wrapHandler(s.ThemeHandler.SaveActivatedThemeSetting))
					themeRouter.POST("/:themeID/settings", s.wrapHandler(s.ThemeHandler.SaveThemeSettingByID))
					themeRouter.DELETE("/:themeID", s.wrapHandler(s.ThemeHandler.DeleteThemeByID))
					themeRouter.POST("upload", s.wrapHandler(s.ThemeHandler.UploadTheme))
					themeRouter.PUT("upload/:themeID", s.wrapHandler(s.ThemeHandler.UpdateThemeByUpload))
					themeRouter.POST("fetching", s.wrapHandler(s.ThemeHandler.FetchTheme))
					themeRouter.PUT("fetching/:themeID", s.wrapHandler(s.ThemeHandler.UpdateThemeByFetching))
					themeRouter.POST("reload", s.wrapHandler(s.ThemeHandler.ReloadTheme))
					themeRouter.GET("activation/template/exists", s.wrapHandler(s.ThemeHandler.TemplateExist))
				}
				{
					emailRouter := authRouter.Group("/mails")
					emailRouter.POST("/test", s.wrapHandler(s.EmailHandler.Test))
				}
			}
		}
		{
			contentRouter := router.Group("")
			contentRouter.Use(s.LogMiddleware.LoggerWithConfig(middleware.GinLoggerConfig{}), s.RecoveryMiddleware.RecoveryWithLogger(), s.InstallRedirectMiddleware.InstallRedirect())

			contentRouter.POST("/content/:type/:slug/authentication", s.wrapHTMLHandler(s.ViewHandler.Authenticate))

			contentRouter.GET("", s.wrapHTMLHandler(s.IndexHandler.Index))
			contentRouter.GET("/page/:page", s.wrapHTMLHandler(s.IndexHandler.IndexPage))
			contentRouter.GET("/robots.txt", s.wrapTextHandler(s.FeedHandler.Robots))
			contentRouter.GET("/atom", s.wrapTextHandler(s.FeedHandler.Atom))
			contentRouter.GET("/atom.xml", s.wrapTextHandler(s.FeedHandler.Atom))
			contentRouter.GET("/rss", s.wrapTextHandler(s.FeedHandler.Feed))
			contentRouter.GET("/rss.xml", s.wrapTextHandler(s.FeedHandler.Feed))
			contentRouter.GET("/feed", s.wrapTextHandler(s.FeedHandler.Feed))
			contentRouter.GET("/feed.xml", s.wrapTextHandler(s.FeedHandler.Feed))
			contentRouter.GET("/feed/categories/:slug", s.wrapTextHandler(s.FeedHandler.CategoryFeed))
			contentRouter.GET("/atom/categories/:slug", s.wrapTextHandler(s.FeedHandler.CategoryAtom))
			contentRouter.GET("/sitemap.xml", s.wrapTextHandler(s.FeedHandler.SitemapXML))
			contentRouter.GET("/sitemap.html", s.wrapHTMLHandler(s.FeedHandler.SitemapHTML))

			contentRouter.GET("/version", s.wrapHandler(s.ViewHandler.Version))
			contentRouter.GET("/install", s.ViewHandler.Install)
			contentRouter.GET("/logo", s.wrapHandler(s.ViewHandler.Logo))
			contentRouter.GET("/favicon", s.wrapHandler(s.ViewHandler.Favicon))
			contentRouter.GET("/search", s.wrapHTMLHandler(s.ContentSearchHandler.Search))
			contentRouter.GET("/search/page/:page", s.wrapHTMLHandler(s.ContentSearchHandler.PageSearch))
			err := s.registerDynamicRouters(contentRouter)
			if err != nil {
				s.logger.DPanic("regiterDynamicRouters err", zap.Error(err))
			}
		}
		{
			contentAPIRouter := router.Group("/api/content")
			contentAPIRouter.Use(s.LogMiddleware.LoggerWithConfig(middleware.GinLoggerConfig{}), s.RecoveryMiddleware.RecoveryWithLogger())

			contentAPIRouter.GET("/archives/years", s.wrapHandler(s.ContentAPIArchiveHandler.ListYearArchives))
			contentAPIRouter.GET("/archives/months", s.wrapHandler(s.ContentAPIArchiveHandler.ListMonthArchives))

			contentAPIRouter.GET("/categories", s.wrapHandler(s.ContentAPICategoryHandler.ListCategories))
			contentAPIRouter.GET("/categories/:slug/posts", s.wrapHandler(s.ContentAPICategoryHandler.ListPosts))

			contentAPIRouter.GET("/journals", s.wrapHandler(s.ContentAPIJournalHandler.ListJournal))
			contentAPIRouter.GET("/journals/:journalID", s.wrapHandler(s.ContentAPIJournalHandler.GetJournal))
			contentAPIRouter.GET("/journals/:journalID/comments/top_view", s.wrapHandler(s.ContentAPIJournalHandler.ListTopComment))
			contentAPIRouter.GET("/journals/:journalID/comments/:parentID/children", s.wrapHandler(s.ContentAPIJournalHandler.ListChildren))
			contentAPIRouter.GET("/journals/:journalID/comments/tree_view", s.wrapHandler(s.ContentAPIJournalHandler.ListCommentTree))
			contentAPIRouter.GET("/journals/:journalID/comments/list_view", s.wrapHandler(s.ContentAPIJournalHandler.ListComment))
			contentAPIRouter.POST("/journals/comments", s.wrapHandler(s.ContentAPIJournalHandler.CreateComment))
			contentAPIRouter.POST("/journals/:journalID/likes", s.wrapHandler(s.ContentAPIJournalHandler.Like))

			contentAPIRouter.POST("/photos/:photoID/likes", s.wrapHandler(s.ContentAPIPhotoHandler.Like))

			contentAPIRouter.GET("/posts/:postID/comments/top_view", s.wrapHandler(s.ContentAPIPostHandler.ListTopComment))
			contentAPIRouter.GET("/posts/:postID/comments/:parentID/children", s.wrapHandler(s.ContentAPIPostHandler.ListChildren))
			contentAPIRouter.GET("/posts/:postID/comments/tree_view", s.wrapHandler(s.ContentAPIPostHandler.ListCommentTree))
			contentAPIRouter.GET("/posts/:postID/comments/list_view", s.wrapHandler(s.ContentAPIPostHandler.ListComment))
			contentAPIRouter.POST("/posts/comments", s.wrapHandler(s.ContentAPIPostHandler.CreateComment))
			contentAPIRouter.POST("/posts/:postID/likes", s.wrapHandler(s.ContentAPIPostHandler.Like))

			contentAPIRouter.GET("/sheets/:sheetID/comments/top_view", s.wrapHandler(s.ContentAPISheetHandler.ListTopComment))
			contentAPIRouter.GET("/sheets/:sheetID/comments/:parentID/children", s.wrapHandler(s.ContentAPISheetHandler.ListChildren))
			contentAPIRouter.GET("/sheets/:sheetID/comments/tree_view", s.wrapHandler(s.ContentAPISheetHandler.ListCommentTree))
			contentAPIRouter.GET("/sheets/:sheetID/comments/list_view", s.wrapHandler(s.ContentAPISheetHandler.ListComment))
			contentAPIRouter.POST("/sheets/comments", s.wrapHandler(s.ContentAPISheetHandler.CreateComment))

			contentAPIRouter.GET("/links", s.wrapHandler(s.ContentAPILinkHandler.ListLinks))
			contentAPIRouter.GET("/links/team_view", s.wrapHandler(s.ContentAPILinkHandler.LinkTeamVO))

			contentAPIRouter.GET("/options/comment", s.wrapHandler(s.ContentAPIOptionHandler.Comment))

			contentAPIRouter.POST("/comments/:commentID/likes", s.wrapHandler(s.ContentAPICommentHandler.Like))
		}
	}
}

func (s *Server) registerDynamicRouters(contentRouter *gin.RouterGroup) error {
	ctx := context.Background()
	ctx = dal.SetCtxQuery(ctx, dal.GetQueryByCtx(ctx).ReplaceDB(dal.GetDB().Session(
		&gorm.Session{Logger: dal.DB.Logger.LogMode(logger.Warn)},
	)))

	archivePath, err := s.OptionService.GetArchivePrefix(ctx)
	if err != nil {
		return err
	}
	categoryPath, err := s.OptionService.GetCategoryPrefix(ctx)
	if err != nil {
		return err
	}
	sheetPermaLinkType, err := s.OptionService.GetSheetPermalinkType(ctx)
	if err != nil {
		return err
	}
	sheetPath, err := s.OptionService.GetSheetPrefix(ctx)
	if err != nil {
		return err
	}
	tagPath, err := s.OptionService.GetTagPrefix(ctx)
	if err != nil {
		return err
	}
	journalPath, err := s.OptionService.GetJournalPrefix(ctx)
	if err != nil {
		return err
	}

	photoPath, err := s.OptionService.GetPhotoPrefix(ctx)
	if err != nil {
		return err
	}
	linkPath, err := s.OptionService.GetLinkPrefix(ctx)
	if err != nil {
		return err
	}
	contentRouter.GET(archivePath, s.wrapHTMLHandler(s.ArchiveHandler.Archives))
	contentRouter.GET(archivePath+"/page/:page", s.wrapHTMLHandler(s.ArchiveHandler.ArchivesPage))
	contentRouter.GET(archivePath+"/:slug", s.wrapHTMLHandler(s.ArchiveHandler.ArchivesBySlug))

	contentRouter.GET(tagPath, s.wrapHTMLHandler(s.ContentTagHandler.Tags))
	contentRouter.GET(tagPath+"/:slug/page/:page", s.wrapHTMLHandler(s.ContentTagHandler.TagPostPage))
	contentRouter.GET(tagPath+"/:slug", s.wrapHTMLHandler(s.ContentTagHandler.TagPost))

	contentRouter.GET(categoryPath, s.wrapHTMLHandler(s.ContentCategoryHandler.Categories))
	contentRouter.GET(categoryPath+"/:slug", s.wrapHTMLHandler(s.ContentCategoryHandler.CategoryDetail))
	contentRouter.GET(categoryPath+"/:slug/page/:page", s.wrapHTMLHandler(s.ContentCategoryHandler.CategoryDetailPage))

	contentRouter.GET(linkPath, s.wrapHTMLHandler(s.ContentLinkHandler.Link))

	contentRouter.GET(photoPath, s.wrapHTMLHandler(s.ContentPhotoHandler.Phtotos))
	contentRouter.GET(photoPath+"/page/:page", s.wrapHTMLHandler(s.ContentPhotoHandler.PhotosPage))

	contentRouter.GET(journalPath, s.wrapHTMLHandler(s.ContentJournalHandler.Journals))
	contentRouter.GET(journalPath+"/page/:page", s.wrapHTMLHandler(s.ContentJournalHandler.JournalsPage))
	contentRouter.GET("admin_preview/"+archivePath+"/:slug", s.wrapHTMLHandler(s.ArchiveHandler.AdminArchivesBySlug))
	if sheetPermaLinkType == consts.SheetPermaLinkTypeRoot {
		contentRouter.GET("/:slug")
	} else {
		contentRouter.GET(sheetPath+"/:slug", s.wrapHTMLHandler(s.ContentSheetHandler.SheetBySlug))
	}
	contentRouter.GET("admin_preview/"+sheetPath+"/:slug", s.wrapHTMLHandler(s.ContentSheetHandler.AdminSheetBySlug))
	return nil
}
