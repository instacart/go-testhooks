package testhooks

// Hook to return fake errors in tests.
type errHook valueHook[error]

// WithError returns a WithSetup that will setup errHook with fakeErr.
//
// When Apply is called and the setup is in effect, this hook will return
// fakeErr.
func (e *errHook) WithError(fakeErr error) WithSetup {
	return (*valueHook[error])(e).WithValue(fakeErr)
}

// Returns a hook to return fake errors in tests.
func NewErrHook() *errHook {
	return &errHook{}
}
