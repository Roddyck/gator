package handlers

import (
	"context"
	"fmt"

	"github.com/Roddyck/gator/command"
	"github.com/Roddyck/gator/state"
)

func HandleFeeds(s *state.State, cmd command.Command) error {
	feeds, err := s.Db.ListFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("error retrieving feeds from Db: %w", err)
	}

	for _, feed := range feeds {
		fmt.Println("Name:", feed.Name)
		fmt.Println("URL:", feed.Url)

		user, err := s.Db.GetUserById(context.Background(), feed.UserID)
		if err != nil {
			return fmt.Errorf("error retrieving user that created the feed: %w", err)
		}

		fmt.Println("Username:", user.Name)
	}

	return nil
}
