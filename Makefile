.PHONY: test run update format install build relase

install:
	go get github.com/aktau/github-release
	go list -f '{{range .Imports}}{{.}} {{end}}' ./... | xargs go get -v
	go list -f '{{range .TestImports}}{{.}} {{end}}' ./... | xargs go get -v

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
	GOOS=linux  GOARCH=386 go build -o build/linux/386/gitlab-ci-helper cli/main.go
	GOOS=linux  GOARCH=arm go build -o build/linux/arm/gitlab-ci-helper cli/main.go
	GOOS=linux  GOARCH=arm64 go build -o build/linux/arm64/gitlab-ci-helper cli/main.go

release: build
	github-release delete  --tag master --user rande --repo gitlab-ci-helper|| exit 0
	github-release release --tag master --user rande --repo gitlab-ci-helper --name "Beta release" --pre-release
	github-release upload  --tag master --user rande --repo gitlab-ci-helper --name "gitlab-ci-helper-osx-amd64"   --file build/darwin/amd64/gitlab-ci-helper
	github-release upload  --tag master --user rande --repo gitlab-ci-helper --name "gitlab-ci-helper-linux-amd64" --file build/linux/amd64/gitlab-ci-helper
	github-release upload  --tag master --user rande --repo gitlab-ci-helper --name "gitlab-ci-helper-linux-386"   --file build/linux/386/gitlab-ci-helper
	github-release upload  --tag master --user rande --repo gitlab-ci-helper --name "gitlab-ci-helper-linux-arm"   --file build/linux/arm/gitlab-ci-helper
	github-release upload  --tag master --user rande --repo gitlab-ci-helper --name "gitlab-ci-helper-linux-arm64" --file build/linux/arm64/gitlab-ci-helper
