package binding

import (
	"net/http"

	"github.com/gin-gonic/gin/binding"
)

const defaultMemory = 32 << 20

// CustomFormBinding  If the type implements the UnmarshalJSON interface, use JSON to bind
// For the purpose of support enum string to turn the enum type binding
var CustomFormBinding = customFormBinding{}
var CustomFormPostBinding = customFormPostBinding{}

type (
	customFormBinding     struct{}
	customFormPostBinding struct{}
)

func (customFormBinding) Name() string {
	return "form"
}

func (customFormBinding) Bind(req *http.Request, obj interface{}) error {
	if err := req.ParseForm(); err != nil {
		return err
	}
	if err := req.ParseMultipartForm(defaultMemory); err != nil {
		if err != http.ErrNotMultipart {
			return err
		}
	}
	if err := mapForm(obj, req.Form); err != nil {
		return err
	}
	return validate(obj)
}

func (customFormPostBinding) Name() string {
	return "form-urlencoded"
}

func (customFormPostBinding) Bind(req *http.Request, obj interface{}) error {
	if err := req.ParseForm(); err != nil {
		return err
	}
	if err := mapForm(obj, req.PostForm); err != nil {
		return err
	}
	return validate(obj)
}

func validate(obj interface{}) error {
	if binding.Validator == nil {
		return nil
	}
	return binding.Validator.ValidateStruct(obj)
}
