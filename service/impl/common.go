package impl

import (
	"errors"
	"reflect"
	"strings"
	"unicode"

	"gorm.io/gen/field"
	"gorm.io/gorm"

	"github.com/go-sonic/sonic/model/param"
	"github.com/go-sonic/sonic/util/xerr"
)

func WrapDBErr(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return xerr.NoRecord.Wrap(err).WithMsg("The resource does not exist or has been deleted").WithStatus(xerr.StatusNotFound)
	}
	return xerr.DB.Wrap(err).WithStatus(xerr.StatusInternalServerError)
}

type Order struct {
	Property string
	Asc      bool
}

func BuildSort(sort *param.Sort, dalStruct interface{}, dalDO interface{}) error {
	if sort == nil || len(sort.Fields) == 0 {
		return nil
	}

	type iDAL interface {
		GetFieldByName(fieldName string) (field.OrderExpr, bool)
	}
	type iField interface {
		Desc() field.Expr
	}

	dal, ok := dalStruct.(iDAL)
	if !ok {
		panic("no GetFieldByName method")
	}

	rDO := reflect.ValueOf(dalDO)
	if rDO.Kind() != reflect.Ptr {
		panic("not pointer type")
	}
	do := rDO.Elem()

	orderMethod := do.MethodByName("Order")
	if orderMethod.IsZero() {
		panic("no order method")
	}

	orders, err := ConvertSort(sort)
	if err != nil {
		return err
	}

	for _, order := range orders {
		expr, ok := dal.GetFieldByName(order.Property)
		if !ok {
			return xerr.WithMsg(nil, "sort parameter error").WithStatus(xerr.StatusBadRequest)
		}
		var params []reflect.Value
		if order.Asc {
			params = append(params, reflect.ValueOf(expr))
			result := orderMethod.Call(params)[0]
			do.Set(result)
		} else {
			params = append(params, reflect.ValueOf(expr.(iField).Desc()))
			result := orderMethod.Call(params)[0]
			do.Set(result)
		}
	}
	return nil
}

func ConvertSort(sorts *param.Sort) ([]*Order, error) {
	if sorts == nil {
		return nil, nil
	}
	result := make([]*Order, 0, len(sorts.Fields))
	for _, sort := range sorts.Fields {
		items := strings.Split(sort, ",")
		if len(items) > 2 {
			return nil, xerr.WithMsg(nil, "").WithStatus(xerr.StatusBadRequest)
		}
		order := Order{}
		if len(items) == 1 {
			order.Property = UnderscoreName(items[0])
			// default asc
			order.Asc = true
			result = append(result, &order)
			continue
		}
		if len(items) == 2 {
			order.Property = UnderscoreName(items[0])
			items[1] = strings.ToLower(items[1])

			switch items[1] {
			case "asc":
				order.Asc = true
			case "desc":
				order.Asc = false
			default:
				return nil, xerr.WithStatus(nil, xerr.StatusBadRequest).WithMsg("sort parameter error")
			}

			result = append(result, &order)
		}
	}
	return result, nil
}

// UnderscoreName 驼峰式写法转为下划线写法
func UnderscoreName(name string) string {
	buffer := strings.Builder{}
	for i, r := range name {
		if unicode.IsUpper(r) {
			if i != 0 {
				buffer.WriteByte('_')
			}
			buffer.WriteRune(unicode.ToLower(r))
		} else {
			buffer.WriteRune(r)
		}
	}
	return buffer.String()
}
