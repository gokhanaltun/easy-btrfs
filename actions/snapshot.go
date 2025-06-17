package actions

import (
	"easy-btrfs/database"
	"easy-btrfs/models"
	"easy-btrfs/store"
	"easy-btrfs/utils"
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/urfave/cli/v2"
)

// Snapshot creates a snapshot of a subvolume based on the provided configuration name.
// Returns an error if the configuration does not exist, or if the snapshot creation fails.
func Snapshot(c *cli.Context) error {
	configName := c.Args().Get(0)
	if configName == "" {
		return errors.New("the 'config name' argument is required and cannot be empty")
	}

	description := "-"
	if len(c.Args().Slice()) > 1 {
		description = strings.Join(c.Args().Slice()[1:], " ")
	}

	db, err := database.GetGormSqliteDb()
	if err != nil {
		return err
	}

	subvolumeConfigStore := store.NewSubvolumeConfigStore(db)
	config, result := subvolumeConfigStore.FindFirstByName(configName)

	if result.Error != nil {
		if result.RowsAffected == 0 {
			return fmt.Errorf("config not found %s", configName)
		}
		return result.Error
	}

	formattedTime := utils.FormattedTimeNow()

	snapshot := models.Snapshot{
		Name:          config.Name,
		Description:   description,
		Path:          utils.SnapshotsPath + formattedTime,
		SubvolumePath: config.SubvolumePath,
		Pre:           false,
	}

	snapshotStore := store.NewSnapshotStore(db)
	saveErr := snapshotStore.Save(&snapshot)
	if saveErr != nil {
		return fmt.Errorf("failed to save snapshot to the database: %v", saveErr)
	}

	snapCmd := exec.Command("btrfs", "subvol", "snap", "-r", utils.MountPoint+config.SubvolumePath, snapshot.Path)
	snapCmdOutput, snapCmdErr := snapCmd.CombinedOutput()
	if snapCmdErr != nil {
		return errors.New("failed to create snapshot: " + string(snapCmdOutput))
	}

	fmt.Printf("%s snapshot created: %s \n", config.Name, utils.SnapshotsPath+formattedTime)

	return nil
}
