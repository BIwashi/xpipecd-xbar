package cli

import (
	"context"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/cockroachdb/errors"
	"github.com/spf13/cobra"
)

type CLI struct {
	rootCmd *cobra.Command
	flags   PersistentFlags
}

type Input struct {
	Logger          slog.Logger
	PersistentFlags PersistentFlags
	Stdin           io.Reader
}

type PersistentFlags struct {
	LogLevel string
}

var defaultPersistentFlags = PersistentFlags{
	LogLevel: "info",
}

func NewCLI(name, desc string) *CLI {
	c := &CLI{
		rootCmd: &cobra.Command{
			Use:           name,
			Short:         desc,
			SilenceErrors: true,
		},
		flags: defaultPersistentFlags,
	}

	c.setGlobalFlags()

	return c
}

func (c *CLI) AddCommands(cmds ...*cobra.Command) {
	for _, cmd := range cmds {
		c.rootCmd.AddCommand(cmd)
	}
}

func (c *CLI) Run() error {
	if err := c.rootCmd.Execute(); err != nil {
		return errors.Wrap(err, "execute command")
	}

	return nil
}

func (c *CLI) setGlobalFlags() {
	c.rootCmd.PersistentFlags().StringVar(&c.flags.LogLevel, "log-level", c.flags.LogLevel,
		"Log level. Available values: debug, info, warn, error. Default is info.",
	)
}

type Runner func(ctx context.Context, input Input) error

func WithContext(runner Runner) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		defer signal.Stop(ch)

		return runWithContext(cmd, ch, runner)
	}
}

func runWithContext(cmd *cobra.Command, signalCh <-chan os.Signal, runner Runner) error {
	flags, err := parsePersistentFlags(cmd)
	if err != nil {
		return err
	}

	input := Input{
		PersistentFlags: flags,
		Stdin:           cmd.InOrStdin(),
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout,
		&slog.HandlerOptions{
			Level: slogLevelFromString(flags.LogLevel),
		}),
	)

	input.Logger = *logger

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		<-signalCh
		cancel()
	}()

	if err := runner(ctx, input); err != nil {
		return errors.Wrap(err, "run command")
	}

	return nil
}

func parsePersistentFlags(cmd *cobra.Command) (PersistentFlags, error) {
	fs := cmd.Flags()
	flags := defaultPersistentFlags

	if fs.Lookup("log-level") != nil {
		s, err := fs.GetString("log-level")
		if err != nil {
			return flags, err
		}
		flags.LogLevel = s
	}

	return flags, nil
}
