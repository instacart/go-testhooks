package notesttests

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/instacart/go-testhooks/testhooks"
)

func TestValueHook(t *testing.T) {
	valueHook := testhooks.NewValueHook[string]()
	originalValue := "original value"
	fakeValue := "fake value"

	// Base case, no fake error set up
	assert.Equal(t, originalValue, valueHook.Apply(originalValue))

	// When fake error is set up, it will be ignored because this is a notest build
	valueHook.WithValue(fakeValue).RunFunc(func() {
		assert.Equal(t, originalValue, valueHook.Apply(originalValue))
	})
}
