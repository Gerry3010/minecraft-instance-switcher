#!/bin/bash

# Minecraft Instance Manager - Installation Script
# This script installs the minecraft-instances tool system-wide

set -e

INSTALL_DIR="/usr/local/bin"
SCRIPT_NAME="minecraft-instances"
CURRENT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if we're in the right directory
if [ ! -f "$CURRENT_DIR/$SCRIPT_NAME" ]; then
    print_error "Cannot find $SCRIPT_NAME in current directory!"
    print_info "Please run this script from the minecraft-instance-manager directory."
    exit 1
fi

print_info "Minecraft Instance Manager Installation"
echo ""

# Check if running as root for system-wide install
if [ "$EUID" -eq 0 ]; then
    print_info "Installing system-wide to $INSTALL_DIR"
    
    # Copy script to system directory
    cp "$CURRENT_DIR/$SCRIPT_NAME" "$INSTALL_DIR/"
    chmod +x "$INSTALL_DIR/$SCRIPT_NAME"
    
    print_success "Installed to $INSTALL_DIR/$SCRIPT_NAME"
    print_info "You can now run 'minecraft-instances' from anywhere!"
    
else
    print_warning "Not running as root - installing to user bin directory"
    
    # Create user bin directory if it doesn't exist
    USER_BIN="$HOME/.local/bin"
    mkdir -p "$USER_BIN"
    
    # Copy script to user directory
    cp "$CURRENT_DIR/$SCRIPT_NAME" "$USER_BIN/"
    chmod +x "$USER_BIN/$SCRIPT_NAME"
    
    print_success "Installed to $USER_BIN/$SCRIPT_NAME"
    
    # Check if user bin is in PATH
    if [[ ":$PATH:" == *":$USER_BIN:"* ]]; then
        print_success "You can now run 'minecraft-instances' from anywhere!"
    else
        print_warning "$USER_BIN is not in your PATH"
        print_info "Add this line to your ~/.bashrc or ~/.zshrc:"
        echo "export PATH=\"\$PATH:$USER_BIN\""
        print_info "Or run the script directly: $USER_BIN/minecraft-instances"
    fi
fi

echo ""
print_info "Quick start:"
echo "  minecraft-instances create vanilla"
echo "  minecraft-instances create modpack"
echo "  minecraft-instances list"
echo "  minecraft-instances switch vanilla"

echo ""
print_success "Installation complete! Happy mining! ⛏️"