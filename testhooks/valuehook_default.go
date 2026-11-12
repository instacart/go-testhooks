//go:build !notest || _gazelle_dummy

// _gazelle_dummy exists to force gazelle to include this file in the BUILD file
// since it's mutually exclusive with _notest.go file.

package testhooks

// Apply applies the hook to the value and returns the result.
//
// If a fakeValue has been setup by calling WithValue, this will
// return the fakeValue instead of the originally given value.
//
// Note: When built with notest build tag, this method implementation
// will be replaced with one that just returns the value back.
func (e *valueHook[T]) Apply(value T) T {
	if len(e.activations) > 0 {
		return e.activations[len(e.activations)-1].value
	}
	return value
}
