package clear

import (
	"devzat/pkg/interfaces"
)

const (
	Name     = "clear"
	argsInfo = ""
	info     = ""
)

type Command struct{}

func (c *Command) Name() string {
	return Name
}

func (c *Command) ArgsInfo() string {
	return argsInfo
}

func (c *Command) Info() string {
	return info
}

func (c *Command) IsRest() bool {
	return false
}

func (c *Command) IsSecret() bool {
	return false
}

func (c *Command) Fn(_ string, u interfaces.User) error {
	_, err := u.Term().Write([]byte("\033[H\033[2J"))

	return err
}
