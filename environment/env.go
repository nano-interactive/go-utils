package environment

import (
	"errors"
	"strings"
)

type Env uint8

var ErrInvalidEnv = errors.New("invalid Environment: prod, production, dev, development, develop, testing, test")

const (
	Development Env = iota
	Production
	Testing
	Staging
)

func Parse(env string) (Env, error) {
	switch strings.ToLower(env) {
	case "prod", "production":
		return Production, nil
	case "dev", "development", "develop":
		return Development, nil
	case "testing", "test":
		return Testing, nil
	case "staging", "stage":
		return Staging, nil
	default:
		return 0, ErrInvalidEnv
	}
}

func MustParse(env string) Env {
	e, err := Parse(env)
	if err != nil {
		panic(err)
	}

	return e
}
