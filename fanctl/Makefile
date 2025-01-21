
GOOS := linux
GOARCH := amd64

.PHONY: build
build:
	@echo "Building..."
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build \
		-o fanctl .
