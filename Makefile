GO=GO111MODULE=on go
GOBUILD=$(GO) build

all: build

build: cli/pkged.go
	$(GOBUILD) ./cli/oscal-fedramp-templater.go

cli/pkged.go: pkger README.md
	pkger -o cli

.PHONY: pkger
pkger:
ifeq ("$(wildcard $(GOPATH)/bin/pkger)","")
	go get -u -v github.com/markbates/pkger/cmd/pkger
endif
