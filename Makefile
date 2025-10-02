GOCMD=go
BINARY_NAME=hbsdsrv-build
SERVICE_PATH=/usr/local/etc/rc.d/hbsdsrv-builder
BUILD_DIR=bin
SRC=.

all: clean build install service

clean:
	@echo "Cleaning..."
	$(GOCMD) clean
	@rm -rfv $(BUILD_DIR)

build: clean
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOCMD) mod tidy
	$(GOCMD) build -o $(BUILD_DIR)/$(BINARY_NAME) $(SRC)

install: build
	@echo "Installing to /usr/local/bin/$(BINARY_NAME)"
	/usr/bin/install $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/$(BINARY_NAME)

service: install
	@echo "Installing new Service File"
	@cp -vf ./hbsdsrv_builder.service $(SERVICE_PATH)
	@chmod +x $(SERVICE_PATH)

.PHONY: all clean build install service
