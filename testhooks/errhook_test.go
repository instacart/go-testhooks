//go:build !notest

package testhooks

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrHook(t *testing.T) {
	errHook := NewErrHook()
	originalError := errors.New("original error")
	fakeError1 := errors.New("fake error1")
	fakeError2 := errors.New("fake error2")

	// Base case, no fake error set up
	assert.Equal(t, originalError, errHook.Apply(originalError))

	// Base case, no fake error set up and nil returned.
	assert.Nil(t, errHook.Apply(nil))

	// When fake error is set up
	errHook.WithError(fakeError1).RunFunc(func() {
		assert.Equal(t, fakeError1, errHook.Apply(originalError))
	})

	// Outside of WithError should revert to originalError
	// This tests that we cleaned up correctly.
	assert.Equal(t, originalError, errHook.Apply(originalError))

	// Nested tests (unlikely to be used but testing just in case)
	errHook.WithError(fakeError1).RunFunc(func() {
		errHook.WithError(fakeError2).RunFunc(func() {
			assert.Equal(t, fakeError2, errHook.Apply(originalError))
		})
		assert.Equal(t, fakeError1, errHook.Apply(originalError))
	})
	assert.Equal(t, originalError, errHook.Apply(originalError))
}
