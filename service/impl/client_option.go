package impl

import (
	"context"

	"go.uber.org/zap"

	"github.com/go-sonic/sonic/cache"
	"github.com/go-sonic/sonic/config"
	"github.com/go-sonic/sonic/event"
	"github.com/go-sonic/sonic/model/dto"
	"github.com/go-sonic/sonic/model/property"
	"github.com/go-sonic/sonic/service"
)

type clientOptionServiceImpl struct {
	service.OptionService
	Cache         cache.Cache
	Config        *config.Config
	Event         event.Bus
	Logger        *zap.Logger
	PrivateOption map[string]struct{}
}

func NewClientOptionService(config *config.Config, cache cache.Cache, event event.Bus, logger *zap.Logger, optionService service.OptionService) service.ClientOptionService {
	co := &clientOptionServiceImpl{
		Cache:         cache,
		Config:        config,
		Event:         event,
		Logger:        logger,
		OptionService: optionService,
	}
	co.PrivateOption = co.getPrivateOption()
	return co
}

func (c *clientOptionServiceImpl) ListAllOption(ctx context.Context) ([]*dto.Option, error) {
	options, err := c.OptionService.ListAllOption(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]*dto.Option, 0)
	for _, option := range options {
		if _, ok := c.PrivateOption[option.Key]; !ok {
			result = append(result, option)
		}
	}
	return result, nil
}

func (c *clientOptionServiceImpl) getPrivateOption() map[string]struct{} {
	privateOption := make(map[string]struct{})
	privateProperty := []property.Property{
		property.EmailProtocol,
		property.EmailSSLPort,
		property.EmailUsername,
		property.EmailPassword,
		property.EmailFromName,
		property.APIAccessKey,
		property.MinioEndpoint,
		property.MinioBucketName,
		property.MinioAccessKey,
		property.MinioAccessSecret,
		property.MinioProtocol,
		property.MinioSource,
		property.MinioRegion,
		property.MinioFrontBase,
		property.AliOssEndpoint,
		property.AliOssBucketName,
		property.AliOssAccessKey,
		property.AliOssDomain,
		property.AliOssProtocol,
		property.AliOssAccessSecret,
		property.HuaweiOssDomain,
		property.HuaweiOssEndpoint,
		property.HuaweiOssBucketName,
		property.HuaweiOssAccessKey,
		property.HuaweiOssAccessSecret,
		property.QiniuOssAccessKey,
		property.QiniuOssAccessSecret,
		property.QiniuOssDomain,
		property.QiniuOssBucket,
		property.QiniuDomainProtocol,
		property.QiniuOssStyleRule,
		property.QiniuOssThumbnailStyleRule,
		property.QiniuOssZone,
		property.TencentCosDomain,
		property.TencentCosProtocol,
		property.TencentCosRegion,
		property.TencentCosBucketName,
		property.TencentCosSecretID,
		property.TencentCosSecretKey,
		property.TencentCosSource,
		property.TencentCosStyleRule,
		property.TencentCosThumbnailStyleRule,
		property.UpOssSource,
		property.UpOssPassword,
		property.UpOssBucket,
		property.UpOssDomain,
		property.UpOssProtocol,
		property.UpOssOperator,
		property.UpOssStyleRule,
		property.UpOssThumbnailStyleRule,
		property.JWTSecret,
	}
	for _, p := range privateProperty {
		privateOption[p.KeyValue] = struct{}{}
	}
	return privateOption
}
