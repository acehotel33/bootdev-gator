package commands

import (
	"fmt"
	"os"

	"github.com/acehotel33/bootdev-gator/internal/state"
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

func HandlerLogin(s *state.State, cmd command) error {
	if len(cmd.arguments) == 0 || len(cmd.arguments) > 1 {
		fmt.Println("invalid arguments")
		os.Exit(1)
	}
	username := cmd.arguments[0]
	if err := s.Cfg.SetUser(username); err != nil {
		return err
	}
	fmt.Printf("\nUser has been set to: %v\n", username)
	return nil
}
