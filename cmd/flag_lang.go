package cmd

import "github.com/spf13/cobra"

func addLangFlag(cmd *cobra.Command) {
	cmd.PersistentFlags().String(
		"lang",
		"en-US",
		"Set the language for the CLI",
	)
}
