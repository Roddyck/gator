package handlers

import (
	"context"
	"errors"
	"fmt"

	"github.com/Roddyck/gator/command"
	"github.com/Roddyck/gator/state"
	"github.com/Roddyck/gator/internal/database"
)

func HandleUnfollow(s *state.State, cmd command.Command, user database.User) error {
	if len(cmd.Args) != 1 {
		return errors.New("you must provide feed url as argument")
	}

	feedURL := cmd.Args[0]

	feed, err := s.Db.GetFeedByUrl(context.Background(), feedURL)
	if err != nil {
		return fmt.Errorf("error retrieving feed with given url: %w", err)
	}

	deleteParams := database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}

	err = s.Db.DeleteFeedFollow(context.Background(), deleteParams)
	if err != nil {
		return fmt.Errorf("error unfollowing feed: %w", err)
	}

	return nil
}
