package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/Roddyck/gator/internal/config"
	"github.com/Roddyck/gator/internal/database"
	"github.com/Roddyck/gator/state"
	"github.com/Roddyck/gator/command"
	"github.com/Roddyck/gator/handlers"
	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("postgres", cfg.DbUrl)
	dbQueries := database.New(db)

	s := &state.State{
		Db: dbQueries,
		Cfg: &cfg,
	}

	cmds := command.Commands{
		Handlers: make(map[string]func(*state.State, command.Command) error),
	}

	cmds.Register("login", handlers.HandleLogin)
	cmds.Register("register", handlers.HandleRegister)
	cmds.Register("reset", handlers.HandleReset)
	cmds.Register("users", handlers.HandleUsers)
	cmds.Register("agg", handlers.HandleAgg)
	cmds.Register("addfeed", middlewareLoggedIn(handlers.HandleAddFeed))
	cmds.Register("feeds", handlers.HandleFeeds)
	cmds.Register("follow", middlewareLoggedIn(handlers.HandleFollow))
	cmds.Register("following", middlewareLoggedIn(handlers.HandleFollowing))
	cmds.Register("unfollow", middlewareLoggedIn(handlers.HandleUnfollow))

	args := os.Args
	if len(args) < 2 {
		log.Fatal("not enough arguments provided")
	}

	cmd := command.Command{
		Name: args[1],
		Args: args[2:],
	}

	err = cmds.Run(s, cmd)
	if err != nil {
		log.Fatal(err)
	}
}
