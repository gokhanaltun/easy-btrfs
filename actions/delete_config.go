package actions

import (
	"easy-btrfs/database"
	"easy-btrfs/models"
	"easy-btrfs/store"
	"errors"
	"fmt"

	"github.com/urfave/cli/v2"
)

// DeleteConfig deletes a subvolume configuration from the database based on the provided name.
// Returns an error if the name is empty or if the configuration cannot be found.
func DeleteConfig(c *cli.Context) error {

	name := c.Args().Get(0)
	if name == "" {
		return errors.New("the 'config name' argument is required and cannot be empty")
	}

	db, err := database.GetGormSqliteDb()
	if err != nil {
		return err
	}

	subvolumeConfigStore := store.NewSubvolumeConfigStore(db)

	result := subvolumeConfigStore.DeleteByField("name", name, &models.SubvolumeConfig{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("config not found " + name)
	}
	fmt.Printf("config deleted %s \n", name)

	return nil
}
