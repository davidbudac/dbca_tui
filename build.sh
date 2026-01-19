#!/bin/bash
#
# Build script for DBCA TUI
# Cross-compiles for multiple platforms
#

APP_NAME="dbca_tui"
VERSION="${VERSION:-1.0.0}"
OUTPUT_DIR="dist"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "========================================"
echo "  DBCA TUI Build Script"
echo "  Version: ${VERSION}"
echo "========================================"
echo ""

# Create output directory
mkdir -p "${OUTPUT_DIR}"

SUCCESS_COUNT=0
FAIL_COUNT=0

# Function to build for a specific platform
build() {
    local os=$1
    local arch=$2
    local suffix=$3
    local output="${OUTPUT_DIR}/${APP_NAME}-${os}-${arch}${suffix}"

    echo -n "Building for ${os}/${arch}... "

    if GOOS="${os}" GOARCH="${arch}" go build -ldflags="-s -w" -o "${output}" . 2>/dev/null; then
        local size=$(ls -lh "${output}" | awk '{print $5}')
        echo -e "${GREEN}OK${NC} (${size})"
        SUCCESS_COUNT=$((SUCCESS_COUNT + 1))
    else
        echo -e "${RED}FAILED${NC}"
        FAIL_COUNT=$((FAIL_COUNT + 1))
    fi
}

echo "Building binaries..."
echo ""

# macOS ARM64 (Apple Silicon)
build darwin arm64 ""

# macOS AMD64 (Intel)
build darwin amd64 ""

# Linux AMD64
build linux amd64 ""

# Linux ARM64
build linux arm64 ""

# Linux x86 (32-bit)
build linux 386 ""

# Windows AMD64
build windows amd64 ".exe"

# Windows x86 (32-bit)
build windows 386 ".exe"

# Solaris AMD64
build solaris amd64 ""

# FreeBSD AMD64
build freebsd amd64 ""

# Note about AIX
echo ""
echo -e "${YELLOW}Note: AIX (ppc64) is not supported due to clipboard library limitations${NC}"

echo ""
echo "========================================"
echo "  Build complete!"
echo "  Success: ${SUCCESS_COUNT}, Failed: ${FAIL_COUNT}"
echo "========================================"
echo ""
echo "Output files in ${OUTPUT_DIR}/:"
ls -lh "${OUTPUT_DIR}/"

echo ""
echo "Platform notes:"
echo "  - darwin-arm64:   macOS on Apple Silicon (M1/M2/M3)"
echo "  - darwin-amd64:   macOS on Intel"
echo "  - linux-amd64:    64-bit x86 Linux"
echo "  - linux-arm64:    ARM64 Linux (AWS Graviton, Raspberry Pi 4)"
echo "  - linux-386:      32-bit x86 Linux"
echo "  - windows-amd64:  64-bit Windows"
echo "  - windows-386:    32-bit Windows"
echo "  - solaris-amd64:  Oracle Solaris on x86-64"
echo "  - freebsd-amd64:  FreeBSD on x86-64"
