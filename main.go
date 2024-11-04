package main

import (
	"fmt"
	cli "gator/internal/cli"
	internal "gator/internal/config"
	state "gator/internal/state"
	"os"
)

func main() {
	config := internal.Read()
	st := state.State{Config: &config}
	commands := cli.Commands{Registered: make(map[string]func(*state.State, cli.Command) error)}

	commands.Register("login", cli.HandlerLogin)

	args := os.Args[1:]

	if len(args) < 2 {
		fmt.Println("minimum 2 arguments required")
		os.Exit(1)
	}

	command := args[0]
	commandArgs := args[1:]

	err := commands.Run(&st, cli.Command{Name: command, Args: commandArgs})

	if err != nil {
		fmt.Println(err)
	}
}
