package validator

import (
	"reflect"
	"strings"

	v "github.com/go-playground/validator/v10"

	"github.com/falentio/skul/internal/pkg/response"
)

var validate = v.New()

func init() {
	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
		if name == "" {
			name = strings.ToLower(field.Name)
		}
		if name == "-" {
			return ""
		}
		return name
	})
}

func Struct(s any) (err error) {
	err = validate.Struct(s)
	if ferr, ok := err.(v.ValidationErrors); ok {
		e := make(map[string]string)
		for _, f := range ferr {
			e[f.Namespace()] = f.Tag()
		}
		err = response.NewBadRequest(e, "failed to validate request body")
	}
	return
}
