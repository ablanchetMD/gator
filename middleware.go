package main

import (
	"context"
	"errors"

	"github.com/ablanchetmd/gator/internal/config"
	"github.com/ablanchetmd/gator/internal/database"
	"github.com/google/uuid"
)



func middlewareLoggedIn(handler func(s *config.State, cmd command, user database.User) error) func(*config.State, command) error {
	return func(s *config.State, cmd command) error {
		if s.Config.CurrentUserName == "" {
			return errors.New("not logged in")
		}
		user, err := s.Db.GetUser(context.Background(), database.GetUserParams{
			Name: s.Config.CurrentUserName,
			ID:   uuid.Nil,
		})
		if err != nil {		
			return err
		}
		
		return handler(s, cmd, user)
	}
}