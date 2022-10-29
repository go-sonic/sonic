package service

import (
	"context"
	"time"

	"github.com/go-sonic/sonic/consts"
	"github.com/go-sonic/sonic/model/dto"
	"github.com/go-sonic/sonic/model/entity"
	"github.com/go-sonic/sonic/model/param"
)

type UserService interface {
	GetAllUser(ctx context.Context) ([]*entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	ConvertToDTO(ctx context.Context, user *entity.User) *dto.User
	GetByUsername(ctx context.Context, username string) (*entity.User, error)
	MustNotExpire(ctx context.Context, expireTime *time.Time) error
	PasswordMatch(ctx context.Context, hashedPassword, plainPassword string) bool
	GetByID(ctx context.Context, id int32) (*entity.User, error)
	CreateByParam(ctx context.Context, userParam param.User) (*entity.User, error)
	Update(ctx context.Context, userParam *param.User) (*entity.User, error)
	UpdatePassword(ctx context.Context, oldPassword string, newPassword string) error
	UpdateMFA(ctx context.Context, mfaKey string, mfaType consts.MFAType, mfaCode string) error
	EncryptPassword(ctx context.Context, plainPassword string) string
}
