package install

import (
	"fmt"

	"github.com/mitchellh/cli"
)

type Cmd struct {
	ui cli.Ui
}

func New(ui cli.Ui) *Cmd {
	cmd := &Cmd{ui: ui}
	return cmd
}

func (c *Cmd) Run(args []string) int {
	c.ui.Output(fmt.Sprintf("Running install command: eos install %v", args))
	return 0
}

func (c *Cmd) Synopsis() string {
	return synopsis
}

func (c *Cmd) Help() string {
	return help
}

const synopsis = "Installs EOS with specified configuration"
const help = `
Usage: eos install [options]
  Installs a new EOS cluster as specified in config.
`
