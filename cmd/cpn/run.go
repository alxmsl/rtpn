package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/alxmsl/cpn"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

const (
	flagDryRun   = "dry-run"
	flagFilename = "filename"
)

type runFlags struct {
	IsDryRun bool
	Filename string
}

var (
	runCmdFlags runFlags

	runCmd = &cobra.Command{
		Use:   "run [flags] [placement value...]",
		Short: "Run Petri-Net plugin with a given markup",
		Run: func(cmd *cobra.Command, args []string) {
			var net, err = lookup(runCmdFlags.Filename)
			if err != nil {
				slog.Error("error on open Petri-Net plugin", "name", runCmdFlags.Filename, "error", err)
				os.Exit(exitCodeWrongPlugin)
				return
			}
			slog.Info("Petri-Net symbol has been found", "filename", runCmdFlags.Filename)
			if runCmdFlags.IsDryRun {
				return
			}

			if len(args)%2 != 0 {
				slog.Error("number of placement values does not match")
				os.Exit(exitCodeWrongRun)
				return
			}

			var (
				ctx, cancel = context.WithTimeout(context.Background(), time.Second)
				eg          = errgroup.Group{}
			)
			defer cancel()
			eg.Go(func() error {
				return net.Run(ctx)
			})
			eg.Go(func() error {
				for i := 0; i < len(args); i += 2 {
					var (
						token = *cpn.NewToken(args[i+1])
						err   = net.Markup(args[i], token)
					)
					if err != nil {
						return err
					}
				}
				return nil
			})
			err = eg.Wait()
			if err != nil {
				slog.Error("error on run Petri-Net plugin", "name", runCmdFlags.Filename, "error", err)
				os.Exit(exitCodeRunError)
				return
			}
			slog.Info("Petri-Net has been completed", "filename", runCmdFlags.Filename)
		},
	}
)

func init() {
	runCmd.Flags().BoolVar(&runCmdFlags.IsDryRun, flagDryRun, true, "Load symbol from a Petri-Net  plugin only.")
	runCmd.Flags().StringVar(&runCmdFlags.Filename, flagFilename, "", "Petri-Net plugin filename.")

	_ = runCmd.MarkFlagRequired(flagFilename)

	rootCmd.AddCommand(runCmd)
}
