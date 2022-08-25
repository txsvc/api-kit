EXAMPLE_NAME = sc

TARGET_LINUX = GOARCH=amd64 GOOS=linux
CONTAINER_REGISTRY = eu.gcr.io/podops

.PHONY: all
all: test

.PHONY: test
test:
	cd api && go test -covermode=atomic
	cd cli && go test -covermode=atomic
	cd config && go test -covermode=atomic
	cd internal && go test -covermode=atomic
	cd internal/cli && go test -covermode=atomic
	cd internal/settings && go test -covermode=atomic
	go test -covermode=atomic

.PHONY: test_build
test_build:
	go mod verify && go mod tidy
	cd examples/auth/api && go build main.go && rm main
	cd examples/auth/cli && go build cli.go && rm cli
	cd examples/appengine && go build main.go && rm main


.PHONY: examples
examples: example_cli example_api

.PHONY: example_appengine
example_appengine:
	cd examples/appengine && gcloud app deploy . --quiet

.PHONY: example_cli
example_cli:
	cd examples/auth/cli && go build -o ${EXAMPLE_NAME} cli.go && mv ${EXAMPLE_NAME} ../../../bin/${EXAMPLE_NAME}
	chmod +x bin/${EXAMPLE_NAME}

.PHONY: example_api
example_api:
	cd examples/auth/cli && go build -o svc cli.go && mv svc ../../../bin/${EXAMPLE_NAME}svc
	chmod +x bin/${EXAMPLE_NAME}svc

#.PHONY example_api_container
#example_api_container:
#	cd examples/simple_api && ${TARGET_LINUX} go build -o svc main.go && podman build -t "" .
#	rm examples/simple_api/svc
