package flag

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"strings"

	"github.com/mitchellh/cli"
	text "github.com/tonnerre/golang-text"
)

type FlagSet struct {
	*flag.FlagSet
	ui cli.Ui
}

func NewFlagSet(name string, ui cli.Ui) *FlagSet {
	f := flag.NewFlagSet(name, flag.ExitOnError)
	result := FlagSet{f, ui}
	result.Usage = result.customUsage

	return &result
}

func (f FlagSet) Help(head string) string {
	out := new(bytes.Buffer)
	out.WriteString(strings.TrimSpace(head))
	out.WriteString("\n\nOptions:\n")

	f.VisitAll(func(f *flag.Flag) {
		helpFlag(out, f)
	})

	return strings.TrimRight(out.String(), "\n")
}

func (f FlagSet) customUsage() {
	out := new(bytes.Buffer)
	fmt.Fprintf(out, "Usage of %s:\n", f.Name())

	f.VisitAll(func(f *flag.Flag) {
		helpFlag(out, f)
	})

	s := strings.TrimRight(out.String(), "\n")
	f.ui.Error(s)
}

func helpFlag(w io.Writer, f *flag.Flag) {
	example, usage := flag.UnquoteUsage(f)
	if example != "" {
		fmt.Fprintf(w, "  -%s=<%s>\n", f.Name, example)
	} else {
		fmt.Fprintf(w, "  -%s\n", f.Name)
	}

	indented := wrapAtLength(usage, 8)
	fmt.Fprintf(w, "%s\n\n", indented)
}

func wrapAtLength(s string, pad int) string {
	wrapped := text.Wrap(s, 120-pad)
	lines := strings.Split(wrapped, "\n")
	for i, line := range lines {
		lines[i] = strings.Repeat(" ", pad) + line
	}
	return strings.Join(lines, "\n")
}
