#!/usr/bin/env make -f

PROJECTNAME := "ngraph"
# Make is verbose in Linux. Make it silent.
MAKEFLAGS += --silent

## test: Run the go test command.
.PHONY: test
test:
	go test -v ./...

## build: Compile the binary.
.PHONY: build
build:
	@mkdir -p bin
	env GOOS=linux GOARCH=amd64 go build -o bin/$(PROJECTNAME) cmd/ngraph/main.go

## clean: Cleanup binary.
clean:
	@rm -f bin/$(PROJECTNAME)

## help: Show this message.
.PHONY: help
help: Makefile
	@echo "Available targets:"
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
