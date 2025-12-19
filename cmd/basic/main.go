package main

import (
	"fmt"

	"github.com/adreasnow/auto-pull/internal/pkg/config"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(cfg.Directories)

}
