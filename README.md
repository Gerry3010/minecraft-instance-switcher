# Minecraft Instance Manager

A modern, lightweight Minecraft instance manager with a beautiful terminal interface. Uses symlinks to instantly switch between different Minecraft setups without copying files.

## ✨ Features

- **⚡ Instant switching** - Uses symlinks, no file copying required
- **💾 Space efficient** - No duplicate files, shared assets
- **🔒 Safe backups** - Automatic backup before switching
- **📊 Rich information** - Shows mods, configs, saves count for each instance
- **🎨 Beautiful TUI** - Interactive terminal interface with Bubble Tea
- **🧹 Modern design** - Written in Go with Cobra CLI framework
- **🔄 Easy restore** - One command to restore original setup
- **🌐 Cross-platform** - Works on Linux, macOS, and Windows

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
| `r` | Refresh instance list |
| `R` | Restore default .minecraft |
| `?` | Toggle help |
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