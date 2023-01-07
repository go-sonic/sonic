package property

import "reflect"

var (
	BlogLocale = Property{
		KeyValue:     "blog_locale",
		DefaultValue: "zh",
		Kind:         reflect.String,
	}
	BlogTitle = Property{
		KeyValue:     "blog_title",
		DefaultValue: "",
		Kind:         reflect.String,
	}
	BlogLogo = Property{
		KeyValue:     "blog_logo",
		DefaultValue: "",
		Kind:         reflect.String,
	}
	BlogURL = Property{
		KeyValue:     "blog_url",
		DefaultValue: "",
		Kind:         reflect.String,
	}
	BlogFavicon = Property{
		KeyValue:     "blog_favicon",
		DefaultValue: "",
		Kind:         reflect.String,
	}
	BlogFooterInfo = Property{
		KeyValue:     "blog_footer_info",
		DefaultValue: "",
		Kind:         reflect.String,
	}
)
