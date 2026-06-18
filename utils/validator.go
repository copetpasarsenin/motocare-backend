package utils

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = newValidator()

func newValidator() *validator.Validate {
	validatorInstance := validator.New()
	validatorInstance.RegisterTagNameFunc(func(field reflect.StructField) string {
		jsonName := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
		if jsonName == "" || jsonName == "-" {
			return field.Name
		}

		return jsonName
	})

	return validatorInstance
}

func ValidateStruct(payload any) map[string]string {
	errors := make(map[string]string)

	if err := validate.Struct(payload); err != nil {
		for _, fieldError := range err.(validator.ValidationErrors) {
			field := fieldError.Field()
			errors[field] = validationMessage(field, fieldError.Tag(), fieldError.Param())
		}
	}

	if len(errors) == 0 {
		return nil
	}

	return errors
}

func validationMessage(field string, tag string, param string) string {
	switch tag {
	case "required":
		return field + " wajib diisi"
	case "email":
		return "email harus valid"
	case "min":
		return field + " minimal " + param + " karakter"
	case "gte":
		return field + " tidak boleh negatif"
	case "oneof":
		return field + " tidak valid"
	default:
		return field + " tidak valid"
	}
}
