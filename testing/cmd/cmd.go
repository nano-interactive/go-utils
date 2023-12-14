package cmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"testing"
	"time"

	"github.com/nano-interactive/go-utils/v2/cmd"
	"github.com/spf13/cobra"
)

type (
	Args struct {
		In           io.Reader
		Out          io.Writer
		Args         map[string]any
		ContextItems map[string]any
	}

	Command struct {
		cancel context.CancelFunc
		done   <-chan *cobra.Command
	}

	Result int

	Option func(options *Options)

	Options struct {
		timeout time.Duration
		args    Args
	}
)

const (
	Done Result = iota
	Timeout
)

func startCommand(t testing.TB, ctx context.Context, root *cobra.Command, args Args) Command {
	ctx, cancel := context.WithCancel(ctx)
	t.Cleanup(func() {
		cancel()
	})

	if args.In == nil {
		args.In = os.Stdin
	}

	if args.Out == nil {
		args.Out = os.Stdout
	}

	if len(args.Args) == 0 {
		args.Args = make(map[string]any)
	}

	if len(args.ContextItems) == 0 {
		args.ContextItems = make(map[string]any)
	}

	arg := make([]string, 0, len(args.Args))

	for k, v := range args.Args {
		switch data := v.(type) {
		case string:
			arg = append(arg, k, data)
		case []string:
			for _, i := range data {
				arg = append(arg, k, i)
			}
		case fmt.Stringer:
			arg = append(arg, k, data.String())
		case nil:
			arg = append(arg, k)
		default:
			t.Fatalf("Failed to append to args %T", v)
		}
	}

	for k, item := range args.ContextItems {
		ctx = context.WithValue(ctx, k, item)
	}

	root.SetArgs(arg)
	root.SetIn(args.In)
	root.SetOut(args.Out)

	done := make(chan *cobra.Command, 1)
	t.Cleanup(func() {
		close(done)
	})

	go func(ctx context.Context) {
		c, err := root.ExecuteContextC(ctx)
		if err != nil {
			t.Error(err)
			t.Fail()
		}

		done <- c
	}(ctx)

	return Command{
		cancel: cancel,
		done:   done,
	}
}

type ()

func WithTimeout(wait time.Duration) Option {
	return func(options *Options) {
		options.timeout = wait
	}
}

func WithArgs(args map[string]any) Option {
	return func(options *Options) {
		options.args.Args = args
	}
}

func WithInput(in io.Reader) Option {
	return func(options *Options) {
		options.args.In = in
	}
}

func WithOutput(in io.Writer) Option {
	return func(options *Options) {
		options.args.Out = in
	}
}

func WithConfig(key string, cfg any) Option {
	return func(options *Options) {
		options.args.ContextItems[key] = cfg
	}
}

func WithConfigDefault(cfg any) Option {
	return func(options *Options) {
		options.args.ContextItems[string(cmd.ConfigContextKey)] = cfg
	}
}

func WithContextItems(items map[string]any) Option {
	return func(options *Options) {
		for key, value := range items {
			options.args.ContextItems[key] = value
		}
	}
}

func StartCommand(tb testing.TB, ctx context.Context, root *cobra.Command, opts ...Option) Result {
	opt := Options{
		args: Args{
			In:           os.Stdin,
			Out:          os.Stdout,
			Args:         make(map[string]any),
			ContextItems: make(map[string]any),
		},
		timeout: 0,
	}

	for _, o := range opts {
		o(&opt)
	}

	c := startCommand(tb, ctx, root, opt.args)
	defer c.cancel()

	if opt.timeout > 0 {
		timer := time.NewTimer(opt.timeout)
		defer timer.Stop()
		select {
		case <-c.done:
			return Done
		case <-timer.C:
			return Timeout
		}
	}

	<-c.done
	return Done
}
