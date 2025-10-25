# Minecraft Instance Manager

A modern, lightweight Minecraft instance manager with a beautiful terminal interface. Uses symlinks to instantly switch between different Minecraft setups without copying files.

> 🤖 **AI Collaboration Notice**: This project was developed in collaboration with the Warp AI assistant (powered by Claude 3.5 Sonnet). The AI helped with code implementation, documentation, and project structure. While the core ideas and direction came from human creativity, the AI's assistance made this project more robust and feature-complete. We believe in transparency about AI usage while celebrating the potential of human-AI collaboration in software development.

## ✨ Features

- **⚡ Instant switching** - Uses symlinks, no file copying required
- **💾 Space efficient** - No duplicate files, shared assets
- **🔒 Safe backups** - Automatic backup before switching
- **📊 Rich information** - Shows mods, configs, saves count for each instance
- **🎨 Beautiful TUI** - Interactive terminal interface with Bubble Tea
- **🧹 Modern design** - Written in Go with Cobra CLI framework
- **🔄 Easy restore** - One command to restore original setup
- **🌐 Cross-platform** - Works on Linux, macOS, and Windows (with considerations)

## 🚀 Quick Start

### Installation

#### Download Pre-built Binary
Download the latest release for your platform from the [releases page](https://github.com/Gerry3010/minecraft-instance-switcher/releases).

#### Install with Go
```bash
go install github.com/Gerry3010/minecraft-instance-switcher/cmd/minecraft-instance-manager@latest
```

#### Build from Source
```bash
git clone https://github.com/Gerry3010/minecraft-instance-switcher.git
cd minecraft-instance-switcher
go build -o minecraft-instance-manager ./cmd/minecraft-instance-manager
```

### Usage

#### Interactive TUI Mode (Default)
```bash
# Launch the beautiful terminal interface
# Automatically detects your OS and sets appropriate Minecraft paths
minecraft-instance-manager
```

#### Command Line Mode
```bash
# Create instances
minecraft-instance-manager create vanilla           # Clean Minecraft
minecraft-instance-manager create modpack-1.20.1    # Your modpack
minecraft-instance-manager create testing           # Testing environment

# Switch between instances
minecraft-instance-manager switch modpack-1.20.1   # Switch to modpack
minecraft-instance-manager switch vanilla           # Switch to vanilla

# List all instances with details
minecraft-instance-manager list

# Delete an instance
minecraft-instance-manager delete old-instance

# Restore original .minecraft
minecraft-instance-manager restore
```

## 🎨 TUI Features

The interactive terminal interface provides:

- **🎨 Beautiful Interface** - Clean, modern terminal UI with colors and styling
- **⌨️ Keyboard Navigation** - Full keyboard control with intuitive shortcuts
- **📋 Instance List** - See all instances with mod/config/save counts at a glance
- **🔍 Instance Details** - View detailed information about any instance
- **⚡ Quick Switching** - Switch instances with just Enter key
- **➕ Create/Delete** - Create new instances or delete existing ones
- **🎆 Real-time Updates** - Interface updates instantly when changes are made

### TUI Keyboard Shortcuts

| Key | Action |
|-----|--------|
| `↑/↓` or `j/k` | Navigate up/down |
| `Enter` | Switch to instance or view details |
| `c` | Create new instance |
| `d` | Delete selected instance |
| `s` | Show detailed file panels (in detail view) |
| `Tab/Shift+Tab` | Switch between panels (in panel view) |
| `F5` | Refresh instance list |
| `r` | Restore default .minecraft |
| `?` | Toggle help |
| `ESC` | Go back / Cancel |
| `q` or `Ctrl+C` | Quit |

## 📋 CLI Commands

| Command | Description | Example |
|---------|-------------|---------|
| `create <name>` | Create a new instance | `minecraft-instance-manager create forge-1.20.1` |
| `switch <name>` | Switch to an instance | `minecraft-instance-manager switch vanilla` |
| `list` | List all instances with details | `minecraft-instance-manager list` |
| `delete <name>` | Delete an instance | `minecraft-instance-manager delete old-instance` |
| `restore` | Restore original .minecraft directory | `minecraft-instance-manager restore` |

## 📁 How It Works

### Directory Structure

```
~/.minecraft-instances/
├── vanilla/
│   ├── mods/           (empty)
│   ├── config/
│   ├── saves/
│   └── ...
├── modpack-1.20.1/
│   ├── mods/           (108 mods)
│   ├── config/
│   ├── saves/
│   └── ...
└── testing/
    ├── mods/           (1 mod)
    ├── config/
    ├── saves/
    └── ...
```

### Symlink Magic

When you switch instances:
1. Current `~/.minecraft` is backed up to `~/.minecraft.backup`
2. A symlink `~/.minecraft -> ~/.minecraft-instances/chosen-instance` is created
3. Minecraft launcher uses the instance transparently

## 🎯 Use Cases

### Mod Development
```bash
./minecraft-instances create dev-environment
# Add your mod to ~/.minecraft-instances/dev-environment/mods/
./minecraft-instances switch dev-environment
```

### Different Minecraft Versions
```bash
./minecraft-instances create mc-1.19.4
./minecraft-instances create mc-1.20.1
./minecraft-instances create mc-1.21
```

### Modpack Testing
```bash
./minecraft-instances create modpack-backup
./minecraft-instances create modpack-experimental
# Test changes in experimental, keep backup safe
```

## 🛡️ Safety Features

- **Automatic backups** - Your original .minecraft is always backed up
- **Safe switching** - Validates instance exists before switching
- **Easy restore** - One command restores original setup
- **Non-destructive** - Never deletes your original data

## 🔧 Advanced Usage

### Adding Mods to an Instance
```bash
# Add mods directly to the instance
cp my-mod.jar ~/.minecraft-instances/my-instance/mods/

# Or switch to instance and use normal mod installation
minecraft-instance-manager switch my-instance
# Now use your launcher's mod management or copy mods to ~/.minecraft/mods/
```

### Sharing Instances
```bash
# Backup an instance
tar -czf my-modpack.tar.gz ~/.minecraft-instances/my-modpack/

# Restore on another machine
tar -xzf my-modpack.tar.gz -C ~/
```

## 📊 Instance Information

The `list` command shows detailed information:

```
Available instances:
  - vanilla               (0 mods)
  - sebi-1.20.1          (107 mods)
  - vanillaplus-test      (1 mods)

Current instance: sebi-1.20.1
```

## 🐛 Troubleshooting

### Instance doesn't appear in list
- Check that `~/.minecraft-instances/instance-name` exists
- Ensure the directory has proper permissions

### Minecraft won't launch
- Verify the instance has all required files (copied during creation)
- Check Minecraft version compatibility
- Use `minecraft-instance-manager restore` to return to original setup

### Lost original .minecraft
- Your original is backed up at `~/.minecraft.backup`
- Run `minecraft-instance-manager restore` to recover it

## 🖥️ Platform Compatibility

The Minecraft Instance Manager is designed to work across different operating systems with some platform-specific considerations:

### ✅ Fully Supported Platforms

| Platform | Status | Notes |
|----------|--------|--------|
| **Linux** | ✅ Full Support | Native symlink support, standard `.minecraft` path |
| **macOS** | ✅ Full Support | Native symlink support, standard `.minecraft` path |
| **WSL/WSL2** | ✅ Full Support | Linux compatibility within Windows |

### ⚠️ Windows Considerations

| Feature | Status | Requirements |
|---------|--------|-------------|
| **Basic Functionality** | ✅ Supported | Windows 10/11 |
| **Symlink Creation** | ⚠️ Requires Privileges | Administrator rights OR Developer Mode |
| **Minecraft Path** | ⚠️ Manual Config | May need to set custom path |

#### Windows Setup Instructions

**Option 1: Enable Developer Mode (Recommended)**
1. Open Settings → Update & Security → For Developers
2. Enable "Developer Mode"
3. Restart your computer
4. Run the application normally

**Option 2: Run as Administrator**
1. Right-click Command Prompt/PowerShell
2. Select "Run as Administrator"
3. Run minecraft-instance-manager commands

**Option 3: Use WSL2 (Best Experience)**
1. Install WSL2 with Ubuntu
2. Install and run minecraft-instance-manager in WSL2
3. Access Windows Minecraft installation via `/mnt/c/Users/.../AppData/Roaming/.minecraft`

### 🗂️ Platform-Specific Paths

#### Automatic Platform Detection

The application **automatically detects your operating system** and sets appropriate default paths:

| Platform | Default Path | Auto-Detected | Configurable |
|----------|-------------|---------------|--------------|
| **Linux** | `~/.minecraft` | ✅ Yes | ✅ Yes |
| **macOS** | `~/Library/Application Support/minecraft` | ✅ Yes | ✅ Yes |
| **Windows** | `%APPDATA%\.minecraft` | ✅ Yes | ✅ Yes |

> 💡 **No manual configuration needed!** The app automatically uses the correct path for your platform.

#### Configuration Locations

| Platform | Config Directory |
|----------|-----------------|
| **Linux** | `~/.config/minecraft-instance/` |
| **macOS** | `~/Library/Application Support/minecraft-instance/` |
| **Windows** | `%APPDATA%\minecraft-instance\` |

### 🔧 Custom Path Configuration

If your Minecraft installation is in a non-standard location:

```bash
# Set custom Minecraft path
minecraft-instance-manager config minecraft-path "C:\Games\Minecraft\.minecraft"

# Set custom instances directory  
minecraft-instance-manager config instances-path "D:\MinecraftInstances"

# Verify configuration
minecraft-instance-manager config show
```

### 🚀 Build Information

Pre-built binaries are available for:
- **Linux**: AMD64, ARM64
- **macOS**: AMD64 (Intel), ARM64 (Apple Silicon)  
- **Windows**: AMD64

Download from the [releases page](https://github.com/Gerry3010/minecraft-instance-switcher/releases).

### 🔍 Platform-Specific Troubleshooting

#### Windows Issues

**"Access Denied" or Symlink Errors:**
- Enable Developer Mode or run as Administrator
- Check if Minecraft path is correct: `%APPDATA%\.minecraft`
- Consider using WSL2 for better compatibility

**Path Issues:**
```bash
# Windows example paths
minecraft-instance-manager config minecraft-path "C:\Users\Username\AppData\Roaming\.minecraft"
minecraft-instance-manager config instances-path "C:\MinecraftInstances"
```

#### macOS Issues

**Permissions:**
```bash
# Give terminal full disk access in Security & Privacy settings
# Or use chmod to fix permissions
chmod -R 755 ~/.minecraft-instances/
```

#### Linux Issues

**Snap/Flatpak Minecraft:**
- May require custom path configuration
- Check if Minecraft runs in sandboxed environment

### 🧪 Testing Your Platform

Verify compatibility on your system:

```bash
# Test basic functionality
minecraft-instance-manager config show

# Test instance creation (safe)
minecraft-instance-manager create test-compatibility

# Test symlink creation (creates backup first)
minecraft-instance-manager switch test-compatibility

# Restore original setup
minecraft-instance-manager restore

# Clean up test
minecraft-instance-manager delete test-compatibility
```

## 🤝 Contributing

This script is simple and focused. Contributions welcome for:
- Bug fixes
- Small enhancements
- Documentation improvements
- Platform compatibility

## 📜 License

MIT License - Feel free to use, modify, and share!

## 🙏 Credits

Originally created for VanillaPlusAdditions mod development workflow.

---

**Happy mining!** ⛏️✨