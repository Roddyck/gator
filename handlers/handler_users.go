package handlers

import (
	"context"
	"fmt"

	"github.com/Roddyck/gator/command"
	"github.com/Roddyck/gator/state"
)

func HandleUsers(s *state.State, cmd command.Command) error {
	users, err := s.Db.ListUsers(context.Background())
	if err != nil {
		return fmt.Errorf("error retrieving users: %w", err)
	}

	for _, user := range users {
		if s.Cfg.CurrentUserName == user.Name {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Printf("* %s\n", user.Name)
		}
	}

	return nil
}
