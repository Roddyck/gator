package utils

import (
	"context"
	"database/sql"
	"log"
	"sync"
	"time"

	"github.com/Roddyck/gator/internal/database"
	"github.com/Roddyck/gator/rss_feed"
	"github.com/Roddyck/gator/state"
	"github.com/google/uuid"
)


func ScrapeFeed(wg *sync.WaitGroup, s *state.State, feed database.Feed) {
	defer wg.Done()

	params := database.MarkFeedFetchedParams{
		LastFetchedAt: sql.NullTime{Time: time.Now().UTC(), Valid: true},
		UpdatedAt:     time.Now().UTC(),
		ID:            feed.ID,
	}
	err := s.Db.MarkFeedFetched(context.Background(), params)
	if err != nil {
		log.Printf("error marking feed fetched: %v", err)
		return
	}

	rssFeed, err := rss_feed.FetchFeed(context.Background(), feed.Url)
	if err != nil {
		log.Printf("error fetching feed: %v", err)
		return
	}

	for _, item := range rssFeed.Channel.Item {
		post, err := s.Db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Title:       item.Title,
			Url:         item.Link,
			Description: sql.NullString{String: item.Description, Valid: item.Description != ""},
			PublishedAt: parsePubDate(item.PubDate),
			FeedID:      feed.ID,
		})
		if err != nil {
			log.Printf("error creating post: %v", err)
		}

		if post.ID == uuid.Nil && post.Title == "" {
			log.Printf("Duplicate post URL, skipping: %s", item.Link)
		}
		log.Printf("Created post: ID: %s, Title: %s, URL: %s", post.ID, post.Title, post.Url)
	}

}

func parsePubDate(pubDate string) time.Time {
	pubDateTime, err := time.Parse(time.RFC1123, pubDate)
	if err != nil {
		log.Printf("error parsing pubDate: %s, error: %v", pubDate, err)
		return time.Now()
	}

	return pubDateTime
}
