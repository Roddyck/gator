package handlers

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Roddyck/gator/command"
	"github.com/Roddyck/gator/internal/database"
	"github.com/Roddyck/gator/state"
)

func HandleBrowse(s *state.State, cmd command.Command, user database.User) error {
	limit := 2
	if len(cmd.Args) == 1 {
		if limitNum, err := strconv.Atoi(cmd.Args[0]); err == nil {
			limit = limitNum
		} else {
			return fmt.Errorf("invalid limit: %s", cmd.Args[0])
		}
	}

	posts, err := s.Db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		ID:    user.ID,
		Limit: int32(limit),
	})
	if err != nil {
		return fmt.Errorf("error retrieving posts: %w", err)
	}

	for _, post := range posts {
		fmt.Printf(" ---\"%s\" (%s)---\n", post.Title, post.FeedName)
		fmt.Printf("\t%s\n", post.Description.String)
		fmt.Printf(" Link: %s\n", post.Url)
		fmt.Println("---------------------------------")
	}

	return nil
}
