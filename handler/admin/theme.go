package admin

import (
	"context"
	"errors"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/go-playground/validator/v10"

	"github.com/go-sonic/sonic/consts"
	"github.com/go-sonic/sonic/handler/trans"
	"github.com/go-sonic/sonic/model/param"
	"github.com/go-sonic/sonic/service"
	"github.com/go-sonic/sonic/util"
	"github.com/go-sonic/sonic/util/xerr"
)

type ThemeHandler struct {
	ThemeService  service.ThemeService
	OptionService service.OptionService
}

func NewThemeHandler(l service.ThemeService, o service.OptionService) *ThemeHandler {
	return &ThemeHandler{
		ThemeService:  l,
		OptionService: o,
	}
}

func (t *ThemeHandler) GetThemeByID(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	themeID, err := util.ParamString(_ctx, ctx, "themeID")
	if err != nil {
		return nil, err
	}
	return t.ThemeService.GetThemeByID(_ctx, themeID)
}

func (t *ThemeHandler) ListAllThemes(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	return t.ThemeService.ListAllTheme(_ctx)
}

func (t *ThemeHandler) ListActivatedThemeFile(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	activatedThemeID, err := t.OptionService.GetActivatedThemeID(_ctx)
	if err != nil {
		return nil, err
	}
	return t.ThemeService.ListThemeFiles(_ctx, activatedThemeID)
}

func (t *ThemeHandler) ListThemeFileByID(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	themeID, err := util.ParamString(_ctx, ctx, "themeID")
	if err != nil {
		return nil, err
	}
	return t.ThemeService.ListThemeFiles(_ctx, themeID)
}

func (t *ThemeHandler) GetThemeFileContent(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	path, err := util.MustGetQueryString(_ctx, ctx, "path")
	if err != nil {
		return nil, err
	}
	activatedThemeID, err := t.OptionService.GetActivatedThemeID(_ctx)
	if err != nil {
		return nil, err
	}
	return t.ThemeService.GetThemeFileContent(_ctx, activatedThemeID, path)
}

func (t *ThemeHandler) GetThemeFileContentByID(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	themeID, err := util.ParamString(_ctx, ctx, "themeID")
	if err != nil {
		return nil, err
	}
	path, err := util.MustGetQueryString(_ctx, ctx, "path")
	if err != nil {
		return nil, err
	}

	return t.ThemeService.GetThemeFileContent(_ctx, themeID, path)
}

func (t *ThemeHandler) UpdateThemeFile(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	themeParam := &param.ThemeContent{}
	err := ctx.BindAndValidate(themeParam)
	if err != nil {
		if err != nil {
			e := validator.ValidationErrors{}
			if errors.As(err, &e) {
				return nil, xerr.WithStatus(e, xerr.StatusBadRequest).WithMsg(trans.Translate(e))
			}
			return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("parameter error")
		}
	}
	activatedThemeID, err := t.OptionService.GetActivatedThemeID(_ctx)
	if err != nil {
		return nil, err
	}
	return nil, t.ThemeService.UpdateThemeFile(_ctx, activatedThemeID, themeParam.Path, themeParam.Content)
}

func (t *ThemeHandler) UpdateThemeFileByID(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	themeID, err := util.ParamString(_ctx, ctx, "themeID")
	if err != nil {
		return nil, err
	}
	themeParam := &param.ThemeContent{}
	err = ctx.BindAndValidate(themeParam)
	if err != nil {
		if err != nil {
			e := validator.ValidationErrors{}
			if errors.As(err, &e) {
				return nil, xerr.WithStatus(e, xerr.StatusBadRequest).WithMsg(trans.Translate(e))
			}
			return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("parameter error")
		}
	}
	return nil, t.ThemeService.UpdateThemeFile(_ctx, themeID, themeParam.Path, themeParam.Content)
}

func (t *ThemeHandler) ListCustomSheetTemplate(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	activatedThemeID, err := t.OptionService.GetActivatedThemeID(_ctx)
	if err != nil {
		return nil, err
	}
	return t.ThemeService.ListCustomTemplates(_ctx, activatedThemeID, consts.ThemeCustomSheetPrefix)
}

func (t *ThemeHandler) ListCustomPostTemplate(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	activatedThemeID, err := t.OptionService.GetActivatedThemeID(_ctx)
	if err != nil {
		return nil, err
	}
	return t.ThemeService.ListCustomTemplates(_ctx, activatedThemeID, consts.ThemeCustomPostPrefix)
}

func (t *ThemeHandler) ActivateTheme(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	themeID, err := util.ParamString(_ctx, ctx, "themeID")
	if err != nil {
		return nil, err
	}
	return t.ThemeService.ActivateTheme(_ctx, themeID)
}

func (t *ThemeHandler) GetActivatedTheme(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	activatedThemeID, err := t.OptionService.GetActivatedThemeID(_ctx)
	if err != nil {
		return nil, err
	}
	return t.ThemeService.GetThemeByID(_ctx, activatedThemeID)
}

func (t *ThemeHandler) GetActivatedThemeConfig(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	activatedThemeID, err := t.OptionService.GetActivatedThemeID(_ctx)
	if err != nil {
		return nil, err
	}
	return t.ThemeService.GetThemeConfig(_ctx, activatedThemeID)
}

func (t *ThemeHandler) GetThemeConfigByID(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	themeID, err := util.ParamString(_ctx, ctx, "themeID")
	if err != nil {
		return nil, err
	}
	return t.ThemeService.GetThemeConfig(_ctx, themeID)
}

func (t *ThemeHandler) GetThemeConfigByGroup(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	themeID, err := util.ParamString(_ctx, ctx, "themeID")
	if err != nil {
		return nil, err
	}
	group, err := util.ParamString(_ctx, ctx, "group")
	if err != nil {
		return nil, err
	}
	themeSettings, err := t.ThemeService.GetThemeConfig(_ctx, themeID)
	if err != nil {
		return nil, err
	}
	for _, setting := range themeSettings {
		if setting.Name == group {
			return setting.Items, nil
		}
	}
	return nil, nil
}

func (t *ThemeHandler) GetThemeConfigGroupNames(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	themeID, err := util.ParamString(_ctx, ctx, "themeID")
	if err != nil {
		return nil, err
	}
	themeSettings, err := t.ThemeService.GetThemeConfig(_ctx, themeID)
	if err != nil {
		return nil, err
	}
	groupNames := make([]string, len(themeSettings))
	for index, setting := range themeSettings {
		groupNames[index] = setting.Name
	}
	return groupNames, nil
}

func (t *ThemeHandler) GetActivatedThemeSettingMap(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	activatedThemeID, err := t.OptionService.GetActivatedThemeID(_ctx)
	if err != nil {
		return nil, err
	}
	return t.ThemeService.GetThemeSettingMap(_ctx, activatedThemeID)
}

func (t *ThemeHandler) GetThemeSettingMapByID(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	themeID, err := util.ParamString(_ctx, ctx, "themeID")
	if err != nil {
		return nil, err
	}
	return t.ThemeService.GetThemeSettingMap(_ctx, themeID)
}

func (t *ThemeHandler) GetThemeSettingMapByGroupAndThemeID(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	themeID, err := util.ParamString(_ctx, ctx, "themeID")
	if err != nil {
		return nil, err
	}
	group, err := util.ParamString(_ctx, ctx, "group")
	if err != nil {
		return nil, err
	}
	return t.ThemeService.GetThemeGroupSettingMap(_ctx, themeID, group)
}

func (t *ThemeHandler) SaveActivatedThemeSetting(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	activatedThemeID, err := t.OptionService.GetActivatedThemeID(_ctx)
	if err != nil {
		return nil, err
	}
	settings := make(map[string]interface{})
	err = ctx.BindAndValidate(&settings)
	if err != nil {
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest)
	}
	return nil, t.ThemeService.SaveThemeSettings(_ctx, activatedThemeID, settings)
}

func (t *ThemeHandler) SaveThemeSettingByID(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	themeID, err := util.ParamString(_ctx, ctx, "themeID")
	if err != nil {
		return nil, err
	}
	settings := make(map[string]interface{})
	err = ctx.BindAndValidate(&settings)
	if err != nil {
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest)
	}
	return nil, t.ThemeService.SaveThemeSettings(_ctx, themeID, settings)
}

func (t *ThemeHandler) DeleteThemeByID(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	themeID, err := util.ParamString(_ctx, ctx, "themeID")
	if err != nil {
		return nil, err
	}
	isDeleteSetting, err := util.GetQueryBool(_ctx, ctx, "deleteSettings", false)
	if err != nil {
		return nil, err
	}
	return nil, t.ThemeService.DeleteTheme(_ctx, themeID, isDeleteSetting)
}

func (t *ThemeHandler) UploadTheme(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	fileHeader, err := ctx.FormFile("file")
	if err != nil {
		return nil, xerr.WithMsg(err, "upload theme error").WithStatus(xerr.StatusBadRequest)
	}
	return t.ThemeService.UploadTheme(_ctx, fileHeader)
}

func (t *ThemeHandler) UpdateThemeByUpload(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	themeID, err := util.ParamString(_ctx, ctx, "themeID")
	if err != nil {
		return nil, err
	}
	fileHeader, err := ctx.FormFile("file")
	if err != nil {
		return nil, xerr.WithMsg(err, "upload theme error").WithStatus(xerr.StatusBadRequest)
	}
	return t.ThemeService.UpdateThemeByUpload(_ctx, themeID, fileHeader)
}

func (t *ThemeHandler) FetchTheme(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	uri, _ := util.MustGetQueryString(_ctx, ctx, "uri")
	return t.ThemeService.Fetch(_ctx, uri)
}

func (t *ThemeHandler) UpdateThemeByFetching(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	return nil, xerr.WithMsg(nil, "not support").WithStatus(xerr.StatusInternalServerError)
}

func (t *ThemeHandler) ReloadTheme(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	return nil, t.ThemeService.ReloadTheme(_ctx)
}

func (t *ThemeHandler) TemplateExist(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	template, err := util.MustGetQueryString(_ctx, ctx, "template")
	if err != nil {
		return nil, err
	}
	return t.ThemeService.TemplateExist(_ctx, template)
}
