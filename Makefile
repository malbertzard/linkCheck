# Name of the binary output
BINARY_NAME=linkcheck

# Main package of your application
MAIN_PACKAGE=./cmd/

# Default target
all: build

build: $(GOFILES)
	@echo "Building the application..."
	@go build -o ./bin/$(BINARY_NAME) $(MAIN_PACKAGE)

tidy: 
	@echo "Tidy up source code..."
	@go fmt $(MAIN_PACKAGE)

run: build
	@echo "Running the application..."
	./bin/$(BINARY_NAME)

clean:
	@echo "Cleaning up..."
	@rm -f ./bin/$(BINARY_NAME)

# Clean, build, and run the application
rebuild: clean build run

.PHONY: all build run clean rebuild tidy

