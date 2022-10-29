package trans

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	chTranslations "github.com/go-playground/validator/v10/translations/zh"
)

var trans ut.Translator

func init() {
	local := "zh"
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
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
	var errList []string
	for _, e := range errs {
		// can translate each error one at a time.
		errList = append(errList, e.Translate(trans))
	}
	return strings.Join(errList, "|")
}
