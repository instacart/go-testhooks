# Test Hooks (testhooks)

Behavior-stubbing hooks that use Go build tags (`!notest`) to compile test code out of
production binaries completely.

Pre-1.0, maintained on a best-effort basis. See [Project status](#project-status).

## Usage

### Injecting errors

In test cases, we sometimes need to exercise error conditions returned by your dependencies that are hard to reproduce in unit tests. The solution to reproduce those scenarios typically involves one of the following options:
- Doing interaction tests with mocks
- Creating a fake version of the dependency

Both of these can however be hard when the interaction with the dependency is complicated (e.g. testing `database/sql` query interaction which involves many steps). This package implements an alternative where instead of replacing your dependencies you use them as is in tests and just add hooks to inject errors / failures.

For example, if you are testing:
```go
func MyFunc() (string, error) {
    // ... other code ...
    result, err := somedependency.DoSomething()
    if err != nil {
        // error handling like:
        return "", err
    }
    // ... other code ...
}
```

Using this library you can do:
```go
var DoSomethingErrHook = testhooks.NewErrHook()

func MyFunc() (string, error) {
    // ... other code ...
    result, err := somedependency.DoSomething()
    err = DoSomethingErrHook.Apply(err) // <------ added line
    if err != nil {
        // error handling like:
        return "", err
    }
    // ... other code ...
}
```

Then in tests you can simulate the error condition with:
```go
import (
    "errors"
    "testing"

	"github.com/stretchr/testify/assert"
)

func TestMyFunc(t *testing.T) {
    fakeErr := errors.New("fake error")
    // Will make DoSomethingErrHook.Apply return fakeErr.
    // This change will be reverted at the end of the test.
    DoSomethingErrHook.WithError(fakeErr).SetupInTest(t)
    _, err := MyFunc()
    assert.Equal(t, fakeErr, err)
}
```

### Injecting values in tests

Similar to the case for errors, there are times we want to be able to inject values in a test. One example is if we have something that is supposed to be a constant value that we want to change just during tests. The solution usually involves making the value not a constant and adding get/set methods to change it during tests.

To solve for that, we can generalize the errhook concept mentioned in the previous section to arbitrary values.

For example, if you are testing:
```go
const maxDbResults = 1000

func MyFunc() {
    results := FetchFromDb()
    if (len(results) >= maxDbResults) {
        // do something
    } else {
        // do something else
    }
}
```

And in tests you would like to change maxDbResults to exercise the two branches, using this library you can do:
```go
var maxDbResultsValueHook = testhooks.NewValueHook[int]()

func maxDbResults() int {
    return maxDbResultsValueHook.Apply(1000)
}

func MyFunc() {
    results := FetchFromDb()
    if (len(results) >= maxDbResults()) {
        // do something
    } else {
        // do something else
    }
}
```

Then in tests you can change the value with:
```go
import (
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestMyFunc(t *testing.T) {
    // Change maxDbResults to 3.
    // This change will be reverted at the end of the test.
    maxDbResultsValueHook.WithValue(3).SetupInTest(t)
    MyFunc()
    // do assertions ...
}
```

The alternative if you didn't use this library would be to make `maxDbResults` a var and change it in tests but making sure you revert it back to the original value.

### Alternate setup usage with RunFunc

By default the hooks return a `WithSetup` object that you can call `.SetupInTest(t)` that will run the setup immediately but also schedule to clean up any changes that were made at the end of the test. There is also an alternative `.RunFunc(func() {})` in case you want the change to the hook to apply just within a small block.

`RunFunc` defines a scope: any setup activated inside the block, including a `SetupInTest`, is cleaned up when the block returns. If an override must remain active until the test ends, register its `SetupInTest` outside the `RunFunc` block.

So in the example from "Injecting values in tests" section above, you could write:

```go
func TestMyFunc(t *testing.T) {
    // Change maxDbResults to 3 within the func() {} block
    // This change will be reverted once RunFunc returns.
    maxDbResultsValueHook.WithValue(3).RunFunc(func() {
        MyFunc()
        // do assertions ...
    })
}
```

## Development

Requires Go 1.25 or later.

Tests also require Python 3 for the disassembly check; the library builds without it.

```bash
make test
```

## Advantages/Disadvantages

Advantages:
- For ErrHooks, avoids the need to create fakes/mocks for dependencies just to test error cases. This can be particularly useful when the dependencies are hard to fake/replicate.
- Avoids the need for you to clean up yourself any changes to values in tests.
- Related to the above, allows you to achieve higher levels of code coverage with less effort.
- When built with the `notest` build tag, there should be no overhead to using the hooks (this is discussed in the [Mitigating runtime cost of the hooks](#mitigating-runtime-cost-of-the-hooks) section).

Disadvantages:
- Extra lines in your code for calling the hooks.
- Need to set up your production build with the `notest` tag to reduce runtime costs of testhooks.

## Mitigating runtime cost of the hooks

To mitigate runtime costs introduced by the hooks, this library leverages Go build tags and the inline behavior of the Go compiler to strip away instructions related to the test hooks themselves.

When the code is built with the `notest` build tag, the `Apply` functions of the hooks are replaced with basic implementations that just return back the value they were given in. This allows the Go compiler to inline the method and basically eliminate any runtime cost of using the hooks. The code generated assembly should be roughly equivalent to that generated by the same code without the call to the `Apply` method of the hooks.

By default (when the `notest` build tag is not used), it assumes the code being built may be used in tests, so it includes implementations of the `Apply` methods with the branches necessary for injecting errors/values.

However production binaries should be built using the `notest` tag like the following:
```bash
go build -tags notest ./...
```

To verify that this mitigation works, we run some checks on the disassembled code. More details can be found in [disassemblychecks](disassemblychecks/README.md).

## Project status

A variant of this library is used at Instacart. We're publishing it so others
can use, adapt, and build on the ideas.

This is not an officially supported Instacart product and comes with no
commitment to ongoing maintenance. Fixes and improvements are made on a
best-effort basis, and the repository may be archived in the future.

The API is pre-1.0 and may change between minor versions; breaking changes
will be noted in the release notes.

Licensed under the Apache License 2.0. See [LICENSE](LICENSE).
