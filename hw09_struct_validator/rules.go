package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type Rules []string

func invalidRule(rule string) error {
	return fmt.Errorf("invalid rule `%s`", rule)
}

func validateString(value reflect.Value, rulesString string) error {
	var err error
	var errs []string

	rules := strings.Split(rulesString, "|")
	for _, rule := range rules {
		switch {
		case strings.HasPrefix(rule, "len:"):
			err = ruleLen(value, rule[4:])
		case strings.HasPrefix(rule, "in:"):
			err = ruleIn(value, rule[3:])
		case strings.HasPrefix(rule, "regexp:"):
			err = ruleRegexp(value, rule[7:])
		default:
			err = invalidRule(rule)
		}

		if err != nil {
			errs = append(errs, fmt.Sprint(err))
		}
	}

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, ", "))
	}

	return nil
}

func validateInt(value reflect.Value, rulesString string) error {
	var err error
	var errs []string

	rules := strings.Split(rulesString, "|")
	for _, rule := range rules {
		switch {
		case strings.HasPrefix(rule, "in:"):
			err = ruleIn(value, rule[3:])
		case strings.HasPrefix(rule, "min:"):
			err = ruleMin(value, rule[4:])
		case strings.HasPrefix(rule, "max:"):
			err = ruleMax(value, rule[4:])
		default:
			err = invalidRule(rule)
		}

		if err != nil {
			errs = append(errs, fmt.Sprint(err))
		}
	}

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, ", "))
	}

	return nil
}

func validateSlice(value reflect.Value, rulesString string) error {
	var err error
	var errs []string

	for i := 0; i < value.Len(); i++ {
		switch value.Type().String()[2:5] {
		case "str":
			err = validateString(value.Index(i), rulesString)
		case "int", "uin":
			err = validateInt(value.Index(i), rulesString)
		default:
			err = invalidRule(rulesString)
		}

		if err != nil {
			errs = append(errs, fmt.Sprint(err))
		}
	}

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, ", "))
	}

	return nil
}

func ruleLen(value reflect.Value, length string) error {
	l, err := strconv.Atoi(length)
	if err != nil {
		return invalidRule("len:" + length)
	}

	if value.Len() != l {
		return errors.New("invalid length")
	}

	return nil
}

func ruleRegexp(value reflect.Value, re string) error {
	regExp, err := regexp.Compile(re)
	if err != nil {
		return invalidRule("regexp:" + re)
	}

	if !regExp.MatchString(value.String()) {
		return errors.New("invalid regexp")
	}

	return nil
}

func ruleIn(value reflect.Value, in string) error {
	var valid bool
	set := strings.Split(in, ",")

	switch kind := value.Kind(); {
	case kind == reflect.String:
		valid = containsStr(value.String(), set)
	case kind >= reflect.Int && kind <= reflect.Int64:
		valid = containsInt(value.Int(), set)
	case kind >= reflect.Uint && kind <= reflect.Uint64:
		valid = containsUint(value.Uint(), set)
	default:
		return invalidRule("in:" + in)
	}

	if !valid {
		return errors.New("invalid contains")
	}

	return nil
}

func containsInt(value int64, list []string) bool {
	var valueStr string
	rv := reflect.ValueOf(value)

	value = rv.Convert(reflect.TypeOf(value)).Int()
	valueStr = strconv.FormatInt(value, 10)

	return containsStr(valueStr, list)
}

func containsUint(value uint64, list []string) bool {
	var valueStr string
	rv := reflect.ValueOf(value)

	value = rv.Convert(reflect.TypeOf(value)).Uint()
	valueStr = strconv.FormatUint(value, 10)

	return containsStr(valueStr, list)
}

func containsStr(value string, list []string) bool {
	for _, el := range list {
		if el == value {
			return true
		}
	}

	return false
}

func ruleMin(value reflect.Value, min string) error {
	var valid bool
	var err error

	switch kind := value.Kind(); {
	case kind >= reflect.Int && kind <= reflect.Int64:
		valid, err = comparisonInt(value.Int(), min, "min")
	case kind >= reflect.Uint && kind <= reflect.Uint64:
		valid, err = comparisonUint(value.Uint(), min, "min")
	default:
		return invalidRule("min:" + min)
	}

	if err != nil {
		return err
	}

	if !valid {
		return errors.New("invalid min")
	}

	return nil
}

func ruleMax(value reflect.Value, max string) error {
	var valid bool
	var err error

	switch kind := value.Kind(); {
	case kind >= reflect.Int && kind <= reflect.Int64:
		valid, err = comparisonInt(value.Int(), max, "max")
	case kind >= reflect.Uint && kind <= reflect.Uint64:
		valid, err = comparisonUint(value.Uint(), max, "max")
	default:
		return invalidRule("max:" + max)
	}

	if err != nil {
		return err
	}

	if !valid {
		return errors.New("invalid max")
	}

	return nil
}

func comparisonInt(value int64, target string, direction string) (bool, error) {
	compared, err := strconv.ParseInt(target, 10, 64)
	if err != nil {
		return false, invalidRule(direction + ":" + target)
	}

	switch direction {
	case "max":
		return value <= compared, nil
	case "min":
		return value >= compared, nil
	default:
		return false, nil
	}
}

func comparisonUint(value uint64, target string, direction string) (bool, error) {
	compared, err := strconv.ParseUint(target, 10, 64)
	if err != nil {
		return false, invalidRule(direction + ":" + target)
	}

	switch direction {
	case "max":
		return value <= compared, nil
	case "min":
		return value >= compared, nil
	default:
		return false, nil
	}
}
