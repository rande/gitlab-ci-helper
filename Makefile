.PHONY: test run update format install

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

test:
	go test -v ./...
	go vet ./...
	gofmt -l -s -e .
	exit `gofmt -l -s -e . | wc -l`
