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

type linkServiceImpl struct{}

// ConvertToLinkTeamVO implements service.LinkService
func (l *linkServiceImpl) ConvertToLinkTeamVO(ctx context.Context, links []*entity.Link) []*vo.LinkTeamVO {
	m := make(map[string][]*dto.Link)
	for _, link := range links {
		m[link.Team] = append(m[link.Team], l.ConvertToDTO(ctx, link))
	}
	result := make([]*vo.LinkTeamVO, 0)
	for team, links := range m {
		result = append(result, &vo.LinkTeamVO{
			Team:  team,
			Links: links,
		})
	}
	return result
}

func NewLinkService() service.LinkService {
	return &linkServiceImpl{}
}

func (l *linkServiceImpl) ListTeams(ctx context.Context) ([]string, error) {
	teams := make([]string, 0)
	linkDAL := dal.Use(dal.GetDBByCtx(ctx)).Link
	err := linkDAL.WithContext(ctx).Select(linkDAL.Team).Distinct(linkDAL.Team).Scan(&teams)
	if err != nil {
		return nil, WrapDBErr(err)
	}
	return teams, nil
}

func (l *linkServiceImpl) Delete(ctx context.Context, id int32) error {
	linkDAL := dal.Use(dal.GetDBByCtx(ctx)).Link
	deleteResult, err := linkDAL.WithContext(ctx).Where(linkDAL.ID.Eq(id)).Delete()
	if err != nil {
		return WrapDBErr(err)
	}
	if deleteResult.RowsAffected != 1 {
		return xerr.DB.New("delete link failed id=%d", id).WithStatus(xerr.StatusInternalServerError)
	}
	return nil
}

func (l *linkServiceImpl) Create(ctx context.Context, linkParam *param.Link) (*entity.Link, error) {
	link := &entity.Link{
		Name:        linkParam.Name,
		Description: linkParam.Description,
		URL:         linkParam.URL,
		Logo:        linkParam.Logo,
		Priority:    linkParam.Priority,
		Team:        linkParam.Team,
	}
	linkDAL := dal.Use(dal.GetDBByCtx(ctx)).Link
	err := linkDAL.WithContext(ctx).Create(link)
	if err != nil {
		return nil, WrapDBErr(err)
	}
	return link, nil
}

func (l *linkServiceImpl) Update(ctx context.Context, id int32, linkParam *param.Link) (*entity.Link, error) {
	linkDAL := dal.Use(dal.GetDBByCtx(ctx)).Link
	updateResult, err := linkDAL.WithContext(ctx).Where(linkDAL.ID.Eq(id)).UpdateSimple(
		linkDAL.Name.Value(linkParam.Name),
		linkDAL.Description.Value(linkParam.Description),
		linkDAL.URL.Value(linkParam.URL),
		linkDAL.Logo.Value(linkParam.Logo),
		linkDAL.Priority.Value(linkParam.Priority),
		linkDAL.Team.Value(linkParam.Team),
	)
	if err != nil {
		return nil, WrapDBErr(err)
	}
	if err != nil {
		return nil, WrapDBErr(err)
	}
	if updateResult.RowsAffected != 1 {
		return nil, xerr.NoType.New("update link failed").WithMsg("update link failed").WithStatus(xerr.StatusInternalServerError)
	}
	link, err := linkDAL.WithContext(ctx).Where(linkDAL.ID.Eq(id)).First()
	if err != nil {
		return nil, WrapDBErr(err)
	}
	return link, nil
}

func (l *linkServiceImpl) List(ctx context.Context, sort *param.Sort) ([]*entity.Link, error) {
	linkDAL := dal.Use(dal.GetDBByCtx(ctx)).Link
	linkDO := linkDAL.WithContext(ctx)
	err := BuildSort(sort, &linkDAL, &linkDO)
	if err != nil {
		return nil, err
	}
	links, err := linkDO.Find()
	if err != nil {
		return nil, WrapDBErr(err)
	}
	return links, nil
}

func (l *linkServiceImpl) GetByID(ctx context.Context, id int32) (*entity.Link, error) {
	linkDAL := dal.Use(dal.GetDBByCtx(ctx)).Link
	link, err := linkDAL.WithContext(ctx).Where(linkDAL.ID.Eq(id)).First()
	if err != nil {
		return nil, err
	}
	return link, nil
}

func (l *linkServiceImpl) ConvertToDTO(ctx context.Context, link *entity.Link) *dto.Link {
	return &dto.Link{
		ID:          link.ID,
		Name:        link.Name,
		URL:         link.URL,
		Logo:        link.Logo,
		Description: link.Description,
		Team:        link.Team,
		Priority:    link.Priority,
	}
}

func (l *linkServiceImpl) ConvertToDTOs(ctx context.Context, links []*entity.Link) []*dto.Link {
	result := make([]*dto.Link, 0, len(links))

	for _, link := range links {
		result = append(result, l.ConvertToDTO(ctx, link))
	}
	return result
}

func (l *linkServiceImpl) Count(ctx context.Context) (int64, error) {
	linkDAL := dal.Use(dal.GetDBByCtx(ctx)).Link
	count, err := linkDAL.WithContext(ctx).Count()
	if err != nil {
		return 0, WrapDBErr(err)
	}
	return count, nil
}
