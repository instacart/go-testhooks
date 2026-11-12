package testhooks

import "sync/atomic"

// activationList is an immutable snapshot of active overrides. Snapshots are
// published with copy-on-write so Apply can read without a lock.
type activationList[T any] []valueActivation[T]

type valueActivation[T any] struct {
	id    int
	value T
}

// Hook to return fake errors in tests.
type valueHook[T any] struct {
	activations atomic.Pointer[activationList[T]]
}

// WithValue returns a WithSetup that will setup valueHook with fakeValue.
//
// When Apply is called and the setup is in effect, this hook will return
// fakeValue.
func (e *valueHook[T]) WithValue(fakeValue T) WithSetup {
	return newWithSetup(
		func(activationID int) {
			for {
				current := e.activations.Load()
				var next activationList[T]
				if current != nil {
					next = make(activationList[T], len(*current), len(*current)+1)
					copy(next, *current)
				}
				next = append(next, valueActivation[T]{id: activationID, value: fakeValue})
				if e.activations.CompareAndSwap(current, newActivationList(next)) {
					return
				}
			}
		},
		func(activationID int) {
			for {
				current := e.activations.Load()
				if current == nil {
					return
				}
				activations := *current
				index := -1
				for i := len(activations) - 1; i >= 0; i-- {
					if activations[i].id == activationID {
						index = i
						break
					}
				}
				if index < 0 {
					return
				}
				// Remove only this activation; other activations, including
				// ones created after it, belong to their own scopes.
				next := make(activationList[T], len(activations)-1)
				copy(next, activations[:index])
				copy(next[index:], activations[index+1:])
				if e.activations.CompareAndSwap(current, newActivationList(next)) {
					return
				}
			}
		},
	)
}

func newActivationList[T any](activations activationList[T]) *activationList[T] {
	if len(activations) == 0 {
		return nil
	}
	return &activations
}

// Returns a hook to return fake errors in tests.
func NewValueHook[T any]() *valueHook[T] {
	return &valueHook[T]{}
}
