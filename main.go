package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/acehotel33/bootdev-gator/internal/commands"
	"github.com/acehotel33/bootdev-gator/internal/config"
	"github.com/acehotel33/bootdev-gator/internal/database"
	"github.com/acehotel33/bootdev-gator/internal/state"
	_ "github.com/lib/pq"
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

	// Connect to database and store in state
	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		log.Fatalf("Failed to connect to database")
	}
	dbQueries := database.New(db)
	state.DB = dbQueries

	// Initialize commands
	cmds, err := commands.InitializeCommands()
	if err != nil {
		log.Fatalf("Failed to initialize commands: %v", err)
	}

	// Run commands
	if err := commands.RunCommand(state, cmds); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	os.Exit(0)
}
