package main

import (
	"github.com/Roddyck/gator/internal/config"
	"github.com/Roddyck/gator/internal/database"
)

type state struct {
	db *database.Queries
	cfg *config.Config
}
