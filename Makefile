APP_NAME := netry
VERSION  := 0.0.1
BUILD_DIR := build

# Targets: OS/ARCH combinations
TARGETS = \
    linux/amd64 \
    linux/arm64 \
    windows/amd64 \
    darwin/amd64 \
    darwin/arm64

run:
	@go run main.go

rebuild: clean build

build:
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@$(foreach target,$(TARGETS), \
		GOOS=$(word 1,$(subst /, ,$(target))) \
		GOARCH=$(word 2,$(subst /, ,$(target))) \
		CGO_ENABLED=0 \
		go build -o $(BUILD_DIR)/$(APP_NAME)-$(subst /,-,$(target))$(if $(findstring windows,$(target)),.exe,) main.go && echo "Built for $(target)"; \
	)

clean:
	@echo "Cleaning build directory..."
	@rm -rf $(BUILD_DIR)

.PHONY: all clean build
