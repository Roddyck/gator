package handlers

import (
	"context"
	"fmt"

	"github.com/Roddyck/gator/command"
	"github.com/Roddyck/gator/state"
)

func HandleReset(s *state.State, cmd command.Command) error {
	err := s.Db.DeleteAllUsers(context.Background())
	if err != nil {
		return fmt.Errorf("error deleting users: %w", err)
	}

	return nil
}
