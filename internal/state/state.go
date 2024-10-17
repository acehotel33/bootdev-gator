package state

import (
	"github.com/acehotel33/bootdev-gator/internal/config"
	"github.com/acehotel33/bootdev-gator/internal/database"
)

type State struct {
	Cfg *config.Config
	DB  *database.Queries
}

func InitializeState(cfg *config.Config) (*State, error) {
	return &State{
		Cfg: cfg,
	}, nil
}
