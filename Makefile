SHELL := /bin/bash

VERSION := `git describe --tags`
GITCOMMIT := `git rev-parse HEAD`
BRANCH := `git branch`
BUILDDATE := `date +%Y-%m-%d`
BUILDUSER := `whoami`

LDFLAGSSTRING :=-X main.Version=$(VERSION)
LDFLAGSSTRING +=-X main.GitCommit=$(GITCOMMIT)
LDFLAGSSTRING +=-X main.Branch=$(BRANCH)
LDFLAGSSTRING +=-X main.BuildDate=$(BUILDDATE)
LDFLAGSSTRING +=-X main.BuildUser=$(BUILDUSER)

LDFLAGS :=-ldflags "$(LDFLAGSSTRING)"

.PHONY: all build

all: build

# Build binary
build:
	go build $(LDFLAGS) 

test:
	go test -v ./...