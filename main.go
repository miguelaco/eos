package main

import (
	"os"

	"github.com/mitchellh/cli"

	"github.com/miguelaco/eos/install"
	"github.com/miguelaco/eos/status"
	"github.com/miguelaco/eos/server"
)

func main() {
	ui := &cli.BasicUi{Writer: os.Stdout, ErrorWriter: os.Stderr}

	c := cli.NewCLI("eos", "1.0.0")
	c.Args = os.Args[1:]
	c.Commands = map[string]cli.CommandFactory{
		"install": func() (cli.Command, error) { return install.New(ui), nil },
		"status": func() (cli.Command, error) { return status.New(ui), nil },
		"server": func() (cli.Command, error) { return server.New(ui), nil },
	}

	exitStatus, err := c.Run()
	if err != nil {
		ui.Error(err.Error())
	}

	os.Exit(exitStatus)
}
