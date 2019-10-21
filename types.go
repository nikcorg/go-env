package env

// Env is in interface for reading the environment with assertions
type Env interface {
	Validate(interface{}) interface{}
}

// Int is an optional int value
type Int int

func (v Int) String() string { return string(v) }

// NonEmptyInt is a required int value
type NonEmptyInt int

func (v NonEmptyInt) String() string { return string(v) }

// URL is an optional URL value
type URL string

func (v URL) String() string { return string(v) }

// NonEmptyURL is a required URL value
type NonEmptyURL string

func (v NonEmptyURL) String() string { return string(v) }

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
