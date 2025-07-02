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

.PHONY: build run test lint clean inspector patch minor major

define bump_version
	@latest=$$(git describe --tags --abbrev=0); \
	base=$$(echo $$latest | sed 's/^v//'); \
	IFS='.' read -r major minor patch <<< "$$base"; \
	if [ "$(1)" = "major" ]; then \
		major=$$((major+1)); minor=0; patch=0; \
	elif [ "$(1)" = "minor" ]; then \
		minor=$$((minor+1)); patch=0; \
	else \
		patch=$$((patch+1)); \
	fi; \
	new_tag="v$$major.$$minor.$$patch"; \
	echo "Tagging $$new_tag"; \
	git tag -a "$$new_tag" -m "Bump version to $$new_tag"; \
	git push origin "$$new_tag"
endef

patch:
	$(call bump_version,patch)

minor:
	$(call bump_version,minor)

major:
	$(call bump_version,major)
