package openthread

import (
	"context"

	"github.com/CPtung/mattercontroller/internal/openthread/simulation"
)

type OpenThread interface {
	Initialize(context.Context) error
	Close()
}

type OpenThreadImpl struct{}

func NewOpenThread() OpenThread {
	return &OpenThreadImpl{}
}

func (ot *OpenThreadImpl) Initialize(ctx context.Context) error {
	// create rcp simulator
	ttyNum, err := simulation.StartRCP()
	if err != nil {
		return err
	}

	// create otbr simulator
	if err := createBorderRouter(ttyNum); err != nil {
		defer simulation.StopRCP()
		return err
	}

	// start OTBR
	if err := setupThread(); err != nil {
		defer closeBorderRouter()
		defer simulation.StopRCP()
		return err
	}

	return nil
}

/*
	func (ot *OpenThreadImpl) Initialize(ctx context.Context) error {
		ot.Add(1)
		chanErr := make(chan error, 1)
		go func() {
			// create rcp simulator
			ttyNum, err := simulation.StartRCP()
			if err != nil {
				chanErr <- err
				return
			}
			defer simulation.StopRCP()

			// create otbr simulator
			if err := createBorderRouter(ttyNum); err != nil {
				chanErr <- err
				return
			}
			defer closeBorderRouter()

			// start OTBR
			if err := setupThread(); err != nil {
				chanErr <- err
				return
			}
			defer tearDownThread()

			// Return initialization result first
			chanErr <- nil
			// Wait for context done
			<-ot.ctx.Done()
			// Leave OTBR
			ot.Done()
		}()
		return <-chanErr
	}
*/
func (ot *OpenThreadImpl) Close() {
	tearDownThread()
	closeBorderRouter()
	simulation.StopRCP()
}
