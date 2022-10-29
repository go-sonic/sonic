package property

import "reflect"

var (
	PostPermalinkType = Property{
		DefaultValue: "DEFAULT",
		KeyValue:     "post_permalink_type",
		Kind:         reflect.String,
	}
	SheetPermalinkType = Property{
		DefaultValue: "SECONDARY",
		KeyValue:     "sheet_permalink_type",
		Kind:         reflect.String,
	}
	CategoriesPrefix = Property{
		DefaultValue: "categories",
		KeyValue:     "categories_prefix",
		Kind:         reflect.String,
	}
	TagsPrefix = Property{
		DefaultValue: "tags",
		KeyValue:     "tags_prefix",
		Kind:         reflect.String,
	}
	ArchivesPrefix = Property{
		DefaultValue: "archives",
		KeyValue:     "archives_prefix",
		Kind:         reflect.String,
	}
	SheetPrefix = Property{
		DefaultValue: "s",
		KeyValue:     "sheet_prefix",
		Kind:         reflect.String,
	}
	LinksPrefix = Property{
		DefaultValue: "links",
		KeyValue:     "links_prefix",
		Kind:         reflect.String,
	}
	PhotosPrefix = Property{
		DefaultValue: "photos",
		KeyValue:     "photos_prefix",
		Kind:         reflect.String,
	}
	JournalsPrefix = Property{
		DefaultValue: "journals",
		KeyValue:     "journals_prefix",
		Kind:         reflect.String,
	}
	PathSuffix = Property{
		DefaultValue: "",
		KeyValue:     "path_suffix",
		Kind:         reflect.String,
	}
)
