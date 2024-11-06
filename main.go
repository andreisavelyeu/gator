package main

import (
	"database/sql"
	"fmt"
	cli "gator/internal/cli"
	internal "gator/internal/config"
	database "gator/internal/database"
	middleware "gator/internal/middleware"
	state "gator/internal/state"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	config := internal.Read()
	st := state.State{Config: &config}
	commands := cli.Commands{Registered: make(map[string]func(*state.State, cli.Command) error)}
	db, err := sql.Open("postgres", st.Config.Db_url)

	if err != nil {
		fmt.Println(err)
	}

	dbQueries := database.New(db)

	st.Db = dbQueries

	commands.Register("login", cli.HandlerLogin)
	commands.Register("register", cli.HandlerRegister)
	commands.Register("reset", cli.HandlerReset)
	commands.Register("users", cli.HandlerGetUsers)
	commands.Register("agg", cli.HandlerAgg)
	commands.Register("addfeed", middleware.LoggedInUserMiddleware(cli.HandlerAddFeed))
	commands.Register("feeds", cli.HandlerGetFeeds)
	commands.Register("follow", middleware.LoggedInUserMiddleware(cli.HandlerFollow))
	commands.Register("following", middleware.LoggedInUserMiddleware(cli.HandlerFollowing))
	commands.Register("unfollow", middleware.LoggedInUserMiddleware(cli.HandlerUnfollow))

	args := os.Args[1:]

	command := args[0]
	commandArgs := args[1:]

	err = commands.Run(&st, cli.Command{Name: command, Args: commandArgs})

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
