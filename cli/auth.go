package cli

import (
	"github.com/urfave/cli/v2"
)

func WithAuthCommand() []*cli.Command {
	return []*cli.Command{
		{
			Name:    "auth",
			Aliases: []string{"a"},
			Usage:   "options to signup and login",
			Subcommands: []*cli.Command{
				{
					Name:   "signup",
					Usage:  "signup with the API service",
					Action: signup,
				},
				{
					Name:   "login",
					Usage:  "authenticate with the API service",
					Action: login,
				},
				{
					Name:   "logout",
					Usage:  "logout from the API service",
					Action: logout,
				},
			},
		},
	}
}

func signup(c *cli.Context) error {
	return NoOpCommand(c)
}

func login(c *cli.Context) error {
	return NoOpCommand(c)
}

func logout(c *cli.Context) error {
	return NoOpCommand(c)
}
