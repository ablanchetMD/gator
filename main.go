package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/ablanchetmd/gator/internal/config"
	"github.com/ablanchetmd/gator/internal/database"
	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.ReadConfig()
	if err != nil {
		fmt.Println("Error reading config file")
	}
	s := config.State{}

	s.Config = &cfg

	commands := commands{}
	commands.register("login", handlerLogin)
	commands.register("register", handlerRegister)
	commands.register("reset", handlerReset)
	commands.register("users", handlerListUsers)
	commands.register("agg", handlerAgg)
	commands.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	commands.register("feeds", handlerFeeds)
	commands.register("follow", middlewareLoggedIn(handlerFollow))
	commands.register("following", middlewareLoggedIn(handlerFollowing))
	commands.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	commands.register("browse", middlewareLoggedIn(handlerBrowse))

	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		fmt.Println("Error connecting to database")
	}
	defer db.Close()
	dbQueries := database.New(db)

	s.Db = dbQueries

	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Println("No command provided")
		os.Exit(1)
	}

	err = commands.run(&s, command{Name: args[0], Args: args[1:]})
	if err != nil {
		fmt.Println("Error running command:", err)
		os.Exit(1)
	}

}
