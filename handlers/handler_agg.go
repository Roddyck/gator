package handlers

import (
	"errors"
	"fmt"
	"time"

	"github.com/Roddyck/gator/command"
	"github.com/Roddyck/gator/state"
	"github.com/Roddyck/gator/utils"
)

func HandleAgg(s *state.State, cmd command.Command) error {
	if len(cmd.Args) != 1 {
		return errors.New("you must provide duration string like 1s, 1m, 1h")
	}

	timeBetweenRequests, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("error parsing duration string: %w", err)
	}

	ticker := time.NewTicker(timeBetweenRequests)

	for ; ; <-ticker.C {
		utils.ScrapeFeeds(s)
	}
}
