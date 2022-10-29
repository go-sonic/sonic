package property

import "reflect"

var (
	SummaryLength = Property{
		KeyValue:     "post_summary_length",
		DefaultValue: 150,
		Kind:         reflect.Int,
	}
	RssPageSize = Property{
		KeyValue:     "rss_page_size",
		DefaultValue: 20,
		Kind:         reflect.Int,
	}
	RssContentType = Property{
		KeyValue:     "rss_content_type",
		DefaultValue: "full",
		Kind:         reflect.String,
	}
	IndexPageSize = Property{
		KeyValue:     "post_index_page_size",
		DefaultValue: 10,
		Kind:         reflect.Int,
	}
	ArchivePageSize = Property{
		KeyValue:     "post_archives_page_size",
		DefaultValue: 10,
		Kind:         reflect.Int,
	}
	IndexSort = Property{
		KeyValue:     "post_index_sort",
		DefaultValue: "createTime",
		Kind:         reflect.String,
	}
	RecycledPostCleaningEnabled = Property{
		KeyValue:     "recycled_post_cleaning_enabled",
		DefaultValue: false,
		Kind:         reflect.Bool,
	}
	RecycledPostRetentionTime = Property{
		KeyValue:     "recycled_post_retention_time",
		DefaultValue: 30,
		Kind:         reflect.Int,
	}
	RecycledPostRetentionTimeunit = Property{
		KeyValue:     "recycled_post_retention_timeunit",
		DefaultValue: "DAY",
		Kind:         reflect.String,
	}
)
