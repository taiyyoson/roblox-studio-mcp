# output binary path; matches the /bin/ rule in .gitignore
BIN := bin/server

.PHONY: build run test test-race vet fmt tidy clean check

# compile the server to bin/server
build:
	go build -o $(BIN) ./cmd/server

# compile and run over stdio, no leftover binary
run:
	go run ./cmd/server

# run all tests
test:
	go test ./...

# sync go.mod/go.sum with imports
tidy:
	go mod tidy

# delete build output
clean:
	rm -rf bin

# pre-commit gate
check: fmt vet test

# run all tests with the race detector
test-race:
	go test -race ./...

# static analysis
vet:
	go vet ./...

# format all code in place
fmt:
	gofmt -w .