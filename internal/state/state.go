package state

import "github.com/acehotel33/bootdev-gator/internal/config"

type State struct {
	Cfg *config.Config
}

func InitializeState(cfg *config.Config) (*State, error) {
	return &State{
		Cfg: cfg,
	}, nil
}
