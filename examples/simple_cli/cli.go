package main

import (
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/urfave/cli/v2"

	kit "github.com/txsvc/apikit/cli"
	"github.com/txsvc/apikit/config"
	"github.com/txsvc/apikit/examples"
)

func init() {
	config.InitConfigProvider(examples.NewExampleConfigProvider())
}

func main() {
	// initialize the CLI
	app := &cli.App{
		Name:      config.ShortName(),
		Version:   config.VersionString(),
		Usage:     config.About(),
		Copyright: config.Copyright(),
		Commands:  setupCommands(),
		Flags:     setupFlags(),
		Before: func(c *cli.Context) error {
			// handle global config
			if path := c.String("config"); path != "" {
				config.SetConfigLocation(path)
			}
			return nil
		},
	}
	sort.Sort(cli.FlagsByName(app.Flags))

	// run the CLI
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

//
// CLI commands and flags
//

// setupCommands returns all custom CLI commands and the default ones
func setupCommands() []*cli.Command {
	cmds := []*cli.Command{
		{
			Name:    "ping",
			Aliases: []string{"p"},
			Usage:   "ping the API service",
			Action:  PingCmd,
		},
	}

	// merge with default commands
	return kit.MergeCommands(cmds, kit.WithAuthCommands())
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
	return kit.MergeFlags(flags, kit.WithGlobalFlags())
}

//
// Commands implementations. Usually this would be in its own package
// but as it is an example, I will keep it in just one file for clarity.
//

func PingCmd(c *cli.Context) error {

	ds := config.GetSettings()

	fmt.Println(config.ResolveConfigLocation())
	fmt.Println(ds)

	return nil
}
