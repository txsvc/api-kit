package main

import (
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/urfave/cli/v2"

	"github.com/txsvc/apikit/api"
	kit "github.com/txsvc/apikit/cli"
	"github.com/txsvc/apikit/config"
	"github.com/txsvc/cloudlib/observer"
)

func init() {
	config.SetProvider(config.NewLocalConfigProvider())
}

func main() {
	// initialize the CLI
	cfg := config.GetConfig()
	app := &cli.App{
		Name:      cfg.Info().ShortName(),
		Version:   cfg.Info().VersionString(),
		Usage:     cfg.Info().About(),
		Copyright: cfg.Info().Copyright(),
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
// but as this is an example, I will keep it in just one file for clarity.
//

func PingCmd(c *cli.Context) error {

	cl := api.NewClient(nil)
	if cl == nil {
		return fmt.Errorf("could not create a client instance")
	}

	var so api.StatusObject
	if status, err := cl.GET("/ping", &so); err != nil {
		observer.LogWithLevel(observer.LevelError, fmt.Sprintf("status: %d: %s", status, err))
		return nil
	}

	observer.Log(fmt.Sprintf("%v\n", so))

	return nil
}
