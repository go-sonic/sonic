package handler

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/dig"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/go-sonic/sonic/config"
	"github.com/go-sonic/sonic/event"
	"github.com/go-sonic/sonic/handler/admin"
	"github.com/go-sonic/sonic/handler/content"
	"github.com/go-sonic/sonic/handler/content/api"
	"github.com/go-sonic/sonic/handler/middleware"
	"github.com/go-sonic/sonic/model/dto"
	"github.com/go-sonic/sonic/service"
	"github.com/go-sonic/sonic/template"
	"github.com/go-sonic/sonic/util/xerr"
)

type Server struct {
	logger                    *zap.Logger
	Config                    *config.Config
	HTTPServer                *http.Server
	Router                    *gin.Engine
	Template                  *template.Template
	AuthMiddleware            *middleware.AuthMiddleware
	LogMiddleware             *middleware.GinLoggerMiddleware
	RecoveryMiddleware        *middleware.RecoveryMiddleware
	InstallRedirectMiddleware *middleware.InstallRedirectMiddleware
	OptionService             service.OptionService
	ThemeService              service.ThemeService
	SheetService              service.SheetService
	AdminHandler              *admin.AdminHandler
	AttachmentHandler         *admin.AttachmentHandler
	BackupHandler             *admin.BackupHandler
	CategoryHandler           *admin.CategoryHandler
	InstallHandler            *admin.InstallHandler
	JournalHandler            *admin.JournalHandler
	JournalCommentHandler     *admin.JournalCommentHandler
	LinkHandler               *admin.LinkHandler
	LogHandler                *admin.LogHandler
	MenuHandler               *admin.MenuHandler
	OptionHandler             *admin.OptionHandler
	PhotoHandler              *admin.PhotoHandler
	PostHandler               *admin.PostHandler
	PostCommentHandler        *admin.PostCommentHandler
	SheetHandler              *admin.SheetHandler
	SheetCommentHandler       *admin.SheetCommentHandler
	StatisticHandler          *admin.StatisticHandler
	TagHandler                *admin.TagHandler
	ThemeHandler              *admin.ThemeHandler
	UserHandler               *admin.UserHandler
	EmailHandler              *admin.EmailHandler
	IndexHandler              *content.IndexHandler
	FeedHandler               *content.FeedHandler
	ArchiveHandler            *content.ArchiveHandler
	ViewHandler               *content.ViewHandler
	ContentCategoryHandler    *content.CategoryHandler
	ContentSheetHandler       *content.SheetHandler
	ContentTagHandler         *content.TagHandler
	ContentLinkHandler        *content.LinkHandler
	ContentPhotoHandler       *content.PhotoHandler
	ContentJournalHandler     *content.JournalHandler
	ContentSearchHandler      *content.SearchHandler
	ContentAPIArchiveHandler  *api.ArchiveHandler
	ContentAPICategoryHandler *api.CategoryHandler
	ContentAPIJournalHandler  *api.JournalHandler
	ContentAPILinkHandler     *api.LinkHandler
	ContentAPIPostHandler     *api.PostHandler
	ContentAPISheetHandler    *api.SheetHandler
	ContentAPIOptionHandler   *api.OptionHandler
	ContentAPIPhotoHandler    *api.PhotoHandler
	ContentAPICommentHandler  *api.CommentHandler
}

type ServerParams struct {
	dig.In
	Config                    *config.Config
	Logger                    *zap.Logger
	Event                     event.Bus
	Template                  *template.Template
	AuthMiddleware            *middleware.AuthMiddleware
	LogMiddleware             *middleware.GinLoggerMiddleware
	RecoveryMiddleware        *middleware.RecoveryMiddleware
	InstallRedirectMiddleware *middleware.InstallRedirectMiddleware
	OptionService             service.OptionService
	ThemeService              service.ThemeService
	SheetService              service.SheetService
	AdminHandler              *admin.AdminHandler
	AttachmentHandler         *admin.AttachmentHandler
	BackupHandler             *admin.BackupHandler
	CategoryHandler           *admin.CategoryHandler
	InstallHandler            *admin.InstallHandler
	JournalHandler            *admin.JournalHandler
	JournalCommentHandler     *admin.JournalCommentHandler
	LinkHandler               *admin.LinkHandler
	LogHandler                *admin.LogHandler
	MenuHandler               *admin.MenuHandler
	OptionHandler             *admin.OptionHandler
	PhotoHandler              *admin.PhotoHandler
	PostHandler               *admin.PostHandler
	PostCommentHandler        *admin.PostCommentHandler
	SheetHandler              *admin.SheetHandler
	SheetCommentHandler       *admin.SheetCommentHandler
	StatisticHandler          *admin.StatisticHandler
	TagHandler                *admin.TagHandler
	ThemeHandler              *admin.ThemeHandler
	UserHandler               *admin.UserHandler
	EmailHandler              *admin.EmailHandler
	IndexHandler              *content.IndexHandler
	FeedHandler               *content.FeedHandler
	ArchiveHandler            *content.ArchiveHandler
	ViewHandler               *content.ViewHandler
	ContentCategoryHandler    *content.CategoryHandler
	ContentSheetHandler       *content.SheetHandler
	ContentTagHandler         *content.TagHandler
	ContentLinkHandler        *content.LinkHandler
	ContentPhotoHandler       *content.PhotoHandler
	ContentJournalHandler     *content.JournalHandler
	ContentSearchHandler      *content.SearchHandler
	ContentAPIArchiveHandler  *api.ArchiveHandler
	ContentAPICategoryHandler *api.CategoryHandler
	ContentAPIJournalHandler  *api.JournalHandler
	ContentAPILinkHandler     *api.LinkHandler
	ContentAPIPostHandler     *api.PostHandler
	ContentAPISheetHandler    *api.SheetHandler
	ContentAPIOptionHandler   *api.OptionHandler
	ContentAPIPhotoHandler    *api.PhotoHandler
	ContentAPICommentHandler  *api.CommentHandler
}

func NewServer(param ServerParams, lifecycle fx.Lifecycle) *Server {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	conf := param.Config

	httpServer := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", conf.Server.Host, conf.Server.Port),
		Handler: router,
	}

	s := &Server{
		logger:                    param.Logger,
		Config:                    param.Config,
		HTTPServer:                httpServer,
		Router:                    router,
		Template:                  param.Template,
		AuthMiddleware:            param.AuthMiddleware,
		LogMiddleware:             param.LogMiddleware,
		RecoveryMiddleware:        param.RecoveryMiddleware,
		InstallRedirectMiddleware: param.InstallRedirectMiddleware,
		AdminHandler:              param.AdminHandler,
		AttachmentHandler:         param.AttachmentHandler,
		BackupHandler:             param.BackupHandler,
		CategoryHandler:           param.CategoryHandler,
		InstallHandler:            param.InstallHandler,
		JournalHandler:            param.JournalHandler,
		JournalCommentHandler:     param.JournalCommentHandler,
		LinkHandler:               param.LinkHandler,
		LogHandler:                param.LogHandler,
		MenuHandler:               param.MenuHandler,
		OptionHandler:             param.OptionHandler,
		PhotoHandler:              param.PhotoHandler,
		PostHandler:               param.PostHandler,
		PostCommentHandler:        param.PostCommentHandler,
		SheetHandler:              param.SheetHandler,
		SheetCommentHandler:       param.SheetCommentHandler,
		StatisticHandler:          param.StatisticHandler,
		TagHandler:                param.TagHandler,
		ThemeHandler:              param.ThemeHandler,
		UserHandler:               param.UserHandler,
		EmailHandler:              param.EmailHandler,
		OptionService:             param.OptionService,
		ThemeService:              param.ThemeService,
		SheetService:              param.SheetService,
		IndexHandler:              param.IndexHandler,
		FeedHandler:               param.FeedHandler,
		ArchiveHandler:            param.ArchiveHandler,
		ViewHandler:               param.ViewHandler,
		ContentCategoryHandler:    param.ContentCategoryHandler,
		ContentSheetHandler:       param.ContentSheetHandler,
		ContentTagHandler:         param.ContentTagHandler,
		ContentLinkHandler:        param.ContentLinkHandler,
		ContentPhotoHandler:       param.ContentPhotoHandler,
		ContentJournalHandler:     param.ContentJournalHandler,
		ContentAPIArchiveHandler:  param.ContentAPIArchiveHandler,
		ContentAPICategoryHandler: param.ContentAPICategoryHandler,
		ContentAPIJournalHandler:  param.ContentAPIJournalHandler,
		ContentAPILinkHandler:     param.ContentAPILinkHandler,
		ContentAPIPostHandler:     param.ContentAPIPostHandler,
		ContentAPISheetHandler:    param.ContentAPISheetHandler,
		ContentAPIOptionHandler:   param.ContentAPIOptionHandler,
		ContentSearchHandler:      param.ContentSearchHandler,
		ContentAPIPhotoHandler:    param.ContentAPIPhotoHandler,
		ContentAPICommentHandler:  param.ContentAPICommentHandler,
	}
	lifecycle.Append(fx.Hook{
		OnStop:  httpServer.Shutdown,
		OnStart: s.Run,
	})
	return s
}

func (s *Server) Run(ctx context.Context) error {
	if config.IsDev() {
		gin.SetMode(gin.DebugMode)
	}
	go func() {
		if err := s.HTTPServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			// print err info when httpServer start failed
			s.logger.Error("unexpected error from ListenAndServe", zap.Error(err))
			fmt.Printf("http server start error:%s\n", err.Error())
			os.Exit(1)
		}
	}()
	return nil
}

type wrapperHandler func(ctx *gin.Context) (interface{}, error)

func (s *Server) wrapHandler(handler wrapperHandler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		data, err := handler(ctx)
		if err != nil {
			s.logger.Error("handler error", zap.Error(err))
			status := xerr.GetHTTPStatus(err)
			ctx.JSON(status, &dto.BaseDTO{Status: status, Message: xerr.GetMessage(err)})
			return
		}

		ctx.JSON(http.StatusOK, &dto.BaseDTO{
			Status:  http.StatusOK,
			Data:    data,
			Message: "OK",
		})
	}
}

type wrapperHTMLHandler func(ctx *gin.Context, model template.Model) (templateName string, err error)

var (
	htmlContentType = []string{"text/html; charset=utf-8"}
	xmlContentType  = []string{"application/xml; charset=utf-8"}
)

func (s *Server) wrapHTMLHandler(handler wrapperHTMLHandler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		model := template.Model{}
		templateName, err := handler(ctx, model)
		if err != nil {
			s.handleError(ctx, err)
			return
		}
		if templateName == "" {
			return
		}
		header := ctx.Writer.Header()
		if val := header["Content-Type"]; len(val) == 0 {
			header["Content-Type"] = htmlContentType
		}
		err = s.Template.ExecuteTemplate(ctx.Writer, templateName, model)
		if err != nil {
			s.logger.Error("render template err", zap.Error(err))
		}
	}
}

func (s *Server) wrapTextHandler(handler wrapperHTMLHandler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		model := template.Model{}
		templateName, err := handler(ctx, model)
		if err != nil {
			s.handleError(ctx, err)
			return
		}
		header := ctx.Writer.Header()
		if val := header["Content-Type"]; len(val) == 0 {
			header["Content-Type"] = xmlContentType
		}
		err = s.Template.ExecuteTextTemplate(ctx.Writer, templateName, model)
		if err != nil {
			s.logger.Error("render template err", zap.Error(err))
		}
	}
}

func (s *Server) handleError(ctx *gin.Context, err error) {
	status := xerr.GetHTTPStatus(err)
	message := xerr.GetMessage(err)
	model := template.Model{}

	templateName, _ := s.ThemeService.Render(ctx, strconv.Itoa(status))
	t := s.Template.HTMLTemplate.Lookup(templateName)
	if t == nil {
		templateName = "common/error/error"
	}

	header := ctx.Writer.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = htmlContentType
	}

	model["status"] = status
	model["message"] = message
	model["err"] = err

	err = s.Template.ExecuteTemplate(ctx.Writer, templateName, model)
	if err != nil {
		s.logger.Error("render error template err", zap.Error(err))
	}
}
