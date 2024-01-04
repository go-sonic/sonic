package impl

import (
	"context"
	"strings"
	"time"

	"github.com/go-sonic/sonic/model/dto"
	"github.com/go-sonic/sonic/model/entity"
	"github.com/go-sonic/sonic/util"
	"github.com/go-sonic/sonic/util/xerr"

	"github.com/go-sonic/sonic/dal"
	"github.com/go-sonic/sonic/model/param"
	"github.com/go-sonic/sonic/service"
)

type applicationPasswordServiceImpl struct {
	OptionService       service.OptionService
	AuthenticateService service.AuthenticateService
}

func NewApplicationPasswordService(optionService service.OptionService, authenticateService service.AuthenticateService) service.ApplicationPasswordService {
	return &applicationPasswordServiceImpl{
		OptionService:       optionService,
		AuthenticateService: authenticateService,
	}
}

func (a *applicationPasswordServiceImpl) CreatePwd(ctx context.Context, param *param.ApplicationPasswordParam) (*dto.ApplicationPasswordDTO, error) {
	var err error
	appPwdDTO := &dto.ApplicationPasswordDTO{
		Name: param.Name,
	}

	user, err := MustGetAuthorizedUser(ctx)
	if err != nil {
		return nil, err
	}

	err = dal.GetQueryByCtx(ctx).Transaction(func(tx *dal.Query) error {
		appPwdDAL := tx.ApplicationPassword

		count, err := appPwdDAL.WithContext(ctx).Select().Where(appPwdDAL.Name.Eq(param.Name), appPwdDAL.UserID.Eq(user.ID)).Count()
		if err != nil {
			return WrapDBErr(err)
		}

		if count > 0 {
			return xerr.BadParam.New("").WithMsg("名称已经存在(Application password name already exists)").WithStatus(xerr.StatusBadRequest)
		}

		token := util.GenUUIDWithOutDash()
		tokenMd5 := util.Md5(token)
		// pass claim token to frond, but save in db after md5
		appPwdDTO.Password = token

		currentTime := time.Now()

		appPwdEntity := &entity.ApplicationPassword{
			CreateTime:       currentTime,
			UpdateTime:       &currentTime,
			Name:             param.Name,
			Password:         tokenMd5,
			UserID:           user.ID,
			LastActivateTime: nil,
			LastActivateIP:   "",
		}
		err = appPwdDAL.WithContext(ctx).Create(appPwdEntity)
		if err != nil {
			return WrapDBErr(err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return appPwdDTO, nil
}

func (a *applicationPasswordServiceImpl) DeletePwd(ctx context.Context, appPwdParam *param.ApplicationPasswordParam) error {
	var err error
	user, err := MustGetAuthorizedUser(ctx)
	if err != nil {
		return err
	}

	if appPwdParam == nil || len(strings.TrimSpace(appPwdParam.Name)) == 0 {
		return xerr.BadParam.New("name参数为空").WithStatus(xerr.StatusBadRequest).
			WithMsg("name 参数不能为空")
	}

	appPwdParam.Name = strings.TrimSpace(appPwdParam.Name)

	appPwdDAL := dal.GetQueryByCtx(ctx).ApplicationPassword
	if _, err = appPwdDAL.WithContext(ctx).Where(appPwdDAL.UserID.Eq(user.ID), appPwdDAL.Name.Eq(appPwdParam.Name)).Delete(); err != nil {
		return WrapDBErr(err)
	}

	return nil
}

func (a *applicationPasswordServiceImpl) List(ctx context.Context) ([]*dto.ApplicationPasswordDTO, error) {
	user, err := MustGetAuthorizedUser(ctx)
	if err != nil {
		return nil, err
	}

	appPwdDAL := dal.GetQueryByCtx(ctx).ApplicationPassword
	entities, err := appPwdDAL.WithContext(ctx).Where(appPwdDAL.UserID.Eq(user.ID)).Find()
	if err != nil {
		return nil, WrapDBErr(err)
	}

	appPwdDTOList := make([]*dto.ApplicationPasswordDTO, len(entities))

	for _, appPwdEntity := range entities {
		appPwdDTOList = append(appPwdDTOList, a.ConvertToDTO(appPwdEntity))
	}

	return appPwdDTOList, nil
}

func (a *applicationPasswordServiceImpl) Verify(ctx context.Context, userID int32, pwd string) (*entity.ApplicationPassword, error) {
	appPwdDAL := dal.GetQueryByCtx(ctx).ApplicationPassword
	entityList, err := appPwdDAL.WithContext(ctx).Where(appPwdDAL.UserID.Eq(userID)).Find()
	if err != nil {
		return nil, WrapDBErr(err)
	}

	pwdMd5 := util.Md5(pwd)

	for _, appPwdEntity := range entityList {
		if appPwdEntity.Password == pwdMd5 {
			return appPwdEntity, nil
		}
	}
	return nil, nil
}

func (a *applicationPasswordServiceImpl) Update(ctx context.Context, entityID int32, ip string) error {
	appPwdDAL := dal.GetQueryByCtx(ctx).ApplicationPassword
	now := time.Now()

	_, err := appPwdDAL.WithContext(ctx).Where(appPwdDAL.ID.Eq(entityID)).Updates(entity.ApplicationPassword{
		LastActivateIP:   ip,
		LastActivateTime: &now,
	})
	if err != nil {
		return WrapDBErr(err)
	}

	return nil
}

func (a *applicationPasswordServiceImpl) ConvertToDTO(appPwdEntity *entity.ApplicationPassword) *dto.ApplicationPasswordDTO {
	var lastActivateTime int64
	if appPwdEntity.LastActivateTime == nil {
		lastActivateTime = appPwdEntity.LastActivateTime.Unix()
	}
	appPwdDTO := &dto.ApplicationPasswordDTO{
		Name:             appPwdEntity.Name,
		LastActiveIP:     appPwdEntity.LastActivateIP,
		LastActivateTime: lastActivateTime,
		CreateTime:       appPwdEntity.CreateTime.Unix(),
	}

	return appPwdDTO
}
