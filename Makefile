.PHONY: build generate clean test install gen gen-test bench bench-all

BINARY=bin/elval-gen
test_invalid_spec_dirs=test/lint/invalid,pkg/parser/testdata/invalid

lint:
	golangci-lint run ./...

lint-spec: install
	elval-gen lint -v -exclude $(test_invalid_spec_dirs)

install:
	go install ./cmd/elval-gen

build:
	@mkdir -p bin
	@go build -o $(BINARY) ./cmd/elval-gen

gen: install
	go generate ./...

gen-spec:
	go run ./cmd/elval-gen gen -i ./..

# unit tests. R=1 for race flag, C=1 for cover
test: gen-spec
	go test $(if $(R),-race,) $(if $(C),-cover,) ./...

# benchmark
bench:
	@go test -bench=. -benchmem ./test/benchmark

# benchmark with longer time (10s)
bench-long:
	@go test -bench=. -benchtime=10s -benchmem ./test/benchmark

# benchmark with memory profile
bench-mem:
	@go test -bench=. -benchmem -memprofile=mem.out -cpuprofile=cpu.out ./test/benchmark

clean:
	@rm -rf bin
	@find ./ -name "*.debug.go" -delete

docs:
	godocdown ./pkg/generator/utils > pkg/generator/utils/README.md
	godocdown ./ > ./elval.README.md
