package server

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"

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
	c.flags = flag.NewFlagSet("server", flag.ExitOnError)
	c.flags.StringVar(&c.addr, "addr", ":1234", "Sets the HTTP API address to listen on")
	c.help = c.Usage()
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

func (c *Cmd) Usage() string {
	s := help + "\n" + "Options:\n"
	c.flags.VisitAll(func(f *flag.Flag) {
		s += fmt.Sprintf("  -%s", f.Name) // Two spaces before -; see next two comments.
		name, usage := flag.UnquoteUsage(f)
		if len(name) > 0 {
			s += " " + name
		}
		// Boolean flags of one ASCII letter are so common we
		// treat them specially, putting their usage on the same line.
		if len(s) <= 4 { // space, space, '-', 'x'.
			s += "\t"
		} else {
			// Four spaces before the tab triggers good alignment
			// for both 4- and 8-space tab stops.
			s += "\n    \t"
		}
		s += strings.Replace(usage, "\n", "\n    \t", -1)
		s += fmt.Sprintf(" (default %q)", f.DefValue)
	})

	return s
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
