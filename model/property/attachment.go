package property

import (
	"reflect"

	"github.com/go-sonic/sonic/consts"
)

var UploadImagePreviewEnable = Property{
	KeyValue:     "attachment_upload_image_preview_enable",
	DefaultValue: true,
	Kind:         reflect.Bool,
}

var UploadMaxParallelUploads = Property{
	KeyValue:     "attachment_upload_max_parallel_uploads",
	DefaultValue: 3,
	Kind:         reflect.Int,
}

var UploadMaxFiles = Property{
	KeyValue:     "attachment_upload_max_files",
	DefaultValue: 50,
	Kind:         reflect.Int,
}

var AttachmentType = Property{
	KeyValue:     "attachment_type",
	DefaultValue: consts.AttachmentTypeLocal.String(),
	Kind:         reflect.String,
}

var MinioEndpoint = Property{
	DefaultValue: "",
	KeyValue:     "minio_endpoint",
	Kind:         reflect.String,
}

var MinioBucketName = Property{
	DefaultValue: "",
	KeyValue:     "minio_bucket_name",
	Kind:         reflect.String,
}

var MinioAccessKey = Property{
	DefaultValue: "",
	KeyValue:     "minio_access_key",
	Kind:         reflect.String,
}

var MinioAccessSecret = Property{
	DefaultValue: "",
	KeyValue:     "minio_access_secret",
	Kind:         reflect.String,
}

var MinioProtocol = Property{
	DefaultValue: "https://",
	KeyValue:     "minio_protocol",
	Kind:         reflect.String,
}

var MinioSource = Property{
	DefaultValue: "",
	KeyValue:     "minio_source",
	Kind:         reflect.String,
}

var MinioRegion = Property{
	DefaultValue: "",
	KeyValue:     "minio_region",
	Kind:         reflect.String,
}

var MinioFrontBase = Property{
	DefaultValue: "",
	KeyValue:     "minio_front_base",
	Kind:         reflect.String,
}

var AliOssEndpoint = Property{
	DefaultValue: "",
	KeyValue:     "oss_ali_endpoint",
	Kind:         reflect.String,
}

var AliOssBucketName = Property{
	DefaultValue: "",
	KeyValue:     "oss_ali_bucket_name",
	Kind:         reflect.String,
}

var AliOssAccessKey = Property{
	DefaultValue: "",
	KeyValue:     "oss_ali_access_key",
	Kind:         reflect.String,
}

var AliOssDomain = Property{
	DefaultValue: "",
	KeyValue:     "oss_ali_domain",
	Kind:         reflect.String,
}

var AliOssProtocol = Property{
	DefaultValue: "https://",
	KeyValue:     "oss_ali_domain_protocol",
	Kind:         reflect.String,
}

var AliOssAccessSecret = Property{
	DefaultValue: "",
	KeyValue:     "oss_ali_access_secret",
	Kind:         reflect.String,
}

var AliOssSource = Property{
	DefaultValue: "",
	KeyValue:     "oss_ali_source",
	Kind:         reflect.String,
}

var AliOssThumbnailStyleRule = Property{
	DefaultValue: "",
	KeyValue:     "oss_ali_thumbnail_style_rule",
	Kind:         reflect.String,
}

var AliOssStyleRule = Property{
	DefaultValue: "",
	KeyValue:     "oss_ali_style_rule",
	Kind:         reflect.String,
}

var HuaweiOssDomain = Property{
	DefaultValue: "",
	KeyValue:     "oss_huawei_domain",
	Kind:         reflect.String,
}

var HuaweiOssEndpoint = Property{
	DefaultValue: "",
	KeyValue:     "oss_huawei_endpoint",
	Kind:         reflect.String,
}

var HuaweiOssBucketName = Property{
	DefaultValue: "",
	KeyValue:     "oss_huawei_bucket_name",
	Kind:         reflect.String,
}

var HuaweiOssAccessKey = Property{
	DefaultValue: "",
	KeyValue:     "oss_huawei_access_key",
	Kind:         reflect.String,
}

var HuaweiOssAccessSecret = Property{
	DefaultValue: "",
	KeyValue:     "oss_huawei_access_secret",
	Kind:         reflect.String,
}

var QiniuOssAccessKey = Property{
	DefaultValue: "",
	KeyValue:     "oss_qiniu_access_key",
	Kind:         reflect.String,
}

var QiniuOssAccessSecret = Property{
	DefaultValue: "",
	KeyValue:     "oss_qiniu_access_secret",
	Kind:         reflect.String,
}

var QiniuOssDomain = Property{
	DefaultValue: "",
	KeyValue:     "oss_qiniu_domain",
	Kind:         reflect.String,
}

var QiniuOssBucket = Property{
	DefaultValue: "",
	KeyValue:     "oss_qiniu_bucket",
	Kind:         reflect.String,
}

var QiniuDomainProtocol = Property{
	DefaultValue: "https://",
	KeyValue:     "oss_qiniu_domain_protocol",
	Kind:         reflect.String,
}

var QiniuOssStyleRule = Property{
	DefaultValue: "",
	KeyValue:     "oss_qiniu_style_rule",
	Kind:         reflect.String,
}

var QiniuOssThumbnailStyleRule = Property{
	DefaultValue: "",
	KeyValue:     "oss_qiniu_thumbnail_style_rule",
	Kind:         reflect.String,
}

var QiniuOssZone = Property{
	DefaultValue: "",
	KeyValue:     "oss_qiniu_zone",
	Kind:         reflect.String,
}

var TencentCosDomain = Property{
	DefaultValue: "",
	KeyValue:     "cos_tencent_domain",
	Kind:         reflect.String,
}

var TencentCosProtocol = Property{
	DefaultValue: "https://",
	KeyValue:     "cos_tencent_protocol",
	Kind:         reflect.String,
}

var TencentCosRegion = Property{
	DefaultValue: "",
	KeyValue:     "cos_tencent_region",
	Kind:         reflect.String,
}

var TencentCosBucketName = Property{
	DefaultValue: "",
	KeyValue:     "cos_tencent_bucket_name",
	Kind:         reflect.String,
}

var TencentCosSecretID = Property{
	DefaultValue: "",
	KeyValue:     "cos_tencent_secret_id",
	Kind:         reflect.String,
}

var TencentCosSecretKey = Property{
	DefaultValue: "",
	KeyValue:     "cos_tencent_secret_key",
	Kind:         reflect.String,
}

var TencentCosSource = Property{
	DefaultValue: "",
	KeyValue:     "cos_tencent_source",
	Kind:         reflect.String,
}

var TencentCosStyleRule = Property{
	DefaultValue: "",
	KeyValue:     "cos_tencent_style_rule",
	Kind:         reflect.String,
}

var TencentCosThumbnailStyleRule = Property{
	DefaultValue: "",
	KeyValue:     "cos_tencent_thumbnail_style_rule",
	Kind:         reflect.String,
}

var UpOssSource = Property{
	DefaultValue: "",
	KeyValue:     "oss_upyun_source",
	Kind:         reflect.String,
}

var UpOssPassword = Property{
	DefaultValue: "",
	KeyValue:     "oss_upyun_password",
	Kind:         reflect.String,
}

var UpOssBucket = Property{
	DefaultValue: "",
	KeyValue:     "oss_upyun_name",
	Kind:         reflect.String,
}

var UpOssDomain = Property{
	DefaultValue: "",
	KeyValue:     "oss_upyun_domain",
	Kind:         reflect.String,
}

var UpOssProtocol = Property{
	DefaultValue: "http://",
	KeyValue:     "oss_upyun_protocol",
	Kind:         reflect.String,
}

var UpOssOperator = Property{
	DefaultValue: "",
	KeyValue:     "oss_upyun_operator",
	Kind:         reflect.String,
}

var UpOssStyleRule = Property{
	DefaultValue: "",
	KeyValue:     "oss_upyun_style_rule",
	Kind:         reflect.String,
}

var UpOssThumbnailStyleRule = Property{
	DefaultValue: "",
	KeyValue:     "oss_upyun_thumbnail_style_rule",
	Kind:         reflect.String,
}
