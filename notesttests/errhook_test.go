//go:build notest

package notesttests

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/instacart/go-testhooks/testhooks"
)

func TestErrHook(t *testing.T) {
	errHook := testhooks.NewErrHook()
	originalError := errors.New("original error")
	fakeError := errors.New("fake error")

	// Base case, no fake error set up
	assert.Equal(t, originalError, errHook.Apply(originalError))

	// Base case, no fake error set up and nil returned.
	assert.Nil(t, errHook.Apply(nil))

	// When fake error is set up, it will be ignored because this is a notest build
	errHook.WithError(fakeError).RunFunc(func() {
		assert.Equal(t, originalError, errHook.Apply(originalError))
	})
}
