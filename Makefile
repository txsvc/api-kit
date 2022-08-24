EXAMPLE_CLI_NAME = sc

TARGET_LINUX = GOARCH=amd64 GOOS=linux
CONTAINER_REGISTRY = eu.gcr.io/podops

.PHONY: all
all: test

.PHONY: test
test:
	cd api && go test
	cd cli && go test
	cd config && go test
	cd internal && go test
	cd internal/cli && go test
	cd internal/settings && go test
	go test

.PHONY: test_build
test_build:
	go mod verify && go mod tidy
	cd examples/auth/api && go build main.go && rm main
	cd examples/auth/cli && go build cli.go && rm cli
	cd examples/appengine && go build main.go && rm main
	
.PHONY: test_coverage
test_coverage:
	go test `go list ./... | grep -v 'hack\|deprecated\|examples'` -coverprofile=coverage.txt -covermode=atomic


.PHONY: examples
examples: example_cli example_appengine

.PHONY: example_appengine
example_appengine:
	cd examples/appengine && gcloud app deploy . --quiet

#.PHONY example_api_container
#example_api_container:
#	cd examples/simple_api && ${TARGET_LINUX} go build -o svc main.go && podman build -t "" .
#	rm examples/simple_api/svc

.PHONY: example_cli
example_cli:
	cd examples/cli && go build -o ${EXAMPLE_CLI_NAME} cli.go && mv ${EXAMPLE_CLI_NAME} ../../bin/${EXAMPLE_CLI_NAME}
	chmod +x bin/${EXAMPLE_CLI_NAME}

