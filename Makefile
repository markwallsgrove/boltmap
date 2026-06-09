GO := go
BINARY := boltmap

.PHONY: build test fmt vet clean

build:
	$(GO) build -o $(BINARY) ./cmd/boltmap

test:
	$(GO) test ./...

fmt:
	$(GO) fmt ./...

vet:
	$(GO) vet ./...

clean:
	rm -f $(BINARY)
