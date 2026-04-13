# unit tests. $R=1 for race flag
.PHONY: test
test:
	go test $(if $(R),-race,) ./...
