package testhooks

// Hook to return fake errors in tests.
type valueHook[T any] struct {
	fakeValue    T
	fakeValueSet bool
}

// WithError returns a WithSetup that will setup valueHook with fakeValue.
//
// When Apply is called and the setup is in effect, this hook will return
// fakeValue.
func (e *valueHook[T]) WithValue(fakeValue T) WithSetup {
	var (
		oldFakeValue    T
		oldFakeValueSet bool
	)
	return NewWithSetup(
		func() {
			oldFakeValue = e.fakeValue
			oldFakeValueSet = e.fakeValueSet
			e.fakeValue = fakeValue
			e.fakeValueSet = true
		},
		func() {
			e.fakeValue = oldFakeValue
			e.fakeValueSet = oldFakeValueSet
		},
	)
}

// Returns a hook to return fake errors in tests.
func NewValueHook[T any]() *valueHook[T] {
	return &valueHook[T]{}
}
