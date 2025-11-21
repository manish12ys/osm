.PHONY: build install clean package install-pkg uninstall help

BINARY_NAME=omarchy-monitor
INSTALL_PATH=/usr/bin/osm
GO=go

help:
	@echo "Omarchy System Monitor - Makefile"
	@echo ""
	@echo "Available targets:"
	@echo "  make build       - Build the binary"
	@echo "  make install     - Install binary to /usr/bin (requires sudo)"
	@echo "  make package     - Create Arch Linux package"
	@echo "  make install-pkg - Build and install Arch package"
	@echo "  make clean       - Remove built files"
	@echo "  make uninstall   - Remove installed binary"
	@echo "  make run         - Build and run the monitor"

build:
	@echo "Building $(BINARY_NAME)..."
	$(GO) build -o $(BINARY_NAME) .

run: build
	@echo "Starting $(BINARY_NAME)..."
	./$(BINARY_NAME)

install: build
	@echo "Installing to $(INSTALL_PATH)..."
	sudo install -Dm755 $(BINARY_NAME) $(INSTALL_PATH)
	@echo "Installed! Run with: osm"

package:
	@echo "Building Arch Linux package..."
	./build-pkg.sh

install-pkg: package
	@echo "Installing Arch Linux package..."
	sudo pacman -U omarchy-monitor*.pkg.tar.zst

uninstall:
	@echo "Removing $(INSTALL_PATH)..."
	sudo rm -f $(INSTALL_PATH)

clean:
	@echo "Cleaning build artifacts..."
	rm -f $(BINARY_NAME)
	rm -rf pkg/ src/
	rm -f *.pkg.tar.zst

# Development targets
dev: build
	@./$(BINARY_NAME)

rebuild: clean build
