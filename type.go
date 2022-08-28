package utils

import (
	"errors"
	"strings"
)

type Type string

const (
	JSON Type = "json"
	YAML Type = "yaml"
	TOML Type = "toml"
)

func ParseConfigType(configType string) (Type, error) {
	switch strings.ToLower(configType) {
	case "json":
		return JSON, nil
	case "yaml", "": // Empty string as default
		return YAML, nil
	case "toml":
		return TOML, nil
	default:
		return "", errors.New("Invalid Configuration Type: JSON, YAML, TOML or \"\"(empty string), Given: " + configType)
	}
}
