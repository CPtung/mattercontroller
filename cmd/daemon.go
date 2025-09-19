package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/CPtung/mattercontroller/internal/openthread"
	"github.com/CPtung/mattercontroller/internal/server"
	"github.com/CPtung/mattercontroller/pkg/restapi"
	"github.com/spf13/cobra"
)

var (
	daemonCmd = &cobra.Command{
		Use:   "daemon",
		Short: "Run daemon",
		Run:   daemon,
	}
)

func init() {
	rootCmd.AddCommand(daemonCmd)
}

func daemon(cmd *cobra.Command, args []string) {
	// Initialize openthread network and border router
	ot := openthread.NewOpenThread()
	if err := ot.Initialize(cmd.Context()); err != nil {
		log.Panicf("Failed to initialize OpenThread: %v", err)
	}
	defer ot.Close()

	// Create a API server to handle REST API request/response
	apiServer := server.NewAPIServer()
	route := apiServer.Router()
	matterAPI := route.Group("/matter")
	{
		matterAPI.POST("/pairing", restapi.PostPairing)
		matterAPI.POST("/unpairing/:deviceID", restapi.PostUnpairing)
	}
	lightAPI := route.Group("/light")
	{
		lightAPI.GET("/:deviceID", restapi.GetLightState)
		lightAPI.PUT("/:deviceID", restapi.PutLightOnOff)
	}
	if err := apiServer.Start(); err != nil {
		log.Panicf("Failed to start API server: %v", err)
	}
	defer apiServer.Stop()

	// Declare a channel to receive OS signals for graceful shutdown
	chanSignal := make(chan os.Signal, 1)
	signal.Notify(chanSignal, syscall.SIGINT, syscall.SIGTERM)
	// Wait for interrupt/terminate system signal
	<-chanSignal
}
