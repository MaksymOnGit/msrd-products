package utils

import "github.com/go-playground/validator/v10"

func Validate[T any](s T) (fields map[string]string) {
	validate := validator.New()
	err := validate.Struct(s)
	if err != nil {
		fields = map[string]string{}
		for _, err := range err.(validator.ValidationErrors) {
			fields[err.Field()] = err.Error()
		}
	}
	return
}
