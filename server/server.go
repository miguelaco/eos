package server

import (
	"log"
	"flag"
	"net/http"
	"encoding/json"

	"github.com/mitchellh/cli"
	"github.com/gorilla/mux"
)

const host = ""
const port = 1234

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
	c.flags.StringVar(&c.addr, "addr", ":1234", "Sets the HTTP API address to listen on")
}

func (c *Cmd) Run(args []string) int {
	if err := c.flags.Parse(args); err != nil {
		log.Fatalf("Error parsing args: %e", err)
	}

	log.Printf("Running server command: eos server %v", args)

	router := mux.NewRouter()
	router.HandleFunc("/v1/sys/health", c.healthHandler).Methods(http.MethodGet)

	log.Printf("Listening on: %s", c.addr)
	log.Fatal(http.ListenAndServe(c.addr, router))

	return 0
}

func (c *Cmd) Synopsis() string {
	return synopsis
}

func (c *Cmd) Help() string {
	return help
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
