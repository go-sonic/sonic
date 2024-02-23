package admin

import (
	"errors"

	"github.com/gin-gonic/gin"
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

func (m *MenuHandler) ListMenus(ctx *gin.Context) (interface{}, error) {
	sort := param.Sort{}
	err := ctx.ShouldBindQuery(&sort)
	if err != nil {
		return nil, xerr.WithMsg(err, "sort parameter error").WithStatus(xerr.StatusBadRequest)
	}
	if len(sort.Fields) == 0 {
		sort.Fields = append(sort.Fields, "team,desc", "priority,asc")
	} else {
		sort.Fields = append(sort.Fields, "priority,asc")
	}
	menus, err := m.MenuService.List(ctx, &sort)
	if err != nil {
		return nil, err
	}
	return m.MenuService.ConvertToDTOs(ctx, menus), nil
}

func (m *MenuHandler) ListMenusAsTree(ctx *gin.Context) (interface{}, error) {
	sort := param.Sort{}
	err := ctx.ShouldBindQuery(&sort)
	if err != nil {
		return nil, xerr.WithMsg(err, "sort parameter error").WithStatus(xerr.StatusBadRequest)
	}
	if len(sort.Fields) == 0 {
		sort.Fields = append(sort.Fields, "team,desc", "priority,asc")
	} else {
		sort.Fields = append(sort.Fields, "priority,asc")
	}
	menus, err := m.MenuService.ListAsTree(ctx, &sort)
	if err != nil {
		return nil, err
	}
	return menus, nil
}

func (m *MenuHandler) ListMenusAsTreeByTeam(ctx *gin.Context) (interface{}, error) {
	sort := param.Sort{}
	err := ctx.ShouldBindQuery(&sort)
	if err != nil {
		return nil, xerr.WithMsg(err, "sort parameter error").WithStatus(xerr.StatusBadRequest)
	}
	if len(sort.Fields) == 0 {
		sort.Fields = append(sort.Fields, "priority,asc")
	}
	team, _ := util.MustGetQueryString(ctx, "team")
	if team == "" {
		menus, err := m.MenuService.ListAsTree(ctx, &sort)
		if err != nil {
			return nil, err
		}
		return menus, nil
	}
	menus, err := m.MenuService.ListAsTreeByTeam(ctx, team, &sort)
	if err != nil {
		return nil, err
	}
	return menus, nil
}

func (m *MenuHandler) GetMenuByID(ctx *gin.Context) (interface{}, error) {
	id, err := util.ParamInt32(ctx, "id")
	if err != nil {
		return nil, err
	}
	menu, err := m.MenuService.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return m.MenuService.ConvertToDTO(ctx, menu), nil
}

func (m *MenuHandler) CreateMenu(ctx *gin.Context) (interface{}, error) {
	menuParam := &param.Menu{}
	err := ctx.ShouldBindJSON(menuParam)
	if err != nil {
		e := validator.ValidationErrors{}
		if errors.As(err, &e) {
			return nil, xerr.WithStatus(e, xerr.StatusBadRequest).WithMsg(trans.Translate(e))
		}
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("parameter error")
	}
	menu, err := m.MenuService.Create(ctx, menuParam)
	if err != nil {
		return nil, err
	}
	return m.MenuService.ConvertToDTO(ctx, menu), nil
}

func (m *MenuHandler) CreateMenuBatch(ctx *gin.Context) (interface{}, error) {
	menuParams := make([]*param.Menu, 0)
	err := ctx.ShouldBindJSON(&menuParams)
	if err != nil {
		e := validator.ValidationErrors{}
		if errors.As(err, &e) {
			return nil, xerr.WithStatus(e, xerr.StatusBadRequest).WithMsg(trans.Translate(e))
		}
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("parameter error")
	}
	menus, err := m.MenuService.CreateBatch(ctx, menuParams)
	if err != nil {
		return nil, err
	}
	return m.MenuService.ConvertToDTOs(ctx, menus), nil
}

func (m *MenuHandler) UpdateMenu(ctx *gin.Context) (interface{}, error) {
	id, err := util.ParamInt32(ctx, "id")
	if err != nil {
		return nil, err
	}
	menuParam := &param.Menu{}
	err = ctx.ShouldBindJSON(menuParam)
	if err != nil {
		e := validator.ValidationErrors{}
		if errors.As(err, &e) {
			return nil, xerr.WithStatus(e, xerr.StatusBadRequest).WithMsg(trans.Translate(e))
		}
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("parameter error")
	}
	menu, err := m.MenuService.Update(ctx, id, menuParam)
	if err != nil {
		return nil, err
	}
	return m.MenuService.ConvertToDTO(ctx, menu), nil
}

func (m *MenuHandler) UpdateMenuBatch(ctx *gin.Context) (interface{}, error) {
	menuParams := make([]*param.Menu, 0)
	err := ctx.ShouldBindJSON(&menuParams)
	if err != nil {
		e := validator.ValidationErrors{}
		if errors.As(err, &e) {
			return nil, xerr.WithStatus(e, xerr.StatusBadRequest).WithMsg(trans.Translate(e))
		}
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("parameter error")
	}
	menus, err := m.MenuService.UpdateBatch(ctx, menuParams)
	if err != nil {
		return nil, err
	}
	return m.MenuService.ConvertToDTOs(ctx, menus), nil
}

func (m *MenuHandler) DeleteMenu(ctx *gin.Context) (interface{}, error) {
	id, err := util.ParamInt32(ctx, "id")
	if err != nil {
		return nil, err
	}
	return nil, m.MenuService.Delete(ctx, id)
}

func (m *MenuHandler) DeleteMenuBatch(ctx *gin.Context) (interface{}, error) {
	menuIDs := make([]int32, 0)
	err := ctx.ShouldBind(&menuIDs)
	if err != nil {
		return nil, xerr.WithMsg(err, "menuIDs error").WithStatus(xerr.StatusBadRequest)
	}
	return nil, m.MenuService.DeleteBatch(ctx, menuIDs)
}

func (m *MenuHandler) ListMenuTeams(ctx *gin.Context) (interface{}, error) {
	return m.MenuService.ListTeams(ctx)
}
