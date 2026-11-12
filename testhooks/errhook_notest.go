//go:build notest || _gazelle_dummy

// _gazelle_dummy exists to force gazelle to include this file in the BUILD file
// since it's mutually exclusive with _default.go file.

package testhooks

// Apply when built with notest build tag will just return the err back.
//
// This allow the Go compiler to inline and get rid of function overhead for
// calling Apply.
//
// See errhook_default.go for documentation.
func (e *errHook) Apply(err error) error {
	return err
}
