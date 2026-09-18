.PHONY: tools build run fmt lint vet check

tools:
	go install golang.org/x/tools/cmd/goimports@latest
ifeq ($(shell uname), Darwin)
	brew install golangci-lint
else
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin latest
endif

build:
	go build -o tfctl .

run:
	go run .

fmt:
	goimports -w .

lint:
	golangci-lint run

vet:
	go vet ./...

check: fmt vet lint