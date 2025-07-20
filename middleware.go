package main

import (
	"context"

	"github.com/Roddyck/gator/internal/database"
	"github.com/Roddyck/gator/state"
	"github.com/Roddyck/gator/command"
)

func middlewareLoggedIn(handler func(s *state.State, cmd command.Command, user database.User) error) func(*state.State, command.Command) error {
	return func(s *state.State, cmd command.Command) error {
		user, err := s.Db.GetUser(context.Background(), s.Cfg.CurrentUserName)
		if err != nil {
		    return err
		}

		return handler(s, cmd, user)
	}
}
