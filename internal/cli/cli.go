package cli

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

// NoOpCommand is just a placeholder
func NoOpCommand(c *cli.Context) error {
	return cli.Exit(fmt.Sprintf("%s: command '%s' is not implemented", c.App.Name, c.Command.Name), 0)
}
