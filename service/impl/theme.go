package impl

import (
	"context"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"go.uber.org/fx"

	"github.com/go-sonic/sonic/config"
	"github.com/go-sonic/sonic/dal"
	"github.com/go-sonic/sonic/event"
	"github.com/go-sonic/sonic/model/dto"
	"github.com/go-sonic/sonic/model/entity"
	"github.com/go-sonic/sonic/model/property"
	"github.com/go-sonic/sonic/service"
	"github.com/go-sonic/sonic/service/theme"
	"github.com/go-sonic/sonic/util"
	"github.com/go-sonic/sonic/util/xerr"
)

type themeServiceImpl struct {
	OptionService   service.OptionService
	Config          *config.Config
	Event           event.Bus
	PropertyScanner theme.PropertyScanner
	FileScanner     theme.FileScanner
	ThemeFetchers   themeFetchers
}

type themeFetchers struct {
	fx.In
	MultipartZipThemeFetcher theme.ThemeFetcher `name:"multipartZipThemeFetcher"`
	GitRepoThemeFetcher      theme.ThemeFetcher `name:"gitRepoThemeFetcher"`
}

func NewThemeService(optionService service.OptionService, config *config.Config, event event.Bus, propertyScanner theme.PropertyScanner, fileScanner theme.FileScanner, themeFetcher themeFetchers) service.ThemeService {
	return &themeServiceImpl{
		OptionService:   optionService,
		Config:          config,
		Event:           event,
		PropertyScanner: propertyScanner,
		FileScanner:     fileScanner,
		ThemeFetchers:   themeFetcher,
	}
}

func (t *themeServiceImpl) GetActivateTheme(ctx context.Context) (*dto.ThemeProperty, error) {
	activatedThemeID, err := t.OptionService.GetActivatedThemeID(ctx)
	if err != nil {
		return nil, err
	}
	return t.GetThemeByID(ctx, activatedThemeID)
}

func (t *themeServiceImpl) GetThemeByID(ctx context.Context, themeID string) (*dto.ThemeProperty, error) {
	themeProperty, err := t.PropertyScanner.GetThemeByThemeID(ctx, themeID)
	if err != nil {
		return nil, err
	}
	if themeProperty == nil {
		return nil, xerr.WithStatus(nil, xerr.StatusBadRequest).WithMsg(themeID + " not exist")
	}
	activatedThemeID, err := t.OptionService.GetActivatedThemeID(ctx)
	if err != nil {
		return nil, err
	}
	if themeProperty.ID == activatedThemeID {
		themeProperty.Activated = true
	}
	return themeProperty, nil
}

func (t *themeServiceImpl) ListAllTheme(ctx context.Context) ([]*dto.ThemeProperty, error) {
	themes, err := t.PropertyScanner.ListAll(ctx, t.Config.Sonic.ThemeDir)
	if err != nil {
		return nil, err
	}
	activatedThemeID, err := t.OptionService.GetActivatedThemeID(ctx)
	if err != nil {
		return nil, err
	}

	for _, t := range themes {
		if t.ID == activatedThemeID {
			t.Activated = true
		}
	}
	return themes, nil
}

func (t *themeServiceImpl) ListThemeFiles(ctx context.Context, themeID string) ([]*dto.ThemeFile, error) {
	themeProperty, err := t.PropertyScanner.GetThemeByThemeID(ctx, themeID)
	if err != nil {
		return nil, err
	}
	if themeProperty == nil {
		return nil, xerr.WithStatus(nil, xerr.StatusBadRequest).WithMsg(themeID + " not exist")
	}
	return t.FileScanner.ListThemeFiles(ctx, themeProperty.ThemePath)
}

func (t *themeServiceImpl) GetThemeFileContent(ctx context.Context, themeID, absPath string) (string, error) {
	themeProperty, err := t.PropertyScanner.GetThemeByThemeID(ctx, themeID)
	if err != nil {
		return "", err
	}
	if themeProperty == nil {
		return "", xerr.WithStatus(nil, xerr.StatusBadRequest).WithMsg(themeID + " not exist")
	}

	if err = t.checkPathValid(themeProperty.ThemePath, absPath); err != nil {
		return "", err
	}

	return t.ReadThemeFile(ctx, absPath)
}

func (t *themeServiceImpl) ReadThemeFile(ctx context.Context, absPath string) (content string, err error) {
	file, err := os.Open(absPath)
	defer func() {
		cerr := file.Close()
		if err == nil && cerr != nil {
			err = xerr.WithStatus(cerr, xerr.StatusInternalServerError).WithMsg("close file err")
		}
	}()

	if os.IsNotExist(err) {
		return "", xerr.WithMsg(err, "file not exist").WithStatus(xerr.StatusBadRequest)
	} else if os.IsPermission(err) {
		return "", xerr.WithMsg(err, "file permission err").WithStatus(xerr.StatusForbidden)
	}

	contentBytes, err := io.ReadAll(file)
	if err != nil {
		err = xerr.WithMsg(err, "read file error").WithStatus(xerr.StatusInternalServerError)
		return
	}
	content = util.BytesToString(contentBytes)
	return
}

func (t *themeServiceImpl) UpdateThemeFile(ctx context.Context, themeID, absPath, content string) error {
	themeProperty, err := t.PropertyScanner.GetThemeByThemeID(ctx, themeID)
	if err != nil {
		return err
	}
	if themeProperty == nil {
		return xerr.WithStatus(nil, xerr.StatusBadRequest).WithMsg(themeID + " not exist")
	}

	if err = t.checkPathValid(themeProperty.ThemePath, absPath); err != nil {
		return err
	}
	file, err := os.OpenFile(absPath, os.O_WRONLY, 0)
	if err != nil {
		return xerr.WithMsg(err, "open file error")
	}
	defer file.Close()
	_, err = file.WriteString(content)
	if err != nil {
		return xerr.WithMsg(err, "write to file err")
	}
	return nil
}

func (t *themeServiceImpl) checkPathValid(themePath, absPath string) error {
	absPath = filepath.Clean(absPath)
	if !filepath.IsAbs(absPath) {
		return xerr.WithMsg(nil, "path error").WithStatus(xerr.StatusForbidden)
	}
	if !strings.HasPrefix(absPath, themePath) {
		return xerr.WithMsg(nil, "path error").WithStatus(xerr.StatusForbidden)
	}
	return nil
}

func (t *themeServiceImpl) ListCustomTemplates(ctx context.Context, themeID, prefix string) ([]string, error) {
	themeProperty, err := t.PropertyScanner.GetThemeByThemeID(ctx, themeID)
	if err != nil {
		return nil, err
	}
	if themeProperty == nil {
		return nil, xerr.WithMsg(nil, "theme does not exist").WithStatus(xerr.StatusInternalServerError)
	}
	files, err := os.ReadDir(themeProperty.ThemePath)
	if err != nil {
		return nil, xerr.WithMsg(err, "read theme files error").WithStatus(xerr.StatusInternalServerError)
	}
	result := make([]string, 0)
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		fileName := file.Name()
		if !strings.HasPrefix(fileName, prefix) {
			continue
		}
		customName := strings.TrimPrefix(fileName, prefix)
		customName = strings.TrimSuffix(customName, ".ftl")
		result = append(result, customName)
	}
	return result, nil
}

func (t *themeServiceImpl) ActivateTheme(ctx context.Context, themeID string) (*dto.ThemeProperty, error) {
	err := t.OptionService.Save(ctx, map[string]string{property.Theme.KeyValue: themeID})
	if err != nil {
		return nil, err
	}
	t.Event.Publish(ctx, &event.ThemeActivatedEvent{})
	return t.GetThemeByID(ctx, themeID)
}

func (t *themeServiceImpl) GetThemeConfig(ctx context.Context, themeID string) ([]*dto.ThemeConfigGroup, error) {
	themeProperty, err := t.GetThemeByID(ctx, themeID)
	if err != nil {
		return nil, err
	}
	return t.PropertyScanner.ReadThemeConfig(ctx, themeProperty.ThemePath)
}

func (t *themeServiceImpl) GetThemeSettingMap(ctx context.Context, themeID string) (map[string]interface{}, error) {
	itemMap, err := t.getThemeConfigItemMap(ctx, themeID)
	if err != nil {
		return nil, err
	}
	return t.getThemeSettingMapByItemMap(ctx, themeID, itemMap)
}

func (t *themeServiceImpl) GetThemeGroupSettingMap(ctx context.Context, themeID, group string) (map[string]interface{}, error) {
	themeConfig, err := t.GetThemeConfig(ctx, themeID)
	if err != nil {
		return nil, err
	}

	var groupConfig *dto.ThemeConfigGroup
	for _, themeGroup := range themeConfig {
		if themeGroup.Name == group {
			groupConfig = themeGroup
			break
		}
	}
	if groupConfig == nil {
		return nil, nil
	}
	itemMap := make(map[string]*dto.ThemeConfigItem)
	for _, item := range groupConfig.Items {
		itemMap[item.Name] = item
	}
	return t.getThemeSettingMapByItemMap(ctx, themeID, itemMap)
}

func (t *themeServiceImpl) getThemeSettingMapByItemMap(ctx context.Context, themeID string, itemMap map[string]*dto.ThemeConfigItem) (map[string]interface{}, error) {
	themeSettingDAL := dal.GetQueryByCtx(ctx).ThemeSetting
	themeSettings, err := themeSettingDAL.WithContext(ctx).Where(themeSettingDAL.ThemeID.Eq(themeID)).Find()
	if err != nil {
		return nil, WrapDBErr(err)
	}
	result := make(map[string]interface{})
	for _, themeSetting := range themeSettings {
		item, ok := itemMap[themeSetting.SettingKey]
		if !ok {
			continue
		}
		value, err := item.DataType.Convert(themeSetting.SettingValue)
		if err != nil {
			return nil, xerr.WithStatus(err, xerr.StatusInternalServerError)
		}
		result[themeSetting.SettingKey] = value
	}

	for _, item := range itemMap {
		if _, ok := result[item.Name]; ok || item.DefaultValue == nil {
			continue
		}
		result[item.Name] = item.DefaultValue
	}
	return result, nil
}

func (t *themeServiceImpl) SaveThemeSettings(ctx context.Context, themeID string, settings map[string]interface{}) error {
	if len(settings) == 0 {
		return nil
	}

	itemMap, err := t.getThemeConfigItemMap(ctx, themeID)
	if err != nil {
		return err
	}

	themeSettingDAL := dal.GetQueryByCtx(ctx).ThemeSetting
	allThemeSetting, err := themeSettingDAL.WithContext(ctx).Where(themeSettingDAL.ThemeID.Eq(themeID)).Find()
	if err != nil {
		return WrapDBErr(err)
	}
	allThemeSettingMap := make(map[string]*entity.ThemeSetting, len(allThemeSetting))
	for _, themeSetting := range allThemeSetting {
		allThemeSettingMap[themeSetting.SettingKey] = themeSetting
	}

	toDeleteIDs := make([]int32, 0)
	toCreate := make([]*entity.ThemeSetting, 0)
	toUpdate := make([]*entity.ThemeSetting, 0)
	now := time.Now()

	for name, value := range settings {
		item, ok := itemMap[name]
		if !ok {
			return xerr.WithStatus(nil, xerr.StatusBadRequest).WithMsg("setting name invalid: " + name)
		}
		if themeSetting, ok := allThemeSettingMap[name]; ok {
			if value == "" {
				toDeleteIDs = append(toDeleteIDs, themeSetting.ID)
			} else {
				valueStr, err := item.DataType.FormatToStr(value)
				if err != nil {
					return xerr.WithErrMsgf(err, "value=%s type invalid ,setting name=%s", value, name).WithStatus(xerr.StatusBadRequest)
				}

				if value != themeSetting.SettingValue {
					themeSetting.SettingValue = valueStr
					themeSetting.UpdateTime = util.TimePtr(now)
					toUpdate = append(toUpdate, themeSetting)
				}
			}
		} else {
			if value != "" {
				valueStr, err := item.DataType.FormatToStr(value)
				if err != nil {
					return xerr.WithErrMsgf(err, "value=%s type invalid ,setting name=%s", value, name).WithStatus(xerr.StatusBadRequest)
				}
				toCreate = append(toCreate, &entity.ThemeSetting{SettingKey: name, SettingValue: valueStr, ThemeID: themeID})
			}
		}
	}
	err = dal.Transaction(ctx, func(txCtx context.Context) error {
		themeSettingDAL := dal.GetQueryByCtx(txCtx).ThemeSetting
		_, err := themeSettingDAL.WithContext(txCtx).Where(themeSettingDAL.ID.In(toDeleteIDs...)).Delete()
		if err != nil {
			return WrapDBErr(err)
		}
		err = themeSettingDAL.WithContext(txCtx).Save(toUpdate...)
		if err != nil {
			return WrapDBErr(err)
		}
		err = themeSettingDAL.WithContext(txCtx).Create(toCreate...)
		return WrapDBErr(err)
	})
	if err != nil {
		return err
	}
	t.Event.Publish(ctx, &event.ThemeUpdateEvent{})
	return nil
}

func (t *themeServiceImpl) DeleteThemeSettings(ctx context.Context, themeID string) error {
	activatedThemeID, err := t.OptionService.GetActivatedThemeID(ctx)
	if err != nil {
		return err
	}
	if activatedThemeID == themeID {
		return xerr.WithStatus(nil, xerr.StatusBadRequest).WithMsg("can not delete the theme being used")
	}

	themeSettingDAL := dal.GetQueryByCtx(ctx).ThemeSetting
	_, err = themeSettingDAL.WithContext(ctx).Where(themeSettingDAL.ThemeID.Eq(themeID)).Delete()
	return WrapDBErr(err)
}

func (t *themeServiceImpl) DeleteTheme(ctx context.Context, themeID string, deleteSettings bool) error {
	activatedThemeID, err := t.OptionService.GetActivatedThemeID(ctx)
	if err != nil {
		return err
	}
	if activatedThemeID == themeID {
		return xerr.WithStatus(nil, xerr.StatusBadRequest).WithMsg("can not delete the theme being used")
	}
	if deleteSettings {
		err = t.DeleteThemeSettings(ctx, themeID)
		if err != nil {
			return err
		}
	}
	themeProperty, err := t.GetThemeByID(ctx, themeID)
	if err != nil {
		return err
	}
	err = os.RemoveAll(themeProperty.ThemePath)
	if err != nil {
		return xerr.WithStatus(err, xerr.StatusInternalServerError).WithMsg("delete theme directory err")
	}
	return nil
}

func (t *themeServiceImpl) UploadTheme(ctx context.Context, file *multipart.FileHeader) (*dto.ThemeProperty, error) {
	themeProperty, err := t.ThemeFetchers.MultipartZipThemeFetcher.FetchTheme(ctx, file)
	if err != nil {
		return nil, err
	}
	return t.addTheme(ctx, themeProperty)
}

func (t *themeServiceImpl) UpdateThemeByUpload(ctx context.Context, themeID string, file *multipart.FileHeader) (*dto.ThemeProperty, error) {
	oldThemeProperty, err := t.GetThemeByID(ctx, themeID)
	if err != nil {
		return nil, err
	}
	newThemeProperty, err := t.ThemeFetchers.MultipartZipThemeFetcher.FetchTheme(ctx, file)
	if err != nil {
		return nil, err
	}
	err = os.RemoveAll(oldThemeProperty.ThemePath)
	if err != nil {
		return nil, xerr.WithMsg(err, "delete old theme err").WithStatus(xerr.StatusInternalServerError)
	}
	return t.addTheme(ctx, newThemeProperty)
}

func (t *themeServiceImpl) ReloadTheme(ctx context.Context) error {
	t.Event.Publish(ctx, &event.ThemeUpdateEvent{})
	t.Event.Publish(ctx, &event.ThemeFileUpdatedEvent{})
	return nil
}

func (t *themeServiceImpl) TemplateExist(ctx context.Context, template string) (bool, error) {
	activatedThemeID, err := t.OptionService.GetActivatedThemeID(ctx)
	if err != nil {
		return false, err
	}
	themeProperty, err := t.GetThemeByID(ctx, activatedThemeID)
	if err != nil {
		return false, err
	}
	return util.FileIsExisted(filepath.Join(themeProperty.ThemePath, template)), nil
}

func (t *themeServiceImpl) getThemeConfigItemMap(ctx context.Context, themeID string) (map[string]*dto.ThemeConfigItem, error) {
	themeConfigGroup, err := t.GetThemeConfig(ctx, themeID)
	if err != nil {
		return nil, err
	}

	itemMap := make(map[string]*dto.ThemeConfigItem)
	for _, themeConfig := range themeConfigGroup {
		for _, item := range themeConfig.Items {
			itemMap[item.Name] = item
		}
	}
	return itemMap, nil
}

func (t *themeServiceImpl) addTheme(ctx context.Context, themeProperty *dto.ThemeProperty) (*dto.ThemeProperty, error) {
	existTheme, err := t.PropertyScanner.GetThemeByThemeID(ctx, themeProperty.ID)
	if err != nil {
		return nil, err
	}
	if existTheme != nil {
		return nil, xerr.WithStatus(nil, xerr.StatusBadRequest).WithMsg("theme already exist")
	}
	err = util.CopyDir(themeProperty.ThemePath, filepath.Join(t.Config.Sonic.ThemeDir, themeProperty.FolderName))
	if err != nil {
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest)
	}
	return t.PropertyScanner.ReadThemeProperty(ctx, filepath.Join(t.Config.Sonic.ThemeDir, themeProperty.FolderName))
}

func (t *themeServiceImpl) Render(ctx context.Context, name string) (string, error) {
	activatedThemeID, err := t.OptionService.GetActivatedThemeID(ctx)
	if err != nil {
		return "", err
	}
	return activatedThemeID + "/" + name, nil
}

func (t *themeServiceImpl) Fetch(ctx context.Context, themeURL string) (*dto.ThemeProperty, error) {
	fetchTheme, err := t.ThemeFetchers.GitRepoThemeFetcher.FetchTheme(ctx, themeURL)
	if err != nil {
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg(err.Error())
	}
	return t.addTheme(ctx, fetchTheme)
}
