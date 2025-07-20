package handlers

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Roddyck/gator/command"
	"github.com/Roddyck/gator/internal/database"
	"github.com/Roddyck/gator/state"
	"github.com/google/uuid"
)

func HandleAddFeed(s *state.State, cmd command.Command, user database.User) error {
	if len(cmd.Args) != 2 {
		return errors.New("you should provide name and url of the feed")
	}

	name, url := cmd.Args[0], cmd.Args[1]

	feedParams := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      name,
		Url:       url,
		UserID:    user.ID,
	}

	DbFeed, err := s.Db.CreateFeed(context.Background(), feedParams)
	if err != nil {
		return err
	}

	feedFollowParams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feedParams.ID,
	}

	_, err = s.Db.CreateFeedFollow(context.Background(), feedFollowParams)
	if err != nil {
		return err
	}

	fmt.Printf("Feed record: %+v", DbFeed)

	return nil
}
