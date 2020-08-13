.PHONY: all clean

COVERAGE_FILE := coverage.txt
GOTEST_FLAGS  := -v -race -coverprofile=${COVERAGE_FILE}

all: build

build:
	go mod tidy
	go build

test:
	go test ${GOTEST_FLAGS} ./...
	go tool cover -html=${COVERAGE_FILE}

clean:
	rm -f ${COVERAGE_FILE}