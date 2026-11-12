package main

import (
	"fmt"

	"github.com/instacart/go-testhooks/testhooks"
)

var sampleValueHook = testhooks.NewValueHook[string]()

func RealFunc() string {
	return "real"
}

func main() {
	sampleValueHook.WithValue("fake").RunFunc(run)
}

// run() is the function that will be inspected during the diassemblytest
// When built with notest, we expect all references to sampleValueHook
// to be removed.
//
// We also use noinline to prevent inlining of the function in the main
// method as we want to look at it separately.
//
//go:noinline
func run() {
	result := RealFunc()
	result = sampleValueHook.Apply(result)
	Print(result)
}

//go:noinline
func Print(result string) {
	fmt.Printf("result = %s\n", result)
}
