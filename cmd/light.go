package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/CPtung/mattercontroller/internal/matter/lighting"
	"github.com/spf13/cobra"
)

var (
	lightCmd = &cobra.Command{
		Use:   "light <interface>",
		Short: "Lighting App",
		Args:  cobra.ExactArgs(1),
		Run:   light,
	}
)

func init() {
	rootCmd.AddCommand(lightCmd)
}

func light(cmd *cobra.Command, args []string) {
	iface := args[0]
	// Declare a channel to receive OS signals for graceful shutdown
	chanSignal := make(chan os.Signal, 1)
	signal.Notify(chanSignal, syscall.SIGINT, syscall.SIGTERM)

	app := lighting.NewApp("", iface)
	if err := app.Start(cmd.Context()); err != nil {
		panic(err)
	}
	defer app.Stop()
	<-chanSignal
}
