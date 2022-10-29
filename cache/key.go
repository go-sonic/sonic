package cache

import (
	"context"
	"strconv"

	"github.com/go-sonic/sonic/consts"
	"github.com/go-sonic/sonic/util/xerr"
)

func BuildTokenAccessKey(accessToken string) string {
	return consts.TokenAccessCachePrefix + accessToken
}

func BuildTokenRefreshKey(refreshToken string) string {
	return consts.TokenRefreshCachePrefix + refreshToken
}

func BuildAccessTokenKey(userID int32) string {
	return consts.TokenAccessCachePrefix + strconv.Itoa(int(userID))
}

func BuildRefreshTokenKey(userID int32) string {
	return consts.TokenRefreshCachePrefix + strconv.Itoa(int(userID))
}

func BuildCodeCacheKey(userID int32) string {
	return consts.CodePrefix + strconv.Itoa(int(userID))
}

func BuildAccessPermissionKey(ctx context.Context) (string, error) {
	sessionID := ctx.Value(consts.SessionID)
	if sessionID == nil {
		return "", xerr.NoType.New("session_id not exist").WithStatus(xerr.StatusInternalServerError)
	}
	sessionIDStr, ok := sessionID.(string)
	if !ok || sessionIDStr == "" {
		return "", xerr.NoType.New("session_id not exist").WithStatus(xerr.StatusInternalServerError)
	}
	return consts.AccessPermissionKeyPrefix + sessionIDStr, nil
}

func BuildCategoryPermissionKey(categoryID int32) string {
	return strconv.Itoa(int(categoryID))
}
