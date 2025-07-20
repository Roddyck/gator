package handlers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/Roddyck/gator/command"
	"github.com/Roddyck/gator/state"
	"github.com/Roddyck/gator/utils"
)

func HandleAgg(s *state.State, cmd command.Command) error {
	if len(cmd.Args) != 1 {
		return errors.New("you must provide duration string like 1s, 1m, 1h")
	}

	concurrentScrapers := 3

	timeBetweenRequests, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("error parsing duration string: %w", err)
	}

	log.Printf("Scraping on %v goroutines every %s duration", concurrentScrapers, timeBetweenRequests)
	ticker := time.NewTicker(timeBetweenRequests)

	for ; ; <-ticker.C {
		feedToFetch, err := s.Db.GetNextFeedToFetch(context.Background(), int32(concurrentScrapers))
		if err != nil {
			log.Printf("error getting next feed to fetch: %v", err)
		}

		wg := &sync.WaitGroup{}
		for _, feed := range feedToFetch {
			wg.Add(1)

			go utils.ScrapeFeed(wg, s, feed)
		}
		wg.Wait()
	}
}
