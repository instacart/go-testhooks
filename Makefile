.PHONY: test vet

# `go test` already runs a subset of vet (printf, tests, ...) against every
# package it builds. `go vet` adds the rest -- copylocks, lostcancel,
# unusedresult -- which do not run during a test build.
test: vet
	go test ./...
	go test -tags notest ./...
	./disassemblychecks/run_checks.sh

# The mode-specific tests carry matching `notest` / `!notest` build tags, so
# both package-wide runs exercise a consistent set of expectations.
vet:
	go vet ./...
	go vet -tags notest ./...
