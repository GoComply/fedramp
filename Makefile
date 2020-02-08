GO=GO111MODULE=on go
GOBUILD=$(GO) build

all: build

build: cli/pkged.go
	$(GOBUILD) ./cli/fedramp.go

cli/pkged.go: pkger README.md
	pkger -o cli

.PHONY: pkger
pkger:
ifeq ("$(wildcard $(GOPATH)/bin/pkger)","")
	go get -u -v github.com/markbates/pkger/cmd/pkger
endif

ci-update-fedramp-templates:
	rm bundled/templates/FedRAMP-SSP-*-Baseline-Template.docx
	wget -P bundled/templates/ https://www.fedramp.gov/assets/resources/templates/FedRAMP-SSP-High-Baseline-Template.docx https://www.fedramp.gov/assets/resources/templates/FedRAMP-SSP-Moderate-Baseline-Template.docx https://www.fedramp.gov/assets/resources/templates/FedRAMP-SSP-Low-Baseline-Template.docx
