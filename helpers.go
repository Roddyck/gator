package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Roddyck/gator/internal/database"
)

func screpeFeeds(s *state) error {
	feedToFetch, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("error getting next feed to fetch: %w", err)
	}

	params := database.MarkFeedFetchedParams{
		LastFetchedAt: sql.NullTime{Time: time.Now(), Valid: true},
		UpdatedAt:     time.Now(),
		ID:            feedToFetch.ID,
	}
	err = s.db.MarkFeedFetched(context.Background(), params)
	if err != nil {
		return fmt.Errorf("error marking feed fetched: %w", err)
	}

	feed, err := fetchFeed(context.Background(), feedToFetch.Url)
	if err != nil {
		return fmt.Errorf("error fetching feed: %w", err)
	}

	for _, item := range feed.Channel.Item {
		fmt.Println(" - ", item.Title)
	}

	return nil
}
