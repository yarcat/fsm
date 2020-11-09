package fsm_test

import (
	"fmt"
	"os"
	"time"

	"github.com/yarcat/fsm-go"
)

const (
	stInit        fsm.StateType = "stInit"
	stWaitTimeout fsm.StateType = "stWaitTimeout"
	stFinal       fsm.StateType = "stFinal"

	evInitialized fsm.EventType = "evInitialized"
	evTimeout     fsm.EventType = "evTimeout"
)

var transitions = fsm.Transitions{
	fsm.When(stInit, evInitialized):    stWaitTimeout,
	fsm.When(stWaitTimeout, evTimeout): stFinal,
}

type (
	printing struct{ name string }
	quitting struct{}
)

func (st printing) Enter() { fmt.Println("ENTER:", st.name) }
func (st printing) Leave() { fmt.Println("LEAVE:", st.name) }

func (quitting) Enter() { os.Exit(0) }
func (quitting) Leave() {}

func Example() {
	p := new(fsm.MachineProvider)
	stateName := func(st fsm.StateType) string {
		return fmt.Sprintf("%s (sync)", st)
	}
	states := fsm.States{
		stInit: printing{stateName(stInit)},
		stWaitTimeout: fsm.Compose(
			printing{stateName(stWaitTimeout)},
			fsm.NewExpiring(p, fsm.After(100*time.Millisecond), evTimeout),
		),
		stFinal: fsm.Compose(
			printing{stateName(stFinal)},
			quitting{},
		),
	}
	fsm := fsm.New(stInit, transitions, states, nil)
	p.Set(fsm)
	fsm.Send(evInitialized)
	select {} // Infinite sleep.

	// Output:
	// ENTER: stInit (sync)
	// LEAVE: stInit (sync)
	// ENTER: stWaitTimeout (sync)
	// LEAVE: stWaitTimeout (sync)
	// ENTER: stFinal (sync)
}

func ExampleNewAsync() {
	p := new(fsm.MachineProvider)
	stateName := func(st fsm.StateType) string {
		return fmt.Sprintf("%s (async)", st)
	}
	states := fsm.States{
		stInit: printing{stateName(stInit)},
		stWaitTimeout: fsm.Compose(
			printing{stateName(stWaitTimeout)},
			fsm.NewExpiring(p, fsm.After(100*time.Millisecond), evTimeout),
		),
		stFinal: fsm.Compose(
			printing{stateName(stFinal)},
			quitting{},
		),
	}
	fsm := fsm.NewAsync(stInit, transitions, states, nil)
	p.Set(fsm)
	fsm.Send(evInitialized)
	fsm.Run() // Infinite loop.

	// Output:
	// ENTER: stInit (async)
	// LEAVE: stInit (async)
	// ENTER: stWaitTimeout (async)
	// LEAVE: stWaitTimeout (async)
	// ENTER: stFinal (async)
}
