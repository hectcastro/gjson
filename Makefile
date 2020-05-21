TEST?=./...
GOFMT_FILES?=$$(find . -not -path "./vendor/*" -type f -name '*.go')

default: test

test: fmtcheck
	go list -mod=vendor $(TEST) | xargs -t -n4 go test $(TESTARGS) -mod=vendor -timeout=2m -parallel=4

testrace: fmtcheck
	go test -mod=vendor -race $(TEST) $(TESTARGS)

cover:
	go test $(TEST) -coverprofile=coverage.out
	go tool cover -html=coverage.out
	rm coverage.out

fmt:
	gofmt -w $(GOFMT_FILES)

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck'"

.NOTPARALLEL:

.PHONY: cover default fmt fmtcheck test testrace 
