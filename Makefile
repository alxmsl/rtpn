.PHONY: build
build:
	go build -buildvcs=false -o ./bin/ ./cmd/cpn/...

.PHONY: bench
bench:
	go test ./... -check.f='!Test' -bench=. -benchmem

.PHONY: fmt
fmt:
	gofmt -l -w `find . -type f -name '*.go'`
	goimports -l -w `find . -type f -name '*.go'`

.PHONY: plugin
plugin:
	go build -buildvcs=false -buildmode=plugin -o ./bin/ ./examples/plugin/

.PHONY: test
test:
	go test ./... -v -race

%:
	@:
