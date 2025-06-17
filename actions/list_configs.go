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

// ListConfigs lists all subvolume configurations from the database.
// Returns an error if the database operation fails.
func ListConfigs(c *cli.Context) error {
	db, err := database.GetGormSqliteDb()
	if err != nil {
		return err
	}

	subvolumeConfigStore := store.NewSubvolumeConfigStore(db)

	configs, result := subvolumeConfigStore.FindAll()
	if result.Error != nil {
		return fmt.Errorf("failed to retrieve configurations from the database: %v", result.Error)
	}

	if len(configs) == 0 {
		return errors.New("no configurations found")
	}

	tab := tabulate.New(tabulate.ASCII)
	tab.Header("ID").SetAlign(tabulate.BL)
	tab.Header("Name").SetAlign(tabulate.BL)
	tab.Header("Subvol Path").SetAlign(tabulate.BL)

	for _, config := range configs {
		row := tab.Row()
		row.Column(fmt.Sprint(config.ID))
		row.Column(config.Name)
		row.Column(config.SubvolumePath)
	}

	tab.Print(os.Stdout)

	return nil
}
