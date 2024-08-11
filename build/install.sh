#!/bin/bash

echo "Starting easy-btrfs installation..."

NAME="ebtrfs"
TARGET_PATH="/usr/local/bin"
MOUNT_POINT="/mnt/@ebtrfs"
DATA_PATH="$MOUNT_POINT/@data"
SNAPSHOTS_PATH="$MOUNT_POINT/@snapshots"
OLD_PATH="$MOUNT_POINT/@old"

echo "Retrieving disk information..."
disk=$(df -Th | grep btrfs | grep /$ | cut -d ' ' -f 1)
if [ $? -ne 0 ]; then
    echo "An error occurred while retrieving disk information."
    exit 1
fi

if [ -z "$disk" ]; then
    echo "Disk information could not be retrieved. The system might not support Btrfs."
    exit 1
fi

echo "Creating mount point..."
mkdir -p "$MOUNT_POINT"
if [ $? -ne 0 ]; then
    echo "An error occurred while creating the mount point."
    exit 1
fi

echo "Mounting disk..."
mount "$disk" "$MOUNT_POINT"
if [ $? -ne 0 ]; then
    echo "An error occurred while mounting the disk."
    exit 1
fi

create_subvol() {
    subvol=$(basename "$1")
    result=$(btrfs subvolume list "$MOUNT_POINT" | grep "$subvol")
    if [ -z "$result" ]; then
        btrfs subvolume create "$1"
        if [ $? -ne 0 ]; then 
            echo "An error occurred while creating the subvolume."
            exit 1     
        fi
    fi
}

echo "Creating necessary directories (subvolumes)..."
create_subvol "$DATA_PATH"
create_subvol "$SNAPSHOTS_PATH"
create_subvol "$OLD_PATH"

echo "Performing umount operation..."
umount "$MOUNT_POINT"
if [ $? -ne 0 ]; then
    echo "An error occurred while performing the umount operation."
    exit 1
fi

echo "Moving the program to the installation directory..."
mv "$NAME" "$TARGET_PATH/$NAME"
if [ $? -ne 0 ]; then
    echo "An error occurred while moving the program to the target directory."
    exit 1
fi

echo "Setting execution permissions..."
chmod +x "$TARGET_PATH/$NAME"
if [ $? -ne 0 ]; then
    echo "An error occurred while setting execution permissions."
    exit 1
fi

echo "Installation completed successfully."
