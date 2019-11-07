package workflow

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestSingleFlow(t *testing.T) {
	d := CreateSingleFlow("s")

	d.PushEvent(Run_event, func(parentCtx context.Context, parentChan chan<- byte) error {
		TraceLog(d.GetUserData().(string) + " on run ")
		time.Sleep(4 * time.Second)

		parentChan <- 0
		return nil
	})

	ch := make(chan byte)
	go d.Run(context.Background(), ch)

	time.Sleep(2 * time.Second)

	d.Cancel()
}

func TestSerialFlow(t *testing.T) {
	serialF := CreateSerialFlow(1)
	for i := 0; i < 3; i++ {
		d := CreateSingleFlow(fmt.Sprintf("s %d", i))

		d.PushEvent(Enter_event, func(parentCtx context.Context, parentChan chan<- byte) error {
			defer func() { parentChan <- 0 }()
			TraceLog("serial " + d.GetUserData().(string) + " on enter ")
			return nil
		})

		d.PushEvent(Run_event, func(parentCtx context.Context, parentChan chan<- byte) error {
			defer func() { parentChan <- 0 }()
			TraceLog("serial " + d.GetUserData().(string) + " on run ")
			time.Sleep(2 * time.Second)
			return nil
		})
		d.PushEvent(Quit_event, func(parentCtx context.Context, parentChan chan<- byte) error {
			defer func() { parentChan <- 0 }()
			TraceLog("serial " + d.GetUserData().(string) + " on quit \n")
			return nil
		})

		serialF.AppendChild(d)
	}

	ch := make(chan byte)
	go serialF.Run(context.Background(), ch)
	time.Sleep(4 * time.Second)

	serialF.Cancel()
}

func TestParallelFlow(t *testing.T) {
	pF := CreateParallelFlow(1)
	for i := 0; i < 3; i++ {
		d := CreateSingleFlow(fmt.Sprintf("p %d", i))

		d.PushEvent(Enter_event, func(parentCtx context.Context, parentChan chan<- byte) error {
			defer func() { parentChan <- 0 }()
			TraceLog("parallel " + d.GetUserData().(string) + " on enter ")
			return nil
		})
		d.PushEvent(Run_event, func(parentCtx context.Context, parentChan chan<- byte) error {
			defer func() { parentChan <- 0 }()
			TraceLog("parallel " + d.GetUserData().(string) + " on run ")
			time.Sleep(2 * time.Second)
			parentChan <- 0
			return nil
		})
		d.PushEvent(Quit_event, func(parentCtx context.Context, parentChan chan<- byte) error {
			defer func() { parentChan <- 0 }()
			TraceLog("parallel " + d.GetUserData().(string) + " on quit \n")
			return nil
		})

		pF.AppendChild(d)
	}
	ch := make(chan byte)
	go pF.Run(context.Background(), ch)
	time.Sleep(3 * time.Second)

	pF.Cancel()

}
