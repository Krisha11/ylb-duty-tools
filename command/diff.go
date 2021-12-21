package ylb_duty

import (
	"context"
	"fmt"

	"bb.yandex-team.ru/cloud/cloud-go/healthcheck/command"
	common "bb.yandex-team.ru/cloud/cloud-go/healthcheck/ylb-duty"
	diffimpl "bb.yandex-team.ru/cloud/cloud-go/healthcheck/ylb-duty/diff"
	"bb.yandex-team.ru/cloud/cloud-go/lib/log/ctxlog"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// New creates new Command.
func New() *Command {
	return &Command{}
}

// Config contains command configuration parameters.
type Config struct {
	logLevel    string
	logJournald bool
	kinds       []string
}

// Command contains command info.
type Command struct {
	config Config
}

// Run runs command.
func (c *Command) Run(args []string) error {
	logger, err := command.InitializeLogger(c.config.logLevel, c.config.logJournald)
	if err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}
	ctx := context.Background()
	ctxlog.Info(ctx, logger, "logger")

	diffCollectorConfig := common.Config{Logger: logger}
	err = viper.UnmarshalKey("endpoint", &diffCollectorConfig.Endpoint)
	if err != nil {
		return fmt.Errorf("failed to extract endpoint from config: %w", err)
	}

	err = viper.UnmarshalKey("database", &diffCollectorConfig.Database)
	if err != nil {
		return fmt.Errorf("failed to extract database name from config: %w", err)
	}

	err = viper.UnmarshalKey("prefix", &diffCollectorConfig.Prefix)
	if err != nil {
		return fmt.Errorf("failed to extract tables prefix from config: %w", err)
	}

	differ, err := diffimpl.NewDiffCollector(&diffCollectorConfig)

	if err != nil {
		return err
	}

	danglings, err := differ.Diff(ctx, c.config.kinds)
	if err != nil {
		return err
	}

	for _, v := range danglings {
		if len(v) > 0 {
			return fmt.Errorf("inconsistent: %v", danglings)
		}
	}

	return nil
}

// Help returns command's help message.
func (c *Command) Help() string {
	return "Index consistency checker"
}

// Synopsis returns command's synopsis message.
func (c *Command) Synopsis() string {
	return "Index consistency checker"
}

// ExportFlags exports command's flags into given flag set.
func (c *Command) ExportFlags(flag *pflag.FlagSet) {
	flag.StringVar(&c.config.logLevel,
		"log-level", "info",
		"log-level: debug/info/error",
	)
	flag.BoolVar(&c.config.logJournald,
		"log-journal", false,
		"log-journal",
	)
	flag.StringSliceVarP(&c.config.kinds,
		"kinds", "k", []string{},
		"kinds of checking",
	)
}
