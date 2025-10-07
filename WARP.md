# WARP.md

This file provides guidance to WARP (warp.dev) when working with code in this repository.

## Project Overview

The Minecraft Instance Manager is a modern Go application that provides both a beautiful terminal user interface (TUI) and command-line interface (CLI) for managing multiple Minecraft installations using symlinks. It enables instant switching between different Minecraft setups (modpacks, versions, configurations) without copying files or wasting disk space.

**Tech Stack**: Go, Cobra (CLI framework), Viper (configuration), Bubble Tea (TUI framework), Lip Gloss (styling)

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
- **Main Application**: `cmd/minecraft-instance-manager/main.go` - Entry point with Cobra CLI setup
- **CLI Commands**: `cmd/minecraft-instance-manager/commands.go` - Cobra command implementations
- **Instance Manager**: `internal/instance/manager.go` - Core business logic for instance management
- **TUI Interface**: `internal/tui/` - Bubble Tea terminal user interface
- **GitHub Actions**: `.github/workflows/build-and-release.yml` - CI/CD pipeline for cross-platform builds
- **Documentation**: `README.md`, `examples/USAGE_EXAMPLES.md`, and `WARP.md`

## Common Commands

### Development and Testing
```bash
# Build the application
go build -o minecraft-instance-manager ./cmd/minecraft-instance-manager

# Run tests
go test ./...

# Test the application locally
./minecraft-instance-manager list

# Test all core functions
./minecraft-instance-manager create test-instance
./minecraft-instance-manager switch test-instance
./minecraft-instance-manager list
./minecraft-instance-manager restore

# Test TUI mode
./minecraft-instance-manager

# Clean up test instance
./minecraft-instance-manager delete test-instance
```

### Build and Release
```bash
# Build for current platform
go build -o minecraft-instance-manager ./cmd/minecraft-instance-manager

# Build for all platforms (requires Go 1.21+)
GOOS=linux GOARCH=amd64 go build -o dist/minecraft-instance-manager-linux-amd64 ./cmd/minecraft-instance-manager
GOOS=windows GOARCH=amd64 go build -o dist/minecraft-instance-manager-windows-amd64.exe ./cmd/minecraft-instance-manager
GOOS=darwin GOARCH=amd64 go build -o dist/minecraft-instance-manager-darwin-amd64 ./cmd/minecraft-instance-manager

# Install locally
go install ./cmd/minecraft-instance-manager
```

### Code Quality
```bash
# Run Go vet
go vet ./...

# Run tests with coverage
go test -v -cover ./...

# Format code
go fmt ./...

# Tidy dependencies
go mod tidy

# Check for updates
go list -u -m all
```

## Application Architecture

### Project Structure
The Go application follows standard Go project layout:
1. **cmd/minecraft-instance-manager/** - Application entry point and CLI setup
2. **internal/instance/** - Core business logic for instance management
3. **internal/tui/** - Bubble Tea terminal user interface components
4. **pkg/** - Reusable packages (configuration, utilities)

### Key Components
- `instance.Manager` - Core instance management with full CRUD operations
- `tui.Model` - Bubble Tea model handling UI state and user interactions
- `cobra.Command` - CLI command definitions with proper argument validation
- GitHub Actions - Automated testing and cross-platform binary builds

### TUI Architecture
- **State Management** - Clean state transitions (list → detail → create → delete)
- **Keyboard Handling** - Intuitive shortcuts with help system
- **Real-time Updates** - Immediate UI refresh after operations
- **Error Handling** - User-friendly error messages and recovery

### Safety Mechanisms
- Always backs up current .minecraft before switching
- Validates instance exists before switching  
- Uses symlinks (non-destructive, easily reversible)
- Provides restore functionality
- Confirmation dialogs for destructive operations
- Proper error handling and user feedback

## Development Guidelines

### Code Style
- Follow Go conventions (gofmt, go vet, golint)
- Use meaningful package and function names
- Include comprehensive error handling with wrapped errors
- Add helpful user messages and feedback in both CLI and TUI modes
- Use structured logging when necessary

### Testing Approach
- Test with different Minecraft directory states (exists/doesn't exist)
- Test error conditions (invalid instance names, missing directories, permission issues)
- Verify symlink creation and backup functionality  
- Test mod/config/save counting accuracy
- Test TUI state transitions and keyboard interactions
- Test CLI command parsing and validation

### File Structure Expectations
- Follow Go project layout standards
- Keep business logic in internal/instance package
- Separate TUI concerns in internal/tui package
- Use dependency injection for testability
- Keep CLI commands thin, delegating to business logic

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