package cli

import (
	"github.com/urfave/cli/v2"

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
					Action: kit.LoginCommand,
				},
				{
					Name:   "logout",
					Usage:  "logout from the API service",
					Action: kit.LogoutCommand,
				},
			},
		},
	}
}
