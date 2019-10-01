PROJECT?=github.com/pabluk/promql-share
RELEASE?=0.0.1

COMMIT := git-$(shell git rev-parse --short HEAD)
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')

test:
	go test --race ./...

build:
	GO111MODULE=on CGO_ENABLED=0 go build \
		-ldflags "-s -w -X ${PROJECT}/internal/diagnostics.Version=${RELEASE} \
		-X ${PROJECT}/internal/diagnostics.Commit=${COMMIT} \
		-X ${PROJECT}/internal/diagnostics.BuildTime=${BUILD_TIME}" \
		-o bin/promql-share ${PROJECT}/cmd/promql-share

docker-build:
	docker build -t promql-share:${RELEASE} .

