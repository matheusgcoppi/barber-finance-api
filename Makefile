BINARY_NAME=goProjectStructure
MAIN_PACKAGE_PATH := ./cmd/goProjectStructure

build:
	@go build -o ~/bin/barber-finance/$(BINARY_NAME) $(MAIN_PACKAGE_PATH)

run: build
	@~/bin/barber-finance/$(BINARY_NAME)

test:
	@go test -v ./...