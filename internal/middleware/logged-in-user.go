package middleware

import (
	"context"
	cli "gator/internal/cli"
	"gator/internal/database"
	"gator/internal/state"
)

func LoggedInUserMiddleware(handler func(s *state.State, cmd cli.Command, user database.User) error) func(*state.State, cli.Command) error {
	return func(s *state.State, cmd cli.Command) error {
		user, err := s.Db.GetUser(context.Background(), s.Config.Current_user_name)
		if err != nil {
			return err
		}
		return handler(s, cmd, user)
	}

}
