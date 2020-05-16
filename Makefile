.PHONY: build
build:
	go build -gcflags '-e'

.PHONY: test
test:
	go test ./...

.PHONY: test-full
test-full:
	go test -race ./...

.PHONY: tags
tags:
	gotags -f tags -R .

.PHONY: lint
lint:
	golint ./...

.PHONY: clean
clean:
	go clean
	rm -f tags

# based on: github.com/koron-go/_skeleton/Makefile
