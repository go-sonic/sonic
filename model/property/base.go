package property

import (
	"reflect"
	"strconv"

	"github.com/go-sonic/sonic/model/entity"
)

type Property struct {
	DefaultValue interface{}
	KeyValue     string
	Kind         reflect.Kind
}

func (p Property) ConvertToOption() *entity.Option {
	var value string
	switch p.Kind {
	case reflect.Bool:
		value = strconv.FormatBool(p.DefaultValue.(bool))
	case reflect.Int:
		value = strconv.FormatInt(int64(p.DefaultValue.(int)), 10)
	case reflect.Int32:
		value = strconv.FormatInt(int64(p.DefaultValue.(int32)), 10)
	case reflect.Int64:
		value = strconv.FormatInt(p.DefaultValue.(int64), 10)
	case reflect.String:
		if p.DefaultValue != nil {
			value = p.DefaultValue.(string)
		}
	}
	return &entity.Option{
		OptionKey:   p.KeyValue,
		OptionValue: value,
	}
}

var AllProperty = []Property{
	UploadImagePreviewEnable,
	UploadMaxParallelUploads,
	UploadMaxFiles,
	AttachmentType,
	BlogLocale,
	BlogTitle,
	BlogLogo,
	BlogURL,
	BlogFavicon,
	BlogFooterInfo,
	EmailHost,
	EmailProtocol,
	EmailSSLPort,
	EmailUsername,
	EmailPassword,
	EmailFromName,
	EmailIsEnabled,
	EmailStarttls,
	CustomHead,
	CustomContentHead,
	StatisticsCode,
	GlobalAbsolutePathEnabled,
	DefaultEditor,
	PostPermalinkType,
	SheetPermalinkType,
	CategoriesPrefix,
	TagsPrefix,
	ArchivesPrefix,
	SheetPrefix,
	LinksPrefix,
	PhotosPrefix,
	JournalsPrefix,
	PathSuffix,
	IsInstalled,
	Theme,
	BirthDay,
	DefaultMenuTeam,
	SeoKeywords,
	SeoDescription,
	SeoSpiderDisabled,
	SummaryLength,
	RssPageSize,
	RssContentType,
	IndexPageSize,
	ArchivePageSize,
	IndexSort,
	RecycledPostCleaningEnabled,
	RecycledPostRetentionTime,
	RecycledPostRetentionTimeunit,
	APIAccessKey,
	CommentGravatarDefault,
	CommentNewNeedCheck,
	CommentNewNotice,
	CommentReplyNotice,
	CommentAPIEnabled,
	CommentPageSize,
	CommentContentPlaceholder,
	CommentInternalPluginJs,
	CommentGravatarSource,
	CommentBanTime,
	CommentRange,
	MinioEndpoint,
	MinioBucketName,
	MinioAccessKey,
	MinioAccessSecret,
	MinioProtocol,
	MinioSource,
	MinioRegion,
	MinioFrontBase,
	AliOssEndpoint,
	AliOssBucketName,
	AliOssAccessKey,
	AliOssDomain,
	AliOssProtocol,
	AliOssAccessSecret,
	AliOssSource,
	AliOssStyleRule,
	AliOssThumbnailStyleRule,
	HuaweiOssDomain,
	HuaweiOssEndpoint,
	HuaweiOssBucketName,
	HuaweiOssAccessKey,
	HuaweiOssAccessSecret,
	QiniuOssAccessKey,
	QiniuOssAccessSecret,
	QiniuOssDomain,
	QiniuOssBucket,
	QiniuDomainProtocol,
	QiniuOssStyleRule,
	QiniuOssThumbnailStyleRule,
	QiniuOssZone,
	TencentCosDomain,
	TencentCosProtocol,
	TencentCosRegion,
	TencentCosBucketName,
	TencentCosSecretID,
	TencentCosSecretKey,
	TencentCosSource,
	TencentCosStyleRule,
	TencentCosThumbnailStyleRule,
	UpOssSource,
	UpOssPassword,
	UpOssBucket,
	UpOssDomain,
	UpOssProtocol,
	UpOssOperator,
	UpOssStyleRule,
	UpOssThumbnailStyleRule,
	PhotoPageSize,
	JournalPageSize,
	JWTSecret,
}
