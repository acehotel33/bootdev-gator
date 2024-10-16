package main

import (
	"fmt"
	"log"

	"github.com/acehotel33/bootdev-gator/internal/commands"
	"github.com/acehotel33/bootdev-gator/internal/config"
	"github.com/acehotel33/bootdev-gator/internal/state"
)

const user = "vakho"

func main() {
	// Initialize config
	cfg, err := config.InitializeConfig(user)
	if err != nil {
		log.Fatalf("Failed to initialize config: %v", err)
	}

	// Initialize State
	state, err := state.InitializeState(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize state: %v", err)
	}

	// Initialize commands
	_, err = commands.InitializeCommands()
	if err != nil {
		log.Fatalf("Failed to initialize commands: %v", err)
	}

	fmt.Println(state.Cfg)
}
