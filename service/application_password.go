package service

import (
	"context"
	"github.com/go-sonic/sonic/model/dto"
	"github.com/go-sonic/sonic/model/entity"
	"github.com/go-sonic/sonic/model/param"
)

type ApplicationPasswordService interface {
	CreatePwd(ctx context.Context, appPwdParam *param.ApplicationPasswordParam) (*dto.ApplicationPasswordDTO, error)
	DeletePwd(ctx context.Context, appPwdParam *param.ApplicationPasswordParam) error
	List(ctx context.Context) ([]*dto.ApplicationPasswordDTO, error)
	Verify(ctx context.Context, userId int32, pwd string) (*entity.ApplicationPassword, error)
	Update(ctx context.Context, entityId int32, ip string) error
}
