package status

import (
	"log"
	"net/http"
	"io/ioutil"

	"github.com/miguelaco/eos/flag"

	"github.com/mitchellh/cli"
)

type Cmd struct {
	ui cli.Ui
	flags *flag.FlagSet
	addr string
	help string
}

func New(ui cli.Ui) *Cmd {
	cmd := &Cmd{ui: ui}
	cmd.init()
	return cmd
}

func (c *Cmd) init() {
	c.flags = flag.NewFlagSet("status", c.ui)
	c.flags.StringVar(&c.addr, "addr", "http://localhost:1234", "Sets the HTTP API `url` to dial to")
	c.help = c.flags.Help(help)
}

func (c *Cmd) Run(args []string) int {
	c.flags.Parse(args)

	log.Printf("Running status command: eos status %v", args)

	c.status()

	return 0
}

func (c * Cmd) status() {
	url := c.addr + "/v1/sys/status"

	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	status, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%s", status)
}

func (c *Cmd) Synopsis() string {
	return synopsis
}

func (c *Cmd) Help() string {
	return c.help
}

const synopsis = "Installs EOS with specified configuration"
const help = `
Usage: eos install [options]
  Installs a new EOS cluster as specified in config.
`
