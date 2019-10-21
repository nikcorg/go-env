# Go Env

A utility for reading the environment with assertions

## Usage

```go
package env

import (
	"log"
	"os"

	env "github.com/nikcorg/go-env"
)

type AppEnv struct {
	Beep env.NonEmptyString `env:"BEEP"`
	Boop env.Enum           `env:"BOOP" enum:"testing,one,two"`
}

var appEnv AppEnv

func main() {
	if err := env.Validate(&appEnv, os.Getenv); err != nil {
		log.Fatalf("Invalid environment: %v", err)
	}
	log.Printf("Beep is %s, Boop is %s", appEnv.Beep, appEnv.Boop)
}
```

## Supported validations

- `Int` / `NonEmptyInt` - ensure the value is parseable as a number
- `URL` / `NonEmptyUrl` - ensure the value is parseable as a URL
- `String` / `NonEmptyString` - no formal validation
- `Enum` / `NonEmptyEnum` - ensure the value matches one of the enumerated set of acceptable values

Each validation's `NonEmpty*` variant adds an additional assertion on the value not being `nil`.
