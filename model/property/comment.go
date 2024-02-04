package property

import "reflect"

var (
	CommentGravatarDefault = Property{
		KeyValue:     "comment_gravatar_default",
		DefaultValue: "identicon",
		Kind:         reflect.String,
	}
	CommentNewNeedCheck = Property{
		KeyValue:     "comment_new_need_check",
		DefaultValue: true,
		Kind:         reflect.Bool,
	}
	CommentNewNotice = Property{
		KeyValue:     "comment_new_notice",
		DefaultValue: true,
		Kind:         reflect.Bool,
	}
	CommentReplyNotice = Property{
		KeyValue:     "comment_reply_notice",
		DefaultValue: false,
		Kind:         reflect.Bool,
	}
	CommentAPIEnabled = Property{
		KeyValue:     "comment_api_enabled",
		DefaultValue: true,
		Kind:         reflect.Bool,
	}
	CommentPageSize = Property{
		KeyValue:     "comment_page_size",
		DefaultValue: 10,
		Kind:         reflect.Int,
	}
	CommentContentPlaceholder = Property{
		KeyValue:     "comment_content_placeholder",
		DefaultValue: "",
		Kind:         reflect.String,
	}
	CommentInternalPluginJs = Property{
		KeyValue:     "comment_internal_plugin_js",
		DefaultValue: "https://cdn.jsdelivr.net/npm/halo-comment@latest/dist/halo-comment.min.js",
		Kind:         reflect.String,
	}
	CommentGravatarSource = Property{
		KeyValue:     "gravatar_source",
		DefaultValue: "https://gravatar.com/avatar/",
		Kind:         reflect.String,
	}
	CommentBanTime = Property{
		KeyValue:     "comment_ban_time",
		DefaultValue: 10,
		Kind:         reflect.Int,
	}
	CommentRange = Property{
		KeyValue:     "comment_range",
		DefaultValue: 30,
		Kind:         reflect.Int,
	}
)
