package cli

import (
	"errors"
	"fmt"
	"gator/internal/state"
)

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	Registered map[string]func(*state.State, Command) error
}

func (c *Commands) Register(name string, f func(*state.State, Command) error) {
	_, ok := c.Registered[name]

	if ok {
		fmt.Println("command has already been registered")
	} else {
		c.Registered[name] = f
	}
}

func (c *Commands) Run(s *state.State, cmd Command) error {
	command, ok := c.Registered[cmd.Name]

	if ok {
		err := command(s, cmd)
		return err
	} else {
		return errors.New("command not found")
	}
}

func HandlerLogin(s *state.State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return errors.New("not enough arguments")
	}

	username := cmd.Args[0]

	s.Config.SetUser(username)

	fmt.Printf("%s has been set", username)
	return nil
}
