package testhooks

import "sync/atomic"

var nextActivationID atomic.Int64

type withSetup struct {
	setupFunc   func(activationID int)
	cleanupFunc func(activationID int)
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
	//
	// RunFunc defines a scope: activations created inside block, including
	// ones registered with SetupInTest, are cleaned up when the block returns.
	RunFunc(block func())

	// SetupInTest runs the setup immediately and schedules the cleanup
	// function to be run when the current test/subtest finishes.
	//
	// If called inside a RunFunc scope, it is scoped to that RunFunc: the
	// cleanup runs when the RunFunc returns, not at the end of the test.
	//
	// Note: it expects to receive *testing.T as an argument,
	// however the current type argument is an interface with
	// a subset of the methods in *testing.T to avoid a direct
	// dependency.
	SetupInTest(t testingT)
}

func (w withSetup) RunFunc(block func()) {
	activationID := newActivationID()
	w.setupFunc(activationID)
	defer w.cleanupFunc(activationID)
	block()
}

func (w withSetup) SetupInTest(t testingT) {
	activationID := newActivationID()
	w.setupFunc(activationID)
	t.Cleanup(func() { w.cleanupFunc(activationID) })
}

func newActivationID() int {
	return int(nextActivationID.Add(1))
}

// NewWithSetup returns a WithSetup.
//
// The callbacks are not activation-aware; use a hook's WithValue/WithError when
// the setup needs to be scoped to the activation that created it.
func NewWithSetup(setupFunc func(), cleanupFunc func()) WithSetup {
	return newWithSetup(
		func(int) { setupFunc() },
		func(int) { cleanupFunc() },
	)
}

func newWithSetup(setupFunc func(activationID int), cleanupFunc func(activationID int)) WithSetup {
	return withSetup{
		setupFunc:   setupFunc,
		cleanupFunc: cleanupFunc,
	}
}
