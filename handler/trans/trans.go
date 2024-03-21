package trans

import (
	"fmt"
	"github.com/hertz-contrib/binding/go_playground"
	"strings"

	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	chTranslations "github.com/go-playground/validator/v10/translations/zh"
)

var trans ut.Translator

func init() {
	local := "zh"
	vd := go_playground.NewValidator()
	if v, ok := vd.Engine().(*validator.Validate); ok {
		zhT := zh.New() // chinese
		enT := en.New() // english
		uni := ut.New(enT, zhT, enT)

		var o bool
		trans, o = uni.GetTranslator(local)
		if !o {
			panic(fmt.Sprintf("uni.GetTranslator(%s) failed", local))
		}

		err := chTranslations.RegisterDefaultTranslations(v, trans)
		if err != nil {
			panic(err)
		}

		return
	}
}

func Translate(errs validator.ValidationErrors) string {
	errList := make([]string, 0, len(errs))
	for _, e := range errs {
		// can translate each error one at a time.
		errList = append(errList, e.Translate(trans))
	}
	return strings.Join(errList, "|")
}
