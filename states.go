package fsm

import (
	"time"
)

type (
	// StateType is a type of FSM states.
	StateType string
	// State is a state event handler.
	State interface {
		// Enter is called upon entering this state.
		Enter()
		// Leave is called upon leaving this state.
		Leave()
	}
	// EventSender wraps FSM.Send.
	EventSender interface{ Send(EventType) }
	// MachineProvider provides a reference to the FSM. The provider is not
	// go-routine safe. We recommend to put a machine into an initialization
	// state and configure it before actually running it. See the Expiring
	// usage example.
	MachineProvider struct{ fsm EventSender }
	// Expiring is a state handler with a timeout. When the state timeout
	// occurres the event is sent to the FSM. We recommend to use the async
	// implementation of the state machine, since the timeout is sent
	// asynchronously.
	//
	// The state machine must be provided (using a MachineProvider) before
	// entering the state.
	//
	// Usage example:
	//
	//	const (
	//		stInit         StateType = "stInit"
	//		stWaitTimeout  StateType = "stWaitTimeout"
	//		stFinal        StateType = "stFinal"
	//		evInitialized  EventType = "evInitialized"
	//		evTimeout      EventType = "evTimeout"
	//	)
	//	transitions := Transitions{
	//		When(stInit, evInitialized):    stWaitTimeout,
	//		When(stWaitTimeout, evTimeout): stFinal,
	//	}
	//	provider := new(MachineProvider)
	//	states := States{
	//		"stMyState": NewExpiring(provider, After(time.Second), evTimeout),
	//	}
	//	fsm := fsm.NewAsync(stInit, transitions, states, nil)
	//	provider.Set(fsm)
	//	fsm.Send(evInitialized)
	//	fsm.Run()
	Expiring struct {
		fsm    EventSender
		event  EventType
		later  AfterFunc
		cancel func()
	}
	// Composite contains multiple states handlers and dispatches state
	// notifications to each of them.
	Composite []State

	// States describes state event handlers.
	States map[StateType]State
)

// DefaultHandler is a state handler that does nothing.
var DefaultHandler defaultHandler

// Enter sends enter notification to every state contained.
func (st Composite) Enter() {
	for _, s := range st {
		s.Enter()
	}
}

// Leave sends leave notification to every state contained.
func (st Composite) Leave() {
	for _, s := range st {
		s.Leave()
	}
}

// Compose wraps all states with a composite state.
func Compose(states ...State) (s Composite) {
	for _, st := range states {
		s = append(s, st)
	}
	return
}

type defaultHandler struct{}

func (defaultHandler) Enter() {}
func (defaultHandler) Leave() {}

// Set sets the FSM reference.
func (p *MachineProvider) Set(fsm EventSender) {
	p.fsm = fsm
}

// Send sends the event to the FSM.
func (p *MachineProvider) Send(e EventType) {
	p.fsm.Send(e)
}

// AfterFunc schedules the callback and allows to cancel the execution.
type AfterFunc func(f func()) (cancel func())

// After returns a scheduler configured with the duration.
func After(d time.Duration) AfterFunc {
	return func(f func()) (cancel func()) {
		timer := time.AfterFunc(d, f)
		return func() { timer.Stop() }
	}
}

// NewExpiring returns the state handler which send the event to the state
// machine after some time interval.
func NewExpiring(p *MachineProvider, after AfterFunc, e EventType) *Expiring {
	return &Expiring{
		fsm:   p,
		event: e,
		later: after,
	}
}

// Enter sets up the timeout timer.
func (st *Expiring) Enter() {
	sendEvent := func() { st.fsm.Send(st.event) }
	st.cancel = st.later(sendEvent)
}

// Leave cancels the timeout timer.
func (st *Expiring) Leave() {
	// TODO(yarcat): Cancel sent event.
	st.cancel()
}
