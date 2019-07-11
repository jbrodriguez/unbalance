#
# Makefile to perform "live code reloading" after changes to .go files.
#
# n.b. you must install fswatch (OS X: `brew install fswatch`)
#
# To start live reloading run the following command:
# $ make serve
#

mb_version := $(shell cat ../VERSION)
mb_count := $(shell git rev-list HEAD --count)
mb_hash := $(shell git rev-parse --short HEAD)
mb_plgver := $(shell cat ../PLGVER)

# binary name to kill/restart
PROG = unbalance

# targets not associated with files
.PHONY: dependencies default build test coverage clean kill restart serve

 # check we have a couple of dependencies
dependencies:
	@command -v fswatch --version >/dev/null 2>&1 || { printf >&2 "fswatch is not installed, please run: brew install fswatch\n"; exit 1; }

# default targets to run when only running `make`
default: dependencies test

# clean up
clean:
	go clean

# run formatting tool and build
build: dependencies clean
	go build fmt
	go build -ldflags "-s -w -X main.version=$(mb_version) -X main.plgver=${mb_plgver}" -gcflags "all=-trimpath=$GOPATH" -v -o ${PROG} ${PROG}.go

server: dependencies clean
	go build fmt
	env GOOS=linux GOARCH=amd64 go build -ldflags "-s -w -X main.version=$(mb_version) -X main.plgver=${mb_plgver}" -gcflags "all=-trimpath=$GOPATH" -v -o ../dist/${PROG} ${PROG}.go

# run unit tests with code coverage
test: dependencies
	go test -v

# generate code coverage report
coverage: test
	go build test -coverprofile=.coverage.out
	go build tool cover -html=.coverage
