package impl

import (
	"context"
	"time"

	"github.com/go-sonic/sonic/dal"
	"github.com/go-sonic/sonic/model/dto"
	"github.com/go-sonic/sonic/model/entity"
	"github.com/go-sonic/sonic/model/param"
	"github.com/go-sonic/sonic/service"
	"github.com/go-sonic/sonic/util"
	"github.com/go-sonic/sonic/util/xerr"
)

type photoServiceImpl struct{}

func NewPhotoService() service.PhotoService {
	return &photoServiceImpl{}
}

func (p *photoServiceImpl) Page(ctx context.Context, page param.Page, sort *param.Sort) ([]*entity.Photo, int64, error) {
	if page.PageNum < 0 || page.PageSize <= 0 || page.PageSize > 100 {
		return nil, 0, xerr.BadParam.New("").WithStatus(xerr.StatusBadRequest).WithMsg("Paging parameter error")
	}
	photoDAL := dal.GetQueryByCtx(ctx).Photo
	photoDO := photoDAL.WithContext(ctx)
	err := BuildSort(sort, &photoDAL, &photoDO)
	if err != nil {
		return nil, 0, err
	}
	photos, totalCount, err := photoDO.FindByPage(page.PageSize*page.PageNum, page.PageSize)
	if err != nil {
		return nil, 0, WrapDBErr(err)
	}
	return photos, totalCount, nil
}

func (p *photoServiceImpl) ListTeams(ctx context.Context) ([]string, error) {
	teams := make([]string, 0)
	photoDAL := dal.GetQueryByCtx(ctx).Photo
	err := photoDAL.WithContext(ctx).Select(photoDAL.Team).Distinct(photoDAL.Team).Scan(&teams)
	if err != nil {
		return nil, WrapDBErr(err)
	}
	return teams, nil
}

func (p *photoServiceImpl) Delete(ctx context.Context, id int32) error {
	photoDAL := dal.GetQueryByCtx(ctx).Photo
	deleteResult, err := photoDAL.WithContext(ctx).Where(photoDAL.ID.Eq(id)).Delete()
	if err != nil {
		return WrapDBErr(err)
	}
	if deleteResult.RowsAffected != 1 {
		return xerr.DB.New("delete photo failed id=%d", id).WithStatus(xerr.StatusInternalServerError)
	}
	return nil
}

func (p *photoServiceImpl) Create(ctx context.Context, photoParam *param.Photo) (*entity.Photo, error) {
	photo := &entity.Photo{
		Name:        photoParam.Name,
		Description: photoParam.Description,
		URL:         photoParam.URL,
		Thumbnail:   photoParam.Thumbnail,
		Location:    photoParam.Location,
		Team:        photoParam.Team,
	}
	if photoParam.TakeTime != nil && *photoParam.TakeTime != 0 {
		photo.CreateTime = time.Unix(*photoParam.TakeTime, 0)
	}
	photoDAL := dal.GetQueryByCtx(ctx).Photo
	err := photoDAL.WithContext(ctx).Create(photo)
	if err != nil {
		return nil, WrapDBErr(err)
	}
	return photo, nil
}

func (*photoServiceImpl) CreateBatch(ctx context.Context, photosParam []*param.Photo) ([]*entity.Photo, error) {
	photos := make([]*entity.Photo, len(photosParam))
	for i, photoParam := range photosParam {
		photo := &entity.Photo{
			Name:        photoParam.Name,
			Description: photoParam.Description,
			URL:         photoParam.URL,
			Thumbnail:   photoParam.Thumbnail,
			Location:    photoParam.Location,
			Team:        photoParam.Team,
		}
		photos[i] = photo
	}
	photoDAL := dal.GetQueryByCtx(ctx).Photo
	err := photoDAL.WithContext(ctx).CreateInBatches(photos, 100)
	if err != nil {
		return nil, WrapDBErr(err)
	}
	return photos, err
}

func (p *photoServiceImpl) Update(ctx context.Context, id int32, photoParam *param.Photo) (*entity.Photo, error) {
	photoDAL := dal.GetQueryByCtx(ctx).Photo
	var takeTime *time.Time
	if photoParam.TakeTime != nil && *photoParam.TakeTime != 0 {
		takeTime = util.TimePtr(time.Unix(*photoParam.TakeTime, 0))
	}
	photo := &entity.Photo{
		ID:          id,
		Name:        photoParam.Name,
		URL:         photoParam.URL,
		Thumbnail:   photoParam.Thumbnail,
		Description: photoParam.Description,
		Location:    photoParam.Location,
		TakeTime:    takeTime,
		Team:        photoParam.Team,
	}
	updateResult, err := photoDAL.WithContext(ctx).Where(photoDAL.ID.Eq(id)).Updates(photo)
	if err != nil {
		return nil, WrapDBErr(err)
	}
	if err != nil {
		return nil, WrapDBErr(err)
	}
	if updateResult.RowsAffected != 1 {
		return nil, xerr.NoType.New("update photo failed").WithMsg("update photo failed").WithStatus(xerr.StatusInternalServerError)
	}
	photo, err = photoDAL.WithContext(ctx).Where(photoDAL.ID.Eq(id)).First()
	if err != nil {
		return nil, WrapDBErr(err)
	}
	return photo, nil
}

func (p *photoServiceImpl) List(ctx context.Context, sort *param.Sort) ([]*entity.Photo, error) {
	photoDAL := dal.GetQueryByCtx(ctx).Photo
	photoDO := photoDAL.WithContext(ctx)
	err := BuildSort(sort, &photoDAL, &photoDO)
	if err != nil {
		return nil, err
	}
	photos, err := photoDO.Find()
	if err != nil {
		return nil, WrapDBErr(err)
	}
	return photos, nil
}

func (p *photoServiceImpl) GetByID(ctx context.Context, id int32) (*entity.Photo, error) {
	photoDAL := dal.GetQueryByCtx(ctx).Photo
	photo, err := photoDAL.WithContext(ctx).Where(photoDAL.ID.Eq(id)).First()
	if err != nil {
		return nil, err
	}
	return photo, nil
}

func (p *photoServiceImpl) ConvertToDTO(ctx context.Context, photo *entity.Photo) *dto.Photo {
	return &dto.Photo{
		ID:          photo.ID,
		Name:        photo.Name,
		URL:         photo.URL,
		Thumbnail:   photo.Thumbnail,
		Description: photo.Description,
		Team:        photo.Team,
		Location:    photo.Location,
	}
}

func (p *photoServiceImpl) ConvertToDTOs(ctx context.Context, photos []*entity.Photo) []*dto.Photo {
	result := make([]*dto.Photo, 0, len(photos))

	for _, photo := range photos {
		result = append(result, p.ConvertToDTO(ctx, photo))
	}
	return result
}

func (p *photoServiceImpl) IncreaseLike(ctx context.Context, id int32) error {
	photoDAL := dal.GetQueryByCtx(ctx).Photo
	updateResult, err := photoDAL.WithContext(ctx).Where(photoDAL.ID.Eq(id)).UpdateSimple(photoDAL.Likes.Add(1))
	if err != nil {
		return WrapDBErr(err)
	}
	if updateResult.RowsAffected != 1 {
		return xerr.NoType.New("increase photo like failed id=%v", id).WithMsg("increase like failed").WithStatus(xerr.StatusInternalServerError)
	}
	return nil
}

func (p *photoServiceImpl) ListByTeam(ctx context.Context, team string, sort *param.Sort) ([]*entity.Photo, error) {
	photoDAL := dal.GetQueryByCtx(ctx).Photo
	photoDO := photoDAL.WithContext(ctx)
	err := BuildSort(sort, &photoDAL, &photoDO)
	if err != nil {
		return nil, err
	}
	photos, err := photoDO.Where(photoDAL.Team.Eq(team)).Find()
	return photos, WrapDBErr(err)
}

func (p *photoServiceImpl) GetPhotoCount(ctx context.Context) (int64, error) {
	photoDAL := dal.GetQueryByCtx(ctx).Photo
	count, err := photoDAL.WithContext(ctx).Count()
	return count, WrapDBErr(err)
}
