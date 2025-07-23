package cmd

import (
	"strings"

	"github.com/spf13/cobra"
)

const logDebugFlag = "debug"
const logDebugDef = "info+:*"

func addLogDebugFlag(cmd *cobra.Command) {
	cmd.Root().PersistentFlags().String(
		logDebugFlag,
		"error",
		`Display detailed log information at the debug level`,
	)
}

func getLogDebugFlag(cmd *cobra.Command) int {
	if result, ok := cmd.Root().PersistentFlags().GetString(logDebugFlag); ok == nil {
		if result != "" {
			return parseDebugLevel(result)
		}
	}

	return LevelErrorLevel
}

const (
	LevelDebug = "debug"
	LevelInfo  = "info"
	LevelWarn  = "warn"
	LevelError = "error"

	LevelDebugLevel = -4
	LevelInfoLevel  = 0
	LevelWarnLevel  = 4
	LevelErrorLevel = 8
)

func parseDebugLevel(result string) int {
	switch strings.ToLower(result) {
	case LevelDebug:
		return LevelDebugLevel
	case LevelInfo:
		return LevelInfoLevel
	case LevelWarn:
		return LevelWarnLevel
	case LevelError:
		return LevelErrorLevel
	}

	return LevelErrorLevel
}
