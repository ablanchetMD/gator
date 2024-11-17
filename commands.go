package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/ablanchetmd/gator/internal/config"
	"github.com/ablanchetmd/gator/internal/database"
	"github.com/google/uuid"
)

type command struct {
	Name string
	Args []string
}

type commands struct {
	library map[string]func(*config.State, command) error
}

func (c *commands) register(name string, f func(s *config.State, c command) error) error {
	if c.library == nil {
		c.library = make(map[string]func(*config.State, command) error) // Initialize the map if not already done
	}
	if _, exists := c.library[name]; exists {
		return errors.New("command already registered")
	}
	c.library[name] = f
	return nil
}

func (c *commands) run(s *config.State, cmd command) error {
	f, ok := c.library[cmd.Name]
	if !ok {
		return errors.New("command not found")
	}
	return f(s, cmd)
}

func handlerRegister(s *config.State, c command) error {
	if len(c.Args) != 1 {
		return errors.New("register command requires 1 argument")
	}
	user, err := s.Db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      c.Args[0],
	})
	if err != nil {
		fmt.Println("Error creating user", err.Error())
		os.Exit(1)
		return err
	}
	fmt.Println("User created with ID", user.ID)
	err = s.Config.SetUser(user.Name)
	if err != nil {
		return err
	}

	return nil
}

func handlerAddFeed(s *config.State, c command, user database.User) error {
	if len(c.Args) != 2 {
		return errors.New("AddFeed command requires 2 arguments (name, url)")
	}
	
	feed, err := s.Db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      c.Args[0],
		Url:       c.Args[1],
		UserID:    user.ID,

	})
	if err != nil {		
		return err
	}

	follow_feed, err := s.Db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,		

	})
	if err != nil {		
		return err
	}
	fmt.Println("Feed created with ID", feed)	
	fmt.Printf("%s is now following %s.", follow_feed.UserName, follow_feed.FeedName)		

	return nil
}

func handlerFollow(s *config.State, c command,user database.User) error {
	if len(c.Args) != 1 {
		return errors.New("follow command requires 1 arguments (url)")
	}	

	feed, err := s.Db.GetFeedByUrl(context.Background(), c.Args[0])
	
	if err != nil {		
		return err
	}

	follow_feed, err := s.Db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,		

	})
	if err != nil {		
		return err
	}
	fmt.Printf("%s is now following %s.", follow_feed.UserName, follow_feed.FeedName)	

	return nil
}

func handlerBrowse(s *config.State, c command, user database.User) error {
	if len(c.Args) > 1 {
		return errors.New("browse command requires 1 argument at most (number of posts to return)")
	}
	ctx := context.Background()
	limit := 2
	if len(c.Args) != 0 {
		new_limit, err := strconv.ParseInt(c.Args[0],10,32)
		if err != nil {			
			fmt.Printf("Error parsing limit: %v, leaving limit to default 2", err)
		}
		limit = int(new_limit)
	}
	
	posts, err := s.Db.GetPostsForUser(ctx, database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		fmt.Println("Error fetching posts", err.Error())
		os.Exit(1)
		return err
	}
	
	for _, post := range posts {
		fmt.Printf("*****\n TITLE: %s\n URL: %s\n DESCRIPTION: %s\n PUBLISHED AT: %s\n FROM FEED: %s\n*****\n", post.Title.String, post.Url, post.Description.String, post.PublishedAt, post.FeedName)
	}
	
	return nil
}

func handlerFeeds(s *config.State, c command) error {
	if len(c.Args) != 0 {
		return errors.New("feeds command requires no argument")
	}
	ctx := context.Background()
	
	feeds, err := s.Db.GetFeeds(ctx)
	if err != nil {
		fmt.Println("Error fetching feeds", err.Error())		
		return err
	}
	
	for _, feed := range feeds {
		user, err := s.Db.GetUser(context.Background(), database.GetUserParams{
			Name: "",
			ID:   feed.UserID,
		})
		if err != nil {
			return err
		}
		
		fmt.Printf("FEED:\n NAME: %s\n URL: %s\n CREATED BY: %s\n", feed.Name, feed.Url, user.Name)
		
		
	}
	
	return nil
}

func handlerFollowing(s *config.State, c command, user database.User) error {
	if len(c.Args) != 0 {
		return errors.New("following command requires no argument")
	}
	ctx := context.Background()	
	
	feeds, err := s.Db.GetFeedFollowsForUser(ctx, user.ID)
	if err != nil {
		fmt.Println("Error fetching followed feeds", err.Error())		
		return err
	}
	fmt.Printf("%s is currently following:\n", user.Name)
	if len(feeds) == 0 {
		fmt.Println("(empty)")
		return nil
	}
	for _, feed := range feeds {
		
		fmt.Printf("- %s\n", feed.FeedName)	
		
	}
	
	return nil
}

func handlerUnfollow(s *config.State, c command, user database.User) error {
	if len(c.Args) != 1 {
		return errors.New("unfollow command requires an url argument")
	}
	ctx := context.Background()	
	
	feed, err := s.Db.GetFeedByUrl(ctx, c.Args[0])
	if err != nil {
		return err
	}

	err = s.Db.UnfollowFeed(ctx, database.UnfollowFeedParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})

	if err != nil {
		return err
	}
	fmt.Printf("%s is no longer following %s.", user.Name, feed.Name)
	
	return nil
}



func handlerReset(s *config.State, c command) error {
	if len(c.Args) != 0 {
		return errors.New("reset command requires no argument")
	}
	ctx := context.Background()
	
	err := s.Db.DeleteUsers(ctx)
	if err != nil {
		fmt.Println("Error reseting users table", err.Error())
		os.Exit(1)
		return err
	}
	
	fmt.Println("User table reset")
	return nil
}

func handlerListUsers(s *config.State, c command) error {
	if len(c.Args) != 0 {
		return errors.New("list command requires no argument")
	}
	ctx := context.Background()
	
	users, err := s.Db.GetUsers(ctx)
	if err != nil {
		fmt.Println("Error fetching users", err.Error())
		os.Exit(1)
		return err
	}
	
	for _, user := range users {
		if user.Name == s.Config.CurrentUserName {
			fmt.Printf("* %s (current)\n", user.Name)
		}else{
			fmt.Printf("* %s\n", user.Name)
		}
		
	}
	
	return nil
}

func handlerLogin(s *config.State, c command) error {
	if len(c.Args) != 1 {
		return errors.New("login command requires 1 argument")
	}
	ctx := context.Background()
	// get the author we just inserted
	fetchedUser, err := s.Db.GetUser(ctx,database.GetUserParams{
		Name: c.Args[0],
		ID:   uuid.Nil,
	})
	if err != nil {
		fmt.Println("Error fetching user", err.Error())
		os.Exit(1)
		return err
	}

	err = s.Config.SetUser(fetchedUser.Name)
	if err != nil {
		return err
	}

	fmt.Println("User set to", c.Args[0])
	return nil
}

func handlerAgg(s *config.State, c command) error {
	if len(c.Args) != 1 {
		return errors.New("agg command requires 1 argument (time string)")
	}
	
	timeBetweenReqs, err := time.ParseDuration(c.Args[0])
	if err != nil {
		fmt.Println("Error parsing time duration", err.Error())
		return err		
	}

	ticker := time.NewTicker(timeBetweenReqs)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}	
}
