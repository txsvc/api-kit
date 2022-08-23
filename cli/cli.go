package cli

import (
	"github.com/urfave/cli/v2"

	"github.com/txsvc/apikit/config"
	kit "github.com/txsvc/apikit/internal/cli"
)

func WithAuthCommands() []*cli.Command {
	return []*cli.Command{
		{
			Name:    "auth",
			Aliases: []string{"a"},
			Usage:   "options to register and login",
			Subcommands: []*cli.Command{
				{
					Name:        "init",
					Usage:       "register with the API service",
					UsageText:   "init email [passphrase]", // FIXME: better description
					Description: "longform description",    // FIXME: better description
					Action:      kit.InitCommand,
				},
				{
					Name:   "login",
					Usage:  "authenticate with the API service",
					Action: kit.NoOpCommand,
				},
				{
					Name:   "logout",
					Usage:  "logout from the API service",
					Action: kit.NoOpCommand,
				},
			},
		},
	}
}

func WithGlobalFlags() []cli.Flag {
	flags := []cli.Flag{
		&cli.StringFlag{
			Name:        "config",
			Usage:       "configuration and secrets directory",
			DefaultText: config.DefaultConfigLocation(),
			Aliases:     []string{"c"},
		},
	}
	return flags
}

// MergeCommands merges all the arrays with CLI commands into one
func MergeCommands(cmds ...[]*cli.Command) []*cli.Command {
	cmd := make([]*cli.Command, 0)
	for _, cc := range cmds {
		cmd = append(cmd, cc...)
	}
	return cmd
}

// MergeFlags merges all the arrays with CLI flags into one
func MergeFlags(flags ...[]cli.Flag) []cli.Flag {
	flag := make([]cli.Flag, 0)
	for _, ff := range flags {
		flag = append(flag, ff...)
	}
	return flag
}
