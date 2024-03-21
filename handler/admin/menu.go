package admin

import (
	"context"
	"errors"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/go-playground/validator/v10"

	"github.com/go-sonic/sonic/handler/trans"
	"github.com/go-sonic/sonic/model/param"
	"github.com/go-sonic/sonic/service"
	"github.com/go-sonic/sonic/util"
	"github.com/go-sonic/sonic/util/xerr"
)

type MenuHandler struct {
	MenuService service.MenuService
}

func NewMenuHandler(menuService service.MenuService) *MenuHandler {
	return &MenuHandler{
		MenuService: menuService,
	}
}

func (m *MenuHandler) ListMenus(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	sort := param.Sort{}
	err := ctx.BindAndValidate(&sort)
	if err != nil {
		return nil, xerr.WithMsg(err, "sort parameter error").WithStatus(xerr.StatusBadRequest)
	}
	if len(sort.Fields) == 0 {
		sort.Fields = append(sort.Fields, "team,desc", "priority,asc")
	} else {
		sort.Fields = append(sort.Fields, "priority,asc")
	}
	menus, err := m.MenuService.List(_ctx, &sort)
	if err != nil {
		return nil, err
	}
	return m.MenuService.ConvertToDTOs(_ctx, menus), nil
}

func (m *MenuHandler) ListMenusAsTree(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	sort := param.Sort{}
	err := ctx.BindAndValidate(&sort)
	if err != nil {
		return nil, xerr.WithMsg(err, "sort parameter error").WithStatus(xerr.StatusBadRequest)
	}
	if len(sort.Fields) == 0 {
		sort.Fields = append(sort.Fields, "team,desc", "priority,asc")
	} else {
		sort.Fields = append(sort.Fields, "priority,asc")
	}
	menus, err := m.MenuService.ListAsTree(_ctx, &sort)
	if err != nil {
		return nil, err
	}
	return menus, nil
}

func (m *MenuHandler) ListMenusAsTreeByTeam(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	sort := param.Sort{}
	err := ctx.BindAndValidate(&sort)
	if err != nil {
		return nil, xerr.WithMsg(err, "sort parameter error").WithStatus(xerr.StatusBadRequest)
	}
	if len(sort.Fields) == 0 {
		sort.Fields = append(sort.Fields, "priority,asc")
	}
	team, _ := util.MustGetQueryString(_ctx, ctx, "team")
	if team == "" {
		menus, err := m.MenuService.ListAsTree(_ctx, &sort)
		if err != nil {
			return nil, err
		}
		return menus, nil
	}
	menus, err := m.MenuService.ListAsTreeByTeam(_ctx, team, &sort)
	if err != nil {
		return nil, err
	}
	return menus, nil
}

func (m *MenuHandler) GetMenuByID(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	id, err := util.ParamInt32(_ctx, ctx, "id")
	if err != nil {
		return nil, err
	}
	menu, err := m.MenuService.GetByID(_ctx, id)
	if err != nil {
		return nil, err
	}
	return m.MenuService.ConvertToDTO(_ctx, menu), nil
}

func (m *MenuHandler) CreateMenu(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	menuParam := &param.Menu{}
	err := ctx.BindAndValidate(menuParam)
	if err != nil {
		e := validator.ValidationErrors{}
		if errors.As(err, &e) {
			return nil, xerr.WithStatus(e, xerr.StatusBadRequest).WithMsg(trans.Translate(e))
		}
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("parameter error")
	}
	menu, err := m.MenuService.Create(_ctx, menuParam)
	if err != nil {
		return nil, err
	}
	return m.MenuService.ConvertToDTO(_ctx, menu), nil
}

func (m *MenuHandler) CreateMenuBatch(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	menuParams := make([]*param.Menu, 0)
	err := ctx.BindAndValidate(&menuParams)
	if err != nil {
		e := validator.ValidationErrors{}
		if errors.As(err, &e) {
			return nil, xerr.WithStatus(e, xerr.StatusBadRequest).WithMsg(trans.Translate(e))
		}
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("parameter error")
	}
	menus, err := m.MenuService.CreateBatch(_ctx, menuParams)
	if err != nil {
		return nil, err
	}
	return m.MenuService.ConvertToDTOs(_ctx, menus), nil
}

func (m *MenuHandler) UpdateMenu(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	id, err := util.ParamInt32(_ctx, ctx, "id")
	if err != nil {
		return nil, err
	}
	menuParam := &param.Menu{}
	err = ctx.BindAndValidate(menuParam)
	if err != nil {
		e := validator.ValidationErrors{}
		if errors.As(err, &e) {
			return nil, xerr.WithStatus(e, xerr.StatusBadRequest).WithMsg(trans.Translate(e))
		}
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("parameter error")
	}
	menu, err := m.MenuService.Update(_ctx, id, menuParam)
	if err != nil {
		return nil, err
	}
	return m.MenuService.ConvertToDTO(_ctx, menu), nil
}

func (m *MenuHandler) UpdateMenuBatch(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	menuParams := make([]*param.Menu, 0)
	err := ctx.BindAndValidate(&menuParams)
	if err != nil {
		e := validator.ValidationErrors{}
		if errors.As(err, &e) {
			return nil, xerr.WithStatus(e, xerr.StatusBadRequest).WithMsg(trans.Translate(e))
		}
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("parameter error")
	}
	menus, err := m.MenuService.UpdateBatch(_ctx, menuParams)
	if err != nil {
		return nil, err
	}
	return m.MenuService.ConvertToDTOs(_ctx, menus), nil
}

func (m *MenuHandler) DeleteMenu(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	id, err := util.ParamInt32(_ctx, ctx, "id")
	if err != nil {
		return nil, err
	}
	return nil, m.MenuService.Delete(_ctx, id)
}

func (m *MenuHandler) DeleteMenuBatch(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	menuIDs := make([]int32, 0)
	err := ctx.BindAndValidate(&menuIDs)
	if err != nil {
		return nil, xerr.WithMsg(err, "menuIDs error").WithStatus(xerr.StatusBadRequest)
	}
	return nil, m.MenuService.DeleteBatch(_ctx, menuIDs)
}

func (m *MenuHandler) ListMenuTeams(_ctx context.Context, ctx *app.RequestContext) (interface{}, error) {
	return m.MenuService.ListTeams(_ctx)
}
