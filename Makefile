BIN := marky
SRC=$(shell find . -name "*.go")

# Build the application
build:
	$(info ******************** building ${BIN} ********************)
	@go build -o bin/${BIN} cmd/${BIN}/main.go

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

.PHONY: build run test lint clean