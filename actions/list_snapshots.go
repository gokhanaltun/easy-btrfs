package actions

import (
	"easy-btrfs/database"
	"easy-btrfs/store"
	"errors"
	"fmt"
	"os"

	"github.com/markkurossi/tabulate"
	"github.com/urfave/cli/v2"
)

// ListSnapshots lists all snapshots from the database or, if a configuration name is provided,
// lists snapshots associated with that configuration.
// Returns an error if snapshot records cannot be retrieved or if no snapshots are found.
func ListSnapshots(c *cli.Context) error {

	db, err := database.GetGormSqliteDb()
	if err != nil {
		return err
	}

	snapshotStore := store.NewSnapshotStore(db)

	configName := c.Args().Get(0)
	if configName == "" {
		snaps, result := snapshotStore.FindAll()
		if result.RowsAffected == 0 {
			return errors.New("no snapshots found")
		}

		if result.Error != nil {
			return result.Error
		}

		tb := tabulate.New(tabulate.ASCII)
		tb.Header("ID").SetAlign(tabulate.BL)
		tb.Header("Name").SetAlign(tabulate.BL)
		tb.Header("Description").SetAlign(tabulate.BL)
		tb.Header("Path").SetAlign(tabulate.BL)
		tb.Header("Pre").SetAlign(tabulate.BL)
		tb.Header("PreFrom").SetAlign(tabulate.BL)

		for _, snap := range snaps {
			if snap.Description == "" {
				snap.Description = "-"
			}

			row := tb.Row()
			row.Column(fmt.Sprint(snap.ID))
			row.Column(snap.Name)
			row.Column(snap.Description)
			row.Column(snap.Path)

			if snap.Pre {
				row.Column("true")
				row.Column(fmt.Sprint(snap.PreviousFrom))
			} else {
				row.Column("false")
				row.Column("-")
			}
		}

		tb.Print(os.Stdout)

	} else {
		subvolumeConfigStore := store.NewSubvolumeConfigStore(db)
		_, result := subvolumeConfigStore.FindFirstByName(configName)
		if result.RowsAffected == 0 {
			errMessage := fmt.Sprintf("%s configuration not found", configName)
			return errors.New(errMessage)
		}
		if result.Error != nil {
			return result.Error
		}

		snaps, snapsResult := snapshotStore.FindAllByName(configName)
		if snapsResult.RowsAffected == 0 {
			errMessage := fmt.Sprintf("no snapshots found for configuration: %s \n", configName)
			return errors.New(errMessage)
		}
		if snapsResult.Error != nil {
			return snapsResult.Error
		}

		tb := tabulate.New(tabulate.ASCII)
		tb.Header("ID").SetAlign(tabulate.BL)
		tb.Header("Name").SetAlign(tabulate.BL)
		tb.Header("Description").SetAlign(tabulate.BL)
		tb.Header("Path").SetAlign(tabulate.BL)

		for _, snap := range snaps {
			row := tb.Row()
			row.Column(fmt.Sprint(snap.ID))
			row.Column(snap.Name)
			row.Column(snap.Description)
			row.Column(snap.Path)
		}

		tb.Print(os.Stdout)
	}

	return nil
}
