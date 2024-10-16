package main

import (
	"fmt"

	"github.com/acehotel33/bootdev-gator/internal/config"
)

const user = "vakho"

func main() {
	cfg, err := config.Read()
	if err != nil {
		panic(err)
	}

	if err := cfg.SetUser(user); err != nil {
		panic(err)
	}

	cfg, err = config.Read()
	if err != nil {
		panic(err)
	}
	fmt.Println(cfg)
}
