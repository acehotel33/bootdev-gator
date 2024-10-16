package main

import (
	"fmt"
	"log"

	"github.com/acehotel33/bootdev-gator/internal/commands"
	"github.com/acehotel33/bootdev-gator/internal/config"
	"github.com/acehotel33/bootdev-gator/internal/state"
)

func main() {
	// Initialize config
	cfg, err := config.InitializeConfig()
	if err != nil {
		log.Fatalf("Failed to initialize config: %v", err)
	}

	// Initialize State
	state, err := state.InitializeState(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize state: %v", err)
	}

	// Initialize commands
	cmds, err := commands.InitializeCommands()
	if err != nil {
		log.Fatalf("Failed to initialize commands: %v", err)
	}

	// Run commands
	commands.RunCommand(state, cmds)

	fmt.Println(state.Cfg)

}
