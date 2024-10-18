package commands

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/acehotel33/bootdev-gator/internal/database"
	"github.com/acehotel33/bootdev-gator/internal/rss"
	"github.com/acehotel33/bootdev-gator/internal/state"
	"github.com/google/uuid"
)

type command struct {
	name      string
	arguments []string
}

type Commands struct {
	commandsMap map[string]func(*state.State, command) error
}

func InitializeCommands() (*Commands, error) {
	cmds := &Commands{
		commandsMap: map[string]func(*state.State, command) error{},
	}
	cmds.Register("login", HandlerLogin)
	cmds.Register("register", HandlerRegister)
	cmds.Register("reset", HandlerReset)
	cmds.Register("users", HandlerGetAllUsers)
	cmds.Register("agg", HandlerAggregator)
	cmds.Register("addfeed", HandlerAddFeed)
	return cmds, nil
}

func RunCommand(state *state.State, cmds *Commands) error {
	if len(os.Args) < 2 {
		return errors.New("not enough arguments")
	}

	commandName := os.Args[1]
	commandArgs := os.Args[2:]

	cmd := command{
		name:      commandName,
		arguments: commandArgs,
	}

	if err := cmds.Run(state, cmd); err != nil {
		return err
	}

	return nil
}

func (c *Commands) Register(name string, f func(*state.State, command) error) {
	c.commandsMap[name] = f
}

func (c *Commands) Run(s *state.State, cmd command) error {
	if f, ok := c.commandsMap[cmd.name]; ok {
		if err := f(s, cmd); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("invalid command")
}

func HandlerRegister(s *state.State, cmd command) error {
	if len(cmd.arguments) != 1 {
		return errors.New("invalid arguments")
	}

	username := cmd.arguments[0]
	createUserParams := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      username,
	}

	user, err := s.DB.CreateUser(context.Background(), createUserParams)
	if err != nil {
		return errors.New("could not register user")
	}

	s.Cfg.SetUser(username)
	fmt.Printf("user %v created\n", username)
	fmt.Println(user)

	return nil
}

func HandlerLogin(s *state.State, cmd command) error {
	if len(cmd.arguments) == 0 || len(cmd.arguments) > 1 {
		return errors.New("invalid arguments")
	}

	username := cmd.arguments[0]
	_, err := s.DB.GetUser(context.Background(), username)
	if err != nil {
		return errors.New("could not log in")
	}
	if err := s.Cfg.SetUser(username); err != nil {
		return err
	}
	fmt.Printf("user has been set to: %v\n", username)
	return nil
}

func HandlerReset(s *state.State, cmd command) error {
	if len(cmd.arguments) != 0 {
		return errors.New("invalid arguments")
	}

	if err := s.DB.ResetUsers(context.Background()); err != nil {
		return errors.New("could not reset users")
	}

	fmt.Println("users reset")
	return nil
}

func HandlerGetAllUsers(s *state.State, cmd command) error {
	users, err := s.DB.GetAllUsers(context.Background())
	if err != nil {
		return errors.New("could not get users")
	}

	for _, user := range users {
		line := "* " + user.Name
		if user.Name == s.Cfg.CurrentUsername {
			line = line + " (current)"
		}
		fmt.Println(line)
	}

	return nil
}

func HandlerAggregator(s *state.State, cmd command) error {
	if len(cmd.arguments) != 0 {
		return errors.New("invalid arguments")
	}

	rssFeed, err := rss.FetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}

	fmt.Println(rssFeed)
	return nil
}

func HandlerAddFeed(s *state.State, cmd command) error {
	if len(cmd.arguments) != 2 {
		return errors.New("invalid arguments")
	}
	feedName := cmd.arguments[0]
	feedURL := cmd.arguments[1]

	currentUsername := s.Cfg.CurrentUsername
	currentUser, err := s.DB.GetUser(context.Background(), currentUsername)
	if err != nil {
		return err
	}
	currentUserID := currentUser.ID

	_, err = rss.FetchFeed(context.Background(), feedURL)
	if err != nil {
		return err
	}

	createFeedParams := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      feedName,
		Url:       feedURL,
		UserID: uuid.NullUUID{
			UUID:  currentUserID,
			Valid: true,
		},
	}

	dbFeed, err := s.DB.CreateFeed(context.Background(), createFeedParams)
	if err != nil {
		return err
	}

	fmt.Println(dbFeed)
	return nil
}
