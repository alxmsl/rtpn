package main

import (
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

const (
	exitCodeGeneral     = 1
	exitCodeWrongPlugin = 2
	exitCodeWrongRun    = 3
	exitCodeRunError    = 4
)

var (
	rootCmd = &cobra.Command{
		Use:   "cpn",
		Short: "Tool to manage Petri-Net.",
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		slog.Error("error on run command", "error", err)
		os.Exit(exitCodeGeneral)
		return
	}
}
