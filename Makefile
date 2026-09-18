.PHONY: tools build run fmt lint vet test coverage-html coverage-check check

COVERAGE_THRESHOLD := 70
COVERAGE_OUT := coverage.out
COVERAGE_HTML := coverage.html

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

test:
	go test ./...


# Generate the browsable HTML report
coverage-html:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "open coverage.html in your browser"

# Enforce a minimum threshold — fails (exit 1) if below it, for CI
coverage-check:
	go test ./... -coverprofile=coverage.out
	@total=$$(go tool cover -func=coverage.out | grep total | awk '{print $$3}' | tr -d '%'); \
	echo "total coverage: $$total%"; \
	if [ $$(echo "$$total < $(COVERAGE_THRESHOLD)" | bc) -eq 1 ]; then \
		echo "coverage $$total% is below threshold $(COVERAGE_THRESHOLD)%"; \
		exit 1; \
	fi

# coverage-check already runs the tests so no need to include "test" here
check: fmt vet lint coverage-check

