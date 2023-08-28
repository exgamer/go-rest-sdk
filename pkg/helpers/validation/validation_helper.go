package validation

import (
	"github.com/gookit/validate"
)

func ValidationErrorsAsMap(validationErrors validate.Errors) *map[string]string {
	eMap := make(map[string]string)

	for k, ve := range validationErrors {
		eMap[k] = ve.String()
	}

	return &eMap
}
