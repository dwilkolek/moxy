package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/dwilkolek/moxy/config"
	"github.com/dwilkolek/moxy/internal/app"
	"github.com/dwilkolek/moxy/internal/logger"
)

func main() {
	// Configuration
	cmd := ""
	if len(os.Args) >= 2 {
		cmd = os.Args[1]
	}
	logger.New("Moxy").Printf("Version %s", app.Version())
	var wg sync.WaitGroup
	if cmd != "update" {
		go func() {
			wg.Add(1)
			isUrlUpdate, _ := app.CheckForUpdateUrl()
			if isUrlUpdate != "" {
				fmt.Println("There is new version available. run program with `update` argument to get it")
			} else {
				fmt.Println("Up to date")
			}
			defer wg.Done()
		}()
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
		wg.Add(1)
		app.Update()
		wg.Done()
	default:
		fmt.Println("Available options:\n \t start [config_file:config.json]/[profile] - to start application with [config_file] or config-[profile].json\n \t update - to upadate application")
	}
	wg.Wait()

}

func defaultConfigFile() string {
	cwd, _ := os.Executable()
	return filepath.Join(filepath.Dir(cwd), "config.json")
}
