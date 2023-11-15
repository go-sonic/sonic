package service

import (
	"context"
	"github.com/go-sonic/sonic/model/dto"
	"github.com/go-sonic/sonic/model/param"
)

type ApplicationPasswordService interface {
	CreatePwd(ctx context.Context, appPwdParam *param.ApplicationPasswordParam) (*dto.ApplicationPasswordDTO, error)
	DeletePwd(ctx context.Context, appPwdParam *param.ApplicationPasswordParam) error
	List(ctx context.Context) ([]*dto.ApplicationPasswordDTO, error)
}
