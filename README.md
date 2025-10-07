# Minecraft Instance Manager

A lightweight and efficient Minecraft instance manager that uses symlinks to instantly switch between different Minecraft setups without copying files.

## ✨ Features

- **⚡ Instant switching** - Uses symlinks, no file copying required
- **💾 Space efficient** - No duplicate files, shared assets
- **🔒 Safe backups** - Automatic backup before switching
- **📊 Mod counting** - Shows mod count for each instance
- **🧹 Clean design** - Simple bash script, no dependencies
- **🔄 Easy restore** - One command to restore original setup

## 🚀 Quick Start

### Installation

```bash
# Clone or download the script
chmod +x minecraft-instances
```

### Basic Usage

```bash
# Create instances
./minecraft-instances create vanilla           # Clean Minecraft
./minecraft-instances create modpack-1.20.1    # Your modpack
./minecraft-instances create testing           # Testing environment

# Switch between instances
./minecraft-instances switch modpack-1.20.1   # Switch to modpack
./minecraft-instances switch vanilla           # Switch to vanilla

# List all instances
./minecraft-instances list

# Restore original .minecraft
./minecraft-instances restore
```

## 📋 Commands

| Command | Description | Example |
|---------|-------------|---------|
| `create <name>` | Create a new instance | `./minecraft-instances create forge-1.20.1` |
| `switch <name>` | Switch to an instance | `./minecraft-instances switch vanilla` |
| `list` | List all instances with mod counts | `./minecraft-instances list` |
| `restore` | Restore original .minecraft directory | `./minecraft-instances restore` |

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
./minecraft-instances switch my-instance
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
- Use `./minecraft-instances restore` to return to original setup

### Lost original .minecraft
- Your original is backed up at `~/.minecraft.backup`
- Run `./minecraft-instances restore` to recover it

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