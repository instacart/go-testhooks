//go:build !notest || _gazelle_dummy

// _gazelle_dummy exists to force gazelle to include this file in the BUILD file
// since it's mutually exclusive with _notest.go file.

package testhooks

// Apply applies the hook to the error and returns the result.
//
// If a fakeErr has been setup by calling WithError, this will
// return the fakeErr instead of the original error.
//
// Note: When built with notest build tag, this method implementation
// will be replaced with one that just returns the value back.
func (e *errHook) Apply(err error) error {
	return (*valueHook[error])(e).Apply(err)
}
