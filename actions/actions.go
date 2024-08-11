package actions

import (
	"easy-btrfs/database"
	"easy-btrfs/models"
	"easy-btrfs/utils"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/markkurossi/tabulate"
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

	db := database.GetGormSqliteDb()

	var count int64
	if err := db.Model(&models.SubvolumeConfig{}).Where("name = ?", name).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return fmt.Errorf("a configuration with the name '%s' already exists", name)
	}

	if err := db.Model(&models.SubvolumeConfig{}).Where("subvolume_path = ?", subvolPath).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return fmt.Errorf("a configuration with the path '%s' already exists", path)
	}

	subvolumeConfig := &models.SubvolumeConfig{
		Name:          name,
		SubvolumePath: subvolPath,
	}

	tx := db.Save(subvolumeConfig)
	if tx == nil {
		return errors.New("database error: configuration could not be saved")
	}

	fmt.Printf("%s config created successfuly \n", name)

	return nil
}

// ListConfigs lists all subvolume configurations from the database.
// Returns an error if the database operation fails.
func ListConfigs(c *cli.Context) error {

	db := database.GetGormSqliteDb()
	configs := []models.SubvolumeConfig{}
	result := db.Find(&configs)
	if result.Error != nil {
		return errors.New("failed to retrieve configurations from the database: " + result.Error.Error())
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

// DeleteConfig deletes a subvolume configuration from the database based on the provided name.
// Returns an error if the name is empty or if the configuration cannot be found.
func DeleteConfig(c *cli.Context) error {

	name := c.Args().Get(0)
	if name == "" {
		return errors.New("the 'config name' argument is required and cannot be empty")
	}

	db := database.GetGormSqliteDb()
	result := db.Model(&models.SubvolumeConfig{}).Where("name = ?", name).Delete(&models.SubvolumeConfig{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("config not found " + name)
	}
	fmt.Printf("config deleted %s \n", name)

	return nil
}

// Snapshot creates a snapshot of a subvolume based on the provided configuration name.
// Returns an error if the configuration does not exist, or if the snapshot creation fails.
func Snapshot(c *cli.Context) error {
	configName := c.Args().Get(0)
	if configName == "" {
		return errors.New("the 'config name' argument is required and cannot be empty")
	}

	config := models.SubvolumeConfig{}
	db := database.GetGormSqliteDb()
	result := db.Where("name = ?", configName).First(&config)
	if result.Error != nil {
		if result.RowsAffected == 0 {
			return fmt.Errorf("config not found %s", configName)
		}
		return result.Error
	}

	formattedTime := utils.FormattedTimeNow()

	snapshot := models.Snapshot{
		Name:          config.Name,
		Path:          utils.SnapshotsPath + formattedTime,
		SubvolumePath: config.SubvolumePath,
		Pre:           false,
	}

	snapSaveResult := db.Save(&snapshot)
	if snapSaveResult.Error != nil {
		return errors.New("failed to save snapshot to the database: " + snapSaveResult.Error.Error())
	}

	snapCmd := exec.Command("btrfs", "subvol", "snap", "-r", utils.MountPoint+config.SubvolumePath, snapshot.Path)
	snapCmdOutput, snapCmdErr := snapCmd.CombinedOutput()
	if snapCmdErr != nil {
		return errors.New("failed to create snapshot: " + string(snapCmdOutput))
	}

	fmt.Printf("%s snapshot created: %s \n", config.Name, utils.SnapshotsPath+formattedTime)

	return nil
}

// ListSnapshots lists all snapshots from the database or, if a configuration name is provided,
// lists snapshots associated with that configuration.
// Returns an error if snapshot records cannot be retrieved or if no snapshots are found.
func ListSnapshots(c *cli.Context) error {

	configName := c.Args().Get(0)
	if configName == "" {
		snaps := []models.Snapshot{}
		db := database.GetGormSqliteDb()
		result := db.Find(&snaps)
		if result.RowsAffected == 0 {
			return errors.New("no snapshots found")
		}

		if result.Error != nil {
			return result.Error
		}

		tb := tabulate.New(tabulate.ASCII)
		tb.Header("ID").SetAlign(tabulate.BL)
		tb.Header("Name").SetAlign(tabulate.BL)
		tb.Header("Path").SetAlign(tabulate.BL)
		tb.Header("Pre").SetAlign(tabulate.BL)
		tb.Header("PreFrom").SetAlign(tabulate.BL)

		for _, snap := range snaps {
			row := tb.Row()
			row.Column(fmt.Sprint(snap.ID))
			row.Column(snap.Name)
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
		db := database.GetGormSqliteDb()
		result := db.Find(&models.SubvolumeConfig{})
		if result.RowsAffected == 0 {
			errMessage := fmt.Sprintf("%s configuration not found", configName)
			return errors.New(errMessage)
		}
		if result.Error != nil {
			return result.Error
		}

		snaps := []models.Snapshot{}
		snapsResult := db.Where("name = ?", configName).Find(&snaps)
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
		tb.Header("Path").SetAlign(tabulate.BL)

		for _, snap := range snaps {
			row := tb.Row()
			row.Column(fmt.Sprint(snap.ID))
			row.Column(snap.Name)
			row.Column(snap.Path)
		}

		tb.Print(os.Stdout)
	}

	return nil
}

// DeleteSnapshots deletes snapshots based on provided snapshot IDs.
// Returns an error if any snapshot cannot be found or deleted.
func DeleteSnapshots(c *cli.Context) error {
	args := c.Args().Slice()

	if len(args) == 0 {
		return errors.New("at least one snapshot ID is expected, provided 0")
	}

	db := database.GetGormSqliteDb()

	for _, id := range args {
		snap := models.Snapshot{}
		result := db.Where("id = ?", id).First(&snap)

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
				message := fmt.Sprintf("snapshot not found at path: %s \n", snap.Path)
				return errors.New(message)

			} else {
				return err
			}
		} else {
			cmd := exec.Command("btrfs", "subvol", "delete", snap.Path)
			cmdOutput, cmdErr := cmd.CombinedOutput()
			if cmdErr != nil {
				return errors.New("failed to delete snapshot: " + string(cmdOutput))
			}

			result := db.Delete(&snap)
			if result.Error != nil {
				return result.Error
			}

			fmt.Printf("snapshot deleted %s \n", snap.Path)
		}
	}

	return nil
}

// RollBack rolls back the subvolume to a snapshot based on the provided snapshot ID.
// Returns an error if the snapshot or configuration cannot be found, or if rollback fails.
func RollBack(c *cli.Context) error {

	snapshotId := c.Args().Get(0)

	if snapshotId == "" {
		return errors.New("the 'snapshot id' argument is required and cannot be empty")
	}

	db := database.GetGormSqliteDb()

	snap := models.Snapshot{}
	snapResult := db.Where("id = ?", snapshotId).First(&snap)
	if snapResult.RowsAffected == 0 {
		message := fmt.Sprintf("snapshot not found: %s \n", snapshotId)
		return errors.New(message)
	}

	if snapResult.Error != nil {
		return snapResult.Error
	}

	config := models.SubvolumeConfig{}
	configResult := db.Where("id = ?", snapshotId).First(&config)
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

	preSnapSaveResult := db.Save(&preSnap)
	if preSnapSaveResult.Error != nil {
		return errors.New("failed to save backup snapshot to the database: " + preSnap.Path + " : " + preSnapSaveResult.Error.Error())
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
