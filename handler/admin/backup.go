package admin

import (
	"errors"
	"net/http"
	"path"
	"path/filepath"

	"github.com/gin-gonic/gin"
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

func (b *BackupHandler) GetWorkDirBackup(ctx *gin.Context) (interface{}, error) {
	filename, err := util.MustGetQueryString(ctx, "filename")
	if err != nil {
		return nil, err
	}
	return b.BackupService.GetBackup(ctx, filepath.Join(config.BackupDir, filename), service.WholeSite)
}

func (b *BackupHandler) GetDataBackup(ctx *gin.Context) (interface{}, error) {
	filename, err := util.MustGetQueryString(ctx, "filename")
	if err != nil {
		return nil, err
	}
	return b.BackupService.GetBackup(ctx, filepath.Join(config.DataExportDir, filename), service.JSONData)
}

func (b *BackupHandler) GetMarkDownBackup(ctx *gin.Context) (interface{}, error) {
	filename, err := util.MustGetQueryString(ctx, "filename")
	if err != nil {
		return nil, err
	}
	return b.BackupService.GetBackup(ctx, filepath.Join(config.BackupMarkdownDir, filename), service.Markdown)
}

func (b *BackupHandler) BackupWholeSite(ctx *gin.Context) (interface{}, error) {
	toBackupItems := make([]string, 0)
	err := ctx.ShouldBindJSON(&toBackupItems)
	if err != nil {
		e := validator.ValidationErrors{}
		if errors.As(err, &e) {
			return nil, xerr.WithStatus(e, xerr.StatusBadRequest).WithMsg(trans.Translate(e))
		}
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest)
	}

	return b.BackupService.BackupWholeSite(ctx, toBackupItems)
}

func (b *BackupHandler) ListBackups(ctx *gin.Context) (interface{}, error) {
	return b.BackupService.ListFiles(ctx, config.BackupDir, service.WholeSite)
}

func (b *BackupHandler) ListToBackupItems(ctx *gin.Context) (interface{}, error) {
	return b.BackupService.ListToBackupItems(ctx)
}

func (b *BackupHandler) HandleWorkDir(ctx *gin.Context) {
	path := ctx.Request.URL.Path
	if path == "/api/admin/backups/work-dir/fetch" {
		wrapHandler(b.GetWorkDirBackup)(ctx)
		return
	}
	if path == "/api/admin/backups/work-dir/options" || path == "/api/admin/backups/work-dir/options/" {
		wrapHandler(b.ListToBackupItems)(ctx)
		return
	}
	b.DownloadBackups(ctx)
}

func (b *BackupHandler) DownloadBackups(ctx *gin.Context) {
	filename := ctx.Param("path")
	if filename == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, &dto.BaseDTO{
			Status:  http.StatusBadRequest,
			Message: "Filename parameter does not exist",
		})
		return
	}
	filePath, err := b.BackupService.GetBackupFilePath(ctx, config.BackupDir, filename)
	if err != nil {
		log.CtxErrorf(ctx, "err=%+v", err)
		status := xerr.GetHTTPStatus(err)
		ctx.JSON(status, &dto.BaseDTO{Status: status, Message: xerr.GetMessage(err)})
	}
	ctx.File(filePath)
}

func (b *BackupHandler) DeleteBackups(ctx *gin.Context) (interface{}, error) {
	filename, err := util.MustGetQueryString(ctx, "filename")
	if err != nil {
		return nil, err
	}
	return nil, b.BackupService.DeleteFile(ctx, config.BackupDir, filename)
}

func (b *BackupHandler) ImportMarkdown(ctx *gin.Context) (interface{}, error) {
	fileHeader, err := ctx.FormFile("file")
	if err != nil {
		return nil, xerr.WithMsg(err, "上传文件错误").WithStatus(xerr.StatusBadRequest)
	}
	filenameExt := path.Ext(fileHeader.Filename)
	if filenameExt != ".md" && filenameExt != ".markdown" && filenameExt != ".mdown" {
		return nil, xerr.WithMsg(err, "Unsupported format").WithStatus(xerr.StatusBadRequest)
	}
	return nil, b.BackupService.ImportMarkdown(ctx, fileHeader)
}

func (b *BackupHandler) ExportData(ctx *gin.Context) (interface{}, error) {
	return b.BackupService.ExportData(ctx)
}

func (b *BackupHandler) HandleData(ctx *gin.Context) {
	path := ctx.Request.URL.Path
	if path == "/api/admin/backups/data/fetch" {
		wrapHandler(b.GetDataBackup)(ctx)
		return
	}
	if path == "/api/admin/backups/data" || path == "/api/admin/backups/data/" {
		wrapHandler(b.ListExportData)(ctx)
		return
	}
	b.DownloadData(ctx)
}

func (b *BackupHandler) ListExportData(ctx *gin.Context) (interface{}, error) {
	return b.BackupService.ListFiles(ctx, config.DataExportDir, service.JSONData)
}

func (b *BackupHandler) DownloadData(ctx *gin.Context) {
	filename := ctx.Param("path")
	if filename == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, &dto.BaseDTO{
			Status:  http.StatusBadRequest,
			Message: "Filename parameter does not exist",
		})
	}
	filePath, err := b.BackupService.GetBackupFilePath(ctx, config.DataExportDir, filename)
	if err != nil {
		log.CtxErrorf(ctx, "err=%+v", err)
		status := xerr.GetHTTPStatus(err)
		ctx.JSON(status, &dto.BaseDTO{Status: status, Message: xerr.GetMessage(err)})
	}
	ctx.File(filePath)
}

func (b *BackupHandler) DeleteDataFile(ctx *gin.Context) (interface{}, error) {
	filename, ok := ctx.GetQuery("filename")
	if !ok || filename == "" {
		return nil, xerr.BadParam.New("no filename param").WithStatus(xerr.StatusBadRequest).WithMsg("no filename param")
	}
	return nil, b.BackupService.DeleteFile(ctx, config.DataExportDir, filename)
}

func (b *BackupHandler) ExportMarkdown(ctx *gin.Context) (interface{}, error) {
	var exportMarkdownParam param.ExportMarkdown
	err := ctx.ShouldBindJSON(&exportMarkdownParam)
	if err != nil {
		e := validator.ValidationErrors{}
		if errors.As(err, &e) {
			return nil, xerr.WithStatus(e, xerr.StatusBadRequest).WithMsg(trans.Translate(e))
		}
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest)
	}
	return b.BackupService.ExportMarkdown(ctx, exportMarkdownParam.NeedFrontMatter)
}

func (b *BackupHandler) ListMarkdowns(ctx *gin.Context) (interface{}, error) {
	return b.BackupService.ListFiles(ctx, config.BackupMarkdownDir, service.Markdown)
}

func (b *BackupHandler) DeleteMarkdowns(ctx *gin.Context) (interface{}, error) {
	filename, err := util.MustGetQueryString(ctx, "filename")
	if err != nil {
		return nil, err
	}
	return nil, b.BackupService.DeleteFile(ctx, config.BackupMarkdownDir, filename)
}

func (b *BackupHandler) DownloadMarkdown(ctx *gin.Context) {
	filename := ctx.Param("filename")
	if filename == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, &dto.BaseDTO{
			Status:  http.StatusBadRequest,
			Message: "Filename parameter does not exist",
		})
		return
	}
	filePath, err := b.BackupService.GetBackupFilePath(ctx, config.BackupMarkdownDir, filename)
	if err != nil {
		log.CtxErrorf(ctx, "err=%+v", err)
		status := xerr.GetHTTPStatus(err)
		ctx.JSON(status, &dto.BaseDTO{Status: status, Message: xerr.GetMessage(err)})
	}
	ctx.File(filePath)
}

type wrapperHandler func(ctx *gin.Context) (interface{}, error)

func wrapHandler(handler wrapperHandler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		data, err := handler(ctx)
		if err != nil {
			log.CtxErrorf(ctx, "err=%+v", err)
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
