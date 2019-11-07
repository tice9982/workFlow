package workflow

import (
	"context"
	"fmt"
	"path/filepath"
	"runtime"
)

type FlowType int

const (
	FlowType_Single   FlowType = 0 // 单一
	FlowType_Serial   FlowType = 1 // 串行
	FlowType_Parallel FlowType = 2 // 并行
)

type funEvent func(parentCtx context.Context, ch chan<- byte) error

type Event struct {
	Name     string
	Function funEvent
}

var Enter_event = Event{Name: "enter"}
var Run_event = Event{Name: "run"}
var Quit_event = Event{Name: "quit"}

type Dispatcher interface {
	Run(parentCtx context.Context, parentChan chan<- byte) // run work flow
	PushEvent(e Event, f funEvent) error
	AppendChild(wfd Dispatcher) error
	Cancel()
	GetUserData() interface{}
}

func TraceLog(content string) {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		return
	}
	fmt.Println(fmt.Sprintf("%s %d>>> \t%s", filepath.Base(file), line, content))
}

//////////
type SingleFlow struct {
	Dispatcher

	cancel   context.CancelFunc
	UserData interface{}
	Events   []Event
}

func (wf *SingleFlow) Cancel() {
	if wf.cancel != nil {
		wf.cancel()
	}
}

func (wf *SingleFlow) PushEvent(e Event, f funEvent) error {
	// if _, ok := wf.Events[e]; ok {
	// 	return fmt.Errorf(fmt.Sprintf("single flow event: %s exist", e.Name))
	// }
	// wf.Events[e]
	e.Function = f
	wf.Events = append(wf.Events, e)
	return nil
}

func (wf *SingleFlow) AppendChild(wfd Dispatcher) error {
	return fmt.Errorf("single flow can not append child")
}

func (wf *SingleFlow) GetUserData() interface{} {
	return wf.UserData
}

func (wf *SingleFlow) Run(parentCtx context.Context, parentChan chan<- byte) {
	defer func() {
		parentChan <- 0
	}()

	ctx, cancel := context.WithCancel(parentCtx)
	wf.cancel = cancel

	for _, e := range wf.Events {
		ch := make(chan byte)
		go e.Function(ctx, ch)
	loop:
		for {
			select {
			case <-ctx.Done():
				TraceLog(wf.GetUserData().(string) + " cancel")
				return
			case <-ch:
				TraceLog(wf.GetUserData().(string) + " done")
				break loop
			default:
			}
		}

	}
}

////////////
type SerialFlow struct {
	Dispatcher

	cancel   context.CancelFunc
	UserData interface{}
	children []Dispatcher
}

func (wf *SerialFlow) Cancel() {
	if wf.cancel != nil {
		wf.cancel()
	}
}

func (wf *SerialFlow) PushEvent(e Event, f funEvent) error {
	return fmt.Errorf("serial flow can not push event")
}

func (wf *SerialFlow) AppendChild(wfd Dispatcher) error {
	wf.children = append(wf.children, wfd)
	return nil
}

func (wf *SerialFlow) GetUserData() interface{} {
	return wf.UserData
}

func (wf *SerialFlow) Run(parentCtx context.Context, parentChan chan<- byte) {
	defer func() {
		parentChan <- 0
	}()

	ctx, cancel := context.WithCancel(parentCtx)
	wf.cancel = cancel

	childrenCount := len(wf.children)
	TraceLog(fmt.Sprintf("serial flow! children count:%d", childrenCount))
	flag := 0

	for {
		select {
		case <-ctx.Done():
			TraceLog(wf.GetUserData().(string) + " cancel")
			return
		default:
			if childrenCount > flag {
				ch := make(chan byte)
				cwf := wf.children[flag]
				go cwf.Run(ctx, ch)
				<-ch
				flag++
			} else {
				TraceLog(wf.GetUserData().(string) + " done")
				return
			}
		}
	}
}

////////////
type ParallelFlow struct {
	Dispatcher

	cancel   context.CancelFunc
	UserData interface{}
	children []Dispatcher
}

func (wf *ParallelFlow) Cancel() {
	if wf.cancel != nil {
		wf.cancel()
	}
}

func (wf *ParallelFlow) PushEvent(e Event, f funEvent) error {
	return fmt.Errorf("parallel flow can not push event")
}

func (wf *ParallelFlow) AppendChild(wfd Dispatcher) error {
	wf.children = append(wf.children, wfd)
	return nil
}

func (wf *ParallelFlow) GetUserData() interface{} {
	return wf.UserData
}

func (wf *ParallelFlow) Run(parentCtx context.Context, parentChan chan<- byte) {
	defer func() {
		parentChan <- 0
	}()

	ctx, cancel := context.WithCancel(parentCtx)
	wf.cancel = cancel
	ch := make(chan byte)

	count := len(wf.children)
	TraceLog(fmt.Sprintf("parallel flow! children count:%d", count))
	for _, cwf := range wf.children {
		go cwf.Run(ctx, ch)
	}

	for {
		select {
		case <-ctx.Done():
			TraceLog(wf.GetUserData().(string) + " cancel")
			return
		case <-ch:
			count--
			if count <= 0 {
				TraceLog(wf.GetUserData().(string) + " done")
				return
			}
		default:
			for i := 0; i < len(wf.children); i++ {
				<-ch
			}
		}
	}
}

///////////////////////
func CreateSingleFlow(userData interface{}) Dispatcher {
	return &SingleFlow{UserData: userData}
}

func CreateSerialFlow(userData interface{}) Dispatcher {
	return &SerialFlow{UserData: userData, children: []Dispatcher{}}
}

func CreateParallelFlow(userData interface{}) Dispatcher {
	return &ParallelFlow{UserData: userData, children: []Dispatcher{}}
}
