package testhooks

type withSetup struct {
	setupFunc   func()
	cleanupFunc func()
}

// testingT is an interface with a subset of the methods in testing.T
//
// Used so we can avoid depending on testing.T directly.
type testingT interface {
	Cleanup(func())
}

// WithSetup is an interface that can be used to represent some setup and
// its associated cleanup function.
type WithSetup interface {
	// RunFunc runs the given block calling the setup function and the
	// cleanup function.
	RunFunc(block func())

	// SetupInTest runs the setup immediately and schedules the cleanup
	// function to be run when the current test/subtest finishes.
	//
	// Note: it expects to receive *testing.T as an argument,
	// however the current type argument is an interface with
	// a subset of the methods in *testing.T to avoid a direct
	// dependency.
	SetupInTest(t testingT)
}

func (w withSetup) RunFunc(block func()) {
	w.setupFunc()
	defer w.cleanupFunc()
	block()
}

func (w withSetup) SetupInTest(t testingT) {
	w.setupFunc()
	t.Cleanup(w.cleanupFunc)
}

// NewWithSetup returns a WithSetup.
func NewWithSetup(setupFunc func(), cleanupFunc func()) WithSetup {
	return withSetup{
		setupFunc:   setupFunc,
		cleanupFunc: cleanupFunc,
	}
}
