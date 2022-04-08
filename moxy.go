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
	"sync"
	"time"

	"github.com/elliotchance/sshtunnel"
)

var wg sync.WaitGroup

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
	Port    int               `json:"port"`
	Headers map[string]string `json:"headers"`
}

func main() {
	config := getConfig(os.Args[1])
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

		hostHeader := service + ".service"

		director := func(req *http.Request) {
			req.Header.Add("X-Forwarded-Host", req.Host)
			req.Header.Add("X-Origin-Host", origin.Host)

			for header, value := range conf.Headers {
				req.Header.Add(header, value)
			}

			req.Host = hostHeader
			req.URL.Scheme = "http"
			req.URL.Host = origin.Host

		}

		proxy := &httputil.ReverseProxy{Director: director}
		server := http.NewServeMux()

		server.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			fmt.Printf("%s: %s %s\n", service, r.Method, r.URL.Path)
			proxy.ServeHTTP(w, r)
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
