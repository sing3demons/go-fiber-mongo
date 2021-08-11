package helper

import "github.com/go-playground/validator"

type ErrorResponse struct {
	FailedField string
	Tag         string
	Value       string
}

func ValidateStruct(v interface{}) []*ErrorResponse {
	var errors []*ErrorResponse
	validate := validator.New()
	err := validate.Struct(v)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}
