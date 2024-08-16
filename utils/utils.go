package utils

import (
	"easy-btrfs/models"
	"easy-btrfs/repository"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	MountPoint    string = "/mnt/@ebtrfs/"
	SnapshotsPath string = MountPoint + "@snapshots/"
	OldPath       string = MountPoint + "@old/"
	DataPath      string = MountPoint + "@data/"
)

// Install sets up the disk by finding the appropriate disk for the btrfs filesystem
// and saving the general configuration to a SQLite database. If the disk cannot
// be found or the database operation fails, it returns an error.
func Install() error {
	disk, diskErr := GetDisk()
	if diskErr != nil {
		return errors.New("disk err: " + diskErr.Error())
	}

	generalConfigRepo := repository.NewGeneralConfigRepository()

	generalConfig := models.GeneralConfig{Disk: disk}
	result := generalConfigRepo.Save(&generalConfig)
	if result.Error != nil {
		return errors.New(result.Error.Error())
	}

	return nil
}

// GetDisk scans the mounted filesystems to find a btrfs disk and returns its
// mount point. If no btrfs disk is found, it returns an error. This function uses
// the 'df -Th' command to list filesystems and checks for btrfs types.
func GetDisk() (string, error) {
	dfCmd := exec.Command("df", "-Th")
	dfCmdOutput, dfCmdErr := dfCmd.Output()
	if dfCmdErr != nil {
		return "", fmt.Errorf("error listing disks: %v", dfCmdErr)
	}

	dfCmdOutputLines := strings.Split(string(dfCmdOutput), "\n")

	disk := ""

	for _, line := range dfCmdOutputLines {

		if strings.Contains(line, "btrfs") && strings.HasSuffix(line, "/") {
			disk = strings.Split(line, " ")[0]
			break
		}
	}

	if disk == "" {
		return "", errors.New("btrfs disk not found")
	}

	return disk, nil
}

// Mount mounts the btrfs disk to the predefined mount point. If the mount point
// directory does not exist, it creates it. It uses the 'mount' command to mount
// the disk and returns an error if the mount operation fails.
func Mount() error {
	_, err := os.Stat(MountPoint)
	if err != nil {
		if os.IsNotExist(err) {
			mkdirCmd := exec.Command("mkdir", "/mnt/@ebtrfs")
			mkdirCmdOutput, mkdirCmdErr := mkdirCmd.CombinedOutput()
			if mkdirCmdErr != nil {
				return errors.New(string(mkdirCmdOutput))
			}
		} else {
			return err
		}
	}

	disk, diskErr := GetDisk()
	if err != nil {
		return errors.New("mount err: " + diskErr.Error())
	}

	mntCmd := exec.Command("mount", disk, MountPoint)
	mntCmdOutput, err := mntCmd.CombinedOutput()
	if err != nil {
		return errors.New(string(mntCmdOutput))
	}

	return nil
}

// Umount unmounts the btrfs filesystem from the mount point. It uses the 'umount'
// command and returns an error if the unmount operation fails.
func Umount() error {
	umntCmd := exec.Command("umount", MountPoint)
	umntCmdOutput, err := umntCmd.CombinedOutput()
	if err != nil {
		return errors.New(string(umntCmdOutput))
	}

	return nil
}

// FormattedTimeNow returns the current date and time formatted as "2006-01-02-15:04:05".
// This format is useful for timestamping files or logs.
func FormattedTimeNow() string {
	now := time.Now()
	formattedTime := now.Format("2006-01-02-15:04:05")

	return formattedTime
}
