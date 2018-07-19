.PHONY: test run update format install build relase
.DEFAULT_GOAL := default


SHA1     = $(shell git rev-parse HEAD)
GO_PKG   = ./,./commands,./integrations/flowdock,./integrations/hipchat
GO_FILES = $(shell find $(GO_PROJECTS_PATHS) -maxdepth 1 -type f -name "*.go")
OS       = $(shell uname)

# COLORS
RED    = $(shell printf "\33[31m")
GREEN  = $(shell printf "\33[32m")
WHITE  = $(shell printf "\33[37m")
YELLOW = $(shell printf "\33[33m")
RESET  = $(shell printf "\33[0m")

help: ## prints help
	@cat $(MAKEFILE_LIST) | grep -e "^[a-zA-Z_\-]*: *.*## *" | awk 'BEGIN {FS = ":.*?## "}; {printf " > \033[36m%-20s\033[0m %s\n", $$1, $$2}'

default: test build ## test and build binaries

install: ## install dependencies
	go list -f '{{range .Imports}}{{.}} {{end}}' ./... | xargs go get -v
	go list -f '{{range .TestImports}}{{.}} {{end}}' ./... | xargs go get -v
	go get -v github.com/wadey/gocovmerge

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
    ifneq ("", "$(shell which docker)")
		docker run --rm -v $(shell pwd):/usr/src/myapp -v $(GOPATH):/usr/src/myapp/vendor -w /usr/src/myapp -e "GOPATH=/usr/src/myapp/vendor:/go" -e GOOS=linux -e GOARCH=amd64 golang:1.9-alpine go build -ldflags "-X main.RefLog=$(SHA1) -s -w" -o build/alpine-amd64-gitlab-ci-helper cli/main.go
    endif

build_checksums: ## generate checksums for binaries
	@echo "${YELLOW}Generating CLI build checksums${RESET}"
	@rm -f build/checksums.txt

    # for OSX users where you have md5 instead of md5sum
    ifeq (${OS}, Darwin)
        # md5 output has the following format:
        #
        # MD5 (darwin-386-gitlab-ci-helper) = 8eb317789e5d08e1c800cc469c20325a
        #
        # that's why sed and awk are used to cleanup
		@cd build && ls . | grep gitlab-ci-helper \
            | xargs md5 \
            | awk '{ printf("%s\n%s\n\n", $$2, $$4) }' \
            | sed 's/[()]//g' \
            >> checksums.txt
    else
        # md5sum output has the following format:
        #
        # 8eb317789e5d08e1c800cc469c20325a darwin-386-gitlab-ci-helper
        #
        # that's why awk is used to cleanup
		@cd build && ls . | grep gitlab-ci-helper \
            | xargs md5sum \
            | awk '{ printf("%s\n%s\n\n", $$2, $$1) }' \
            >> checksums.txt
    endif
	@echo "${GREEN}✔ successfully generated CLI build checksums to ${WHITE}build/checksums.txt${RESET}\n"

coverage-backend: ## run coverage tests
	mkdir -p build/coverage
	rm -rf build/coverage/*.cov
	go test -v -timeout 60s -coverpkg $(GO_PKG) -covermode count -coverprofile=build/coverage/main.cov ./
	go test -v -timeout 60s -coverpkg $(GO_PKG) -covermode count -coverprofile=build/coverage/commands.cov ./commands
	go test -v -timeout 60s -coverpkg $(GO_PKG) -covermode count -coverprofile=build/coverage/integration_flowdock.cov ./integrations/flowdock
	go test -v -timeout 60s -coverpkg $(GO_PKG) -covermode count -coverprofile=build/coverage/integration_hipchat.cov ./integrations/hipchat
	gocovmerge build/coverage/* > build/gitlabcihelper.coverage
	go tool cover -html=./build/gitlabcihelper.coverage -o build/gitlabcihelper.html
