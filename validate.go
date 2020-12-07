package env

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
)

const envTag = "env"
const fallbackTag = "default"

// Known error outcomes
var (
	ErrExpectedAtLeastOneValue = errors.New("expected at least one value")
	ErrUnexpectedEmptyValue    = errors.New("unexpected empty value")
	ErrExpectedPointerValue    = errors.New("expected a pointer value")
	ErrUnexpectedNilPointer    = errors.New("unexpected nil-pointer")
	ErrExpectedStructValue     = errors.New("expected a struct")
	ErrUnsettableField         = errors.New("unsettable field")
	ErrUntaggedField           = errors.New("untagged field")
	ErrUnknownFieldType        = errors.New("unknown field type")
	ErrPartialURLValue         = errors.New("expected url to have Scheme and Host")
	ErrInvalidEnumValue        = errors.New("invalid enum value")
)

// Options represents the library's configurable traits
type Options struct {
	Getenv Getter
}

// AssertedEnvironment represents an environment configuration and a value getter
type AssertedEnvironment struct {
	config interface{}
	opts   *Options
}

// Getter is used to retrieve values for populating an environment structure
type Getter func(string) string

var defaultConfig = Options{os.Getenv}

// New constructs a new AssertedEnvironment using a provided value getter
func New(config interface{}, opts ...*Options) *AssertedEnvironment {
	options := &Options{defaultConfig.Getenv}

	for _, o := range opts {
		options.Getenv = o.Getenv
	}

	return &AssertedEnvironment{config, options}
}

// Validate reads and validates the environment values
func (e *AssertedEnvironment) Validate() error {
	return validate(e.config, e.opts.Getenv)
}

// MustValidate validates the environment and panics on any validation error
func (e *AssertedEnvironment) MustValidate() {
	if err := validate(e.config, e.opts.Getenv); err != nil {
		panic(err)
	}
}

func validate(a interface{}, getenv Getter) error {
	reflectType := reflect.TypeOf(a)

	if reflectType.Kind() != reflect.Ptr {
		return ErrExpectedPointerValue
	}

	if reflect.ValueOf(a).IsNil() {
		return ErrUnexpectedNilPointer
	}

	rval := reflect.ValueOf(a)

	finalValue, err := getValue(reflectType.Elem(), getenv)
	if err != nil {
		return err
	}

	rval.Elem().Set(finalValue)
	return nil
}

func getValue(t reflect.Type, getenv Getter) (reflect.Value, error) {
	k := t.Kind()

	if k != reflect.Struct {
		return reflect.Value{}, ErrExpectedStructValue
	}

	v := reflect.New(t).Elem()

	for i := 0; i < v.NumField(); i++ {
		f := t.Field(i)

		if !v.Field(i).CanSet() {
			return reflect.Value{}, fmt.Errorf("%w: %s", ErrUnsettableField, f.Name)
		}

		var (
			ok                           bool
			candidate, envName, fallback string
			typ                          reflect.Type
		)

		if envName, ok = f.Tag.Lookup(envTag); !ok {
			return reflect.Value{}, fmt.Errorf("%w: %s", ErrUntaggedField, f.Name)
		}

		candidate = getenv(envName)

		if candidate == "" {
			if fallback, ok = f.Tag.Lookup(fallbackTag); ok {
				candidate = fallback
			}
		}

		typ = v.Field(i).Type()

		switch typ.String() {
		case "env.Int":
			valid, err := asInt(candidate)
			if err != nil {
				return reflect.Value{}, err
			}
			v.Field(i).Set(reflect.ValueOf(valid).Convert(typ))

		case "env.NonEmptyInt":
			valid, err := asNotEmptyInt(candidate)
			if err != nil {
				return reflect.Value{}, err
			}
			v.Field(i).Set(reflect.ValueOf(valid).Convert(typ))

		case "env.String":
			v.Field(i).Set(reflect.ValueOf(candidate).Convert(typ))

		case "env.NonEmptyString":
			valid, err := asNotEmpty(candidate)
			if err != nil {
				return reflect.Value{}, err
			}
			v.Field(i).Set(reflect.ValueOf(valid).Convert(typ))

		case "env.URL":
			valid, err := asURL(candidate)
			if err != nil {
				return reflect.Value{}, err
			}
			v.Field(i).Set(reflect.ValueOf(valid).Convert(typ))

		case "env.NonEmptyURL":
			valid, err := asNotEmptyURL(candidate)
			if err != nil {
				return reflect.Value{}, err
			}
			v.Field(i).Set(reflect.ValueOf(valid).Convert(typ))

		case "env.Enum":
			valid, err := asEnum(candidate, f.Tag.Get("enum"))
			if err != nil {
				return reflect.Value{}, err
			}
			v.Field(i).Set(reflect.ValueOf(valid).Convert(typ))

		case "env.NonEmptyEnum":
			valid, err := asNotEmptyEnum(candidate, f.Tag.Get("enum"))
			if err != nil {
				return reflect.Value{}, err
			}
			v.Field(i).Set(reflect.ValueOf(valid).Convert(typ))

		case "env.StringSlice":
			valid, err := asStringSlice(candidate, f.Tag.Get("separator"))
			if err != nil {
				return reflect.Value{}, err
			}
			v.Field(i).Set(reflect.ValueOf(valid).Convert(typ))

		case "env.NonEmptyStringSlice":
			valid, err := asNotEmptyStringSlice(candidate, f.Tag.Get("separator"))
			if err != nil {
				return reflect.Value{}, err
			}
			v.Field(i).Set(reflect.ValueOf(valid).Convert(typ))

		case "env.IntSlice":
			valid, err := asIntSlice(candidate, f.Tag.Get("separator"))
			if err != nil {
				return reflect.Value{}, err
			}
			v.Field(i).Set(reflect.ValueOf(valid).Convert(typ))

		case "env.NonEmptyIntSlice":
			valid, err := asNotEmptyIntSlice(candidate, f.Tag.Get("separator"))
			if err != nil {
				return reflect.Value{}, err
			}
			v.Field(i).Set(reflect.ValueOf(valid).Convert(typ))

		case "env.HostPort":
			valid, err := asHostPort(candidate)
			if err != nil {
				return reflect.Value{}, err
			}
			v.Field(i).Set(reflect.ValueOf(valid).Convert(typ))

		case "env.NonEmptyHostPort":
			valid, err := asNotEmptyHostPort(candidate)
			if err != nil {
				return reflect.Value{}, err
			}
			v.Field(i).Set(reflect.ValueOf(valid).Convert(typ))

		default:
			return reflect.Value{}, fmt.Errorf("%w: %s", ErrUnknownFieldType, typ)
		}
	}

	return v, nil
}

// asNotEmpty validates input is not the empty string
func asNotEmpty(s string) (string, error) {
	if s == "" {
		return s, ErrUnexpectedEmptyValue
	}

	return s, nil
}

// asInt validates the input can be parsed into a number value
// If an input value is not present the returned value is 0
func asInt(s string) (int, error) {
	if s == "" {
		return 0, nil
	}
	num, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}

	return num, nil
}

// asNotEmptyInt validates the input is not empty and is a number
func asNotEmptyInt(s string) (int, error) {
	_, err := asNotEmpty(s)
	if err != nil {
		return 0, err
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
		return "", err
	}
	if u.Scheme == "" || u.Host == "" {
		return "", fmt.Errorf("%w: %s", ErrPartialURLValue, s)
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
		return "", fmt.Errorf("%w: %s", ErrInvalidEnumValue, s)
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

// asStringSlice splits a string value by a separator
// If a separator is not provided the comma is used
func asStringSlice(s, separator string) ([]string, error) {
	if len(s) == 0 {
		return []string{}, nil
	}

	splitBy := separator
	if splitBy == "" {
		splitBy = ","
	}
	return strings.Split(s, splitBy), nil
}

// asNotEmptyStringSlice validates that the input is not empty
// It does not validate non-zero length values
func asNotEmptyStringSlice(s, separator string) ([]string, error) {
	if _, err := asNotEmpty(s); err != nil {
		return nil, err
	}

	v, err := asStringSlice(s, separator)
	if err != nil {
		return nil, err
	}
	if len(v) == 0 {
		return nil, ErrExpectedAtLeastOneValue
	}
	return v, nil
}

type intParser func(string) (int, error)

// asIntSlice splits a string value by a separator and validates that the values can be parsed as an int
// If a separator is not provided the comma is used
func asIntSlice(s, separator string) ([]int, error) {
	stringVals, err := asStringSlice(s, separator)
	if err != nil {
		return nil, err
	}
	intVals := make([]int, len(stringVals))
	for n := range stringVals {
		if intVals[n], err = asNotEmptyInt(stringVals[n]); err != nil {
			return nil, err
		}
	}
	return intVals, nil
}

// asNotEmptyIntSlice validates that the input is not empty
func asNotEmptyIntSlice(s, separator string) ([]int, error) {
	if _, err := asNotEmpty(s); err != nil {
		return nil, err
	}

	v, err := asIntSlice(s, separator)
	if err != nil {
		return nil, err
	} else if len(v) == 0 {
		return nil, ErrExpectedAtLeastOneValue
	}
	return v, nil
}

// asHostPort validates that a value is successfully parsed by net.SplitHostPort
func asHostPort(s string) (HostPort, error) {
	if s == "" {
		return HostPort{}, nil
	}

	host, port, err := net.SplitHostPort(s)
	if err != nil {
		return HostPort{}, err
	}

	return HostPort{host, port}, nil
}

// asNotEmptyHostPort validates the the input is not empty
func asNotEmptyHostPort(s string) (NonEmptyHostPort, error) {
	if _, err := asNotEmpty(s); err != nil {
		return NonEmptyHostPort{}, err
	}
	hp, err := asHostPort(s)
	if err != nil {
		return NonEmptyHostPort{}, err
	}
	return NonEmptyHostPort(hp), nil
}

func contains(xs []string, x string) bool {
	for _, v := range xs {
		if v == x {
			return true
		}
	}
	return false
}
