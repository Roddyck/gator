package handlers

import (
	"context"
	"errors"
	"fmt"

	"github.com/Roddyck/gator/state"
	"github.com/Roddyck/gator/command"
)

func HandleLogin(s *state.State, cmd command.Command) error {
	if len(cmd.Args) == 0 {
		return errors.New("you should provide username argument")
	}

	name := cmd.Args[0]

	_, err := s.Db.GetUser(context.Background(), name)
	if err != nil {
		return fmt.Errorf("user doesn't exist: %w", err)
	}

	err = s.Cfg.SetUser(name)
	if err != nil {
		return err
	}

	fmt.Println("user has been set to:", name)

	return nil
}

