package install

import (
	"fmt"
    "flag"

	"github.com/mitchellh/cli"
)

type Cmd struct {
	ui cli.Ui

    // flags
    flags *flag.FlagSet
    verbose bool
    securize bool
    principal string
}

func New(ui cli.Ui) *Cmd {
	cmd := &Cmd{ui: ui}
    cmd.init()
	return cmd
}

func (c *Cmd) init() {
    c.flags = flag.NewFlagSet("install", flag.ContinueOnError)
    c.flags.BoolVar(&c.verbose, "v", false,
        "Increase output info.")
    c.flags.BoolVar(&c.securize, "securize", false,
        "Securize cluster after install")
    c.flags.StringVar(&c.principal, "principal", "root/admin",
        "Principal to use when generating kerberos credentials")
}

func (c *Cmd) Run(args []string) int {
    if err := c.flags.Parse(args); err != nil {
        return 1
    }

    c.flags.VisitAll(func (flag *flag.Flag) { fmt.Printf("Visiting flag %v = %v\n", flag.Name, flag.Value) })
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
