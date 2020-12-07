package main

import (
	"log"

	env "github.com/nikcorg/go-env"
)

type AppEnv struct {
	Beep env.NonEmptyString   `env:"BEEP"`
	Boop env.Enum             `env:"BOOP" enum:"testing,one,two"`
	Brrt env.NonEmptyString   `env:"BRRT" default:"fallback value"`
	Bzzt env.StringSlice      `env:"BZZT" separator:":"`
	Bomf env.NonEmptyIntSlice `env:"BOMF"`
}

var appEnv AppEnv

func main() {
	if err := env.NewFromEnv(&appEnv).Validate(); err != nil {
		log.Fatalf("Invalid environment: %v", err)
	}
	log.Printf(
		"Beep is %s, Boop is %s, Brrt is %s, Bzzt (%d) is %s, Bomf (%d) is %s\n",
		appEnv.Beep, appEnv.Boop, appEnv.Brrt, len(appEnv.Bzzt), appEnv.Bzzt, len(appEnv.Bomf), appEnv.Bomf,
	)
}
