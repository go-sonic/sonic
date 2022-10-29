package impl

import (
	"context"
	"time"

	"github.com/go-sonic/sonic/cache"
	"github.com/go-sonic/sonic/model/entity"
	"github.com/go-sonic/sonic/service"
	"github.com/go-sonic/sonic/util/xerr"
)

type authenticateServiceImpl struct {
	CategoryService service.CategoryService
	Cache           cache.Cache
}

func NewAuthenticateService(categoryService service.CategoryService, cache cache.Cache) service.AuthenticateService {
	return &authenticateServiceImpl{
		CategoryService: categoryService,
		Cache:           cache,
	}
}

func (a *authenticateServiceImpl) PostAuthenticate(ctx context.Context, post *entity.Post, password string) (bool, error) {
	panic("implement me")
}

func (a *authenticateServiceImpl) CategoryAuthenticate(ctx context.Context, categoryID int32, password string) (bool, error) {
	categories, err := a.CategoryService.ListAll(ctx, nil)
	if err != nil {
		return false, nil
	}
	categoryIDMap := make(map[int32]*entity.Category)
	for _, category := range categories {
		categoryIDMap[category.ID] = category
	}
	return a.doCategoryAuthenticate(ctx, categoryIDMap, categoryID, password)
}

func (a *authenticateServiceImpl) doCategoryAuthenticate(ctx context.Context, categoryIDMap map[int32]*entity.Category, categoryID int32, password string) (bool, error) {
	category, ok := categoryIDMap[categoryID]
	if !ok || category == nil {
		return false, nil
	}
	permissionsMap, err := a.getAccessPermission(ctx)
	if err != nil {
		return false, err
	}
	if category.Password != "" {
		if _, ok := permissionsMap[cache.BuildCategoryPermissionKey(categoryID)]; ok {
			return true, nil
		}
		if category.Password == password {
			err := a.setAccessPermission(ctx, cache.BuildCategoryPermissionKey(categoryID))
			return false, xerr.WithErrMsgf(err, "set category permission cache failed")
		}
		return false, nil
	}
	if category.ParentID == 0 {
		return true, nil
	}
	return a.doCategoryAuthenticate(ctx, categoryIDMap, category.ParentID, password)
}

func (a *authenticateServiceImpl) getAccessPermission(ctx context.Context) (map[string]struct{}, error) {
	accessKey, err := cache.BuildAccessPermissionKey(ctx)
	if err != nil {
		return nil, err
	}
	permissions, ok := a.Cache.Get(accessKey)
	if !ok {
		return make(map[string]struct{}), nil
	}
	permissionsMap, ok := permissions.(map[string]struct{})
	if !ok {
		return nil, xerr.NoType.New("cache value is not map[string]struct{}")
	}
	return permissionsMap, nil
}

func (a *authenticateServiceImpl) setAccessPermission(ctx context.Context, permissionKey string) error {
	accessKey, err := cache.BuildAccessPermissionKey(ctx)
	if err != nil {
		return err
	}
	permissions, ok := a.Cache.Get(accessKey)
	if !ok {
		return nil
	}
	permissionsMap, ok := permissions.(map[string]struct{})
	if !ok {
		return xerr.NoType.New("cache value is not map[string]struct{}")
	}
	permissionsMap[permissionKey] = struct{}{}
	a.Cache.Set(accessKey, permissionsMap, time.Hour*24)
	return nil
}
