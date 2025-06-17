package actions

import (
	"easy-btrfs/database"
	"easy-btrfs/store"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"github.com/urfave/cli/v2"
)

// DeleteSnapshots deletes snapshots based on provided snapshot IDs.
// Returns an error if any snapshot cannot be found or deleted.
func DeleteSnapshots(c *cli.Context) error {
	args := c.Args().Slice()

	if len(args) == 0 {
		return errors.New("at least one snapshot ID is expected, provided 0")
	}

	db, err := database.GetGormSqliteDb()
	if err != nil {
		return err
	}

	snapshotStore := store.NewSnapshotStore(db)

	for _, id := range args {
		intResult, strconvErr := strconv.Atoi(id)
		if strconvErr != nil {
			errMessage := fmt.Sprintf("input contains string value %s; please enter only positive integers", id)
			return errors.New(errMessage)
		}

		snap, result := snapshotStore.FindFirstById(intResult)

		if result.RowsAffected == 0 {
			message := fmt.Sprintf("snapshot not found: %s \n", id)
			return errors.New(message)
		}

		if result.Error != nil {
			return result.Error
		}

		_, err := os.Stat(snap.Path)
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Printf("snapshot not found at path: %s \n", snap.Path)

				err := snapshotStore.Delete(&snap)
				if err != nil {
					return err
				}
			} else {
				return err
			}
		} else {
			cmd := exec.Command("btrfs", "subvol", "delete", snap.Path)
			cmdOutput, cmdErr := cmd.CombinedOutput()
			if cmdErr != nil {
				return errors.New("failed to delete snapshot: " + string(cmdOutput))
			}

			err := snapshotStore.Delete(&snap)
			if err != nil {
				return err
			}

			fmt.Printf("snapshot deleted %s \n", snap.Path)
		}
	}

	return nil
}
