package cli

import (
	"github.com/urfave/cli/v2"
)

func WithGlobalFlags() []cli.Flag {
	return make([]cli.Flag, 0) // FIXME add global flags if any
}
