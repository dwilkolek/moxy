package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/elliotchance/sshtunnel"
	"github.com/inconshreveable/go-update"
)

var wg sync.WaitGroup

var version = "1.0.1"

type TunnelConfig struct {
	UserAndHost      string `json:"userAndHost"`
	PathToPrivateKey string `json:"pathToPrivateKey"`
	Destination      string `json:"destination"`
}
type MoxyConfig struct {
	Tunnel   TunnelConfig           `json:"tunnel"`
	Services map[string]ProxyConfig `json:"services"`
}

type ProxyConfig struct {
	Port      int               `json:"port"`
	Headers   map[string]string `json:"headers"`
	AllowCors bool              `json:"allowCors"`
}

func doUpdate() error {
	var url string

	var client = &http.Client{}
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		url = req.URL.String()
		return nil
	}
	client.Get("https://github.com/dwilkolek/moxy/releases/latest")
	if runtime.GOOS == "windows" {
		url = url + "/moxy-windows.exe"
	} else if runtime.GOOS == "darwin" {
		url = url + "/moxy-mac"
	} else if runtime.GOOS == "linux" {
		url = url + "/moxy-linux"
	} else {
		return nil
	}
	url = strings.Replace(url, "/tag/", "/download/", -1)
	fmt.Printf("Checking for updates %s \n", url)

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	err = update.Apply(resp.Body, update.Options{})
	if err != nil {
		//
	}
	return err
}

func main() {
	doUpdate()
	fmt.Printf("Moxy version: %s \n", version)

	configFile := "config.json"
	if len(os.Args) > 1 {
		configFile = os.Args[1]
	}
	fmt.Printf("Config file: %s \n", configFile)

	config := getConfig(configFile)
	tunnelPort := setupTunnel(config.Tunnel)
	wg.Add(1)

	for service, conf := range config.Services {
		setupHttpServerForService(service, conf, fmt.Sprintf("http://localhost:%d", tunnelPort))
	}
	wg.Wait()
}

func setupTunnel(tunnelConfig TunnelConfig) int {
	tunnel := sshtunnel.NewSSHTunnel(
		tunnelConfig.UserAndHost,
		sshtunnel.PrivateKeyFile(tunnelConfig.PathToPrivateKey),
		tunnelConfig.Destination,
		"0",
	)

	tunnel.Log = log.New(os.Stdout, "", log.Ldate|log.Lmicroseconds)

	go tunnel.Start()
	time.Sleep(100 * time.Millisecond)
	fmt.Printf("Tunnel started and exposed on port: %d\n", tunnel.Local.Port)
	return tunnel.Local.Port
}

func setupHttpServerForService(service string, conf ProxyConfig, to string) {
	go func() {
		origin, _ := url.Parse(to)
		fmt.Printf("Setup %s \n", service)

		director := func(req *http.Request) {
			req.Header.Add("X-Forwarded-Host", req.Host)
			req.Header.Add("X-Origin-Host", origin.Host)

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
			fmt.Printf("%s: %s %s\n", service, r.Method, r.URL.Path)

			fmt.Println(conf.AllowCors)
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
			}

		})

		fmt.Printf("Starting server for %s at port %d\n", service, conf.Port)
		if err := http.ListenAndServe(fmt.Sprintf(":%d", conf.Port), server); err != nil {
			log.Fatal(err)
		}
	}()
}

func getConfig(file string) MoxyConfig {
	jsonFile, _ := ioutil.ReadFile(file)
	var objmap MoxyConfig
	json.Unmarshal([]byte(jsonFile), &objmap)
	return objmap
}
