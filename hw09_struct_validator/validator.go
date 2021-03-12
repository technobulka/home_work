package hw09structvalidator

import (
	"fmt"
	"reflect"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	panic("implement me")
}

func Validate(v interface{}) error {
	rv := reflect.ValueOf(v)

	for i := 0; i < rv.NumField(); i++ {
		field := rv.Type().Field(i)
		//var fieldPath string
		//if currentPath == "" {
		//	fieldPath = field.Name
		//} else {
		//	fieldPath = currentPath + "." + field.Name
		//}
		//if field.Type.Kind() == reflect.Struct {
		//	validateFields(rv.Field(i).Interface(), fieldPath, errs)
		//	continue
		//}

		if field.Type.Kind() == reflect.Struct {
			_ = validateStruct(rv.Field(i).Interface())
			continue
		}

		if validate, ok := field.Tag.Lookup("validate"); ok {
			fmt.Println(validate)
		}
	}

	return nil
}

func validateStruct(s interface{}) error {
	v := reflect.ValueOf(s)

	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)

		if validate, ok := field.Tag.Lookup("validate"); ok {
			fmt.Println(validate)
		}
	}

	return nil
}
