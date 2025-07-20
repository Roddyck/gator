package utils

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/Roddyck/gator/internal/database"
	"github.com/Roddyck/gator/rss_feed"
	"github.com/Roddyck/gator/state"
	"github.com/google/uuid"
)

func ScrapeFeeds(s *state.State) error {
	feedToFetch, err := s.Db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("error getting next feed to fetch: %w", err)
	}

	params := database.MarkFeedFetchedParams{
		LastFetchedAt: sql.NullTime{Time: time.Now().UTC(), Valid: true},
		UpdatedAt:     time.Now().UTC(),
		ID:            feedToFetch.ID,
	}
	err = s.Db.MarkFeedFetched(context.Background(), params)
	if err != nil {
		return fmt.Errorf("error marking feed fetched: %w", err)
	}

	feed, err := rss_feed.FetchFeed(context.Background(), feedToFetch.Url)
	if err != nil {
		return fmt.Errorf("error fetching feed: %w", err)
	}

	for _, item := range feed.Channel.Item {
		post, err := s.Db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Title:       item.Title,
			Url:         item.Link,
			Description: sql.NullString{String: item.Description, Valid: item.Description != ""},
			PublishedAt: parsePubDate(item.PubDate),
			FeedID:      feedToFetch.ID,
		})
		if err != nil {
			return fmt.Errorf("error creating post: %w", err)
		}

		if post.ID == uuid.Nil && post.Title == "" {
			log.Printf("Duplicate post URL, skipping: %s", item.Link)
		}
		log.Printf("Created post: ID: %s, Title: %s, URL: %s", post.ID, post.Title, post.Url)
	}

	return nil
}

func parsePubDate(pubDate string) time.Time {
	pubDateTime, err := time.Parse(time.RFC1123, pubDate)
	if err != nil {
		log.Printf("error parsing pubDate: %s, error: %v", pubDate, err)
		return time.Now()
	}

	return pubDateTime
}
