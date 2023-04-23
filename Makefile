# Makefile


# Linter variables
# Set linter version and binary location
LINTER_VERSION := v1.52.2
LINTER_BINARY := ./bin/golangci-lint

# Set the variables.
PROTO_DIR=./proto/ping
PROTO_FILE=$(PROTO_DIR)/pong.proto
GO_OUT_DIR=./proto/ping

# Generate the Go code from the protobuf file.
.PHONY: generate
generate:
	@echo "Generating protobuf Go code..."
	@protoc \
		-I $(PROTO_DIR) \
		--go_out=$(GO_OUT_DIR) \
		--go_opt=paths=source_relative \
		--go-grpc_out=$(GO_OUT_DIR) \
		--go-grpc_opt=paths=source_relative \
		$(PROTO_FILE)

# Clean the generated files.
.PHONY: clean
clean:
	@echo "Cleaning generated files..."
	@rm -f $(GO_OUT_DIR)/*.pb.go $(GO_OUT_DIR)/*.grpc.go
	@rm -rf bin
	@echo "Done."


# Install dependencies.
.PHONY: deps
deps:
	go mod download

# Run all the tests.
.PHONY: test
test:
	go test -race -cover -timeout 5m -v ./...

# Run linter
.PHONY: lint
lint: $(LINTER_BINARY)
	$(LINTER_BINARY) run ./...

# Install linter if necessary
$(LINTER_BINARY):
	@mkdir -p $(@D)
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh \
        | sh -s -- -b $(dir $@) $(LINTER_VERSION)
