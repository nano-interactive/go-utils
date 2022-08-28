package utils

import (
	"errors"
	"strings"
)

type Env uint8

var ErrInvalidEnvironment = errors.New("invalid environment")

const (
	Testing Env = iota
	Development
	Production
)

func ParseEnvironment(env string) (Env, error) {
	switch strings.ToLower(env) {
	case "prod", "production":
		return Production, nil
	case "dev", "development", "develop":
		return Development, nil
	case "testing", "test":
		return Testing, nil
	default:
		return 0, errors.New("Invalid Configuration Type: JSON, YAML, TOML or \"\"(empty string), Given: " + env)
	}
}
