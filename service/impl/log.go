package impl

import (
	"context"

	"github.com/go-sonic/sonic/consts"
	"github.com/go-sonic/sonic/dal"
	"github.com/go-sonic/sonic/model/dto"
	"github.com/go-sonic/sonic/model/entity"
	"github.com/go-sonic/sonic/model/param"
	"github.com/go-sonic/sonic/service"
)

type logServiceImpl struct{}

func NewLogService() service.LogService {
	return &logServiceImpl{}
}

func (l *logServiceImpl) Clear(ctx context.Context) error {
	logDAL := dal.GetQueryByCtx(ctx).Log
	_, err := logDAL.WithContext(ctx).Delete()
	if err != nil {
		return WrapDBErr(err)
	}
	return nil
}

func (l *logServiceImpl) PageLog(ctx context.Context, page param.Page, sort *param.Sort) ([]*entity.Log, int64, error) {
	logDAL := dal.GetQueryByCtx(ctx).Log
	logDO := logDAL.WithContext(ctx)
	err := BuildSort(sort, &logDAL, &logDO)
	if err != nil {
		return nil, 0, err
	}

	logs, totalCount, err := logDO.FindByPage(page.PageNum*page.PageSize, page.PageSize)
	if err != nil {
		return nil, 0, WrapDBErr(err)
	}
	return logs, totalCount, nil
}

func (l *logServiceImpl) ConvertToDTO(log *entity.Log) *dto.Log {
	return &dto.Log{
		ID:         log.ID,
		LogKey:     log.LogKey,
		LogType:    consts.LogType(log.Type),
		Content:    log.Content,
		IPAddress:  log.IPAddress,
		CreateTime: log.CreateTime.UnixMilli(),
	}
}
