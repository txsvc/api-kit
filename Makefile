.PHONY: all
all: test

.PHONY: test
test:
	go test
	cd settings && go test


.PHONY: test_build
test_build:
	cd examples/simple_api && go build main.go && rm main
	

.PHONY: test_coverage
test_coverage:
	go test `go list ./... | grep -v 'hack\|deprecated\|examples'` -coverprofile=coverage.txt -covermode=atomic
