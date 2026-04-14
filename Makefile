.PHONY: build generate clean test

BINARY=bin/elval-gen

build:
	@mkdir -p bin
	@go build -o $(BINARY) ./cmd/elval-gen

# gen $V=1 for verbose flag
gen: build
	@./$(BINARY) -input ./test/integration/person $(if $(V),-v,)
	@./$(BINARY) -input ./test/integration/product $(if $(V),-v,)

# unit tests. $R=1 for race flag
test: gen
	go test $(if $(R),-race,) ./...

clean:
	@rm -rf bin
	@find ./test/integration -name "*.gen.go" -delete
