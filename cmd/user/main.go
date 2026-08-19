package main

import (
	"log"
	"os"
	"path"

	"github.com/abozorov/school_online/cmd/user/internal/app"
	"github.com/abozorov/school_online/cmd/user/internal/config"
)

func main() {
	dir, err := os.Getwd()
	if err != nil {
		log.Println("Error in path determination: ", err)
		return
	}
	cfg, err := config.NewConfig(path.Join(dir, "cmd", "user", "internal", "config", "user_config.env"))

	if err != nil {
		log.Println("Error when load config file: ", err)
		return
	}

	app.Run(cfg)
}
