.PHONY: test run update format install build relase
.DEFAULT_GOAL := default

SHA1=$(shell git rev-parse HEAD)
GO_PKG = ./,./commands,./integrations/flowdock,./integrations/hipchat
GO_FILES = $(shell find $(GO_PROJECTS_PATHS) -maxdepth 1 -type f -name "*.go")

help: ## prints help
	@cat $(MAKEFILE_LIST) | grep -e "^[a-zA-Z_\-]*: *.*## *" | awk 'BEGIN {FS = ":.*?## "}; {printf " > \033[36m%-20s\033[0m %s\n", $$1, $$2}'

default: test build ## test and build binaries

install: ## install dependencies
	go list -f '{{range .Imports}}{{.}} {{end}}' ./... | xargs go get -v
	go list -f '{{range .TestImports}}{{.}} {{end}}' ./... | xargs go get -v
	go get github.com/wadey/gocovmerge

update: ## update dependencies
	go get -u all

run: ## run the command
	go run cli/main.go

format: ## format the code and generate commands.md file
	gofmt -l -w -s .
	go fix ./...
	go run cli/main.go dump:readme > commands.md

test: ## run tests and cs tools
	go test -v ./...
	go vet ./...
	gofmt -l -s -e .
	exit `gofmt -l -s -e . | wc -l`

build: ## build binaries
	GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.RefLog=$(SHA1) -s -w" -o build/darwin-amd64-gitlab-ci-helper cli/main.go
	GOOS=linux  GOARCH=amd64 go build -ldflags "-X main.RefLog=$(SHA1) -s -w" -o build/linux-amd64-gitlab-ci-helper  cli/main.go
	GOOS=linux  GOARCH=386   go build -ldflags "-X main.RefLog=$(SHA1) -s -w" -o build/linux-386-gitlab-ci-helper    cli/main.go
	GOOS=linux  GOARCH=arm   go build -ldflags "-X main.RefLog=$(SHA1) -s -w" -o build/linux-arm-gitlab-ci-helper    cli/main.go
	GOOS=linux  GOARCH=arm64 go build -ldflags "-X main.RefLog=$(SHA1) -s -w" -o build/linux-arm64-gitlab-ci-helper  cli/main.go

coverage-backend: ## run coverage tests
	mkdir -p build/coverage
	rm -rf build/coverage/*.cov
	go test -v -timeout 60s -coverpkg $(GO_PKG) -covermode count -coverprofile=build/coverage/main.cov ./
	go test -v -timeout 60s -coverpkg $(GO_PKG) -covermode count -coverprofile=build/coverage/commands.cov ./commands
	go test -v -timeout 60s -coverpkg $(GO_PKG) -covermode count -coverprofile=build/coverage/integration_flowdock.cov ./integrations/flowdock
	go test -v -timeout 60s -coverpkg $(GO_PKG) -covermode count -coverprofile=build/coverage/integration_hipchat.cov ./integrations/hipchat
	gocovmerge build/coverage/* > build/gitlabcihelper.coverage
	go tool cover -html=./build/gitlabcihelper.coverage -o build/gitlabcihelper.html
