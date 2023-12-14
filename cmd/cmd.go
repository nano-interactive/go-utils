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
		cfg, err := cfgFn()
		if err != nil {
			return err
		}

		ctx := context.WithValue(context.Background(), ConfigContextKey, cfg)
		ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
		ctx = context.WithValue(ctx, CancelContextKey, cancel)
		cmd.SetContext(ctx)
		return persistentPreRunE(ctx, cfg, cmd, args)
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
