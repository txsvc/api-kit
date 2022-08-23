EXAMPLE_CLI_NAME = sc

TARGET_LINUX = GOARCH=amd64 GOOS=linux
CONTAINER_REGISTRY = eu.gcr.io/podops

.PHONY: all
all: test

.PHONY: test
test:
	cd cli && go test
	cd config && go test
	cd internal && go test
	cd internal/cli && go test
	cd internal/settings && go test
	go test

.PHONY: test_build
test_build:
	go mod verify && go mod tidy
	cd example/api && go build main.go && rm main
	cd example/cli && go build cli.go && rm cli
	
.PHONY: test_coverage
test_coverage:
	go test `go list ./... | grep -v 'hack\|deprecated\|examples'` -coverprofile=coverage.txt -covermode=atomic


.PHONY: examples
examples: example_cli example_api

.PHONY: example_api
example_api:
	cd example/api && gcloud app deploy . --quiet

#.PHONY example_api_container
#example_api_container:
#	cd examples/simple_api && ${TARGET_LINUX} go build -o svc main.go && podman build -t "" .
#	rm examples/simple_api/svc

.PHONY: example_cli
example_cli:
	cd example/cli && go build -o ${EXAMPLE_CLI_NAME} cli.go && mv ${EXAMPLE_CLI_NAME} ../../bin/${EXAMPLE_CLI_NAME}
	chmod +x bin/${EXAMPLE_CLI_NAME}

