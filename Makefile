.PHONY: test run update format install build relase
.DEFAULT_GOAL := default

SHA1=$(shell git rev-parse HEAD)

help: ## prints help
	@cat $(MAKEFILE_LIST) | grep -e "^[a-zA-Z_\-]*: *.*## *" | awk 'BEGIN {FS = ":.*?## "}; {printf " > \033[36m%-20s\033[0m %s\n", $$1, $$2}'

default: test build ## test and build binaries

install: ## install dependencies
	go get github.com/aktau/github-release
	go list -f '{{range .Imports}}{{.}} {{end}}' ./... | xargs go get -v
	go list -f '{{range .TestImports}}{{.}} {{end}}' ./... | xargs go get -v

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
	GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.RefLog=$(SHA1)" -o build/darwin/amd64/gitlab-ci-helper cli/main.go
	GOOS=linux  GOARCH=amd64 go build -ldflags "-X main.RefLog=$(SHA1)" -o build/linux/amd64/gitlab-ci-helper  cli/main.go
	GOOS=linux  GOARCH=386   go build -ldflags "-X main.RefLog=$(SHA1)" -o build/linux/386/gitlab-ci-helper    cli/main.go
	GOOS=linux  GOARCH=arm   go build -ldflags "-X main.RefLog=$(SHA1)" -o build/linux/arm/gitlab-ci-helper    cli/main.go
	GOOS=linux  GOARCH=arm64 go build -ldflags "-X main.RefLog=$(SHA1)" -o build/linux/arm64/gitlab-ci-helper  cli/main.go
	build/linux/amd64/gitlab-ci-helper version -e

release: build ## build and release binaries on github
	github-release delete  --tag master --user rande --repo gitlab-ci-helper|| exit 0
	github-release release --tag master --user rande --repo gitlab-ci-helper --name "Beta release" --pre-release
	github-release upload  --tag master --user rande --repo gitlab-ci-helper --name "gitlab-ci-helper-osx-amd64"   --file build/darwin/amd64/gitlab-ci-helper
	github-release upload  --tag master --user rande --repo gitlab-ci-helper --name "gitlab-ci-helper-linux-amd64" --file build/linux/amd64/gitlab-ci-helper
	github-release upload  --tag master --user rande --repo gitlab-ci-helper --name "gitlab-ci-helper-linux-386"   --file build/linux/386/gitlab-ci-helper
	github-release upload  --tag master --user rande --repo gitlab-ci-helper --name "gitlab-ci-helper-linux-arm"   --file build/linux/arm/gitlab-ci-helper
	github-release upload  --tag master --user rande --repo gitlab-ci-helper --name "gitlab-ci-helper-linux-arm64" --file build/linux/arm64/gitlab-ci-helper
