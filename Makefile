test:
	go test -v ./...

fmt:
	go run golang.org/x/tools/cmd/goimports@v0.3.0 -w .

gomodtidy:
	go mod tidy -v

generate:
	go generate ./...

update-licenses:
	go run github.com/elastic/go-licenser@v0.4.1 .
