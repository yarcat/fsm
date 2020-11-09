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
	states := fsm.States{
		stInit: printing{string(stInit)},
		stWaitTimeout: fsm.Compose(
			printing{string(stWaitTimeout)},
			fsm.NewExpiring(p, fsm.After(time.Second), evTimeout),
		),
		stFinal: fsm.Compose(
			printing{string(stFinal)},
			quitting{},
		),
	}
	fsm := fsm.New(stInit, transitions, states, nil)
	p.Set(fsm)
	fsm.Send(evInitialized)
	select {} // Infinite sleep.

	// Output:
	// ENTER: stInit
	// LEAVE: stInit
	// ENTER: stWaitTimeout
	// LEAVE: stWaitTimeout
	// ENTER: stFinal
}
