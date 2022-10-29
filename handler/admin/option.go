package admin

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/go-sonic/sonic/model/param"
	"github.com/go-sonic/sonic/service"
	"github.com/go-sonic/sonic/util/xerr"
)

type OptionHandler struct {
	OptionService service.OptionService
}

func NewOptionHandler(optionService service.OptionService) *OptionHandler {
	return &OptionHandler{
		OptionService: optionService,
	}
}

func (o *OptionHandler) ListAllOptions(ctx *gin.Context) (interface{}, error) {
	return o.OptionService.ListAllOption(ctx)
}

func (o *OptionHandler) SaveOption(ctx *gin.Context) (interface{}, error) {
	optionParams := make([]*param.Option, 0)
	err := ctx.ShouldBindJSON(&optionParams)
	if err != nil {
		return nil, xerr.WithMsg(err, "param error").WithStatus(xerr.StatusBadRequest)
	}
	optionMap := make(map[string]string, 0)
	for _, option := range optionParams {
		optionMap[option.Key] = option.Value
	}
	return nil, o.OptionService.Save(ctx, optionMap)
}

func (o *OptionHandler) ListAllOptionsAsMap(ctx *gin.Context) (interface{}, error) {
	options, err := o.OptionService.ListAllOption(ctx)
	if err != nil {
		return nil, err
	}
	result := make(map[string]interface{})
	for _, option := range options {
		result[option.Key] = option.Value
	}
	return result, nil
}

func (o *OptionHandler) ListAllOptionsAsMapWithKey(ctx *gin.Context) (interface{}, error) {
	keys := make([]string, 0)
	err := ctx.ShouldBindJSON(&keys)
	if err != nil {
		return nil, xerr.WithMsg(err, "option key error").WithStatus(xerr.StatusBadRequest)
	}
	options, err := o.OptionService.ListAllOption(ctx)
	if err != nil {
		return nil, err
	}
	keyMap := make(map[string]struct{})
	for _, key := range keys {
		keyMap[key] = struct{}{}
	}
	result := make(map[string]interface{})
	for _, option := range options {
		if _, ok := keyMap[option.Key]; ok {
			result[option.Key] = option.Value
		}
	}
	return result, nil
}

func (o *OptionHandler) SaveOptionWithMap(ctx *gin.Context) (interface{}, error) {
	optionMap := make(map[string]interface{}, 0)
	err := ctx.ShouldBind(&optionMap)
	if err != nil {
		return nil, xerr.WithMsg(err, "parameter error").WithStatus(xerr.StatusBadRequest)
	}
	temp := make(map[string]string)
	for key, value := range optionMap {
		var v string
		switch value := value.(type) {
		case int32:
			v = strconv.Itoa(int(value))
		case int64:
			v = strconv.FormatInt(value, 10)
		case int:
			v = strconv.Itoa(value)
		case string:
			v = value
		case bool:
			v = strconv.FormatBool(value)
		case float64:
			v = strconv.FormatFloat(value, 'f', -1, 64)
		case float32:
			v = strconv.FormatFloat(float64(value), 'f', -1, 32)
		default:
			return nil, xerr.BadParam.New("key=%v,value=%v", key, value).WithStatus(xerr.StatusBadRequest).WithMsg("Parameter type is incorrect")
		}
		temp[key] = v
	}
	return nil, o.OptionService.Save(ctx, temp)
}
