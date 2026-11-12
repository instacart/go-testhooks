package testhooks

// Hook to return fake errors in tests.
type valueHook[T any] struct {
	fakeValue    T
	fakeValueSet bool

	// previous holds the state saved by each activation. A single WithSetup
	// can be reused in nested scopes, so restores must happen in LIFO order.
	previous []valueHookState[T]
}

// valueHookState is the state saved by one activation of a WithSetup.
type valueHookState[T any] struct {
	value T
	set   bool
}

// WithValue returns a WithSetup that will setup valueHook with fakeValue.
//
// When Apply is called and the setup is in effect, this hook will return
// fakeValue.
func (e *valueHook[T]) WithValue(fakeValue T) WithSetup {
	return NewWithSetup(
		func() {
			e.previous = append(e.previous, valueHookState[T]{
				value: e.fakeValue,
				set:   e.fakeValueSet,
			})
			e.fakeValue = fakeValue
			e.fakeValueSet = true
		},
		func() {
			last := len(e.previous) - 1
			state := e.previous[last]
			// Clear the slot before reslicing so the backing array does not
			// keep the restored value reachable.
			e.previous[last] = valueHookState[T]{}
			e.previous = e.previous[:last]
			e.fakeValue = state.value
			e.fakeValueSet = state.set
		},
	)
}

// Returns a hook to return fake errors in tests.
func NewValueHook[T any]() *valueHook[T] {
	return &valueHook[T]{}
}
