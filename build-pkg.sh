#!/bin/bash
# Build Arch Linux package for Omarchy System Monitor

set -e

echo "🔧 Building Arch Linux package for Omarchy System Monitor..."

# Check if we're on Arch-based system
if ! command -v makepkg &> /dev/null; then
    echo "❌ Error: makepkg not found. This script requires an Arch-based system."
    exit 1
fi

# Clean previous builds
echo "🧹 Cleaning previous builds..."
rm -rf pkg/ src/ *.pkg.tar.zst

# Check if Go is available (command or binary)
if command -v go &> /dev/null; then
    # Build binary first, then package it
    echo "📦 Building binary with Go..."
    go build -o omarchy-monitor .
    
    echo "📦 Creating Arch package from binary..."
    makepkg -f -p PKGBUILD.bin
else
    echo "❌ Error: Go not found in PATH."
    echo "Please install Go or ensure it's in your PATH."
    echo ""
    echo "To install Go on Arch: sudo pacman -S go"
    exit 1
fi

# Find the built package
PKG_FILE=$(ls -t omarchy-monitor*.pkg.tar.zst 2>/dev/null | head -n1)

if [ -z "$PKG_FILE" ]; then
    echo "❌ Error: Package file not found!"
    exit 1
fi

echo ""
echo "✅ Package built successfully: $PKG_FILE"
echo ""
echo "To install the package, run:"
echo "  sudo pacman -U $PKG_FILE"
