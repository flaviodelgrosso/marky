BIN := marky
BIN_MCP := marky-mcp
SRC=$(shell find . -name "*.go")

# Build the application
build:
	$(info ******************** building ${BIN} ********************)
	@go build -o bin/${BIN} cmd/${BIN}/main.go

build-mcp:
	$(info ******************** building ${BIN_MCP} ********************)
	@go build -o bin/${BIN_MCP} marky-mcp/main.go

# Run the application
run:
	@go run cmd/${BIN}/main.go $(ARGS)

# Run tests
test:
	$(info ******************** running tests ********************)
	@go test -v ./...

lint:
	$(info ******************** running lint tools ********************)
	golangci-lint run -v

# Clean the binary
clean:
	$(info ******************** cleaning up ********************)
	@rm -f bin/${BIN}
	@rm -f bin/${BIN_MCP}

inspector:
	$(info ******************** running inspector ********************)
	@npx @modelcontextprotocol/inspector go run mcp/main.go

.PHONY: build run test lint clean