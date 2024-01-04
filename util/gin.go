package util

import (
	"context"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/go-sonic/sonic/util/xerr"
)

func GetClientIP(ctx context.Context) string {
	ginCtx, ok := ctx.(*gin.Context)
	if !ok {
		return ""
	}
	return ginCtx.ClientIP()
}

func GetUserAgent(ctx context.Context) string {
	ginCtx, ok := ctx.(*gin.Context)
	if !ok {
		return ""
	}
	return ginCtx.GetHeader("User-Agent")
}

func MustGetQueryString(ctx *gin.Context, key string) (string, error) {
	str, ok := ctx.GetQuery(key)
	if !ok || str == "" {
		return "", xerr.WithStatus(nil, xerr.StatusBadRequest).WithMsg(fmt.Sprintf("%s parameter does not exisit", key))
	}
	return str, nil
}

func MustGetQueryInt32(ctx *gin.Context, key string) (int32, error) {
	str, ok := ctx.GetQuery(key)
	if !ok {
		return 0, xerr.WithStatus(nil, xerr.StatusBadRequest).WithMsg(fmt.Sprintf("%s parameter does not exisit", key))
	}
	value, err := strconv.ParseInt(str, 10, 32)
	if err != nil {
		return 0, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg(fmt.Sprintf("The parameter %s type is incorrect", key))
	}
	return int32(value), nil
}

func MustGetQueryInt64(ctx *gin.Context, key string) (int64, error) {
	str, ok := ctx.GetQuery(key)
	if !ok {
		return 0, xerr.WithStatus(nil, xerr.StatusBadRequest).WithMsg(fmt.Sprintf("%s parameter does not exisit", key))
	}
	value, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg(fmt.Sprintf("The parameter %s type is incorrect", key))
	}
	return value, nil
}

func MustGetQueryInt(ctx *gin.Context, key string) (int, error) {
	str, ok := ctx.GetQuery(key)
	if !ok {
		return 0, xerr.WithStatus(nil, xerr.StatusBadRequest).WithMsg(fmt.Sprintf("%s parameter does not exisit", key))
	}
	value, err := strconv.Atoi(str)
	if err != nil {
		return 0, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg(fmt.Sprintf("The parameter %s type is incorrect", key))
	}
	return value, nil
}

func MustGetQueryBool(ctx *gin.Context, key string) (bool, error) {
	str, ok := ctx.GetQuery(key)
	if !ok {
		return false, xerr.WithStatus(nil, xerr.StatusBadRequest).WithMsg(fmt.Sprintf("%s parameter does not exisit", key))
	}
	value, err := strconv.ParseBool(str)
	if err != nil {
		return false, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg(fmt.Sprintf("The parameter %s type is incorrect", key))
	}
	return value, nil
}

func GetQueryBool(ctx *gin.Context, key string, defaultValue bool) (bool, error) {
	str, ok := ctx.GetQuery(key)
	if !ok {
		return defaultValue, nil
	}
	value, err := strconv.ParseBool(str)
	if err != nil {
		return false, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg(fmt.Sprintf("The parameter %s type is incorrect", key))
	}
	return value, nil
}

func GetQueryInt32(ctx *gin.Context, key string, defaultValue int32) (int32, error) {
	str, ok := ctx.GetQuery(key)
	if !ok {
		return defaultValue, nil
	}
	value, err := strconv.ParseInt(str, 10, 32)
	if err != nil {
		return 0, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg(fmt.Sprintf("The parameter %s type is incorrect", key))
	}
	return int32(value), nil
}

func ParamString(ctx *gin.Context, key string) (string, error) {
	str := ctx.Param(key)
	if str == "" {
		return "", xerr.WithStatus(nil, xerr.StatusBadRequest).WithMsg(fmt.Sprintf("%s parameter does not exisit", key))
	}
	return str, nil
}

func ParamInt32(ctx *gin.Context, key string) (int32, error) {
	str := ctx.Param(key)
	if str == "" {
		return 0, xerr.WithStatus(nil, xerr.StatusBadRequest).WithMsg(fmt.Sprintf("%s parameter does not exisit", key))
	}
	value, err := strconv.ParseInt(str, 10, 32)
	if err != nil {
		return 0, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg(fmt.Sprintf("The parameter %s type is incorrect", key))
	}
	return int32(value), nil
}

func ParamInt64(ctx *gin.Context, key string) (int64, error) {
	str := ctx.Param(key)
	if str == "" {
		return 0, xerr.WithStatus(nil, xerr.StatusBadRequest).WithMsg(fmt.Sprintf("%s parameter does not exisit", key))
	}
	value, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg(fmt.Sprintf("The parameter %s type is incorrect", key))
	}
	return value, nil
}

func ParamBool(ctx *gin.Context, key string) (bool, error) {
	str := ctx.Param(key)
	if str == "" {
		return false, xerr.WithStatus(nil, xerr.StatusBadRequest).WithMsg(fmt.Sprintf("%s parameter does not exisit", key))
	}
	value, err := strconv.ParseBool(str)
	if err != nil {
		return false, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg(fmt.Sprintf("The parameter %s type is incorrect", key))
	}
	return value, nil
}
