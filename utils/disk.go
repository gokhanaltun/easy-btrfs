package utils

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// getDisk scans the system's mounted filesystems using the 'df -Th' command
// and searches for a filesystem of type 'btrfs' that is mounted on a root path (ends with '/').
// It returns the device name of the first matching btrfs disk found.
// Returns an error if no btrfs disk is detected or if the command fails.
func getDisk() (string, error) {
	dfCmd := exec.Command("df", "-Th")
	dfCmdOutput, dfCmdErr := dfCmd.CombinedOutput()
	if dfCmdErr != nil {
		return "", fmt.Errorf("failed to listing disks: %s", string(dfCmdOutput))
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

// Mount ensures that the predefined mount point directory exists (creates it if missing),
// then mounts the detected btrfs disk onto that mount point using the 'mount' command.
// Returns an error if directory creation or mounting fails.
func Mount() error {
	_, err := os.Stat(MountPoint)
	if err != nil {
		if os.IsNotExist(err) {
			mkdirCmd := exec.Command("mkdir", MountPoint)
			mkdirCmdOutput, mkdirCmdErr := mkdirCmd.CombinedOutput()
			if mkdirCmdErr != nil {
				return fmt.Errorf("failed to create @ebtrfs path: %s", string(mkdirCmdOutput))
			}
		} else {
			return err
		}
	}

	disk, diskErr := getDisk()
	if diskErr != nil {
		return errors.New(diskErr.Error())
	}

	mntCmd := exec.Command("mount", disk, MountPoint)
	mntCmdOutput, err := mntCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to mount disk: %s", string(mntCmdOutput))
	}

	return nil
}

// Umount unmounts the btrfs filesystem from the predefined mount point using the 'umount' command.
// Returns an error if the unmount operation fails.
func Umount() error {
	umntCmd := exec.Command("umount", MountPoint)
	umntCmdOutput, err := umntCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to unmount disk: %s", string(umntCmdOutput))
	}

	return nil
}
