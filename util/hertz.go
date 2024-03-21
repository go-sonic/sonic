package util

import (
	"context"
	"fmt"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"

	"github.com/go-sonic/sonic/util/xerr"
)

func GetClientIP(ctx context.Context) string {
	value, ok := ctx.Value("clientIP").(string)
	if !ok {
		return ""
	}
	return value
}

func GetUserAgent(ctx context.Context) string {
	value, ok := ctx.Value("userAgent").(string)
	if !ok {
		return ""
	}
	return value
}

func MustGetQueryString(_ctx context.Context, ctx *app.RequestContext, key string) (string, error) {
	str, ok := ctx.GetQuery(key)
	if !ok || str == "" {
		return "", xerr.WithStatus(nil, xerr.StatusBadRequest).WithMsg(fmt.Sprintf("%s parameter does not exisit", key))
	}
	return str, nil
}

func MustGetQueryInt32(_ctx context.Context, ctx *app.RequestContext, key string) (int32, error) {
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

func MustGetQueryInt64(_ctx context.Context, ctx *app.RequestContext, key string) (int64, error) {
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

func MustGetQueryInt(_ctx context.Context, ctx *app.RequestContext, key string) (int, error) {
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

func MustGetQueryBool(_ctx context.Context, ctx *app.RequestContext, key string) (bool, error) {
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

func GetQueryBool(_ctx context.Context, ctx *app.RequestContext, key string, defaultValue bool) (bool, error) {
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

func ParamString(_ctx context.Context, ctx *app.RequestContext, key string) (string, error) {
	str := ctx.Param(key)
	if str == "" {
		return "", xerr.WithStatus(nil, xerr.StatusBadRequest).WithMsg(fmt.Sprintf("%s parameter does not exisit", key))
	}
	return str, nil
}

func ParamInt32(_ctx context.Context, ctx *app.RequestContext, key string) (int32, error) {
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

func ParamInt64(_ctx context.Context, ctx *app.RequestContext, key string) (int64, error) {
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

func ParamBool(_ctx context.Context, ctx *app.RequestContext, key string) (bool, error) {
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
