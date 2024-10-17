package commands

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/acehotel33/bootdev-gator/internal/database"
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
	return cmds, nil
}

func RunCommand(state *state.State, cmds *Commands) {
	if len(os.Args) < 2 {
		fmt.Println("Error: Not enough arguments")
		os.Exit(1)
	}

	commandName := os.Args[1]
	commandArgs := os.Args[2:]

	cmd := command{
		name:      commandName,
		arguments: commandArgs,
	}

	cmds.Run(state, cmd)
}

func (c *Commands) Register(name string, f func(*state.State, command) error) {
	c.commandsMap[name] = f
}

func (c *Commands) Run(s *state.State, cmd command) error {
	if f, ok := c.commandsMap[cmd.name]; ok {
		f(s, cmd)
		return nil
	}
	return fmt.Errorf("invalid command")
}

func HandlerRegister(s *state.State, cmd command) error {
	if len(cmd.arguments) != 1 {
		log.Fatalf("invalid arguments")
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
		log.Fatalf("could not register user")
	}

	s.Cfg.SetUser(username)
	log.Printf("user %v created", username)
	log.Println(user)

	return nil
}

func HandlerLogin(s *state.State, cmd command) error {
	if len(cmd.arguments) == 0 || len(cmd.arguments) > 1 {
		fmt.Println("invalid arguments")
		os.Exit(1)
	}
	username := cmd.arguments[0]
	_, err := s.DB.GetUser(context.Background(), username)
	if err != nil {
		log.Fatalf("could not log in")
	}
	if err := s.Cfg.SetUser(username); err != nil {
		return err
	}
	fmt.Printf("\nUser has been set to: %v\n", username)
	return nil
}
