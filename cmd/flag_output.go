package cmd

import (
	"github.com/spf13/cobra"
)

const outputFlag = "output"
const helpFormatter = "help"
const defaultFormatter = "yaml"

type OutputFormatter interface {
	Format(value any, options string, isRaw bool) error
	Description() string
}

var outputFormatters = map[string]OutputFormatter{}

func addOutputFlag(cmd *cobra.Command) {
	cmd.Root().PersistentFlags().StringP(
		outputFlag,
		"o",
		"",
		`Change the output format. Use '--output=help' to know more details.`)
}

func getOutputFlag(cmd *cobra.Command) string {
	return cmd.Root().PersistentFlags().Lookup(outputFlag).Value.String()
}
