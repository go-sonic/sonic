package impl

import (
	"context"
	"net/url"

	"github.com/go-sonic/sonic/dal"
	"github.com/go-sonic/sonic/model/entity"
	"github.com/go-sonic/sonic/model/param"
	"github.com/go-sonic/sonic/service"
	"github.com/go-sonic/sonic/util/xerr"
)

type scrapServiceImpl struct{}

func NewScrapService() service.ScrapService {
	return &scrapServiceImpl{}
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

func (impl *scrapServiceImpl) Create(ctx context.Context, pageParam *param.ScrapPage) error {
	pageEntity, err := convertToModel(pageParam)
	if err != nil {
		return xerr.BadParam.Wrap(err)
	}

	scrapPageDAL := dal.GetQueryByCtx(ctx).ScrapPage
	err = scrapPageDAL.WithContext(ctx).Create(pageEntity)
	if err != nil {
		return WrapDBErr(err)
	}

	return nil
}

func convertToModel(pageParam *param.ScrapPage) (*entity.ScrapPage, error) {
	pageURL, err := url.Parse(pageParam.URL)
	if err != nil {
		return nil, err
	}

	return &entity.ScrapPage{
		ID:       nil,
		Title:    pageParam.Title,
		Md5:      pageParam.Md5,
		URL:      pageParam.URL,
		Content:  pageParam.Content,
		Summary:  pageParam.Summary,
		CreateAt: pageParam.AddAt,
		Domain:   pageURL.Hostname(),
		OutLink:  nil,
	}, nil
}
