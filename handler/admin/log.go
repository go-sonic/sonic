package admin

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/go-sonic/sonic/model/dto"
	"github.com/go-sonic/sonic/model/param"
	"github.com/go-sonic/sonic/service"
	"github.com/go-sonic/sonic/util"
	"github.com/go-sonic/sonic/util/xerr"
)

type LogHandler struct {
	LogService service.LogService
}

func NewLogHandler(logService service.LogService) *LogHandler {
	return &LogHandler{
		LogService: logService,
	}
}

func (l *LogHandler) PageLatestLog(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	top, err := util.MustGetQueryInt32(_ctx, ctx, "top")
	if err != nil {
		top = 10
	}
	logs, _, err := l.LogService.PageLog(_ctx, param.Page{PageSize: int(top)}, &param.Sort{Fields: []string{"createTime,desc"}})
	if err != nil {
		return nil, err
	}
	logDTOs := make([]*dto.Log, 0, len(logs))
	for _, log := range logs {
		logDTOs = append(logDTOs, l.LogService.ConvertToDTO(log))
	}
	return logDTOs, nil
}

func (l *LogHandler) PageLog(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	type LogParam struct {
		param.Page
		*param.Sort
	}
	var logParam LogParam
	err := ctx.BindAndValidate(&logParam)
	if err != nil {
		return nil, xerr.WithMsg(err, "parameter error").WithStatus(xerr.StatusBadRequest)
	}
	if logParam.Sort == nil {
		logParam.Sort = &param.Sort{
			Fields: []string{"createTime,desc"},
		}
	}
	logs, totalCount, err := l.LogService.PageLog(_ctx, logParam.Page, logParam.Sort)
	if err != nil {
		return nil, err
	}
	logDTOs := make([]*dto.Log, 0, len(logs))
	for _, log := range logs {
		logDTOs = append(logDTOs, l.LogService.ConvertToDTO(log))
	}
	return dto.NewPage(logDTOs, totalCount, logParam.Page), nil
}

func (l *LogHandler) ClearLog(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	return nil, l.LogService.Clear(_ctx)
}
