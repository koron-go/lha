GO_SUBPKGS = $(shell go list ./... | grep -v /vendor/ | sed -e "s!$$(go list)!.!")

default: test

test:
	go test $(GO_SUBPKGS)

test-full:
	go test -v -race $(GO_SUBPKGS)

vet:
	@echo "go vet"
	@go vet $(GO_SUBPKGS)
	@echo ""

lint:
	@echo "golint"
	@for f in $(GO_SUBPKGS) ; do golint $$f ; done
	@echo ""

cyclo:
	-gocyclo -top 10 -avg $(GO_SUBPKGS)
	@echo ""

cyclo-report:
	@echo gocyclo -over 14 -avg
	-@gocyclo -over 14 -avg $(GO_SUBPKGS)
	@echo ""

misspell:
	@echo misspell
	@find $(GO_SUBPKGS) -maxdepth 1 -type f | xargs misspell
	@echo ""

report: misspell cyclo-report vet lint

deps:
	go get -v -u -d -t ./...

tags:
	ctags *.go

clean:
	go clean

.PHONY: test test-full lint cyclo report deps tags clean
