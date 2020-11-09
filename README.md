# fsm-go
Golang finite-state machine.

# Examples

## Comprehensive

In [this example] we do the following:

* utilize a back-reference to our state maching using `MachineProvider`;
* use `Expiring` state handler to timeout a state;
* create custom state handlers that print the states and terminate the app;
* join the handlers using `Compose`.

[this example]: example_test.go

## Simple

In this example we start with `stInit`, and then send `evInitialized`, which
switches the state into `stFinal`.

```
const (
	stInit         StateType = "stInit"
	stFinal        StateType = "stFinal"
	evInitialized  EventType = "evInitialized"
)

transitions := Transitions{
	When(stInit, evInitialized): stFinal,
}

fsm := fsm.New(stInit, transitions, nil, nil)
fsm.Send(evInitialized)
```

## Expiring states

In this example we use an `Expiring` state, which sends events to the FSM, and
requires a back-reference.

```golang
const (
	stInit         StateType = "stInit"
	stWaitTimeout  StateType = "stWaitTimeout"
	stFinal        StateType = "stFinal"
	evInitialized  EventType = "evInitialized"
	evTimeout      EventType = "evTimeout"
)

transitions := Transitions{
	When(stInit, evInitialized):    stWaitTimeout,
	When(stWaitTimeout, evTimeout): stFinal,
}

provider := new(MachineProvider)

states := States{
	"stMyState": NewExpiring(provider, After(time.Second), evTimeout),
}

fsm := fsm.New(stInit, transitions, states, nil)
provider.Set(fsm)
fsm.Send(evInitialized)
```
