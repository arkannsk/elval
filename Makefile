.PHONY: build generate clean test install gen gen-test bench bench-all

BINARY=bin/elval-gen

lint:
	golangci-lint run ./...

lint-spec: install
	elval-gen lint

install:
	go install ./cmd/elval-gen

build:
	@mkdir -p bin
	@go build -o $(BINARY) ./cmd/elval-gen

gen: install
	go generate ./...

# gen $V=1 for verbose flag
gen-test: build
	@./$(BINARY) -input ./test/integration/person $(if $(V),-v,)
	@./$(BINARY) -input ./test/integration/product $(if $(V),-v,)
	@./$(BINARY) -input ./test/integration/mixed $(if $(V),-v,)

# unit tests. R=1 for race flag, C=1 for cover
test: gen
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
	@find ./test/integration -name "*.gen.go" -delete
