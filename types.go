package env

import (
	"fmt"
	"strconv"
	"strings"
)

// Env is in interface for reading the environment with assertions
type Env interface {
	Validate(interface{}) interface{}
}

// Int is an optional int value
type Int int

func (v Int) String() string { return fmt.Sprintf("%d", v) }

// NonEmptyInt is a required int value
type NonEmptyInt int

func (v NonEmptyInt) String() string { return fmt.Sprintf("%d", v) }

// URL is an optional URL value
type URL string

func (v URL) String() string { return string(v) }

// NonEmptyURL is a required URL value
type NonEmptyURL string

func (v NonEmptyURL) String() string { return string(v) }

// HostPort is an optional value that is successfully parsed using net.SplitHostPort
type HostPort struct {
	Host string
	Port string
}

func (v HostPort) String() string {
	if strings.Contains(v.Host, ":") {
		return fmt.Sprintf("[%s]:%s", v.Host, v.Port)
	}
	return fmt.Sprintf("%s:%s", v.Host, v.Port)
}

// NonEmptyHostPort is a required HostPort value
type NonEmptyHostPort HostPort

func (v NonEmptyHostPort) String() string { return HostPort(v).String() }

// String is an optional string value
type String string

func (v String) String() string { return string(v) }

// NonEmptyString is a required string value
type NonEmptyString string

func (v NonEmptyString) String() string { return string(v) }

// Enum is an enumerated set of valid string values
type Enum string

func (x Enum) String() string { return string(x) }

// NonEmptyEnum is a required Enum value
type NonEmptyEnum string

func (x NonEmptyEnum) String() string { return string(x) }

// StringSlice is a CSV value
type StringSlice []string

func (x StringSlice) String() string { return strings.Join(x, ",") }

// NonEmptyStringSlice is a StringSlice value with a length > 0 requirement
type NonEmptyStringSlice []string

func (x NonEmptyStringSlice) String() string { return strings.Join(x, ",") }

// IntSlice is a CSV value
type IntSlice []int

func (x IntSlice) String() string {
	out := ""
	for _, s := range x {
		out += "," + strconv.Itoa(s)
	}
	return out[1:]
}

// NonEmptyIntSlice is an IntSlice value with a length > 0 requirement
type NonEmptyIntSlice []int

func (x NonEmptyIntSlice) String() string {
	out := ""
	for _, s := range x {
		out += "," + strconv.Itoa(s)
	}
	return out[1:]
}
