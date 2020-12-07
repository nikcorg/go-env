package env

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type intTestCase struct {
	input         string
	expectedValue int
	shouldError   bool
}

func TestAsInt(t *testing.T) {
	tests := []intTestCase{
		{"", 0, false},
		{"0", 0, false},
		{"123", 123, false},
		{"hello", 0, true},
	}

	runIntTestCases(t, tests, asInt)
}

func TestAsNotEmptyInt(t *testing.T) {
	tests := []intTestCase{
		{"", 0, true},
		{"0", 0, false},
		{"123", 123, false},
		{"hello", 0, true},
	}

	runIntTestCases(t, tests, asNotEmptyInt)
}

func runIntTestCases(t *testing.T, tests []intTestCase, f func(string) (int, error)) {
	for _, x := range tests {
		v, err := f(x.input)

		if x.shouldError {
			assert.Error(t, err)
			assert.Zero(t, v)
		} else {
			assert.Nil(t, err)
			assert.Equal(t, x.expectedValue, v)
		}
	}
}

type stringTestCase struct {
	input         string
	expectedValue string
	shouldError   bool
}

func TestAsNotEmptyString(t *testing.T) {
	tests := []stringTestCase{
		{"", "", true},
		{"hello world", "hello world", false},
	}

	runStringTestCases(t, tests, asNotEmpty)
}

func runStringTestCases(t *testing.T, tests []stringTestCase, f func(string) (string, error)) {
	for _, x := range tests {
		v, err := f(x.input)

		if x.shouldError {
			assert.Error(t, err)
			assert.Zero(t, v)
		} else {
			assert.Nil(t, err)
			assert.Equal(t, x.expectedValue, v)
		}
	}
}

type enumTestCase struct {
	input         string
	enum          string
	expectedValue string
	shouldError   bool
}

func TestEnum(t *testing.T) {
	tests := []enumTestCase{
		{"", "one,two", "", false},
		{"one", "one,two", "one", false},
		{"hello", "one,two", "", true},
	}

	runEnumTestCases(t, tests, asEnum)
}

func TestNotEmptyEnum(t *testing.T) {
	tests := []enumTestCase{
		{"", "one,two", "", true},
		{"one", "one,two", "one", false},
		{"hello", "one,two", "", true},
	}

	runEnumTestCases(t, tests, asNotEmptyEnum)
}

func runEnumTestCases(t *testing.T, tests []enumTestCase, f func(string, string) (string, error)) {
	for _, x := range tests {
		v, err := f(x.input, x.enum)
		if x.shouldError {
			assert.Error(t, err)
			assert.Zero(t, v)
		} else {
			assert.Nil(t, err)
			assert.Equal(t, x.expectedValue, v)
		}
	}
}

type urlTestCase struct {
	input         string
	expectedValue string
	shouldError   bool
}

func TestURL(t *testing.T) {
	tests := []urlTestCase{
		{"", "", false},
		{"hello world", "", true},
		{"localhost:8080", "", true},
		{"https://github.com/nikcorg/go-env", "https://github.com/nikcorg/go-env", false},
		{"http://localhost:8080", "http://localhost:8080", false},
	}

	runURLTestCases(t, tests, asURL)
}

func TestNotEmptyURL(t *testing.T) {
	tests := []urlTestCase{
		{"", "", true},
		{"hello world", "", true},
		{"localhost:8080", "", true},
		{"https://github.com/nikcorg/go-env", "https://github.com/nikcorg/go-env", false},
		{"http://localhost:8080", "http://localhost:8080", false},
	}
	runURLTestCases(t, tests, asNotEmptyURL)
}

func runURLTestCases(t *testing.T, tests []urlTestCase, f func(string) (string, error)) {
	for _, x := range tests {
		v, err := f(x.input)
		if x.shouldError {
			assert.Error(t, err)
			assert.Zero(t, v)
		} else {
			assert.Nil(t, err)
			assert.Equal(t, x.expectedValue, v)
		}
	}
}

type stringSliceTestCase struct {
	input         string
	separator     string
	expectedValue []string
	shouldError   bool
}

func TestStringSlice(t *testing.T) {
	tests := []stringSliceTestCase{
		{"", "", []string{}, false},
		{"hello,world", "", []string{"hello", "world"}, false},
		{"beep,boop:bleep", ":", []string{"beep,boop", "bleep"}, false},
		{",,,", "", []string{"", "", "", ""}, false},
	}

	runStringSliceTestCases(t, tests, asStringSlice)
}

func TestNonEmptyStringSlice(t *testing.T) {
	tests := []stringSliceTestCase{
		{"", "", []string{}, true},
		{"hello,world", "", []string{"hello", "world"}, false},
		{"beep,boop:bleep", ":", []string{"beep,boop", "bleep"}, false},
		{",,,", "", []string{"", "", "", ""}, false},
	}

	runStringSliceTestCases(t, tests, asNotEmptyStringSlice)
}

func runStringSliceTestCases(t *testing.T, tests []stringSliceTestCase, f func(string, string) ([]string, error)) {
	for _, x := range tests {
		v, err := f(x.input, x.separator)
		if x.shouldError {
			assert.Error(t, err)
			assert.Zero(t, v)
		} else {
			assert.Nil(t, err)
			assert.Equal(t, x.expectedValue, v)
		}
	}
}

type intSliceTestCase struct {
	input         string
	separator     string
	expectedValue []int
	shouldError   bool
}

func TestIntSlice(t *testing.T) {
	tests := []intSliceTestCase{
		{"", "", []int{}, false},
		{"1,2,3", "", []int{1, 2, 3}, false},
		{"1:2:3", ":", []int{1, 2, 3}, false},
		{"1,,3", "", []int{}, true},
	}

	runIntSliceTestCases(t, tests, asIntSlice)
}

func TestNonEmptyIntSlice(t *testing.T) {
	tests := []intSliceTestCase{
		{"", "", []int{}, true},
		{"1,2,3", "", []int{1, 2, 3}, false},
		{"1:2:3", ":", []int{1, 2, 3}, false},
		{"1,,3", "", []int{}, true},
	}

	runIntSliceTestCases(t, tests, asNotEmptyIntSlice)
}

func runIntSliceTestCases(t *testing.T, tests []intSliceTestCase, f func(string, string) ([]int, error)) {
	for _, x := range tests {
		v, err := f(x.input, x.separator)
		if x.shouldError {
			assert.Error(t, err, x)
			assert.Zero(t, v, x)
		} else {
			assert.Nil(t, err, x)
			assert.Equal(t, x.expectedValue, v, x)
		}
	}
}

func TestHostPort(t *testing.T) {
	tests := []struct {
		input         string
		expectedValue HostPort
		shouldError   bool
	}{
		{"", HostPort{}, false},
		{"localhost", HostPort{}, true},
		{"localhost:1234", HostPort{"localhost", "1234"}, false},
	}

	for _, test := range tests {
		v, err := asHostPort(test.input)
		if test.shouldError {
			assert.Error(t, err, test)
			assert.Zero(t, v, test)
		} else {
			assert.Nil(t, err, test)
			assert.Equal(t, test.expectedValue, v, test)
		}
	}
}

func TestNonEmptyHostPort(t *testing.T) {
	tests := []struct {
		input         string
		expectedValue NonEmptyHostPort
		shouldError   bool
	}{
		{"", NonEmptyHostPort{}, true},
		{"localhost", NonEmptyHostPort{}, true},
		{"localhost:1234", NonEmptyHostPort{"localhost", "1234"}, false},
	}

	for _, test := range tests {
		v, err := asNotEmptyHostPort(test.input)
		if test.shouldError {
			assert.Error(t, err, test)
			assert.Zero(t, v, test)
		} else {
			assert.Nil(t, err, test)
			assert.Equal(t, test.expectedValue, v, test)
		}
	}
}

func TestSimple(t *testing.T) {
	type simple struct {
		Beep NonEmptyString   `env:"BEEP"`
		Boop Enum             `env:"BOOP" enum:"testing,one,two"`
		Brrt NonEmptyString   `env:"BRRT" default:"fallback value"`
		Bzzt StringSlice      `env:"BZZT" separator:":"`
		Bomf NonEmptyIntSlice `env:"BOMF"`
	}

	tests := []struct {
		name        string
		config      simple
		shouldError bool
		env         map[string]string
		verify      func(*simple, map[string]string)
	}{
		{
			"empty environment",
			simple{},
			true,
			map[string]string{},
			nil,
		},
		{
			"invalid environment",
			simple{},
			true,
			map[string]string{"BEEP": "hello world", "BOMF": ""},
			nil,
		},
		{
			"minimal initialisation",
			simple{},
			false,
			map[string]string{"BEEP": "hello world", "BOMF": "1"},
			func(cf *simple, env map[string]string) {
				assert.Equal(t, NonEmptyString(env["BEEP"]), cf.Beep)
				assert.Equal(t, NonEmptyString("fallback value"), cf.Brrt)
				assert.Equal(t, NonEmptyIntSlice([]int{1}), cf.Bomf)
			},
		},
		{
			"full init",
			simple{},
			false,
			map[string]string{
				"BEEP": "hello world",
				"BOOP": "one",
				"BRRT": "no more fallback",
				"BZZT": "bee:goes:buzz",
				"BOMF": "1,2,3",
			},
			func(cf *simple, env map[string]string) {
				// Beep matches the env value and Brrt matches the fallback value
				assert.Equal(t, NonEmptyString(env["BEEP"]), cf.Beep)
				assert.Equal(t, Enum(env["BOOP"]), cf.Boop)
				assert.Equal(t, NonEmptyString("no more fallback"), cf.Brrt)
				assert.Equal(t, StringSlice(strings.Split(env["BZZT"], ":")), cf.Bzzt)
				assert.Equal(t, NonEmptyIntSlice([]int{1, 2, 3}), cf.Bomf)
			},
		},
	}

	configForEnv := func(e map[string]string) *Options {
		getter := func(k string) string { return e[k] }
		return &Options{getter}
	}

	for _, test := range tests {
		e := New(&test.config, configForEnv(test.env))
		err := e.Validate()

		if test.shouldError {
			assert.Error(t, err, test)
		} else {
			assert.Nil(t, err, test)
			test.verify(&test.config, test.env)
		}
	}
}
