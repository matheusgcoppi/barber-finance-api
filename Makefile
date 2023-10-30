BINARY_NAME=go_project_structure
MAIN_PACKAGE_PATH := ./cmd/go_project_structure

build:
	@go build -o ~/bin/barber-finance/$(BINARY_NAME) $(MAIN_PACKAGE_PATH)

run: build
	@~/bin/barber-finance/$(BINARY_NAME)

test:
	@go test -v ./tests