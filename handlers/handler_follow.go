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

func HandleFollow(s *state.State, cmd command.Command, user database.User) error {
	if len(cmd.Args) != 1 {
		return errors.New("you must provide url as argument")
	}

	url := cmd.Args[0]

	feed, err := s.Db.GetFeedByUrl(context.Background(), url)
	if err != nil {
		return err
	}

	params := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}

	feedFollow, err := s.Db.CreateFeedFollow(context.Background(), params)
	if err != nil {
		return err
	}

	fmt.Println("Name of the feed:", feedFollow.FeedName)
	fmt.Println("Name of current user:", feedFollow.UserName)

	return nil
}
