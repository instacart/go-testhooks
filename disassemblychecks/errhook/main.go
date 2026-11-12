package main

import (
	"errors"
	"fmt"

	"github.com/instacart/go-testhooks/testhooks"
)

var sampleErrHook = testhooks.NewErrHook()

func RealFunc() (string, error) {
	return "real", nil
}

func main() {
	sampleErrHook.WithError(errors.New("fake error")).RunFunc(run)
}

// run() is the function that will be inspected during the diassemblytest
// When built with notest, we expect all references to sampleErrHook
// to be removed.
//
// We also use noinline to prevent inlining of the function in the main
// method as we want to look at it separately.
//
//go:noinline
func run() {
	result, err := RealFunc()
	err = sampleErrHook.Apply(err)
	Print(result, err)
}

//go:noinline
func Print(result string, err error) {
	if err == nil {
		fmt.Printf("result = %s\n", result)
	} else {
		fmt.Printf("error = %s\n", err.Error())
	}
}
