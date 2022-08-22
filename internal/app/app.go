package app

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/dwilkolek/moxy/internal/config"
	"github.com/dwilkolek/moxy/internal/logger"
	"github.com/elliotchance/sshtunnel"
)

func Run(file string) {
	filePath := findConfigFileLocation(file)
	resolvedConfig, err := config.NewConfig(filePath)
	if err != nil {
		panic("Failed to parse config file " + filePath + " .")
	}
	run(resolvedConfig)
}

func run(config *config.MoxyConfig) {
	var wg sync.WaitGroup

	tunnelPort := setupTunnel(config.Tunnel)

	wg.Add(1)
	for service, conf := range config.Services {
		setupHttpServerForService(service, conf, fmt.Sprintf("http://localhost:%d", tunnelPort))
	}
	wg.Wait()
}

func setupTunnel(tunnelConfig config.TunnelConfig) int {
	tunnelLogger := logger.NewOnPort("Tunnel", 0)
	privateKeyLocation := resolveFileLocation(tunnelConfig.PathToPrivateKey, tunnelLogger)
	privateKeyLocation = assertFileExists(privateKeyLocation, fmt.Sprintf("Failed to read private key from %s", privateKeyLocation))

	tunnelLogger.Printf("Using private key from file: %s", privateKeyLocation)
	tunnel := sshtunnel.NewSSHTunnel(
		tunnelConfig.UserAndHost,
		sshtunnel.PrivateKeyFile(privateKeyLocation),
		tunnelConfig.Destination,
		"0",
	)

	go tunnel.Start()
	time.Sleep(100 * time.Millisecond)
	tunnelLogger = logger.NewOnPort("Tunnel", tunnel.Local.Port)
	tunnel.Log = tunnelLogger
	tunnel.Log.Printf("Started and exposed on port: %d\n", tunnel.Local.Port)

	return tunnel.Local.Port
}

func setupHttpServerForService(service string, conf config.ProxyConfig, to string) {
	go func() {
		logger := logger.NewOnPort(service, conf.Port)
		origin, _ := url.Parse(to)

		director := func(req *http.Request) {
			for header, value := range conf.Headers {
				req.Header.Add(header, value)
				if strings.ToLower(header) == "host" {
					req.Host = value
				}
			}

			req.URL.Scheme = "http"
			req.URL.Host = origin.Host
		}

		proxy := &httputil.ReverseProxy{Director: director}
		server := http.NewServeMux()

		server.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			logger.Printf("%s %s %s\n", r.Method, r.URL.Path, r.URL.RawQuery)
			defer r.Body.Close()
			if conf.AllowCors {
				if r.Method == http.MethodOptions {
					w.Header().Set("Access-Control-Allow-Origin", "*")
					w.Header().Set("Access-Control-Allow-Methods", r.Header.Get("Access-Control-Request-Method"))
					w.Header().Set("Access-Control-Allow-Headers", r.Header.Get("Access-Control-Request-Headers"))
					w.WriteHeader(http.StatusNoContent)
				} else {
					// Set CORS headers for the main request.
					w.Header().Set("Access-Control-Allow-Origin", "*")
					proxy.ServeHTTP(w, r)
				}
			} else {
				proxy.ServeHTTP(w, r)
			}

		})

		logger.Printf("Starting server")
		if err := http.ListenAndServe(fmt.Sprintf(":%d", conf.Port), server); err != nil {
			log.Fatal(err)
		}
	}()
}

func resolveFileLocation(filePath string, logger *log.Logger) string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	executablePath := filepath.Dir(ex)
	if filepath.IsAbs(filePath) {
		return filePath
	}

	filePathToCheck := filepath.Join(executablePath, filePath)
	if fileExists(filePathToCheck) {
		return filePathToCheck
	} else {
		logger.Printf("File %s not found. Trying with current directory as base.\n", filePathToCheck)
	}

	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	filePathToCheck = filepath.Join(dir, filePath)
	if fileExists(filePathToCheck) {
		return filePathToCheck
	}
	return filePathToCheck
}

func fileExists(filepath string) bool {
	if _, err := os.Stat(filepath); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

func assertFileExists(filepath string, errorMsg string) string {
	if _, err := os.Stat(filepath); errors.Is(err, os.ErrNotExist) {
		panic(errorMsg)
	}
	return filepath
}

func findConfigFileLocation(arg string) string {
	logger := logger.New("Moxy")
	var configFile string
	if len(arg) > 0 {
		if strings.Index(arg, ".") != -1 {
			configFile = arg
		} else {
			configFile = fmt.Sprintf("config-%s.json", arg)
		}
	} else {
		panic("Config file not set!")
	}
	configFile = resolveFileLocation(configFile, logger)
	configFile = assertFileExists(configFile, fmt.Sprintf("Config file: %s doesn't exists.", configFile))
	logger.Printf("Config file: %s \n", configFile)
	return configFile
}
