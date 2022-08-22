package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/dwilkolek/moxy/internal/app"
	"github.com/dwilkolek/moxy/internal/logger"
)

func main() {
	// Configuration
	cmd := ""
	if len(os.Args) >= 2 {
		cmd = os.Args[1]
	}
	appLog := logger.New("Moxy")
	if cmd != "version" {
		appLog.Printf("Version %s", app.Version())
	}
	var wg sync.WaitGroup
	if cmd != "update" && cmd != "version" {
		go func() {
			wg.Add(1)
			isUrlUpdate, _ := app.CheckForUpdateUrl()
			if isUrlUpdate != "" {
				appLog.Println("There is new version available. run program with `update` argument to get it")
			} else {
				appLog.Println("Up to date")
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
		app.Run(file)
	case "update":
		wg.Add(1)
		app.Update()
		wg.Done()
	case "version":
		appLog.Printf("Version %s", app.Version())
	default:
		fmt.Println("Available options:\n \t start [config_file:config.json]/[profile] - to start application with [config_file] or config-[profile].json\n \t update - to upadate application\n \t version - to get version information")
	}
	wg.Wait()

}

func defaultConfigFile() string {
	cwd, _ := os.Executable()
	return filepath.Join(filepath.Dir(cwd), "config.json")
}
