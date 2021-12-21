package main

import (
	"bb.yandex-team.ru/cloud/cloud-go/healthcheck/command"
	diffimpl "bb.yandex-team.ru/cloud/cloud-go/healthcheck/command/ylb-duty"

	"github.com/spf13/cobra"
)

func main() {
	cmd := &cobra.Command{
		Use:  "ylb-duty",
		Long: "ylb-duty-tools",
	}
	command.SetupConfig(cmd)

	cmd.AddCommand(
		command.Cobra("diff", diffimpl.New()),
	)

	command.Run(cmd)
}
