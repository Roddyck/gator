package handlers

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Roddyck/gator/internal/database"
	"github.com/Roddyck/gator/state"
	"github.com/Roddyck/gator/command"
	"github.com/google/uuid"
)

func HandleRegister(s *state.State, cmd command.Command) error {
	if len(cmd.Args) == 0 {
		return errors.New("you should provide username argument")
	}

	id := uuid.New()
	name := cmd.Args[0]
	params := database.CreateUserParams{
		ID:        id,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      name,
	}

	user, err := s.Db.CreateUser(context.Background(), params)
	if err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}

	err = s.Cfg.SetUser(name)
	if err != nil {
		return err
	}

	fmt.Printf("user was created\n%+v", user)

	return nil
}
