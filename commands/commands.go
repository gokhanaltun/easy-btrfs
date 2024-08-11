package commands

import (
	"easy-btrfs/actions"

	"github.com/urfave/cli/v2"
)

// CreateConfig creates a new configuration for a subvolume.
// Usage: ebtrfs create-config [configName] [subvolumePath]
// Example: ebtrfs create-config root /path/to/subvolume
func CreateConfig() *cli.Command {
	return &cli.Command{
		Name:    "create-config",
		Aliases: []string{"cconf", "cc"},
		Usage:   " - creates a new configuration for a subvolume. \n - Usage: ebtrfs create-config [configName] [subvolumePath] \n - Example: ebtrfs create-config root /path/to/subvolume \n",
		Action:  actions.CreateConfig,
	}
}

// ListConfigs lists all configurations for subvolumes.
// Usage: ebtrfs list-configs
// Example: ebtrfs list-configs
func ListConfigs() *cli.Command {
	return &cli.Command{
		Name:    "list-configs",
		Aliases: []string{"ls-c", "lsc", "lc"},
		Usage:   " - lists all configurations for subvolumes. \n - Usage: ebtrfs list-configs \n - Example: ebtrfs list-configs \n",
		Action:  actions.ListConfigs,
	}
}

// DeleteConfig deletes a specific configuration.
// Usage: ebtrfs delete-config [configName]
// Example: ebtrfs delete-config root
func DeleteConfig() *cli.Command {
	return &cli.Command{
		Name:    "delete-config",
		Aliases: []string{"del-conf", "d-conf", "dc"},
		Usage:   " - deletes a specific configuration. \n - Usage: ebtrfs delete-config [configName] \n - Example: ebtrfs delete-config root \n",
		Action:  actions.DeleteConfig,
	}
}

// Snapshot takes a snapshot of a subvolume.
// Usage: ebtrfs snapshot [configName]
// Example: ebtrfs snapshot root
func Snapshot() *cli.Command {
	return &cli.Command{
		Name:    "snapshot",
		Aliases: []string{"snap", "s"},
		Usage:   " - takes a snapshot of a subvolume. \n - Usage: ebtrfs snapshot [configName] \n - Example: ebtrfs snapshot root \n",
		Action:  actions.Snapshot,
	}
}

// ListSnapshots lists all snapshots for a specific configuration.
// Usage: ebtrfs list-snapshots [configName]
// Example: ebtrfs list-snapshots root
func ListSnapshots() *cli.Command {
	return &cli.Command{
		Name:    "list-snapshots",
		Aliases: []string{"list-snaps", "ls-s", "lss"},
		Usage:   " - lists all snapshots for a specific configuration. \n - Usage: ebtrfs list-snapshots [configName] \n - Example: ebtrfs list-snapshots root \n",
		Action:  actions.ListSnapshots,
	}
}

// DeleteSnapshots deletes specific snapshots.
// Usage: ebtrfs delete-snapshots [configName] [snapshotIds...]
// Example: ebtrfs delete-snapshots root 1 2 3
func DeleteSnapshots() *cli.Command {
	return &cli.Command{
		Name:    "delete-snapshots",
		Aliases: []string{"del-snaps", "d-snaps", "ds"},
		Usage:   " - deletes specific snapshots. \n - Usage: ebtrfs delete-snapshots [configName] [snapshotIds...] \n - Example: ebtrfs delete-snapshots root 1 2 3 \n",
		Action:  actions.DeleteSnapshots,
	}
}

// RollBack rolls back to a specific snapshot.
// Usage: ebtrfs rollback [configName] [snapshotId]
// Example: ebtrfs rollback root 1
func RollBack() *cli.Command {
	return &cli.Command{
		Name:    "rollback",
		Aliases: []string{"rb", "r"},
		Usage:   " - rolls back to a specific snapshot. \n - Usage: ebtrfs rollback [configName] [snapshotId] \n - Example: ebtrfs rollback root 1 \n",
		Action:  actions.RollBack,
	}
}
