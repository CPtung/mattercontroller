package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	Version, BuildTime string

	infoCmd = &cobra.Command{
		Use:   "info",
		Short: "Show version/build info",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Version:    %s\n", Version)
			fmt.Printf("Build Time: %s\n", BuildTime)
		},
	}

	rootCmd = &cobra.Command{
		Use:   "matter",
		Short: "Matter Controller",
	}
)

func init() {
	rootCmd.AddCommand(infoCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
