package handlers

import (
	"context"
	"fmt"

	"github.com/Roddyck/gator/internal/database"
	"github.com/Roddyck/gator/state"
	"github.com/Roddyck/gator/command"
)

func HandleFollowing(s *state.State, cmd command.Command, user database.User) error {
	feedFollows, err := s.Db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("error getting feed follows for current user: %w", err)
	}

	for _, feedFollow := range feedFollows {
		fmt.Println("Feed name:", feedFollow.FeedName)
	}

	return nil
}
