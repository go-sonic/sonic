package consts

import (
	"database/sql/driver"
	"strconv"
	"strings"

	"github.com/go-sonic/sonic/util/xerr"
)

type DBType string

const (
	DBTypeMySQL  = "MySQL"
	DBTypeSQLite = "SQLite"
)

type AttachmentType int32

const (
	AttachmentTypeLocal AttachmentType = iota
	// AttachmentTypeUpOSS 又拍云
	AttachmentTypeUpOSS
	// AttachmentTypeQiNiuOSS 七牛云
	AttachmentTypeQiNiuOSS
	// AttachmentTypeSMMS sm.ms
	AttachmentTypeSMMS
	// AttachmentTypeAliOSS 阿里云OSS
	AttachmentTypeAliOSS
	// AttachmentTypeBaiDuOSS 百度云OSS
	AttachmentTypeBaiDuOSS
	// AttachmentTypeTencentCOS 腾讯COS
	AttachmentTypeTencentCOS
	// AttachmentTypeHuaweiOBS 华为OBS
	AttachmentTypeHuaweiOBS
	// AttachmentTypeMinIO AttachmentTypeMinIO
	AttachmentTypeMinIO
)

func (a AttachmentType) String() string {
	switch a {
	case AttachmentTypeLocal:
		return "LOCAL"
	case AttachmentTypeUpOSS:
		return "UPOSS"
	case AttachmentTypeQiNiuOSS:
		return "QINIUOSS"
	case AttachmentTypeSMMS:
		return "AttachmentTypeSMMS"
	case AttachmentTypeAliOSS:
		return "ALIOSS"
	case AttachmentTypeBaiDuOSS:
		return "BAIDUOSS"
	case AttachmentTypeTencentCOS:
		return "TENCENTOSS"
	case AttachmentTypeHuaweiOBS:
		return "HUAWEIOBS"
	case AttachmentTypeMinIO:
		return "MINIO"
	default:
		return "UNKNOWN"
	}
}

func (a AttachmentType) MarshalJSON() ([]byte, error) {
	switch a {
	case AttachmentTypeLocal:
		return []byte(`"LOCAL"`), nil
	case AttachmentTypeUpOSS:
		return []byte(`"UPOSS"`), nil
	case AttachmentTypeQiNiuOSS:
		return []byte(`"QINIUOSS"`), nil
	case AttachmentTypeSMMS:
		return []byte(`"AttachmentTypeSMMS"`), nil
	case AttachmentTypeAliOSS:
		return []byte(`"ALIOSS"`), nil
	case AttachmentTypeBaiDuOSS:
		return []byte(`"BAIDUOSS"`), nil
	case AttachmentTypeTencentCOS:
		return []byte(`"TENCENTOSS"`), nil
	case AttachmentTypeHuaweiOBS:
		return []byte(`"HUAWEIOBS"`), nil
	case AttachmentTypeMinIO:
		return []byte(`"MINIO"`), nil
	default:
		return []byte(`"UNKNOWN"`), nil
	}
}

func (a *AttachmentType) UnmarshalJSON(data []byte) error {
	str := string(data)
	switch str {
	case `"LOCAL"`:
		*a = AttachmentTypeLocal
	case `"UPOSS"`:
		*a = AttachmentTypeUpOSS
	case `"QINIUOSS"`:
		*a = AttachmentTypeQiNiuOSS
	case `"AttachmentTypeSMMS"`:
		*a = AttachmentTypeSMMS
	case `"ALIOSS"`:
		*a = AttachmentTypeAliOSS
	case `"BAIDUBOS"`:
		*a = AttachmentTypeBaiDuOSS
	case `"TENCENTCOS"`:
		*a = AttachmentTypeTencentCOS
	case `"HUAWEIOBS"`:
		*a = AttachmentTypeHuaweiOBS
	case `"MINIO"`:
		*a = AttachmentTypeMinIO
	default:
		return xerr.BadParam.New("").WithMsg("unknown AttachmentType")
	}
	return nil
}

func (a *AttachmentType) Scan(src interface{}) error {
	if src == nil {
		return xerr.BadParam.New("").WithMsg("field nil")
	}
	switch data := src.(type) {
	case int64:
		*a = AttachmentType(data)
	case int32:
		*a = AttachmentType(data)
	case int:
		*a = AttachmentType(data)
	default:
		return xerr.BadParam.New("").WithMsg("bad type")
	}
	return nil
}

func (a AttachmentType) Value() (driver.Value, error) {
	return int64(a), nil
}

type LogType int32

const (
	LogTypeBlogInitialized LogType = iota
	LogTypePostPublished
	LogTypePostEdited
	LogTypePostDeleted
	LogTypeLoggedIn
	LogTypeLoggedOut
	LogTypeLoginFailed
	LogTypePasswordUpdated
	LogTypeProfileUpdated
	LogTypeSheetPublished
	LogTypeSheetEdited
	LogTypeSheetDeleted
	LogTypeMfaUpdated
	LogTypeLoggedPreCheck
)

func (l LogType) MarshalJSON() ([]byte, error) {
	switch l {
	case LogTypeBlogInitialized:
		return []byte(`"BLOG_INITIALIZED"`), nil
	case LogTypePostPublished:
		return []byte(`"POST_PUBLISHED"`), nil
	case LogTypePostEdited:
		return []byte(`"POST_EDITED"`), nil
	case LogTypePostDeleted:
		return []byte(`"POST_DELETED"`), nil
	case LogTypeLoggedIn:
		return []byte(`"LOGGED_IN"`), nil
	case LogTypeLoggedOut:
		return []byte(`"LOGGED_OUT"`), nil
	case LogTypeLoginFailed:
		return []byte(`"LOGIN_FAILED"`), nil
	case LogTypePasswordUpdated:
		return []byte(`"PASSWORD_UPDATED"`), nil
	case LogTypeProfileUpdated:
		return []byte(`"PROFILE_UPDATED"`), nil
	case LogTypeSheetPublished:
		return []byte(`"SHEET_PUBLISHED"`), nil
	case LogTypeSheetEdited:
		return []byte(`"SHEET_EDITED"`), nil
	case LogTypeSheetDeleted:
		return []byte(`"SHEET_DELETED"`), nil
	case LogTypeMfaUpdated:
		return []byte(`"MFA_UPDATED"`), nil
	case LogTypeLoggedPreCheck:
		return []byte(`"LOGGED_PRE_CHECK"`), nil
	}
	return nil, nil
}

func (l *LogType) Scan(src interface{}) error {
	if src == nil {
		return xerr.BadParam.New("").WithMsg("field nil")
	}
	switch data := src.(type) {
	case int64:
		*l = LogType(data)
	case int32:
		*l = LogType(data)
	case int:
		*l = LogType(data)
	default:
		return xerr.BadParam.New("").WithMsg("bad type")
	}
	return nil
}

func (l LogType) Value() (driver.Value, error) {
	return int64(l), nil
}

type MFAType int32

const (
	MFANone MFAType = iota
	// MFATFATotp Time-based One-time Password (rfc6238).
	// see: https://tools.ietf.org/html/rfc6238
	MFATFATotp
)

func (m MFAType) MarshalJSON() ([]byte, error) {
	if m == MFANone {
		return []byte(`"NONE"`), nil
	} else if m == MFATFATotp {
		return []byte(`"TFA_TOTP"`), nil
	}
	return nil, nil
}

func (m *MFAType) UnmarshalJSON(data []byte) error {
	str := string(data)
	switch str {
	case `"NONE"`:
		*m = MFANone
	case `"TFA_TOTP"`:
		*m = MFATFATotp
	default:
		return xerr.BadParam.New("").WithMsg("unknown MFAType")
	}
	return nil
}

func (m *MFAType) Scan(src interface{}) error {
	if src == nil {
		return xerr.BadParam.New("").WithMsg("field nil")
	}
	switch data := src.(type) {
	case int64:
		*m = MFAType(data)
	case int32:
		*m = MFAType(data)
	case int:
		*m = MFAType(data)
	default:
		return xerr.BadParam.New("").WithMsg("bad type")
	}
	return nil
}

func (m MFAType) Value() (driver.Value, error) {
	return int64(m), nil
}

type PostStatus int32

const (
	PostStatusPublished PostStatus = iota
	PostStatusDraft
	PostStatusRecycle
	PostStatusIntimate
)

func (c PostStatus) MarshalJSON() ([]byte, error) {
	switch c {
	case PostStatusPublished:
		return []byte(`"PUBLISHED"`), nil
	case PostStatusDraft:
		return []byte(`"DRAFT"`), nil
	case PostStatusRecycle:
		return []byte(`"RECYCLE"`), nil
	case PostStatusIntimate:
		return []byte(`"INTIMATE"`), nil
	}
	return nil, nil
}

func (c *PostStatus) UnmarshalJSON(data []byte) error {
	str := string(data)
	switch str {
	case `"PUBLISHED"`:
		*c = PostStatusPublished
	case `"DRAFT"`:
		*c = PostStatusDraft
	case `"RECYCLE"`:
		*c = PostStatusRecycle
	case `"INTIMATE"`:
		*c = PostStatusIntimate
	case "":
		*c = PostStatusDraft
	default:
		return xerr.BadParam.New("").WithMsg("unknown PostStatus")
	}
	return nil
}

func (c *PostStatus) Scan(src interface{}) error {
	if src == nil {
		return xerr.BadParam.New("").WithMsg("field nil")
	}
	switch data := src.(type) {
	case int64:
		*c = PostStatus(data)
	case int32:
		*c = PostStatus(data)
	case int:
		*c = PostStatus(data)
	default:
		return xerr.BadParam.New("").WithMsg("bad type")
	}
	return nil
}

func (c PostStatus) Value() (driver.Value, error) {
	return int64(c), nil
}

func (c PostStatus) Ptr() *PostStatus {
	return &c
}

func PostStatusFromString(str string) (PostStatus, error) {
	switch str {
	case "PUBLISHED":
		return PostStatusPublished, nil
	case "DRAFT":
		return PostStatusDraft, nil
	case "RECYCLE":
		return PostStatusRecycle, nil
	case "INTIMATE":
		return PostStatusIntimate, nil
	default:
		return PostStatusDraft, xerr.BadParam.New("").WithMsg("unknown PostStatus")
	}
}

type CommentStatus int32

const (
	CommentStatusPublished CommentStatus = iota
	CommentStatusAuditing
	CommentStatusRecycle
)

func (c CommentStatus) MarshalJSON() ([]byte, error) {
	switch c {
	case CommentStatusPublished:
		return []byte(`"PUBLISHED"`), nil
	case CommentStatusAuditing:
		return []byte(`"AUDITING"`), nil
	case CommentStatusRecycle:
		return []byte(`"RECYCLE"`), nil
	}
	return nil, nil
}

func (c *CommentStatus) UnmarshalJSON(data []byte) error {
	str := string(data)

	switch str {
	case `"PUBLISHED"`:
		*c = CommentStatusPublished
	case `"AUDITING"`:
		*c = CommentStatusAuditing
	case `"RECYCLE"`:
		*c = CommentStatusRecycle
	default:
		return xerr.BadParam.New("").WithMsg("unknown CommentStatus")
	}
	return nil
}

func (c *CommentStatus) Scan(src interface{}) error {
	if src == nil {
		return xerr.BadParam.New("").WithMsg("field nil")
	}
	switch data := src.(type) {
	case int64:
		*c = CommentStatus(data)
	case int32:
		*c = CommentStatus(data)
	case int:
		*c = CommentStatus(data)
	default:
		return xerr.BadParam.New("").WithMsg("bad type")
	}
	return nil
}

func CommentStatusFromString(str string) (CommentStatus, error) {
	switch str {
	case "PUBLISHED":
		return CommentStatusPublished, nil
	case "AUDITING":
		return CommentStatusAuditing, nil
	case "RECYCLE":
		return CommentStatusRecycle, nil
	default:
		return CommentStatusPublished, xerr.BadParam.New("").WithMsg("unknown CommentStatus")
	}
}

func (c CommentStatus) Value() (driver.Value, error) {
	return int64(c), nil
}

func (c CommentStatus) Ptr() *CommentStatus {
	return &c
}

type PostPermalinkType string

const (
	PostPermalinkTypeDefault PostPermalinkType = "DEFAULT"
	PostPermalinkTypeDate    PostPermalinkType = "DATE"
	PostPermalinkTypeDay     PostPermalinkType = "DAY"
	PostPermalinkTypeID      PostPermalinkType = "ID"
	PostPermalinkTypeYear    PostPermalinkType = "YEAR"
	PostPermalinkTypeIDSlug  PostPermalinkType = "ID_SLUG"
)

type EditorType int32

const (
	EditorTypeMarkdown EditorType = iota
	EditorTypeRichText
)

func (e EditorType) Ptr() *EditorType {
	return &e
}

func (e *EditorType) Scan(src interface{}) error {
	if src == nil {
		return xerr.BadParam.New("").WithMsg("field nil")
	}
	switch data := src.(type) {
	case int64:
		*e = EditorType(data)
	case int32:
		*e = EditorType(data)
	case int:
		*e = EditorType(data)
	default:
		return xerr.BadParam.New("").WithMsg("bad type")
	}
	return nil
}

func (e EditorType) Value() (driver.Value, error) {
	return int64(e), nil
}

func (e EditorType) MarshalJSON() ([]byte, error) {
	if e == EditorTypeMarkdown {
		return []byte(`"MARKDOWN"`), nil
	} else if e == EditorTypeRichText {
		return []byte(`"RICHTEXT"`), nil
	}
	return nil, nil
}

func (e *EditorType) UnmarshalJSON(data []byte) error {
	str := string(data)
	switch str {
	case `"MARKDOWN"`:
		*e = EditorTypeMarkdown
	case `"RICHTEXT"`:
		*e = EditorTypeRichText
	case "":
		*e = EditorTypeMarkdown
	default:
		return xerr.BadParam.New("").WithMsg("unknown editorType")
	}
	return nil
}

type OptionType int32

const (
	OptionTypeInternal = iota
	OptionTypeCustom
)

func (o OptionType) MarshalJSON() ([]byte, error) {
	if o == OptionTypeInternal {
		return []byte(`"INTERNAL"`), nil
	} else if o == OptionTypeCustom {
		return []byte(`"CUSTOM"`), nil
	}
	return nil, nil
}

func (o *OptionType) UnmarshalJSON(data []byte) error {
	switch string(data) {
	case `"INTERNAL"`:
		*o = OptionTypeInternal
	case `"CUSTOM"`:
		*o = OptionTypeCustom
	default:
		return xerr.BadParam.New("").WithMsg("unknown OptionType")
	}
	return nil
}

func (o *OptionType) Scan(src interface{}) error {
	if src == nil {
		return xerr.BadParam.New("").WithMsg("field nil")
	}
	switch data := src.(type) {
	case int64:
		*o = OptionType(data)
	case int32:
		*o = OptionType(data)
	case int:
		*o = OptionType(data)
	default:
		return xerr.BadParam.New("").WithMsg("bad type")
	}
	return nil
}

func (o OptionType) Value() (driver.Value, error) {
	return int64(o), nil
}

type PostType int32

const (
	PostTypePost PostType = iota
	PostTypeSheet
)

func (p *PostType) Scan(src interface{}) error {
	if src == nil {
		return xerr.BadParam.New("").WithMsg("field nil")
	}
	switch data := src.(type) {
	case int64:
		*p = PostType(data)
	case int32:
		*p = PostType(data)
	case int:
		*p = PostType(data)
	default:
		return xerr.BadParam.New("").WithMsg("bad type")
	}
	return nil
}

func (p PostType) Value() (driver.Value, error) {
	return int64(p), nil
}

type CommentType int32

const (
	CommentTypePost CommentType = iota
	CommentTypeSheet
	CommentTypeJournal
)

func (ct *CommentType) Scan(src interface{}) error {
	if src == nil {
		return xerr.BadParam.New("").WithMsg("unknown OptionType")
	}
	switch data := src.(type) {
	case int64:
		*ct = CommentType(data)
	case int32:
		*ct = CommentType(data)
	case int:
		*ct = CommentType(data)
	default:
		return xerr.BadParam.New("").WithMsg("bad type")
	}
	return nil
}

func (ct CommentType) Value() (driver.Value, error) {
	return int64(ct), nil
}

type SheetPermaLinkType string

const (
	SheetPermaLinkTypeSecondary = "SECONDARY"
	SheetPermaLinkTypeRoot      = "ROOT"
)

type JournalType int32

const (
	JournalTypePublic JournalType = iota
	JournalTypeIntimate
)

func (j JournalType) Ptr() *JournalType {
	return &j
}

func (j JournalType) MarshalJSON() ([]byte, error) {
	if j == JournalTypePublic {
		return []byte(`"PUBLIC"`), nil
	} else if j == JournalTypeIntimate {
		return []byte(`"INTIMATE"`), nil
	}
	return nil, nil
}

func (j *JournalType) UnmarshalJSON(data []byte) error {
	switch string(data) {
	case `"PUBLIC"`:
		*j = JournalTypePublic
	case `"INTIMATE"`:
		*j = JournalTypeIntimate
	default:
		return xerr.BadParam.New("").WithMsg("unknown JournalType")
	}
	return nil
}

func (j *JournalType) Scan(src interface{}) error {
	if src == nil {
		return xerr.BadParam.New("").WithMsg("field nil")
	}
	switch data := src.(type) {
	case int64:
		*j = JournalType(data)
	case int32:
		*j = JournalType(data)
	case int:
		*j = JournalType(data)
	default:
		return xerr.BadParam.New("").WithMsg("bad type")
	}
	return nil
}

func (j JournalType) Value() (driver.Value, error) {
	return int64(j), nil
}

type MetaType int32

const (
	MetaTypePost  = iota
	MetaTypeSheet = iota
)

func (m *MetaType) Scan(src interface{}) error {
	if src == nil {
		return xerr.BadParam.New("").WithMsg("field nil")
	}
	switch data := src.(type) {
	case int64:
		*m = MetaType(data)
	case int32:
		*m = MetaType(data)
	case int:
		*m = MetaType(data)
	default:
		return xerr.BadParam.New("").WithMsg("bad type")
	}
	return nil
}

func (m MetaType) Value() (driver.Value, error) {
	return int64(m), nil
}

type ThemeUpdateStrategy int32

const (
	ThemeUpdateStrategyBranch = iota
	ThemeUpdateStrategyRelease
)

type ThemeConfigInputType int32

const (
	// ThemeConfigInputTypeTEXT Text input type
	ThemeConfigInputTypeTEXT = iota

	// ThemeConfigInputTypeNUMBER Number input type
	ThemeConfigInputTypeNUMBER

	// ThemeConfigInputTypeRADIO Radio box input type
	ThemeConfigInputTypeRADIO

	// ThemeConfigInputTypeSELECT Select input type
	ThemeConfigInputTypeSELECT

	// ThemeConfigInputTypeTEXTAREA Textarea input type
	ThemeConfigInputTypeTEXTAREA

	// ThemeConfigInputTypeCOLOR Color picker input type
	ThemeConfigInputTypeCOLOR

	// ThemeConfigInputTypeATTACHMENT Attachment picker input type
	ThemeConfigInputTypeATTACHMENT

	// ThemeConfigInputTypeSWITCH Switch input type, only true or false
	ThemeConfigInputTypeSWITCH
)

func (t ThemeConfigInputType) MarshalJSON() ([]byte, error) {
	switch t {
	case ThemeConfigInputTypeTEXT:
		return []byte(`"TEXT"`), nil
	case ThemeConfigInputTypeNUMBER:
		return []byte(`"NUMBER"`), nil
	case ThemeConfigInputTypeRADIO:
		return []byte(`"RADIO"`), nil
	case ThemeConfigInputTypeSELECT:
		return []byte(`"SELECT"`), nil
	case ThemeConfigInputTypeTEXTAREA:
		return []byte(`"TEXTAREA"`), nil
	case ThemeConfigInputTypeCOLOR:
		return []byte(`"COLOR"`), nil
	case ThemeConfigInputTypeATTACHMENT:
		return []byte(`"ATTACHMENT"`), nil
	case ThemeConfigInputTypeSWITCH:
		return []byte(`"SWITCH"`), nil
	default:
		return nil, xerr.BadParam.New("").WithMsg("unknown ThemeConfigInputType")
	}
}

func (t *ThemeConfigInputType) UnmarshalJSON(data []byte) error {
	str := string(data)
	switch str {
	case `"TEXT"`:
		*t = ThemeConfigInputTypeTEXT
		return nil
	case `"NUMBER"`:
		*t = ThemeConfigInputTypeNUMBER
		return nil
	case `"RADIO"`:
		*t = ThemeConfigInputTypeRADIO
		return nil
	case `"SELECT"`:
		*t = ThemeConfigInputTypeSELECT
		return nil
	case `"TEXTAREA"`:
		*t = ThemeConfigInputTypeTEXTAREA
		return nil
	case `"COLOR"`:
		*t = ThemeConfigInputTypeCOLOR
		return nil
	case `"SWITCH"`:
		*t = ThemeConfigInputTypeSWITCH
		return nil
	case `"ATTACHMENT"`:
		*t = ThemeConfigInputTypeATTACHMENT
		return nil
	default:
		return xerr.BadParam.New("").WithMsg("unknown ThemeConfigInputType")
	}
}

func (t *ThemeConfigInputType) UnmarshalYAML(unmarshal func(interface{}) error) error {
	strType := ""
	err := unmarshal(&strType)
	if err != nil {
		return xerr.BadParam.New("").WithMsg("ThemeConfigInputType yaml unmarshal err")
	}
	strType = strings.ToUpper(strType)
	switch strType {
	case "TEXT":
		*t = ThemeConfigInputTypeTEXT
		return nil
	case "NUMBER":
		*t = ThemeConfigInputTypeNUMBER
		return nil
	case "RADIO":
		*t = ThemeConfigInputTypeRADIO
		return nil
	case "SELECT":
		*t = ThemeConfigInputTypeSELECT
		return nil
	case "TEXTAREA":
		*t = ThemeConfigInputTypeTEXTAREA
		return nil
	case "COLOR":
		*t = ThemeConfigInputTypeCOLOR
		return nil
	case "SWITCH":
		*t = ThemeConfigInputTypeSWITCH
		return nil
	case "ATTACHMENT":
		*t = ThemeConfigInputTypeATTACHMENT
		return nil
	default:
		return xerr.BadParam.New("").WithMsg("unknown ThemeConfigInputType")
	}
}

type ThemeConfigDataType int32

const (
	ThemeConfigDataTypeString ThemeConfigDataType = iota
	ThemeConfigDataTypeLong
	ThemeConfigDataTypeDouble
	ThemeConfigDataTypeBool
)

func (t ThemeConfigDataType) Convert(value string) (interface{}, error) {
	switch t {
	case ThemeConfigDataTypeString:
		return value, nil
	case ThemeConfigDataTypeLong:
		result, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return nil, xerr.WithErrMsgf(err, "value invalid ThemeConfigDataType")
		}
		return result, nil
	case ThemeConfigDataTypeDouble:
		result, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return nil, xerr.WithErrMsgf(err, "value invalid ThemeConfigDataType")
		}
		return result, nil
	case ThemeConfigDataTypeBool:
		result, err := strconv.ParseBool(value)
		if err != nil {
			return nil, xerr.WithErrMsgf(err, "value invalid ThemeConfigDataType")
		}
		return result, nil
	default:
		return nil, xerr.WithErrMsgf(nil, "invalid ThemeConfigDataType")
	}
}

func (t ThemeConfigDataType) FormatToStr(value interface{}) (string, error) {
	switch t {
	case ThemeConfigDataTypeString:
		valueStr, ok := value.(string)
		if !ok {
			return "", xerr.WithErrMsgf(nil, "value invalid ThemeConfigDataType")
		}
		return valueStr, nil
	case ThemeConfigDataTypeLong:
		var valueStr string
		switch data := value.(type) {
		case int:
			valueStr = strconv.FormatInt(int64(data), 10)
		case int64:
			valueStr = strconv.FormatInt(data, 10)
		case int32:
			valueStr = strconv.FormatInt(int64(data), 10)
		default:
			return "", xerr.WithErrMsgf(nil, "value invalid ThemeConfigDataType")
		}
		return valueStr, nil
	case ThemeConfigDataTypeDouble:
		var valueStr string
		switch data := value.(type) {
		case float32:
			valueStr = strconv.FormatFloat(float64(data), 'f', 5, 32)
		case float64:
			valueStr = strconv.FormatFloat(data, 'f', 5, 64)
		default:
			return "", xerr.WithErrMsgf(nil, "value invalid ThemeConfigDataType")
		}
		return valueStr, nil
	case ThemeConfigDataTypeBool:
		var valueStr string
		switch data := value.(type) {
		case bool:
			valueStr = strconv.FormatBool(data)
		default:
			return "", xerr.WithErrMsgf(nil, "value invalid ThemeConfigDataType")
		}
		return valueStr, nil
	default:
		return "", xerr.WithErrMsgf(nil, "invalid ThemeConfigDataType")
	}
}

func (t ThemeConfigDataType) MarshalJSON() ([]byte, error) {
	switch t {
	case ThemeConfigDataTypeString:
		return []byte(`"STRING"`), nil
	case ThemeConfigDataTypeLong:
		return []byte(`"LONG"`), nil
	case ThemeConfigDataTypeDouble:
		return []byte(`"DOUBLE"`), nil
	case ThemeConfigDataTypeBool:
		return []byte(`"BOOL"`), nil
	default:
		return nil, xerr.BadParam.New("").WithMsg("unknown ThemeConfigDataType")
	}
}

func (t *ThemeConfigDataType) UnmarshalJSON(data []byte) error {
	str := string(data)
	switch str {
	case `"STRING"`:
		*t = ThemeConfigDataTypeString
		return nil
	case `"LONG"`:
		*t = ThemeConfigDataTypeLong
		return nil
	case `"DOUBLE"`:
		*t = ThemeConfigDataTypeDouble
		return nil
	case `"BOOL"`:
		*t = ThemeConfigDataTypeBool
		return nil
	default:
		return xerr.BadParam.New("").WithMsg("unknown ThemeConfigInputType")
	}
}

func (t *ThemeConfigDataType) UnmarshalYAML(unmarshal func(interface{}) error) error {
	strType := ""
	err := unmarshal(&strType)
	if err != nil {
		return xerr.BadParam.New("").WithMsg("ThemeConfigDataType yaml unmarshal err")
	}
	strType = strings.ToUpper(strType)
	switch strType {
	case "STRING":
		*t = ThemeConfigDataTypeString
		return nil
	case "LONG":
		*t = ThemeConfigDataTypeLong
		return nil
	case "DOUBLE":
		*t = ThemeConfigDataTypeDouble
		return nil
	case "BOOL":
		*t = ThemeConfigDataTypeBool
		return nil
	default:
		return xerr.BadParam.New("").WithMsg("unknown ThemeConfigDataType")
	}
}

type EncryptType int32

const (
	EncryptTypePost EncryptType = iota
	EncryptTypeCategory
)

func (e EncryptType) Name() string {
	if e == EncryptTypePost {
		return "post"
	}
	if e == EncryptTypeCategory {
		return "category"
	}
	return ""
}

type CategoryType int32

const (
	CategoryTypeNormal CategoryType = iota
	CategoryTypeIntimate
)

func (c CategoryType) MarshalJSON() ([]byte, error) {
	if c == CategoryTypeNormal {
		return []byte(`"NORMAL"`), nil
	} else if c == CategoryTypeIntimate {
		return []byte(`"INTIMATE"`), nil
	}
	return nil, nil
}

func (c *CategoryType) UnmarshalJSON(data []byte) error {
	switch string(data) {
	case `"NORMAL"`:
		*c = CategoryTypeNormal
	case `"INTIMATE"`:
		*c = CategoryTypeIntimate
	default:
		return xerr.BadParam.New("").WithMsg("unknown PostStatus")
	}
	return nil
}

func (c *CategoryType) Scan(src interface{}) error {
	if src == nil {
		return xerr.BadParam.New("").WithMsg("field nil")
	}
	switch data := src.(type) {
	case int64:
		*c = CategoryType(data)
	case int32:
		*c = CategoryType(data)
	case int:
		*c = CategoryType(data)
	default:
		return xerr.BadParam.New("").WithMsg("bad type")
	}
	return nil
}

func (c CategoryType) Value() (driver.Value, error) {
	return int64(c), nil
}

func (c CategoryType) Ptr() *CategoryType {
	return &c
}
