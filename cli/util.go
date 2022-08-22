package cli

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

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

// NoOpCommand is just a placeholder
func NoOpCommand(c *cli.Context) error {
	return cli.Exit(fmt.Sprintf("%s: command '%s' is not implemented", c.App.Name, c.Command.Name), 0)
}
