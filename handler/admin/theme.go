package admin

import (
	"errors"

	"github.com/gin-gonic/gin"
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

func (t *ThemeHandler) GetThemeByID(ctx *gin.Context) (interface{}, error) {
	themeID, err := util.ParamString(ctx, "themeID")
	if err != nil {
		return nil, err
	}
	return t.ThemeService.GetThemeByID(ctx, themeID)
}

func (t *ThemeHandler) ListAllThemes(ctx *gin.Context) (interface{}, error) {
	return t.ThemeService.ListAllTheme(ctx)
}

func (t *ThemeHandler) ListActivatedThemeFile(ctx *gin.Context) (interface{}, error) {
	activatedThemeID, err := t.OptionService.GetActivatedThemeID(ctx)
	if err != nil {
		return nil, err
	}
	return t.ThemeService.ListThemeFiles(ctx, activatedThemeID)
}

func (t *ThemeHandler) ListThemeFileByID(ctx *gin.Context) (interface{}, error) {
	themeID, err := util.ParamString(ctx, "themeID")
	if err != nil {
		return nil, err
	}
	return t.ThemeService.ListThemeFiles(ctx, themeID)
}

func (t *ThemeHandler) GetThemeFileContent(ctx *gin.Context) (interface{}, error) {
	path, err := util.MustGetQueryString(ctx, "path")
	if err != nil {
		return nil, err
	}
	activatedThemeID, err := t.OptionService.GetActivatedThemeID(ctx)
	if err != nil {
		return nil, err
	}
	return t.ThemeService.GetThemeFileContent(ctx, activatedThemeID, path)
}

func (t *ThemeHandler) GetThemeFileContentByID(ctx *gin.Context) (interface{}, error) {
	themeID, err := util.ParamString(ctx, "themeID")
	if err != nil {
		return nil, err
	}
	path, err := util.MustGetQueryString(ctx, "path")
	if err != nil {
		return nil, err
	}

	return t.ThemeService.GetThemeFileContent(ctx, themeID, path)
}

func (t *ThemeHandler) UpdateThemeFile(ctx *gin.Context) (interface{}, error) {
	themeParam := &param.ThemeContent{}
	err := ctx.ShouldBindJSON(themeParam)
	if err != nil {
		if err != nil {
			e := validator.ValidationErrors{}
			if errors.As(err, &e) {
				return nil, xerr.WithStatus(e, xerr.StatusBadRequest).WithMsg(trans.Translate(e))
			}
			return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("parameter error")
		}
	}
	activatedThemeID, err := t.OptionService.GetActivatedThemeID(ctx)
	if err != nil {
		return nil, err
	}
	return nil, t.ThemeService.UpdateThemeFile(ctx, activatedThemeID, themeParam.Path, themeParam.Content)
}

func (t *ThemeHandler) UpdateThemeFileByID(ctx *gin.Context) (interface{}, error) {
	themeID, err := util.ParamString(ctx, "themeID")
	if err != nil {
		return nil, err
	}
	themeParam := &param.ThemeContent{}
	err = ctx.ShouldBindJSON(themeParam)
	if err != nil {
		if err != nil {
			e := validator.ValidationErrors{}
			if errors.As(err, &e) {
				return nil, xerr.WithStatus(e, xerr.StatusBadRequest).WithMsg(trans.Translate(e))
			}
			return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("parameter error")
		}
	}
	return nil, t.ThemeService.UpdateThemeFile(ctx, themeID, themeParam.Path, themeParam.Content)
}

func (t *ThemeHandler) ListCustomSheetTemplate(ctx *gin.Context) (interface{}, error) {
	activatedThemeID, err := t.OptionService.GetActivatedThemeID(ctx)
	if err != nil {
		return nil, err
	}
	return t.ThemeService.ListCustomTemplates(ctx, activatedThemeID, consts.ThemeCustomSheetPrefix)
}

func (t *ThemeHandler) ListCustomPostTemplate(ctx *gin.Context) (interface{}, error) {
	activatedThemeID, err := t.OptionService.GetActivatedThemeID(ctx)
	if err != nil {
		return nil, err
	}
	return t.ThemeService.ListCustomTemplates(ctx, activatedThemeID, consts.ThemeCustomPostPrefix)
}

func (t *ThemeHandler) ActivateTheme(ctx *gin.Context) (interface{}, error) {
	themeID, err := util.ParamString(ctx, "themeID")
	if err != nil {
		return nil, err
	}
	return t.ThemeService.ActivateTheme(ctx, themeID)
}

func (t *ThemeHandler) GetActivatedTheme(ctx *gin.Context) (interface{}, error) {
	activatedThemeID, err := t.OptionService.GetActivatedThemeID(ctx)
	if err != nil {
		return nil, err
	}
	return t.ThemeService.GetThemeByID(ctx, activatedThemeID)
}

func (t *ThemeHandler) GetActivatedThemeConfig(ctx *gin.Context) (interface{}, error) {
	activatedThemeID, err := t.OptionService.GetActivatedThemeID(ctx)
	if err != nil {
		return nil, err
	}
	return t.ThemeService.GetThemeConfig(ctx, activatedThemeID)
}

func (t *ThemeHandler) GetThemeConfigByID(ctx *gin.Context) (interface{}, error) {
	themeID, err := util.ParamString(ctx, "themeID")
	if err != nil {
		return nil, err
	}
	return t.ThemeService.GetThemeConfig(ctx, themeID)
}

func (t *ThemeHandler) GetThemeConfigByGroup(ctx *gin.Context) (interface{}, error) {
	themeID, err := util.ParamString(ctx, "themeID")
	if err != nil {
		return nil, err
	}
	group, err := util.ParamString(ctx, "group")
	if err != nil {
		return nil, err
	}
	themeSettings, err := t.ThemeService.GetThemeConfig(ctx, themeID)
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

func (t *ThemeHandler) GetThemeConfigGroupNames(ctx *gin.Context) (interface{}, error) {
	themeID, err := util.ParamString(ctx, "themeID")
	if err != nil {
		return nil, err
	}
	themeSettings, err := t.ThemeService.GetThemeConfig(ctx, themeID)
	if err != nil {
		return nil, err
	}
	groupNames := make([]string, len(themeSettings))
	for index, setting := range themeSettings {
		groupNames[index] = setting.Name
	}
	return groupNames, nil
}

func (t *ThemeHandler) GetActivatedThemeSettingMap(ctx *gin.Context) (interface{}, error) {
	activatedThemeID, err := t.OptionService.GetActivatedThemeID(ctx)
	if err != nil {
		return nil, err
	}
	return t.ThemeService.GetThemeSettingMap(ctx, activatedThemeID)
}

func (t *ThemeHandler) GetThemeSettingMapByID(ctx *gin.Context) (interface{}, error) {
	themeID, err := util.ParamString(ctx, "themeID")
	if err != nil {
		return nil, err
	}
	return t.ThemeService.GetThemeSettingMap(ctx, themeID)
}

func (t *ThemeHandler) GetThemeSettingMapByGroupAndThemeID(ctx *gin.Context) (interface{}, error) {
	themeID, err := util.ParamString(ctx, "themeID")
	if err != nil {
		return nil, err
	}
	group, err := util.ParamString(ctx, "group")
	if err != nil {
		return nil, err
	}
	return t.ThemeService.GetThemeGroupSettingMap(ctx, themeID, group)
}

func (t *ThemeHandler) SaveActivatedThemeSetting(ctx *gin.Context) (interface{}, error) {
	activatedThemeID, err := t.OptionService.GetActivatedThemeID(ctx)
	if err != nil {
		return nil, err
	}
	settings := make(map[string]interface{})
	err = ctx.ShouldBindJSON(&settings)
	if err != nil {
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest)
	}
	return nil, t.ThemeService.SaveThemeSettings(ctx, activatedThemeID, settings)
}

func (t *ThemeHandler) SaveThemeSettingByID(ctx *gin.Context) (interface{}, error) {
	themeID, err := util.ParamString(ctx, "themeID")
	if err != nil {
		return nil, err
	}
	settings := make(map[string]interface{})
	err = ctx.ShouldBindJSON(&settings)
	if err != nil {
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest)
	}
	return nil, t.ThemeService.SaveThemeSettings(ctx, themeID, settings)
}

func (t *ThemeHandler) DeleteThemeByID(ctx *gin.Context) (interface{}, error) {
	themeID, err := util.ParamString(ctx, "themeID")
	if err != nil {
		return nil, err
	}
	isDeleteSetting, err := util.GetQueryBool(ctx, "deleteSettings", false)
	if err != nil {
		return nil, err
	}
	return nil, t.ThemeService.DeleteTheme(ctx, themeID, isDeleteSetting)
}

func (t *ThemeHandler) UploadTheme(ctx *gin.Context) (interface{}, error) {
	fileHeader, err := ctx.FormFile("file")
	if err != nil {
		return nil, xerr.WithMsg(err, "upload theme error").WithStatus(xerr.StatusBadRequest)
	}
	return t.ThemeService.UploadTheme(ctx, fileHeader)
}

func (t *ThemeHandler) UpdateThemeByUpload(ctx *gin.Context) (interface{}, error) {
	themeID, err := util.ParamString(ctx, "themeID")
	if err != nil {
		return nil, err
	}
	fileHeader, err := ctx.FormFile("file")
	if err != nil {
		return nil, xerr.WithMsg(err, "upload theme error").WithStatus(xerr.StatusBadRequest)
	}
	return t.ThemeService.UpdateThemeByUpload(ctx, themeID, fileHeader)
}

func (t *ThemeHandler) FetchTheme(ctx *gin.Context) (interface{}, error) {
	uri, _ := util.MustGetQueryString(ctx, "uri")
	return t.ThemeService.Fetch(ctx, uri)
}

func (t *ThemeHandler) UpdateThemeByFetching(ctx *gin.Context) (interface{}, error) {
	return nil, xerr.WithMsg(nil, "not support").WithStatus(xerr.StatusInternalServerError)
}

func (t *ThemeHandler) ReloadTheme(ctx *gin.Context) (interface{}, error) {
	return nil, t.ThemeService.ReloadTheme(ctx)
}

func (t *ThemeHandler) TemplateExist(ctx *gin.Context) (interface{}, error) {
	template, err := util.MustGetQueryString(ctx, "template")
	if err != nil {
		return nil, err
	}
	return t.ThemeService.TemplateExist(ctx, template)
}
