package testhooks

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValueHook(t *testing.T) {
	valueHook := NewValueHook[string]()
	originalValue := "original value"
	fakeValue1 := "fake value1"
	fakeValue2 := "fake value2"

	// Base case, no fake error set up
	assert.Equal(t, originalValue, valueHook.Apply(originalValue))

	// When fake error is set up
	valueHook.WithValue(fakeValue1).RunFunc(func() {
		assert.Equal(t, fakeValue1, valueHook.Apply(originalValue))
	})

	// Outside of WithValue should revert to originalValue
	// This tests that we cleaned up correctly.
	assert.Equal(t, originalValue, valueHook.Apply(originalValue))

	// Nested tests (unlikely to be used but testing just in case)
	valueHook.WithValue(fakeValue1).RunFunc(func() {
		valueHook.WithValue(fakeValue2).RunFunc(func() {
			assert.Equal(t, fakeValue2, valueHook.Apply(originalValue))
		})
		assert.Equal(t, fakeValue1, valueHook.Apply(originalValue))
	})
	assert.Equal(t, originalValue, valueHook.Apply(originalValue))

	// Reusing the same setup in nested scopes must restore the outer state.
	reused := valueHook.WithValue(fakeValue1)
	reused.RunFunc(func() {
		reused.RunFunc(func() {
			assert.Equal(t, fakeValue1, valueHook.Apply(originalValue))
		})
		assert.Equal(t, fakeValue1, valueHook.Apply(originalValue))
	})
	assert.Equal(t, originalValue, valueHook.Apply(originalValue))
}
