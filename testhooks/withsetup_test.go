package testhooks

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithSetup(t *testing.T) {
	var events []string

	setupFunc := func() {
		events = append(events, "setup")
	}

	cleanupFunc := func() {
		events = append(events, "cleanup")
	}

	sampleFunc := func() {
		events = append(events, "sample")
	}

	ws := NewWithSetup(setupFunc, cleanupFunc)

	t.Run("RunFunc", func(t *testing.T) {
		events = nil
		ws.RunFunc(func() {
			assert.Equal(t, []string{"setup"}, events)
			sampleFunc()
			assert.Equal(t, []string{"setup", "sample"}, events)
		})
		assert.Equal(t, []string{"setup", "sample", "cleanup"}, events)
	})

	t.Run("SetupInTest", func(t *testing.T) {
		events = nil
		t.Run("SetupInTest", func(t *testing.T) {
			ws.SetupInTest(t)
			assert.Equal(t, []string{"setup"}, events)
			sampleFunc()
			assert.Equal(t, []string{"setup", "sample"}, events)
		})
		assert.Equal(t, []string{"setup", "sample", "cleanup"}, events)
	})
}
