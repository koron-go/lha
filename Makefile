.PHONY: build
build:
	go build -v -i

.PHONY: test
test:
	go test ./...

.PHONY: test-full
test-full:
	go test -race ./...

.PHONY: tags
tags:
	gotags -f tags -R .
