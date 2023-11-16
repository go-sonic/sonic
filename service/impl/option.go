package impl

import (
	"context"
	"fmt"
	"reflect"
	"strconv"

	"go.uber.org/zap"

	"github.com/go-sonic/sonic/cache"
	"github.com/go-sonic/sonic/config"
	"github.com/go-sonic/sonic/consts"
	"github.com/go-sonic/sonic/dal"
	"github.com/go-sonic/sonic/event"
	"github.com/go-sonic/sonic/log"
	"github.com/go-sonic/sonic/model/dto"
	"github.com/go-sonic/sonic/model/entity"
	"github.com/go-sonic/sonic/model/param"
	"github.com/go-sonic/sonic/model/property"
	"github.com/go-sonic/sonic/service"
	"github.com/go-sonic/sonic/util/xerr"
)

type optionServiceImpl struct {
	Cache  cache.Cache
	Config *config.Config
	Event  event.Bus
	Logger *zap.Logger
}

func NewOptionService(config *config.Config, cache cache.Cache, event event.Bus, logger *zap.Logger) service.OptionService {
	return &optionServiceImpl{
		Cache:  cache,
		Config: config,
		Event:  event,
		Logger: logger,
	}
}

func (o *optionServiceImpl) GetOrByDefault(ctx context.Context, p property.Property) interface{} {
	value, err := o.getFromCacheMissFromDB(ctx, p)
	if xerr.GetType(err) == xerr.NoRecord {
		o.Cache.SetDefault(p.KeyValue, p.DefaultValue)
		return p.DefaultValue
	}
	if err != nil {
		o.Logger.Error("get option", zap.String("key", p.KeyValue), zap.Error(err))
		return p.DefaultValue
	}
	return value
}

func (o *optionServiceImpl) GetPostSummaryLength(ctx context.Context) int {
	p := property.SummaryLength
	value, err := o.getFromCacheMissFromDB(ctx, p)

	if xerr.GetType(err) == xerr.NoRecord {
		o.Cache.SetDefault(p.KeyValue, p.DefaultValue)
		return p.DefaultValue.(int)
	} else if err != nil {
		log.CtxErrorf(ctx, "query option key=%v err=%v", p.KeyValue, err)
		return p.DefaultValue.(int)
	}
	return value.(int)
}

func (o *optionServiceImpl) GetPostSort(ctx context.Context) param.Sort {
	p := property.IndexSort
	value, err := o.getFromCacheMissFromDB(ctx, p)
	sort := p.DefaultValue.(string)

	//nolint:gocritic
	if xerr.GetType(err) == xerr.NoRecord {
		o.Cache.SetDefault(p.KeyValue, p.DefaultValue)
	} else if err != nil {
		log.CtxErrorf(ctx, "query option err=%v", err)
	} else {
		sort = value.(string)
	}
	return param.Sort{
		Fields: []string{"topPriority,desc", sort + ",desc", "id,desc"},
	}
}

func (o *optionServiceImpl) GetIndexPageSize(ctx context.Context) int {
	p := property.IndexPageSize
	value, err := o.getFromCacheMissFromDB(ctx, p)
	if xerr.GetType(err) == xerr.NoRecord {
		o.Cache.SetDefault(p.KeyValue, p.DefaultValue)
		return p.DefaultValue.(int)
	} else if err != nil {
		log.CtxErrorf(ctx, "query option err=%v", err)
		return p.DefaultValue.(int)
	}
	return value.(int)
}

func (o *optionServiceImpl) GetPostPermalinkType(ctx context.Context) (consts.PostPermalinkType, error) {
	p := property.PostPermalinkType
	value, err := o.getFromCacheMissFromDB(ctx, p)
	if xerr.GetType(err) == xerr.NoRecord {
		o.Cache.SetDefault(p.KeyValue, p.DefaultValue)
		return consts.PostPermalinkType(p.DefaultValue.(string)), nil
	} else if err != nil {
		return "", err
	}
	return consts.PostPermalinkType(value.(string)), nil
}

func (o *optionServiceImpl) GetOrByDefaultWithErr(ctx context.Context, p property.Property, defaultValue interface{}) (interface{}, error) {
	value, err := o.getFromCacheMissFromDB(ctx, p)
	if xerr.GetType(err) == xerr.NoRecord {
		o.Cache.SetDefault(p.KeyValue, defaultValue)
		return defaultValue, nil
	}
	if err != nil {
		return nil, xerr.WithStatus(err, xerr.StatusInternalServerError)
	}
	return value, nil
}

func (o *optionServiceImpl) GetBlogBaseURL(ctx context.Context) (string, error) {
	blogURL, err := o.GetOrByDefaultWithErr(ctx, property.BlogURL, "")
	if err != nil {
		return "", err
	}
	if blogURL != "" {
		return blogURL.(string), nil
	}
	if o.Config.Server.Host == "0.0.0.0" {
		return fmt.Sprintf("http://127.0.0.1:%s", o.Config.Server.Port), nil
	} else {
		return fmt.Sprintf("http://%s:%s", o.Config.Server.Host, o.Config.Server.Port), nil
	}
}

func (o *optionServiceImpl) IsEnabledAbsolutePath(ctx context.Context) (bool, error) {
	isEnabled, err := o.GetOrByDefaultWithErr(ctx, property.GlobalAbsolutePathEnabled, true)
	if err != nil {
		return true, err
	}
	return isEnabled.(bool), nil
}

func (o *optionServiceImpl) GetPathSuffix(ctx context.Context) (string, error) {
	p := property.PathSuffix
	value, err := o.getFromCacheMissFromDB(ctx, p)
	if xerr.GetType(err) == xerr.NoRecord {
		o.Cache.SetDefault(p.KeyValue, p.DefaultValue)
		return p.DefaultValue.(string), nil
	} else if err != nil {
		return "", err
	}
	return value.(string), nil
}

func (o *optionServiceImpl) getFromCacheMissFromDB(ctx context.Context, p property.Property) (interface{}, error) {
	value, ok := o.Cache.Get(p.KeyValue)
	if ok {
		return value, nil
	}

	optionDAL := dal.GetQueryByCtx(ctx).Option
	option, err := optionDAL.WithContext(ctx).Where(optionDAL.OptionKey.Eq(p.KeyValue)).Take()
	if err != nil {
		return nil, WrapDBErr(err)
	}

	value, err = o.convert(option, p)
	if err != nil {
		return nil, err
	}
	o.Cache.SetDefault(p.KeyValue, value)
	return value, nil
}

func (o *optionServiceImpl) convert(option *entity.Option, p property.Property) (interface{}, error) {
	var err error
	var result interface{}
	switch p.Kind {
	case reflect.Bool:
		result, err = strconv.ParseBool(option.OptionValue)
	case reflect.Int:
		result, err = strconv.Atoi(option.OptionValue)
	case reflect.Int32:
		v, e := strconv.ParseInt(option.OptionValue, 10, 32)
		result = int32(v)
		err = e
	case reflect.Int64:
		result, err = strconv.ParseInt(option.OptionValue, 10, 64)
	case reflect.String:
		result, err = option.OptionValue, nil
	}
	if err != nil {
		return nil, xerr.BadParam.Wrapf(err, "option 类型错误 optionValue=%v kind=%v", option.OptionValue, p.Kind)
	}
	return result, nil
}

func (o *optionServiceImpl) Save(ctx context.Context, saveMap map[string]string) (err error) {
	propertyMap := o.OptionMap()
	optionDAL := dal.GetQueryByCtx(ctx).Option
	options, err := optionDAL.WithContext(ctx).Find()
	if err != nil {
		return WrapDBErr(err)
	}
	optionKeyMap := make(map[string]*entity.Option)
	for _, option := range options {
		optionKeyMap[option.OptionKey] = option
	}
	toCreates := make([]*entity.Option, 0)
	toUpdates := make([]*entity.Option, 0)
	for key, value := range saveMap {
		p, ok := propertyMap[key]
		if !ok {
			return xerr.BadParam.New("key=%v", key).WithMsg("option key not exist").WithStatus(xerr.StatusBadRequest)
		}
		temp := &entity.Option{
			OptionKey:   key,
			OptionValue: value,
		}
		// check type
		_, err := o.convert(temp, p)
		if err != nil {
			return err
		}
		option, ok := optionKeyMap[key]
		if ok {
			option.OptionValue = value
			toUpdates = append(toUpdates, option)
		} else {
			toCreates = append(toCreates, &entity.Option{
				OptionKey:   key,
				OptionValue: value,
			})
		}
	}

	// Update the database before deleting the cache.
	// Although there is a very small probability that this will lead to temporary cache and database data inconsistency,
	// but it is acceptable

	deleteKeys := make([]string, 0, len(toUpdates)+len(toCreates))
	for _, option := range toCreates {
		deleteKeys = append(deleteKeys, option.OptionKey)
	}
	for _, option := range toUpdates {
		deleteKeys = append(deleteKeys, option.OptionKey)
	}
	o.Cache.BatchDelete(deleteKeys)

	err = dal.GetQueryByCtx(ctx).Transaction(func(tx *dal.Query) error {
		optionDAL := tx.Option
		for _, toUpdate := range toUpdates {
			_, err := optionDAL.WithContext(ctx).Where(optionDAL.ID.Eq(toUpdate.ID), optionDAL.OptionKey.Eq(toUpdate.OptionKey)).UpdateColumnSimple(optionDAL.OptionValue.Value(toUpdate.OptionValue))
			if err != nil {
				return WrapDBErr(err)
			}
		}
		err := optionDAL.WithContext(ctx).Create(toCreates...)
		if err != nil {
			return WrapDBErr(err)
		}
		return nil
	})
	if err != nil {
		return err
	}
	o.Event.Publish(ctx, &event.OptionUpdateEvent{})
	return nil
}

func (o *optionServiceImpl) GetArchivePrefix(ctx context.Context) (string, error) {
	p := property.ArchivesPrefix
	value, err := o.getFromCacheMissFromDB(ctx, p)
	if xerr.GetType(err) == xerr.NoRecord {
		o.Cache.SetDefault(p.KeyValue, p.DefaultValue)
		return p.DefaultValue.(string), nil
	} else if err != nil {
		return "", err
	}
	return value.(string), nil
}

func (o *optionServiceImpl) ListAllOption(ctx context.Context) ([]*dto.Option, error) {
	optionDAL := dal.GetQueryByCtx(ctx).Option
	options, err := optionDAL.WithContext(ctx).Find()
	if err != nil {
		return nil, WrapDBErr(err)
	}

	result := make([]*dto.Option, 0)
	optionKeyMap := make(map[string]struct{})
	propertyMap := o.OptionMap()
	for _, option := range options {
		p, ok := propertyMap[option.OptionKey]
		if !ok {
			panic("property not exist")
		}
		value, err := o.convert(option, p)
		if err != nil {
			return nil, err
		}
		result = append(result, &dto.Option{
			Key:   option.OptionKey,
			Value: value,
		})
		optionKeyMap[option.OptionKey] = struct{}{}
	}
	for _, p := range propertyMap {
		if _, ok := optionKeyMap[p.KeyValue]; ok {
			continue
		}
		result = append(result, &dto.Option{
			Key:   p.KeyValue,
			Value: p.DefaultValue,
		})
	}
	return result, nil
}

func (o *optionServiceImpl) OptionMap() map[string]property.Property {
	result := make(map[string]property.Property)
	for _, p := range property.AllProperty {
		result[p.KeyValue] = p
	}
	return result
}

func (o *optionServiceImpl) GetLinksPrefix(ctx context.Context) (string, error) {
	p := property.LinksPrefix
	value, err := o.getFromCacheMissFromDB(ctx, p)
	if xerr.GetType(err) == xerr.NoRecord {
		o.Cache.SetDefault(p.KeyValue, p.DefaultValue)
		return p.DefaultValue.(string), nil
	} else if err != nil {
		return "", err
	}
	return value.(string), nil
}

func (o *optionServiceImpl) GetPhotoPrefix(ctx context.Context) (string, error) {
	p := property.PhotosPrefix
	value, err := o.getFromCacheMissFromDB(ctx, p)
	if xerr.GetType(err) == xerr.NoRecord {
		o.Cache.SetDefault(p.KeyValue, p.DefaultValue)
		return p.DefaultValue.(string), nil
	} else if err != nil {
		return "", err
	}
	return value.(string), nil
}

func (o *optionServiceImpl) GetJournalPrefix(ctx context.Context) (string, error) {
	p := property.JournalsPrefix
	value, err := o.getFromCacheMissFromDB(ctx, p)
	if xerr.GetType(err) == xerr.NoRecord {
		o.Cache.SetDefault(p.KeyValue, p.DefaultValue)
		return p.DefaultValue.(string), nil
	} else if err != nil {
		return "", err
	}
	return value.(string), nil
}

func (o *optionServiceImpl) GetActivatedThemeID(ctx context.Context) (string, error) {
	p := property.Theme
	value, err := o.getFromCacheMissFromDB(ctx, p)
	if xerr.GetType(err) == xerr.NoRecord {
		o.Cache.SetDefault(p.KeyValue, p.DefaultValue)
		return p.DefaultValue.(string), nil
	} else if err != nil {
		return "", err
	}
	return value.(string), nil
}

func (o *optionServiceImpl) GetCategoryPrefix(ctx context.Context) (string, error) {
	p := property.CategoriesPrefix
	value, err := o.getFromCacheMissFromDB(ctx, p)
	if xerr.GetType(err) == xerr.NoRecord {
		o.Cache.SetDefault(p.KeyValue, p.DefaultValue)
		return p.DefaultValue.(string), nil
	} else if err != nil {
		return "", err
	}
	return value.(string), nil
}

func (o *optionServiceImpl) GetTagPrefix(ctx context.Context) (string, error) {
	p := property.TagsPrefix
	value, err := o.getFromCacheMissFromDB(ctx, p)
	if xerr.GetType(err) == xerr.NoRecord {
		o.Cache.SetDefault(p.KeyValue, p.DefaultValue)
		return p.DefaultValue.(string), nil
	} else if err != nil {
		return "", err
	}
	return value.(string), nil
}

func (o *optionServiceImpl) GetLinkPrefix(ctx context.Context) (string, error) {
	p := property.LinksPrefix
	value, err := o.getFromCacheMissFromDB(ctx, p)
	if xerr.GetType(err) == xerr.NoRecord {
		o.Cache.SetDefault(p.KeyValue, p.DefaultValue)
		return p.DefaultValue.(string), nil
	} else if err != nil {
		return "", err
	}
	return value.(string), nil
}

func (o *optionServiceImpl) GetSheetPermalinkType(ctx context.Context) (consts.SheetPermaLinkType, error) {
	p := property.SheetPermalinkType
	value, err := o.getFromCacheMissFromDB(ctx, p)
	if xerr.GetType(err) == xerr.NoRecord {
		o.Cache.SetDefault(p.KeyValue, p.DefaultValue)
		return consts.SheetPermaLinkType(p.DefaultValue.(string)), nil
	} else if err != nil {
		return "", err
	}
	return consts.SheetPermaLinkType(value.(string)), nil
}

func (o *optionServiceImpl) GetSheetPrefix(ctx context.Context) (string, error) {
	p := property.SheetPrefix
	value, err := o.getFromCacheMissFromDB(ctx, p)
	if xerr.GetType(err) == xerr.NoRecord {
		o.Cache.SetDefault(p.KeyValue, p.DefaultValue)
		return p.DefaultValue.(string), nil
	} else if err != nil {
		return "", err
	}
	return value.(string), nil
}

func (o *optionServiceImpl) GetAttachmentType(ctx context.Context) consts.AttachmentType {
	p := property.AttachmentType
	value, err := o.getFromCacheMissFromDB(ctx, p)
	if xerr.GetType(err) == xerr.NoRecord {
		o.Cache.SetDefault(p.KeyValue, p.DefaultValue)
		return consts.AttachmentTypeLocal
	} else if err != nil {
		return consts.AttachmentTypeLocal
	}

	switch value.(string) {
	case "LOCAL":
		return consts.AttachmentTypeLocal
	case "UPOSS":
		return consts.AttachmentTypeUpOSS
	case "QINIUOSS":
		return consts.AttachmentTypeQiNiuOSS
	case "AttachmentTypeSMMS":
		return consts.AttachmentTypeSMMS
	case "ALIOSS":
		return consts.AttachmentTypeAliOSS
	case "BAIDUOSS":
		return consts.AttachmentTypeBaiDuOSS
	case "TENCENTOSS":
		return consts.AttachmentTypeTencentCOS
	case "HUAWEIOBS":
		return consts.AttachmentTypeHuaweiOBS
	case "MINIO":
		return consts.AttachmentTypeMinIO
	default:
		return consts.AttachmentTypeLocal
	}
}

func (o *optionServiceImpl) GetAdminURLPath(ctx context.Context) (string, error) {
	return o.Config.Sonic.AdminURLPath, nil
}
