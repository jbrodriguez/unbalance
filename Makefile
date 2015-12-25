#
# Makefile to perform "live code reloading" after changes to .go files.
#
# n.b. you must install fswatch (OS X: `brew install fswatch`)
#
# To start live reloading run the following command:
# $ make serve
#

mb_version := $(shell cat VERSION)
mb_count := $(shell git rev-list HEAD --count)
mb_hash := $(shell git rev-parse --short HEAD)

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
	go build -ldflags "-X main.Version=$(mb_version)-$(mb_count).$(mb_hash)" -v -o dist/unbalance server/unbalance.go

server: dependencies clean
	go build fmt
	env GOOS=linux GOARCH=amd64 go build -ldflags "-X main.Version=$(mb_version)-$(mb_count).$(mb_hash)" -v -o dist/unbalance server/unbalance.go

client: dependencies
	npm run build


# buildx: clean
# 	go build fmt
# 	env GOOS=linux GOARCH=amd64 go build -ldflags "-X main.Version=$(mb_version)-$(mb_count).$(mb_hash)" -v -o dist/unbalance server/unbalance.go
 
# run unit tests with code coverage
test: dependencies 
	go test -v
 
# generate code coverage report
coverage: test
	go build test -coverprofile=.coverage.out
	go build tool cover -html=.coverage.out
 
publish: dependencies client server
	rsync -avzP -e "ssh" dist/* $(SERVER):/boot/custom/unbalance
