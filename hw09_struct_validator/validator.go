package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var errString string

	for _, err := range v {
		errString += fmt.Sprintf("%s: %s\n", err.Err.Error(), err.Field)
	}

	return errString
}

func Validate(v interface{}) error {
	rv := reflect.ValueOf(v)

	if rv.Kind() != reflect.Struct {
		return errors.New("not struct")
	}

	rt := reflect.TypeOf(v)
	path := rt.Name()
	var errList ValidationErrors
	validateStruct(v, &path, &errList)

	return errList
}

func validateStruct(s interface{}, path *string, errList *ValidationErrors) {
	v := reflect.ValueOf(s)

	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		fieldPath := *path + "." + field.Name

		if field.Type.Kind() == reflect.Struct {
			rule, ok := field.Tag.Lookup("validate")
			if ok && rule == "nested" {
				validateStruct(v.Field(i).Interface(), &fieldPath, errList)
			}

			continue
		}

		if rules, ok := field.Tag.Lookup("validate"); ok {
			err := validateField(v.Field(i), rules)

			if err != nil {
				*errList = append(*errList, ValidationError{
					Field: fieldPath,
					Err:   err,
				})
			}
		}
	}
}

func validateField(field reflect.Value, rules string) error {
	var err error

	switch kind := field.Kind(); {
	case kind == reflect.String:
		err = validateString(field, rules)
	case kind == reflect.Slice:
		err = validateSlice(field, rules)
	case kind >= reflect.Int && kind <= reflect.Uint64:
		err = validateInt(field, rules)
	default:
		return errors.New("not supported")
	}

	return err
}
