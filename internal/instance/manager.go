package instance

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const (
	InstancesDir  = ".minecraft-instances"
	MinecraftDir  = ".minecraft"
	BackupSuffix  = ".backup"
)

type Manager struct {
	HomeDir       string
	InstancesPath string
	MinecraftPath string
	BackupPath    string
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
	ModsDir     []string
	ConfigsDir  []string
	SavesDir    []string
	OtherFiles  []string
}

func NewManager() (*Manager, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user home directory: %w", err)
	}

	return &Manager{
		HomeDir:       homeDir,
		InstancesPath: filepath.Join(homeDir, InstancesDir),
		MinecraftPath: filepath.Join(homeDir, MinecraftDir),
		BackupPath:    filepath.Join(homeDir, MinecraftDir+BackupSuffix),
	}, nil
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

	// Copy base minecraft structure if it exists
	if _, err := os.Stat(m.MinecraftPath); err == nil {
		if err := copyDir(m.MinecraftPath, instancePath); err != nil {
			return fmt.Errorf("failed to copy minecraft directory: %w", err)
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
			return os.MkdirAll(dstPath, d.Type().Perm())
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