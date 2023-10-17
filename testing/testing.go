package testing

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/nano-interactive/go-utils"
	"github.com/nano-interactive/go-utils/config"
	"github.com/spf13/viper"
)

type AppCreater[TServer, TContainer any] interface {
	Create(context.Context, *viper.Viper) (TServer, TContainer)
}

type AppCreaterFunc[TServer, TContainer any] func(context.Context, *viper.Viper) (TServer, TContainer)

func (h AppCreaterFunc[TServer, TContainer]) Create(ctx context.Context, config *viper.Viper) (TServer, TContainer) {
	return h(ctx, config)
}

func CreateApplicationFunc[TServer, TContainer any](t testing.TB, creater AppCreaterFunc[TServer, TContainer]) (TServer, TContainer) {
	return CreateApplication[TServer, TContainer](t, creater)
}

func CreateApplication[TServer, TContainer any](t testing.TB, creater AppCreater[TServer, TContainer], configName ...string) (TServer, TContainer) {
	configPath := FindConfig(t, configName...)

	cfg := viper.New()

	cfg.SetConfigName("config")
	cfg.AddConfigPath(configPath)
	cfg.SetConfigType("yaml")

	if err := cfg.ReadInConfig(); err != nil {
		panic(err)
	}

	return creater.Create(context.Background(), cfg)
}

func FindFile(t testing.TB, fileName string) string {
	t.Helper()
	workingDir := WorkingDir(t)
	root := ProjectRootDir(t)

	for entries, err := os.ReadDir(workingDir); err == nil; {
		for _, entry := range entries {
			if entry.Name() == fileName {
				return workingDir
			}
		}

		if workingDir == root {
			t.Error("got to he project root, file not found")
			t.FailNow()
		}

		workingDir, err = utils.GetAbsolutePath(filepath.Join(workingDir, ".."))

		if err != nil {
			t.Errorf("failed to get absolute path from %s", filepath.Join(workingDir, ".."))
			t.FailNow()
		}

		entries, err = os.ReadDir(workingDir)
	}

	t.Errorf("failed to find file %s", fileName)
	t.FailNow()
	return ""
}

func ProjectRootDir(t testing.TB) string {
	t.Helper()
	workingDir := WorkingDir(t)

	const gomod = "go.mod"

	for entries, err := os.ReadDir(workingDir); err == nil; {
		for _, entry := range entries {
			if entry.Name() == gomod {
				return workingDir
			}
		}

		if workingDir == "/" {
			t.Error("got to FS Root, file not found")
			t.FailNow()
		}

		workingDir, err = utils.GetAbsolutePath(filepath.Join(workingDir, ".."))

		if err != nil {
			t.Errorf("failed to get absolute path from %s", filepath.Join(workingDir, ".."))
			t.FailNow()
		}

		entries, err = os.ReadDir(workingDir)
	}

	t.Errorf("%s not found", gomod)
	t.FailNow()
	return ""
}

func FindConfig(t testing.TB, configName ...string) string {
	t.Helper()
	cfgName := "config.yml"

	if len(configName) > 0 {
		cfgName = configName[0]
	}

	return FindFile(t, cfgName)
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
	configPath := FindConfig(t)

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
