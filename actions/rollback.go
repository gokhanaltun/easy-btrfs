package actions

import (
	"easy-btrfs/database"
	"easy-btrfs/models"
	"easy-btrfs/store"
	"easy-btrfs/utils"
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/urfave/cli/v2"
)

// RollBack rolls back the subvolume to a snapshot based on the provided snapshot ID.
// Returns an error if the snapshot or configuration cannot be found, or if rollback fails.
func RollBack(c *cli.Context) error {

	snapshotId := c.Args().Get(0)

	if snapshotId == "" {
		return errors.New("the 'snapshot id' argument is required and cannot be empty")
	}

	intResult, strconvErr := strconv.Atoi(snapshotId)
	if strconvErr != nil {
		errMessage := fmt.Sprintf("input contains string value %s; please enter only positive integers", snapshotId)
		return errors.New(errMessage)
	}

	db, err := database.GetGormSqliteDb()
	if err != nil {
		return err
	}

	snapshotStore := store.NewSnapshotStore(db)

	snap, snapResult := snapshotStore.FindFirstById(intResult)
	if snapResult.RowsAffected == 0 {
		message := fmt.Sprintf("snapshot not found: %s \n", snapshotId)
		return errors.New(message)
	}

	if snapResult.Error != nil {
		return snapResult.Error
	}

	subvolumeConfigStore := store.NewSubvolumeConfigStore(db)
	config, configResult := subvolumeConfigStore.FindFirstByName(snap.Name)
	if configResult.RowsAffected == 0 {
		message := fmt.Sprintf("configuration not found for snapshot: %s \n", snap.Name)
		return errors.New(message)
	}

	if configResult.Error != nil {
		return configResult.Error
	}

	formattedTime := utils.FormattedTimeNow()

	preSnap := models.Snapshot{
		Name:          config.Name,
		Path:          utils.OldPath + formattedTime,
		SubvolumePath: config.SubvolumePath,
		Pre:           true,
		PreviousFrom:  snap.ID,
	}

	moveCurrentSubvolCmd := exec.Command("mv", utils.MountPoint+config.SubvolumePath, preSnap.Path)
	moveCurrentSubvolCmdOutput, moveCurrentSubvolCmdErr := moveCurrentSubvolCmd.CombinedOutput()
	if moveCurrentSubvolCmdErr != nil {
		return errors.New("failed to move current subvolume: " + string(moveCurrentSubvolCmdOutput))
	}

	saveErr := snapshotStore.Save(&preSnap)
	if saveErr != nil {
		return fmt.Errorf("failed to save backup snapshot to the database: %s : %v", preSnap.Path, saveErr)
	}

	fmt.Println("backup snapshot created before rollback: " + preSnap.Path)

	rollBackCmd := exec.Command("btrfs", "subvol", "snap", snap.Path, utils.MountPoint+config.SubvolumePath)
	rollBackCmdOutput, rollBackCmdErr := rollBackCmd.CombinedOutput()
	if rollBackCmdErr != nil {
		return errors.New("failed to perform rollback: " + string(rollBackCmdOutput))
	}

	fmt.Println("system successfully rolled back")
	fmt.Println("\n a system reboot is required to apply the changes")
	fmt.Println("do you want to reboot now? (y/n)")

	var input string
	fmt.Scanln(&input)

	input = strings.ToLower(input)

	if input == "y" {
		rebootCmd := exec.Command("reboot")
		rebootCmdOutput, rebootCmdErr := rebootCmd.CombinedOutput()
		if rebootCmdErr != nil {
			return errors.New(string(rebootCmdOutput))
		}
	}

	return nil
}
