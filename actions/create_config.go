package actions

import (
	"easy-btrfs/database"
	"easy-btrfs/models"
	"easy-btrfs/store"
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/urfave/cli/v2"
)

// CreateConfig creates a new subvolume configuration in the database based on provided name and path.
// Returns an error if the name or path is empty, or if the configuration already exists.
func CreateConfig(c *cli.Context) error {
	name := c.Args().Get(0)
	path := c.Args().Get(1)

	if name == "" {
		return errors.New("the 'name' argument is required and cannot be empty")
	}

	if path == "" {
		return errors.New("the 'path' argument is required and cannot be empty")
	}

	showCmd := exec.Command("sudo", "btrfs", "subvolume", "show", path)
	showCmdOutput, showCmdErr := showCmd.CombinedOutput()
	if showCmdErr != nil {
		return errors.New("failed to show" + path + "subvolume details: " + string(showCmdOutput))
	}

	if len(showCmdOutput) == 0 {
		return errors.New("the specified path does not exist or could not be found")
	}

	showCmdOutputLines := strings.Split(string(showCmdOutput), "\n")
	subvolPath := strings.TrimSpace(showCmdOutputLines[0])

	db, err := database.GetGormSqliteDb()
	if err != nil {
		return err
	}

	subvolumeConfigStore := store.NewSubvolumeConfigStore(db)

	countByName, err := subvolumeConfigStore.CountByField("name", name, &models.SubvolumeConfig{})
	if err != nil {
		return err
	}

	if countByName != 0 {
		return fmt.Errorf("a configuration with the name '%s' already exists", name)
	}

	countBySubvolPath, err := subvolumeConfigStore.CountByField("subvolume_path", subvolPath, &models.SubvolumeConfig{})
	if err != nil {
		return err
	}

	if countBySubvolPath != 0 {
		return fmt.Errorf("a configuration with the path '%s' already exists", path)
	}

	subvolumeConfig := &models.SubvolumeConfig{
		Name:          name,
		SubvolumePath: subvolPath,
	}

	saveErr := subvolumeConfigStore.Save(subvolumeConfig)
	if saveErr != nil {
		return fmt.Errorf("database error: configuration could not be saved: %v", saveErr)
	}

	fmt.Printf("%s config created successfuly \n", name)

	return nil
}
