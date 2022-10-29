package impl

import (
	"context"
	"strings"

	"gorm.io/gen/field"

	"github.com/go-sonic/sonic/consts"
	"github.com/go-sonic/sonic/dal"
	"github.com/go-sonic/sonic/model/dto"
	"github.com/go-sonic/sonic/model/entity"
	"github.com/go-sonic/sonic/model/param"
	"github.com/go-sonic/sonic/model/projection"
	"github.com/go-sonic/sonic/model/property"
	"github.com/go-sonic/sonic/model/vo"
	"github.com/go-sonic/sonic/service"
	"github.com/go-sonic/sonic/util"
	"github.com/go-sonic/sonic/util/xerr"
)

type categoryServiceImpl struct {
	OptionService service.OptionService
}

func NewCategoryService(optionService service.OptionService) service.CategoryService {
	return &categoryServiceImpl{
		OptionService: optionService,
	}
}

func (c categoryServiceImpl) GetByID(ctx context.Context, id int32) (*entity.Category, error) {
	categoryDAL := dal.Use(dal.GetDBByCtx(ctx)).Category
	category, err := categoryDAL.WithContext(ctx).Where(categoryDAL.ID.Eq(id)).Take()
	return category, WrapDBErr(err)
}

func (c categoryServiceImpl) GetBySlug(ctx context.Context, slug string) (*entity.Category, error) {
	categoryDAL := dal.Use(dal.GetDBByCtx(ctx)).Category
	category, err := categoryDAL.WithContext(ctx).Where(categoryDAL.Slug.Eq(slug)).Take()
	return category, WrapDBErr(err)
}

func (c categoryServiceImpl) ListCategoryWithPostCountDTO(ctx context.Context, sort *param.Sort) ([]*dto.CategoryWithPostCount, error) {
	categoryDAL := dal.Use(dal.GetDBByCtx(ctx)).Category
	categoryDO := categoryDAL.WithContext(ctx)

	err := BuildSort(sort, &categoryDAL, &categoryDO)
	if err != nil {
		return nil, err
	}

	categories, err := categoryDO.Find()
	if err != nil {
		return nil, WrapDBErr(err)
	}
	postCategoryDAL := dal.Use(dal.GetDBByCtx(ctx)).PostCategory
	postCounts := make([]*projection.CategoryPostCountProjection, 0)
	err = postCategoryDAL.WithContext(ctx).Select(postCategoryDAL.CategoryID, postCategoryDAL.PostID.Count().As("post_count")).Group(postCategoryDAL.CategoryID).Scan(&postCounts)
	if err != nil {
		return nil, WrapDBErr(err)
	}
	postCountMap := make(map[int32]int32)
	for _, postCount := range postCounts {
		postCountMap[postCount.CategoryID] = postCount.PostCount
	}
	results := make([]*dto.CategoryWithPostCount, 0)
	for _, category := range categories {
		categoryDTO, err := c.ConvertToCategoryDTO(ctx, category)
		if err != nil {
			return nil, err
		}
		results = append(results, &dto.CategoryWithPostCount{
			CategoryDTO: categoryDTO,
			PostCount:   postCountMap[category.ID],
		})
	}
	return results, nil
}

func (c categoryServiceImpl) ListAll(ctx context.Context, sort *param.Sort) ([]*entity.Category, error) {
	categoryDAL := dal.Use(dal.GetDBByCtx(ctx)).Category
	categoryDO := categoryDAL.WithContext(ctx)
	err := BuildSort(sort, &categoryDAL, &categoryDO)
	if err != nil {
		return nil, err
	}
	categories, err := categoryDO.Find()
	return categories, err
}

func (c categoryServiceImpl) ListAsTree(ctx context.Context, sort *param.Sort, fillPassword bool) ([]*vo.CategoryVO, error) {
	allCategory, err := c.ListAll(ctx, sort)
	if err != nil {
		return nil, err
	}
	return c.buildTree(ctx, allCategory, fillPassword)
}

func (c categoryServiceImpl) buildTree(ctx context.Context, categories []*entity.Category, fillPassword bool) ([]*vo.CategoryVO, error) {
	categoryTree := make([]*vo.CategoryVO, 0)
	parentIDToCategoryMap := make(map[int32][]*vo.CategoryVO)
	categoryDTOs, err := c.ConvertToCategoryDTOs(ctx, categories)
	if err != nil {
		return nil, err
	}
	categoryVOs := make([]*vo.CategoryVO, len(categories))
	for i, categoryDTO := range categoryDTOs {
		categoryVO := &vo.CategoryVO{
			CategoryDTO: *categoryDTO,
		}
		if !fillPassword {
			categoryVO.Password = ""
		}
		if categoryDTO.ParentID != 0 {
			parentIDToCategoryMap[categoryDTO.ParentID] = append(parentIDToCategoryMap[categoryDTO.ParentID], categoryVO)
		}
		categoryVOs[i] = categoryVO
	}
	for _, categoryVO := range categoryVOs {
		children, ok := parentIDToCategoryMap[categoryVO.ID]
		if ok {
			categoryVO.Children = children
		}
		if categoryVO.ParentID == 0 {
			categoryTree = append(categoryTree, categoryVO)
		}
	}
	return categoryTree, nil
}

func (c categoryServiceImpl) ConvertToCategoryDTO(ctx context.Context, e *entity.Category) (*dto.CategoryDTO, error) {
	categoryDTO := &dto.CategoryDTO{}
	categoryDTO.ID = e.ID
	categoryDTO.Thumbnail = e.Thumbnail
	categoryDTO.ParentID = e.ParentID
	categoryDTO.Password = e.Password
	categoryDTO.Name = e.Name
	categoryDTO.CreateTime = e.CreateTime.UnixMilli()
	categoryDTO.Description = e.Description
	categoryDTO.Slug = e.Slug
	categoryDTO.Priority = e.Priority
	isEnabled, err := c.OptionService.IsEnabledAbsolutePath(ctx)
	if err != nil {
		return nil, err
	}
	fullPath := strings.Builder{}
	if isEnabled {
		blogBaseUrl, err := c.OptionService.GetBlogBaseURL(ctx)
		if err != nil {
			return nil, err
		}
		fullPath.WriteString(blogBaseUrl)
	}
	fullPath.WriteString("/")
	categoryPrefix, err := c.OptionService.GetOrByDefaultWithErr(ctx, property.CategoriesPrefix, "categories")
	if err != nil {
		return nil, err
	}
	fullPath.WriteString(categoryPrefix.(string))
	fullPath.WriteString("/")
	fullPath.WriteString(e.Slug)
	pathSuffix, err := c.OptionService.GetPathSuffix(ctx)
	if err != nil {
		return nil, err
	}
	fullPath.WriteString(pathSuffix)
	categoryDTO.FullPath = fullPath.String()
	return categoryDTO, nil
}

func (c categoryServiceImpl) ConvertToCategoryDTOs(ctx context.Context, categories []*entity.Category) ([]*dto.CategoryDTO, error) {
	result := make([]*dto.CategoryDTO, len(categories))
	isEnabled, err := c.OptionService.IsEnabledAbsolutePath(ctx)
	if err != nil {
		return nil, err
	}
	blogBaseUrl, err := c.OptionService.GetBlogBaseURL(ctx)
	if err != nil {
		return nil, err
	}
	categoryPrefix, err := c.OptionService.GetOrByDefaultWithErr(ctx, property.CategoriesPrefix, "categories")
	if err != nil {
		return nil, err
	}
	pathSuffix, err := c.OptionService.GetPathSuffix(ctx)
	if err != nil {
		return nil, err
	}
	for i, category := range categories {
		categoryDTO := &dto.CategoryDTO{}
		categoryDTO.ID = category.ID
		categoryDTO.Thumbnail = category.Thumbnail
		categoryDTO.ParentID = category.ParentID
		categoryDTO.Password = category.Password
		categoryDTO.Name = category.Name
		categoryDTO.CreateTime = category.CreateTime.UnixMilli()
		categoryDTO.Description = category.Description
		categoryDTO.Slug = category.Slug
		categoryDTO.Priority = category.Priority

		fullPath := strings.Builder{}
		if isEnabled {
			fullPath.WriteString(blogBaseUrl)
		}
		fullPath.WriteString("/")
		fullPath.WriteString(categoryPrefix.(string))
		fullPath.WriteString("/")
		fullPath.WriteString(category.Slug)
		fullPath.WriteString(pathSuffix)
		categoryDTO.FullPath = fullPath.String()
		result[i] = categoryDTO
	}
	return result, nil
}

func (c categoryServiceImpl) Create(ctx context.Context, categoryParam *param.Category) (*entity.Category, error) {
	categoryDAL := dal.Use(dal.GetDBByCtx(ctx)).Category

	count, err := categoryDAL.WithContext(ctx).Where(categoryDAL.Name.Eq(categoryParam.Name)).Count()
	if err != nil {
		return nil, WrapDBErr(err)
	}
	if count > 0 {
		return nil, xerr.BadParam.New("").WithMsg("Category has exist already").WithStatus(xerr.StatusBadRequest)
	}
	if categoryParam.Slug != "" {
		slugCount, err := categoryDAL.WithContext(ctx).Where(categoryDAL.Slug.Eq(categoryParam.Slug)).Count()
		if err != nil {
			return nil, WrapDBErr(err)
		}
		if slugCount > 0 {
			return nil, xerr.BadParam.New("").WithStatus(xerr.StatusBadRequest).WithMsg("category slug has exist already")
		}
	}
	var parentCategory *entity.Category
	if categoryParam.ParentID > 0 {
		parentCategory, err = categoryDAL.WithContext(ctx).Where(categoryDAL.ID.Eq(categoryParam.ParentID)).Take()
		if err != nil {
			return nil, WrapDBErr(err)
		}
		if parentCategory == nil {
			return nil, xerr.BadParam.New("parentID=%d", categoryParam.ParentID).WithStatus(xerr.StatusBadRequest).WithMsg("Parent Category was not found")
		}
	}
	if categoryParam.Password != "" {
		categoryParam.Password = strings.TrimSpace(categoryParam.Password)
	}
	if categoryParam.Slug == "" {
		categoryParam.Slug = util.Slug(categoryParam.Name)
	} else {
		categoryParam.Slug = util.Slug(categoryParam.Slug)
	}

	category := &entity.Category{
		Name:        categoryParam.Name,
		Slug:        categoryParam.Slug,
		Description: categoryParam.Description,
		Password:    categoryParam.Password,
		Thumbnail:   categoryParam.Thumbnail,
		ParentID:    categoryParam.ParentID,
		Priority:    categoryParam.Priority,
		Type:        util.IfElse(categoryParam.Password == "", consts.CategoryTypeNormal, consts.CategoryTypeIntimate).(consts.CategoryType),
	}
	if parentCategory != nil && parentCategory.Type == consts.CategoryTypeIntimate {
		category.Type = consts.CategoryTypeIntimate
	}
	err = categoryDAL.WithContext(ctx).Create(category)
	if err != nil {
		return nil, WrapDBErr(err)
	}
	return category, nil
}

func (c *categoryServiceImpl) Update(ctx context.Context, categoryParam *param.Category) (*entity.Category, error) {

	executor := newCategoryUpdateExecutor(ctx)
	if err := executor.Update(ctx, categoryParam); err != nil {
		return nil, err
	}

	categoryDAL := dal.Use(dal.GetDBByCtx(ctx)).Category
	category, err := categoryDAL.WithContext(ctx).Where(categoryDAL.ID.Eq(categoryParam.ID)).First()
	return category, WrapDBErr(err)
}

func (c categoryServiceImpl) UpdateBatch(ctx context.Context, categoryParams []*param.Category) ([]*entity.Category, error) {

	executor := newCategoryUpdateExecutor(ctx)
	if err := executor.UpdateBatch(ctx, categoryParams); err != nil {
		return nil, err
	}

	categoryDAL := dal.Use(dal.GetDBByCtx(ctx)).Category
	categoryIDs := make([]int32, 0)
	for _, category := range categoryParams {
		categoryIDs = append(categoryIDs, category.ID)
	}
	categories, err := categoryDAL.WithContext(ctx).Where(categoryDAL.ID.In(categoryIDs...)).Find()
	if err != nil {
		return nil, WrapDBErr(err)
	}
	return categories, nil
}

func (c categoryServiceImpl) Delete(ctx context.Context, categoryID int32) (err error) {
	return newCategoryUpdateExecutor(ctx).Delete(ctx, categoryID)
}

func (c categoryServiceImpl) ListByIDs(ctx context.Context, categoryIDs []int32) ([]*entity.Category, error) {
	categoryDAL := dal.Use(dal.GetDBByCtx(ctx)).Category
	categories, err := categoryDAL.WithContext(ctx).Where(categoryDAL.ID.In(categoryIDs...)).Find()
	if err != nil {
		return nil, WrapDBErr(err)
	}
	return categories, nil
}

func (c categoryServiceImpl) IsCategoriesEncrypt(ctx context.Context, categoryIDs ...int32) (bool, error) {
	if len(categoryIDs) == 0 {
		return false, nil
	}
	categoryDAL := dal.Use(dal.GetDBByCtx(ctx)).Category
	count, err := categoryDAL.WithContext(ctx).Where(categoryDAL.ID.In(categoryIDs...), categoryDAL.Type.Eq(consts.CategoryTypeIntimate)).Count()
	return count > 0, WrapDBErr(err)
}

func (c categoryServiceImpl) Count(ctx context.Context) (int64, error) {
	categoryDAL := dal.Use(dal.GetDBByCtx(ctx)).Category
	count, err := categoryDAL.WithContext(ctx).Count()
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (c *categoryServiceImpl) GetChildCategory(ctx context.Context, parentCategoryID int32) ([]*entity.Category, error) {
	allCategory, err := c.ListAll(ctx, nil)
	if err != nil {
		return nil, err
	}
	allCategoryMap := make(map[int32]*entity.Category)
	parentIDToChild := make(map[int32][]*entity.Category)
	for _, category := range allCategory {
		parentIDToChild[category.ParentID] = append(parentIDToChild[category.ParentID], category)
		allCategoryMap[category.ID] = category
	}

	q := util.NewQueueCap[int32](len(allCategory))

	for _, category := range parentIDToChild[parentCategoryID] {
		q.Push(category.ID)
	}
	childs := make([]*entity.Category, 0)
	for !q.IsEmpty() {
		categoryID := q.Next()
		childs = append(childs, allCategory[categoryID])
		for _, category := range parentIDToChild[categoryID] {
			q.Push(category.ID)
		}
	}
	return childs, nil

}

type categoryUpdateExecutor struct {
	AllCategory    map[int32]*entity.Category
	CategoryPosts  map[int32][]*entity.Post
	PostMap        map[int32]*entity.Post
	PostToCategory map[int32][]*entity.Category
}

func newCategoryUpdateExecutor(ctx context.Context) *categoryUpdateExecutor {
	return &categoryUpdateExecutor{}
}

func (c *categoryUpdateExecutor) Update(ctx context.Context, categoryParam *param.Category) error {
	if err := c.prepare(ctx, []*param.Category{categoryParam}); err != nil {
		return err
	}

	parentCategory, ok := c.AllCategory[categoryParam.ParentID]
	oldCategory := c.AllCategory[categoryParam.ID]
	if categoryParam.ParentID > 0 && (!ok || parentCategory == nil) {
		return xerr.BadParam.New("parentID=%d", categoryParam.ParentID).WithStatus(xerr.StatusBadRequest).WithMsg("Parent Category was not found")
	}

	category := c.convertParam(categoryParam)
	if parentCategory != nil && parentCategory.Type == consts.CategoryTypeIntimate {
		category.Type = consts.CategoryTypeIntimate
	}

	err := dal.Transaction(ctx, func(txCtx context.Context) error {
		categoryDAL := dal.Use(dal.GetDBByCtx(txCtx)).Category

		resultInfo, err := categoryDAL.WithContext(txCtx).Where(categoryDAL.ID.Eq(categoryParam.ID)).Select(field.Star).Omit(categoryDAL.CreateTime).Updates(category)
		if err != nil {
			return WrapDBErr(err)
		}
		if resultInfo.RowsAffected != 1 {
			return xerr.DB.New("").WithMsg("update failed")
		}

		if oldCategory.Type != category.Type {
			c.AllCategory[category.ID].Type = category.Type
			if err := c.refreshChildsType(txCtx, category.ID, category.Type); err != nil {
				return err
			}
			if err := c.refreshPostStatus(txCtx); err != nil {
				return nil
			}
		}

		return nil
	})
	return err
}

func (c *categoryUpdateExecutor) UpdateBatch(ctx context.Context, categoryParams []*param.Category) error {

	categories := make([]*entity.Category, 0)
	for _, categoryParam := range categoryParams {
		categories = append(categories, c.convertParam(categoryParam))
	}

	err := dal.Transaction(ctx, func(txCtx context.Context) error {
		categoryDAL := dal.Use(dal.GetDBByCtx(txCtx)).Category
		err := categoryDAL.WithContext(txCtx).Omit(categoryDAL.CreateTime).Save(categories...)
		if err != nil {
			return WrapDBErr(err)
		}
		if err := c.prepare(txCtx, categoryParams); err != nil {
			return err
		}
		if err := c.refreshAllType(txCtx); err != nil {
			return err
		}
		if err := c.refreshPostStatus(txCtx); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *categoryUpdateExecutor) Delete(ctx context.Context, categoryID int32) error {
	if err := c.prepare(ctx, []*param.Category{{ID: categoryID}}); err != nil {
		return err
	}
	curCategory, ok := c.AllCategory[categoryID]
	if !ok || curCategory == nil {
		return xerr.BadParam.New("category=%d", categoryID).WithStatus(xerr.StatusBadRequest).WithMsg("Category was not found")
	}

	parent, ok := c.AllCategory[curCategory.ParentID]

	err := dal.Transaction(ctx, func(txCtx context.Context) error {

		if ok && parent != nil {
			if parent.Type == consts.CategoryTypeNormal && curCategory.Type == consts.CategoryTypeIntimate {
				if err := c.refreshChildsType(txCtx, categoryID, consts.CategoryTypeNormal); err != nil {
					return err
				}
			}
		} else if curCategory.Type == consts.CategoryTypeIntimate {
			if err := c.refreshChildsType(txCtx, categoryID, consts.CategoryTypeNormal); err != nil {
				return err
			}
		}

		if err := c.removePostCategory(txCtx, categoryID); err != nil {
			return err
		}

		if err := c.removeCategory(txCtx, categoryID); err != nil {
			return err
		}
		if err := c.refreshPostStatus(txCtx); err != nil {
			return err
		}
		return nil
	})
	return err
}

func (c *categoryUpdateExecutor) prepare(ctx context.Context, categoryParams []*param.Category) error {
	categoryDAL := dal.Use(dal.GetDBByCtx(ctx)).Category
	categoryDO := categoryDAL.WithContext(ctx)

	categories, err := categoryDO.Find()
	if err != nil {
		return err
	}
	categoryMap := make(map[int32]*entity.Category)

	for _, category := range categories {
		categoryMap[category.ID] = category
	}
	c.AllCategory = categoryMap

	categoryIDMap := make(map[int32]struct{})
	for _, category := range categoryParams {
		categoryIDMap[category.ID] = struct{}{}
		childs := c.getChildCategory(category.ID)
		for _, child := range childs {
			categoryIDMap[child.ID] = struct{}{}
		}
	}
	categoryIDs := util.MapKeyToArray(categoryIDMap)

	postDAL := dal.Use(dal.GetDBByCtx(ctx)).Post
	postCategoryDAL := dal.Use(dal.GetDBByCtx(ctx)).PostCategory

	posts, err := postDAL.WithContext(ctx).Where(
		postDAL.Password.Null()).Or(postDAL.Password.Zero()).Where(
		postDAL.WithContext(ctx).Columns(postDAL.ID).In(
			postCategoryDAL.WithContext(ctx).Select(postCategoryDAL.PostID).Where(postCategoryDAL.CategoryID.In(categoryIDs...)),
		),
		postDAL.Status.Neq(consts.PostStatusRecycle),
	).Select(postDAL.Status, postDAL.Password, postDAL.ID).Find()
	if err != nil {
		return WrapDBErr(err)
	}

	postMap := make(map[int32]*entity.Post)
	for _, post := range posts {
		postMap[post.ID] = post
	}
	c.PostMap = postMap

	postIDs := make([]int32, 0)
	for _, post := range posts {
		postIDs = append(postIDs, post.ID)
	}
	postCategorys, err := postCategoryDAL.WithContext(ctx).Where(postCategoryDAL.PostID.In(postIDs...)).Find()
	if err != nil {
		return err
	}
	postToCategory := make(map[int32][]*entity.Category, 0)
	for _, postCategory := range postCategorys {
		category := categoryMap[postCategory.CategoryID]
		postToCategory[postCategory.PostID] = append(postToCategory[postCategory.PostID], category)
	}
	c.PostToCategory = postToCategory
	return nil
}

func (c *categoryUpdateExecutor) getChildCategory(parentCategoryID int32) []*entity.Category {

	parentIDToChild := make(map[int32][]*entity.Category)
	for _, category := range c.AllCategory {
		parentIDToChild[category.ParentID] = append(parentIDToChild[category.ParentID], category)
	}

	q := util.NewQueueCap[int32](len(c.AllCategory))

	for _, category := range parentIDToChild[parentCategoryID] {
		q.Push(category.ID)
	}
	childs := make([]*entity.Category, 0)
	for !q.IsEmpty() {
		categoryID := q.Next()
		childs = append(childs, c.AllCategory[categoryID])
		for _, category := range parentIDToChild[categoryID] {
			q.Push(category.ID)
		}
	}
	return childs

}

func (c *categoryUpdateExecutor) convertParam(categoryParam *param.Category) *entity.Category {
	if categoryParam.Slug != "" {
		categoryParam.Slug = util.Slug(categoryParam.Slug)
	} else {
		categoryParam.Slug = util.Slug(categoryParam.Name)
	}

	return &entity.Category{
		ID:          categoryParam.ID,
		Name:        categoryParam.Name,
		Description: categoryParam.Description,
		Thumbnail:   categoryParam.Thumbnail,
		ParentID:    categoryParam.ParentID,
		Password:    strings.TrimSpace(categoryParam.Password),
		Slug:        categoryParam.Slug,
		Priority:    categoryParam.Priority,
		Type:        util.IfElse(categoryParam.Password == "", consts.CategoryTypeNormal, consts.CategoryTypeIntimate).(consts.CategoryType),
	}
}

func (c *categoryUpdateExecutor) refreshChildsType(ctx context.Context, parentID int32, categoryType consts.CategoryType) error {
	childs := c.getChildCategory(parentID)
	for _, category := range childs {
		category.Type = consts.CategoryTypeNormal
	}

	needEncryptMap := make(map[int32]struct{})

	for _, child := range childs {
		if categoryType == consts.CategoryTypeIntimate {
			child.Type = consts.CategoryTypeIntimate
			needEncryptMap[child.ID] = struct{}{}
			continue
		}
		if child.Password != "" {
			child.Type = consts.CategoryTypeIntimate
			needEncryptMap[child.ID] = struct{}{}
			childs := c.getChildCategory(child.ID)
			for _, child := range childs {
				child.Type = consts.CategoryTypeIntimate
				needEncryptMap[child.ID] = struct{}{}
			}
		}
	}
	needDecrypt := make([]int32, 0)
	for _, child := range childs {
		if _, ok := needEncryptMap[child.ID]; !ok {
			needDecrypt = append(needDecrypt, child.ID)
			child.Type = consts.CategoryTypeNormal
		}
	}
	needEncrypt := util.MapKeyToArray(needEncryptMap)
	if len(needEncrypt) > 0 {
		categoryDAL := dal.Use(dal.GetDBByCtx(ctx)).Category
		_, err := categoryDAL.WithContext(ctx).Where(categoryDAL.ID.In(needEncrypt...)).UpdateColumnSimple(categoryDAL.Type.Value(consts.CategoryTypeIntimate))
		if err != nil {
			return WrapDBErr(err)
		}
	}

	if len(needDecrypt) > 0 {
		categoryDAL := dal.Use(dal.GetDBByCtx(ctx)).Category
		_, err := categoryDAL.WithContext(ctx).Where(categoryDAL.ID.In(needDecrypt...)).UpdateColumnSimple(categoryDAL.Type.Value(consts.CategoryTypeNormal))
		if err != nil {
			return WrapDBErr(err)
		}
	}
	return nil
}

func (c *categoryUpdateExecutor) refreshAllType(ctx context.Context) error {
	for _, category := range c.AllCategory {
		category.Type = consts.CategoryTypeNormal
	}

	needEncryptMap := make(map[int32]struct{})
	for _, category := range c.AllCategory {
		if category.Password != "" {
			category.Type = consts.CategoryTypeIntimate
			needEncryptMap[category.ID] = struct{}{}
			childs := c.getChildCategory(category.ID)
			for _, child := range childs {
				child.Type = consts.CategoryTypeIntimate
				needEncryptMap[child.ID] = struct{}{}
			}
		}
	}
	needEncrypt := util.MapKeyToArray(needEncryptMap)

	needDecrypt := make([]int32, 0)
	for id, category := range c.AllCategory {
		if _, ok := needEncryptMap[id]; !ok {
			needDecrypt = append(needDecrypt, id)
			category.Type = consts.CategoryTypeNormal
		}
	}

	if len(needEncrypt) > 0 {
		categoryDAL := dal.Use(dal.GetDBByCtx(ctx)).Category
		_, err := categoryDAL.WithContext(ctx).Where(categoryDAL.ID.In(needEncrypt...)).UpdateColumnSimple(categoryDAL.Type.Value(consts.CategoryTypeIntimate))
		if err != nil {
			return WrapDBErr(err)
		}
	}

	if len(needDecrypt) > 0 {
		categoryDAL := dal.Use(dal.GetDBByCtx(ctx)).Category
		_, err := categoryDAL.WithContext(ctx).Where(categoryDAL.ID.In(needDecrypt...)).UpdateColumnSimple(categoryDAL.Type.Value(consts.CategoryTypeNormal))
		if err != nil {
			return WrapDBErr(err)
		}
	}
	return nil
}

func (c *categoryUpdateExecutor) refreshPostStatus(ctx context.Context) error {

	needEncryptPostID := make([]int32, 0)
	needDecryptPostID := make([]int32, 0)
	for id, post := range c.PostMap {
		var status consts.PostStatus
		for _, category := range c.PostToCategory[id] {
			if category.Type == consts.CategoryTypeIntimate {
				status = consts.PostStatusIntimate
				break
			}
		}
		if status == consts.PostStatusIntimate {
			needEncryptPostID = append(needEncryptPostID, id)
		} else {
			if post.Status == consts.PostStatusIntimate && post.Password == "" {
				needDecryptPostID = append(needDecryptPostID, id)
			}
		}
	}
	if len(needEncryptPostID) > 0 {
		postDAL := dal.Use(dal.GetDBByCtx(ctx)).Post
		_, err := postDAL.WithContext(ctx).Where(postDAL.ID.In(needEncryptPostID...), postDAL.Status.Neq(consts.PostStatusDraft)).UpdateColumnSimple(postDAL.Status.Value(consts.PostStatusIntimate))
		if err != nil {
			return WrapDBErr(err)
		}
	}
	if len(needDecryptPostID) > 0 {
		postDAL := dal.Use(dal.GetDBByCtx(ctx)).Post
		_, err := postDAL.WithContext(ctx).Where(postDAL.ID.In(needDecryptPostID...), postDAL.Status.Neq(consts.PostStatusDraft)).UpdateColumnSimple(postDAL.Status.Value(consts.PostStatusPublished))
		if err != nil {
			return WrapDBErr(err)
		}
	}
	return nil
}

func (c *categoryUpdateExecutor) removePostCategory(ctx context.Context, categoryID int32) error {
	parent, ok := c.AllCategory[c.AllCategory[categoryID].ParentID]
	postCategory := make([]*entity.PostCategory, 0)
	for postID, categories := range c.PostToCategory {
		if len(categories) == 1 && categories[0].ID == categoryID {
			if ok && parent != nil {
				postCategory = append(postCategory, &entity.PostCategory{
					PostID:     postID,
					CategoryID: parent.ID,
				})
				c.PostToCategory[postID] = []*entity.Category{parent}
			} else {
				c.PostToCategory[postID] = nil
			}
		} else {
			t := make([]*entity.Category, 0)
			for _, category := range categories {
				if category.ID == categoryID {
					continue
				}
				t = append(t, category)
			}
			c.PostToCategory[postID] = t
		}
	}
	postCategoryDAL := dal.Use(dal.GetDBByCtx(ctx)).PostCategory
	err := postCategoryDAL.WithContext(ctx).Create(postCategory...)

	if err != nil {
		return WrapDBErr(err)
	}
	_, err = postCategoryDAL.WithContext(ctx).Where(postCategoryDAL.CategoryID.Eq(categoryID)).Delete()

	return WrapDBErr(err)

}

func (c *categoryUpdateExecutor) removeCategory(ctx context.Context, categoryID int32) error {
	childs := c.getChildCategory(categoryID)
	toUpdate := make([]int32, 0)
	for _, child := range childs {
		if child.ParentID == categoryID {
			toUpdate = append(toUpdate, child.ID)
		}
	}
	parent, ok := c.AllCategory[c.AllCategory[categoryID].ParentID]
	parentID := 0
	if ok && parent != nil {
		parentID = int(parent.ID)
	}
	categoryDAL := dal.Use(dal.GetDBByCtx(ctx)).Category
	if len(toUpdate) > 0 {
		_, err := categoryDAL.WithContext(ctx).Where(categoryDAL.ID.In(toUpdate...)).UpdateColumnSimple(categoryDAL.ParentID.Value(int32(parentID)))
		if err != nil {
			return WrapDBErr(err)
		}
	}
	_, err := categoryDAL.WithContext(ctx).Where(categoryDAL.ID.Eq(categoryID)).Delete()
	return WrapDBErr(err)
}
