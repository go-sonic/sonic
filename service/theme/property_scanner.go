package theme

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/go-sonic/sonic/config"
	"github.com/go-sonic/sonic/consts"
	"github.com/go-sonic/sonic/model/dto"
	"github.com/go-sonic/sonic/util/xerr"
)

type PropertyScanner interface {
	ListAll(ctx context.Context, themeRootPath string) ([]*dto.ThemeProperty, error)
	GetThemeByThemeID(ctx context.Context, themeID string) (*dto.ThemeProperty, error)
	ReadThemeProperty(ctx context.Context, themePath string) (*dto.ThemeProperty, error)
	ReadThemeConfig(ctx context.Context, themePath string) ([]*dto.ThemeConfigGroup, error)
	UnmarshalProperty(ctx context.Context, themePropertyContent []byte) (*dto.ThemeProperty, error)
	UnmarshalConfig(ctx context.Context, themeSettingContent []byte) ([]*dto.ThemeConfigGroup, error)
}

type propertyScannerImpl struct {
	Config *config.Config
}

func NewPropertyScanner(config *config.Config) PropertyScanner {
	return &propertyScannerImpl{
		Config: config,
	}
}

func (s *propertyScannerImpl) GetThemeByThemeID(ctx context.Context, themeID string) (*dto.ThemeProperty, error) {
	allThemes, err := s.ListAll(ctx, s.Config.Sonic.ThemeDir)
	if err != nil {
		return nil, err
	}

	for _, themeProperty := range allThemes {
		if themeProperty.ID == themeID {
			return themeProperty, nil
		}
	}
	return nil, nil
}

func (s *propertyScannerImpl) ListAll(ctx context.Context, themeRootPath string) ([]*dto.ThemeProperty, error) {
	themeRootDir, err := os.Open(themeRootPath)
	defer themeRootDir.Close()
	if err != nil {
		return nil, xerr.NoType.Wrap(err)
	}

	themeDirs, err := themeRootDir.ReadDir(0)
	result := make([]*dto.ThemeProperty, 0)

	for _, themeDir := range themeDirs {
		if themeDir.IsDir() {
			themeProperty, err := s.ReadThemeProperty(ctx, filepath.Join(themeRootPath, themeDir.Name()))
			if err != nil {
				return nil, err
			}
			result = append(result, themeProperty)
		}
	}

	return result, nil
}

func (s *propertyScannerImpl) ReadThemeProperty(ctx context.Context, themePath string) (*dto.ThemeProperty, error) {
	var (
		err               error
		themePropertyFile *os.File
		fileStat          os.FileInfo
	)
	for _, themePropertyFilename := range consts.ThemePropertyFilenames {
		themePropertyFile, err = os.Open(filepath.Join(themePath, themePropertyFilename))
		defer themePropertyFile.Close()
		if os.IsNotExist(err) {
			continue
		}
		fileStat, err = themePropertyFile.Stat()
		if err != nil {
			continue
		}
		if fileStat.IsDir() {
			err = xerr.WithErrMsgf(err, "%s is dir", fileStat.Name())
			continue
		}
		break
	}
	if err != nil {
		return nil, xerr.NoType.Wrapf(err, "theme.yaml not exist")
	}
	propertyContent, err := io.ReadAll(themePropertyFile)
	if err != nil {
		return nil, xerr.WithMsg(err, "read theme property file err")
	}
	themeProperty, err := s.UnmarshalProperty(ctx, propertyContent)
	if err != nil {
		return nil, err
	}
	themeProperty.ThemePath = themePath
	themeProperty.FolderName = filepath.Base(themePath)
	themeProperty.Activated = false
	hasOptions, _ := s.hasSettingFile(ctx, themePath)
	themeProperty.HasOptions = hasOptions
	screenshotFilename, err := s.GetThemeScreenshotAbsPath(ctx, themePath)
	if err != nil {
		return nil, err
	}
	themeProperty.ScreenShots = strings.Join([]string{"/themes", filepath.Base(themePath), screenshotFilename}, "/")
	return themeProperty, nil
}

func (s *propertyScannerImpl) UnmarshalProperty(ctx context.Context, themePropertyContent []byte) (*dto.ThemeProperty, error) {
	property := dto.ThemeProperty{}
	err := yaml.Unmarshal(themePropertyContent, &property)
	if err != nil {
		return nil, xerr.WithMsg(err, "unmarshal yaml file err")
	}
	return &property, nil
}

func (s *propertyScannerImpl) GetThemeScreenshotAbsPath(ctx context.Context, themePath string) (string, error) {
	themeDir, err := os.Open(themePath)
	defer themeDir.Close()
	if err != nil {
		return "", xerr.NoType.Wrapf(err, "open theme path error themePath=%s", themePath)
	}
	themeFiles, err := themeDir.ReadDir(0)
	if err != nil {
		return "", xerr.NoType.Wrapf(err, "read theme file error themePath=%s", themePath)
	}
	for _, themeFile := range themeFiles {
		if themeFile.IsDir() {
			continue
		}
		if !themeFile.Type().IsRegular() {
			continue
		}
		if strings.HasPrefix(themeFile.Name(), consts.ThemeScreenshotsName) {
			return themeFile.Name(), nil
		}
	}
	return "", nil
}

func (s *propertyScannerImpl) hasSettingFile(ctx context.Context, themePath string) (bool, error) {
	var (
		err                  error
		themeSettingFileInfo os.FileInfo
	)
	for _, themeSettingFilename := range consts.ThemeSettingFilenames {
		themeSettingFileInfo, err = os.Stat(filepath.Join(themePath, themeSettingFilename))
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			continue
		}
		if themeSettingFileInfo.IsDir() {
			err = xerr.WithErrMsgf(err, "%s is dir", themeSettingFileInfo.Name())
			continue
		}
		return true, nil
	}
	return false, err
}

func (s *propertyScannerImpl) ReadThemeConfig(ctx context.Context, themePath string) ([]*dto.ThemeConfigGroup, error) {
	var (
		err              error
		themeSettingFile *os.File
		fileStat         os.FileInfo
	)
	for _, themeSettingFilename := range consts.ThemeSettingFilenames {
		themeSettingFile, err = os.Open(filepath.Join(themePath, themeSettingFilename))
		defer themeSettingFile.Close()
		if os.IsNotExist(err) {
			continue
		}
		fileStat, err = themeSettingFile.Stat()
		if err != nil {
			continue
		}
		if fileStat.IsDir() {
			err = xerr.WithErrMsgf(err, "%s is dir", fileStat.Name())
			continue
		}
		break
	}
	if err != nil {
		return nil, xerr.NoType.Wrapf(err, "setting.yaml not exist")
	}
	settingContent, err := io.ReadAll(themeSettingFile)
	if err != nil {
		return nil, xerr.WithMsg(err, "read theme setting file err")
	}
	themeSetting, err := s.UnmarshalConfig(ctx, settingContent)
	if err != nil {
		return nil, err
	}
	return themeSetting, nil
}

func (s *propertyScannerImpl) UnmarshalConfig(ctx context.Context, themeSettingContent []byte) ([]*dto.ThemeConfigGroup, error) {
	settingMap := make(map[string]*dto.ThemeConfigGroup, 0)
	err := yaml.Unmarshal(themeSettingContent, &settingMap)
	if err != nil {
		return nil, xerr.WithMsg(err, "unmarshal yaml file err")
	}

	settings := make([]*dto.ThemeConfigGroup, 0, len(settingMap))
	for name, setting := range settingMap {
		setting.Name = name
		for _, item := range setting.ItemMap {
			setting.Items = append(setting.Items, item)
		}

		settings = append(settings, setting)
	}
	return settings, nil
}
