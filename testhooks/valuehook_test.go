//go:build !notest

package testhooks

import (
	"sync"
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

	// RunFunc removes only its own activation. A SetupInTest created inside
	// the block stays active until that test's cleanup runs.
	setupHook := NewValueHook[string]()
	setupHook.WithValue(fakeValue1).RunFunc(func() {
		setupHook.WithValue(fakeValue2).SetupInTest(t)
		assert.Equal(t, fakeValue2, setupHook.Apply(originalValue))
	})
	assert.Equal(t, fakeValue2, setupHook.Apply(originalValue))

	// Other hooks are independent: a SetupInTest on another hook keeps its
	// test-lifetime cleanup even when created inside this hook's RunFunc.
	otherHook := NewValueHook[string]()
	valueHook.WithValue(fakeValue1).RunFunc(func() {
		otherHook.WithValue(fakeValue2).SetupInTest(t)
	})
	assert.Equal(t, fakeValue2, otherHook.Apply(originalValue))

	// Apply must be safe while another goroutine is setting up and tearing
	// down overrides. Run with -race to verify the snapshots are immutable.
	var wg sync.WaitGroup
	stop := make(chan struct{})
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-stop:
					return
				default:
					_ = valueHook.Apply(originalValue)
				}
			}
		}()
	}
	for i := 0; i < 100; i++ {
		valueHook.WithValue(fakeValue1).RunFunc(func() {})
	}
	close(stop)
	wg.Wait()
}
