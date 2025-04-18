.PHONY: install build run clean

# Install dependencies, build the program, and install it to /usr/local/bin
install: build
	@echo "Installing dependencies..."
	@go mod download
	@echo "Building the program..."
	@go build -o myprogram ./...
	@echo "Installing program to /usr/local/bin (sudo required)..."
	sudo install -m 0755 myprogram /usr/local/bin/myprogram

# Build the program
build:
	@go build -o myprogram ./...

# Run the program
run: build
	@./myprogram

# Clean up build artifacts
clean:
	@rm -f myprogram
