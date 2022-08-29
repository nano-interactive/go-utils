package testing

import (
	"context"
	"errors"
	"os"
	"path/filepath"

	"github.com/spf13/viper"

	utils "github.com/nano-interactive/go-utils"
)

type AppCreater[TServer, TContainer any] interface {
	Create(context.Context, *viper.Viper) (TServer, TContainer)
}

type AppCreaterFunc[TServer, TContainer any] func(context.Context, *viper.Viper) (TServer, TContainer)

func (h AppCreaterFunc[TServer, TContainer]) Create(ctx context.Context, config *viper.Viper) (TServer, TContainer) {
	return h(ctx, config)
}

func CreateApplicationFunc[TServer, TContainer any](creater AppCreaterFunc[TServer, TContainer]) (TServer, TContainer) {
	return CreateApplication[TServer, TContainer](creater)
}

func CreateApplication[TServer, TContainer any](creater AppCreater[TServer, TContainer]) (TServer, TContainer) {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	configPath, err := findConfig(wd)
	if err != nil {
		panic(err)
	}

	cfg := viper.New()

	cfg.SetConfigName("config")
	cfg.AddConfigPath(configPath)
	cfg.SetConfigType("yaml")

	if err := cfg.ReadInConfig(); err != nil {
		panic(err)
	}

	return creater.Create(context.Background(), cfg)
}

func findConfig(workingDir string, configName ...string) (string, error) {
	cfgName := "config.yml"

	if len(configName) == 0 {
		cfgName = configName[0]
	}

	for entries, err := os.ReadDir(workingDir); err == nil; {
		for _, entry := range entries {
			if !entry.IsDir() && entry.Name() == cfgName {
				return workingDir, nil
			}
		}

		workingDir, err = utils.GetAbsolutePath(filepath.Join(workingDir, ".."))

		if err != nil {
			return "", err
		}

		entries, err = os.ReadDir(workingDir)
	}

	return "", errors.New("config file not found")
}
