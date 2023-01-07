package property

import "reflect"

var (
	CustomHead = Property{
		DefaultValue: "",
		KeyValue:     "blog_custom_head",
		Kind:         reflect.String,
	}
	CustomContentHead = Property{
		DefaultValue: "",
		KeyValue:     "blog_custom_content_head",
		Kind:         reflect.String,
	}
	StatisticsCode = Property{
		DefaultValue: "",
		KeyValue:     "blog_statistics_code",
		Kind:         reflect.String,
	}
	GlobalAbsolutePathEnabled = Property{
		DefaultValue: true,
		KeyValue:     "global_absolute_path_enabled",
		Kind:         reflect.Bool,
	}
	DefaultEditor = Property{
		DefaultValue: "MARKDOWN",
		KeyValue:     "default_editor",
		Kind:         reflect.String,
	}
	APIAccessKey = Property{
		DefaultValue: nil,
		KeyValue:     "api_access_key",
		Kind:         reflect.String,
	}
	PhotoPageSize = Property{
		DefaultValue: 10,
		KeyValue:     "photos_page_size",
		Kind:         reflect.Int,
	}
	JournalPageSize = Property{
		DefaultValue: 10,
		KeyValue:     "journals_page_size",
		Kind:         reflect.Int,
	}
	JWTSecret = Property{
		DefaultValue: "",
		KeyValue:     "jwt_secret",
		Kind:         reflect.String,
	}
)
