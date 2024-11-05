package cli

import (
	"context"
	"errors"
	"fmt"
	"gator/internal/database"
	"gator/internal/state"
	"time"

	"github.com/google/uuid"
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

	user, err := s.Db.GetUser(context.Background(), username)

	if err != nil {
		return err
	}

	s.Config.SetUser(user.Name)

	fmt.Printf("%s has been set", user.Name)
	return nil
}

func HandlerRegister(s *state.State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return errors.New("not enough arguments")
	}

	user := cmd.Args[0]

	newUser := database.CreateUserParams{
		Name:      user,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		ID:        uuid.New(),
	}

	createdUser, err := s.Db.CreateUser(context.Background(), newUser)

	if err != nil {
		return err
	}

	s.Config.SetUser(createdUser.Name)

	fmt.Printf("%s has been created, id: %v", createdUser.Name, createdUser.ID)
	return nil
}

func HandlerReset(s *state.State, cmd Command) error {
	err := s.Db.DeleteAllUsers(context.Background())

	if err != nil {
		return err
	}

	fmt.Println("users have been deleted")
	return nil
}

func HandlerGetUsers(s *state.State, cmd Command) error {
	users, err := s.Db.GetUsers(context.Background())

	if err != nil {
		return err
	}

	for _, user := range users {
		if user.Name == s.Config.Current_user_name {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Printf("* %s\n", user.Name)
		}
	}

	return nil

}
