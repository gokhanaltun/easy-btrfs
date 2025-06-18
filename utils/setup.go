package utils

import (
	"easy-btrfs/database"
	"easy-btrfs/models"
	"easy-btrfs/store"
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

// createSubvol tries to create a Btrfs subvolume at the given path if it does not already exist.
// It lists existing subvolumes under the defined MountPoint, checks if the target subvolume is present,
// and creates it using the `btrfs subvolume create` command if missing.
// Returns an error if listing or creation fails.
func createSubvol(path string) error {
	subvol := filepath.Base(path)

	listCmd := exec.Command("btrfs", "subvolume", "list", MountPoint)
	listOutput, err := listCmd.Output()
	if err != nil {
		return fmt.Errorf("failed to list subvolumes: %s", string(listOutput))
	}

	if !strings.Contains(string(listOutput), subvol) {
		createCmd := exec.Command("btrfs", "subvolume", "create", path)
		output, err := createCmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("an error occurred while creating the subvolume: %s", string(output))
		}
	}

	return nil
}

// Setup ensures that the necessary Btrfs subvolumes exist by calling createSubvol
// for predefined paths (DataPath, SnapshotsPath, OldPath).
// It returns an error if any of the subvolume creations fail.
func Setup() error {
	paths := []string{
		DataPath,
		SnapshotsPath,
		OldPath,
	}

	for _, p := range paths {
		err := createSubvol(p)
		if err != nil {
			return err
		}
	}
	return nil
}

// CreateDiskConfig finds the btrfs disk device by scanning the system's mounted filesystems,
// then saves this disk information into a SQLite database as a general configuration.
// It returns an error if the disk cannot be found or if the database operations fail.
func CreateDiskConfig() error {
	disk, diskErr := getDisk()
	if diskErr != nil {
		return errors.New(diskErr.Error())
	}

	db, err := database.GetGormSqliteDb()
	if err != nil {
		return err
	}

	generalConfigStore := store.NewGeneralConfigStore(db)

	generalConfig := models.GeneralConfig{Disk: disk}
	saveErr := generalConfigStore.Save(&generalConfig)
	if saveErr != nil {
		return fmt.Errorf("failed to save general config: %w", saveErr)
	}

	return nil
}
