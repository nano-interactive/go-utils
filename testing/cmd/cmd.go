package cmd

import (
	"context"
	"io"
	"os"
	"testing"
	"time"

	"github.com/spf13/cobra"
)

type (
	Args struct {
		In               io.Reader
		Out              io.Writer
		Args             map[string]string
		ExecutingCommand string
		Commands         []*cobra.Command
	}

	Command struct {
		cancel context.CancelFunc
		done   <-chan struct{}
	}
)

func (c Command) Wait() {
	c.cancel()
	<-c.done
}

func StartCommand[T any](t testing.TB, ctx context.Context, root *cobra.Command, args ...Args) Command {
	ctx, cancel := context.WithCancel(ctx)
	t.Cleanup(func() {
		cancel()
	})

	defaultArgs := Args{
		Args:     make(map[string]string),
		Commands: []*cobra.Command{},
		In:       os.Stdin,
		Out:      os.Stdout,
	}

	if len(args) > 0 {
		defaultArgs = args[0]
		if defaultArgs.In == nil {
			defaultArgs.In = os.Stdin
		}

		if defaultArgs.Out == nil {
			defaultArgs.Out = os.Stdout
		}

		if len(defaultArgs.Args) == 0 {
			defaultArgs.Args = make(map[string]string)
		}
	}

	arg := make([]string, 0)

	for k, v := range defaultArgs.Args {
		arg = append(arg, "--"+k, v)
	}

	arg = append(arg, defaultArgs.ExecutingCommand)

	root.SetArgs(arg)
	root.SetIn(defaultArgs.In)
	root.SetOut(defaultArgs.Out)

	done := make(chan struct{}, 1)
	t.Cleanup(func() {
		close(done)
	})

	go func(ctx context.Context) {
		if err := root.ExecuteContext(ctx); err != nil {
			t.Error(err)
			t.Fail()
		}

		done <- struct{}{}
	}(ctx)

	return Command{
		cancel: cancel,
		done:   done,
	}
}

type Result int

const (
	Done Result = iota
	Timeout
)

func StartCommandWithWait[T any](t testing.TB, ctx context.Context, root *cobra.Command, wait time.Duration, args ...Args) Result {
	c := StartCommand[T](t, ctx, root, args...)
	timer := time.NewTimer(wait)
	defer c.cancel()
	defer timer.Stop()

	select {
	case <-c.done:
		return Done
	case <-timer.C:
		return Timeout
	}
}
