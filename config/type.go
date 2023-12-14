package config

import (
	"errors"
	"strings"
)

type Type string

var ErrInvalidConfigType = errors.New("invalid Configuration Type: JSON, YAML, TOML or \"\"(empty string)")

const (
	JSON Type = "json"
	YAML Type = "yaml"
	TOML Type = "toml"
)

func ParseType(configType string) (Type, error) {
	switch strings.ToLower(configType) {
	case "json":
		return JSON, nil
	case "yaml", "": // Empty string as default
		return YAML, nil
	case "toml":
		return TOML, nil
	default:
		return "", ErrInvalidConfigType
	}
}

func MustParseType(configType string) Type {
	t, err := ParseType(configType)
	if err != nil {
		panic(err)
	}

	return t
}
