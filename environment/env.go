package environment

import (
	"errors"
	"strings"
)

type Env uint8

const (
	Development Env = iota
	Production
	Testing
)

// Parses environment and returns constant
func Parse(env string) (Env, error) {
	switch strings.ToLower(env) {
	case "prod", "production":
		return Production, nil
	case "dev", "development", "develop":
		return Development, nil
	case "testing", "test":
		return Testing, nil
	default:
		return 0, errors.New("Invalid Environment: prod, production, dev, development, develop, testing, test, Given: " + env)
	}
}
