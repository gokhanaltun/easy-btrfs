# easy-btrfs
A user-friendly Btrfs CLI tool for managing snapshots.

## Table of Contents
- [Important Note](#important-note)
- [Installation AMD64 (x86_64)](#installation-amd64)
- [Manual Installation](#manual-installation)
- [Commands and Usage](#commands-and-usage)
  - [Create Config](#create-config)
  - [List Configs](#list-configs)
  - [Delete Config](#delete-config)
  - [Snapshot](#snapshot)
  - [List Snapshots](#list-snapshots)
  - [Delete Snapshots](#delete-snapshots)
  - [RollBack](#rollback)
- [Old Backups](#old-backups)
- Latest version: **v1.1.2** — [See what's new](#-whats-new-in-v112)


## Important Note

For this program to function correctly, the default setting of the BTRFS file system must be `ID 5 (FS_TREE)`. Please follow the steps below to verify or configure this setting:

To check the default setting, use the following command:
```bash
sudo btrfs subvolume get-default /
```
If the output of this command is not `ID 5 (FS_TREE)`, use the following command to set it and then restart your system:
```bash
sudo btrfs subvolume set-default 5 /
```

## Arch Linux (via AUR)
```bash
yay -S easy-btrfs
```

## Installation via DEB Package (Debian-based systems)

You can download the latest .deb package from the [Releases](https://github.com/gokhanaltun/easy-btrfs/releases)
 page and install it using:
```bash
sudo dpkg -i easy-btrfs_1.1.2_amd64.deb
```

## Manual Installation (For other distros or custom installs)

- Download the .tar.gz archive from the [Releases](https://github.com/gokhanaltun/easy-btrfs/releases) page.

- Extract and install manually:

```bash
tar -xzf easy-btrfs_1.1.2_amd64.tar.gz
cd easy-btrfs_1.1.2/
sudo chmod +x ./install.sh
sudo ./install.sh
```

After this, you can run ebtrfs from your terminal.


## Commands and Usage
This section will explain the commands and how to use them.

### Create Config
The `create-config` command creates a subvolume configuration that allows you to easily take snapshots.  
Aliases: `cc`, `cconf`  
It requires two arguments: `name` and `path`.  
- **name:** A name for the configuration (you define this name).  
- **path:** The path to the subvolume for which the configuration will be created, e.g., `/home` for the home directory.

Example usage: Let's create a config for the root.
```bash
sudo ebtrfs create-config root /
```

### List Configs
The `list-configs` command lists all the configurations you have created. It does not take any parameters.

Aliases: `ls-c`, `lsc`, `lc`

Example usage:
```bash
sudo ebtrfs list-configs
```

### Delete Config
The `delete-config` command allows you to delete a configuration. It takes a `name` parameter.

- **name:** The name of the configuration you want to delete.
- **Aliases:** `del-conf`, `d-conf`, `dc`

Example usage: Let's delete the root config
```bash
sudo ebtrfs delete-config root
```

### Snapshot
The `snapshot` command creates a read-only snapshot for a subvolume. It takes one required argument and one optional description argument:

- **name**: The name of the config for which you want to create the snapshot.
- **description**: (optional) A description for the snapshot, which can be multiple words without needing quotes.
- **Aliases:** `snap`, `s`

Example usage: Let's create a snapshot of root config
```bash
sudo ebtrfs snapshot root This is a root snapshot
```

or with quotes

```bash
sudo ebtrfs snapshot root "This is a root snapshot"
```

### List Snapshots
The `list-snapshots` command lists all snapshots or snapshots for a specific configuration.

- **name:** Optional. If provided, only snapshots for this configuration will be listed.
- **Aliases:** `list-snaps`, `ls-s`, `lss`

**Example usage:**
```bash
# List all snapshots
sudo ebtrfs list-snapshots

# List snapshots for the 'root' configuration
sudo ebtrfs list-snapshots root
```

### Delete Snapshots
The `delete-snapshots` command deletes snapshots. It takes the `id` parameter as an argument:
- **id:** The ID of the snapshot to be deleted. You can specify multiple `id` arguments to delete multiple snapshots at once.

Aliases: `del-snaps`, `d-snaps`, `ds`

**Example Usage:**

```bash
#To delete a single snapshot:
sudo ebtrfs delete-snapshots 8

#To delete multiple snapshots:
sudo ebtrfs delete-snapshots 1 2 3 4 5 6 7 8
```

### RollBack
The `rollback` command is used to restore a snapshot. It takes an `id` parameter.
- `id`: The ID of the snapshot you want to restore.
- Aliases: `rb`, `r`

**Example Usage:**

```bash
sudo ebtrfs rollback 5
```

During the rollback process, the current subvolume is moved to the @old directory, and the selected snapshot for the rollback is loaded in place of the moved subvolume. To apply these changes, the system needs to be restarted. After the rollback is complete, the program will ask if you want to restart, and if you choose to restart, the system will automatically reboot.

## Old Backups
For the rollback operation, when you run the command, the program will move the subvolume to be restored before the rollback to the `@old` directory and record which snapshot it was created before. This is done to reduce the risk of data loss and to allow the user to revert to the pre-rollback state at any time.

To list old backups, use the `list-snapshots` command.

In the list, snapshots with the `Pre` column set to `true` are old backups, and the `PreFrom` column contains a snapshot ID indicating which snapshot this backup was created before.

For example, in the following output, the snapshot with ID 4 is an old backup, and the snapshot with ID 3 is a backup taken just before the rollback without any rollback.


```bash
+----+------+---------------------------------------------+---------------------------+-------+---------+
| ID | Name | Path                                        | Description               | Pre   | PreFrom |
+----+------+---------------------------------------------+---------------------------+-------+---------+
| 3  | root | /mnt/@ebtrfs/@snapshots/2024-08-11-19:14:25 | Daily backup before update | false | -       |
| 4  | root | /mnt/@ebtrfs/@old/2024-08-11-19:14:45       | Backup before rollback     | true  | 3       |
+----+------+---------------------------------------------+---------------------------+-------+---------+

```


## 📦 What's New in v1.1.2
### 🛠 Improvements
* ✍️ Added missing descriptions to command help texts.