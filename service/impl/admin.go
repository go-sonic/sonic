package impl

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	uuid2 "github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/go-sonic/sonic/cache"
	"github.com/go-sonic/sonic/config"
	"github.com/go-sonic/sonic/consts"
	"github.com/go-sonic/sonic/dal"
	"github.com/go-sonic/sonic/event"
	"github.com/go-sonic/sonic/log"
	"github.com/go-sonic/sonic/model/dto"
	"github.com/go-sonic/sonic/model/entity"
	"github.com/go-sonic/sonic/model/param"
	"github.com/go-sonic/sonic/model/property"
	"github.com/go-sonic/sonic/service"
	"github.com/go-sonic/sonic/util"
	"github.com/go-sonic/sonic/util/xerr"
)

type adminServiceImpl struct {
	UserService      service.UserService
	OptionService    service.OptionService
	Cache            cache.Cache
	Config           *config.Config
	Event            event.Bus
	TwoFactorTOTPMFA service.TwoFactorTOTPMFAService
	EmailService     service.EmailService
}

func NewAdminService(userService service.UserService, cache cache.Cache, config *config.Config, event event.Bus, twoFactorMFA service.TwoFactorTOTPMFAService, emailService service.EmailService) service.AdminService {
	return &adminServiceImpl{
		UserService:      userService,
		Cache:            cache,
		Config:           config,
		Event:            event,
		TwoFactorTOTPMFA: twoFactorMFA,
		EmailService:     emailService,
	}
}

func (a *adminServiceImpl) Authenticate(ctx context.Context, loginParam param.LoginParam) (*entity.User, error) {
	missMatchTip := "用户名或密码不正确"

	var user *entity.User
	err := util.Validate.Var(loginParam.Username, "email")

	if err != nil {
		user, err = a.UserService.GetByUsername(ctx, loginParam.Username)
	} else {
		user, err = a.UserService.GetByEmail(ctx, loginParam.Username)
	}

	if xerr.GetType(err) == xerr.NoRecord {
		return nil, xerr.WithMsg(err, missMatchTip).WithStatus(xerr.StatusBadRequest)
	}
	if err != nil {
		return nil, err
	}

	err = a.UserService.MustNotExpire(ctx, user.ExpireTime)
	if err != nil {
		return nil, err
	}

	if !a.UserService.PasswordMatch(ctx, user.Password, loginParam.Password) {
		return nil, xerr.BadParam.New("").WithMsg(missMatchTip).WithStatus(xerr.StatusBadRequest)
	}
	return user, nil
}

func (a *adminServiceImpl) Auth(ctx context.Context, loginParam param.LoginParam) (*dto.AuthTokenDTO, error) {
	user, err := a.Authenticate(ctx, loginParam)
	if err != nil {
		return nil, err
	}
	if a.TwoFactorTOTPMFA.UseMFA(user.MfaType) {
		if len(loginParam.AuthCode) != 6 {
			return nil, xerr.WithMsg(nil, "请输入6位两步验证码").WithStatus(xerr.StatusBadRequest)
		}
		mfaAuth := a.TwoFactorTOTPMFA.ValidateTFACode(user.MfaKey, loginParam.AuthCode)
		if !mfaAuth {
			return nil, xerr.WithStatus(nil, xerr.StatusBadRequest).WithMsg("两步验证码验证错误")
		}
	}
	a.Event.Publish(ctx, &event.LogEvent{
		LogKey:    user.Username,
		LogType:   consts.LogTypeLoggedIn,
		Content:   user.Nickname,
		IPAddress: util.GetClientIP(ctx),
	})
	return a.buildAuthToken(user), nil
}

func (a *adminServiceImpl) ClearToken(ctx context.Context) error {
	user, ok := GetAuthorizedUser(ctx)
	if !ok || user == nil {
		return xerr.Forbidden.New("").WithStatus(xerr.StatusForbidden).WithMsg("未登录")
	}
	accessToken, _ := a.Cache.Get(cache.BuildAccessTokenKey(user.ID))
	refreshToken, _ := a.Cache.Get(cache.BuildRefreshTokenKey(user.ID))

	a.Cache.Delete(cache.BuildTokenAccessKey(accessToken.(string)))
	a.Cache.Delete(cache.BuildTokenRefreshKey(refreshToken.(string)))

	a.Cache.Delete(cache.BuildAccessTokenKey(user.ID))
	a.Cache.Delete(cache.BuildRefreshTokenKey(user.ID))
	a.Event.Publish(ctx, &event.LogEvent{
		LogKey:    user.Username,
		LogType:   consts.LogTypeLoggedOut,
		Content:   user.Nickname,
		IPAddress: util.GetClientIP(ctx),
	})
	return nil
}

func (a *adminServiceImpl) SendResetPasswordCode(ctx context.Context, resetParam param.ResetPasswordParam) error {
	user, ok := GetAuthorizedUser(ctx)
	if !ok || user == nil {
		return xerr.Forbidden.New("").WithStatus(xerr.StatusForbidden).WithMsg("未登录")
	}
	_, ok = a.Cache.Get(cache.BuildCodeCacheKey(user.ID))
	if ok {
		return xerr.NoType.New("").WithMsg("已经获取过验证码，不能重复获取").WithStatus(xerr.StatusInternalServerError)
	}
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	code := fmt.Sprintf("%06v", rnd.Int31n(1000000))

	log.CtxInfof(ctx, "重置密码验证码: %v", code)
	a.Cache.Set(cache.BuildCodeCacheKey(user.ID), code, consts.CodeValidDuration)
	emailEnabled, err := a.OptionService.GetOrByDefaultWithErr(ctx, property.EmailIsEnabled, false)
	if err != nil {
		return err
	}
	if !emailEnabled.(bool) {
		return xerr.NoType.New("未启用 SMTP 服务").WithMsg("未启用 SMTP 服务，无法发送邮件，但是你可以通过系统日志找到验证码")
	}
	content := "您正在进行密码重置操作，如不是本人操作，请尽快做好相应措施。密码重置验证码如下（五分钟有效）：\n" + code
	return a.EmailService.SendTextEmail(ctx, resetParam.Email, "找回密码验证码", content)
}

func (a *adminServiceImpl) buildAuthToken(user *entity.User) *dto.AuthTokenDTO {
	accessToken := uuid2.New().String()
	refreshToken := uuid2.New().String()

	authToken := &dto.AuthTokenDTO{}
	authToken.AccessToken = accessToken
	authToken.ExpiredIn = consts.AccessTokenExpiredSeconds
	authToken.RefreshToken = refreshToken

	a.Cache.Set(cache.BuildTokenAccessKey(accessToken), user.ID, time.Second*consts.AccessTokenExpiredSeconds)
	a.Cache.Set(cache.BuildTokenRefreshKey(refreshToken), user.ID, consts.RefreshTokenExpiredDays*24*3600*time.Second)

	a.Cache.Set(cache.BuildAccessTokenKey(user.ID), accessToken, time.Second*consts.AccessTokenExpiredSeconds)
	a.Cache.Set(cache.BuildRefreshTokenKey(user.ID), refreshToken, consts.RefreshTokenExpiredDays*24*3600*time.Second)

	return authToken
}

func (a *adminServiceImpl) RefreshToken(ctx context.Context, refreshToken string) (*dto.AuthTokenDTO, error) {
	userID, ok := a.Cache.Get(cache.BuildTokenRefreshKey(refreshToken))
	if !ok {
		return nil, xerr.WithMsg(nil, "登录状态已失效，请重新登录").WithStatus(xerr.StatusBadRequest)
	}
	userDAL := dal.GetQueryByCtx(ctx).User
	user, err := userDAL.WithContext(ctx).Where(userDAL.ID.Eq(userID.(int32))).First()
	if err != nil {
		return nil, err
	}
	return a.buildAuthToken(user), nil
}

func (a *adminServiceImpl) GetEnvironments(ctx context.Context) *dto.EnvironmentDTO {
	environments := &dto.EnvironmentDTO{
		Database:  string(dal.DBType) + " " + consts.DatabaseVersion,
		Version:   consts.SonicVersion,
		StartTime: consts.StartTime.UnixMilli(),
		Mode:      util.IfElse(a.Config.Sonic.Mode == "", "production", a.Config.Sonic.Mode).(string),
	}
	return environments
}

func (a *adminServiceImpl) GetLogFiles(ctx context.Context, lineNum int64) (string, error) {
	errTips := "读取日志文件错误"
	lineEndByte := []byte("\n")[0]
	linesTotalByteNum := 0
	var linesCount int64

	fileName := filepath.Join(a.Config.Sonic.LogDir, a.Config.Log.FileName)
	file, err := os.Open(fileName)
	if err != nil {
		return "", xerr.WithStatus(err, xerr.StatusInternalServerError).WithMsg(errTips)
	}

	defer func() {
		err = file.Close()
		if err != nil {
			log.CtxError(ctx, "关闭日志文件错误", zap.Error(err))
		}
	}()

	fileInfo, err := file.Stat()
	if err != nil {
		return "", xerr.WithStatus(err, xerr.StatusInternalServerError).WithMsg(errTips)
	}

	fileBytesLength := fileInfo.Size()
	position := fileBytesLength - 1

	bufReader := bufio.NewReader(file)
	curLine := bytes.Buffer{}

	globalIsPrefix := false
	lines := make([][]byte, 0, lineNum)

	for position > 0 {
		if !globalIsPrefix {
			position--

			_, err = file.Seek(position, 0)
			if err != nil {
				return "", xerr.WithMsg(err, errTips).WithStatus(xerr.StatusInternalServerError)
			}

			bufReader.Reset(file)

			curByte, err := bufReader.ReadByte()
			if err != nil {
				return "", xerr.WithMsg(err, errTips).WithStatus(xerr.StatusInternalServerError)
			}
			if curByte != lineEndByte {
				continue
			}
		}

		content, isPrefix, err := bufReader.ReadLine()
		if err != nil {
			return "", xerr.WithMsg(err, errTips).WithStatus(xerr.StatusInternalServerError)
		}

		globalIsPrefix = isPrefix

		if !isPrefix {
			curLine.Write(content)
			lines = append(lines, curLine.Bytes())
			linesTotalByteNum += len(content) + 1
			linesCount++
			curLine = bytes.Buffer{}
		} else {
			curLine.Write(content)
		}
		if linesCount == lineNum {
			break
		}
	}
	result := bytes.Buffer{}
	result.Grow(linesTotalByteNum)

	for i := len(lines) - 1; i >= 0; i-- {
		result.Write(lines[i])
		result.WriteByte(lineEndByte)
	}
	return result.String(), nil
}
