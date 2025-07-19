package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Roddyck/gator/internal/database"
	"github.com/google/uuid"
)

type command struct {
	name string
	args []string
}

type commands struct {
	handlers map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	handler, ok := c.handlers[cmd.name]
	if !ok {
		return errors.New("given commands does not exist")
	}

	err := handler(s, cmd)
	if err != nil {
		return err
	}

	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.handlers[name] = f
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("you should provide username argument")
	}

	name := cmd.args[0]

	_, err := s.db.GetUser(context.Background(), name)
	if err != nil {
		return fmt.Errorf("user doesn't exist: %w", err)
	}

	err = s.cfg.SetUser(name)
	if err != nil {
		return err
	}

	fmt.Println("user has been set to:", name)

	return nil
}

func handleRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("you should provide username argument")
	}

	id := uuid.New()
	name := cmd.args[0]
	params := database.CreateUserParams{
		ID:        id,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
	}

	user, err := s.db.CreateUser(context.Background(), params)
	if err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}

	err = s.cfg.SetUser(name)
	if err != nil {
		return err
	}

	fmt.Printf("user was created\n%+v", user)

	return nil
}

func handleReset(s *state, cmd command) error {
	err := s.db.DeleteAllUsers(context.Background())
	if err != nil {
		return fmt.Errorf("error deleting users: %w", err)
	}

	return nil
}

func handleUsers(s *state, cmd command) error {
	users, err := s.db.ListUsers(context.Background())
	if err != nil {
		return fmt.Errorf("error retrieving users: %w", err)
	}

	for _, user := range users {
		if s.cfg.CurrentUserName == user.Name {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Printf("* %s\n", user.Name)
		}
	}

	return nil
}

func handleAgg(s *state, cmd command) error {
	url := "https://www.wagslane.dev/index.xml"

	feed, err := fetchFeed(context.Background(), url)
	if err != nil {
		return fmt.Errorf("error fetching feed: %w", err)
	}

	fmt.Printf("Feed: %+v\n", feed)

	return nil
}

func handleAddFeed(s *state, cmd command) error {
	if len(cmd.args) != 2 {
		return errors.New("you should provide name and url of the feed")
	}

	name, url := cmd.args[0], cmd.args[1]

	username := s.cfg.CurrentUserName
	currentUser, err := s.db.GetUser(context.Background(), username)
	if err != nil {
		return fmt.Errorf("error retrieving current user: %w", err)
	}

	feedParams := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
		Url:       url,
		UserID:    currentUser.ID,
	}

	dbFeed, err := s.db.CreateFeed(context.Background(), feedParams)
	if err != nil {
		return err
	}

	feedFollowParams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    currentUser.ID,
		FeedID:    feedParams.ID,
	}

	_, err = s.db.CreateFeedFollow(context.Background(), feedFollowParams)
	if err != nil {
		return err
	}

	fmt.Printf("Feed record: %+v", dbFeed)

	return nil
}

func handleFeeds(s *state, cmd command) error {
	feeds, err := s.db.ListFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("error retrieving feeds from db: %w", err)
	}

	for _, feed := range feeds {
		fmt.Println("Name:", feed.Name)
		fmt.Println("URL:", feed.Url)

		user, err := s.db.GetUserById(context.Background(), feed.UserID)
		if err != nil {
			return fmt.Errorf("error retrieving user that created the feed: %w", err)
		}

		fmt.Println("Username:", user.Name)
	}

	return nil
}

func handleFollow(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return errors.New("you must provide url as argument")
	}

	url := cmd.args[0]

	feed, err := s.db.GetFeedByUrl(context.Background(), url)
	if err != nil {
		return err
	}

	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return err
	}

	params := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}

	feedFollow, err := s.db.CreateFeedFollow(context.Background(), params)
	if err != nil {
		return err
	}

	fmt.Println("Name of the feed:", feedFollow.FeedName)
	fmt.Println("Name of current user:", feedFollow.UserName)

	return nil
}

func handleFollowing(s *state, cmd command) error {
	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("error retrieving current user: %w", err)
	}

	feedFollows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("error getting feed follows for current user: %w", err)
	}

	for _, feedFollow := range feedFollows {
		fmt.Println("Feed name:", feedFollow.FeedName)
	}

	return nil
}
