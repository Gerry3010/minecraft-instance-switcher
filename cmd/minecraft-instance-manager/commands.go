package main

import (
	"fmt"
	"os"

	"github.com/Gerry3010/minecraft-instance-switcher/internal/instance"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(switchCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(restoreCmd)
	rootCmd.AddCommand(deleteCmd)
}

var createCmd = &cobra.Command{
	Use:   "create <instance-name>",
	Short: "Create a new Minecraft instance",
	Long: `Create a new Minecraft instance with the given name.
This will copy your current .minecraft directory structure to create a new instance.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		manager, err := instance.NewManager()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error initializing manager: %v\n", err)
			os.Exit(1)
		}

		instanceName := args[0]
		if err := manager.CreateInstance(instanceName); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating instance: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Created instance: %s\n", instanceName)
		fmt.Printf("Add mods to: %s/.minecraft-instances/%s/mods/\n", manager.HomeDir, instanceName)
	},
}

var switchCmd = &cobra.Command{
	Use:   "switch <instance-name>",
	Short: "Switch to a Minecraft instance",
	Long: `Switch to the specified Minecraft instance.
This will backup your current .minecraft directory and create a symlink to the instance.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		manager, err := instance.NewManager()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error initializing manager: %v\n", err)
			os.Exit(1)
		}

		instanceName := args[0]
		if err := manager.SwitchInstance(instanceName); err != nil {
			fmt.Fprintf(os.Stderr, "Error switching instance: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Switched to instance: %s\n", instanceName)
		fmt.Println("Launch Minecraft normally - it will use this instance")
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all Minecraft instances",
	Long:  `List all available Minecraft instances with their mod counts and status.`,
	Run: func(cmd *cobra.Command, args []string) {
		manager, err := instance.NewManager()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error initializing manager: %v\n", err)
			os.Exit(1)
		}

		instances, err := manager.ListInstances()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error listing instances: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Available instances:")
		if len(instances) == 0 {
			fmt.Println("  No instances found")
		} else {
			for _, inst := range instances {
				status := "Inactive"
				if inst.IsActive {
					status = "ACTIVE"
				}
				fmt.Printf("  - %-20s (%d mods, %d configs, %d saves) [%s]\n",
					inst.Name, inst.ModCount, inst.ConfigCount, inst.SaveCount, status)
			}
		}

		fmt.Printf("\nCurrent instance: %s\n", manager.GetActiveInstance())
	},
}

var restoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "Restore default .minecraft directory",
	Long: `Restore the original .minecraft directory by removing the current symlink
and restoring from the backup.`,
	Run: func(cmd *cobra.Command, args []string) {
		manager, err := instance.NewManager()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error initializing manager: %v\n", err)
			os.Exit(1)
		}

		if err := manager.RestoreDefault(); err != nil {
			fmt.Fprintf(os.Stderr, "Error restoring default: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Restored default .minecraft directory")
	},
}

var deleteCmd = &cobra.Command{
	Use:   "delete <instance-name>",
	Short: "Delete a Minecraft instance",
	Long: `Delete the specified Minecraft instance permanently.
Note: You cannot delete the currently active instance.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		manager, err := instance.NewManager()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error initializing manager: %v\n", err)
			os.Exit(1)
		}

		instanceName := args[0]
		
		// Confirm deletion
		fmt.Printf("Are you sure you want to delete instance '%s'? This cannot be undone. (y/N): ", instanceName)
		var response string
		fmt.Scanln(&response)
		
		if response != "y" && response != "Y" {
			fmt.Println("Deletion cancelled")
			return
		}

		if err := manager.DeleteInstance(instanceName); err != nil {
			fmt.Fprintf(os.Stderr, "Error deleting instance: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Deleted instance: %s\n", instanceName)
	},
}