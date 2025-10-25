package instance

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
)

const (
	InstancesDir  = ".minecraft-instances"
	MinecraftDir  = ".minecraft"
	BackupSuffix  = ".backup"
	AppFolderName = "minecraft-instance" // folder inside OS config/app support dir
)

type Config struct {
	InstancesPath string `json:"instances_path"`
	MinecraftPath string `json:"minecraft_path"`
	BackupPath    string `json:"backup_path"`
}

type Manager struct {
	HomeDir       string
	AppDir        string
	ConfigFile    string
	InstancesPath string
	MinecraftPath string
	BackupPath    string
	cfg           Config
}

type Instance struct {
	Name        string
	Path        string
	ModCount    int
	ConfigCount int
	SaveCount   int
	IsActive    bool
}

type InstanceInfo struct {
	ModsDir    []string
	ConfigsDir []string
	SavesDir   []string
	OtherFiles []string
}

func NewManager() (*Manager, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user home directory: %w", err)
	}

	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user config directory: %w", err)
	}

	appDir := filepath.Join(userConfigDir, AppFolderName)
	if err := os.MkdirAll(appDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create app config dir: %w", err)
	}

	configFile := filepath.Join(appDir, "config.json")

	// Get platform-specific default Minecraft path
	defaultMinecraftPath, err := getDefaultMinecraftPath()
	if err != nil {
		return nil, fmt.Errorf("failed to determine default Minecraft path: %w", err)
	}

	// Other defaults
	defaultInstancesPath := filepath.Join(appDir, "instances")
	defaultBackupPath := filepath.Join(appDir, "backup")

	m := &Manager{
		HomeDir:       homeDir,
		AppDir:        appDir,
		ConfigFile:    configFile,
		InstancesPath: defaultInstancesPath,
		MinecraftPath: defaultMinecraftPath,
		BackupPath:    defaultBackupPath,
		cfg: Config{
			InstancesPath: defaultInstancesPath,
			MinecraftPath: defaultMinecraftPath,
			BackupPath:    defaultBackupPath,
		},
	}

	// Load config if exists, otherwise create it with defaults
	if err := m.loadConfig(); err != nil {
		// if load failed because file doesn't exist, save defaults
		if os.IsNotExist(err) {
			if err := os.MkdirAll(m.InstancesPath, 0755); err != nil {
				return nil, fmt.Errorf("failed to create default instances dir: %w", err)
			}
			if err := m.saveConfig(); err != nil {
				return nil, fmt.Errorf("failed to write default config: %w", err)
			}
		} else {
			return nil, err
		}
	} else {
		// apply loaded config
		m.InstancesPath = m.cfg.InstancesPath
		m.MinecraftPath = m.cfg.MinecraftPath
		m.BackupPath = m.cfg.BackupPath
		// ensure instances dir exists
		if err := os.MkdirAll(m.InstancesPath, 0755); err != nil {
			return nil, fmt.Errorf("failed to create instances dir from config: %w", err)
		}
	}

	return m, nil
}

func (m *Manager) loadConfig() error {
	data, err := os.ReadFile(m.ConfigFile)
	if err != nil {
		return err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}

	// expand ~ in paths
	cfg.InstancesPath = expandPath(cfg.InstancesPath)
	cfg.MinecraftPath = expandPath(cfg.MinecraftPath)
	cfg.BackupPath = expandPath(cfg.BackupPath)

	// if any are empty, set defaults relative to app dir / platform-specific paths
	if cfg.InstancesPath == "" {
		cfg.InstancesPath = filepath.Join(m.AppDir, "instances")
	}
	if cfg.MinecraftPath == "" {
		if defaultPath, err := getDefaultMinecraftPath(); err == nil {
			cfg.MinecraftPath = defaultPath
		} else {
			// Fallback to the old behavior if detection fails
			cfg.MinecraftPath = filepath.Join(m.HomeDir, MinecraftDir)
		}
	}
	if cfg.BackupPath == "" {
		cfg.BackupPath = filepath.Join(m.AppDir, "backup")
	}

	m.cfg = cfg
	return nil
}

func (m *Manager) saveConfig() error {
	m.cfg.InstancesPath = m.InstancesPath
	m.cfg.MinecraftPath = m.MinecraftPath
	m.cfg.BackupPath = m.BackupPath

	data, err := json.MarshalIndent(m.cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to encode config: %w", err)
	}
	if err := os.WriteFile(m.ConfigFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}
	return nil
}

func expandPath(p string) string {
	if p == "" {
		return p
	}
	if strings.HasPrefix(p, "~") {
		if home, err := os.UserHomeDir(); err == nil {
			return filepath.Join(home, strings.TrimPrefix(p, "~"))
		}
	}
	return p
}

// getDefaultMinecraftPath returns the platform-specific default Minecraft directory path
// This function automatically detects the operating system and returns the appropriate
// default path where Minecraft is typically installed on each platform.
func getDefaultMinecraftPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}

	switch runtime.GOOS {
	case "windows":
		// Windows: %APPDATA%\.minecraft (e.g., C:\Users\Username\AppData\Roaming\.minecraft)
		appDataDir := os.Getenv("APPDATA")
		if appDataDir != "" {
			return filepath.Join(appDataDir, ".minecraft"), nil
		}
		// Fallback to user profile if APPDATA environment variable is not set
		userProfile := os.Getenv("USERPROFILE")
		if userProfile != "" {
			return filepath.Join(userProfile, "AppData", "Roaming", ".minecraft"), nil
		}
		// Ultimate fallback - construct path from home directory
		return filepath.Join(homeDir, "AppData", "Roaming", ".minecraft"), nil

	case "darwin":
		// macOS: ~/Library/Application Support/minecraft
		// Note: Minecraft on macOS uses "minecraft" (lowercase) not ".minecraft"
		return filepath.Join(homeDir, "Library", "Application Support", "minecraft"), nil

	default:
		// Linux and other Unix-like systems: ~/.minecraft
		// This includes most Linux distributions, BSD variants, etc.
		return filepath.Join(homeDir, ".minecraft"), nil
	}
}

// UpdateConfig updates one of the supported config keys and persists the file.
// Supported keys: "minecraft-path", "instances-path", "backup-path"
func (m *Manager) UpdateConfig(key, value string) error {
	value = expandPath(value)
	switch key {
	case "minecraft-path", "minecraft-dir", "minecraft":
		m.MinecraftPath = value
	case "instances-path", "instances-dir", "instances":
		m.InstancesPath = value
		// ensure instances dir exists
		if err := os.MkdirAll(m.InstancesPath, 0755); err != nil {
			return fmt.Errorf("failed to create instances dir: %w", err)
		}
	case "backup-path", "backup-dir", "backup":
		m.BackupPath = value
	default:
		return fmt.Errorf("unknown config key: %s", key)
	}
	if err := m.saveConfig(); err != nil {
		return err
	}
	return nil
}

// GetConfig returns current configuration as a map
func (m *Manager) GetConfig() map[string]string {
	return map[string]string{
		"minecraft-path": m.MinecraftPath,
		"instances-path": m.InstancesPath,
		"backup-path":    m.BackupPath,
		"app-dir":        m.AppDir,
		"config-file":    m.ConfigFile,
	}
}

func (m *Manager) CreateInstance(name string) error {
	if name == "" {
		return fmt.Errorf("instance name cannot be empty")
	}

	instancePath := filepath.Join(m.InstancesPath, name)

	// Check if instance already exists
	if _, err := os.Stat(instancePath); err == nil {
		return fmt.Errorf("instance '%s' already exists", name)
	}

	// Create instances directory if it doesn't exist
	if err := os.MkdirAll(m.InstancesPath, 0755); err != nil {
		return fmt.Errorf("failed to create instances directory: %w", err)
	}

	// Create instance directory
	if err := os.MkdirAll(instancePath, 0755); err != nil {
		return fmt.Errorf("failed to create instance directory: %w", err)
	}

	// Copy base minecraft structure if it exists and is not a symlink
	if info, err := os.Lstat(m.MinecraftPath); err == nil {
		// If it's a symlink, resolve it and copy from the actual directory
		if info.Mode()&os.ModeSymlink != 0 {
			if target, err := os.Readlink(m.MinecraftPath); err == nil {
				if err := copyDir(target, instancePath); err != nil {
					return fmt.Errorf("failed to copy minecraft directory: %w", err)
				}
			}
		} else {
			// It's a regular directory
			if err := copyDir(m.MinecraftPath, instancePath); err != nil {
				return fmt.Errorf("failed to copy minecraft directory: %w", err)
			}
		}
	}

	// Create essential directories
	essentialDirs := []string{"mods", "config", "saves", "resourcepacks", "shaderpacks"}
	for _, dir := range essentialDirs {
		dirPath := filepath.Join(instancePath, dir)
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return nil
}

func (m *Manager) SwitchInstance(name string) error {
	if name == "" {
		return fmt.Errorf("instance name cannot be empty")
	}

	instancePath := filepath.Join(m.InstancesPath, name)

	// Check if instance exists
	if _, err := os.Stat(instancePath); os.IsNotExist(err) {
		return fmt.Errorf("instance '%s' does not exist", name)
	}

	// Backup current minecraft directory if it exists and is not a symlink
	if info, err := os.Lstat(m.MinecraftPath); err == nil {
		if info.Mode()&os.ModeSymlink == 0 {
			// It's a regular directory, back it up
			if err := os.RemoveAll(m.BackupPath); err != nil {
				return fmt.Errorf("failed to remove old backup: %w", err)
			}
			if err := os.Rename(m.MinecraftPath, m.BackupPath); err != nil {
				return fmt.Errorf("failed to backup minecraft directory: %w", err)
			}
		} else {
			// It's already a symlink, just remove it
			if err := os.Remove(m.MinecraftPath); err != nil {
				return fmt.Errorf("failed to remove existing symlink: %w", err)
			}
		}
	}

	// Create symlink to instance
	if err := os.Symlink(instancePath, m.MinecraftPath); err != nil {
		return fmt.Errorf("failed to create symlink: %w", err)
	}

	return nil
}

func (m *Manager) RestoreDefault() error {
	// Check if minecraft path is a symlink
	if info, err := os.Lstat(m.MinecraftPath); err == nil && info.Mode()&os.ModeSymlink != 0 {
		// Remove the symlink
		if err := os.Remove(m.MinecraftPath); err != nil {
			return fmt.Errorf("failed to remove symlink: %w", err)
		}
	}

	// Restore backup if it exists
	if _, err := os.Stat(m.BackupPath); err == nil {
		if err := os.Rename(m.BackupPath, m.MinecraftPath); err != nil {
			return fmt.Errorf("failed to restore backup: %w", err)
		}
	}

	return nil
}

func (m *Manager) ListInstances() ([]Instance, error) {
	var instances []Instance

	// Check if instances directory exists
	if _, err := os.Stat(m.InstancesPath); os.IsNotExist(err) {
		return instances, nil
	}

	// Read instances directory
	entries, err := os.ReadDir(m.InstancesPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read instances directory: %w", err)
	}

	// Get current active instance
	activeInstance := m.GetActiveInstance()

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		name := entry.Name()
		instancePath := filepath.Join(m.InstancesPath, name)

		instance := Instance{
			Name:     name,
			Path:     instancePath,
			IsActive: name == activeInstance,
		}

		// Count mods
		modsPath := filepath.Join(instancePath, "mods")
		instance.ModCount = countJarFiles(modsPath)

		// Count configs
		configPath := filepath.Join(instancePath, "config")
		instance.ConfigCount = countFiles(configPath)

		// Count saves
		savesPath := filepath.Join(instancePath, "saves")
		instance.SaveCount = countDirectories(savesPath)

		instances = append(instances, instance)
	}

	// Sort instances alphabetically
	sort.Slice(instances, func(i, j int) bool {
		return instances[i].Name < instances[j].Name
	})

	return instances, nil
}

func (m *Manager) GetActiveInstance() string {
	if info, err := os.Lstat(m.MinecraftPath); err == nil && info.Mode()&os.ModeSymlink != 0 {
		if target, err := os.Readlink(m.MinecraftPath); err == nil {
			return filepath.Base(target)
		}
	}
	return "default"
}

func (m *Manager) GetInstanceInfo(name string) (*InstanceInfo, error) {
	instancePath := filepath.Join(m.InstancesPath, name)

	if _, err := os.Stat(instancePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("instance '%s' does not exist", name)
	}

	info := &InstanceInfo{}

	// Get mods
	modsPath := filepath.Join(instancePath, "mods")
	info.ModsDir = getJarFiles(modsPath)

	// Get configs
	configPath := filepath.Join(instancePath, "config")
	info.ConfigsDir = getConfigFiles(configPath)

	// Get saves
	savesPath := filepath.Join(instancePath, "saves")
	info.SavesDir = getDirectoryNames(savesPath)

	return info, nil
}

func (m *Manager) DeleteInstance(name string) error {
	if name == "" {
		return fmt.Errorf("instance name cannot be empty")
	}

	instancePath := filepath.Join(m.InstancesPath, name)

	// Check if instance exists
	if _, err := os.Stat(instancePath); os.IsNotExist(err) {
		return fmt.Errorf("instance '%s' does not exist", name)
	}

	// Check if it's the active instance
	if m.GetActiveInstance() == name {
		return fmt.Errorf("cannot delete active instance '%s'. Switch to another instance first", name)
	}

	// Remove the instance directory
	return os.RemoveAll(instancePath)
}

// Helper functions

func copyDir(src, dst string) error {
	return filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		dstPath := filepath.Join(dst, relPath)

		if d.IsDir() {
			return os.MkdirAll(dstPath, 0755)
		}

		// Skip copying certain files/directories
		if strings.Contains(relPath, ".git") || strings.Contains(relPath, ".DS_Store") {
			return nil
		}

		return copyFile(path, dstPath)
	})
}

func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	return os.WriteFile(dst, data, 0644)
}

func countJarFiles(dir string) int {
	count := 0
	if entries, err := os.ReadDir(dir); err == nil {
		for _, entry := range entries {
			if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".jar") {
				count++
			}
		}
	}
	return count
}

func countFiles(dir string) int {
	count := 0
	if entries, err := os.ReadDir(dir); err == nil {
		for _, entry := range entries {
			if !entry.IsDir() {
				count++
			}
		}
	}
	return count
}

func countDirectories(dir string) int {
	count := 0
	if entries, err := os.ReadDir(dir); err == nil {
		for _, entry := range entries {
			if entry.IsDir() {
				count++
			}
		}
	}
	return count
}

func getJarFiles(dir string) []string {
	var files []string
	if entries, err := os.ReadDir(dir); err == nil {
		for _, entry := range entries {
			if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".jar") {
				files = append(files, entry.Name())
			}
		}
	}
	sort.Strings(files)
	return files
}

func getConfigFiles(dir string) []string {
	var files []string
	if entries, err := os.ReadDir(dir); err == nil {
		for _, entry := range entries {
			if !entry.IsDir() {
				files = append(files, entry.Name())
			}
		}
	}
	sort.Strings(files)
	return files
}

func getDirectoryNames(dir string) []string {
	var dirs []string
	if entries, err := os.ReadDir(dir); err == nil {
		for _, entry := range entries {
			if entry.IsDir() {
				dirs = append(dirs, entry.Name())
			}
		}
	}
	sort.Strings(dirs)
	return dirs
}
