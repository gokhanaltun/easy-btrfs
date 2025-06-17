package main

import (
	"easy-btrfs/commands"
	"easy-btrfs/database"
	"easy-btrfs/store"
	"easy-btrfs/utils"
	"errors"
	"fmt"
	"log"
	"os"
	"os/user"

	"github.com/urfave/cli/v2"
	"gorm.io/gorm"
)

func main() {
	app := &cli.App{
		Name:  "easy-btrfs",
		Usage: "A user-friendly Btrfs CLI tool for managing snapshots and subvolume configurations.",
		Before: func(ctx *cli.Context) error {

			currentUser, err := user.Current()
			if err != nil {
				fmt.Fprintln(os.Stderr, "Failed to get user info:", err)
				os.Exit(1)
			}

			if currentUser.Uid != "0" {
				fmt.Fprintln(os.Stderr, "This program must be run as root. Please use sudo.")
				os.Exit(1)
			}

			if err := utils.Mount(); err != nil {
				fmt.Fprintln(os.Stderr, "Mount error:", err)
				os.Exit(1)
			}

			db, err := database.GetGormSqliteDb()
			if err != nil {
				fmt.Fprintln(os.Stderr, "Failed to initialize database:", err)
				os.Exit(1)
			}

			generalConfigStore := store.NewGeneralConfigStore(db)
			count, err := generalConfigStore.Count()
			if (err != nil && errors.Is(err, gorm.ErrRecordNotFound)) || count == 0 {
				if err := utils.Install(); err != nil {
					fmt.Fprintln(os.Stderr, err)
					os.Exit(1)
				}
			}

			return nil
		},
		After: func(ctx *cli.Context) error {
			dbInstance, err := database.GetGormSqliteDb()
			if err != nil {
				return fmt.Errorf("failed to initialize database: %v", err)
			}

			db, err := dbInstance.DB()
			if err != nil {
				return err
			}
			db.Close()

			umntErr := utils.Umount()
			if umntErr != nil {
				return errors.New("umount err: " + umntErr.Error())
			}

			return nil
		},
		Commands: []*cli.Command{
			commands.CreateConfig(),
			commands.ListConfigs(),
			commands.DeleteConfig(),
			commands.Snapshot(),
			commands.ListSnapshots(),
			commands.DeleteSnapshots(),
			commands.RollBack(),
		},
		CommandNotFound: func(ctx *cli.Context, cmd string) {
			fmt.Fprintf(os.Stderr, "Error: '%s' is not a valid command.\n", cmd)
			fmt.Println("Use 'easy-btrfs --help' or 'easy-btrfs <command> --help' for more information.")
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
