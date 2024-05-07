# Paths to the source code and build artifacts
SRC_PATH=.
DIST_PATH=./dist

# Nme of the binary and Docker image
BINARY_NAME=substrate-publisher

# Build flags for go build
BUILD_FLAGS=-ldflags="-s -w"

.PHONY: build
build:
	# Build the production binary
	CGO_ENABLED=0 go build $(BUILD_FLAGS) -o $(DIST_PATH)/$(BINARY_NAME) $(SRC_PATH)

.PHONY: build-debug
build-debug:
	# Build the debug binary
	go build -o $(DIST_PATH)/$(BINARY_NAME) $(SRC_PATH)
