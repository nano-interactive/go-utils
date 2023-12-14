package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
)

type ContextKey string

const (
	ConfigContextKey ContextKey = "go-utils-cmd-config-key"
	CancelContextKey ContextKey = "go-utils-cmd-cancel-key"
)

func PersistentPreRunE[T any](cfgFn func() (T, error), persistentPreRunE func(context.Context, T, *cobra.Command, []string) error) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		cfgVal := ctx.Value(ConfigContextKey)
		cfgVal2 := ctx.Value(string(ConfigContextKey))

		if cfgVal2 != nil {
			cfgVal = cfgVal2
		}

		if cfgVal == nil {
			cfg, err := cfgFn()
			if err != nil {
				return err
			}

			cfgVal = cfg
		}

		ctx = context.WithValue(context.Background(), ConfigContextKey, cfgVal)
		ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
		ctx = context.WithValue(ctx, CancelContextKey, cancel)
		cmd.SetContext(ctx)
		return persistentPreRunE(ctx, cfgVal.(T), cmd, args)
	}
}

func PersistentPostRunE(persistentPostRunE func(*cobra.Command, []string) error) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		cancel := ctx.Value(CancelContextKey).(context.CancelFunc)
		cancel()
		return persistentPostRunE(cmd, args)
	}
}
