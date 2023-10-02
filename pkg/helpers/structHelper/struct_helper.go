package structHelper

import (
	"reflect"
	"strings"
)

func GetFieldsAsJsonTags(str interface{}) []string {
	result := make([]string, 0)

	val := reflect.ValueOf(str).Elem()
	t := val.Type()

	for i := 0; i < t.NumField(); i++ {
		result = append(result, t.Field(i).Tag.Get("json"))
	}

	return result
}

func GetFieldsAsUpperSnake(str interface{}) []string {
	result := make([]string, 0)

	for _, v := range GetFieldsAsJsonTags(str) {
		result = append(result, strings.ToUpper(v))
	}

	return result
}
