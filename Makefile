# Makefile

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


# Install dependencies.
.PHONY: deps
deps:
	go mod download

# Run all the tests.
.PHONY: test
test:
	go test -race -cover -timeout 5m -v ./...
