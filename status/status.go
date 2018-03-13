package status

import (
	"log"
	"flag"
	"net/http"
	"io/ioutil"

	"github.com/mitchellh/cli"
)

type Cmd struct {
	ui cli.Ui
	flags *flag.FlagSet
	addr string
}

func New(ui cli.Ui) *Cmd {
	cmd := &Cmd{ui: ui}
	cmd.init()
	return cmd
}

func (c *Cmd) init() {
	c.flags = flag.NewFlagSet("server", flag.ContinueOnError)
	c.flags.StringVar(&c.addr, "addr", "http://localhost:1234", "Sets the HTTP API address to dial")
}

func (c *Cmd) Run(args []string) int {
	if err := c.flags.Parse(args); err != nil {
		log.Fatalf("Error parsing args: %e", err)
	}

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
	return help
}

const synopsis = "Installs EOS with specified configuration"
const help = `
Usage: eos install [options]
  Installs a new EOS cluster as specified in config.
`
