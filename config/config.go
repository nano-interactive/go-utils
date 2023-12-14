package config

import (
	"fmt"

	"github.com/nano-interactive/go-utils/environment"
	"github.com/spf13/viper"
)

type (
	Config struct {
		ProjectName string
		Env         string
		Name        string
		Type        string
		Paths       []string
	}

	// Modifier function
	Modifier func(*viper.Viper)
)

var DefaultConfig = Config{
	Env:  "development",
	Name: "config",
	Type: "yaml",
}

func NewWithModifier(cfg Config, modifiers ...Modifier) (*viper.Viper, error) {
	if cfg.Env == "" {
		cfg.Env = DefaultConfig.Env
	}

	if cfg.Name == "" {
		cfg.Name = DefaultConfig.Name
	}

	if cfg.Type == "" {
		cfg.Type = DefaultConfig.Type
	}

	env, err := environment.Parse(cfg.Env)
	if err != nil {
		return nil, err
	}

	if len(cfg.Paths) == 0 {
		cfg.Paths = DefaultConfig.Paths

		if env == environment.Production {
			cfg.Paths = append(cfg.Paths, ".", fmt.Sprintf("/etc/%s", cfg.ProjectName))
		} else {
			cfg.Paths = append(cfg.Paths, ".")
		}
	}

	configType, err := ParseType(cfg.Type)
	if err != nil {
		return nil, err
	}

	v := viper.New()

	v.SetConfigName(cfg.Name)
	v.SetConfigType(string(configType))

	for _, path := range cfg.Paths {
		v.AddConfigPath(path)
	}

	for _, modifier := range modifiers {
		modifier(v)
	}

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	return v, nil
}

func New(c ...Config) (*viper.Viper, error) {
	cfg := DefaultConfig

	if len(c) > 0 {
		cfg = c[0]
	}

	return NewWithModifier(cfg)
}
