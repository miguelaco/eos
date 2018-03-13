package install

import (
	"log"
	"time"
	"math/rand"

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
	log.Printf("Running install command: eos install %v", args)

	rand.Seed(time.Now().UnixNano())
	start := make(chan int, 5)
	done := make(chan int, 10)

	for i := 0; i < 10; i++ {
		go func(i int) {
			start <- i
			log.Printf("%d: Start\n", i)
			c.install(i)
			done <- i
			<-start
		}(i)
	}


	for i := 0; i < 10; i++ {
		j := <- done
		log.Printf("%d: Done\n", j)
	}

	log.Println("All done")

	return 0
}

func (c * Cmd) install(id int) {
	s := time.Duration(rand.Intn(9) + 1)
	log.Printf("%d: Sleeping %d seconds\n", id, s)
	time.Sleep(s * time.Second)
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
