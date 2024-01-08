package impl

import (
	"context"
	"mime/multipart"
	"net/url"
	"time"

	"github.com/go-sonic/sonic/dal"
	"github.com/go-sonic/sonic/model/dto"
	"github.com/go-sonic/sonic/model/entity"
	"github.com/go-sonic/sonic/model/param"
	"github.com/go-sonic/sonic/service"
	"github.com/go-sonic/sonic/util/xerr"
)

type scrapServiceImpl struct {
	AttachmentService service.AttachmentService
}

type ScrapPageWithAttach struct {
	entity.ScrapPage
	entity.Attachment
}

func NewScrapService(attachmentService service.AttachmentService) service.ScrapService {
	return &scrapServiceImpl{
		AttachmentService: attachmentService,
	}
}

func (impl *scrapServiceImpl) QueryMd5List(ctx context.Context) ([]string, error) {
	scrapPageDAL := dal.GetQueryByCtx(ctx).ScrapPage
	entities, err := scrapPageDAL.WithContext(ctx).Select(dal.ScrapPage.Md5).Find()
	if err != nil {
		return nil, WrapDBErr(err)
	}

	md5List := make([]string, len(entities))
	for _, v := range entities {
		md5List = append(md5List, v.Md5)
	}

	return md5List, nil
}

func (impl *scrapServiceImpl) Create(ctx context.Context, pageParam *param.ScrapPage, file *multipart.FileHeader) (*dto.ScrapPageDTO, error) {
	attachmentDTO, err := impl.AttachmentService.Upload(ctx, file)
	if err != nil {
		return nil, err
	}

	pageEntity, err := convertToModel(pageParam)
	if err != nil {
		return nil, xerr.BadParam.Wrap(err)
	}

	pageEntity.Attachment = &attachmentDTO.ID

	scrapPageDAL := dal.GetQueryByCtx(ctx).ScrapPage
	err = scrapPageDAL.WithContext(ctx).Create(pageEntity)
	if err != nil {
		return nil, WrapDBErr(err)
	}

	return convertToDTO(pageEntity, false), nil
}

func (impl *scrapServiceImpl) Get(ctx context.Context, pageID int32) (*dto.ScrapPageDTO, error) {
	scrapPageDAL := dal.GetQueryByCtx(ctx).ScrapPage
	pageEntity, err := scrapPageDAL.WithContext(ctx).Where(scrapPageDAL.ID.Eq(pageID)).First()
	if err != nil {
		return nil, WrapDBErr(err)
	}

	return convertToDTO(pageEntity, true), nil
}

func (impl *scrapServiceImpl) Query(ctx context.Context, query *param.ScrapPageQuery) ([]*dto.ScrapPageDTO, int64, error) {
	if query == nil || query.PageNum < 0 || query.PageSize <= 0 {
		return nil, 0, xerr.BadParam.New("").WithStatus(xerr.StatusBadRequest).WithMsg("Paging parameter error")
	}

	scrapPageDAL := dal.GetQueryByCtx(ctx).ScrapPage

	scrapDo := scrapPageDAL.WithContext(ctx)
	if len(query.KeyWord) > 0 {
		scrapDo = scrapDo.Where(scrapPageDAL.Title.Like(query.KeyWord))
	}

	total, err := scrapDo.Count()
	if err != nil {
		return nil, 0, WrapDBErr(err)
	}

	dtoList := make([]*dto.ScrapPageDTO, 0)

	db := dal.GetDB().Table("scrap_page").Select("scrap_page.*, attachment.id as attachment_id, attachment.file_key, attachment.media_type, attachment.name, attachment.path, attachment.size, attachment.suffix, attachment.thumb_path").Joins("left join attachment on attachment.id = scrap_page.attachment")

	if len(query.KeyWord) > 0 {
		db = db.Where("scrap_page.title like ?", "%"+query.KeyWord+"%")
	}

	db.Scan(&dtoList)

	return dtoList, total, nil
}

func (impl *scrapServiceImpl) Update(ctx context.Context, pageId int32, pageParam *param.ScrapPage) (*dto.ScrapPageDTO, error) {
	scrapPageDAL := dal.GetQueryByCtx(ctx).ScrapPage

	updatePage, err := convertToModel(pageParam)
	if err != nil {
		return nil, xerr.BadParam.New("").WithStatus(xerr.StatusBadRequest).WithMsg("invalid param")
	}
	updateResult, err := scrapPageDAL.WithContext(ctx).Where(scrapPageDAL.ID.Eq(pageId)).Omit(scrapPageDAL.CreateTime).Updates(updatePage)
	if err != nil {
		return nil, WrapDBErr(err)
	}
	if updateResult.RowsAffected != 1 {
		return nil, xerr.NoType.New("update scrap page failed").WithMsg("update scrap page failed").WithStatus(xerr.StatusInternalServerError)
	}

	scrapPage, err := scrapPageDAL.WithContext(ctx).Where(scrapPageDAL.ID.Eq(pageId)).First()
	if err != nil {
		return nil, WrapDBErr(err)
	}
	return convertToDTO(scrapPage, false), nil
}

func convertToDTO(pageEntity *entity.ScrapPage, withContent bool) *dto.ScrapPageDTO {
	var content string
	var summary string
	if withContent && pageEntity.Content != nil {
		content = *pageEntity.Content
	}

	if pageEntity.Summary != nil {
		summary = *pageEntity.Summary
	}

	return &dto.ScrapPageDTO{
		ID:        pageEntity.ID,
		Content:   content,
		Title:     pageEntity.Title,
		Summary:   summary,
		URL:       pageEntity.URL,
		OriginURL: pageEntity.OriginURL,
		Domain:    pageEntity.Domain,
		Resource:  pageEntity.Resource,
	}
}

func convertToModel(pageParam *param.ScrapPage) (*entity.ScrapPage, error) {
	pageURL, err := url.Parse(pageParam.URL)
	if err != nil {
		return nil, err
	}

	createTime := time.Now()
	if pageParam.AddAt != nil {
		createTime = time.Unix(*pageParam.AddAt, 0)
	}

	return &entity.ScrapPage{
		Title:      pageParam.Title,
		Md5:        pageParam.Md5,
		URL:        pageParam.URL,
		Content:    pageParam.Content,
		Summary:    pageParam.Summary,
		CreateTime: createTime,
		Domain:     pageURL.Hostname(),
		Resource:   pageParam.Resource,
	}, nil
}
