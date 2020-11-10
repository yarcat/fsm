package fsm

import (
	"log"
)

type (
	// FSM is a finite-state machine. It implements state transitions and events.
	FSM struct {
		current StateType

		transitions    Transitions
		states         States
		defaultHandler State
	}
	// AsyncFSM is an asynchronous version of FSM. It is go-routine safe.
	AsyncFSM struct {
		fsm           *FSM
		queue         chan []EventType
		done, updated chan struct{}
	}
)

// New returns new configured finite-state machine.
func New(startState StateType, transitions Transitions, states States, defaultHandler State) *FSM {
	if defaultHandler == nil {
		defaultHandler = DefaultHandler
	}
	fsm := &FSM{
		transitions:    transitions,
		states:         states,
		defaultHandler: defaultHandler,
	}
	fsm.handler(startState).Enter()
	fsm.current = startState
	log.Printf("fsm %p: initialized: state %v", fsm, startState)
	return fsm
}

// NewAsync returns new configured asynchronous finite-state machine.
func NewAsync(startState StateType, transitions Transitions, states States, defaultHandler State) *AsyncFSM {
	fsm := New(startState, transitions, states, defaultHandler)
	queue := make(chan []EventType, 1)
	queue <- nil
	return &AsyncFSM{
		fsm:     fsm,
		done:    make(chan struct{}),
		queue:   queue,
		updated: make(chan struct{}, 1),
	}
}

// Send sends the event to the finite-state machine. The machine changes its
// state accordingly to the transitions table.
func (fsm *FSM) Send(event EventType) {
	log.Printf("fsm %p: state %v: received %v ", fsm, fsm.current, event)
	next, ok := fsm.transitions[StateEvent{State: fsm.current, Event: event}]
	if !ok {
		return
	}
	fsm.change(next)
}

// change changes the state calling Leave and Enter state methods. The method
// does not fire the state handler if next and current states are the same.
func (fsm *FSM) change(next StateType) {
	if fsm.current == next {
		return
	}
	log.Printf("fsm %p: state %v: next state %v", fsm, fsm.current, next)
	fsm.handler(fsm.current).Leave()
	fsm.handler(next).Enter()
	fsm.current = next
}

func (fsm *FSM) handler(state StateType) State {
	if h, ok := fsm.states[state]; ok {
		return h
	}
	return fsm.defaultHandler
}

// Send sends the event to the state machine asynchronously. The state machine
// changes its state accordingly to the transition table.
func (fsm *AsyncFSM) Send(e EventType) {
	select {
	case <-fsm.done:
		return
	case q := <-fsm.queue:
		fsm.queue <- append(q, e)
	}
	select { // Non-blocking notify.
	case fsm.updated <- struct{}{}:
	default:
	}
}

// Run processes the event queue in a loop.
func (fsm *AsyncFSM) Run() {
	for {
		select {
		case <-fsm.done:
			break
		case <-fsm.updated:
			q := <-fsm.queue
			fsm.queue <- nil
			for _, e := range q {
				fsm.fsm.Send(e)
			}
		}
	}
}

// Stop stops the event queue loop. Once stopped the FSM cannot be resumed.
func (fsm *AsyncFSM) Stop() {
	close(fsm.done)
}
