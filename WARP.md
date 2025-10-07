# WARP.md

This file provides guidance to WARP (warp.dev) when working with code in this repository.

## Project Overview

The Minecraft Instance Manager is a lightweight bash-based tool that allows users to manage multiple Minecraft installations using symlinks. It enables instant switching between different Minecraft setups (modpacks, versions, configurations) without copying files or wasting disk space.

## Core Architecture

### Symlink-Based Design
The system works by:
1. Storing instances in `~/.minecraft-instances/`
2. Creating symlinks from `~/.minecraft` to the active instance
3. Backing up the original `.minecraft` directory to `.minecraft.backup`

### Directory Structure
```
~/.minecraft-instances/
├── vanilla/           # Clean Minecraft instance
├── modpack-1.20.1/   # Modded instance
└── testing/          # Development instance
```

### Key Components
- **Main Script**: `minecraft-instances` - The core bash script containing all functionality
- **Installation Script**: `install.sh` - Handles system-wide or user-local installation
- **Documentation**: `README.md` and `examples/USAGE_EXAMPLES.md`

## Common Commands

### Development and Testing
```bash
# Test the script locally (make executable first)
chmod +x minecraft-instances
./minecraft-instances list

# Test all core functions
./minecraft-instances create test-instance
./minecraft-instances switch test-instance
./minecraft-instances list
./minecraft-instances restore

# Clean up test instance
rm -rf ~/.minecraft-instances/test-instance
```

### Installation
```bash
# Install system-wide (requires sudo)
sudo ./install.sh

# Install for current user only
./install.sh

# Manual installation to custom location
cp minecraft-instances ~/bin/
chmod +x ~/bin/minecraft-instances
```

### Code Quality
```bash
# Check bash syntax
bash -n minecraft-instances

# Run with debug output
bash -x minecraft-instances list

# Check for common bash issues with shellcheck (if available)
shellcheck minecraft-instances
shellcheck install.sh
```

## Script Architecture

### Function Organization
The main script follows a simple pattern:
1. **Configuration Variables** - Define paths and constants
2. **Core Functions** - Each command is a separate function
3. **Command Router** - Case statement dispatching to functions

### Key Functions
- `create_instance()` - Creates new instance by copying current .minecraft
- `switch_instance()` - Backs up current .minecraft and creates symlink
- `list_instances()` - Shows instances with mod counts and current status
- `restore_default()` - Removes symlink and restores original .minecraft

### Safety Mechanisms
- Always backs up current .minecraft before switching
- Validates instance exists before switching
- Uses symlinks (non-destructive, easily reversible)
- Provides restore functionality

## Development Guidelines

### Code Style
- Use bash best practices (proper quoting, error handling)
- Include helpful echo messages for user feedback
- Use local variables in functions
- Check for required arguments and provide usage help

### Testing Approach
- Test with different Minecraft directory states (exists/doesn't exist)
- Test error conditions (invalid instance names, missing directories)
- Verify symlink creation and backup functionality
- Test mod counting accuracy

### File Structure Expectations
- Main script should remain self-contained
- Installation script should handle both system and user installs
- Documentation should be comprehensive but concise

## Platform Considerations

### Linux/macOS
- Uses standard Unix tools (ln, mv, cp, find)
- Symlink behavior is consistent
- Path handling uses standard shell expansion

### Windows/WSL
- May require different symlink handling
- Path separators and permissions might differ
- Consider Windows Minecraft launcher behavior

## Instance Management Patterns

### Development Workflow
- Create clean test instances for mod development
- Use descriptive naming (e.g., `dev-1.20.1`, `compatibility-test`)
- Keep separate instances for different Minecraft versions

### Backup Strategy
- Original .minecraft is always preserved as .minecraft.backup
- Create dated backup instances before major changes
- Use tar/zip for sharing instances between systems

### Mod Organization
- Each instance maintains its own mods/ directory
- Mod counts are displayed for quick reference
- Easy to add/remove mods per instance