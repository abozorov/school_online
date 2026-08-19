package main

import (
	"log"
	"os"
	"path"

	"github.com/abozorov/school_online/cmd/classroom/internal/app"
	"github.com/abozorov/school_online/cmd/classroom/internal/config"
)

func main() {
	dir, err := os.Getwd()
	if err != nil {
		log.Println("Error in path determination: ", err)
		return
	}
	cfg, err := config.NewConfig(path.Join(dir, "cmd", "classroom", "internal", "config", "classroom_config.env"))

	if err != nil {
		log.Println("Error when load config file: ", err)
		return
	}

	app.Run(cfg)
}
