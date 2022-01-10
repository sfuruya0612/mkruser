NAME := $(shell basename $(CURDIR))
MODULE := $(shell head -1 go.mod | awk '{print $$2}')

DATE := $(shell TZ=Asia/Tokyo date +%Y%m%d-%H%M%S)
HASH := $(shell git rev-parse --short HEAD)
GO_VERSION := $(shell go version)
LDFLAGS := -s -w -X 'main.name=${NAME}' -X 'main.date=${DATE}' -X 'main.hash=${HASH}' -X 'main.goversion=${GO_VERSION}'

.PHONY:

init:
	asdf install

tidy:
	go mod tidy -go=1.17

fmt:
	go fmt ./...

vet:
	go vet ./...

test:
	go test -v -race --cover ./...

install: fmt vet tidy
	go install -ldflags "${LDFLAGS}" ${MODULE}

clean:
	-rm ${GOPATH}/bin/${NAME}

