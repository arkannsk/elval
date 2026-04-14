.PHONY: build generate clean test install gen

BINARY=bin/elval-gen

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

# unit tests. $R=1 for race flag
test: gen
	go test $(if $(R),-race,) ./...

clean:
	@rm -rf bin
	@find ./test/integration -name "*.gen.go" -delete
