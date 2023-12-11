#
# Makefile to perform "live code reloading" after changes to .go files.
#
# To start live reloading run the following command:
# $ make serve
#

mb_date := $(shell date '+%Y.%m.%d')
mb_hash := $(shell git rev-parse --short HEAD)

# binary name to kill/restart
PROG = unbalance
 
# targets not associated with files
.PHONY: default build test coverage clean kill restart serve
 
# default targets to run when only running `make`
default: test
 
# clean up
clean:
	go clean

protobuf:
	protoc --go_out=. --go_opt=paths=source_relative --go-drpc_out=. --go-drpc_opt=paths=source_relative import.proto

local: clean
	pushd ./ui && npm run build && popd
	go build fmt
	go build -ldflags "-X main.Version=$(mb_date)-$(mb_hash)" -v -o ${PROG}

release: clean
	pushd ./ui && npm run build && popd
	go build fmt
	GOOS=linux GOARCH=amd64 go build -ldflags "-X main.Version=$(mb_date)-$(mb_hash)" -v -o ${PROG}

# run unit tests with code coverage
test: 
	go test -v
 
# generate code coverage report
coverage: test
	go build test -coverprofile=.coverage.out
	go build tool cover -html=.coverage.out
 
# attempt to kill running server
kill:
	-@killall -9 $(PROG) 2>/dev/null || true
 
# attempt to build and start server
restart:
	@make kill
	@make build; (if [ "$$?" -eq 0 ]; then (env GIN_MODE=debug ./${PROG} &); fi)

publish: build
	cp ./${PROG} ~/bin