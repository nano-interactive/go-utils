package testing

import (
	"context"
	"errors"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	"github.com/valyala/fasthttp"

	utils "github.com/nano-interactive/go-utils"
)

type AppCreater[T any] interface {
	Create(context.Context, *viper.Viper) (*fasthttp.Server, T)
}

func CreateApplication[T any](creater AppCreater[T]) (*fasthttp.Server, T) {
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
