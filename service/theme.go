package service

import (
	"context"
	"mime/multipart"

	"github.com/go-sonic/sonic/model/dto"
)

type ThemeService interface {
	GetThemeByID(ctx context.Context, themeID string) (*dto.ThemeProperty, error)
	GetActivateTheme(ctx context.Context) (*dto.ThemeProperty, error)
	ListAllTheme(ctx context.Context) ([]*dto.ThemeProperty, error)
	ListThemeFiles(ctx context.Context, themeID string) ([]*dto.ThemeFile, error)
	GetThemeFileContent(ctx context.Context, themeID, absPath string) (string, error)
	UpdateThemeFile(ctx context.Context, themeID, absPath, content string) error
	ListCustomTemplates(ctx context.Context, themeID, prefix string) ([]string, error)
	ActivateTheme(ctx context.Context, themeID string) (*dto.ThemeProperty, error)
	GetThemeConfig(ctx context.Context, themeID string) ([]*dto.ThemeConfigGroup, error)
	GetThemeSettingMap(ctx context.Context, themeID string) (map[string]interface{}, error)
	GetThemeGroupSettingMap(ctx context.Context, themeID, group string) (map[string]interface{}, error)
	SaveThemeSettings(ctx context.Context, themeID string, settings map[string]interface{}) error
	DeleteThemeSettings(ctx context.Context, themeID string) error
	DeleteTheme(ctx context.Context, themeID string, deleteSettings bool) error
	UploadTheme(ctx context.Context, file *multipart.FileHeader) (*dto.ThemeProperty, error)
	UpdateThemeByUpload(ctx context.Context, themeID string, file *multipart.FileHeader) (*dto.ThemeProperty, error)
	ReloadTheme(ctx context.Context) error
	TemplateExist(ctx context.Context, template string) (bool, error)
	Render(ctx context.Context, name string) (string, error)
	Fetch(ctx context.Context, themeURL string) (*dto.ThemeProperty, error)
}
