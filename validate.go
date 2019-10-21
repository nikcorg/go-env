package env

import (
	"errors"
	"fmt"
	"math"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

const tagName = "env"

// Validate reads and validates the environment values
func Validate(a interface{}, ValueFromEnv func(string) string) error {
	reflectType := reflect.TypeOf(a)

	if reflectType.Kind() != reflect.Ptr {
		return errors.New("Not a pointer value")
	}

	if reflect.ValueOf(a).IsNil() {
		return fmt.Errorf("Unsupported pointer type: %s", reflectType.Elem().String())
	}

	rval := reflect.ValueOf(a)

	finalValue, err := getValue(reflectType.Elem(), ValueFromEnv)
	if err != nil {
		return err
	}

	rval.Elem().Set(finalValue)
	return nil
}

func getValue(t reflect.Type, ValueFromEnv func(string) string) (reflect.Value, error) {
	k := t.Kind()

	if k != reflect.Struct {
		return reflect.Value{}, fmt.Errorf("expected a Struct, got %s", k)
	}

	v := reflect.New(t).Elem()

	for i := 0; i < v.NumField(); i++ {
		f := t.Field(i)

		if !v.Field(i).CanSet() {
			return reflect.Value{}, fmt.Errorf("unexported field error: %s", f.Name)
		}

		tag := f.Tag.Get(tagName)
		val := ValueFromEnv(tag)
		typ := v.Field(i).Type()

		switch typ.String() {
		case "env.Int":
			valid, err := asInt(val)
			if err != nil {
				return reflect.Value{}, err
			}
			v.Field(i).Set(reflect.ValueOf(valid).Convert(typ))
			break
		case "env.NonEmptyInt":
			valid, err := asNotEmptyInt(val)
			if err != nil {
				return reflect.Value{}, err
			}
			v.Field(i).Set(reflect.ValueOf(valid).Convert(typ))
			break
		case "env.String":
			v.Field(i).Set(reflect.ValueOf(val).Convert(typ))
			break
		case "env.NonEmptyString":
			valid, err := asNotEmpty(val)
			if err != nil {
				return reflect.Value{}, err
			}
			v.Field(i).Set(reflect.ValueOf(valid).Convert(typ))
			break
		case "env.URL":
			valid, err := asURL(val)
			if err != nil {
				return reflect.Value{}, err
			}
			v.Field(i).Set(reflect.ValueOf(valid).Convert(typ))
			break
		case "env.NonEmptyURL":
			valid, err := asNotEmptyURL(val)
			if err != nil {
				return reflect.Value{}, err
			}
			v.Field(i).Set(reflect.ValueOf(valid).Convert(typ))
			break
		case "env.Enum":
			valid, err := asEnum(val, f.Tag.Get("enum"))
			if err != nil {
				return reflect.Value{}, err
			}
			v.Field(i).Set(reflect.ValueOf(valid).Convert(typ))
			break
		case "env.NonEmptyEnum":
			valid, err := asNotEmptyEnum(val, f.Tag.Get("enum"))
			if err != nil {
				return reflect.Value{}, err
			}
			v.Field(i).Set(reflect.ValueOf(valid).Convert(typ))
			break
		default:
			return reflect.Value{}, fmt.Errorf("unknown field type: %s", typ)
		}
	}

	return v, nil
}

// asNotEmpty validates input is not the empty string
func asNotEmpty(s string) (string, error) {
	if s == "" {
		return s, fmt.Errorf("Expected string to be non-empty")
	}

	return s, nil
}

// asInt validates the input can be parsed into a number value
// If an input value is not present the returned value is math.NaN
func asInt(s string) (int, error) {
	if s == "" {
		return int(math.NaN()), nil
	}
	num, err := strconv.Atoi(s)
	if err != nil {
		return int(math.NaN()), err
	}

	return num, nil
}

// asNotEmptyInt validates the input is not empty and is a number
func asNotEmptyInt(s string) (int, error) {
	_, err := asNotEmpty(s)
	if err != nil {
		return int(math.NaN()), err
	}
	return asInt(s)
}

// asURL validates that the input can be parsed as a URL
func asURL(s string) (string, error) {
	if s == "" {
		return "", nil
	}
	u, err := url.Parse(s)
	if err != nil {
		return "", fmt.Errorf("expected a valid URL, got %v", s)
	}
	if u.Scheme == "" || u.Host == "" {
		return "", fmt.Errorf("expected a valid URL including a scheme and a host, got %v", s)
	}
	return u.String(), nil
}

// asNotEmptyURL validates that the input is not empty and can be parsed as a URL
func asNotEmptyURL(s string) (string, error) {
	_, err := asNotEmpty(s)
	if err != nil {
		return "", err
	}
	return asURL(s)
}

// asEnum validates that the input exists in a defined set of values
func asEnum(s string, vals string) (string, error) {
	if s == "" {
		return "", nil
	}
	vs := strings.Split(vals, ",")
	if !contains(vs, s) {
		return "", fmt.Errorf("invalid enum value: %s", s)
	}
	return s, nil
}

// asNotEmptyEnum validates that the input is not empty
func asNotEmptyEnum(s, vals string) (string, error) {
	_, err := asNotEmpty(s)
	if err != nil {
		return "", err
	}
	return asEnum(s, vals)
}

func contains(xs []string, x string) bool {
	for _, v := range xs {
		if v == x {
			return true
		}
	}
	return false
}
