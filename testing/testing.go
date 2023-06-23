package testing

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/nano-interactive/go-utils"
	"github.com/nano-interactive/go-utils/config"
	"github.com/spf13/viper"
)

// interface AppCreater
// Retunrs Server and Container for testing purposes
type AppCreater[TServer, TContainer any] interface {
	Create(context.Context, *viper.Viper) (TServer, TContainer)
}

// Delegates functions for AppCreater
type AppCreaterFunc[TServer, TContainer any] func(context.Context, *viper.Viper) (TServer, TContainer)

// TODO: Check if it's needed to document this
func (h AppCreaterFunc[TServer, TContainer]) Create(ctx context.Context, config *viper.Viper) (TServer, TContainer) {
	return h(ctx, config)
}

// TODO: Check if it's needed to document this
func CreateApplicationFunc[TServer, TContainer any](creater AppCreaterFunc[TServer, TContainer]) (TServer, TContainer) {
	return CreateApplication[TServer, TContainer](creater)
}

// Creates a new instance of application for testing purposes
func CreateApplication[TServer, TContainer any](creater AppCreater[TServer, TContainer], configName ...string) (TServer, TContainer) {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	configPath, err := FindConfig(wd, configName...)
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

// Returns directory of config files and error if file doesent exist
func FindConfig(workingDir string, configName ...string) (string, error) {
	cfgName := "config.yml"

	if len(configName) > 0 {
		cfgName = configName[0]
	}

	return FindFile(workingDir, cfgName)
}

// Returns directory of config file name and error if file doesent exist
func FindFile(workingDir string, fileName string) (string, error) {
	for entries, err := os.ReadDir(workingDir); err == nil; {
		for _, entry := range entries {
			if entry.Name() == fileName {
				return workingDir, nil
			}
		}

		workingDir, err = utils.GetAbsolutePath(filepath.Join(workingDir, ".."))

		if err != nil {
			return "", err
		}

		entries, err = os.ReadDir(workingDir)
	}

	return "", errors.New("file or directory not found")
}

func WorkingDir(t testing.TB) string {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	return wd
}

func GetConfig[T any](t testing.TB, create func(*viper.Viper) (T, error)) T {
	t.Helper()
	configPath, err := FindConfig(WorkingDir(t))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	viperCfg, err := config.New(config.Config{
		Env:   "testing",
		Name:  "config",
		Type:  "yaml",
		Paths: []string{configPath},
	})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	cfg, err := create(viperCfg)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	return cfg
}


func Timeout(t testing.TB, timeout time.Duration) context.Context {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)

	t.Cleanup(func() {
		cancel()
	})

	return ctx
}