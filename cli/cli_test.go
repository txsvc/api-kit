package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
)

func TestMergeCommands(t *testing.T) {
	c1 := []*cli.Command{
		{
			Name:    "add",
			Aliases: []string{"a"},
			Usage:   "add a task to the list",
		},
	}
	c2 := []*cli.Command{
		{
			Name:    "sub",
			Aliases: []string{"s"},
		},
	}
	c3 := []*cli.Command{
		{
			Name:    "add",
			Aliases: []string{"a"},
			Usage:   "add a task to the list",
		},
		{
			Name:    "sub",
			Aliases: []string{"s"},
		},
	}

	// simple merge of ONE
	cmd1 := MergeCommands(c1)
	assert.NotEmpty(t, cmd1)
	assert.Equal(t, 1, len(cmd1))

	// simple merge of TWO
	cmd2 := MergeCommands(c1, c2)
	assert.NotEmpty(t, cmd2)
	assert.Equal(t, 2, len(cmd2))

	// merge of all of them
	cmd3 := MergeCommands(c1, c2, c3)
	assert.NotEmpty(t, cmd3)
	assert.Equal(t, 4, len(cmd3))

	// merge with an empty array
	cmd4 := MergeCommands(make([]*cli.Command, 0), c3)
	assert.NotEmpty(t, cmd4)
	assert.Equal(t, 2, len(cmd4))
}

func TestMergeFlags(t *testing.T) {
	f1 := []cli.Flag{
		&cli.StringFlag{
			Name: "flag1",
		},
	}
	f2 := []cli.Flag{
		&cli.StringFlag{
			Name: "flag2",
		},
	}

	flags1 := MergeFlags(f1)
	assert.NotEmpty(t, flags1)
	assert.Equal(t, 1, len(flags1))

	flags2 := MergeFlags(f1, f2)
	assert.NotEmpty(t, flags2)
	assert.Equal(t, 2, len(flags2))

	flags3 := MergeFlags(f1, make([]cli.Flag, 0), f2)
	assert.NotEmpty(t, flags3)
	assert.Equal(t, 2, len(flags3))

}
