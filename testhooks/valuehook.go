package testhooks

// Hook to return fake errors in tests.
type valueHook[T any] struct {
	// activations is the stack of active overrides. The last entry is the
	// effective value. Activations carry an ID so a cleanup can remove the
	// activation it created rather than whichever activation happens to be on
	// top.
	activations []valueActivation[T]
}

type valueActivation[T any] struct {
	id    int
	value T
}

// WithValue returns a WithSetup that will setup valueHook with fakeValue.
//
// When Apply is called and the setup is in effect, this hook will return
// fakeValue.
func (e *valueHook[T]) WithValue(fakeValue T) WithSetup {
	return newWithSetup(
		func(activationID int) {
			e.activations = append(e.activations, valueActivation[T]{
				id:    activationID,
				value: fakeValue,
			})
		},
		func(activationID int) {
			for i := len(e.activations) - 1; i >= 0; i-- {
				if e.activations[i].id != activationID {
					continue
				}
				// RunFunc is a scope for this hook: clean this activation
				// and everything activated on the same hook inside it. clear
				// releases the values before the backing array is shortened.
				clear(e.activations[i:])
				e.activations = e.activations[:i]
				return
			}
		},
	)
}

// Returns a hook to return fake errors in tests.
func NewValueHook[T any]() *valueHook[T] {
	return &valueHook[T]{}
}
