package main

import (
	"log/slog"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

type listFlags struct {
	Dir       string
	NeedCheck bool
}

var (
	listCmdFlags listFlags

	listCmd = &cobra.Command{
		Use:   "list",
		Short: "List available Petri-Net plugins",
		Run: func(cmd *cobra.Command, args []string) {
			entries, err := os.ReadDir(listCmdFlags.Dir)
			if err != nil {
				slog.Error("error on read dir", "dir", listCmdFlags.Dir, "error", err)
				os.Exit(exitCodeGeneral)
			}

			var pluginCounter int
			for _, e := range entries {
				if strings.HasSuffix(e.Name(), ".so") {
					pluginCounter += 1
					slog.Info("found plugin", "name", e.Name())
				}
			}

			//@todo: add check for plugin file

			slog.Info("total plugins", "dir", listCmdFlags.Dir, "value", pluginCounter)
		},
	}
)

func init() {
	listCmd.Flags().StringVar(&listCmdFlags.Dir, "dir", ".", "Directory which contains Petri-Net plugins")
	listCmd.Flags().BoolVar(&listCmdFlags.NeedCheck, "check", false, "Check Petri-Net plugins")

	rootCmd.AddCommand(listCmd)
}
