package state

import (
	"github.com/Roddyck/gator/internal/config"
	"github.com/Roddyck/gator/internal/database"
)

type State struct {
	Db *database.Queries
	Cfg *config.Config
}
