package main

import (
	"errors"
	"fmt"
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

	err := s.cfg.SetUser(cmd.args[0])
	if err != nil {
	    return err
	}

	fmt.Println("user has been set to:", cmd.args[0])

	return nil
}
