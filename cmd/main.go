package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/CPtung/mtctrl/internal/config"
	"github.com/CPtung/mtctrl/pkg/otbr"
	"github.com/CPtung/mtctrl/pkg/server"
	"github.com/CPtung/mtctrl/pkg/simulation"
)

func main() {
	config.Init()
	// create rcp simulator
	ttyNum, err := simulation.StartRCP()
	if err != nil {
		panic(err)
	}
	defer simulation.StopRCP()

	// create otbr simulator
	otbr.CreateBorderRouter(ttyNum)
	defer otbr.CloseBorderRouter()

	// start OTBR
	otbr.SetupThreadNetwork()
	defer otbr.TearDownThreadNetwork()

	// Declare a channel to receive OS signals for graceful shutdown
	chanSignal := make(chan os.Signal, 1)
	signal.Notify(chanSignal, syscall.SIGINT, syscall.SIGTERM)

	apiServer := server.NewAPIServer()
	if err := apiServer.Start(); err != nil {
		panic(err)
	}
	defer apiServer.Stop()

	<-chanSignal
}
