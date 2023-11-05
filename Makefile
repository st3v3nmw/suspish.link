.PHONY: format
format:
	gofmt -w .

.PHONY: test
test:
	go test ./...  -coverpkg=./... -coverprofile ./coverage.out
	go tool cover -func ./coverage.out

.PHONY: build
build:
	go build -o suspish ./...

.PHONY: run
run: build
	./suspish --verbose
