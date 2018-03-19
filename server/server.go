package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/miguelaco/eos/flag"

	"github.com/gorilla/mux"
	"github.com/mitchellh/cli"
)

const host = ""
const port = 1234

type Cmd struct {
	ui    cli.Ui
	flags *flag.FlagSet
	addr  string
	help  string
}

func New(ui cli.Ui) *Cmd {
	cmd := &Cmd{ui: ui}
	cmd.init()
	return cmd
}

func (c *Cmd) init() {
	c.flags = flag.NewFlagSet("server", c.ui)
	c.flags.StringVar(&c.addr, "addr", ":1234", "Sets the HTTP API address to listen on")
	c.help = c.flags.Help(help)
}

func (c *Cmd) Run(args []string) int {
	c.flags.Parse(args)

	log.Printf("Running server command: eos server %v", args)

	router := mux.NewRouter()
	router.HandleFunc("/v1/sys/status", c.healthHandler).Methods(http.MethodGet)

	log.Printf("Listening on: %s", c.addr)
	log.Fatal(http.ListenAndServe(c.addr, router))

	return 0
}

func (c *Cmd) Synopsis() string {
	return synopsis
}

func (c *Cmd) Help() string {
	return c.help
}

func (c *Cmd) healthHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Request: %s", r.URL.Path)

	type Status struct {
		Status string `json:"status"`
	}

	status := Status{"ok"}
	json.NewEncoder(w).Encode(status)

	s, _ := json.Marshal(status)
	log.Printf("Response: %s", s)
}

const synopsis = "Start EOS server"
const help = `
Usage: eos server [options]
  Start EOS operations API server.
`
