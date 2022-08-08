package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/dwilkolek/moxy/config"
	"github.com/dwilkolek/moxy/internal/app"
)

func main() {
	// Configuration
	cmd := ""
	if len(os.Args) >= 2 {
		cmd = os.Args[1]
	}
	switch cmd {
	case "start":
		file := defaultConfigFile()
		if len(os.Args) >= 3 {
			file = os.Args[2]
		}
		cfg, err := config.NewConfig(file)
		if err != nil {
			log.Fatalf("Config error: %s", err)
			panic("Failed to start application")
		}
		app.Run(cfg)
	case "update":
		app.Update()
	default:
		fmt.Println("Available options:\n \t start [config_file:config.json]/[profile] - to start application with [config_file] or config-[profile].json\n \t update - to upadate application")
	}

}

func defaultConfigFile() string {
	cwd, _ := os.Executable()
	return filepath.Join(filepath.Dir(cwd), "config.json")
}
