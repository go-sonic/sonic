package impl

import (
	"context"
	"time"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	"github.com/go-sonic/sonic/consts"
	"github.com/go-sonic/sonic/dal"
	"github.com/go-sonic/sonic/event"
	"github.com/go-sonic/sonic/log"
	"github.com/go-sonic/sonic/model/dto"
	"github.com/go-sonic/sonic/model/entity"
	"github.com/go-sonic/sonic/model/param"
	"github.com/go-sonic/sonic/service"
	"github.com/go-sonic/sonic/util"
	"github.com/go-sonic/sonic/util/xerr"
)

type userServiceImpl struct {
	TwoFactorMFAService service.TwoFactorTOTPMFAService
	Event               event.Bus
}

func NewUserService(twoFactorMFAService service.TwoFactorTOTPMFAService, event event.Bus) service.UserService {
	return &userServiceImpl{
		TwoFactorMFAService: twoFactorMFAService,
		Event:               event,
	}
}

func (u *userServiceImpl) GetAllUser(ctx context.Context) ([]*entity.User, error) {
	userDAL := dal.GetQueryByCtx(ctx).User
	users, err := userDAL.WithContext(ctx).Find()
	if err != nil {
		return nil, WrapDBErr(err)
	}
	return users, nil
}

func (u *userServiceImpl) UpdatePassword(ctx context.Context, oldPassword string, newPassword string) error {
	user, err := MustGetAuthorizedUser(ctx)
	if err != nil {
		return err
	}
	passwordMatch := u.PasswordMatch(ctx, user.Password, oldPassword)
	if !passwordMatch {
		return xerr.WithStatus(nil, xerr.StatusBadRequest).WithMsg("old password error")
	}
	if newPassword == oldPassword {
		return xerr.WithStatus(nil, xerr.StatusBadRequest).WithMsg("The new password and the old password cannot be the same")
	}
	userDal := dal.GetQueryByCtx(ctx).User
	updateResult, err := userDal.WithContext(ctx).Where(userDal.ID.Eq(user.ID)).UpdateSimple(userDal.Password.Value(u.EncryptPassword(ctx, newPassword)))
	if err != nil {
		return WrapDBErr(err)
	}
	if updateResult.RowsAffected != 1 {
		return xerr.NoType.New("update password error userId=%d", user.ID).WithMsg("update password error").WithStatus(xerr.StatusInternalServerError)
	}
	u.Event.Publish(ctx, &event.UserUpdateEvent{})
	return nil
}

func (u *userServiceImpl) Update(ctx context.Context, userParam *param.User) (*entity.User, error) {
	user, err := MustGetAuthorizedUser(ctx)
	if err != nil {
		return nil, err
	}
	userDal := dal.GetQueryByCtx(ctx).User
	_, err = userDal.WithContext(ctx).Where(userDal.ID.Eq(user.ID)).UpdateSimple(
		userDal.Nickname.Value(userParam.Nickname),
		userDal.Description.Value(userParam.Description),
		userDal.Username.Value(userParam.Username),
		userDal.Email.Value(userParam.Email),
		userDal.Avatar.Value(userParam.Avatar))
	if err != nil {
		return nil, WrapDBErr(err)
	}
	u.Event.Publish(ctx, &event.UserUpdateEvent{})
	return u.GetByID(ctx, user.ID)
}

func (u *userServiceImpl) UpdateMFA(ctx context.Context, mfaKey string, mfaType consts.MFAType, mfaCode string) error {
	user, err := MustGetAuthorizedUser(ctx)
	if err != nil {
		return err
	}

	switch mfaType {
	case consts.MFATFATotp:
		ok := u.TwoFactorMFAService.ValidateTFACode(mfaKey, mfaCode)
		if !ok {
			return xerr.WithStatus(nil, xerr.StatusBadRequest).WithMsg("Invalid Validation Code")
		}
	case consts.MFANone:
		ok := u.TwoFactorMFAService.ValidateTFACode(user.MfaKey, mfaCode)
		if !ok {
			return xerr.WithStatus(nil, xerr.StatusBadRequest).WithMsg("Invalid Validation Code")
		}
	default:
		return xerr.WithMsg(nil, "Not supported authentication").WithStatus(xerr.StatusBadRequest)
	}

	userDal := dal.GetQueryByCtx(ctx).User
	updateResult, err := userDal.WithContext(ctx).Where(userDal.ID.Eq(user.ID)).UpdateSimple(
		userDal.MfaKey.Value(mfaKey),
		userDal.MfaType.Value(mfaType),
	)
	if err != nil {
		return WrapDBErr(err)
	}
	if updateResult.RowsAffected != 1 {
		return xerr.DB.New("").WithMsg("update mfaKey error").WithStatus(xerr.StatusInternalServerError)
	}
	u.Event.Publish(ctx, &event.UserUpdateEvent{})
	return nil
}

func (u *userServiceImpl) ConvertToDTO(ctx context.Context, user *entity.User) *dto.User {
	userDTO := dto.User{
		ID:          user.ID,
		Username:    user.Username,
		Nickname:    user.Nickname,
		Email:       user.Email,
		Avatar:      user.Avatar,
		Description: user.Description,
		MFAType:     user.MfaType,
		CreateTime:  user.CreateTime.UnixMilli(),
	}
	if user.UpdateTime != nil {
		userDTO.UpdateTime = user.UpdateTime.UnixMilli()
	}
	return &userDTO
}

func (u *userServiceImpl) PasswordMatch(ctx context.Context, hashedPassword, plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	return err == nil
}

func (u *userServiceImpl) MustNotExpire(ctx context.Context, expireTime *time.Time) error {
	if expireTime == nil {
		return nil
	}
	now := time.Now()
	if expireTime.After(now) {
		return xerr.Forbidden.New("账号已被停用，请 %s 后重试", util.TimeFormat(int(expireTime.Sub(now).Seconds()))).WithStatus(xerr.StatusForbidden)
	}
	return nil
}

func (u *userServiceImpl) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	userDal := dal.GetQueryByCtx(ctx).User
	user, err := userDal.WithContext(ctx).Where(userDal.Email.Eq(email)).Take()
	if err != nil {
		return nil, WrapDBErr(err)
	}
	return user, nil
}

func (u *userServiceImpl) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	userDal := dal.GetQueryByCtx(ctx).User
	user, err := userDal.WithContext(ctx).Where(userDal.Username.Eq(username)).Take()
	if err != nil {
		return nil, WrapDBErr(err)
	}
	return user, nil
}

func (u *userServiceImpl) GetByID(ctx context.Context, id int32) (*entity.User, error) {
	userDal := dal.GetQueryByCtx(ctx).User
	user, err := userDal.WithContext(ctx).Where(userDal.ID.Eq(id)).Take()
	if err != nil {
		return nil, WrapDBErr(err)
	}
	return user, nil
}

func (u *userServiceImpl) CreateByParam(ctx context.Context, userParam param.User) (*entity.User, error) {
	if len(userParam.Password) < 8 || len(userParam.Password) > 100 {
		return nil, xerr.BadParam.Wrap(nil).WithMsg("password length err")
	}
	user := &entity.User{
		Description: userParam.Description,
		Email:       userParam.Email,
		Password:    u.EncryptPassword(ctx, userParam.Password),
		Username:    userParam.Username,
		Nickname:    userParam.Nickname,
		MfaKey:      "",
		MfaType:     consts.MFANone,
		Avatar:      userParam.Avatar,
	}
	userDAL := dal.GetQueryByCtx(ctx).User
	err := userDAL.WithContext(ctx).Create(user)
	if err != nil {
		return nil, WrapDBErr(err)
	}
	u.Event.Publish(ctx, &event.UserUpdateEvent{
		UserID: user.ID,
	})
	return user, nil
}

func (u *userServiceImpl) EncryptPassword(ctx context.Context, plainPassword string) string {
	password, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		log.CtxError(ctx, "encrypt password", zap.Error(err))
	}
	return string(password)
}
