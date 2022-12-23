test:
	go test -v ./...

fmt:
	go run golang.org/x/tools/cmd/goimports@v0.3.0 -w .

gomodtidy:
	go mod tidy -v

generate:
	go generate ./...

fieldalignment:
	go run golang.org/x/tools/go/analysis/passes/fieldalignment/cmd/fieldalignment@v0.4.0 -test=false $(shell go list ./... | grep -v modeldecoder/generator | grep -v test)

update-licenses:
	go run github.com/elastic/go-licenser@v0.4.1 .
