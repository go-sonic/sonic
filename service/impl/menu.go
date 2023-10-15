package impl

import (
	"context"

	"github.com/go-sonic/sonic/dal"
	"github.com/go-sonic/sonic/model/dto"
	"github.com/go-sonic/sonic/model/entity"
	"github.com/go-sonic/sonic/model/param"
	"github.com/go-sonic/sonic/model/vo"
	"github.com/go-sonic/sonic/service"
	"github.com/go-sonic/sonic/util/xerr"
)

type menuServiceImpl struct{}

func NewMenuService() service.MenuService {
	return &menuServiceImpl{}
}

func (m *menuServiceImpl) DeleteBatch(ctx context.Context, ids []int32) error {
	menuDAL := dal.GetQueryByCtx(ctx).Menu
	_, err := menuDAL.WithContext(ctx).Where(menuDAL.ID.In(ids...)).Delete()
	if err != nil {
		return WrapDBErr(err)
	}
	return nil
}

func (m *menuServiceImpl) UpdateBatch(ctx context.Context, menuParams []*param.Menu) ([]*entity.Menu, error) {
	ids := make([]int32, 0, len(menuParams))
	err := dal.GetQueryByCtx(ctx).Transaction(func(tx *dal.Query) error {
		menuDAL := tx.Menu
		for _, menuParam := range menuParams {
			ids = append(ids, menuParam.ID)
			_, err := menuDAL.WithContext(ctx).Where(menuDAL.ID.Eq(menuParam.ID)).UpdateSimple(
				menuDAL.Team.Value(menuParam.Team),
				menuDAL.Priority.Value(menuParam.Priority),
				menuDAL.Name.Value(menuParam.Name),
				menuDAL.URL.Value(menuParam.URL),
				menuDAL.Target.Value(menuParam.Target),
				menuDAL.Icon.Value(menuParam.Icon),
				menuDAL.ParentID.Value(menuParam.ParentID),
			)
			if err != nil {
				return WrapDBErr(err)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	menuDAL := dal.GetQueryByCtx(ctx).Menu
	menus, err := menuDAL.WithContext(ctx).Where(menuDAL.ID.In(ids...)).Find()
	if err != nil {
		return nil, err
	}
	return menus, nil
}

func (m *menuServiceImpl) CreateBatch(ctx context.Context, menuParams []*param.Menu) ([]*entity.Menu, error) {
	menuDAL := dal.GetQueryByCtx(ctx).Menu
	menus := make([]*entity.Menu, 0, len(menuParams))
	for _, menuParam := range menuParams {
		menu := &entity.Menu{
			Name:     menuParam.Name,
			URL:      menuParam.URL,
			Icon:     menuParam.Icon,
			Priority: menuParam.Priority,
			Team:     menuParam.Team,
			ParentID: menuParam.ParentID,
			Target:   menuParam.Target,
		}
		menus = append(menus, menu)
	}
	err := menuDAL.WithContext(ctx).Create(menus...)
	if err != nil {
		return nil, WrapDBErr(err)
	}
	return menus, nil
}

func (m *menuServiceImpl) ListAsTree(ctx context.Context, sort *param.Sort) ([]*vo.Menu, error) {
	allMenus, err := m.List(ctx, sort)
	if err != nil {
		return nil, err
	}
	return m.buildTree(ctx, allMenus), nil
}

func (m *menuServiceImpl) List(ctx context.Context, sort *param.Sort) ([]*entity.Menu, error) {
	menuDAL := dal.GetQueryByCtx(ctx).Menu
	menuDO := menuDAL.WithContext(ctx)
	err := BuildSort(sort, &menuDAL, &menuDO)
	if err != nil {
		return nil, err
	}
	menus, err := menuDO.Find()
	if err != nil {
		return nil, WrapDBErr(err)
	}
	return menus, nil
}

func (m *menuServiceImpl) ListByTeam(ctx context.Context, team string, sort *param.Sort) ([]*entity.Menu, error) {
	menuDAL := dal.GetQueryByCtx(ctx).Menu
	menuDO := menuDAL.WithContext(ctx)
	err := BuildSort(sort, &menuDAL, &menuDO)
	if err != nil {
		return nil, err
	}
	menus, err := menuDO.Where(menuDAL.Team.Eq(team)).Find()
	if err != nil {
		return nil, WrapDBErr(err)
	}
	return menus, nil
}

func (m *menuServiceImpl) ListAsTreeByTeam(ctx context.Context, team string, sort *param.Sort) ([]*vo.Menu, error) {
	menuDAL := dal.GetQueryByCtx(ctx).Menu
	menuDO := menuDAL.WithContext(ctx)
	err := BuildSort(sort, &menuDAL, &menuDO)
	if err != nil {
		return nil, err
	}
	menus, err := menuDO.Where(menuDAL.Team.Eq(team)).Find()
	if err != nil {
		return nil, err
	}
	return m.buildTree(ctx, menus), nil
}

func (m *menuServiceImpl) GetByID(ctx context.Context, id int32) (*entity.Menu, error) {
	menuDAL := dal.GetQueryByCtx(ctx).Menu
	menu, err := menuDAL.WithContext(ctx).Where(menuDAL.ID.Eq(id)).First()
	if err != nil {
		return nil, WrapDBErr(err)
	}
	return menu, nil
}

func (m *menuServiceImpl) Create(ctx context.Context, menuParam *param.Menu) (*entity.Menu, error) {
	menu := &entity.Menu{
		Name:     menuParam.Name,
		URL:      menuParam.URL,
		Icon:     menuParam.Icon,
		Priority: menuParam.Priority,
		Team:     menuParam.Team,
		ParentID: menuParam.ParentID,
		Target:   menuParam.Target,
	}
	menuDAL := dal.GetQueryByCtx(ctx).Menu
	err := menuDAL.WithContext(ctx).Create(menu)
	if err != nil {
		return nil, WrapDBErr(err)
	}
	return menu, nil
}

func (m *menuServiceImpl) Update(ctx context.Context, id int32, menuParam *param.Menu) (*entity.Menu, error) {
	menuDAL := dal.GetQueryByCtx(ctx).Menu
	updateResult, err := menuDAL.WithContext(ctx).Where(menuDAL.ID.Eq(id)).UpdateSimple(
		menuDAL.Team.Value(menuParam.Team),
		menuDAL.Priority.Value(menuParam.Priority),
		menuDAL.Name.Value(menuParam.Name),
		menuDAL.URL.Value(menuParam.URL),
		menuDAL.Target.Value(menuParam.Target),
		menuDAL.Icon.Value(menuParam.Icon),
		menuDAL.ParentID.Value(menuParam.ParentID),
	)
	if err != nil {
		return nil, WrapDBErr(err)
	}
	if updateResult.RowsAffected != 1 {
		return nil, xerr.NoType.New("update menu failed").WithMsg("update menu failed").WithStatus(xerr.StatusInternalServerError)
	}
	menu, err := menuDAL.WithContext(ctx).Where(menuDAL.ID.Eq(id)).First()
	if err != nil {
		return nil, WrapDBErr(err)
	}
	return menu, nil
}

func (m *menuServiceImpl) Delete(ctx context.Context, id int32) error {
	menuDAL := dal.GetQueryByCtx(ctx).Menu
	deleteResult, err := menuDAL.WithContext(ctx).Where(menuDAL.ID.Eq(id)).Delete()
	if err != nil {
		return WrapDBErr(err)
	}
	if deleteResult.RowsAffected != 1 {
		return xerr.DB.New("delete menu failed id=%d", id).WithStatus(xerr.StatusInternalServerError)
	}
	return nil
}

func (m *menuServiceImpl) ConvertToDTO(ctx context.Context, menu *entity.Menu) *dto.Menu {
	return &dto.Menu{
		ID:       menu.ID,
		Name:     menu.Name,
		URL:      menu.URL,
		Priority: menu.Priority,
		Target:   menu.Target,
		Icon:     menu.Icon,
		ParentID: menu.ParentID,
		Team:     menu.Team,
	}
}

func (m *menuServiceImpl) ConvertToDTOs(ctx context.Context, menus []*entity.Menu) []*dto.Menu {
	result := make([]*dto.Menu, 0, len(menus))

	for _, link := range menus {
		result = append(result, m.ConvertToDTO(ctx, link))
	}
	return result
}

func (m *menuServiceImpl) ListTeams(ctx context.Context) ([]string, error) {
	teams := make([]string, 0)
	menuDAL := dal.GetQueryByCtx(ctx).Menu
	err := menuDAL.WithContext(ctx).Select(menuDAL.Team).Distinct(menuDAL.Team).Scan(&teams)
	if err != nil {
		return nil, WrapDBErr(err)
	}
	return teams, nil
}

func (m *menuServiceImpl) buildTree(ctx context.Context, menus []*entity.Menu) []*vo.Menu {
	menuTree := make([]*vo.Menu, 0)
	idToMenuMap := make(map[int32]*vo.Menu)

	for _, menu := range menus {
		menuDTO := m.ConvertToDTO(ctx, menu)
		menuVO := &vo.Menu{
			Menu:     *menuDTO,
			Children: make([]*vo.Menu, 0),
		}
		idToMenuMap[menuDTO.ID] = menuVO
	}
	for _, menu := range menus {
		menuVO := idToMenuMap[menu.ID]
		parentID := menuVO.ParentID
		if parentID == 0 {
			menuTree = append(menuTree, menuVO)
		} else if parent, ok := idToMenuMap[parentID]; ok {
			parent.Children = append(parent.Children, menuVO)
		} else {
			menuTree = append(menuTree, menuVO)
		}
	}
	return menuTree
}

func (m *menuServiceImpl) GetMenuCount(ctx context.Context) (int64, error) {
	menuDAL := dal.GetQueryByCtx(ctx).Menu
	count, err := menuDAL.WithContext(ctx).Count()
	return count, WrapDBErr(err)
}
