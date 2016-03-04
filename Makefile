.PHONY: test run update format install build relase

install:
	go list -f '{{range .Imports}}{{.}} {{end}}' ./... | xargs go get -v
	go list -f '{{range .TestImports}}{{.}} {{end}}' ./... | xargs go get -v
	go build -v ./...

update:
	go get -u all

run:
	go run cli/main.go

format:
	gofmt -l -w -s .
	go fix ./...
	go run cli/main.go dump:readme > commands.md

test:
	go test -v ./...
	go vet ./...
	gofmt -l -s -e .
	exit `gofmt -l -s -e . | wc -l`

build:
	GOOS=darwin GOARCH=amd64 go build -o build/darwin/amd64/gitlab-ci-helper cli/main.go
	GOOS=linux  GOARCH=amd64 go build -o build/linux/amd64/gitlab-ci-helper cli/main.go

release:
	github-release delete  --tag master --user rande --repo gitlab-ci-helper|| exit 0
	github-release release --tag master --user rande --repo gitlab-ci-helper --name "Beta release" --pre-release
	github-release upload  --tag master --user rande --repo gitlab-ci-helper --name "gitlab-ci-helper-osx-amd64"   --file build/darwin/amd64/gitlab-ci-helper
	github-release upload  --tag master --user rande --repo gitlab-ci-helper --name "gitlab-ci-helper-linux-amd64" --file build/linux/amd64/gitlab-ci-helper
