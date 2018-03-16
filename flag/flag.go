package flag

import (
	"flag"
	"strings"
	"bytes"

	text "github.com/tonnerre/golang-text"
)

type FlagSet struct {
	*flag.FlagSet
}

func NewFlagSet(name string, errorHandling flag.ErrorHandling) *FlagSet {
	f := flag.NewFlagSet(name, errorHandling)
	result := FlagSet(*f)
	return &result
}

func (f FlagSet) Help(head string) string {
	out := new(bytes.Buffer)
	out.WriteString(strings.TrimSpace(head))

	f.VisitAll(func(f *flag.Flag) {
		helpFlag(out, f)
	})

	return strings.TrimRight(out.String(), "\n")
}

func helpFlag(w io.Writer, f *flag.Flag) string {
	example, _ := flag.UnquoteUsage(f)
	if example != "" {
		fmt.Fprintf(w, "  -%s=<%s>\n", f.Name, example)
	} else {
		fmt.Fprintf(w, "  -%s\n", f.Name)
	}

	indented := wrapAtLength(f.Usage, 5)
	fmt.Fprintf(w, "%s\n\n", indented)
}

func wrapAtLength(s string, pad int) string {
	wrapped := text.Wrap(s, maxLineLength-pad)
	lines := strings.Split(wrapped, "\n")
	for i, line := range lines {
		lines[i] = strings.Repeat(" ", pad) + line
	}
	return strings.Join(lines, "\n")
}