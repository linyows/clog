TEST ?= ./...
ifeq ("$(shell uname)","Darwin")
NCPU ?= $(shell sysctl hw.ncpu | cut -f2 -d' ')
else
NCPU ?= $(shell cat /proc/cpuinfo | grep processor | wc -l)
endif
TEST_OPTIONS=-timeout 30s -parallel $(NCPU)

default: build

deps:
	go get golang.org/x/lint/golint
	go get github.com/goreleaser/goreleaser

build:
	go build .

test:
	go test $(TEST) $(TESTARGS) $(TEST_OPTIONS)
	go test -race $(TEST) $(TESTARGS) -coverprofile=coverage.txt -covermode=atomic

lint:
	golint -set_exit_status $(TEST)

ci: deps test lint
	go mod tidy

dist:
	@test -z $(GITHUB_TOKEN) || goreleaser

clean:
	rm -rf coverage.txt
	git checkout go.*
