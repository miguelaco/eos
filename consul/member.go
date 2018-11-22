package consul

import (
	"strings"
)

const (
	StatusAlive  = 1
	StatusFailed = 4
	TypeMaster   = "master"
	TypeGosec    = "gosec"
	TypeAgent    = "agent"
	TypeNA       = "N/A"
)

var statusText = map[int]string{
	StatusAlive:  "alive",
	StatusFailed: "failed",
}

type Member struct {
	Name   string
	Addr   string
	Status int
}

func (m *Member) StatusText() string {
	if text, ok := statusText[m.Status]; ok {
		return text
	}

	return statusText[StatusFailed]
}

func (m *Member) Type() string {
	switch {
	case strings.HasPrefix(m.Name, TypeMaster):
		return TypeMaster
	case strings.HasPrefix(m.Name, TypeGosec):
		return TypeGosec
	case strings.HasPrefix(m.Name, TypeAgent):
		return TypeAgent
	default:
		return TypeNA
	}
}

type MemberMap map[string]Member

type MemberSlice []Member

func (m MemberSlice) Len() int           { return len(m) }
func (m MemberSlice) Less(i, j int) bool { return m[i].Name < m[j].Name }
func (m MemberSlice) Swap(i, j int)      { m[i], m[j] = m[j], m[i] }
