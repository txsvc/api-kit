package main

import (
	"log"
	"os"
	"sort"

	"github.com/urfave/cli/v2"

	kitcli "github.com/txsvc/apikit/cli"
)

func main() {
	// initialize the CLI
	app := &cli.App{
		Name:     "sc", // simple cli
		Usage:    "github.com/txsvc/apikit demo CLI",
		Commands: setupCommands(),
		Flags:    setupFlags(),
	}
	sort.Sort(cli.FlagsByName(app.Flags))

	// run the CLI
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

// setupCommands returns all custom CLI commands and the default ones
func setupCommands() []*cli.Command {
	cmds := []*cli.Command{
		{
			Name:    "ping",
			Aliases: []string{"p"},
			Usage:   "ping the API service",
			Action:  kitcli.NoOpCommand,
		},
	}

	// merge with default commands
	return kitcli.MergeCommands(cmds, kitcli.WithAuthCommand())
}

// setupCommands returns all global CLI flags and some default ones
func setupFlags() []cli.Flag {
	flags := []cli.Flag{
		&cli.BoolFlag{
			Name:    "quiet",
			Aliases: []string{"q"},
		},
	}

	// merge with global flags
	return kitcli.MergeFlags(flags, kitcli.WithGlobalFlags())
}
