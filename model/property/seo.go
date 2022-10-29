package property

import "reflect"

var SeoKeywords = Property{
	KeyValue:     "seo_keywords",
	DefaultValue: "",
	Kind:         reflect.String,
}

var SeoDescription = Property{
	KeyValue:     "seo_description",
	DefaultValue: "",
	Kind:         reflect.String,
}

var SeoSpiderDisabled = Property{
	KeyValue:     "seo_spider_disabled",
	DefaultValue: false,
	Kind:         reflect.Bool,
}
