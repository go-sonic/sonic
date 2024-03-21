package admin

import (
	"context"
	"errors"
	"net/http"
	"path"
	"path/filepath"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/go-playground/validator/v10"

	"github.com/go-sonic/sonic/config"
	"github.com/go-sonic/sonic/handler/trans"
	"github.com/go-sonic/sonic/log"
	"github.com/go-sonic/sonic/model/dto"
	"github.com/go-sonic/sonic/model/param"
	"github.com/go-sonic/sonic/service"
	"github.com/go-sonic/sonic/util"
	"github.com/go-sonic/sonic/util/xerr"
)

type BackupHandler struct {
	BackupService service.BackupService
}

func NewBackupHandler(backupService service.BackupService) *BackupHandler {
	return &BackupHandler{
		BackupService: backupService,
	}
}

func (b *BackupHandler) GetWorkDirBackup(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	filename, err := util.MustGetQueryString(_ctx, ctx, "filename")
	if err != nil {
		return nil, err
	}
	return b.BackupService.GetBackup(_ctx, filepath.Join(config.BackupDir, filename), service.WholeSite)
}

func (b *BackupHandler) GetDataBackup(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	filename, err := util.MustGetQueryString(_ctx, ctx, "filename")
	if err != nil {
		return nil, err
	}
	return b.BackupService.GetBackup(_ctx, filepath.Join(config.DataExportDir, filename), service.JSONData)
}

func (b *BackupHandler) GetMarkDownBackup(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	filename, err := util.MustGetQueryString(_ctx, ctx, "filename")
	if err != nil {
		return nil, err
	}
	return b.BackupService.GetBackup(_ctx, filepath.Join(config.BackupMarkdownDir, filename), service.Markdown)
}

func (b *BackupHandler) BackupWholeSite(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	toBackupItems := make([]string, 0)
	err := ctx.BindAndValidate(&toBackupItems)
	if err != nil {
		e := validator.ValidationErrors{}
		if errors.As(err, &e) {
			return nil, xerr.WithStatus(e, xerr.StatusBadRequest).WithMsg(trans.Translate(e))
		}
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest)
	}

	return b.BackupService.BackupWholeSite(_ctx, toBackupItems)
}

func (b *BackupHandler) ListBackups(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	return b.BackupService.ListFiles(_ctx, config.BackupDir, service.WholeSite)
}

func (b *BackupHandler) ListToBackupItems(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	return b.BackupService.ListToBackupItems(_ctx)
}

func (b *BackupHandler) HandleWorkDir(_ctx context.Context, ctx *app.RequestContext) {
	path := string(ctx.URI().Path())
	if path == "/api/admin/backups/work-dir/fetch" {
		wrapHandler(b.GetWorkDirBackup)(_ctx, ctx)
		return
	}
	if path == "/api/admin/backups/work-dir/options" || path == "/api/admin/backups/work-dir/options/" {
		wrapHandler(b.ListToBackupItems)(_ctx, ctx)
		return
	}
	b.DownloadBackups(_ctx, ctx)
}

func (b *BackupHandler) DownloadBackups(_ctx context.Context, ctx *app.RequestContext) {
	filename := ctx.Param("path")
	if filename == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, &dto.BaseDTO{
			Status:  http.StatusBadRequest,
			Message: "Filename parameter does not exist",
		})
		return
	}
	filePath, err := b.BackupService.GetBackupFilePath(_ctx, config.BackupDir, filename)
	if err != nil {
		log.CtxErrorf(_ctx, "err=%+v", err)
		status := xerr.GetHTTPStatus(err)
		ctx.JSON(status, &dto.BaseDTO{Status: status, Message: xerr.GetMessage(err)})
	}
	ctx.File(filePath)
}

func (b *BackupHandler) DeleteBackups(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	filename, err := util.MustGetQueryString(_ctx, ctx, "filename")
	if err != nil {
		return nil, err
	}
	return nil, b.BackupService.DeleteFile(_ctx, config.BackupDir, filename)
}

func (b *BackupHandler) ImportMarkdown(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	fileHeader, err := ctx.FormFile("file")
	if err != nil {
		return nil, xerr.WithMsg(err, "上传文件错误").WithStatus(xerr.StatusBadRequest)
	}
	filenameExt := path.Ext(fileHeader.Filename)
	if filenameExt != ".md" && filenameExt != ".markdown" && filenameExt != ".mdown" {
		return nil, xerr.WithMsg(err, "Unsupported format").WithStatus(xerr.StatusBadRequest)
	}
	return nil, b.BackupService.ImportMarkdown(_ctx, fileHeader)
}

func (b *BackupHandler) ExportData(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	return b.BackupService.ExportData(_ctx)
}

func (b *BackupHandler) HandleData(_ctx context.Context, ctx *app.RequestContext) {
	path := string(ctx.URI().Path())
	if path == "/api/admin/backups/data/fetch" {
		wrapHandler(b.GetDataBackup)(_ctx, ctx)
		return
	}
	if path == "/api/admin/backups/data" || path == "/api/admin/backups/data/" {
		wrapHandler(b.ListExportData)(_ctx, ctx)
		return
	}
	b.DownloadData(_ctx, ctx)
}

func (b *BackupHandler) ListExportData(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	return b.BackupService.ListFiles(_ctx, config.DataExportDir, service.JSONData)
}

func (b *BackupHandler) DownloadData(_ctx context.Context, ctx *app.RequestContext) {
	filename := ctx.Param("path")
	if filename == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, &dto.BaseDTO{
			Status:  http.StatusBadRequest,
			Message: "Filename parameter does not exist",
		})
	}
	filePath, err := b.BackupService.GetBackupFilePath(_ctx, config.DataExportDir, filename)
	if err != nil {
		log.CtxErrorf(_ctx, "err=%+v", err)
		status := xerr.GetHTTPStatus(err)
		ctx.JSON(status, &dto.BaseDTO{Status: status, Message: xerr.GetMessage(err)})
	}
	ctx.File(filePath)
}

func (b *BackupHandler) DeleteDataFile(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	filename, ok := ctx.GetQuery("filename")
	if !ok || filename == "" {
		return nil, xerr.BadParam.New("no filename param").WithStatus(xerr.StatusBadRequest).WithMsg("no filename param")
	}
	return nil, b.BackupService.DeleteFile(_ctx, config.DataExportDir, filename)
}

func (b *BackupHandler) ExportMarkdown(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	var exportMarkdownParam param.ExportMarkdown
	err := ctx.BindAndValidate(&exportMarkdownParam)
	if err != nil {
		e := validator.ValidationErrors{}
		if errors.As(err, &e) {
			return nil, xerr.WithStatus(e, xerr.StatusBadRequest).WithMsg(trans.Translate(e))
		}
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest)
	}
	return b.BackupService.ExportMarkdown(_ctx, exportMarkdownParam.NeedFrontMatter)
}

func (b *BackupHandler) ListMarkdowns(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	return b.BackupService.ListFiles(_ctx, config.BackupMarkdownDir, service.Markdown)
}

func (b *BackupHandler) DeleteMarkdowns(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	filename, err := util.MustGetQueryString(_ctx, ctx, "filename")
	if err != nil {
		return nil, err
	}
	return nil, b.BackupService.DeleteFile(_ctx, config.BackupMarkdownDir, filename)
}

func (b *BackupHandler) DownloadMarkdown(_ctx context.Context, ctx *app.RequestContext) {
	filename := ctx.Param("filename")
	if filename == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, &dto.BaseDTO{
			Status:  http.StatusBadRequest,
			Message: "Filename parameter does not exist",
		})
		return
	}
	filePath, err := b.BackupService.GetBackupFilePath(_ctx, config.BackupMarkdownDir, filename)
	if err != nil {
		log.CtxErrorf(_ctx, "err=%+v", err)
		status := xerr.GetHTTPStatus(err)
		ctx.JSON(status, &dto.BaseDTO{Status: status, Message: xerr.GetMessage(err)})
	}
	ctx.File(filePath)
}

type wrapperHandler func(_ctx context.Context, ctx *app.RequestContext) (interface{}, error)

func wrapHandler(handler wrapperHandler) app.HandlerFunc {
	return func(_ctx context.Context, ctx *app.RequestContext) {
		data, err := handler(_ctx, ctx)
		if err != nil {
			log.CtxErrorf(_ctx, "err=%+v", err)
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
