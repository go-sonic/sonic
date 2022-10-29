package consts

import (
	"time"
)

const (
	AccessTokenExpiredSeconds = 24 * 3600
	RefreshTokenExpiredDays   = 30
	TokenAccessCachePrefix    = "admin_access_token_"
	TokenRefreshCachePrefix   = "admin_refresh_token_"
	AdminTokenHeaderName      = "Admin-Authorization"
	AuthorizedUser            = "authorized_user"
	CodePrefix                = "code_"
	CodeValidDuration         = time.Second
	OneTimeTokenQueryName     = "ott"
	SessionID                 = "session_id"
	AccessPermissionKeyPrefix = "access_permission_"
)

const (
	SonicVersion              = "1.0.0"
	SonicBackupPrefix         = "sonic-backup-"
	SonicDataExportPrefix     = "sonic-data-export-"
	SonicBackupMarkdownPrefix = "sonic-backup-markdown-"
	SonicDefaultTagColor      = "#cfd3d7"
	SonicUploadDir            = "upload"
	SonicDefaultThemeDirName  = "default-theme-anatole"
)

var (
	ThemePropertyFilenames = [2]string{"theme.yaml", "theme.yml"}
	ThemeSettingFilenames  = [2]string{"settings.yaml", "settings.yml"}
)

const (
	DefaultThemeId         = "caicai_anatole"
	ThemeScreenshotsName   = "screenshot"
	ThemeCustomSheetPrefix = "sheet_"
	ThemeCustomPostPrefix  = "post_"
)

// StartTime 系统启动时间
var StartTime time.Time

var DatabaseVersion string
