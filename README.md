# Go Env

A utility for reading the environment with assertions

## Usage

```go
package main

import (
	"log"
	"os"

	env "github.com/nikcorg/go-env"
)

type AppEnv struct {
	// Require the value `BEEP` is not empty
	Beep env.NonEmptyString `env:"BEEP"`

	// Require the value `BOOP` to be one of `tesing`, `one` or `two`
	Boop env.Enum           `env:"BOOP" enum:"testing,one,two"`

	// Use the value from `BRRT` or default to `ding dong`
	Brrt env.String         `env:"BRRT" default:"ding dong"`

	// Split the value from `BZZT` using `:` and parse as int values
	Bzzt env.IntSlice       `env:"BZZT" separator:":"`
}

var appEnv AppEnv

func main() {
	if err := env.New(&appEnv).Validate(); err != nil {
		log.Fatalf("Invalid environment: %v", err)
	}
	log.Printf("Beep is %s, Boop is %s, Brrt is %s", appEnv.Beep, appEnv.Boop, appEnv.Brrt)
}
```

## Supported validations

- `Enum` - ensure the value matches one of the enumerated set of acceptable values
- `HostPort` - takes a string value and ensures it can be parsed by `net.SplitHostPort`
- `IntSlice` takes a CSV value and splits it into an int slice using a separator
- `Int` - ensures the value is parseable as a number
- `StringSlice` - takes a CSV value and splits it into a string slice using a separator
- `String` - no formal validation
- `URL` - ensures the value is parseable as a `url.URL` and has a non-empty `Scheme` and a `Host` value

Each validation `T` has a `NonEmptyT` variant, which adds an additional assertion on the value not being unset.
